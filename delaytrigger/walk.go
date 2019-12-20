package delaytrigger

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// WalkStop 停止遍历的判断
type WalkStop func() bool

// WalkHandle 目标处理的回调
type WalkHandle func(targetID string, eventData interface{}) error

// WalkByParam 遍历
func (trigger *DelayTrigger) WalkByParam(walkID string, ids []string, statuses []EventStatus, handle WalkHandle, stop WalkStop) (err error) {
	// 获取所有事件信息
	events, err := trigger.GetEventsByParam(ids, statuses)
	if err != nil {
		return err
	}

	conn := trigger.pool.Get()
	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()

	// 轮询事件
	ctx := context.Background()
	for _, event := range events {
		for {
			// 判断是否退出
			if stop() {
				trigger.logger.Info(ctx, "delaytrigger %v walk %v stop", trigger.name, walkID)
				return nil
			}

			// 取出对象
			key := trigger.getEventKey(event.ID)
			targetID, err := redigo.String(conn.Do("SPOP", key))
			if err != nil && err != redigo.ErrNil {
				return err
			}

			// 没有新对象 事件处理完毕
			if err == redigo.ErrNil {
				trigger.logger.Info(ctx, "delaytrigger %v walk %v event %v finish", trigger.name, walkID, event.ID)
				if !event.IsKeep {
					err = trigger.EventUnregister(event.ID)
					if err != nil {
						trigger.logger.Warn(ctx, "delaytrigger %v walk %v event %v Unregister err: %v", trigger.name, walkID, event.ID, err)
					} else {
						trigger.logger.Info(ctx, "delaytrigger %v walk %v event %v Unregister success", trigger.name, walkID, event.ID)
					}
				}
				break
			}

			// 处理对象
			err = handle(targetID, event.Data)
			if err != nil {
				trigger.logger.Warn(ctx, "delaytrigger %v walk %v event %v target %v data %v handle err: %v", trigger.name, walkID, event.ID, targetID, event.Data, err)
			} else {
				trigger.logger.Info(ctx, "delaytrigger %v walk %v event %v target %v data %v handle success", trigger.name, walkID, event.ID, targetID, event.Data)
			}
		}
	}

	trigger.logger.Info(ctx, "delaytrigger %v walk %v finish", trigger.name, walkID)
	return nil
}

// 以下为快捷方式

// WalkActived 遍历活跃
func (trigger *DelayTrigger) WalkActived(walkID string, handle WalkHandle, stop WalkStop) (err error) {
	return trigger.WalkByParam(walkID, []string{}, []EventStatus{EventStatusActived}, handle, stop)
}

// WalkActivedByEventID 遍历活跃 按 eventid
func (trigger *DelayTrigger) WalkActivedByEventID(walkID string, ids []string, handle WalkHandle, stop WalkStop) (err error) {
	return trigger.WalkByParam(walkID, ids, []EventStatus{EventStatusActived}, handle, stop)
}
