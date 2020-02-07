package redigo

import (
	redigo "github.com/gomodule/redigo/redis"
)

// LPush LPush Deprecated: Use Do instead.
func LPush(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	args = append(args, v...)
	return redigo.Int(Do(conn, "LPUSH", args...))
}

// RPush RPush Deprecated: Use Do instead.
func RPush(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	args = append(args, v...)
	return redigo.Int(Do(conn, "RPUSH", args...))
}

// LPushX LPushX Deprecated: Use Do instead.
func LPushX(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	args = append(args, v...)
	return redigo.Int(Do(conn, "LPUSHX", args...))
}

// RPushX RPushX Deprecated: Use Do instead.
func RPushX(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	args = append(args, v...)
	return redigo.Int(Do(conn, "RPUSHX", args...))
}

// LPop LPop Deprecated: Use ToStruct instead.
func LPop(conn redigo.Conn, key string, dest interface{}) (bool, error) {
	return Result(conn.Do("LPOP", key)).ToStruct(dest)
}

// RPop RPop Deprecated: Use ToStruct instead.
func RPop(conn redigo.Conn, key string, dest interface{}) (bool, error) {
	return Result(conn.Do("RPOP", key)).ToStruct(dest)
}

// LRange LRange v []interface{}{} 的指针 Deprecated: Use ToList instead.
func LRange(conn redigo.Conn, key string, start, stop int, dest interface{}) error {
	return Result(conn.Do("LRANGE", key, start, stop)).ToList(dest)
}
