package json

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToJSON ToJSON
func ToJSON(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("v is nil")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON FromJSON
func FromJSON(data string, v interface{}) error {
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
	return FromJSON(stringsToListJSON(datas), v)
}

// StringsToMap v map[string]interface{}{} 的指针
func StringsToMap(datas []string, v interface{}) error {
	if len(datas)%2 != 0 {
		return fmt.Errorf("length is not even")
	}
	return FromJSON(stringsToMapJSON(datas), v)
}
