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

func isMatchID(id string, ids []string) bool {
	if len(ids) == 0 {
		return true
	}
	for _, i := range ids {
		if id == i {
			return true
		}
	}
	return false
}

func isMathStatus(status EventStatus, statuses []EventStatus) bool {
	if len(statuses) == 0 {
		return true
	}
	for _, s := range statuses {
		if status == s {
			return true
		}
	}
	return false
}

// GetEventsWithParam 获取指定事件 ids statuses 为空表示 全匹配
func (trigger *DelayTrigger) GetEventsWithParam(ids []string, statuses []EventStatus) ([]*Event, error) {
	conn := trigger.pool.Get()
	defer conn.Close()

	key := trigger.getTriggerKey()
	events := []*Event{}
	if len(ids) == 0 {
		if err := redigoplus.HVals(conn, key, &events); err != nil {
			return nil, err
		}
	} else {
		v := map[string]interface{}{}
		for _, i := range ids {
			v[i] = &Event{}
		}
		if err := redigoplus.HMGet(conn, key, v); err != nil {
			return nil, err
		}
		for _, vv := range v {
			events = append(events, vv.(*Event))
		}
	}

	rs := []*Event{}
	for _, event := range events {
		if isMathStatus(event.Status, statuses) {
			rs = append(rs, event)
		}
	}
	return rs, nil
}

// EventAndCount 事件信息和目标数
type EventAndCount struct {
	*Event
	Count int
}

// GetEventsAndCountsWithParam 获取所有事件信息
func (trigger *DelayTrigger) GetEventsAndCountsWithParam(ids []string, statuses []EventStatus) (int, []*EventAndCount, error) {
	events, err := trigger.GetEventsWithParam(ids, statuses)
	if err != nil {
		return 0, nil, err
	}

	newIDS := []string{}
	for _, event := range events {
		newIDS = append(newIDS, event.ID)
	}

	counts, err := trigger.getTargetCounts(newIDS)
	if err != nil {
		return 0, nil, err
	}

	total := 0
	rs := []*EventAndCount{}
	for i := 0; i < len(events); i++ {
		total += counts[i]
		rs = append(rs, &EventAndCount{
			Event: events[i],
			Count: counts[i],
		})
	}

	return total, rs, nil
}

// getTargetCounts 获取事件的目标数
func (trigger *DelayTrigger) getTargetCounts(ids []string) ([]int, error) {
	conn := trigger.pool.Get()
	defer conn.Close()

	for _, id := range ids {
		if err := conn.Send("SCARD", trigger.getEventKey(id)); err != nil {
			return nil, err
		}
	}
	if err := conn.Flush(); err != nil {
		return nil, err
	}

	rs := []int{}
	for i := 0; i < len(ids); i++ {
		count, err := redigo.Int(conn.Receive())
		if err != nil {
			return nil, err
		}
		rs = append(rs, count)
	}
	return rs, nil
}

// 以下为快捷方式

// GetEvents 获得所有事件
func (trigger *DelayTrigger) GetEvents() ([]*Event, error) {
	return trigger.GetEventsWithParam([]string{}, []EventStatus{})
}

// GetActivedEvents 获得所有活跃事件
func (trigger *DelayTrigger) GetActivedEvents() ([]*Event, error) {
	return trigger.GetEventsWithParam([]string{}, []EventStatus{EventStatusActived})
}

// GetEventsAndCounts  获得所有事件 和 目标数
func (trigger *DelayTrigger) GetEventsAndCounts() (int, []*EventAndCount, error) {
	return trigger.GetEventsAndCountsWithParam([]string{}, []EventStatus{})
}

// GetActivedEventsAndCounts  获得所有活跃事件 和 目标数
func (trigger *DelayTrigger) GetActivedEventsAndCounts() (int, []*EventAndCount, error) {
	return trigger.GetEventsAndCountsWithParam([]string{}, []EventStatus{EventStatusActived})
}
