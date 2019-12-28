package http

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ToURLPath 结构体转换成 url path
func ToURLPath(v interface{}) (string, error) {
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
		value, ok := m[k]
		if !ok {
			continue
		}
		switch value.(type) {
		case float64:
			splits = append(splits, fmt.Sprintf("%s=%d", k, int(value.(float64))))
		case string:
			splits = append(splits, fmt.Sprintf("%s=%s", k, value.(string)))
		default:
			continue
		}
	}
	return strings.Join(splits, "&"), nil
}
