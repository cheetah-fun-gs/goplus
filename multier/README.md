# 同类对象管理器
快捷调用多个同类对象的方法，并提供了一些增强方法。

## 使用方式
1. ```mobj.Init(obj)``` 初始化该类管理器，该对象名为```default```，只能初始化一次。
2. ```mobj.Register(name, obj)``` 管理器添加一个对象，对象名为```name```，重复会抛异常。
3. ```mobj.Test(args)``` 通过管理器直接调用```default```对象的```Test```方法。
4. ```mobj.TestN(name, args)``` 通过管理器直接调用```name```对象的```Test```方法。

## 对象列表
- [日志器](#日志器)
- [配置器](#配置器)
- [redis连接池](#redis连接池)
- [mongo数据库](#mongo数据库)
- [sql数据库](#sql数据库)

### 日志器
```import mlogger "github.com/cheetah-fun-gs/goplus/multier/multilogger"```  
PS: 考虑到三方库需要打印机控制台日志，可不先Init直接Register，Init后会被清除。  

#### 定义
```golang
// Logger 日志器
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debugc(ctx context.Context, format string, v ...interface{})
	Infoc(ctx context.Context, format string, v ...interface{})
	Warnc(ctx context.Context, format string, v ...interface{})
	Errorc(ctx context.Context, format string, v ...interface{})
}
```

#### 实现
1. DefaultLogger ```import "github.com/cheetah-fun-gs/goplus/logger"```
2. log4goplus.Logger ```import log4goplus "github.com/cheetah-fun-gs/goplus/logger/log4go"```

#### 增强方法
无

### 配置器
```import mconfiger "github.com/cheetah-fun-gs/goplus/multier/multiconfiger"```  

#### 定义
```golang
// Configer 配置器
type Configer interface {
	Get(key string) (ok bool, val interface{}, err error)
	GetD(key string, def interface{}) interface{} // 带默认值获取, 异常或不存在则返回默认值
}
```

#### 实现
1. viperplus.Viper ```import viperplus "github.com/cheetah-fun-gs/goplus/configer/viper"```

#### 增强方法
1. GetBool/GetInt/GetString 指定获取类型
2. GetAny 传递指针, 通过json.Unmarshal方法解析进指针

### redis连接池
```import mredigopool "github.com/cheetah-fun-gs/goplus/multier/multiredigopool"```

#### 定义
redigo.Pool ```import redigo "github.com/gomodule/redigo/redis"```

#### 实现
无

#### 增强方法
无

### mongo数据库
```import mmgodb "github.com/cheetah-fun-gs/goplus/multier/multimgodb"```

#### 定义
mgo.Database ```import "github.com/globalsign/mgo"```

#### 实现
无

#### 增强方法
无

### sql数据库
```import mmgodb "github.com/cheetah-fun-gs/goplus/multier/multisqldb"```

#### 定义
sql.DB ```import "database/sql"```

#### 实现
无

#### 增强方法
1. Get/Select 类似sqlx的Get/Select，传入指针，结果直接导入指针。
2. Insert 指定表名和列对象，直接插入。
