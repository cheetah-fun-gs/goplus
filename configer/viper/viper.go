package viper

import (
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

// GetD GetD
func (vip *Viper) GetD(key string, def interface{}) interface{} {
	ok, v, err := vip.Get(key)
	if !ok || err != nil {
		return def
	}
	return v
}
