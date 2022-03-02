package obj

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Revolyssup/ape/ast"
)

type DataType string

const (
	INTEGER_OBJ      = "Integer"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "Bool"
	NULL_OBJ         = "Null"
	RETURN_OBJ       = "Return"
	ERROR_OBJ        = "Error"
	FUNCTION_OBJ     = "Function"
	BUILTIN_FUNC_OBJ = "Builtin_function"
	ARRAYS_OBJ       = "Array"
	OBJECT_OBJ       = "Object"
)

//All variables will be wrapped inside of an object-like struct.

type Object interface {
	DataType() DataType
	Inspect() string
}

//Implementing Integers

type Integer struct {
	Value int64
}

func (integer *Integer) DataType() DataType {
	return INTEGER_OBJ
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

//Implementing String
type String struct {
	Value string
}

func (s *String) DataType() DataType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

//Implementing Booleans
type Boolean struct {
	Value bool
}

func (boolean *Boolean) DataType() DataType {
	return BOOLEAN_OBJ
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

//Implementing Null
type Null struct{} //it holds no value
func (null *Null) DataType() DataType {
	return NULL_OBJ
}

func (null *Null) Inspect() string {
	return "null"
}

//Implementing  Return
type Return struct {
	Value Object
}

func (ret *Return) DataType() DataType {
	return RETURN_OBJ
}

func (ret *Return) Inspect() string {
	return ret.Value.Inspect()
}

//Implementing Error object is similar to Return as they both stop the execution of program and return something
type Error struct {
	ErrMsg string
}

func (err *Error) DataType() DataType {
	return ERROR_OBJ
}

func (err *Error) Inspect() string {
	return "[MONKE ANGRY:] " + err.ErrMsg
}

//Environment object will passed around recursively in Eval

type Env struct {
	variables map[string]Object
	outer     *Env
}

func (env *Env) Get(s string) (Object, bool) {
	ob, ok := env.variables[s]
	return ob, ok
}

func (env *Env) Set(s string, ob Object) Object {
	env.variables[s] = ob
	return ob
}

func NewEnvironment() *Env {
	s := make(map[string]Object)
	env := &Env{variables: s, outer: nil}
	return env
}

//This function will populate outer environments of function's environment object
func NewEnclosedEnvironment(outer_env *Env) *Env {
	env := NewEnvironment()
	env.outer = outer_env
	return env
}

/*****************/
//FUNCTIONS
type Function struct {
	Args []*ast.Identifier
	Body *ast.BlockStatement
	Env  *Env
}

func (fn *Function) DataType() DataType {
	return FUNCTION_OBJ
}

func (fn *Function) Inspect() string { //Returns all params
	var out bytes.Buffer

	params := []string{}

	for _, p := range fn.Args {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.String())
	out.WriteString("}")

	return out.String()
}

/***********/
//Builtin Functions
type BuiltinFn func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) DataType() DataType {
	return BUILTIN_FUNC_OBJ
}
func (b *Builtin) Inspect() string {
	return "monkey in-built function"
}

/**************/
//Array- Different data types can be added to array.
type Array struct {
	Arr []Object
}

func (a *Array) DataType() DataType {
	return ARRAYS_OBJ
}
func (a *Array) Inspect() string {
	var out bytes.Buffer
	out.WriteString("[")
	for _, ele := range a.Arr {
		out.WriteString(ele.Inspect() + ",")

	}
	out.WriteString("]")
	return out.String()
}

/**************/
//Object
type Obj struct {
	OBJ map[string]Object
}

func (o *Obj) DataType() DataType {
	return OBJECT_OBJ
}
func (o *Obj) Inspect() string {
	var out bytes.Buffer
	out.WriteString("{")
	for key, val := range o.OBJ {
		out.WriteString(key + ":" + val.Inspect() + ",\n")

	}
	out.WriteString("}")
	return out.String()
}
