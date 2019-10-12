package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/knocknote/vitess-sqlparser/sqlparser"
)

// ExecSQL 执行语句
type ExecSQL struct {
	Query string
	Args  []interface{}
	Parse sqlparser.Statement // Query 的解析 不含 args
}

// BeforeExec exec 前调用 发生异常会会中断执行
type BeforeExec func(ctx context.Context, execSQL *ExecSQL) error

// BehindExec exec 后调用 不会中断 v 可能是 sql.Result
type BehindExec func(ctx context.Context, execSQL *ExecSQL, result sql.Result)

// Interceptor 拦截器
type Interceptor struct {
	beforeExeces []BeforeExec
	behindExeces []BehindExec
}

// New ...
func New() *Interceptor {
	return &Interceptor{
		beforeExeces: []BeforeExec{},
		behindExeces: []BehindExec{},
	}
}

// NewSafe ...
func NewSafe() *Interceptor {
	in := &Interceptor{
		beforeExeces: []BeforeExec{},
		behindExeces: []BehindExec{},
	}
	in.RegisterBeforeExec(beforeExecSafe)
	return in
}

// RegisterBeforeExec ...
func (in *Interceptor) RegisterBeforeExec(f ...BeforeExec) {
	in.beforeExeces = append(in.beforeExeces, f...)
}

// RegisterBehindExec ...
func (in *Interceptor) RegisterBehindExec(f ...BehindExec) {
	in.behindExeces = append(in.behindExeces, f...)
}

// BeforeExec ...
func (in *Interceptor) BeforeExec(ctx context.Context, execSQL *ExecSQL) error {
	for _, f := range in.beforeExeces {
		if err := f(ctx, execSQL); err != nil {
			return err
		}
	}
	return nil
}

// BehindExec ...
func (in *Interceptor) BehindExec(ctx context.Context, execSQL *ExecSQL, result sql.Result) {
	for _, f := range in.behindExeces {
		f(ctx, execSQL, result)
	}
}

func beforeExecSafe(ctx context.Context, execSQL *ExecSQL) error {
	switch v := execSQL.Parse.(type) {
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
