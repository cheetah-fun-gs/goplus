package redigo

import (
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet Deprecated: Use Do instead.
func HSet(conn redigo.Conn, key, field string, v interface{}) error {
	_, err := Do(conn, "HSET", key, field, v)
	return err
}

// HSetNX HSetNX Deprecated: Use Do instead.
func HSetNX(conn redigo.Conn, key, field string, v interface{}) (bool, error) {
	ok, err := redigo.Int(Do(conn, "HSETNX", key, field, v))
	return ok > 0, err
}

// HGet HGet Deprecated: Use ToStruct instead.
func HGet(conn redigo.Conn, key, field string, dest interface{}) (bool, error) {
	return Result(conn.Do("HGET", key, field)).ToStruct(dest)
}

// HMSet HMSet v map[string]interface{}{} 的指针
func HMSet(conn redigo.Conn, key string, v interface{}) (int, error) {
	vv, ok := v.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("v must be map[string]interface{}")
	}

	args := []interface{}{key}
	for field, val := range vv {
		args = append(args, field, val)
	}
	return redigo.Int(Do(conn, "HSET", args...))
}

// HMGet HMGet v map[string]interface{}{} 的指针
func HMGet(conn redigo.Conn, key string, dest interface{}) error {
	v, ok := dest.(map[string]interface{})
	if !ok {
		return fmt.Errorf("dest must be map[string]interface{}")
	}

	args := []interface{}{key}
	for field := range v {
		args = append(args, field)
	}
	return Result(conn.Do("HMGET", args...)).ToMap(dest)
}

// HGetAll HGetAll  v map[string]***{} 的指针 Deprecated: Use ToMap instead.
func HGetAll(conn redigo.Conn, key string, dest interface{}) error {
	return Result(conn.Do("HGETALL", key)).ToMap(dest)
}

// HVals HVals v []***{} 的指针 Deprecated: Use ToList instead.
func HVals(conn redigo.Conn, key string, dest interface{}) error {
	return Result(conn.Do("HVALS", key)).ToList(dest)
}
