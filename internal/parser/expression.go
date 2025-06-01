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
	Evaluate() *Value
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

func (b *Binary) Evaluate() *Value {
	left := b.left.Evaluate()
	right := b.right.Evaluate()

	if !left.IsNumber() || !right.IsNumber() {
		panic("Type mismatch")
	}

	leftVal := left.ToFloat64()
	rightVal := right.ToFloat64()

	switch b.operator.Type {
	case ls.MINUS:
		return NewFloatValue(leftVal - rightVal)
	case ls.SLASH:
		return NewFloatValue(leftVal / rightVal)
	case ls.STAR:
		return NewFloatValue(leftVal * rightVal)
	case ls.PLUS:
		return NewFloatValue(leftVal + rightVal)
	}

	return nil
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

func (u *Unary) Evaluate() *Value {
	right := u.right.Evaluate()
	op := u.operator.Type
	if op == ls.MINUS && !right.IsNumber() {
		panic("Type mismatch")
	}

	rightVal := right.ToFloat64()

	switch op {
	case ls.MINUS:
		return NewFloatValue(-rightVal)
	case ls.BANG:
		return NewBoolValue(!right.IsTruthy())
	}

	return nil
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

func (l *Literal) Evaluate() *Value {
	return l.value
}
