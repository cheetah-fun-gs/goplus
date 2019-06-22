package redigo

import (
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	uuidplus "gitlab.liebaopay.com/mikezhang/goplus/uuid"
)

// ErrorLocked 被占用
var ErrorLocked = fmt.Errorf("locked")

// Lock redis锁 秒级 只锁不解
func Lock(conn redigo.Conn, name string, timeout int) error {
	_, err := conn.Do("SET", name, "", "EX", fmt.Sprintf("%d", timeout), "NX")
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	return nil
}

// Locker 锁
type Locker struct {
	pool        *redigo.Pool
	name        string
	uuid        string
	millisecond int
	ticker      *time.Ticker
}

// NewLocker 1. 需要解锁; 2. 进程退出自动解锁
func NewLocker(pool *redigo.Pool, name string, millisecond int) (*Locker, error) {
	if millisecond == 0 {
		millisecond = 200
	}

	uid := uuidplus.GenerateUUID4()
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

// Close 解锁
func (locker *Locker) Close() {
	locker.ticker.Stop()
}

func (locker *Locker) lock() error {
	conn := locker.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", locker.name, locker.uuid, "PX", locker.millisecond, "NX")
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
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
	_, err := script.Do(conn, locker.name, locker.uuid, locker.millisecond)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	return nil
}
