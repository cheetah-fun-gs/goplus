package redigo

import (
	"encoding/json"

	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet
func HSet(conn redigo.Conn, key, field string, value interface{}) error {

	storeValue := ""
	if value != nil {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		storeValue = string(data)
	}

	_, err := conn.Do("HSET", key, field, storeValue)
	return err
}

// HSetNX HSetNX
func HSetNX(conn redigo.Conn, key, field string, value interface{}) (bool, error) {

	storeValue := ""
	if value != nil {
		data, err := json.Marshal(value)
		if err != nil {
			return false, err
		}
		storeValue = string(data)
	}

	ok, err := redigo.Int(conn.Do("HSETNX", key, field, storeValue))
	return ok > 0, err
}

// HGet HGet
func HGet(conn redigo.Conn, key, field string, value interface{}) (bool, error) {
	storeValue, err := redigo.String(conn.Do("HGET", key, field))
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
