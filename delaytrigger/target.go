package delaytrigger

import redigo "github.com/gomodule/redigo/redis"

// TargetRegister 注册目标
func (trigger *DelayTrigger) TargetRegister(eventID string, targetIDS ...string) error {
	conn := trigger.pool.Get()
	defer conn.Close()

	eventKey := trigger.getEventKey(eventID)

	args := []interface{}{eventKey}
	for _, targetID := range targetIDS {
		args = append(args, targetID)
	}

	_, err := conn.Do("SADD", args...)
	return err
}

// TargetUnregister 取消注册目标
func (trigger *DelayTrigger) TargetUnregister(eventID string, targetIDS ...string) error {
	conn := trigger.pool.Get()
	defer conn.Close()

	eventKey := trigger.getEventKey(eventID)

	args := []interface{}{eventKey}
	for _, targetID := range targetIDS {
		args = append(args, targetID)
	}

	_, err := conn.Do("SREM", args...)
	return err
}

// TargetExists 是否存在目标
func (trigger *DelayTrigger) TargetExists(eventID, targetID string) (bool, error) {
	conn := trigger.pool.Get()
	defer conn.Close()

	eventKey := trigger.getEventKey(eventID)
	ok, err := redigo.Int(conn.Do("SISMEMBER", eventKey, targetID))
	if err != nil {
		return false, err
	}
	return ok > 0, nil
}
