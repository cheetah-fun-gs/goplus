package reflect

import (
	"reflect"
	"strings"
)

// DeepElemValue 追踪指针获取最终的 elem value
func DeepElemValue(v reflect.Value) reflect.Value {
	for {
		switch v.Kind() {
		case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
			if v.Elem().IsValid() {
				v = v.Elem()
				continue
			}
			v = reflect.New(v.Type().Elem())
		default:
			return v
		}
	}
}

func structFieldName(field reflect.StructField, tagName string) string {
	if field.Tag == "" {
		return field.Name
	} else if tagValue := field.Tag.Get(tagName); tagValue == "" {
		return field.Name
	} else if splits := strings.Split(tagValue, ","); splits[0] == "" {
		return field.Name
	} else {
		return splits[0]
	}
}
