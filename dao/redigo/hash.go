package redigo

import (
	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet
func HSet(conn redigo.Conn, key, field string, v interface{}) error {
	data, err := toJSON(v)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", key, field, data)
	return err
}

// HSetNX HSetNX
func HSetNX(conn redigo.Conn, key, field string, v interface{}) (bool, error) {
	data, err := toJSON(v)
	if err != nil {
		return false, err
	}
	ok, err := redigo.Int(conn.Do("HSETNX", key, field, data))
	return ok > 0, err
}

// HGet HGet
func HGet(conn redigo.Conn, key, field string, v interface{}) (bool, error) {
	data, err := redigo.String(conn.Do("HGET", key, field))
	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}
	if err := fromJSON(data, v); err != nil {
		return false, err
	}
	return true, nil
}
