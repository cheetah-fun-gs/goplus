// Package redigo 基于redigo的redis方法
package redigo

import (
	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
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

// ToStruct ...
func (res *Res) ToStruct(dest interface{}) (bool, error) {
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

// ToList ...
func (res *Res) ToList(dest interface{}) error {
	datas, err := redigo.Strings(res.reply, res.err)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToList(datas, dest)
}

// ToMap ...
func (res *Res) ToMap(dest interface{}) error {
	datas, err := redigo.Strings(res.reply, res.err)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return jsonplus.StringsToMap(datas, dest)
}
