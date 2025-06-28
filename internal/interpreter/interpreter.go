package interpreter

import (
	"fmt"

	"github.com/Piyush01Bhatt/interpreter_go/internal/parser"
	ls "github.com/Piyush01Bhatt/interpreter_go/internal/scanner"
)

type ExecutionMode int

const (
	ModePrompt ExecutionMode = iota
	ModeFile
)

type Interpreter struct {
	environment *Env
	mode        ExecutionMode
}

func NewInterpreter(mode ExecutionMode) *Interpreter {
	return &Interpreter{
		environment: NewEnv(),
		mode:        mode,
	}
}

// Implement ExprVisitor
func (i *Interpreter) VisitBinary(expr *parser.Binary) *parser.Value {
	left := expr.Left.Accept(i)
	right := expr.Right.Accept(i)

	return i.evaluateBinaryOp(left, right, expr.Operator)
}

func (i *Interpreter) VisitUnary(expr *parser.Unary) *parser.Value {
	right := expr.Right.Accept(i)
	return i.evaluateUnaryOp(right, expr.Operator)
}

func (i *Interpreter) VisitLiteral(expr *parser.Literal) *parser.Value {
	return expr.Value
}

func (i *Interpreter) VisitVariable(expr *parser.Variable) *parser.Value {
	value := i.environment.Get(expr.Name)
	if value.IsNil() {
		// Could raise an error here instead
		return parser.NewNilValue()
	}
	return value
}

// Implement StmtVisitor
func (i *Interpreter) VisitExpressionStmt(stmt *parser.ExpressionStmt) *parser.Value {
	return stmt.Expr.Accept(i)
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) *parser.Value {
	value := stmt.Expr.Accept(i)
	fmt.Println(value.String())
	return value
}

func (i *Interpreter) VisitVarStmt(stmt *parser.VarStmt) *parser.Value {
	var value *parser.Value
	if stmt.Expr != nil {
		value = stmt.Expr.Accept(i)
	} else {
		value = parser.NewNilValue()
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return value
}

// Helper methods for operations
func (i *Interpreter) evaluateBinaryOp(left, right *parser.Value, operator *ls.Token) *parser.Value {
	switch operator.Type {
	case ls.PLUS:
		return i.add(left, right)
	case ls.MINUS:
		return i.subtract(left, right)
	case ls.STAR:
		return i.multiply(left, right)
	case ls.SLASH:
		return i.divide(left, right)
	case ls.GREATER:
		return i.greater(left, right)
	case ls.GREATER_EQUAL:
		return i.greaterEqual(left, right)
	case ls.LESS:
		return i.less(left, right)
	case ls.LESS_EQUAL:
		return i.lessEqual(left, right)
	case ls.EQUAL_EQUAL:
		return i.equal(left, right)
	case ls.BANG_EQUAL:
		return i.notEqual(left, right)
	default:
		panic(fmt.Sprintf("Unknown binary operator: %s", operator.Lexeme))
	}
}

func (i *Interpreter) evaluateUnaryOp(right *parser.Value, operator *ls.Token) *parser.Value {
	switch operator.Type {
	case ls.MINUS:
		return i.negate(right)
	case ls.BANG:
		return i.logicalNot(right)
	default:
		panic(fmt.Sprintf("Unknown unary operator: %s", operator.Lexeme))
	}
}

// Operation implementations
func (i *Interpreter) add(left, right *parser.Value) *parser.Value {
	if left.IsNumber() && right.IsNumber() {
		return parser.NewFloatValue(left.ToFloat64() + right.ToFloat64())
	}
	if left.IsString() && right.IsString() {
		return parser.NewStringValue(*left.StrVal + *right.StrVal)
	}
	panic("Operands must be two numbers or two strings")
}

func (i *Interpreter) subtract(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewFloatValue(left.ToFloat64() - right.ToFloat64())
}

func (i *Interpreter) multiply(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewFloatValue(left.ToFloat64() * right.ToFloat64())
}

func (i *Interpreter) divide(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewFloatValue(left.ToFloat64() / right.ToFloat64())
}

func (i *Interpreter) greater(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewBoolValue(left.ToFloat64() > right.ToFloat64())
}

func (i *Interpreter) greaterEqual(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewBoolValue(left.ToFloat64() >= right.ToFloat64())
}

func (i *Interpreter) less(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewBoolValue(left.ToFloat64() < right.ToFloat64())
}

func (i *Interpreter) lessEqual(left, right *parser.Value) *parser.Value {
	i.checkNumberOperands(left, right)
	return parser.NewBoolValue(left.ToFloat64() <= right.ToFloat64())
}

func (i *Interpreter) equal(left, right *parser.Value) *parser.Value {
	if left.IsNil() && right.IsNil() {
		return parser.NewBoolValue(true)
	}
	if left.IsNil() {
		return parser.NewBoolValue(false)
	}
	return parser.NewBoolValue(left.String() == right.String())
}

func (i *Interpreter) notEqual(left, right *parser.Value) *parser.Value {
	return parser.NewBoolValue(!i.equal(left, right).IsTruthy())
}

func (i *Interpreter) negate(value *parser.Value) *parser.Value {
	i.checkNumberOperand(value)
	return parser.NewFloatValue(-value.ToFloat64())
}

func (i *Interpreter) logicalNot(value *parser.Value) *parser.Value {
	return parser.NewBoolValue(!value.IsTruthy())
}

// Helper methods
func (i *Interpreter) checkNumberOperand(value *parser.Value) {
	if !value.IsNumber() {
		panic("Operand must be a number")
	}
}

func (i *Interpreter) checkNumberOperands(left, right *parser.Value) {
	if !left.IsNumber() || !right.IsNumber() {
		panic("Operands must be numbers")
	}
}

// Main interpret method
func (i *Interpreter) Interpret(statements []parser.Stmt) {
	for _, stmt := range statements {
		result := stmt.Accept(i)
		if i.mode == ModePrompt && stmt.Type() == parser.EXPRESSION_STMT {
			if result == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(result.String())
			}
		}
	}
}
