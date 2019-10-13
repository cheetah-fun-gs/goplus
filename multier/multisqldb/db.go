package multisqldb

import (
	"database/sql"
	"fmt"
	"sync"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
)

const (
	d = "default"
)

type dbWithInterceptor struct {
	*sql.DB
	*sqlplus.Interceptor
}

type mutilDB map[string]*dbWithInterceptor

var (
	once  sync.Once
	mutil mutilDB
)

// InitDB 初始化db
func InitDB(defaultDB *sql.DB) {
	once.Do(func() {
		mutil = mutilDB{
			d: &dbWithInterceptor{
				DB: defaultDB,
			},
		}
	})
}

// RegisterDB 注册 sql db
func RegisterDB(name string, db *sql.DB) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = &dbWithInterceptor{
		DB: db,
	}
	return nil
}

// InitDBWithInterceptor 初始化db
func InitDBWithInterceptor(defaultDB *sql.DB, in *sqlplus.Interceptor) {
	once.Do(func() {
		mutil = mutilDB{
			d: &dbWithInterceptor{
				DB:          defaultDB,
				Interceptor: in,
			},
		}
	})
}

// RegisterDBWithInterceptor 注册 sql db
func RegisterDBWithInterceptor(name string, db *sql.DB, in *sqlplus.Interceptor) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = &dbWithInterceptor{
		DB:          db,
		Interceptor: in,
	}
	return nil
}

// Retrieve 获取 *redigo.Pool
func Retrieve() *sql.DB {
	return mutil[d].DB
}

// RetrieveN 获取 *sql.DB
func RetrieveN(name string) (*sql.DB, error) {
	if c, ok := mutil[name]; ok {
		return c.DB, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// MustRetrieveN 获取 *sql.DB
func MustRetrieveN(name string) *sql.DB {
	c, err := RetrieveN(name)
	if err != nil {
		panic((err))
	}
	return c
}
