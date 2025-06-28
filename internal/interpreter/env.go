package interpreter

import (
	"github.com/Piyush01Bhatt/interpreter_go/internal/parser"
)

type Env struct {
	values map[string]*parser.Value
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]*parser.Value),
	}
}

func (e *Env) Define(name string, value *parser.Value) {
	e.values[name] = value
}

func (e *Env) Get(name string) *parser.Value {
	return e.values[name]
}
