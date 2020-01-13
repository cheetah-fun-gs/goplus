package sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
)

func structFieldName(field reflect.StructField) string {
	if field.Tag == "" {
		return field.Name
	} else if jsonTag := field.Tag.Get("json"); jsonTag == "" { // 复用json tag
		return field.Name
	} else if splits := strings.Split(jsonTag, ","); splits[0] == "" {
		return field.Name
	} else {
		return splits[0]
	}
}

func mapToFields(v reflect.Value) (map[string]reflect.Value, error) {
	for {
		v = v.Elem()
		switch v.Kind() {
		case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
			continue
		case reflect.Map:
		default:
			return nil, fmt.Errorf("elem is not a map")
		}

		fields := map[string]reflect.Value{}

		iter := v.MapRange()
		for iter.Next() {
			if iter.Key().Kind() != reflect.String {
				return nil, fmt.Errorf("map key is not string")
			}
			fields[iter.Key().String()] = iter.Value()
		}

		if len(fields) == 0 {
			return nil, fmt.Errorf("map is blank")
		}

		return fields, nil
	}
}

func structToFields(t reflect.Type) (map[string]reflect.Value, error) {
	for {
		switch t.Kind() {
		case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
			t = t.Elem()
			continue
		case reflect.Struct:
		default:
			return nil, fmt.Errorf("elem is not a struct")
		}

		fields := map[string]reflect.Value{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fields[structFieldName(field)] = reflect.New(field.Type)
		}

		if len(fields) == 0 {
			return nil, fmt.Errorf("struct is blank")
		}
		return fields, nil
	}
}

func scanOne(rows *sql.Rows, columns []string, fields map[string]reflect.Value) (map[string]interface{}, error) {
	var dest []interface{}
	for _, name := range columns {
		val, ok := fields[name]
		if !ok {
			val = reflect.ValueOf(sql.RawBytes{})
		}

		dest = append(dest, val.Interface())
	}

	if err := rows.Scan(dest...); err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	for i := 0; i < len(columns); i++ {
		result[columns[i]] = dest[i]
	}
	return result, nil
}

// Get ...
func Get(rows *sql.Rows, v interface{}) error {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
	default:
		return fmt.Errorf("v is not a pointer")
	}

	val = val.Elem()

	var fields map[string]reflect.Value
	var err error

	switch val.Kind() {
	case reflect.Map:
		fields, err = mapToFields(val)
	case reflect.Struct:
		fields, err = structToFields(val.Type())
	default:
		return fmt.Errorf("elem is not a map or a struct")
	}

	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	rows.Next()
	data, err := scanOne(rows, columns, fields)
	if err != nil {
		return err
	}

	return jsonplus.Convert(data, v)
}

// Select ...
func Select(rows *sql.Rows, v interface{}) error {
	val := reflect.TypeOf(v)
	switch val.Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
	default:
		return fmt.Errorf("v is not a pointer")
	}

	val = val.Elem()
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("elem is not a slice")
	}

	fields, err := structToFields(val.Elem())
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	data := []map[string]interface{}{}
	for rows.Next() {
		result, err := scanOne(rows, columns, fields)
		if err != nil {
			return err
		}
		data = append(data, result)
	}
	return jsonplus.Convert(data, v)
}
