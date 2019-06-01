package redigo

import (
	"encoding/json"

	redigo "github.com/gomodule/redigo/redis"
)

// LPush LPush
func LPush(conn redigo.Conn, key string, values ...interface{}) error {

	args := make([]interface{}, 0)

	args = append(args, key)
	for _, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		storeValue := string(data)
		args = append(args, storeValue)
	}

	_, err := conn.Do("LPUSH", args...)

	return err
}

// LPushX LPushX
func LPushX(conn redigo.Conn, key string, value interface{}) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	storeValue := string(data)

	_, err = conn.Do("LPUSHX", key, storeValue)

	return err
}

// RPush RPush
func RPush(conn redigo.Conn, key string, values ...interface{}) error {

	args := make([]interface{}, 0)

	args = append(args, key)
	for _, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		storeValue := string(data)
		args = append(args, storeValue)
	}

	_, err := conn.Do("RPUSH", args...)

	return err
}

// RPushX RPushX
func RPushX(conn redigo.Conn, key string, value interface{}) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	storeValue := string(data)

	_, err = conn.Do("RPUSHX", key, storeValue)

	return err
}

// LPop LPop
func LPop(conn redigo.Conn, key string, value interface{}) (bool, error) {

	storeValue, err := redigo.String(conn.Do("LPOP", key))

	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}

	if storeValue == "" {
		return true, nil
	}

	err = json.Unmarshal([]byte(storeValue), value)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RPop RPop
func RPop(conn redigo.Conn, key string, value interface{}) (bool, error) {

	storeValue, err := redigo.String(conn.Do("RPOP", key))

	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil { // 找不到 返回空
		return false, nil
	}

	if storeValue == "" {
		return true, nil
	}

	err = json.Unmarshal([]byte(storeValue), value)
	if err != nil {
		return false, err
	}
	return true, nil
}
