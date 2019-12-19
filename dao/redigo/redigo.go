// Package redigo 基于redigo的redis方法
package redigo

import (
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// AssertConn 断言redigo对象
func AssertConn(v interface{}) (bool, redigo.Conn, error) {
	switch v.(type) {
	case redigo.Conn:
		return false, v.(redigo.Conn), nil
	case *redigo.Pool:
		return true, v.(*redigo.Pool).Get(), nil
	default:
		return false, nil, fmt.Errorf("assert redigo.Conn fail")
	}
}
