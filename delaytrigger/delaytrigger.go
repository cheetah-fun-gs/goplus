// Package delaytrigger 延迟触发器
package delaytrigger

import (
	"fmt"

	"github.com/cheetah-fun-gs/goplus/logger"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	triggerTsMin = 120
	triggerTsMax = 3600 * 24 * 10
)

// DelayTrigger 延迟触发器
type DelayTrigger struct {
	pool   *redigo.Pool
	name   string // 触发器名称
	logger logger.Logger
}

// New 获取一个新的触发器
func New(pool *redigo.Pool, name string) *DelayTrigger {
	return &DelayTrigger{
		pool: pool,
		name: name,
	}
}

// SetLogger 指定日志器
func (trigger *DelayTrigger) SetLogger(logger logger.Logger) {
	trigger.logger = logger
}

// eventID			表示一类事件
// targetID			目标ID
// eventData		事件数据

// 触发器事件 hset key为eventID value为 eventData + TriggerTs
func (trigger *DelayTrigger) getTriggerKey() string {
	return fmt.Sprintf("%s:delaytriggerinfo", trigger.name)
}

// 事件目标 Sets key为targetID
func (trigger *DelayTrigger) getEventKey(eventID string) string {
	return fmt.Sprintf("%s:%s:delaytriggerevent", trigger.name, eventID)
}

// GetEvents 获取所有事件信息
func (trigger *DelayTrigger) GetEvents() ([]*Event, error) {
	return nil, nil
}

// GetActivedEvents 获取已激活的所有事件信息
func (trigger *DelayTrigger) GetActivedEvents() ([]*Event, error) {
	events, err := trigger.GetEvents()
	if err != nil {
		return nil, err
	}

	rs := []*Event{}
	for _, event := range events {
		if event.Status == EventStatusActived {
			rs = append(rs, event)
		}
	}
	return rs, nil
}

// EventWithCount 事件信息和目标数
type EventWithCount struct {
	*Event
	Count int
}

// GetEventsWithCount 获取所有事件信息
func (trigger *DelayTrigger) GetEventsWithCount() ([]*Event, error) {
	return nil, nil
}

// GetActivedEventsWithCount 获取已激活的所有事件信息
func (trigger *DelayTrigger) GetActivedEventsWithCount() ([]*Event, error) {
	return nil, nil
}
