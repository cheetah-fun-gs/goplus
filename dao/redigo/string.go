package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// Set Set
func Set(redigoAny interface{}, key string, v interface{}, expire int) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := jsonplus.ToJSON(v)
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
func Get(redigoAny interface{}, key string, v interface{}) (bool, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := redigo.String(conn.Do("GET", key))
	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}

	if err := jsonplus.FromJSON(data, v); err != nil {
		return false, err
	}
	return true, nil
}
