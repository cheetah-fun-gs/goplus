package multimgo

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
	mDB    mutilDB
)

// InitDB 初始化
func InitDB(defaultDB *mgo.Database) {
	onceDB.Do(func() {
		mDB = mutilDB{
			d: defaultDB,
		}
	})
}

// RegisterDB 注册连接池
func RegisterDB(name string, db *mgo.Database) error {
	if _, ok := mDB[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mDB[name] = db
	return nil
}

// GetDB 获取连接池
func GetDB() (*mgo.Database, error) {
	return GetDBN(d)
}

// GetDBN get db with name 获取连接池
func GetDBN(name string) (*mgo.Database, error) {
	if db, ok := mDB[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}
