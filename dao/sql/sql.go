package sql

import (
	"database/sql"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	reflectplus "github.com/cheetah-fun-gs/goplus/reflect"
)

func scanOne(rows *sql.Rows, columns []string, fields map[string]interface{}) (map[string]interface{}, error) {
	var dest []interface{}

	for _, name := range columns {
		val, ok := fields[name]
		if !ok {
			val = sql.RawBytes{}
		}
		dest = append(dest, val)
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

// Get v map[string]interface{} or struct
func Get(rows *sql.Rows, v interface{}) error {
	fields, ok := v.(map[string]interface{})
	if !ok {
		fields = reflectplus.MockStruct(v, true, false)
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
	vv := reflectplus.MockSlice(v, true, false)[0]

	fields, ok := vv.(map[string]interface{})
	if !ok {
		fields = reflectplus.MockStruct(vv, true, false)
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
