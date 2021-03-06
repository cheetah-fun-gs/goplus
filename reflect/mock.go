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
	return mockAny(
		reflect.TypeOf(mocker.data),
		reflect.ValueOf(mocker.data),
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

func newValue(typ reflect.Type, v reflect.Value, isRandom, isPointer bool) interface{} {
	switch v.Kind() {
	case reflect.Bool:
		var val bool
		if isRandom {
			if rand.Int31()%2 == 0 {
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
			val = rand.Int()
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
			val = int8(randplus.MustRandint(1, 127))
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
			val = int16(randplus.MustRandint(1, 32767))
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
			val = rand.Int31()
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
			val = rand.Int63()
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
			val = uint(rand.Uint64())
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
			val = rand.Uint32()
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
			val = rand.Uint64()
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
				return reflect.New(typ).Interface()
			}
			return reflect.New(typ).Elem().Interface()
		}

		// 返回原值
		if isPointer {
			dest := reflect.New(typ) // 创建一个指针
			dest.Elem().Set(v)       // 为指针赋值
			return dest.Interface()
		}
		return v.Interface()
	default:
		return nil
	}
}

func mockAny(typ reflect.Type, v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	typ = DeepElemType(typ)
	v = DeepElemValue(v)
	switch v.Kind() {
	case reflect.Map:
		return mockMap(typ, v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Struct:
		return mockStruct(typ, v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Slice:
		return mockSlice(typ, v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Array:
		return mockArray(typ, v, isRandom, isPointer, isRecurse, skipRecurse)
	case reflect.Chan, reflect.Func, reflect.Interface:
		return nil
	default:
		return newValue(typ, v, isRandom, isPointer)
	}
}

// mockStruct ..
func mockStruct(typ reflect.Type, v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := DeepElemValue(v.Field(i))
		fieldStructType := v.Type().Field(i)
		fieldType := DeepElemType(fieldStructType.Type)

		if fieldStructType.PkgPath == "" { // 只处理导出的字段
			key := structFieldName(fieldStructType, "json") // 使用json tag
			var val interface{}
			if isRecurse && !skipRecurse(fieldType) {
				val = mockAny(fieldType, fieldValue, isRandom, isPointer, isRecurse, skipRecurse)
			} else {
				val = newValue(fieldType, fieldValue, isRandom, isPointer)
			}
			result[key] = val
		}
	}
	return result
}

// mockSlice ..
func mockSlice(typ reflect.Type, v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	if v.Len() > 0 {
		return mockArray(typ, v, isRandom, isPointer, isRecurse, skipRecurse).([]interface{})
	}
	// 长度为空 填充一个
	var val interface{}
	itemType := DeepElemType(typ.Elem())
	itemValue := reflect.New(itemType).Elem()
	if isRecurse && !skipRecurse(itemType) {
		val = mockAny(itemType, itemValue, isRandom, isPointer, isRecurse, skipRecurse)
	} else {
		val = newValue(itemType, itemValue, isRandom, isPointer)
	}
	return []interface{}{val}
}

// mockArray ..
func mockArray(typ reflect.Type, v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		var val interface{}
		itemValue := v.Index(i)
		itemType := itemValue.Type()
		if isRecurse && !skipRecurse(itemType) {
			val = mockAny(itemType, itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemType, itemValue, isRandom, isPointer)
		}
		result = append(result, val)
	}
	return result
}

// mockMap ..
func mockMap(typ reflect.Type, v reflect.Value, isRandom, isPointer, isRecurse bool, skipRecurse func(typ reflect.Type) bool) interface{} {
	result := map[string]interface{}{}
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		var val interface{}
		itemValue := iter.Value()
		itemType := itemValue.Type()
		if isRecurse && !skipRecurse(itemType) {
			val = mockAny(itemType, itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemType, itemValue, isRandom, isPointer)
		}
		result[key] = val
	}
	// 没有key 填充一个
	if len(result) == 0 {
		var val interface{}
		itemType := DeepElemType(v.Type().Elem())
		itemValue := reflect.New(itemType).Elem()
		if isRecurse && !skipRecurse(itemType) {
			val = mockAny(itemType, itemValue, isRandom, isPointer, isRecurse, skipRecurse)
		} else {
			val = newValue(itemType, itemValue, isRandom, isPointer)
		}
		result[""] = val
	}
	return result
}
