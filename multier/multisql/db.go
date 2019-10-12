package multisql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
	"github.com/knocknote/vitess-sqlparser/sqlparser"
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
	onceDB sync.Once
	mDB    mutilDB
)

// InitDB 初始化db
func InitDB(defaultDB *sql.DB) {
	onceDB.Do(func() {
		mDB = mutilDB{
			d: &dbWithInterceptor{
				DB: defaultDB,
			},
		}
	})
}

// RegisterDB 注册 sql db
func RegisterDB(name string, db *sql.DB) error {
	if _, ok := mDB[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mDB[name] = &dbWithInterceptor{
		DB: db,
	}
	return nil
}

// InitDBWithInterceptor 初始化db
func InitDBWithInterceptor(defaultDB *sql.DB, in *sqlplus.Interceptor) {
	onceDB.Do(func() {
		mDB = mutilDB{
			d: &dbWithInterceptor{
				DB:          defaultDB,
				Interceptor: in,
			},
		}
	})
}

// RegisterDBWithInterceptor 注册 sql db
func RegisterDBWithInterceptor(name string, db *sql.DB, in *sqlplus.Interceptor) error {
	if _, ok := mDB[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mDB[name] = &dbWithInterceptor{
		DB:          db,
		Interceptor: in,
	}
	return nil
}

// GetDB ...
func GetDB() (*sql.DB, error) {
	return GetDBN(d)
}

// GetDBN ...
func GetDBN(name string) (*sql.DB, error) {
	if db, ok := mDB[name]; ok {
		return db.DB, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// Begin ...
func Begin() (*sql.Tx, error) {
	return BeginN(d)
}

// BeginTx ...
func BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
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
func Prepare(query string) (*sql.Stmt, error) {
	return PrepareN(d, query)
}

// PrepareContext ...
func PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return PrepareContextN(ctx, d, query)
}

// Query ...
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return QueryN(d, query, args...)
}

// QueryContext ...
func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return QueryContextN(ctx, d, query, args...)
}

// QueryRow ...
func QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	return QueryRowN(d, query, args...)
}

// QueryRowContext ...
func QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	return QueryRowContextN(ctx, d, query, args...)
}

// Stats ...
func Stats() sql.DBStats {
	return StatsN(d)
}

// BeginN ...
func BeginN(name string) (*sql.Tx, error) {
	if db, ok := mDB[name]; ok {
		return db.Begin()
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// BeginTxN ...
func BeginTxN(ctx context.Context, name string, opts *sql.TxOptions) (*sql.Tx, error) {
	if db, ok := mDB[name]; ok {
		return db.BeginTx(ctx, opts)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// ExecN ...
func ExecN(name, query string, args ...interface{}) (sql.Result, error) {
	return ExecContextN(context.Background(), name, query, args...)
}

// ExecContextN ...
func ExecContextN(ctx context.Context, name, query string, args ...interface{}) (sql.Result, error) {
	if db, ok := mDB[name]; ok {
		var execSQL *sqlplus.ExecSQL
		if db.Interceptor != nil {
			parse, _ := sqlparser.Parse(query)
			execSQL = &sqlplus.ExecSQL{
				Query: query,
				Args:  args,
				Parse: parse,
			}
			if err := db.Interceptor.BeforeExec(ctx, execSQL); err != nil {
				return nil, err
			}
		}
		result, err := db.ExecContext(ctx, query, args...)
		if db.Interceptor != nil && err == nil {
			db.Interceptor.BehindExec(ctx, execSQL, result)
		}
		return result, err
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// PrepareN ...
func PrepareN(name, query string) (*sql.Stmt, error) {
	return PrepareContextN(context.Background(), name, query)
}

// PrepareContextN ...
func PrepareContextN(ctx context.Context, name, query string) (*sql.Stmt, error) {
	if db, ok := mDB[name]; ok {
		if db.Interceptor != nil {
			parse, _ := sqlparser.Parse(query)
			execSQL := &sqlplus.ExecSQL{
				Query: query,
				Parse: parse,
			}
			if err := db.Interceptor.BeforeExec(ctx, execSQL); err != nil {
				return nil, err
			}
		}
		return db.PrepareContext(ctx, query)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// QueryN ...
func QueryN(name, query string, args ...interface{}) (*sql.Rows, error) {
	return QueryContextN(context.Background(), name, query, args...)
}

// QueryContextN ...
func QueryContextN(ctx context.Context, name, query string, args ...interface{}) (*sql.Rows, error) {
	if db, ok := mDB[name]; ok {
		var execSQL *sqlplus.ExecSQL
		if db.Interceptor != nil {
			parse, _ := sqlparser.Parse(query)
			execSQL = &sqlplus.ExecSQL{
				Query: query,
				Args:  args,
				Parse: parse,
			}
			if err := db.Interceptor.BeforeExec(ctx, execSQL); err != nil {
				return nil, err
			}
		}
		rows, err := db.QueryContext(ctx, query, args...)
		if db.Interceptor != nil && err == nil {
			db.Interceptor.BehindExec(ctx, execSQL, nil)
		}
		return rows, err
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// QueryRowN ...
func QueryRowN(name, query string, args ...interface{}) (*sql.Row, error) {
	return QueryRowContextN(context.Background(), name, query, args...)
}

// QueryRowContextN ...
func QueryRowContextN(ctx context.Context, name, query string, args ...interface{}) (*sql.Row, error) {
	if db, ok := mDB[name]; ok {
		var execSQL *sqlplus.ExecSQL
		if db.Interceptor != nil {
			parse, _ := sqlparser.Parse(query)
			execSQL = &sqlplus.ExecSQL{
				Query: query,
				Args:  args,
				Parse: parse,
			}
			if err := db.Interceptor.BeforeExec(ctx, execSQL); err != nil {
				return nil, err
			}
		}
		row := db.QueryRowContext(ctx, query, args...)
		if db.Interceptor != nil {
			db.Interceptor.BehindExec(ctx, execSQL, nil)
		}
		return row, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// StatsN ...
func StatsN(name string) sql.DBStats {
	if db, ok := mDB[name]; ok {
		return db.Stats()
	}
	return sql.DBStats{}
}
