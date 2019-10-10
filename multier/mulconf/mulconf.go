package mullog

import (
	"fmt"
	"sync"
)

const (
	d = "default"
)

// Configer 配置器
type Configer interface {
	Get(key string) (ok bool, val interface{}, err error)
	GetBool(key string) (ok bool, val bool, err error)
	GetInt(key string) (ok bool, val int, err error)
	GetString(key string) (ok bool, val string, err error)
	GetMap(key string) (ok bool, val map[string]interface{}, err error)
	GetStruct(key string, v interface{}) (ok bool, err error)
	GetBoolD(key string, def bool) bool
	GetIntD(key string, def int) int
	GetStringD(key string, def string) string
	GetMapD(key string, def map[string]interface{}) map[string]interface{}
}

type mutilConfiger map[string]Configer

var (
	once  sync.Once
	mutil mutilConfiger
)

// Init 初始化
func Init(defaultConfiger Configer) {
	once.Do(func() {
		mutil = mutilConfiger{
			d: defaultConfiger,
		}
	})
}

// Register 注册配置器
func Register(name string, configer Configer) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = configer
	return nil
}

// Get get key from default
func Get(key string) (ok bool, val interface{}, err error) {
	return GetN(d, key)
}

// GetBool get bool key from default
func GetBool(key string) (ok bool, val bool, err error) {
	return GetBoolN(d, key)
}

// GetInt get int from default
func GetInt(key string) (ok bool, val int, err error) {
	return GetIntN(d, key)
}

// GetString get string from default
func GetString(key string) (ok bool, val string, err error) {
	return GetStringN(d, key)
}

// GetMap get map from default
func GetMap(key string) (ok bool, val map[string]interface{}, err error) {
	return GetMapN(d, key)

}

// GetStruct get struct from default
func GetStruct(key string, v interface{}) (ok bool, err error) {
	return GetStructN(d, key, v)

}

// GetBoolD get bool with def from default
func GetBoolD(key string, def bool) bool {
	return GetBoolDN(d, key, def)

}

// GetIntD get int with def from default
func GetIntD(key string, def int) int {
	return GetIntDN(d, key, def)

}

// GetStringD get map with def from default
func GetStringD(key string, def string) string {
	return GetStringDN(d, key, def)

}

// GetMapD get map with def from default
func GetMapD(key string, def map[string]interface{}) map[string]interface{} {
	return GetMapDN(d, key, def)
}

// GetN get with name
func GetN(name, key string) (ok bool, val interface{}, err error) {
	if c, ok := mutil[name]; ok {
		return c.Get(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetBoolN get bool with name
func GetBoolN(name, key string) (ok bool, val bool, err error) {
	if c, ok := mutil[name]; ok {
		return c.GetBool(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetIntN get int with name
func GetIntN(name, key string) (ok bool, val int, err error) {
	if c, ok := mutil[name]; ok {
		return c.GetInt(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetStringN get string with name
func GetStringN(name, key string) (ok bool, val string, err error) {
	if c, ok := mutil[name]; ok {
		return c.GetString(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetMapN get map with name
func GetMapN(name, key string) (ok bool, val map[string]interface{}, err error) {
	if c, ok := mutil[name]; ok {
		return c.GetMap(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetStructN get struct with name
func GetStructN(name, key string, v interface{}) (ok bool, err error) {
	if c, ok := mutil[name]; ok {
		return c.GetStruct(key, v)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetBoolDN get bool with name
func GetBoolDN(name, key string, def bool) bool {
	if c, ok := mutil[name]; ok {
		return c.GetBoolD(d, def)
	}
	return def
}

// GetIntDN get int with name
func GetIntDN(name, key string, def int) int {
	if c, ok := mutil[name]; ok {
		return c.GetIntD(d, def)
	}
	return def
}

// GetStringDN get string with name
func GetStringDN(name, key string, def string) string {
	if c, ok := mutil[name]; ok {
		return c.GetStringD(d, def)
	}
	return def
}

// GetMapDN get map with name
func GetMapDN(name, key string, def map[string]interface{}) map[string]interface{} {
	if c, ok := mutil[name]; ok {
		return c.GetMapD(d, def)
	}
	return def
}
