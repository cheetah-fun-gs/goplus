package delaytrigger

// TargetRegister 注册目标 isCover 是否覆盖
func (trigger *DelayTrigger) TargetRegister(eventID, targetID string, targetData interface{}, isCover bool) error {
	return nil
}

// TargetUnregister 取消注册目标
func (trigger *DelayTrigger) TargetUnregister(eventID, targetID string) error {
	return nil
}

// TargetExists 是否存在目标
func (trigger *DelayTrigger) TargetExists(eventID, targetID string) error {
	return nil
}
