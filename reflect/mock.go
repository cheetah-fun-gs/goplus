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

func newValue(v reflect.Value, isRandom, isPointer bool) interface{} {
	switch v.Kind() {
	case reflect.Bool:
		var val bool
		if isRandom {
			if randplus.MustRandint(0, 1) == 1 {
				val = true
			}
		} else {
			val = v.Bool()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int:
		var val int
		if isRandom {
			val = int(randplus.MustRandint(-2147483648, 2147483647))
		} else {
			val = int(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int8:
		var val int8
		if isRandom {
			val = int8(randplus.MustRandint(-128, 127))
		} else {
			val = int8(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int16:
		var val int16
		if isRandom {
			val = int16(randplus.MustRandint(-32768, 32767))
		} else {
			val = int16(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int32:
		var val int32
		if isRandom {
			val = int32(randplus.MustRandint(-2147483648, 2147483647))
		} else {
			val = int32(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Int64:
		var val int64
		if isRandom {
			val = int64(randplus.MustRandint(-9223372036854775808, 9223372036854775807))
		} else {
			val = v.Int()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint:
		var val uint
		if isRandom {
			val = uint(randplus.MustRandint(1, 4294967295))
		} else {
			val = uint(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint8:
		var val uint8
		if isRandom {
			val = uint8(randplus.MustRandint(1, 255))
		} else {
			val = uint8(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint16:
		var val uint16
		if isRandom {
			val = uint16(randplus.MustRandint(1, 65535))
		} else {
			val = uint16(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint32:
		var val uint32
		if isRandom {
			val = uint32(randplus.MustRandint(1, 4294967295))
		} else {
			val = uint32(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Uint64:
		var val uint64
		if isRandom {
			val = uint64(randplus.MustRandint(1, 9223372036854775807) + randplus.MustRandint(0, 9223372036854775807))
		} else {
			val = uint64(v.Int())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Float32:
		var val float32
		if isRandom {
			val = rand.Float32()
		} else {
			val = float32(v.Float())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Float64:
		var val float64
		if isRandom {
			val = rand.Float64()
		} else {
			val = v.Float()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Complex64:
		var val complex64
		if isRandom {
			val = complex(rand.Float32(), rand.Float32())
		} else {
			val = complex64(v.Complex())
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Complex128:
		var val complex128
		if isRandom {
			val = complex(rand.Float64(), rand.Float64())
		} else {
			val = v.Complex()
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
		} else {
			val = v.String()
		}
		if isPointer {
			return &val
		}
		return val
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
		// 此时随机值是返回初始值
		if isRandom {
			if isPointer {
				return reflect.New(v.Type()).Interface()
			}
			return reflect.New(v.Type()).Elem().Interface()
		}
		// 返回原值
		if isPointer {
			val := reflect.PtrTo(v.Type())     // 创建一个指针
			reflect.ValueOf(val).Elem().Set(v) // 为指针赋值
			return val
		}
		return v.Interface()
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
		return newValue(v, isRandom, isPointer)
	}
}

// mockStruct ..
func mockStruct(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := DeepElemValue(v.Field(i))
		fieldType := fieldValue.Type()
		fieldStructType := fieldType.Field(i)

		if fieldStructType.PkgPath == "" { // 只处理导出的字段
			key := structFieldName(fieldStructType, "json") // 使用json tag
			var val interface{}
			if isRecurse && !skipRecurse(fieldType) {
				val = mockAny(reflect.New(fieldType), isRandom, isPointer, isRecurse, skipRecurse)
			} else {
				val = newValue(fieldValue, isRandom, isPointer)
			}
			result[key] = val
		}
	}
	return result
}

// mockSlice ..
func mockSlice(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	if v.Len() > 0 {
		return mockArray(v, isRandom, isPointer, isRecurse, skipRecurse).([]interface{})
	}
	// 长度为空 填充一个
	var val interface{}
	itemValue := reflect.New(DeepElemType(v.Type().Elem()))
	if isRecurse && !skipRecurse(itemValue.Type()) {
		val = mockAny(itemValue, isRandom, isPointer, isRecurse, skipRecurse)
	} else {
		val = newValue(itemValue, isRandom, isPointer)
	}
	return []interface{}{val}
}

// mockArray ..
func mockArray(v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		var val interface{}
		itemValue := v.Index(i)
		if isRecurse && !skipRecurse(itemValue.Type()) {
			val = mockAny(itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemValue, isRandom, isPointer)
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
		itemValue := iter.Value()
		if isRecurse && !skipRecurse(itemValue.Type()) {
			val = mockAny(itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemValue, isRandom, isPointer)
		}
		result[key] = val
	}
	// 没有key 填充一个
	if len(result) == 0 {
		var val interface{}
		itemValue := reflect.New(DeepElemType(v.Type().Elem()))
		if isRecurse && !skipRecurse(itemValue.Type()) {
			val = mockAny(itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemValue, isRandom, isPointer)
		}
		result[""] = val
	}
	return result
}
