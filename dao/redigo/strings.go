package redigo

import (
	redigo "github.com/gomodule/redigo/redis"
)

// Set Set Deprecated: Use Do instead.
func Set(conn redigo.Conn, key string, v interface{}, expire int) error {
	args := []interface{}{key, v}

	if expire != 0 {
		args = append(args, "EX", expire)
	}

	_, err := Do(conn, "SET", args...)
	return err
}

// Get Get Deprecated: Use Result().StringToJSON instead.
func Get(conn redigo.Conn, key string, dest interface{}) (bool, error) {
	return Result(conn.Do("GET", key)).StringToJSON(dest)
}
