package object

// Package object implements the object system (or value system) of Monkey
// used to both represent values as the evaluator encounters and constructs
// them as well as how the user interacts with values.

import (
	"bytes"
	"fmt"
	"strings"

	"../ast"
)

const (
	// INTEGER is the Integer object type
	INTEGER = "INTEGER"

	// STRING is the String object type
	STRING = "STRING"

	// BOOLEAN is the Boolean object type
	BOOLEAN = "BOOLEAN"

	// NULL is the Null object type
	NULLT = "NULL"

	// RETURN is the Return object type
	RETURN = "RETURN"

	// ERROR is the ERROR object
	ERROR = "ERROR"

	// FUNCTION is the Function object
	FUNCTION = "FUNCTION"

	// BUILTIN is the Builtin object type
	BUILTIN = "BUILTIN"
)

var (
	TRUE = &Boolean{Value: true}

	FALSE = &Boolean{Value: false}

	NULL = &Null{}
)

var noAttr = map[string]Object{}

// BuiltinFunction represents the builtin function type
type BuiltinFunction func(self Object, line int, args ...Object) Object

// Type represents the type of an object
type Type string

// Object represents a value and implementations are expected to implement
// `Type()` and `Inspect()` functions
type Object interface {
	Type() Type
	Inspect() string
	Line() int
	Attributes() map[string]Object
	Parent() Object
	SetParent(x Object)
}

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value   int64
	SLine   int
	SParent Object
	Attrs   map[string]Object
}

// Type returns the type of the object
func (i *Integer) Type() Type { return INTEGER }

// Inspect returns a stringified version of the object for debugging
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Line() int { return i.SLine }

func (i *Integer) Attributes() map[string]Object { return i.Attrs }
func (i *Integer) Parent() Object                { return i.SParent }
func (i *Integer) SetParent(x Object)            { i.SParent = x }

// String
type String struct {
	Value   string
	SLine   int
	Attrs   map[string]Object
	SParent Object
}

func (s *String) Type() Type                    { return STRING }
func (s *String) Inspect() string               { return s.Value }
func (s *String) Line() int                     { return s.SLine }
func (s *String) Attributes() map[string]Object { return s.Attrs }
func (s *String) Parent() Object                { return s.SParent }
func (s *String) SetParent(x Object)            { s.SParent = x }

// Boolean is the boolean type and used to represent boolean literals and holds an interval bool value
type Boolean struct {
	Value   bool
	SLine   int
	SParent Object
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Line() int                     { return b.SLine }
func (b *Boolean) Attributes() map[string]Object { return noAttr }
func (b *Boolean) Parent() Object                { return b.SParent }
func (b *Boolean) SetParent(x Object)            { b.SParent = x }

// Null is the null type and used to represent the absence of a value
type Null struct {
	SLine   int
	SParent Object
}

// Type returns the type of the object
func (n *Null) Type() Type { return NULLT }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }

func (n *Null) Line() int                     { return n.SLine }
func (n *Null) Attributes() map[string]Object { return noAttr }
func (n *Null) Parent() Object                { return n.SParent }
func (n *Null) SetParent(x Object)            { n.SParent = x }

// Return is the return statement
type Return struct {
	Value   Object
	SLine   int
	SParent Object
}

// Type returns the type of the object
func (rv *Return) Type() Type { return RETURN }

// Inspect returns a stringified version of the object for debugging
func (rv *Return) Inspect() string { return rv.Value.Inspect() }

func (rv *Return) Line() int                     { return rv.SLine }
func (rv *Return) Attributes() map[string]Object { return noAttr }
func (rv *Return) Parent() Object                { return rv.SParent }
func (rv *Return) SetParent(x Object)            { rv.SParent = x }

// Error the the error object
type Error struct {
	Message string
	SLine   int
	SParent Object
}

// Type returns the type of the object
func (e *Error) Type() Type { return ERROR }

// Inspect returns a stringified version of the object for debugging
func (e *Error) Inspect() string { return " Runtime Error:\n\t" + e.Message }

func (e *Error) Line() int                     { return e.SLine }
func (e *Error) Attributes() map[string]Object { return noAttr }
func (e *Error) Parent() Object                { return e.SParent }
func (e *Error) SetParent(x Object)            { e.SParent = x }

// Function is the base function object type
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	SLine      int
	SParent    Object
}

// Type returns the type of the object
func (f *Function) Type() Type { return FUNCTION }

// Inspect returns a stringified version of the object for debugging
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (f *Function) Line() int                     { return f.SLine }
func (f *Function) Attributes() map[string]Object { return noAttr }
func (f *Function) Parent() Object                { return f.SParent }
func (f *Function) SetParent(x Object)            { f.SParent = x }

type Builtin struct {
	Fn      BuiltinFunction
	SLine   int
	SParent Object
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns a stringified version of the object for debugging
func (b *Builtin) Inspect() string { return "Builtin Function" }

func (b *Builtin) Line() int                     { return b.SLine }
func (b *Builtin) Attributes() map[string]Object { return noAttr }
func (b *Builtin) Parent() Object                { return b.SParent }
func (b *Builtin) SetParent(x Object)            { b.SParent = x }
