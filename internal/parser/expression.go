package parser

import (
	"fmt"
	"strconv"

	ls "github.com/Piyush01Bhatt/interpreter_go/internal/scanner"
)

type ExprType int

const (
	BINARY ExprType = iota
	UNARY
	LITERAL
	VARIABLE
	ASSIGN
)

type Value struct {
	StrVal   *string
	IntVal   *int
	FloatVal *float64
	BoolVal  *bool
	NilVal   *struct{}
}

func NewStringValue(s string) *Value {
	return &Value{StrVal: &s}
}

func NewIntValue(i int) *Value {
	return &Value{IntVal: &i}
}

func NewFloatValue(f float64) *Value {
	return &Value{FloatVal: &f}
}

func NewBoolValue(b bool) *Value {
	return &Value{BoolVal: &b}
}

func NewNilValue() *Value {
	return &Value{NilVal: &struct{}{}}
}

func (v *Value) String() string {
	switch {
	case v.StrVal != nil:
		return fmt.Sprintf("%q", *v.StrVal) // Quote strings
	case v.IntVal != nil:
		return strconv.Itoa(*v.IntVal)
	case v.FloatVal != nil:
		return fmt.Sprintf("%g", *v.FloatVal) // Avoid unnecessary trailing zeros
	case v.BoolVal != nil:
		return strconv.FormatBool(*v.BoolVal)
	default:
		return "nil"
	}
}

func (v Value) GetType() string {
	switch {
	case v.IntVal != nil:
		return "int"
	case v.FloatVal != nil:
		return "float"
	case v.StrVal != nil:
		return "string"
	case v.BoolVal != nil:
		return "bool"
	default:
		return "nil"
	}
}

func (v Value) IsNumber() bool {
	return v.IntVal != nil || v.FloatVal != nil
}

func (v Value) IsString() bool {
	return v.StrVal != nil
}

func (v Value) ToFloat64() float64 {
	if v.FloatVal != nil {
		return *v.FloatVal
	}
	if v.IntVal != nil {
		return float64(*v.IntVal)
	}
	if v.BoolVal != nil {
		if *v.BoolVal {
			return 1.0
		}
		return 0.0
	}
	panic("Not a numeric value")
}

func (v Value) IsBool() bool {
	return v.BoolVal != nil
}

func (v Value) IsNil() bool {
	return v.StrVal == nil && v.IntVal == nil && v.FloatVal == nil && v.BoolVal == nil
}

func (v Value) IsTruthy() bool {
	switch {
	case v.BoolVal != nil:
		return *v.BoolVal
	case v.IntVal != nil:
		return *v.IntVal != 0
	case v.FloatVal != nil:
		return *v.FloatVal != 0.0
	case v.StrVal != nil:
		return *v.StrVal != ""
	default:
		return false // nil is false
	}
}

type Expr interface {
	Type() ExprType
	String() string
	Accept(visitor ExprVisitor) *Value
}

type ExprVisitor interface {
	VisitBinary(binary *Binary) *Value
	VisitUnary(unary *Unary) *Value
	VisitLiteral(literal *Literal) *Value
	VisitVariable(variable *Variable) *Value
	VisitAssign(assign *Assign) *Value
}

type Binary struct {
	Left     Expr
	Operator *ls.Token
	Right    Expr
}

func (b *Binary) Type() ExprType {
	return BINARY
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator.Lexeme, b.Right)
}

func (b *Binary) Accept(visitor ExprVisitor) *Value {
	return visitor.VisitBinary(b)
}

type Unary struct {
	Operator *ls.Token
	Right    Expr
}

func (*Unary) Type() ExprType {
	return UNARY
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Right)
}

func (u *Unary) Accept(visitor ExprVisitor) *Value {
	return visitor.VisitUnary(u)
}

type Literal struct {
	Value *Value
}

func (l *Literal) Type() ExprType {
	return LITERAL
}

func (l *Literal) String() string {
	return l.Value.String()
}

func (l *Literal) Accept(visitor ExprVisitor) *Value {
	return visitor.VisitLiteral(l)
}

type Variable struct {
	Name string
}

func (v *Variable) Type() ExprType {
	return VARIABLE
}

func (v *Variable) String() string {
	return v.Name
}

func (v *Variable) Accept(visitor ExprVisitor) *Value {
	return visitor.VisitVariable(v)
}

type Assign struct {
	Name string
	Expr Expr
}

func (a *Assign) Type() ExprType {
	return ASSIGN
}

func (a *Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Name, a.Expr)
}

func (a *Assign) Accept(visitor ExprVisitor) *Value {
	return visitor.VisitAssign(a)
}
