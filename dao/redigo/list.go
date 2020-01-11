package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	redigo "github.com/gomodule/redigo/redis"
)

// LPush LPush
func LPush(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.Dump(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("LPUSH", args...))
}

// RPush RPush
func RPush(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.Dump(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("RPUSH", args...))
}

// LPushX LPushX
func LPushX(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.Dump(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("LPUSHX", args...))
}

// RPushX RPushX
func RPushX(conn redigo.Conn, key string, v ...interface{}) (int, error) {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := jsonplus.Dump(vv)
		if err != nil {
			return 0, err
		}
		args = append(args, data)
	}
	return redigo.Int(conn.Do("RPUSHX", args...))
}

// LPop LPop
func LPop(conn redigo.Conn, key string, v interface{}) (bool, error) {
	data, err := redigo.String(conn.Do("LPOP", key))
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

// RPop RPop
func RPop(conn redigo.Conn, key string, v interface{}) (bool, error) {
	data, err := redigo.String(conn.Do("RPOP", key))
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

// LRange LRange v []interface{}{} 的指针
func LRange(conn redigo.Conn, key string, start, stop int, v interface{}) error {
	datas, err := redigo.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, v)
}
