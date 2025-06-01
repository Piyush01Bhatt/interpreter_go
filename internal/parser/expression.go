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
)

type Value struct {
	StrVal   *string
	IntVal   *int
	FloatVal *float64
	BoolVal  *bool
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

type Expr interface {
	Type() ExprType
	String() string
}

type Binary struct {
	left     Expr
	operator *ls.Token
	right    Expr
}

func (b *Binary) Type() ExprType {
	return BINARY
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.left, b.operator.Lexeme, b.right)
}

func TestBinary() Expr {
	binaryExp := Binary{
		left: &Unary{
			operator: &ls.Token{
				Type:    ls.MINUS,
				Lexeme:  "-",
				Literal: "-",
				Line:    0,
			},
			right: &Literal{
				value: NewIntValue(123),
			},
		},
		operator: &ls.Token{
			Type:    ls.STAR,
			Lexeme:  "*",
			Literal: "*",
			Line:    0,
		},
		right: &Literal{
			value: NewIntValue(15),
		},
	}
	return &binaryExp
}

type Unary struct {
	operator *ls.Token
	right    Expr
}

func (*Unary) Type() ExprType {
	return UNARY
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.operator.Lexeme, u.right)
}

type Literal struct {
	value *Value
}

func (l *Literal) Type() ExprType {
	return LITERAL
}

func (l *Literal) String() string {
	return l.value.String()
}
