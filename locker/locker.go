package redigo

import (
	"fmt"
	"time"

	uuidplus "github.com/cheetah-fun-gs/goplus/uuid"
	redigo "github.com/gomodule/redigo/redis"
)

// ErrorLocked 错误: 已锁
var ErrorLocked = fmt.Errorf("locked")

// Lock 简单锁: 超时释放, 秒级, 无需解锁
func Lock(conn redigo.Conn, name string, expire int) error {
	ok, err := redigo.String(conn.Do("SET", name, "1", "EX", expire, "NX"))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil || ok != "OK" {
		return ErrorLocked
	}
	return nil
}

// Locker 守护锁: 需解锁, 进程退出自动解锁
type Locker struct {
	pool     *redigo.Pool
	name     string // 锁名称 唯一
	nonce    string // 随机字符串
	interval int    // 锁间隔
	ticker   *time.Ticker
}

// New 获取一个守护锁
func New(pool *redigo.Pool, name string, v ...interface{}) (*Locker, error) {
	interval := 200
	if len(v) > 0 {
		interval = v[0].(int)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	locker := &Locker{
		pool:     pool,
		name:     name,
		nonce:    uuidplus.NewV4().Base62(),
		interval: interval,
		ticker:   ticker,
	}
	err := locker.lock()
	if err != nil {
		ticker.Stop()
		return nil, err
	}

	go func() {
		for {
			<-locker.ticker.C
			err := locker.extend()
			if err != nil {
				locker.ticker.Stop()
			}
		}
	}()

	return locker, nil
}

// Close 守护锁解锁
func (locker *Locker) Close() {
	locker.ticker.Stop()
}

func (locker *Locker) lock() error {
	conn := locker.pool.Get()
	defer conn.Close()

	ok, err := redigo.String(conn.Do("SET", locker.name, locker.nonce, "PX", 2*locker.interval, "NX"))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil || ok != "OK" {
		return ErrorLocked
	}
	return nil
}

func (locker *Locker) extend() error {
	// 脚本统一返回 OK 成功; nil 失败
	scriptContext := `local v = redis.call("GET", KEYS[1])
	if (v == nil or (type(v) == 'boolean' and v == false))
	then
		return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2], "NX")
	elseif v == ARGV[1]
	then
		redis.call("PEXPIRE", KEYS[1], ARGV[2])
		return "OK"
	else
		return nil
	end`
	conn := locker.pool.Get()
	defer conn.Close()

	script := redigo.NewScript(1, scriptContext)
	ok, err := redigo.String(script.Do(conn, locker.name, locker.nonce, 2*locker.interval))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil || ok != "OK" {
		return ErrorLocked
	}
	return nil
}
