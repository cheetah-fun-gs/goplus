package reflect

import (
	"math/rand"
	"reflect"
	"strings"

	randplus "github.com/cheetah-fun-gs/goplus/math/rand"
)

var (
	chars = strings.Split("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_", "")
)

// Mocker ...
type Mocker struct {
	data        interface{}
	isRandom    bool
	isPointer   bool
	isRecurse   bool                        // 是否递归
	skipRecurse func(typ reflect.Type) bool // 跳过递归的判定
}

// Random ...
func (mocker *Mocker) Random() *Mocker {
	mocker.isRandom = true
	return mocker
}

// Pointer ...
func (mocker *Mocker) Pointer() *Mocker {
	mocker.isPointer = true
	return mocker
}

// DisableRecurse ...
func (mocker *Mocker) DisableRecurse() *Mocker {
	mocker.isRecurse = false
	return mocker
}

// SkipRecurse ...
func (mocker *Mocker) SkipRecurse(f func(typ reflect.Type) bool) *Mocker {
	mocker.skipRecurse = f
	return mocker
}

// Value ...
func (mocker *Mocker) Value() interface{} {
	return mockAny(reflect.ValueOf(mocker.data),
		mocker.isRandom,
		mocker.isPointer,
		mocker.isRecurse,
		mocker.skipRecurse)
}

// Mock ...
func Mock(v interface{}) *Mocker {
	return &Mocker{
		data:        v,
		isRandom:    false,
		isPointer:   false,
		isRecurse:   true,
		skipRecurse: func(typ reflect.Type) bool { return false },
	}
}

func newValue(typ reflect.Type, isRandom, isPointer bool) interface{} {
	switch typ.Kind() {
	case reflect.Bool:
		var val bool
		if isRandom && randplus.MustRandint(0, 1) == 1 {
			val = true
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int:
		var val int
		if isRandom {
			val = int(randplus.MustRandint(-2147483648, 2147483647))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int8:
		var val int8
		if isRandom {
			val = int8(randplus.MustRandint(-128, 127))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int16:
		var val int16
		if isRandom {
			val = int16(randplus.MustRandint(-32768, 32767))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int32:
		var val int32
		if isRandom {
			val = int32(randplus.MustRandint(-2147483648, 2147483647))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int64:
		var val int64
		if isRandom {
			val = int64(randplus.MustRandint(-9223372036854775808, 9223372036854775807))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint:
		var val uint
		if isRandom {
			val = uint(randplus.MustRandint(1, 4294967295))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint8:
		var val uint8
		if isRandom {
			val = uint8(randplus.MustRandint(1, 255))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint16:
		var val uint16
		if isRandom {
			val = uint16(randplus.MustRandint(1, 65535))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint32:
		var val uint32
		if isRandom {
			val = uint32(randplus.MustRandint(1, 4294967295))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint64:
		var val uint64
		if isRandom {
			val = uint64(randplus.MustRandint(1, 9223372036854775807) + randplus.MustRandint(0, 9223372036854775807))
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Float32:
		var val float32
		if isRandom {
			val = rand.Float32()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Float64:
		var val float64
		if isRandom {
			val = rand.Float64()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Complex64:
		var val complex64
		if isRandom {
			val = complex(rand.Float32(), rand.Float32())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Complex128:
		var val complex128
		if isRandom {
			val = complex(rand.Float64(), rand.Float64())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.String:
		var val string
		if isRandom {
			splits := []string{}
			for i := 5; i <= randplus.MustRandint(5, 16); i++ {
				splits = append(splits, chars[randplus.MustRandint(0, len(chars)-1)])
			}
			val = strings.Join(splits, "")
		}
		if isPointer {
			return &val
		}
		return val
	default:
		return nil
	}
}

func mockAny(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	v = DeepElemValue(v)
	switch v.Kind() {
	case reflect.Map:
		return mockMap(v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Struct:
		return mockStruct(v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Slice:
		return mockSlice(v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Array:
		return mockArray(v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Chan, reflect.Func, reflect.Interface:
		return nil
	default:
		return newValue(v.Type(), isRandom, isPointer)
	}
}

// mockStruct ..
func mockStruct(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Field(i).Type()
		fieldStructType := v.Type().Field(i)
		if fieldStructType.PkgPath == "" { // 只处理导出的字段
			key := structFieldName(fieldStructType, "json") // 使用json tag
			var val interface{}
			if isRecurse && !skipRecurse(fieldType) {
				val = mockAny(reflect.New(fieldType), isRandom, isPointer, isRecurse, skipRecurse)
			} else if isPointer {
				val = reflect.New(DeepElemType(fieldType)).Interface()
			} else {
				val = reflect.New(DeepElemType(fieldType)).Elem().Interface()
			}
			result[key] = val
		}
	}
	return result
}

// mockSlice ..
func mockSlice(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := mockArray(v, isRandom, isPointer, isRecurse, skipRecurse).([]interface{})
	if len(result) == 0 {
		var val interface{}
		typ := v.Type().Elem()
		if isRecurse && !skipRecurse(typ) {
			val = mockAny(reflect.New(typ), isRandom, isPointer, isRecurse, skipRecurse)
		} else if isPointer {
			val = reflect.New(DeepElemType(typ)).Interface()
		} else {
			val = reflect.New(DeepElemType(typ)).Elem().Interface()
		}
		result = append(result, val)
	}
	return result
}

// mockArray ..
func mockArray(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		var val interface{}
		typ := v.Index(i).Type()
		if isRecurse && !skipRecurse(typ) {
			val = mockAny(v.Index(i), isRandom, isPointer, isRecurse, skipRecurse)
		} else if isPointer {
			val = reflect.New(DeepElemType(typ)).Interface()
		} else {
			val = reflect.New(DeepElemType(typ)).Elem().Interface()
		}
		result = append(result, val)
	}
	return result
}

// mockMap ..
func mockMap(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := map[string]interface{}{}
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		var val interface{}
		typ := iter.Value().Type()
		if isRecurse && !skipRecurse(typ) {
			val = mockAny(iter.Value(), isRandom, isPointer, isRecurse, skipRecurse)
		} else if isPointer {
			val = reflect.New(DeepElemType(typ)).Interface()
		} else {
			val = reflect.New(DeepElemType(typ)).Elem().Interface()
		}
		result[key] = val
	}
	if len(result) == 0 {
		var val interface{}
		typ := v.Type().Elem()
		if isRecurse && !skipRecurse(typ) {
			val = mockAny(reflect.New(typ), isRandom, isPointer, isRecurse, skipRecurse)
		} else if isPointer {
			val = reflect.New(DeepElemType(typ)).Interface()
		} else {
			val = reflect.New(DeepElemType(typ)).Elem().Interface()
		}
		result[""] = val
	}
	return result
}
