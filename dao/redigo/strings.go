package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// Set Set
func Set(conn redigo.Conn, key string, v interface{}, expire int) error {
	data, err := jsonplus.Dump(v)
	if err != nil {
		return err
	}

	if expire != 0 {
		_, err := conn.Do("SET", key, data, "EX", expire)
		return err
	}

	_, err = conn.Do("SET", key, data)
	return err
}

// Get Get
func Get(conn redigo.Conn, key string, v interface{}) (bool, error) {
	data, err := redigo.String(conn.Do("GET", key))
	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}

	if err := jsonplus.Load(data, v); err != nil {
		return false, err
	}
	return true, nil
}
