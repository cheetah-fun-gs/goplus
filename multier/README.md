# 同类对象管理器
快捷调用多个同类对象的方法，并提供了一些增强方法。

## 对象列表
1. 日志器
2. 配置器
3. redis连接池 redigo.Pool ```import redigo "github.com/gomodule/redigo/redis"```
4. mongo数据库 mgo.Database ```import "github.com/globalsign/mgo"```
5. sql数据库 sql.DB ```import "database/sql"```

## 使用方式
1. ```mobj.Init(obj)``` 初始化该类管理器，该对象名为```default```，只能初始化一次。
2. ```mobj.Register(name, obj)``` 管理器添加一个对象，对象名为```name```，重复会抛异常。
3. ```mobj.Test(args)``` 通过管理器直接调用```default```对象的```Test```方法。
4. ```mobj.TestN(name, args)``` 通过管理器直接调用```name```对象的```Test```方法。
