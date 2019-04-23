package redigo

import (
	"encoding/json"
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

// Set Set
func Set(conn redigo.Conn, key string, value interface{}, timeOut int) error {
	storeValue := ""
	if value != nil {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		storeValue = string(data)
	}
	_, err := conn.Do("SET", key, storeValue, "EX", timeOut)
	return err
}

// Get Get
func Get(conn redigo.Conn, key string, value interface{}) (bool, error) {
	storeValue, err := redigo.String(conn.Do("GET", key))
	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}

	if storeValue == "" {
		return true, nil
	}

	err = json.Unmarshal([]byte(storeValue), value)
	if err != nil {
		return false, err
	}
	return true, nil
}
