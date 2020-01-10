package url

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func toString(v interface{}) string {
	switch v.(type) {
	case float64:
		return strconv.Itoa(int(v.(float64)))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToValues 转换成 Values; 不要使用浮点数
func ToValues(v interface{}) (url.Values, error) {
	m := map[string]interface{}{}
	// struct to map
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	for key, val := range m {
		values[key] = []string{}
		switch val.(type) {
		case []interface{}:
			for _, vv := range val.([]interface{}) {
				values[key] = append(values[key], toString(vv))
			}
		default:
			values[key] = append(values[key], toString(val))
		}
	}
	return values, nil
}

// ToRawQueryEncoded 转换成 URL.RawQuery的urlencode格式; 不要使用浮点数
func ToRawQueryEncoded(v interface{}) (string, error) {
	values, err := ToValues(v)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

// ToRawQuery 转换 URL.RawQuery; 不要使用浮点数
func ToRawQuery(v interface{}) (string, error) {
	values, err := ToValues(v)
	if err != nil {
		return "", err
	}
	// 改 values.Encode()
	if values == nil {
		return "", nil
	}
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		vs := values[key]
		for _, vv := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(vv)
		}
	}
	return buf.String(), nil
}
