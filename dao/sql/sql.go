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
func Open(driverName, dataSourceName string) (*DBPlus, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DBPlus{DB: db}, nil
}

// OpenSafe OpenSafe 禁止 select * from 和 insert ino table values 等未指定列名的sql语句
func OpenSafe(driverName, dataSourceName string) (*DBPlus, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	d := &DBPlus{DB: db}
	d.RegisterPreHook(preHookSafe)
	return d, nil
}

// DBPlus sql plus
type DBPlus struct {
	*sql.DB
	preHooks  []PreHook
	postHooks []PostHook
}

// RegisterPreHook 注册
func (db *DBPlus) RegisterPreHook(preHooks ...PreHook) {
	db.preHooks = append(db.preHooks, preHooks...)
}

// RegisterPostHook 注册
func (db *DBPlus) RegisterPostHook(postHooks ...PostHook) {
	db.postHooks = append(db.postHooks, postHooks...)
}

// Exec Exec
func (db *DBPlus) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

// ExecContext ExecContext
func (db *DBPlus) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
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
func (db *DBPlus) Prepare(query string) (*StmtPlus, error) {
	return db.PrepareContext(context.Background(), query)
}

// PrepareContext PrepareContext
func (db *DBPlus) PrepareContext(ctx context.Context, query string) (*StmtPlus, error) {
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
	return &StmtPlus{
		Stmt:      stmt,
		queryStr:  query,
		queryStmt: queryStmt,
		preHooks:  db.preHooks,
		postHooks: db.postHooks,
	}, err
}

// Query Query
func (db *DBPlus) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

// QueryContext QueryContext
func (db *DBPlus) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
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

// RowPlus RowPlus
type RowPlus struct {
	err error
	*sql.Row
}

// Scan Scan
func (r *RowPlus) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	return r.Row.Scan(dest...)
}

// QueryRow QueryRow
func (db *DBPlus) QueryRow(query string, args ...interface{}) *RowPlus {
	return db.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext QueryRowContext
func (db *DBPlus) QueryRowContext(ctx context.Context, query string, args ...interface{}) *RowPlus {
	var queryStmt sqlparser.Statement
	var err error
	if len(db.preHooks) != 0 || len(db.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return &RowPlus{err: err, Row: nil}
		}
	}

	if err := loopPreHook(ctx, db.preHooks, queryStmt, query, args); err != nil {
		return &RowPlus{err: err, Row: nil}
	}
	row := db.DB.QueryRowContext(ctx, query, args...)
	loopPostHook(ctx, db.postHooks, queryStmt, query, args, nil)
	return &RowPlus{err: nil, Row: row}
}

// Begin Begin
func (db *DBPlus) Begin() (*TxPlus, error) {
	return db.BeginTx(context.Background(), nil)
}

// BeginTx BeginTx
func (db *DBPlus) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TxPlus, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	return &TxPlus{
		Tx:        tx,
		preHooks:  db.preHooks,
		postHooks: db.postHooks,
	}, err
}

// StmtPlus StmtPlus
type StmtPlus struct {
	*sql.Stmt
	queryStmt sqlparser.Statement
	queryStr  string
	preHooks  []PreHook
	postHooks []PostHook
}

// Exec Exec
func (s *StmtPlus) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

// ExecContext ExecContext
func (s *StmtPlus) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	result, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return result, err
	}

	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, result)
	return result, nil
}

// Query Query
func (s *StmtPlus) Query(args ...interface{}) (*sql.Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryContext QueryContext
func (s *StmtPlus) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return rows, err
	}
	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, nil)
	return rows, nil
}

// QueryRow QueryRow
func (s *StmtPlus) QueryRow(args ...interface{}) *RowPlus {
	return s.QueryRowContext(context.Background(), args...)
}

// QueryRowContext QueryRowContext
func (s *StmtPlus) QueryRowContext(ctx context.Context, args ...interface{}) *RowPlus {
	row := s.Stmt.QueryRowContext(ctx, args...)
	loopPostHook(ctx, s.postHooks, s.queryStmt, s.queryStr, args, nil)
	return &RowPlus{err: nil, Row: row}
}

// TxPlus TxPlus
type TxPlus struct {
	*sql.Tx
	preHooks  []PreHook
	postHooks []PostHook
}

// Exec Exec
func (tx *TxPlus) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

// ExecContext ExecContext
func (tx *TxPlus) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
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
func (tx *TxPlus) Prepare(query string) (*StmtPlus, error) {
	return tx.PrepareContext(context.Background(), query)
}

// PrepareContext PrepareContext
func (tx *TxPlus) PrepareContext(ctx context.Context, query string) (*StmtPlus, error) {
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
	return &StmtPlus{
		Stmt:      stmt,
		queryStmt: queryStmt,
		preHooks:  tx.preHooks,
		postHooks: tx.postHooks,
	}, err
}

// Query Query
func (tx *TxPlus) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

// QueryContext QueryContext
func (tx *TxPlus) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
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
func (tx *TxPlus) QueryRow(query string, args ...interface{}) *RowPlus {
	return tx.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext QueryRowContext
func (tx *TxPlus) QueryRowContext(ctx context.Context, query string, args ...interface{}) *RowPlus {
	var queryStmt sqlparser.Statement
	var err error
	if len(tx.preHooks) != 0 || len(tx.postHooks) != 0 {
		queryStmt, err = sqlparser.Parse(query)
		if err != nil {
			return &RowPlus{err: err, Row: nil}
		}
	}

	if err := loopPreHook(ctx, tx.preHooks, queryStmt, query, args); err != nil {
		return &RowPlus{err: err, Row: nil}
	}
	row := tx.Tx.QueryRowContext(ctx, query, args...)
	loopPostHook(ctx, tx.postHooks, queryStmt, query, args, nil)
	return &RowPlus{err: nil, Row: row}
}
