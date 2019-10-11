package viper

import (
	"encoding/json"

	"github.com/spf13/viper"
)

// Viper Viper plus
type Viper struct {
	*viper.Viper
}

// New 获取一个Viper plus 对象
func New(configNameNoExt string, dirs ...string) (*Viper, error) {
	vip := viper.New()
	vip.SetConfigName(configNameNoExt)
	for _, dir := range dirs {
		vip.AddConfigPath(dir)
	}
	if err := vip.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Viper{Viper: vip}, nil
}

// Get Get
func (vip *Viper) Get(key string) (ok bool, val interface{}, err error) {
	if !vip.Viper.IsSet(key) {
		return false, nil, nil
	}
	return true, vip.Viper.Get(key), nil
}

// GetBool GetBool
func (vip *Viper) GetBool(key string) (ok bool, val bool, err error) {
	if !vip.Viper.IsSet(key) {
		return false, false, nil
	}
	return true, vip.Viper.GetBool(key), nil
}

// GetInt GetInt
func (vip *Viper) GetInt(key string) (ok bool, val int, err error) {
	if !vip.Viper.IsSet(key) {
		return false, 0, nil
	}
	return true, vip.Viper.GetInt(key), nil
}

// GetString GetString
func (vip *Viper) GetString(key string) (ok bool, val string, err error) {
	if !vip.Viper.IsSet(key) {
		return false, "", nil
	}
	return true, vip.Viper.GetString(key), nil
}

// GetMap GetMap
func (vip *Viper) GetMap(key string) (ok bool, val map[string]interface{}, err error) {
	if !vip.Viper.IsSet(key) {
		return false, nil, nil
	}
	return true, vip.Viper.GetStringMap(key), nil
}

// GetStruct GetStruct
func (vip *Viper) GetStruct(key string, v interface{}) (ok bool, err error) {
	ok, data, err := vip.GetMap(key)
	if !ok {
		return false, nil
	}

	jsonByte, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(jsonByte, v); err != nil {
		return false, err
	}

	return true, nil
}

// GetBoolD GetBoolD
func (vip *Viper) GetBoolD(key string, def bool) bool {
	if !vip.Viper.IsSet(key) {
		return def
	}
	return vip.Viper.GetBool(key)
}

// GetIntD GetIntD
func (vip *Viper) GetIntD(key string, def int) int {
	if !vip.Viper.IsSet(key) {
		return def
	}
	return vip.Viper.GetInt(key)
}

// GetStringD GetStringD
func (vip *Viper) GetStringD(key string, def string) string {
	if !vip.Viper.IsSet(key) {
		return def
	}
	return vip.Viper.GetString(key)
}

// GetMapD GetMapD
func (vip *Viper) GetMapD(key string, def map[string]interface{}) map[string]interface{} {
	if !vip.Viper.IsSet(key) {
		return def
	}
	return vip.Viper.GetStringMap(key)
}
