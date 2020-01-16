package configer

// Configer 配置器
type Configer interface {
	Get(key string) (ok bool, val interface{}, err error)
	GetD(key string, def interface{}) interface{}
}
