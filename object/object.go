package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/takeru-a/golang_interpreterlang/ast"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ = "STRING"
	BUILTIN_OBJ = "BUILTIN"
	NULLSTRING_OBJ = "NULLSTRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// 整数
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value ) }

// 真偽値
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value )}

// NULL
type NULL struct{}

func (n *NULL) Type() ObjectType { return NULL_OBJ }
func (n *NULL) Inspect() string { return "null" }

type NULLSTRING struct{}

func (ns *NULLSTRING) Type() ObjectType { return NULLSTRING_OBJ }
func (ns *NULLSTRING) Inspect() string { return "" }

// Return
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function
type Function struct {
	Parameters []*ast.Indetifier
	Body *ast.BlockStatement
	Env *Environment
}

func (fn *Function) Type() ObjectType { return FUNCTION_OBJ }
func (fn *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// 文字列
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

// 組み込み関数
type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Builtin) Inspect() string { return "builtin function" }
 