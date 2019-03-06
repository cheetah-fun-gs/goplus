package redigo

import (
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// ErrorLocked 被占用
var ErrorLocked = fmt.Errorf("locked")

// Lock redis锁 秒级 只锁不解
func Lock(conn redigo.Conn, lockKey string, timeOut int) error {
	_, err := conn.Do("SET", lockKey, "", "EX", fmt.Sprintf("%d", timeOut), "NX")
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err == redigo.ErrNil {
		return ErrorLocked
	}
	return nil
}
