package redigo

import (
	"fmt"
	"reflect"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	reflectplus "github.com/cheetah-fun-gs/goplus/reflect"
	redigo "github.com/gomodule/redigo/redis"
)

// HSet HSet Deprecated: Use Do instead.
func HSet(conn redigo.Conn, key, field string, v interface{}) error {
	_, err := Do(conn, "HSET", key, field, v)
	return err
}

// HSetNX HSetNX Deprecated: Use Do instead.
func HSetNX(conn redigo.Conn, key, field string, v interface{}) (bool, error) {
	ok, err := redigo.Int(Do(conn, "HSETNX", key, field, v))
	return ok > 0, err
}

// HGet HGet Deprecated: Use Result().StringToJSON instead.
func HGet(conn redigo.Conn, key, field string, dest interface{}) (bool, error) {
	return Result(conn.Do("HGET", key, field)).StringToJSON(dest)
}

// HMSet HMSet v map[string]***{} or struct
func HMSet(conn redigo.Conn, key string, v interface{}) (int, error) {
	data := map[string]interface{}{}

	typ := reflectplus.DeepElemType(reflect.TypeOf(v))
	switch typ.Kind() {
	case reflect.Map:
		if err := jsonplus.Convert(v, &data); err != nil {
			return 0, err
		}
	case reflect.Struct:
		data = reflectplus.MockStruct(v, true, false)
	default:
		return 0, fmt.Errorf("v must be Map or Struct")
	}

	args := []interface{}{key}
	for field, val := range data {
		args = append(args, field, val)
	}
	return redigo.Int(Do(conn, "HSET", args...))
}

// HMGet HMGet dest map[string]***{} or struct 的指针
func HMGet(conn redigo.Conn, key string, dest interface{}) error {
	data := map[string]interface{}{}

	typ := reflectplus.DeepElemType(reflect.TypeOf(dest))
	switch typ.Kind() {
	case reflect.Map:
		if err := jsonplus.Convert(dest, &data); err != nil {
			return err
		}
	case reflect.Struct:
		data = reflectplus.MockStruct(dest, true, false)
	default:
		return fmt.Errorf("dest must be Map or Struct")
	}

	args := []interface{}{key}
	for field := range data {
		args = append(args, field)
	}
	return Result(conn.Do("HMGET", args...)).StringsToMap(dest)
}

// HGetAll HGetAll  v map[string]***{} 的指针 Deprecated: Use Result().StringsToMap instead.
func HGetAll(conn redigo.Conn, key string, dest interface{}) error {
	return Result(conn.Do("HGETALL", key)).StringsToMap(dest)
}

// HVals HVals v []***{} 的指针 Deprecated: Use Result().StringsToList instead.
func HVals(conn redigo.Conn, key string, dest interface{}) error {
	return Result(conn.Do("HVALS", key)).StringsToList(dest)
}
