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
	Accept(visitor StmtVisitor) *Value
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) *Value
	VisitPrintStmt(stmt *PrintStmt) *Value
	VisitVarStmt(stmt *VarStmt) *Value
}

type ExpressionStmt struct {
	Expr Expr
}

func (es *ExpressionStmt) Type() StmtType {
	return EXPRESSION_STMT
}

func (es *ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt: %s", es.Expr.String())
}

func (es *ExpressionStmt) Accept(visitor StmtVisitor) *Value {
	return visitor.VisitExpressionStmt(es)
}

type PrintStmt struct {
	Expr Expr
}

func (ps *PrintStmt) Type() StmtType {
	return PRINT_STMT
}

func (ps *PrintStmt) String() string {
	return fmt.Sprintf("PrintStmt: %s", ps.Expr.String())
}

func (ps *PrintStmt) Accept(visitor StmtVisitor) *Value {
	return visitor.VisitPrintStmt(ps)
}

type VarStmt struct {
	Name *ls.Token
	Expr Expr
}

func (vs *VarStmt) Type() StmtType {
	return VAR_STMT
}

func (vs *VarStmt) String() string {
	exprStr := "nil"
	if vs.Expr != nil {
		exprStr = vs.Expr.String()
	}
	return fmt.Sprintf("VarStmt: %s = %s", vs.Name.Lexeme, exprStr)
}

func (vs *VarStmt) Accept(visitor StmtVisitor) *Value {
	return visitor.VisitVarStmt(vs)
}
