package mulsql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
)

const (
	d = "default"
)

type mutilDB map[string]sqlplus.DB

var (
	once  sync.Once
	mutil mutilDB
)

// Init 初始化
func Init(defaultDB sqlplus.DB) {
	once.Do(func() {
		mutil = mutilDB{
			d: defaultDB,
		}
	})
}

// Register 注册 sql db
func Register(name string, db sqlplus.DB) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = db
	return nil
}

// GetDB ...
func GetDB() (sqlplus.DB, error) {
	return GetDBN(d)
}

// GetDBN ...
func GetDBN(name string) (sqlplus.DB, error) {
	if db, ok := mutil[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// Begin ...
func Begin() (sqlplus.Tx, error) {
	return BeginN(d)
}

// BeginTx ...
func BeginTx(ctx context.Context, opts *sql.TxOptions) (sqlplus.Tx, error) {
	return BeginTxN(ctx, d, opts)
}

// Exec ...
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return ExecN(d, query, args...)
}

// ExecContext ...
func ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return ExecContextN(ctx, d, query, args...)
}

// Prepare ...
func Prepare(query string) (sqlplus.Stmt, error) {
	return PrepareN(d, query)
}

// PrepareContext ...
func PrepareContext(ctx context.Context, query string) (sqlplus.Stmt, error) {
	return PrepareContextN(ctx, d, query)
}

// Query ...
func Query(query string, args ...interface{}) (sqlplus.Rows, error) {
	return QueryN(d, query, args...)
}

// QueryContext ...
func QueryContext(ctx context.Context, query string, args ...interface{}) (sqlplus.Rows, error) {
	return QueryContextN(ctx, d, query, args...)
}

// QueryRow ...
func QueryRow(query string, args ...interface{}) sqlplus.Row {
	return QueryRowN(d, query, args...)
}

// QueryRowContext ...
func QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlplus.Row {
	return QueryRowContextN(ctx, d, query, args...)
}

// Stats ...
func Stats() sql.DBStats {
	return StatsN(d)
}

// BeginN ...
func BeginN(name string) (sqlplus.Tx, error) {
	if db, ok := mutil[name]; ok {
		return db.Begin()
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// BeginTxN ...
func BeginTxN(ctx context.Context, name string, opts *sql.TxOptions) (sqlplus.Tx, error) {
	if db, ok := mutil[name]; ok {
		return db.BeginTx(ctx, opts)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// ExecN ...
func ExecN(name, query string, args ...interface{}) (sql.Result, error) {
	if db, ok := mutil[name]; ok {
		return db.Exec(query, args...)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// ExecContextN ...
func ExecContextN(ctx context.Context, name, query string, args ...interface{}) (sql.Result, error) {
	if db, ok := mutil[name]; ok {
		return db.ExecContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// PrepareN ...
func PrepareN(name, query string) (sqlplus.Stmt, error) {
	if db, ok := mutil[name]; ok {
		return db.Prepare(query)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// PrepareContextN ...
func PrepareContextN(ctx context.Context, name, query string) (sqlplus.Stmt, error) {
	if db, ok := mutil[name]; ok {
		return db.PrepareContext(ctx, query)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// QueryN ...
func QueryN(name, query string, args ...interface{}) (sqlplus.Rows, error) {
	if db, ok := mutil[name]; ok {
		return db.Query(query, args...)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// QueryContextN ...
func QueryContextN(ctx context.Context, name, query string, args ...interface{}) (sqlplus.Rows, error) {
	if db, ok := mutil[name]; ok {
		return db.QueryContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// QueryRowN ...
func QueryRowN(name, query string, args ...interface{}) sqlplus.Row {
	if db, ok := mutil[name]; ok {
		return db.QueryRow(query, args...)
	}
	return &sqlplus.ErrRow{Err: fmt.Errorf("name not found: %v", name)}
}

// QueryRowContextN ...
func QueryRowContextN(ctx context.Context, name, query string, args ...interface{}) sqlplus.Row {
	if db, ok := mutil[name]; ok {
		return db.QueryRowContext(ctx, query, args...)
	}
	return &sqlplus.ErrRow{Err: fmt.Errorf("name not found: %v", name)}
}

// StatsN ...
func StatsN(name string) sql.DBStats {
	if db, ok := mutil[name]; ok {
		return db.Stats()
	}
	return sql.DBStats{}
}
