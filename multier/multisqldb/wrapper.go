package multisqldb

import (
	"context"
	"database/sql"
	"fmt"

	sqlplus "github.com/cheetah-fun-gs/goplus/dao/sql"
	"github.com/knocknote/vitess-sqlparser/sqlparser"
)

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
	if db, ok := mutil[name]; ok {
		return db.Begin()
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// BeginTxN ...
func BeginTxN(ctx context.Context, name string, opts *sql.TxOptions) (*sql.Tx, error) {
	if db, ok := mutil[name]; ok {
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
	if db, ok := mutil[name]; ok {
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
	if db, ok := mutil[name]; ok {
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
	if db, ok := mutil[name]; ok {
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
	if db, ok := mutil[name]; ok {
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
	if db, ok := mutil[name]; ok {
		return db.Stats()
	}
	return sql.DBStats{}
}

// Get ...
func Get(query string, v interface{}, args ...interface{}) error {
	return GetContextN(context.Background(), d, query, v, args...)
}

// GetContext ...
func GetContext(ctx context.Context, query string, v interface{}, args ...interface{}) error {
	return GetContextN(ctx, d, query, v, args...)
}

// GetN ...
func GetN(name, query string, v interface{}, args ...interface{}) error {
	return GetContextN(context.Background(), name, query, v, args...)
}

// GetContextN ...
func GetContextN(ctx context.Context, name, query string, v interface{}, args ...interface{}) error {
	if db, ok := mutil[name]; ok {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		return sqlplus.Get(rows, v)
	}
	return fmt.Errorf("name not found: %v", name)
}

// Select ...
func Select(query string, v interface{}, args ...interface{}) error {
	return SelectContextN(context.Background(), d, query, v, args...)
}

// SelectContext ...
func SelectContext(ctx context.Context, query string, v interface{}, args ...interface{}) error {
	return SelectContextN(ctx, d, query, v, args...)
}

// SelectN ...
func SelectN(name, query string, v interface{}, args ...interface{}) error {
	return SelectContextN(context.Background(), name, query, v, args...)
}

// SelectContextN ...
func SelectContextN(ctx context.Context, name, query string, v interface{}, args ...interface{}) error {
	if db, ok := mutil[name]; ok {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		return sqlplus.Select(rows, v)
	}
	return fmt.Errorf("name not found: %v", name)
}

// Insert ...
func Insert(tableName string, v interface{}) (sql.Result, error) {
	return InsertContextN(context.Background(), d, tableName, v)
}

// InsertContext ...
func InsertContext(ctx context.Context, name, tableName string, v interface{}) (sql.Result, error) {
	return InsertContextN(ctx, d, tableName, v)
}

// InsertN ...
func InsertN(name, tableName string, v interface{}) (sql.Result, error) {
	return InsertContextN(context.Background(), name, tableName, v)
}

// InsertContextN ...
func InsertContextN(ctx context.Context, name, tableName string, v interface{}) (sql.Result, error) {
	if db, ok := mutil[name]; ok {
		query, args := sqlplus.GenInsert(tableName, v)
		return db.ExecContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("name not found: %v", name)
}
