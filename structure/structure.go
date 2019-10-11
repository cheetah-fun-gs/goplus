package structure

import "encoding/json"

// MapToStruct MapToStruct
func MapToStruct(data map[string]interface{}, v interface{}) error {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonByte, v); err != nil {
		return err
	}

	return nil
}
