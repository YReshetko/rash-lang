package objects

import (
	"bytes"
	"fmt"
	"github.com/YReshetko/monkey-language/ast"
	"strings"
)

type ObjectType string
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	STRING_OBJ       ObjectType = "STRING"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}

	NULL = &Null{}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

type Error struct {
	Message string
	Stack   []string
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) AddStackLine(line string) {
	e.Stack = append([]string{line}, e.Stack...)
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	out := bytes.Buffer{}

	params := []string{}
	for _, parameter := range f.Parameters {
		params = append(params, parameter.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}
