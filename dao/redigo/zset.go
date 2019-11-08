package redigo

import (
	"fmt"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

// Zfind 查找 member
func Zfind(conn redigo.Conn, key, v interface{}, isReverse bool) (ok bool, rank int, score float64, err error) {
	member, err := toJSON(v)
	if err != nil {
		return
	}
	var commandName1, commandName2 string
	if !isReverse {
		commandName1 = "ZRANK"
		commandName2 = "ZRANGE"
	} else {
		commandName1 = "ZREVRANK"
		commandName2 = "ZREVRANGE"
	}
	scriptContext := fmt.Sprintf(`local v = redis.call("%s", KEYS[1], ARGV[1])
	if (v == nil or (type(v) == 'boolean' and v == false))
	then
		return nil
	else
		local vv = redis.call("%s", KEYS[1], v, v, "WITHSCORES")
		return {v, vv[2]}
	end`, commandName1, commandName2)
	script := redigo.NewScript(1, scriptContext)
	r, err := redigo.Values(script.Do(conn, key, member))
	if err != nil && err != redigo.ErrNil {
		return
	}
	if err == redigo.ErrNil {
		return false, 0, 0, nil
	}

	rank = int(r[0].(int64))
	score, err = strconv.ParseFloat(string(r[1].([]byte)), 64)
	if err != nil {
		return
	}
	ok = true
	return
}
