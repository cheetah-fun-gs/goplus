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
func Lock(conn redigo.Conn, name string, timeout int) error {
	ok, err := redigo.String(conn.Do("SET", name, "", "EX", fmt.Sprintf("%d", timeout), "NX"))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	if ok == "" {
		return ErrorLocked
	}
	return nil
}

// Locker 守护锁: 需解锁, 进程退出自动解锁
type Locker struct {
	pool        *redigo.Pool
	name        string
	uuid        string
	millisecond int
	ticker      *time.Ticker
}

// NewLocker 获取一个守护锁
func NewLocker(pool *redigo.Pool, name string, millisecond int) (*Locker, error) {
	if millisecond == 0 {
		millisecond = 200
	}

	uid := uuidplus.NewV4()
	ticker := time.NewTicker(time.Duration(millisecond) * time.Millisecond)

	locker := &Locker{
		pool:        pool,
		name:        name,
		uuid:        uid.Base62(),
		millisecond: 2 * millisecond,
		ticker:      ticker,
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

	ok, err := redigo.String(conn.Do("SET", locker.name, locker.uuid, "PX", locker.millisecond, "NX"))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	if ok == "" {
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
	ok, err := redigo.String(script.Do(conn, locker.name, locker.uuid, locker.millisecond))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	if ok == "" {
		return ErrorLocked
	}
	return nil
}
