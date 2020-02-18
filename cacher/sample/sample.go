package main

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/cheetah-fun-gs/goplus/cacher"
	redigo "github.com/gomodule/redigo/redis"
)

type sample map[string]interface{}

func (s sample) Set(data interface{}, args ...interface{}) error {
	s[args[0].(string)] = data
	return nil
}

func (s sample) Del(args ...interface{}) error {
	delete(s, args[0].(string))
	return nil
}

func (s sample) Get(dest interface{}, args ...interface{}) (bool, error) {
	v, ok := s[args[0].(string)]
	if !ok {
		return false, nil
	}
	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v))
	return true, nil
}

type testData struct {
	A int    `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

func dial() (redigo.Conn, error) {
	return redigo.DialTimeout("tcp", "127.0.0.1:6379", 2*time.Second, 2*time.Second, 2*time.Second)
}

func main() {
	pool := &redigo.Pool{
		Dial: dial,
	}
	defer pool.Close()

	c := cacher.New("test", pool, &sample{})
	c.DisableGoroutine()

	// int
	test1Key := "test1"
	test1Val := int(1)
	var val1 int

	if err := c.Set(test1Val, test1Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok1, err := c.Get(&val1, test1Key); err != nil {
		panic(err)
	} else if !ok1 {
		panic(fmt.Sprintf("%v not found", test1Key))
	}
	if val1 != test1Val {
		panic(fmt.Sprintf("%v not equal %v", test1Key, test1Val))
	}

	if err := c.Del(test1Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok1, err := c.Get(&val1, test1Key); err != nil {
		panic(err)
	} else if ok1 {
		panic(fmt.Sprintf("%v not delete", test1Key))
	}

	// []byte
	test2Key := "test2"
	test2Val := []byte("1231")
	var val2 []byte

	if err := c.Set(test2Val, test2Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok2, err := c.Get(&val2, test2Key); err != nil {
		panic(err)
	} else if !ok2 {
		panic(fmt.Sprintf("%v not found", test2Key))
	}
	if !bytes.Equal(val2, test2Val) {
		panic(fmt.Sprintf("%v not equal %v", test2Key, test2Val))
	}

	if err := c.Del(test2Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok2, err := c.Get(&val2, test2Key); err != nil {
		panic(err)
	} else if ok2 {
		panic(fmt.Sprintf("%v not delete", test2Key))
	}

	// struct
	test3Key := "test3"
	test3Val := testData{
		A: 1,
		B: "1",
	}
	val3 := testData{}

	if err := c.Set(test3Val, test3Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok3, err := c.Get(&val3, test3Key); err != nil {
		panic(err)
	} else if !ok3 {
		panic(fmt.Sprintf("%v not found", test3Key))
	}
	if fmt.Sprintf("%v", test3Val) != fmt.Sprintf("%v", val3) {
		panic(fmt.Sprintf("%v not equal %v", test3Key, test3Val))
	}

	if err := c.Del(test3Key); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	if ok3, err := c.Get(&val3, test3Key); err != nil {
		panic(err)
	} else if ok3 {
		panic(fmt.Sprintf("%v not delete", test3Key))
	}
}
