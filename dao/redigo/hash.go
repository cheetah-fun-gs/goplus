package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet
func HSet(redigoAny interface{}, key, field string, v interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := jsonplus.ToJSON(v)
	if err != nil {
		return 0, err
	}
	return redigo.Int(conn.Do("HSET", key, field, data))
}

// HSetNX HSetNX
func HSetNX(redigoAny interface{}, key, field string, v interface{}) (bool, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := jsonplus.ToJSON(v)
	if err != nil {
		return false, err
	}
	ok, err := redigo.Int(conn.Do("HSETNX", key, field, data))
	return ok > 0, err
}

// HGet HGet
func HGet(redigoAny interface{}, key, field string, v interface{}) (bool, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := redigo.String(conn.Do("HGET", key, field))
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

// HMSet HMSet
func HMSet(redigoAny interface{}, key string, v map[string]interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for field, vv := range v {
		data, err := jsonplus.ToJSON(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, field, data)
	}
	return redigo.Int(conn.Do("HSET", args...))
}

// HMGet HMGet
func HMGet(redigoAny interface{}, key string, v map[string]interface{}) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for field := range v {
		args = append(args, field)
	}
	datas, err := redigo.Strings(conn.Do("HMGET", args...))
	if err != nil {
		return err
	}
	for i := 0; i < len(v); i++ {
		field := datas[2*i]
		val := datas[2*i+1]
		if val != "" {
			if err := jsonplus.FromJSON(datas[2*i+1], v[field]); err != nil {
				return err
			}
		}
	}
	return nil
}

// HGetAll HGetAll  v map[string]interface{}{} 的指针
func HGetAll(redigoAny interface{}, key string, v interface{}) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	datas, err := redigo.Strings(conn.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return jsonplus.StringsToMap(datas, v)
}

// HVals HVals v []interface{}{} 的指针
func HVals(redigoAny interface{}, key string, v interface{}) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	datas, err := redigo.Strings(conn.Do("HVALS", key))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, v)
}
