package redigo

import (
	"fmt"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet
func HSet(conn redigo.Conn, key, field string, v interface{}) error {
	data, err := jsonplus.Dump(v)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", key, field, data)
	return err
}

// HSetNX HSetNX
func HSetNX(conn redigo.Conn, key, field string, v interface{}) (bool, error) {
	data, err := jsonplus.Dump(v)
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
	if err := jsonplus.Load(data, v); err != nil {
		return false, err
	}
	return true, nil
}

// HMSet HMSet v map[string]interface{}{} 的指针
func HMSet(conn redigo.Conn, key string, v interface{}) (int, error) {
	vv, ok := v.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("v must be map[string]interface{}")
	}

	args := []interface{}{key}
	for field, val := range vv {
		data, err := jsonplus.Dump(val)
		if err != nil {
			return 0, err
		}
		args = append(args, field, data)
	}
	return redigo.Int(conn.Do("HSET", args...))
}

// HMGet HMGet v map[string]interface{}{} 的指针
func HMGet(conn redigo.Conn, key string, v interface{}) error {
	vv, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("v must be map[string]interface{}")
	}

	args := []interface{}{key}
	for field := range vv {
		args = append(args, field)
	}
	datas, err := redigo.Strings(conn.Do("HMGET", args...))
	if err != nil {
		return err
	}
	for i := 0; i < len(vv); i++ {
		field := datas[2*i]
		val := datas[2*i+1]
		if val != "" {
			if err := jsonplus.Load(datas[2*i+1], vv[field]); err != nil {
				return err
			}
		}
	}
	return nil
}

// HGetAll HGetAll  v map[string]interface{}{} 的指针
func HGetAll(conn redigo.Conn, key string, v interface{}) error {
	datas, err := redigo.Strings(conn.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return jsonplus.StringsToMap(datas, v)
}

// HVals HVals v []interface{}{} 的指针
func HVals(conn redigo.Conn, key string, v interface{}) error {
	datas, err := redigo.Strings(conn.Do("HVALS", key))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, v)
}
