package parser

import (
	"fmt"

	ls "github.com/Piyush01Bhatt/interpreter_go/internal/scanner"
)

type StmtType int

const (
	EXPRESSION_STMT StmtType = iota
	PRINT_STMT
	VAR_STMT
)

type Stmt interface {
	Type() StmtType
	String() string
	Execute() *Value
}

type ExpressionStmt struct {
	expr Expr
}

func (es *ExpressionStmt) Type() StmtType {
	return EXPRESSION_STMT
}

func (es *ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt: %s", es.expr.String())
}

func (es *ExpressionStmt) Execute() *Value {
	return es.expr.Evaluate()
}

type PrintStmt struct {
	expr Expr
}

func (ps *PrintStmt) Type() StmtType {
	return PRINT_STMT
}

func (ps *PrintStmt) String() string {
	return fmt.Sprintf("PrintStmt: %s", ps.expr.String())
}

func (ps *PrintStmt) Execute() *Value {
	return ps.expr.Evaluate()
}

type VarStmt struct {
	name *ls.Token
	expr Expr
}

func (vs *VarStmt) Type() StmtType {
	return VAR_STMT
}

func (vs *VarStmt) String() string {
	exprStr := "nil"
	if vs.expr != nil {
		exprStr = vs.expr.String()
	}
	return fmt.Sprintf("VarStmt: %s = %s", vs.name.Lexeme, exprStr)
}

func (vs *VarStmt) Execute() *Value {
	if vs.expr == nil {
		return NewNilValue()
	}
	return vs.expr.Evaluate()
}
