package sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	jsonplus "github.com/cheetah-fun-gs/goplus/encoding/json"
	reflectplus "github.com/cheetah-fun-gs/goplus/reflect"
)

// RowsAffected ...
func RowsAffected(result sql.Result, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

// LastInsertId ...
func LastInsertId(result sql.Result, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func scanOne(rows *sql.Rows, columns []string, fields map[string]interface{}) (map[string]interface{}, error) {
	var dest []interface{}

	for _, name := range columns {
		val, ok := fields[name]
		if !ok {
			return nil, fmt.Errorf("column not in fields: %v", name)
		}
		fmt.Printf("%+v, %v\n", val, reflect.TypeOf(val))
		dest = append(dest, val)
	}

	if err := rows.Scan(dest...); err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	for i := 0; i < len(columns); i++ {
		result[columns[i]] = dest[i]
	}
	fmt.Printf("****\n")
	for _, val := range result {
		fmt.Println(reflect.ValueOf(val).Elem())
	}

	return result, nil
}

// Get v map[string]interface{} or struct
func Get(rows *sql.Rows, v interface{}) error {
	fields, ok := v.(map[string]interface{})
	if !ok {
		fields = reflectplus.Mock(v).DisableRecurse().Pointer().Value().(map[string]interface{})
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
	vv := reflectplus.Mock(v).DisableRecurse().Pointer().Value()
	vvv, ok := vv.([]interface{})
	if !ok {
		return fmt.Errorf("v is not []interface{}")
	}

	fields, ok := vvv[0].(map[string]interface{})
	if !ok {
		fields = reflectplus.Mock(vvv[0]).DisableRecurse().Pointer().Value().(map[string]interface{})
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

// GenInsert 生成insert sql
func GenInsert(tableName string, v interface{}) (string, []interface{}) {
	fields, ok := v.(map[string]interface{})
	if !ok {
		fields = reflectplus.Mock(v).DisableRecurse().Pointer().Value().(map[string]interface{})
	}

	coloums := []string{}
	args := []interface{}{}
	marks := []string{}
	for key, val := range fields {
		coloums = append(coloums, key)
		args = append(args, val)
		marks = append(marks, "?")
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName,
		strings.Join(coloums, ", "), strings.Join(marks, ", ")), args
}
