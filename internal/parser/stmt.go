package parser

import (
	"fmt"
)

type StmtType int

const (
	EXPRESSION_STMT StmtType = iota
	PRINT_STMT
)

type Stmt interface {
	Type() StmtType
	String() string
	Execute() *Value
}

type ExpressionStmt struct {
	expr Expr
}

func (s *ExpressionStmt) Type() StmtType {
	return EXPRESSION_STMT
}

func (s *ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt: %s", s.expr.String())
}

func (s *ExpressionStmt) Execute() *Value {
	return s.expr.Evaluate()
}

type PrintStmt struct {
	expr Expr
}

func (s *PrintStmt) Type() StmtType {
	return PRINT_STMT
}

func (s *PrintStmt) String() string {
	return fmt.Sprintf("PrintStmt: %s", s.expr.String())
}

func (s *PrintStmt) Execute() *Value {
	return s.expr.Evaluate()
}
