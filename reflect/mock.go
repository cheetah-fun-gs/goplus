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
		return mockMap(v)
	case reflect.Struct:
		return mockStruct(v)
	case reflect.Slice:
		return mockSlice(v)
	case reflect.Array:
		return mockArray(v)
	case reflect.Chan, reflect.Func, reflect.Interface:
		return nil
	default:
		if !v.IsValid() {
			return reflect.New(v.Type()).Interface()
		}
		return v.Interface()
	}
}

// mockStruct ..
func mockStruct(v reflect.Value) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.PkgPath == "" { // 未导出的字段忽略
			result[structFieldName(field, "json")] = mockAny(reflect.New(field.Type))
		}
	}
	if len(result) == 0 {
		return v.Interface()
	}
	return result
}

// mockSlice ..
func mockSlice(v reflect.Value) interface{} {
	result := mockArray(v).([]interface{})
	if len(result) == 0 {
		result = append(result, mockAny(DeepElemValue(reflect.New(v.Type().Elem()))))
	}
	return result
}

// mockArray ..
func mockArray(v reflect.Value) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		result = append(result, mockAny(v.Index(i)))
	}
	return result
}

// mockMap ..
func mockMap(v reflect.Value) interface{} {
	result := map[string]interface{}{}
	iter := v.MapRange()
	for iter.Next() {
		result[iter.Key().String()] = mockAny(iter.Value())
	}
	if len(result) == 0 {
		result["mock"] = mockAny(DeepElemValue(reflect.New(v.Type().Elem())))
	}
	return result
}
