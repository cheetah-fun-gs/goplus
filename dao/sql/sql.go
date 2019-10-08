// Package sql 对sql封装 支持在执行sql前后添加钩子
// 注意：如果在钩子里也有sql操作，注意不要死循环了
package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/knocknote/vitess-sqlparser/sqlparser"
)

func preHookSafe(ctx context.Context, queryStmt sqlparser.Statement, queryStr string, args []interface{}) error {
	switch v := queryStmt.(type) {
	case *sqlparser.Select:
		if _, ok := v.SelectExprs[0].(*sqlparser.StarExpr); ok {
			return fmt.Errorf("select must have column name, not allow *")
		}
	case *sqlparser.Insert:
		if len(v.Columns) == 0 {
			return fmt.Errorf("insert must have column name")
		}
	case *sqlparser.Update:
		if v.Where == nil {
			return fmt.Errorf("update must have where")
		}
	case *sqlparser.Delete:
		if v.Where == nil {
			return fmt.Errorf("delete must have where")
		}
	}
	return nil
}

// PreHook 前置钩子 err 非空会中断
type PreHook func(ctx context.Context, queryStmt sqlparser.Statement, queryStr string, args []interface{}) error

// PostHook 后置钩子 不中断
type PostHook func(ctx context.Context, queryStmt sqlparser.Statement, queryStr string, args []interface{}, result sql.Result)

func loopPreHook(ctx context.Context, hooks []PreHook, queryStmt sqlparser.Statement, queryStr string, args []interface{}) error {
	for _, hook := range hooks {
		if err := hook(ctx, queryStmt, queryStr, args); err != nil {
			return err
		}
	}
	return nil
}

func loopPostHook(ctx context.Context, hooks []PostHook, queryStmt sqlparser.Statement, queryStr string, args []interface{}, result sql.Result) {
	for _, hook := range hooks {
		hook(ctx, queryStmt, queryStr, args, result)
	}
	return
}

// Open Open
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

// OpenSafe OpenSafe 禁止 select * from 和 insert ino table values 等未指定列名的sql语句
func OpenSafe(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	d := &DB{DB: db}
	d.RegisterPreHook(preHookSafe)
	return d, nil
}

// DB sql plus
type DB struct {
	*sql.DB
	preHooks  []PreHook
	postHooks []PostHook
}

// RegisterPreHook 注册
func (db *DB) RegisterPreHook(preHooks ...PreHook) {
	db.preHooks = append(db.preHooks, preHooks...)
}

// RegisterPostHook 注册
func (db *DB) RegisterPostHook(postHooks ...PostHook) {
	db.postHooks = append(db.postHooks, postHooks...)
}

// Exec Exec
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

// ExecContext ExecContext
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(db.preHooks) != 0 || len(db.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, db.preHooks, queryStmt, query, args); err != nil {
		return nil, err
	}
	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return result, err
	}
	loopPostHook(ctx, db.postHooks, queryStmt, query, args, result)
	return result, nil
}

// Prepare Prepare
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}

// PrepareContext PrepareContext
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(db.preHooks) != 0 || len(db.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, db.preHooks, queryStmt, query, nil); err != nil {
		return nil, err
	}
	stmt, err := db.DB.PrepareContext(ctx, query)
	return &Stmt{
		Stmt:      stmt,
		queryStr:  query,
		queryStmt: queryStmt,
		preHooks:  db.preHooks,
		postHooks: db.postHooks,
	}, err
}

// Query Query
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

// QueryContext QueryContext
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(db.preHooks) != 0 || len(db.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, db.preHooks, queryStmt, query, args); err != nil {
		return nil, err
	}
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return rows, err
	}
	loopPostHook(ctx, db.postHooks, queryStmt, query, args, nil)
	return rows, nil
}

// Row Row
type Row struct {
	err error
	*sql.Row
}

// Scan Scan
func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	return r.Row.Scan(dest...)
}

// QueryRow QueryRow
func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext QueryRowContext
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	var queryStmt sqlparser.Statement
	var err error
	if len(db.preHooks) != 0 || len(db.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return &Row{err: err, Row: nil}
		}
	}

	if err := loopPreHook(ctx, db.preHooks, queryStmt, query, args); err != nil {
		return &Row{err: err, Row: nil}
	}
	row := db.DB.QueryRowContext(ctx, query, args...)
	loopPostHook(ctx, db.postHooks, queryStmt, query, args, nil)
	return &Row{err: nil, Row: row}
}

// Begin Begin
func (db *DB) Begin() (*Tx, error) {
	return db.BeginTx(context.Background(), nil)
}

// BeginTx BeginTx
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	return &Tx{
		Tx:        tx,
		preHooks:  db.preHooks,
		postHooks: db.postHooks,
	}, err
}

// Stmt plus
type Stmt struct {
	*sql.Stmt
	queryStmt sqlparser.Statement
	queryStr  string
	preHooks  []PreHook
	postHooks []PostHook
}

// Exec Exec
func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

// ExecContext ExecContext
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	result, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return result, err
	}

	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, result)
	return result, nil
}

// Query Query
func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryContext QueryContext
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return rows, err
	}
	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, nil)
	return rows, nil
}

// QueryRow QueryRow
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

// QueryRowContext QueryRowContext
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	row := s.Stmt.QueryRowContext(ctx, args...)
	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, nil)
	return &Row{err: nil, Row: row}
}

// Tx plus
type Tx struct {
	*sql.Tx
	preHooks  []PreHook
	postHooks []PostHook
}

// Exec Exec
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

// ExecContext ExecContext
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(tx.preHooks) != 0 || len(tx.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, tx.preHooks, queryStmt, query, args); err != nil {
		return nil, err
	}
	result, err := tx.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return result, err
	}
	loopPostHook(ctx, tx.postHooks, queryStmt, query, args, result)
	return result, nil
}

// Prepare Prepare
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

// PrepareContext PrepareContext
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(tx.preHooks) != 0 || len(tx.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, tx.preHooks, queryStmt, query, nil); err != nil {
		return nil, err
	}
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	return &Stmt{
		Stmt:      stmt,
		queryStmt: queryStmt,
		preHooks:  tx.preHooks,
		postHooks: tx.postHooks,
	}, err
}

// Query Query
func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

// QueryContext QueryContext
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var queryStmt sqlparser.Statement
	var err error
	if len(tx.preHooks) != 0 || len(tx.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return nil, err
		}
	}

	if err := loopPreHook(ctx, tx.preHooks, queryStmt, query, args); err != nil {
		return nil, err
	}
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return rows, err
	}
	loopPostHook(ctx, tx.postHooks, queryStmt, query, args, nil)
	return rows, nil
}

// QueryRow QueryRow
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext QueryRowContext
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	var queryStmt sqlparser.Statement
	var err error
	if len(tx.preHooks) != 0 || len(tx.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return &Row{err: err, Row: nil}
		}
	}

	if err := loopPreHook(ctx, tx.preHooks, queryStmt, query, args); err != nil {
		return &Row{err: err, Row: nil}
	}
	row := tx.Tx.QueryRowContext(ctx, query, args...)
	loopPostHook(ctx, tx.postHooks, queryStmt, query, args, nil)
	return &Row{err: nil, Row: row}
}
