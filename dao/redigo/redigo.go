// Package redigo 基于redigo的redis方法
package redigo

import (
	"reflect"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	reflectplus "github.com/cheetah-fun-gs/goplus/reflect"
	redigo "github.com/gomodule/redigo/redis"
)

func converArgs(args ...interface{}) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg.(type) {
		case string, []byte, int, int64, float64, bool, nil, redigo.Argument:
		default:
			data, err := jsonplus.Dump(arg)
			if err != nil {
				return err
			}
			args[i] = data
		}
	}
	return nil
}

// Do 参数默认使用json格式
func Do(conn redigo.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	if err := converArgs(args...); err != nil {
		return nil, err
	}
	return conn.Do(commandName, args...)
}

// Send 参数默认使用json格式
func Send(conn redigo.Conn, commandName string, args ...interface{}) error {
	if err := converArgs(args...); err != nil {
		return err
	}
	return conn.Send(commandName, args...)
}

// Result ...
func Result(reply interface{}, err error) *Res {
	return &Res{
		reply: reply,
		err:   err,
	}
}

// Res ...
type Res struct {
	reply interface{}
	err   error
}

// StringToJSON ...
func (res *Res) StringToJSON(dest interface{}) (bool, error) {
	data, err := redigo.String(res.reply, res.err)
	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil {
		return false, nil
	}
	if err = jsonplus.Load(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

// Value ...
func (res *Res) Value(dest interface{}) (bool, error) {
	var err error
	var v reflect.Value
	var jsonData string

	typ := reflectplus.DeepElemType(reflect.TypeOf(dest))
	switch typ.Kind() {
	case reflect.String:
		var data string
		data, err = redigo.String(res.reply, res.err)
		v = reflect.ValueOf(data)
	case reflect.Uint8:
		var data []byte
		data, err = redigo.Bytes(res.reply, res.err)
		v = reflect.ValueOf(data)
	case reflect.Int:
		var data int
		data, err = redigo.Int(res.reply, res.err)
		v = reflect.ValueOf(data)
	case reflect.Int64:
		var data int64
		data, err = redigo.Int64(res.reply, res.err)
		v = reflect.ValueOf(data)
	case reflect.Float64:
		var data float64
		data, err = redigo.Float64(res.reply, res.err)
		v = reflect.ValueOf(data)
	case reflect.Bool:
		var data bool
		data, err = redigo.Bool(res.reply, res.err)
		v = reflect.ValueOf(data)
	default:
		jsonData, err = redigo.String(res.reply, res.err)
	}

	if err != nil && err != redigo.ErrNil {
		return false, err
	}
	if err == redigo.ErrNil {
		return false, nil
	}

	if jsonData != "" {
		if err = jsonplus.Load(jsonData, dest); err != nil {
			return false, err
		}
	} else {
		reflect.ValueOf(dest).Elem().Set(v)
	}
	return true, nil
}

// StringsToList ...
func (res *Res) StringsToList(dest interface{}) error {
	datas, err := redigo.Strings(res.reply, res.err)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, dest)
}

// StringsToMap ...
func (res *Res) StringsToMap(dest interface{}) error {
	datas, err := redigo.Strings(res.reply, res.err)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToMap(datas, dest)
}
