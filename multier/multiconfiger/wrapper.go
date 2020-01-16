package multiconfiger

import (
	"fmt"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
)

// Get 获取默认配置文件的key值
func Get(key string) (ok bool, val interface{}, err error) {
	return GetN(d, key)
}

// GetD 带默认值获取默认配置文件的key值, 不存在或者返回异常则使用默认值
func GetD(key string, def interface{}) interface{} {
	return GetND(d, key, def)
}

// GetN 获取指定配置文件的key值
func GetN(name, key string) (ok bool, val interface{}, err error) {
	if c, ok := mutil[name]; ok {
		return c.Get(key)
	}
	err = fmt.Errorf("name not found: %v", name)
	return
}

// GetND 带默认值获取指定配置文件的key值, 不存在或者返回异常则使用默认值
func GetND(name, key string, def interface{}) interface{} {
	ok, v, err := GetN(name, key)
	if !ok || err != nil {
		return def
	}
	return v
}

// GetBool ...
func GetBool(key string) (ok bool, val bool, err error) {
	return GetBoolN(d, key)
}

// GetBoolD ...
func GetBoolD(key string, def bool) bool {
	return GetBoolND(d, key, def)
}

// GetBoolN ...
func GetBoolN(name, key string) (ok bool, val bool, err error) {
	ok, v, err := GetN(name, key)
	if !ok || err != nil {
		return ok, false, err
	}
	vv, ok := v.(bool)
	if !ok {
		return true, false, fmt.Errorf("val is not bool")
	}
	return true, vv, nil
}

// GetBoolND ...
func GetBoolND(name, key string, def bool) bool {
	ok, v, err := GetBoolN(name, key)
	if !ok || err != nil {
		return def
	}
	return v
}

// GetInt ...
func GetInt(key string) (ok bool, val int, err error) {
	return GetIntN(d, key)
}

// GetIntD ...
func GetIntD(key string, def int) int {
	return GetIntND(d, key, def)
}

// GetIntN ...
func GetIntN(name, key string) (ok bool, val int, err error) {
	ok, v, err := GetN(name, key)
	if !ok || err != nil {
		return ok, 0, err
	}
	vv, ok := v.(int)
	if !ok {
		return true, 0, fmt.Errorf("val is not int")
	}
	return true, vv, nil
}

// GetIntND ...
func GetIntND(name, key string, def int) int {
	ok, v, err := GetIntN(name, key)
	if !ok || err != nil {
		return def
	}
	return v
}

// GetString ...
func GetString(key string) (ok bool, val string, err error) {
	return GetStringN(d, key)
}

// GetStringD ...
func GetStringD(key string, def string) string {
	return GetStringND(d, key, def)
}

// GetStringN ...
func GetStringN(name, key string) (ok bool, val string, err error) {
	ok, v, err := GetN(name, key)
	if !ok || err != nil {
		return ok, "", err
	}
	vv, ok := v.(string)
	if !ok {
		return true, "", fmt.Errorf("val is not string")
	}
	return true, vv, nil
}

// GetStringND ...
func GetStringND(name, key string, def string) string {
	ok, v, err := GetStringN(name, key)
	if !ok || err != nil {
		return def
	}
	return v
}

// GetAny ...
func GetAny(key string, v interface{}) (ok bool, err error) {
	return GetAnyN(d, key, v)
}

// GetAnyD ...
func GetAnyD(key string, def interface{}) {
	GetAnyND(d, key, def)
}

// GetAnyN ...
func GetAnyN(name, key string, v interface{}) (ok bool, err error) {
	ok, data, err := GetN(name, key)
	if !ok || err != nil {
		return ok, err
	}

	if err = jsonplus.Convert(data, v); err != nil {
		return false, err
	}

	return true, nil
}

// GetAnyND ...
func GetAnyND(name, key string, def interface{}) {
	data, err := jsonplus.Dump(def)
	if err != nil {
		return
	}

	ok, err := GetAnyN(name, key, def)
	if ok && err == nil {
		return
	}

	jsonplus.Load(data, def)
}
