// Package delaytrigger 延迟触发器
package delaytrigger

import (
	"fmt"
	"time"

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

// TriggerTsRange 触发时间范围
type TriggerTsRange struct {
	Min int64
	Max int64
}

// Param 参数
type Param struct {
	IDS            []string
	Statuses       []EventStatus
	TriggerTsRange *TriggerTsRange
}

// GetEventsByParam 获取指定事件 ids statuses 为空表示 全匹配
func (trigger *DelayTrigger) GetEventsByParam(param *Param) ([]*Event, error) {
	// 格式化参数
	if param == nil {
		param = &Param{
			IDS:            []string{},
			Statuses:       []EventStatus{},
			TriggerTsRange: &TriggerTsRange{},
		}
	} else if param.TriggerTsRange == nil {
		param.TriggerTsRange = &TriggerTsRange{}
	}

	conn := trigger.pool.Get()
	defer conn.Close()

	key := trigger.getTriggerKey()
	events := []*Event{}
	if len(param.IDS) == 0 {
		if err := redigoplus.HVals(conn, key, &events); err != nil {
			return nil, err
		}
	} else {
		v := map[string]interface{}{}
		for _, i := range param.IDS {
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
		if isMathStatus(event.Status, param.Statuses) &&
			event.TriggerTs >= param.TriggerTsRange.Min &&
			(param.TriggerTsRange.Max == 0 || event.TriggerTs < param.TriggerTsRange.Max) {
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

// GetEventsAndCountsByParam 获取所有事件信息
func (trigger *DelayTrigger) GetEventsAndCountsByParam(param *Param) (int, []*EventAndCount, error) {
	events, err := trigger.GetEventsByParam(param)
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
	return trigger.GetEventsByParam(nil)
}

// GetActivedEvents 获得所有活跃事件
func (trigger *DelayTrigger) GetActivedEvents() ([]*Event, error) {
	now := time.Now()
	param := &Param{
		Statuses: []EventStatus{EventStatusActived},
		TriggerTsRange: &TriggerTsRange{
			Max: now.Unix(),
		},
	}
	return trigger.GetEventsByParam(param)
}

// GetActivedEventsByID 按事件id获得所有活跃事件
func (trigger *DelayTrigger) GetActivedEventsByID(ids []string) ([]*Event, error) {
	now := time.Now()
	param := &Param{
		IDS:      ids,
		Statuses: []EventStatus{EventStatusActived},
		TriggerTsRange: &TriggerTsRange{
			Max: now.Unix(),
		},
	}
	return trigger.GetEventsByParam(param)
}

// GetEventsAndCounts  获得所有事件 和 目标数
func (trigger *DelayTrigger) GetEventsAndCounts() (int, []*EventAndCount, error) {
	return trigger.GetEventsAndCountsByParam(nil)
}

// GetActivedEventsAndCounts  获得所有活跃事件 和 目标数
func (trigger *DelayTrigger) GetActivedEventsAndCounts() (int, []*EventAndCount, error) {
	now := time.Now()
	param := &Param{
		Statuses: []EventStatus{EventStatusActived},
		TriggerTsRange: &TriggerTsRange{
			Max: now.Unix(),
		},
	}
	return trigger.GetEventsAndCountsByParam(param)
}

// GetActivedEventsAndCountsByID  获得所有活跃事件 和 目标数
func (trigger *DelayTrigger) GetActivedEventsAndCountsByID(ids []string) (int, []*EventAndCount, error) {
	now := time.Now()
	param := &Param{
		IDS:      ids,
		Statuses: []EventStatus{EventStatusActived},
		TriggerTsRange: &TriggerTsRange{
			Max: now.Unix(),
		},
	}
	return trigger.GetEventsAndCountsByParam(param)
}
