package redigo

import (
	"encoding/json"

	redigo "github.com/gomodule/redigo/redis"
)

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
