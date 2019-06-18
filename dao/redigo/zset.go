package redigo

import (
	"fmt"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

// ZsetMember zset的成员
type ZsetMember struct {
	Member string
	Rank   int64
	Score  float64
}

// ZrankWithScore 查找排名附带score
func ZrankWithScore(conn redigo.Conn, key, member string, isReverse bool) (bool, *ZsetMember, error) {
	var commandName1, commandName2 string
	if !isReverse {
		commandName1 = "ZRANK"
		commandName2 = "ZRANGE"
	} else {
		commandName1 = "ZREVRANK"
		commandName2 = "ZREVRANGE"
	}

	scriptContext := fmt.Sprintf(`local v = redis.call("%s", KEYS[1], ARGV[1])
	if v == nil
	then
		return nil
	else
		local vv = redis.call("%s", KEYS[1], v, v, "WITHSCORES")
		return {v, vv[2]}
	end`, commandName1, commandName2)

	script := redigo.NewScript(1, scriptContext)
	r, err := redigo.Values(script.Do(conn, key, member))
	if err != nil && err != redigo.ErrNil {
		return false, nil, err
	}
	if err == redigo.ErrNil {
		return false, nil, nil
	}
	rank := r[0].(int64)
	scoreStr := string(r[1].([]byte))
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return false, nil, err
	}
	return true, &ZsetMember{
		Member: member,
		Rank:   rank,
		Score:  score,
	}, nil
}
