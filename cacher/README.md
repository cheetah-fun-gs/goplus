# 缓存器
按照指定格式提供回源方法，对结果进行缓存

## 如何使用
```golang
import	"github.com/cheetah-fun-gs/goplus/cacher"

// Source 源对象
type Source interface {
	Get(dest interface{}, args ...interface{}) (bool, error) // 获取 必须, PS: 没有结果也是一种结果 用bool表示
	Set(data interface{}, args ...interface{}) error         // 设置
	Del(args ...interface{}) error                           // 删除
}
```

1. ```Get(dest interface{}, args ...interface{}) (bool, error)``` 回源方法，必须提供。dest是结果的指针，args是指向该结果的参数列表，bool表示是否有结果
2. ```Set(data interface{}, args ...interface{}) error``` 同时设置源和缓存，不使用可实现空函数
3. ```Del(args ...interface{}) error ``` 同时设置源和缓存，不使用可实现空函数
4. ```c := cacher.New("test", pool, &source{})``` 创建一个缓存器
5. ```c := cacher.New("test", pool, &source{}), 600``` 创建一个缓存器，缓存失效时间为600秒
6. ```c := cacher.New("test", pool, &source{}), 600, 30``` 创建一个缓存器，缓存失效时间为600秒，而且在缓存失效的前30秒内会提前回源，更加平滑
7. 使用c.Get，c.Set，c.Del替代Source的对应方法

## 示例
```golang
package main

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/cheetah-fun-gs/goplus/cacher"
	redigo "github.com/gomodule/redigo/redis"
)

type example map[string]interface{}

func (e example) Set(data interface{}, args ...interface{}) error {
	e[args[0].(string)] = data
	return nil
}

func (e example) Del(args ...interface{}) error {
	delete(e, args[0].(string))
	return nil
}

func (e example) Get(dest interface{}, args ...interface{}) (bool, error) {
	v, ok := e[args[0].(string)]
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
	c := cacher.New("test", pool, &example{})
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
```
