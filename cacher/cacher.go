// Package cacher 缓存库 管理回源
package cacher

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	redigoplus "github.com/cheetah-fun-gs/goplus/dao/redigo"
	"github.com/cheetah-fun-gs/goplus/locker"
	mlogger "github.com/cheetah-fun-gs/goplus/multier/multilogger"
	redigo "github.com/gomodule/redigo/redis"
)

// ErrorLocked 错误: 已锁
var ErrorLocked = fmt.Errorf("locked")

// Cacher ...
type Cacher struct {
	name               string
	pool               *redigo.Pool
	getFunc            func(dest interface{}, args ...interface{}) error // 源获取 要处理空值, 不能抛异常
	setFunc            func(data interface{}, args ...interface{}) error // 源更新
	expire             int                                               // 缓存超时时间
	safety             int                                               // 回源安全时间 在缓存时间不足safety时, 开始回源
	isDisableGoroutine bool                                              // 是否禁用goroutine  faas中需要禁用
	mLogName           string
}

// New 一个新的缓存器
// v[0]: expire, 缓存超时时间, 默认10分钟
// v[1]: safety, 回源安全时间, 在缓存时间不足safety时, 开始回源, 默认30秒
func New(name string, pool *redigo.Pool,
	getFunc func(dest interface{}, args ...interface{}) error,
	setFunc func(data interface{}, args ...interface{}) error,
	v ...int) *Cacher {
	cacher := &Cacher{
		name:     name,
		pool:     pool,
		getFunc:  getFunc,
		setFunc:  setFunc,
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

func (cacher *Cacher) getKey(args ...interface{}) string {
	splits := []string{cacher.name, "cacher"}
	for _, arg := range splits {
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
func (cacher *Cacher) backToSource(dest interface{}, args ...interface{}) error {
	lock, err := cacher.getLocker(args...)
	if err != nil {
		return err
	}
	defer lock.Close()

	if err = cacher.getFunc(dest, args...); err != nil {
		return err
	}

	return cacher.cacheSet(dest, args...)
}

// Set ...
func (cacher *Cacher) Set(data interface{}, args ...interface{}) error {
	lock, err := cacher.getLocker(args...)
	if err != nil {
		return err
	}
	defer lock.Close()

	if err := cacher.setFunc(data, args...); err != nil {
		return err
	}

	return cacher.cacheSet(data, args...)
}

func (cacher *Cacher) cacheSet(data interface{}, args ...interface{}) error {
	conn := cacher.pool.Get()
	defer conn.Close()

	key := cacher.getKey(args...)
	_, err := redigoplus.Do(conn, "SET", key, data, "EX", cacher.expire)
	return err
}

// Get ...
func (cacher *Cacher) Get(dest interface{}, args ...interface{}) error {
	ok, deadline, err := cacher.cacheGet(dest, args...)
	if err != nil {
		return err
	}

	now := time.Now()

	// 从缓存中取到 并且无需提前回源
	if ok && deadline-now.Unix() > int64(cacher.safety) {
		return nil
	}

	// 从缓存中取到 提前回源
	if ok && deadline-now.Unix() <= int64(cacher.safety) {
		destCopy := reflect.New(reflect.TypeOf(dest).Elem()).Interface()
		if cacher.isDisableGoroutine { // 同步回源
			if err = cacher.backToSource(destCopy, args...); err != nil {
				mlogger.WarnN(cacher.mLogName, "safety sync cacher.backToSource, key: %v, err: %v", cacher.getKey(args...), err)
			} else {
				reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(destCopy).Elem())
			}
		} else { // 异步回源
			go func() {
				if err = cacher.backToSource(destCopy, args...); err != nil {
					mlogger.WarnN(cacher.mLogName, "safety async cacher.backToSource, key: %v, err: %v", cacher.getKey(args...), err)
				}
			}()
		}
		return nil
	}

	// 缓存中取不到 强制回源
	return cacher.backToSource(dest, args...)
}

func (cacher *Cacher) cacheGet(dest interface{}, args ...interface{}) (ok bool, deadline int64, err error) {
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

	ok, err = redigoplus.Result(conn.Receive()).Value(dest)
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
