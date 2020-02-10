# 缓存器
按照指定格式提供回源方法，对结果进行缓存

## 如何使用
```golang
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
