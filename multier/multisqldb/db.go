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

// Init 初始化db
func Init(defaultDB *sql.DB) {
	once.Do(func() {
		mutil = mutilDB{
			d: &dbWithInterceptor{
				DB: defaultDB,
			},
		}
	})
}

// Register 注册 sql db
func Register(name string, db *sql.DB) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = &dbWithInterceptor{
		DB: db,
	}
	return nil
}

// InitWithInterceptor 初始化db
func InitWithInterceptor(defaultDB *sql.DB, in *sqlplus.Interceptor) {
	once.Do(func() {
		mutil = mutilDB{
			d: &dbWithInterceptor{
				DB:          defaultDB,
				Interceptor: in,
			},
		}
	})
}

// RegisterWithInterceptor 注册 sql db
func RegisterWithInterceptor(name string, db *sql.DB, in *sqlplus.Interceptor) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = &dbWithInterceptor{
		DB:          db,
		Interceptor: in,
	}
	return nil
}

// Retrieve 获取 *sql.DB
func Retrieve() *sql.DB {
	return mutil[d].DB
}

// RetrieveAll 获取所有*sql.DB
func RetrieveAll() map[string]*sql.DB {
	r := map[string]*sql.DB{}
	for k, v := range mutil {
		r[k] = v.DB
	}
	return r
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
		panic(err)
	}
	return c
}
