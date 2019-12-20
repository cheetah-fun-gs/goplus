// Package delaytrigger 延迟触发器
package delaytrigger

import (
	"fmt"

	redigoplus "github.com/cheetah-fun-gs/goplus/dao/redigo"
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
		pool:   pool,
		name:   name,
		logger: &logger.DefaultLogger{},
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
	conn := trigger.pool.Get()
	defer conn.Close()

	key := trigger.getTriggerKey()
	rs := []*Event{}
	if err := redigoplus.HVals(conn, key, &rs); err != nil {
		return nil, err
	}
	return rs, nil
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

// GetEventWithCounts 获取所有事件信息
func (trigger *DelayTrigger) GetEventWithCounts() ([]*EventWithCount, error) {
	events, err := trigger.GetEvents()
	if err != nil {
		return nil, err
	}
	return trigger.GetCountForEvents(events)
}

// GetActivedEventWithCounts 获取已激活的所有事件信息
func (trigger *DelayTrigger) GetActivedEventWithCounts() ([]*EventWithCount, error) {
	events, err := trigger.GetActivedEvents()
	if err != nil {
		return nil, err
	}
	return trigger.GetCountForEvents(events)
}

// GetCountForEvents 获取事件数量
func (trigger *DelayTrigger) GetCountForEvents(events []*Event) ([]*EventWithCount, error) {
	conn := trigger.pool.Get()
	defer conn.Close()

	for _, event := range events {
		if err := conn.Send("SCARD", trigger.getEventKey(event.ID)); err != nil {
			return nil, err
		}
	}
	if err := conn.Flush(); err != nil {
		return nil, err
	}

	rs := []*EventWithCount{}
	for _, event := range events {
		count, err := redigo.Int(conn.Receive())
		if err != nil {
			return nil, err
		}
		rs = append(rs, &EventWithCount{
			Event: event,
			Count: count,
		})
	}
	return rs, nil
}

// GetTargetCount 获得所有目标数
func (trigger *DelayTrigger) GetTargetCount() (int, error) {
	eventWithCounts, err := trigger.GetEventWithCounts()
	if err != nil {
		return 0, err
	}
	total := 0
	for _, e := range eventWithCounts {
		total += e.Count
	}
	return total, nil
}

// GetActivedTargetCount 获得活跃的目标数
func (trigger *DelayTrigger) GetActivedTargetCount() (int, error) {
	eventWithCounts, err := trigger.GetActivedEventWithCounts()
	if err != nil {
		return 0, err
	}
	total := 0
	for _, e := range eventWithCounts {
		total += e.Count
	}
	return total, nil
}

// GetTargetCountForEventIDS 获得所有目标数
func (trigger *DelayTrigger) GetTargetCountForEventIDS(eventIDS []string) (int, error) {
	events := []*Event{}
	for _, eID := range eventIDS {
		events = append(events, &Event{
			ID: eID,
		})
	}

	eventWithCounts, err := trigger.GetCountForEvents(events)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, e := range eventWithCounts {
		total += e.Count
	}
	return total, nil
}

// GetActivedTargetCountForEventIDS 获得活跃的目标数
func (trigger *DelayTrigger) GetActivedTargetCountForEventIDS(eventIDS []string) (int, error) {
	eventWithCounts, err := trigger.GetActivedEventWithCounts()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, e := range eventWithCounts {
		for _, eID := range eventIDS {
			if eID == e.ID {
				total += e.Count
			}
		}
	}
	return total, nil
}
