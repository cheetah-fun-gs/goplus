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
func XAdd(redigoConn redigo.Conn, key string, maxlen int, data interface{}) (string, error) {
	if maxlen == 0 {
		maxlen = 10000
	}
	args := []interface{}{key, "MAXLEN", "~", maxlen, "*"}

	for k, v := range structs.Map(data) {
		args = append(args, k, v)
	}
	return redigo.String(redigoConn.Do("XADD", args...))
}

// XAddPipeline xadd的pipeline模式
func XAddPipeline(redigoConn redigo.Conn, key string, maxlen int, datas ...interface{}) ([]string, error) {
	if maxlen == 0 {
		maxlen = 10000
	}

	for _, data := range datas {
		args := []interface{}{key, "MAXLEN", "~", maxlen, "*"}
		for k, v := range structs.Map(data) {
			args = append(args, k, v)
		}
		err := redigoConn.Send("XADD", args...)
		if err != nil {
			return nil, err
		}
	}
	err := redigoConn.Flush()
	if err != nil {
		return nil, err
	}
	r := []string{}
	for i := 0; i < len(datas); i++ {
		reply, _ := redigo.String(redigoConn.Receive())
		r = append(r, reply)
	}
	return r, nil
}

// XRead v 结构体指针, id 有2个特殊值: 0-0 从头开始读, $ 从加入时开始读
func XRead(redigoConn redigo.Conn, key, id string, v interface{}, block int64) (bool, string, error) {
	if block == 0 {
		block = 1000
	}
	if id == "" {
		id = "$"
	}

	args := []interface{}{"COUNT", 1, "BLOCK", block, "STREAMS", key, id}
	reply, err := redigoConn.Do("XREAD", args...)
	if err != nil && err != redigo.ErrNil {
		return false, "", err
	}
	if err == redigo.ErrNil {
		return false, "", nil
	}
	id, err = decodeReply(reply, v)
	if err != nil {
		return false, "", err
	}
	return false, id, err
}

// XReadGroup v 结构体指针
func XReadGroup(redigoConn redigo.Conn, key, groupName, consumerName, id string, v interface{}, block int, isAck bool) (bool, string, error) {
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
	reply, err := redigoConn.Do("XREADGROUP", args...)
	if err != nil && err != redigo.ErrNil {
		return false, "", err
	}
	if err == redigo.ErrNil {
		return false, "", nil
	}
	id, err = decodeReply(reply, v)
	if err != nil {
		return false, "", err
	}
	return true, id, err
}

// XGroupCreate 创建消费者组
func XGroupCreate(redigoConn redigo.Conn, key, groupName string) error {
	_, err := redigoConn.Do("XGROUP", "CREATE", key, groupName, "$", "MKSTREAM")
	return err
}

// XGroupDestroy 销毁消费者组
func XGroupDestroy(redigoConn redigo.Conn, key, groupName string) error {
	_, err := redigoConn.Do("XGROUP", "DESTROY", key, groupName)
	if err != nil && err != redigo.ErrNil {
		return err
	}
	return nil
}
