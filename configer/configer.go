package configer

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
