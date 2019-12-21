# 延迟触发器
## 使用场景
- 定时通知指定人群
- 接收信号再通知指定人群
## 名词解释
### 事件 Event
- ID 触发器内唯一
- Status 状态，禁用/活跃/完成
- TriggerTs 触发时间戳
- Data 数据，透传
### 目标 Target
- ID 事件内唯一
## 数据结构
### 触发器信息 redis.hset
- field EventID
- value Event
### 事件信息 redis.sets
- targets
## 如何使用
### 简介
1. 启动 Walk，Walk会遍历所有生效事件，将事件的数据和目标传递给回调函数处理
2. 注册事件 EventRegister，可通过参数配置是否有效、触发时间
3. 注册事件目标 TargetRegister
4. 【可选】再次注册事件，激活事件或者关闭事件
### 其他
1. 支持取消目标的注册 TargetUnregister
2. Walk支持多实例，无需加锁
3. 可通过 GetEvents 或 GetEventsAndCounts 查看事件信息和目标数量
