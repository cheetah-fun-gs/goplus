package redigo

import (
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// Lock redis锁 秒级 只锁不解
func Lock(conn redigo.Conn, lockKey string, timeOut int) (bool, error) {
	value, err := redigo.String(conn.Do("SET", lockKey, "", "EX", fmt.Sprintf("%d", timeOut), "NX"))
	if err != nil {
		return false, err
	}
	return value == "OK", nil
}
