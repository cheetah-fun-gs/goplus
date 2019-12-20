package delaytrigger

// WalkStop 停止遍历的判断
type WalkStop func() bool

// WalkHandle 目标处理的回调
type WalkHandle func(targetID string, eventData interface{}) error

// Walk 遍历所有有效事件
func (trigger *DelayTrigger) Walk(handle WalkHandle, stop WalkStop) error {
	return nil
}

// WalkWithEventID 根据EventID 遍历有效事件
func (trigger *DelayTrigger) WalkWithEventID(eventIDS []string, handle WalkHandle, stop WalkStop) error {
	return nil
}
