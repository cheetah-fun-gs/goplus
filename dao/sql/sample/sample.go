package main

import (
	"database/sql"
	"fmt"
	"time"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	tableName      = "test_sqlplus_data"
	tableCreateSQL = `CREATE TABLE IF NOT EXISTS %v (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		a varchar(45) DEFAULT NULL,
		b int(11) DEFAULT NULL,
		c bigint(20) DEFAULT NULL,
		d datetime DEFAULT NULL,
		e timestamp NULL DEFAULT NULL,
		f float DEFAULT NULL,
		g decimal(10,2) DEFAULT NULL,
		h blob,
		i text,
		j date DEFAULT NULL,
		k time DEFAULT NULL,
		l binary(1) DEFAULT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='sqlplus 测试'`
)

type testData struct {
	ID int            `json:"id,omitempty"`
	A  sql.NullString `json:"a,omitempty"`
	B  int            `json:"b,omitempty"`
	C  int64          `json:"c,omitempty"`
	D  time.Time      `json:"d,omitempty"`
	E  time.Time      `json:"e,omitempty"`
	F  float64        `json:"f,omitempty"`
	G  float32        `json:"g,omitempty"`
	H  sql.NullString `json:"h,omitempty"`
	I  sql.NullString `json:"i,omitempty"`
	J  time.Time      `json:"j,omitempty"`
	K  sql.NullTime   `json:"k,omitempty"`
	L  []byte         `json:"l,omitempty"`
}

func main() {
	db, err := sql.Open("mysql", "admin:admin123@tcp(127.0.0.1:3306)/test?parseTime=true&charset=utf8mb4&loc=Asia%2FShanghai")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf(tableCreateSQL, tableName)); err != nil {
		panic(err)
	}

	if _, err = db.Exec(fmt.Sprintf("truncate table %v;", tableName)); err != nil {
		panic(err)
	}

	data1 := &testData{}
	data2 := &testData{B: 10}
	data3 := &testData{D: time.Now()}

	for _, data := range []*testData{data1, data2, data3} {
		query, args := sqlplus.GenInsert(tableName, data)
		if _, err = db.Exec(query, args...); err != nil {
			panic(err)
		}
	}

	query := fmt.Sprintf("SELECT * FROM %v LIMIT 10;", tableName)
	// Get
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rs := testData{}
	err = sqlplus.Get(rows, &rs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get Result: %+v\n", rs)

	// Select
	rows, err = db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rs2 := []*testData{}
	err = sqlplus.Select(rows, &rs2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Select Result:\n")
	for _, r := range rs2 {
		fmt.Printf("%+v\n", r)
	}
}
