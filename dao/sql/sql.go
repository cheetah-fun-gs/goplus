package sql

import (
	"database/sql"
	"fmt"

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

// Get ...
func Get(rows *sql.Rows, v interface{}) error {
	fields, ok := reflectplus.Mock(v).(map[string]interface{})
	if !ok {
		return fmt.Errorf("v must be a map or a struct")
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
	l, ok := reflectplus.Mock(v).([]interface{})
	if !ok {
		return fmt.Errorf("v must be a slice")
	}

	if len(l) == 0 {
		return fmt.Errorf("elem is blank")
	}

	fields, ok := l[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("elem must be a struct")
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
