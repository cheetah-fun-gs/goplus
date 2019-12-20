package delaytrigger

import (
	"fmt"
	"time"

	redigoplus "github.com/cheetah-fun-gs/goplus/dao/redigo"
	redigo "github.com/gomodule/redigo/redis"
)

// EventStatus 事件状态
type EventStatus int

// 常量
const (
	EventStatusDisable  EventStatus = iota // 已关闭
	EventStatusActived                     // 已激活
	EventStatusFinished                    // 已完成
)

// Event 事件数据
type Event struct {
	ID        string      `json:"id,omitempty"`         // ID
	Status    EventStatus `json:"status,omitempty"`     // 状态
	TriggerTs int64       `json:"trigger_ts,omitempty"` // 触发时间戳
	Data      interface{} `json:"data,omitempty"`       // 事件数据
}

// EventRegister 注册自定义事件 isCover 是否覆盖
func (trigger *DelayTrigger) EventRegister(event *Event, isCover bool) error {
	if event.Status == EventStatusFinished {
		return fmt.Errorf("Status is finished")
	}

	now := time.Now()
	if event.Status == EventStatusActived && event.TriggerTs > now.Unix()+triggerTsMax {
		return fmt.Errorf("TriggerTs is too later")
	}

	conn := trigger.pool.Get()
	defer conn.Close()

	var err error
	key := trigger.getTriggerKey()
	if isCover {
		err = redigoplus.HSet(conn, key, event.ID, event)
	} else {
		_, err = redigoplus.HSetNX(conn, key, event.ID, event)
	}
	return err
}

// EventRegisterTimer 注册定时触发事件 isCover 是否覆盖
func (trigger *DelayTrigger) EventRegisterTimer(eventID string, triggerTs int64, eventData interface{}, isCover bool) error {
	now := time.Now()
	if triggerTs > now.Unix()+triggerTsMax || triggerTs < now.Unix()+triggerTsMin {
		return fmt.Errorf("TriggerTs is our range")
	}

	event := &Event{
		ID:        eventID,
		TriggerTs: triggerTs,
		Data:      eventData,
		Status:    EventStatusActived,
	}
	return trigger.EventRegister(event, isCover)
}

// EventUnregister 取消注册事件 清除改事件的所有信息
func (trigger *DelayTrigger) EventUnregister(eventID string) error {
	conn := trigger.pool.Get()
	defer conn.Close()

	key := trigger.getTriggerKey()
	eventKey := trigger.getEventKey(eventID)

	if err := conn.Send("DEL", eventKey); err != nil {
		return err
	}
	if err := conn.Send("HDEL", key, eventID); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	if _, err := conn.Receive(); err != nil {
		return err
	}
	if _, err := conn.Receive(); err != nil {
		return err
	}
	return nil
}

// EventExists 是否存在事件
func (trigger *DelayTrigger) EventExists(eventID string) (bool, error) {
	conn := trigger.pool.Get()
	defer conn.Close()

	key := trigger.getTriggerKey()
	ok, err := redigo.Int(conn.Do("HEXISTS", key, eventID))
	if err != nil {
		return false, err
	}
	return ok > 0, nil
}
