package http

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ToRawQuery 结构体转换成 Request.URL.RawQuery; 不要使用浮点数
func ToRawQuery(v interface{}) (string, error) {
	// struct to map
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}

	// 字典序排序
	sortedKeys := make([]string, 0)
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// 对key=value的键值对用&连接起来，略过空值
	splits := []string{}
	for _, k := range sortedKeys {
		v := m[k]
		switch v.(type) {
		case float64:
			splits = append(splits, fmt.Sprintf("%s=%d", k, int(v.(float64))))
		default:
			splits = append(splits, fmt.Sprintf("%s=%v", k, v))
		}
	}
	return strings.Join(splits, "&"), nil
}
