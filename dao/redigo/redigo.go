// Package redigo 基于redigo的redis方法
package redigo

import (
	"encoding/json"
	"fmt"
	"strings"
)

func toJSON(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("v is nil")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fromJSON(data string, v interface{}) error {
	if data == "" {
		return fmt.Errorf("data is blank")
	}
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return err
	}
	return nil
}

func stringsToJSON(datas []string) string {
	split := []string{}
	for _, data := range datas {
		if data == "" || data == `""` {
			continue
		}
		split = append(split, data)
	}
	return fmt.Sprintf("[%s]", strings.Join(split, ","))
}
