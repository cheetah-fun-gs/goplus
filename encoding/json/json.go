package json

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 尝试处理json异常
func tryHandleError(v interface{}, err error) (interface{}, bool) {
	if strings.HasPrefix(err.Error(), "json: unsupported type: map[interface {}]") {
		return keyToString(v), true
	}
	return v, false
}

// 强制把json的key转为string
func keyToString(v interface{}) interface{} {
	switch v.(type) {
	case []interface{}:
		r := []interface{}{}
		for _, vv := range v.([]interface{}) {
			r = append(r, keyToString(vv))
		}
		return r
	case map[interface{}]interface{}:
		r := map[string]interface{}{}
		for key, val := range v.(map[interface{}]interface{}) {
			r[fmt.Sprintf("%v", key)] = keyToString(val)
		}
		return r
	}
	return v
}

// Dump ...
func Dump(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("v is nil")
	}

	var err error
	var ok bool
	for {
		data, errNew := json.Marshal(v)
		if errNew == nil {
			return string(data), nil
		}
		// 和上一次异常相同表示这类异常 tryHandleError 处理不了
		if err != nil && err.Error() == errNew.Error() {
			return "", errNew
		}

		v, ok = tryHandleError(v, errNew)
		if !ok {
			return "", errNew
		}
		err = errNew
	}
}

// Load ...
func Load(data string, v interface{}) error {
	if data == "" {
		return fmt.Errorf("data is blank")
	}
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return err
	}
	return nil
}

func stringsToListJSON(datas []string) string {
	split := []string{}
	for _, data := range datas {
		if data == "" || data == `""` {
			continue
		}
		split = append(split, data)
	}
	return fmt.Sprintf("[%s]", strings.Join(split, ","))
}

func stringsToMapJSON(datas []string) string {
	split := []string{}
	for i := 0; i < len(datas)/2; i++ {
		field := datas[2*i]
		val := datas[2*i+1]
		if val != "" {
			split = append(split, fmt.Sprintf(`"%s":%s`, field, val))
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(split, ","))
}

// StringsToList v []interface{}{} 的指针
func StringsToList(datas []string, v interface{}) error {
	return Load(stringsToListJSON(datas), v)
}

// StringsToMap v map[string]interface{}{} 的指针
func StringsToMap(datas []string, v interface{}) error {
	if len(datas)%2 != 0 {
		return fmt.Errorf("length is not even")
	}
	return Load(stringsToMapJSON(datas), v)
}

// Convert 利用json转换结构
func Convert(src, dst interface{}) error {
	if v, err := Dump(src); err != nil {
		return err
	} else if err = Load(v, dst); err != nil {
		return err
	} else {
		return nil
	}
}
