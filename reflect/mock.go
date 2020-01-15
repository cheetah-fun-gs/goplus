package reflect

import (
	"math/rand"
	"reflect"

	randplus "github.com/cheetah-fun-gs/goplus/math/rand"
)

// Mock 生成 随机值的 mock 数据
func Mock(v interface{}) interface{} {
	return mockAny(reflect.ValueOf(v), false)
}

// MockZero 生成 0值的 mock 数据
func MockZero(v interface{}) interface{} {
	return mockAny(reflect.ValueOf(v), true)
}

func randomNew(typ reflect.Type) interface{} {
	switch typ.Kind() {
	case reflect.Bool:
		return false
	case reflect.Int:
		return int(randplus.MustRandint(-2147483648, 2147483647))
	case reflect.Int8:
		return int8(randplus.MustRandint(-128, 127))
	case reflect.Int16:
		return int16(randplus.MustRandint(-32768, 32767))
	case reflect.Int32:
		return int32(randplus.MustRandint(-2147483648, 2147483647))
	case reflect.Int64:
		return int64(randplus.MustRandint(-9223372036854775808, 9223372036854775807))
	case reflect.Uint:
		return uint(randplus.MustRandint(1, 4294967295))
	case reflect.Uint8:
		return uint8(randplus.MustRandint(1, 255))
	case reflect.Uint16:
		return uint16(randplus.MustRandint(1, 65535))
	case reflect.Uint32:
		return uint32(randplus.MustRandint(1, 4294967295))
	case reflect.Uint64:
		return uint64(randplus.MustRandint(1, 9223372036854775807))
	case reflect.Float32:
		return float32(rand.Float32())
	case reflect.Float64:
		return float64(rand.Float64())
	case reflect.Complex64:
		return complex64(complex(rand.Float32(), rand.Float32()))
	case reflect.Complex128:
		return complex128(complex(rand.Float64(), rand.Float64()))
	case reflect.String:
		samples := []string{"abc", "test", "this is a string"}
		weights := []int{}
		for i := 0; i < len(samples); i++ {
			weights = append(weights, 1)
		}
		index := randplus.MustWeightSample(weights)
		return samples[index]
	default:
		return nil
	}
}

func mockAny(v reflect.Value, isZeroValue bool) interface{} {
	v = DeepElemValue(v)
	switch v.Kind() {
	case reflect.Map:
		return mockMap(v, isZeroValue, true)
	case reflect.Struct:
		return mockStruct(v, isZeroValue, true)
	case reflect.Slice:
		return mockSlice(v, isZeroValue, true)
	case reflect.Array:
		return mockArray(v, isZeroValue, true)
	case reflect.Chan, reflect.Func, reflect.Interface:
		return nil
	default:
		if isZeroValue {
			if !v.IsValid() {
				return reflect.New(v.Type()).Interface()
			}
			return v.Interface()
		}
		return randomNew(v.Type())
	}
}

// MockStruct ...
func MockStruct(v interface{}, isZeroValue, isRecurse bool) map[string]interface{} {
	return mockStruct(DeepElemValue(reflect.ValueOf(v)), isZeroValue, isRecurse).(map[string]interface{})
}

// mockStruct ..
func mockStruct(v reflect.Value, isZeroValue, isRecurse bool) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.PkgPath == "" { // 只处理导出的字段
			key := structFieldName(field, "json")
			var val interface{}
			if isRecurse {
				val = mockAny(reflect.New(field.Type), isZeroValue)
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
func MockSlice(v interface{}, isZeroValue, isRecurse bool) []interface{} {
	return mockSlice(DeepElemValue(reflect.ValueOf(v)), isZeroValue, isRecurse).([]interface{})
}

// mockSlice ..
func mockSlice(v reflect.Value, isZeroValue, isRecurse bool) interface{} {
	result := mockArray(v, isZeroValue, isRecurse).([]interface{})
	if len(result) == 0 {
		var val interface{}
		if isRecurse {
			val = mockAny(DeepElemValue(reflect.New(v.Type().Elem())), isZeroValue)
		} else {
			val = DeepElemValue(reflect.New(v.Type().Elem())).Interface()
		}
		result = append(result, val)
	}
	return result
}

// MockArray ...
func MockArray(v interface{}, isZeroValue, isRecurse bool) []interface{} {
	return mockArray(DeepElemValue(reflect.ValueOf(v)), isZeroValue, isRecurse).([]interface{})
}

// mockArray ..
func mockArray(v reflect.Value, isZeroValue, isRecurse bool) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		var val interface{}
		if isRecurse {
			val = mockAny(v.Index(i), isZeroValue)
		} else {
			val = v.Index(i).Interface()
		}
		result = append(result, val)
	}
	return result
}

// MockMap ...
func MockMap(v interface{}, isZeroValue, isRecurse bool) map[string]interface{} {
	return mockMap(DeepElemValue(reflect.ValueOf(v)), isZeroValue, isRecurse).(map[string]interface{})
}

// mockMap ..
func mockMap(v reflect.Value, isZeroValue, isRecurse bool) interface{} {
	result := map[string]interface{}{}
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		var val interface{}
		if isRecurse {
			val = mockAny(iter.Value(), isZeroValue)
		} else {
			val = iter.Value().Interface()
		}
		result[key] = val
	}
	if len(result) == 0 {
		var val interface{}
		if isRecurse {
			val = mockAny(DeepElemValue(reflect.New(v.Type().Elem())), isZeroValue)
		} else {
			val = DeepElemValue(reflect.New(v.Type().Elem())).Interface()
		}
		result[""] = val
	}
	return result
}
