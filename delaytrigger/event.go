package delaytrigger

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
	return nil
}

// EventRegisterTimer 注册定时触发事件 isCover 是否覆盖
func (trigger *DelayTrigger) EventRegisterTimer(eventID string, triggerTs int64, eventData interface{}, isCover bool) error {
	return nil
}

// EventUnregister 取消注册事件 清除改事件的所有信息
func (trigger *DelayTrigger) EventUnregister(eventID string) error {
	return nil
}

// EventExists 是否存在事件
func (trigger *DelayTrigger) EventExists(eventID string) (bool, error) {
	return false, nil
}
