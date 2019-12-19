package redigo

import (
	"github.com/fatih/structs"
	redigo "github.com/gomodule/redigo/redis"
)

// 1) 1) "mystream"
//    2) 1) 1) 1526984818136-0
//          2) 1) "duration"
//             2) "1532"
//             3) "event-id"
//             4) "5"
//             5) "user-id"
//             6) "7782813"
//       2) 1) 1526999352406-0
//          2) 1) "duration"
//             2) "812"
//             3) "event-id"
//             4) "9"
//             5) "user-id"
//             6) "388234"

// 一个 key 一条 消息
func decodeReply(reply interface{}, v interface{}) (string, error) {
	streamReply := reply.([]interface{})[0] // 第一个 stream key 的返回
	// streamName := string(streamReply.([]interface{})[0].([]byte))
	streamData := streamReply.([]interface{})[1]
	msg := streamData.([]interface{})[0] // 第一条消息
	msgID := string(msg.([]interface{})[0].([]byte))
	msgBody := msg.([]interface{})
	err := redigo.ScanStruct(msgBody, v)
	if err != nil {
		return "", err
	}
	return msgID, nil
}

// XAdd xadd  data:struct
func XAdd(redigoAny interface{}, key string, maxlen int, v interface{}) (string, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return "", err
	}
	if isPool {
		defer conn.Close()
	}

	if maxlen == 0 {
		maxlen = 10000
	}
	args := []interface{}{key, "MAXLEN", "~", maxlen, "*"}

	for key, val := range structs.Map(v) {
		args = append(args, key, val)
	}
	return redigo.String(conn.Do("XADD", args...))
}

// XAddPipeline xadd的pipeline模式
func XAddPipeline(redigoAny interface{}, key string, maxlen int, v ...interface{}) ([]string, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return nil, err
	}
	if isPool {
		defer conn.Close()
	}

	if maxlen == 0 {
		maxlen = 10000
	}

	for _, vv := range v {
		args := []interface{}{key, "MAXLEN", "~", maxlen, "*"}
		for key, val := range structs.Map(vv) {
			args = append(args, key, val)
		}
		err := conn.Send("XADD", args...)
		if err != nil {
			return nil, err
		}
	}
	err = conn.Flush()
	if err != nil {
		return nil, err
	}
	r := []string{}
	for i := 0; i < len(v); i++ {
		reply, _ := redigo.String(conn.Receive())
		r = append(r, reply)
	}
	return r, nil
}

// XRead v 结构体指针, id 有2个特殊值: 0-0 从头开始读, $ 从加入时开始读
func XRead(redigoAny interface{}, key, id string, v interface{}, block int64) (bool, string, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, "", err
	}
	if isPool {
		defer conn.Close()
	}

	if block == 0 {
		block = 1000
	}
	if id == "" {
		id = "$"
	}

	args := []interface{}{"COUNT", 1, "BLOCK", block, "STREAMS", key, id}
	reply, err := conn.Do("XREAD", args...)
	if err != nil {
		return false, "", err
	}
	if reply == nil {
		return false, "", nil
	}
	id, err = decodeReply(reply, v)
	if err != nil {
		return false, "", err
	}
	return false, id, err
}

// XReadGroup v 结构体指针
func XReadGroup(redigoAny interface{}, key, groupName, consumerName, id string, v interface{}, block int, isAck bool) (bool, string, error) {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return false, "", err
	}
	if isPool {
		defer conn.Close()
	}

	if block == 0 {
		block = 1000
	}
	if consumerName == "" {
		consumerName = groupName
	}
	if id == "" {
		id = ">"
	}

	args := []interface{}{"GROUP", groupName, consumerName, "COUNT", 1, "BLOCK", block}
	if !isAck {
		args = append(args, "NOACK")
	}
	args = append(args, "STREAMS", key, ">")
	reply, err := conn.Do("XREADGROUP", args...)
	if err != nil {
		return false, "", err
	}
	if reply == nil {
		return false, "", nil
	}
	id, err = decodeReply(reply, v)
	if err != nil {
		return false, "", err
	}
	return true, id, err
}

// XGroupCreate 创建消费者组
func XGroupCreate(redigoAny interface{}, key, groupName string) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	_, err = conn.Do("XGROUP", "CREATE", key, groupName, "$", "MKSTREAM")
	return err
}

// XGroupDestroy 销毁消费者组
func XGroupDestroy(redigoAny interface{}, key, groupName string) error {
	isPool, conn, err := AssertConn(redigoAny)
	if err != nil {
		return err
	}
	if isPool {
		defer conn.Close()
	}

	_, err = redigo.Bytes(conn.Do("XGROUP", "DESTROY", key, groupName))
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return nil
}
