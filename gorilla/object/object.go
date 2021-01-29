package object

// Package object implements the object system (or value system) of Monkey
// used to both represent values as the evaluator encounters and constructs
// them as well as how the user interacts with values.

import (
	"bytes"
	"fmt"
	"strings"

	"../ast"
	"../code"
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

	// CompiledFunction is the CompiledFunction object
	COMPILEDFUNCTION = "COMPILEDFUNCTION"

	// BUILTIN is the Builtin object type
	BUILTIN = "BUILTIN"

	// CLOSURE
	CLOSURE = "CLOSURE"

	// ARRAY
	ARRAY = "ARRAY"
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
	SetAttribute(name string, value Object) Object
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
func (i *Integer) SetAttribute(name string, value Object) Object {
	i.Attrs[name] = value
	return i
}
func (i *Integer) Parent() Object     { return i.SParent }
func (i *Integer) SetParent(x Object) { i.SParent = x }

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
func (s *String) SetAttribute(name string, value Object) Object {
	s.Attrs[name] = value
	return s
}
func (s *String) Parent() Object     { return s.SParent }
func (s *String) SetParent(x Object) { s.SParent = x }

func init() {

}

// Array
type Array struct {
	Value   []Object
	SLine   int
	Attrs   map[string]Object
	SParent Object
}

func (a *Array) Type() Type { return ARRAY }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Value {
		var v string
		switch e.(type) {
		case *String:
			v = "\"" + e.Inspect() + "\""
		default:
			v = e.Inspect()
		}
		elements = append(elements, v)
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Line() int                     { return a.SLine }
func (a *Array) Attributes() map[string]Object { return a.Attrs }
func (a *Array) SetAttribute(name string, value Object) Object {
	a.Attrs[name] = value
	return a
}
func (a *Array) Parent() Object     { return a.SParent }
func (a *Array) SetParent(x Object) { a.SParent = x }

func (a *Array) Push(x Object) Object {
	a.Value = append(a.Value, x)
	return a
}

func (a *Array) PushAll(x []Object) Object {
	a.Value = append(a.Value, x...)
	return a
}

func (a *Array) PopLast() Object {
	k := a.Value[len(a.Value)-1]
	a.Value = a.Value[:len(a.Value)-1]
	return k
}

func (a *Array) PopFirst() Object {
	k := a.Value[0]
	a.Value = a.Value[1:]
	return k
}

func (a *Array) SetIndex(i int, v Object) Object {
	a.Value[i] = v
	return a
}

// Boolean is the boolean type and used to represent boolean literals and holds an interval bool value
type Boolean struct {
	Value   bool
	SLine   int
	SParent Object
	Attrs   map[string]Object
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Line() int                     { return b.SLine }
func (b *Boolean) Attributes() map[string]Object { return b.Attrs }
func (b *Boolean) SetAttribute(name string, value Object) Object {
	b.Attrs[name] = value
	return b
}
func (b *Boolean) Parent() Object     { return b.SParent }
func (b *Boolean) SetParent(x Object) { b.SParent = x }

// Null is the null type and used to represent the absence of a value
type Null struct {
	SLine   int
	SParent Object
	Attrs   map[string]Object
}

// Type returns the type of the object
func (n *Null) Type() Type { return NULLT }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }

func (n *Null) Line() int                     { return n.SLine }
func (n *Null) Attributes() map[string]Object { return n.Attrs }
func (n *Null) SetAttribute(name string, value Object) Object {
	n.Attrs[name] = value
	return n
}
func (n *Null) Parent() Object     { return n.SParent }
func (n *Null) SetParent(x Object) { n.SParent = x }

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

func (rv *Return) Line() int                          { return rv.SLine }
func (rv *Return) Attributes() map[string]Object      { return noAttr }
func (rv *Return) SetAttribute(string, Object) Object { return rv }
func (rv *Return) Parent() Object                     { return rv.SParent }
func (rv *Return) SetParent(x Object)                 { rv.SParent = x }

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

func (e *Error) Line() int                          { return e.SLine }
func (e *Error) Attributes() map[string]Object      { return noAttr }
func (e *Error) SetAttribute(string, Object) Object { return e }
func (e *Error) Parent() Object                     { return e.SParent }
func (e *Error) SetParent(x Object)                 { e.SParent = x }

// Function is the base function object type
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Attrs      map[string]Object
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
func (f *Function) Attributes() map[string]Object { return f.Attrs }
func (f *Function) SetAttribute(name string, value Object) Object {
	f.Attrs[name] = value
	return f
}
func (f *Function) Parent() Object     { return f.SParent }
func (f *Function) SetParent(x Object) { f.SParent = x }

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
	SLine         int
	SParent       Object
	Attrs         map[string]Object
}

func (cf *CompiledFunction) Type() Type { return COMPILEDFUNCTION }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("Compiled Function [%p]", cf)
}

func (cf *CompiledFunction) Line() int                     { return cf.SLine }
func (cf *CompiledFunction) Attributes() map[string]Object { return cf.Attrs }
func (cf *CompiledFunction) SetAttribute(name string, value Object) Object {
	cf.Attrs[name] = value
	return cf
}
func (cf *CompiledFunction) Parent() Object     { return cf.SParent }
func (cf *CompiledFunction) SetParent(x Object) { cf.SParent = x }

type Builtin struct {
	Fn      BuiltinFunction
	SLine   int
	SParent Object
	Attrs   map[string]Object
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns a stringified version of the object for debugging
func (b *Builtin) Inspect() string { return "Builtin Function" }

func (b *Builtin) Line() int                     { return b.SLine }
func (b *Builtin) Attributes() map[string]Object { return b.Attrs }
func (b *Builtin) SetAttribute(name string, value Object) Object {
	b.Attrs[name] = value
	return b
}
func (b *Builtin) Parent() Object     { return b.SParent }
func (b *Builtin) SetParent(x Object) { b.SParent = x }

type Closure struct {
	Fn      *CompiledFunction
	Free    []Object
	SLine   int
	SParent Object
	Attrs   map[string]Object
}

func (c *Closure) Type() Type { return CLOSURE }

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure [%p]", c)
}

func (c *Closure) Line() int                     { return c.SLine }
func (c *Closure) Attributes() map[string]Object { return c.Attrs }
func (c *Closure) SetAttribute(name string, value Object) Object {
	c.Attrs[name] = value
	return c
}
func (c *Closure) Parent() Object     { return c.SParent }
func (c *Closure) SetParent(x Object) { c.SParent = x }
