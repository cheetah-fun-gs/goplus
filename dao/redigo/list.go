package redigo

import (
	redigo "github.com/gomodule/redigo/redis"
)

// LPush LPush
func LPush(conn redigo.Conn, key string, v ...interface{}) error {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := toJSON(vv)
		if err != nil {
			return err
		}
		args = append(args, data)
	}
	_, err := conn.Do("LPUSH", args...)
	return err
}

// RPush RPush
func RPush(conn redigo.Conn, key string, v ...interface{}) error {
	args := []interface{}{key}
	for _, vv := range v {
		data, err := toJSON(vv)
		if err != nil {
			return err
		}
		args = append(args, data)
	}
	_, err := conn.Do("RPUSH", args...)
	return err
}

// LPushX LPushX
func LPushX(conn redigo.Conn, key string, v interface{}) error {
	data, err := toJSON(v)
	if err != nil {
		return err
	}
	_, err = conn.Do("LPUSHX", key, data)
	return err
}

// RPushX RPushX
func RPushX(conn redigo.Conn, key string, v interface{}) error {
	data, err := toJSON(v)
	if err != nil {
		return err
	}
	_, err = conn.Do("RPUSHX", key, data)
	return err
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

	if err := fromJSON(data, v); err != nil {
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

	if err := fromJSON(data, v); err != nil {
		return false, err
	}
	return true, nil
}

// LRange LRange
func LRange(conn redigo.Conn, key string, start, stop int, v interface{}) error {
	datas, err := redigo.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	if err := fromJSON(stringsToJSON(datas), v); err != nil {
		return err
	}
	return nil
}
