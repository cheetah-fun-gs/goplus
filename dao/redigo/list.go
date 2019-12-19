package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// LPush LPush
func LPush(redigoAny interface{}, key string, v ...interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.ToJSON(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("LPUSH", args...))
}

// RPush RPush
func RPush(redigoAny interface{}, key string, v ...interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.ToJSON(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("RPUSH", args...))
}

// LPushX LPushX
func LPushX(redigoAny interface{}, key string, v ...interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.ToJSON(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("LPUSHX", args...))
}

// RPushX RPushX
func RPushX(redigoAny interface{}, key string, v ...interface{}) (int, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return 0, err
	}
	if isPool {
		defer conn.Close()
	}

	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.ToJSON(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("RPUSHX", args...))
}

// LPop LPop
func LPop(redigoAny interface{}, key string, v interface{}) (bool, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := redigo.String(conn.Do("LPOP", key))
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

// RPop RPop
func RPop(redigoAny interface{}, key string, v interface{}) (bool, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, err
	}
	if isPool {
		defer conn.Close()
	}

	data, err := redigo.String(conn.Do("RPOP", key))
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

// LRange LRange v []interface{}{} 的指针
func LRange(redigoAny interface{}, key string, start, stop int, v interface{}) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	datas, err := redigo.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, v)
}
