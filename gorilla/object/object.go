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
	NULL = "NULL"

	// RETURN is the Return object type
	RETURN = "RETURN"

	// ERROR is the ERROR object
	ERROR = "ERROR"

	// FUNCTION is the Function object
	FUNCTION = "FUNCTION"
)

// Type represents the type of an object
type Type string

// Object represents a value and implementations are expected to implement
// `Type()` and `Inspect()` functions
type Object interface {
	Type() Type
	Inspect() string
	Line() int
}

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value int64
	SLine int
}

// Type returns the type of the object
func (i *Integer) Type() Type { return INTEGER }

// Inspect returns a stringified version of the object for debugging
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Line() int { return i.SLine }

// String
type String struct {
	Value string
	SLine int
}

func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return "\"" + s.Value + "\"" }
func (s *String) Line() int       { return s.SLine }

// Boolean is the boolean type and used to represent boolean literals and holds an interval bool value
type Boolean struct {
	Value bool
	SLine int
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Line() int { return b.SLine }

// Null is the null type and used to represent the absence of a value
type Null struct {
	SLine int
}

// Type returns the type of the object
func (n *Null) Type() Type { return NULL }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }

func (n *Null) Line() int { return n.SLine }

// Return is the return statement
type Return struct {
	Value Object
	SLine int
}

// Type returns the type of the object
func (rv *Return) Type() Type { return RETURN }

// Inspect returns a stringified version of the object for debugging
func (rv *Return) Inspect() string { return rv.Value.Inspect() }

func (rv *Return) Line() int { return rv.SLine }

// Error the the error object
type Error struct {
	Message string
	SLine   int
}

// Type returns the type of the object
func (e *Error) Type() Type { return ERROR }

// Inspect returns a stringified version of the object for debugging
func (e *Error) Inspect() string { return " Runtime Error:\n\t" + e.Message }

func (e *Error) Line() int { return e.SLine }

// Function is the base function object type
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	SLine      int
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

func (f *Function) Line() int { return f.SLine }
