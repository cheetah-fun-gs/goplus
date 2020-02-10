// Package cacher 缓存库 管理回源
package cacher

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	redigoplus "github.com/cheetah-fun-gs/goplus/dao/redigo"
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	"github.com/cheetah-fun-gs/goplus/locker"
	mlogger "github.com/cheetah-fun-gs/goplus/multier/multilogger"
	redigo "github.com/gomodule/redigo/redis"
)

// ErrorLocked 错误: 已锁
var ErrorLocked = fmt.Errorf("locked")

// Source 源方法
type Source interface {
	Get(dest interface{}, args ...interface{}) (bool, error) // 获取 必须, PS: 没有结果也是一种结果 用bool表示
	Set(data interface{}, args ...interface{}) error         // 设置
	Del(args ...interface{}) error                           // 删除
}

type cacheValue struct {
	IsNil bool   `json:"is_nil,omitempty"`
	Data  string `json:"data,omitempty"`
}

func (val *cacheValue) parse(dest interface{}) (bool, error) {
	if val.IsNil {
		return false, nil
	}
	if err := jsonplus.Load(val.Data, dest); err != nil {
		return false, err
	}
	return true, nil
}

// Cacher ...
type Cacher struct {
	name               string
	pool               *redigo.Pool
	source             Source
	expire             int  // 缓存超时时间
	safety             int  // 回源安全时间 在缓存时间不足safety时, 开始回源
	isDisableGoroutine bool // 是否禁用goroutine  faas中需要禁用
	mLogName           string
}

// New 一个新的缓存器
// v[0]: expire, 缓存超时时间, 默认10分钟
// v[1]: safety, 回源安全时间, 在缓存时间不足safety时, 开始回源, 默认30秒
func New(name string, pool *redigo.Pool, source Source, v ...int) *Cacher {
	cacher := &Cacher{
		name:     name,
		pool:     pool,
		source:   source,
		expire:   600,
		safety:   30,
		mLogName: "default",
	}
	if len(v) >= 1 && v[0] != 0 {
		cacher.expire = v[0]
	}
	if len(v) == 2 && v[1] != 0 {
		cacher.safety = v[0]
	}
	if cacher.expire <= cacher.safety {
		panic("expire is below safety")
	}
	return cacher
}

// SetMLogName 设置日志器名称
func (cacher *Cacher) SetMLogName(name string) {
	cacher.mLogName = name
}

// DisableGoroutine 禁用协程 比如faas无法使用协程
func (cacher *Cacher) DisableGoroutine() {
	cacher.isDisableGoroutine = true
}

func (cacher *Cacher) getKey(args ...interface{}) string {
	splits := []string{cacher.name, "cacher"}
	for _, arg := range args {
		splits = append(splits, fmt.Sprintf("%v", arg))
	}
	return strings.Join(splits, ":")
}

func (cacher *Cacher) getLocker(args ...interface{}) (*locker.Locker, error) {
	lockerKey := cacher.getKey(args...) + ":locker"
	lock, err := locker.New(cacher.pool, lockerKey)
	if err != nil && err != locker.ErrorLocked {
		return nil, err
	}
	if err == locker.ErrorLocked {
		return nil, ErrorLocked
	}
	return lock, nil
}

// 回源
func (cacher *Cacher) backToSource(dest interface{}, args ...interface{}) (bool, error) {
	lock, err := cacher.getLocker(args...)
	if err != nil {
		return false, err
	}
	defer lock.Close()

	ok, err := cacher.source.Get(dest, args...)
	if err != nil {
		return false, err
	}

	if err = cacher.cacheSet(ok, dest, args...); err != nil {
		return false, err
	}
	return ok, nil
}

// Set ...
func (cacher *Cacher) Set(data interface{}, args ...interface{}) error {
	lock, err := cacher.getLocker(args...)
	if err != nil {
		return err
	}
	defer lock.Close()

	if err := cacher.source.Set(data, args...); err != nil {
		return err
	}

	return cacher.cacheSet(true, data, args...)
}

func (cacher *Cacher) cacheSet(vaild bool, data interface{}, args ...interface{}) error {
	val := cacheValue{
		IsNil: !vaild,
	}

	if data != nil {
		cacheData, err := jsonplus.Dump(data)
		if err != nil {
			return err
		}
		val.Data = cacheData
	}

	conn := cacher.pool.Get()
	defer conn.Close()

	key := cacher.getKey(args...)

	_, err := redigoplus.Do(conn, "SET", key, val, "EX", cacher.expire)
	return err
}

// Del ...
func (cacher *Cacher) Del(args ...interface{}) error {
	lock, err := cacher.getLocker(args...)
	if err != nil {
		return err
	}
	defer lock.Close()

	if err := cacher.source.Del(args...); err != nil {
		return err
	}

	return cacher.cacheSet(false, nil, args...)
}

// Get ...
func (cacher *Cacher) Get(dest interface{}, args ...interface{}) (bool, error) {
	val := &cacheValue{}
	ok, deadline, err := cacher.cacheGet(val, args...)
	if err != nil {
		return false, err
	}

	now := time.Now()

	// 从缓存中取到 并且无需提前回源
	if ok && deadline-now.Unix() > int64(cacher.safety) {
		return val.parse(dest)
	}

	// 从缓存中取到 提前回源
	if ok && deadline-now.Unix() <= int64(cacher.safety) {
		if cacher.isDisableGoroutine { // 同步回源
			vaild, err := cacher.backToSource(dest, args...)
			if err != nil {
				mlogger.WarnN(cacher.mLogName, "safety sync cacher.backToSource, key: %v, err: %v", cacher.getKey(args...), err)
				return val.parse(dest) // 回源失败 使用缓存
			}
			return vaild, nil // 回源成功 使用源
		}

		// 异步回源
		go func() {
			destCopy := reflect.New(reflect.TypeOf(dest).Elem()).Interface() // 拷贝一个指针
			if _, err = cacher.backToSource(destCopy, args...); err != nil {
				mlogger.WarnN(cacher.mLogName, "safety async cacher.backToSource, key: %v, err: %v", cacher.getKey(args...), err)
			}
		}()
		return val.parse(dest) // 使用缓存
	}

	// 缓存中取不到 强制回源
	return cacher.backToSource(dest, args...)
}

func (cacher *Cacher) cacheGet(val *cacheValue, args ...interface{}) (ok bool, deadline int64, err error) {
	conn := cacher.pool.Get()
	defer conn.Close()

	key := cacher.getKey(args...)
	if err = redigoplus.Send(conn, "GET", key); err != nil {
		return
	}
	if err = redigoplus.Send(conn, "TTL", key); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}

	ok, err = redigoplus.Result(conn.Receive()).Value(val)
	if err != nil {
		return false, 0, err
	}
	if !ok {
		return
	}

	var expire int64
	expire, err = redigo.Int64(conn.Receive())
	if err != nil {
		return
	}

	deadline = time.Now().Unix() + expire
	return
}
