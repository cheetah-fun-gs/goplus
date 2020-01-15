package reflect

import (
	"reflect"
)

// Mock 生成 0值的 mock 数据
func Mock(v interface{}) interface{} {
	return mockAny(reflect.ValueOf(v))
}

func mockAny(v reflect.Value) interface{} {
	v = DeepElemValue(v)
	switch v.Kind() {
	case reflect.Map:
		return mockMap(v, true)
	case reflect.Struct:
		return mockStruct(v, true)
	case reflect.Slice:
		return mockSlice(v, true)
	case reflect.Array:
		return mockArray(v, true)
	case reflect.Chan, reflect.Func, reflect.Interface:
		return nil
	default:
		if !v.IsValid() {
			return reflect.New(v.Type()).Interface()
		}
		return v.Interface()
	}
}

// MockStruct ...
func MockStruct(v interface{}, isRecurse bool) map[string]interface{} {
	return mockStruct(reflect.ValueOf(v), isRecurse).(map[string]interface{})
}

// mockStruct ..
func mockStruct(v reflect.Value, isRecurse bool) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.PkgPath == "" { // 只处理导出的字段
			key := structFieldName(field, "json")
			var val interface{}
			if isRecurse {
				val = mockAny(reflect.New(field.Type))
			} else {
				val = reflect.New(field.Type).Interface()
			}
			result[key] = val
		}
	}
	if len(result) == 0 {
		return v.Interface()
	}
	return result
}

// MockSlice ...
func MockSlice(v interface{}, isRecurse bool) []interface{} {
	return mockSlice(reflect.ValueOf(v), isRecurse).([]interface{})
}

// mockSlice ..
func mockSlice(v reflect.Value, isRecurse bool) interface{} {
	result := mockArray(v, isRecurse).([]interface{})
	if len(result) == 0 {
		var val interface{}
		if isRecurse {
			val = mockAny(DeepElemValue(reflect.New(v.Type().Elem())))
		} else {
			val = DeepElemValue(reflect.New(v.Type().Elem())).Interface()
		}
		result = append(result, val)
	}
	return result
}

// MockArray ...
func MockArray(v interface{}, isRecurse bool) []interface{} {
	return mockArray(reflect.ValueOf(v), isRecurse).([]interface{})
}

// mockArray ..
func mockArray(v reflect.Value, isRecurse bool) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		var val interface{}
		if isRecurse {
			val = mockAny(v.Index(i))
		} else {
			val = v.Index(i).Interface()
		}
		result = append(result, val)
	}
	return result
}

// MockMap ...
func MockMap(v interface{}, isRecurse bool) map[string]interface{} {
	return mockMap(reflect.ValueOf(v), isRecurse).(map[string]interface{})
}

// mockMap ..
func mockMap(v reflect.Value, isRecurse bool) interface{} {
	result := map[string]interface{}{}
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		var val interface{}
		if isRecurse {
			val = mockAny(iter.Value())
		} else {
			val = iter.Value().Interface()
		}
		result[key] = val
	}
	if len(result) == 0 {
		var val interface{}
		if isRecurse {
			val = mockAny(DeepElemValue(reflect.New(v.Type().Elem())))
		} else {
			val = DeepElemValue(reflect.New(v.Type().Elem())).Interface()
		}
		result[""] = val
	}
	return result
}
