package multimgodb

import (
	"fmt"
	"sync"

	"github.com/globalsign/mgo"
)

const (
	d = "default"
)

type mutilDB map[string]*mgo.Database

var (
	onceDB sync.Once
	mutil  mutilDB
)

// Init 初始化
func Init(defaultDB *mgo.Database) {
	onceDB.Do(func() {
		mutil = mutilDB{
			d: defaultDB,
		}
	})
}

// Register 注册连接池
func Register(name string, db *mgo.Database) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = db
	return nil
}

// Retrieve 获取 mgo.Database
func Retrieve() *mgo.Database {
	return mutil[d]
}

// RetrieveN 获取 mgo.Database
func RetrieveN(name string) (*mgo.Database, error) {
	if c, ok := mutil[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// MustRetrieveN 获取 mgo.Database
func MustRetrieveN(name string) *mgo.Database {
	c, err := RetrieveN(name)
	if err != nil {
		panic((err))
	}
	return c
}
