package redigo

import (
	"bytes"
	"encoding/json"
	"fmt"

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

// LRange LRange
// Usage:
//     value := make([]*A, 0)
//     LRange(conn, key, start, stop, &value)
func LRange(conn redigo.Conn, key string, start, stop int, value interface{}) error {

	storeVals, err := redigo.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil && err != redigo.ErrNil {
		return err
	}

	valStr, err := strListToStr(storeVals)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(valStr), value)
	if err != nil {
		return err
	}

	return nil
}

func strListToStr(strList []string) (string, error) {
	var buffer bytes.Buffer

	for _, v := range strList {
		if v == `""` || v == "" {
			continue
		}
		buffer.WriteString(v)

		buffer.WriteString(",")

	}

	valStr := buffer.String()
	// 去掉最后的","
	if len(valStr) > 0 {
		valStr = valStr[:len(valStr)-1]
	}

	return fmt.Sprintf("[%s]", valStr), nil
}
