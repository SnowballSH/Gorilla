package object

import "fmt"

const (
	INTEGER = "Integer"
)

// Every Gorilla Object implements this
type BaseObject interface {
	Type() string
	Inspect() string
	Line() int
	Value() interface{}
	FindMethod(name string) (BaseObject, error)
	SetMethod(name string, value BaseObject)
}

// Implements BaseObject
type Object struct {
	TT            string
	InternalValue interface{}
	InspectValue  func(self BaseObject) string
	SLine         int
	Methods       map[string]BaseObject
}

func (o *Object) Type() string {
	return o.TT
}

func (o *Object) Inspect() string {
	return o.InspectValue(o)
}

func (o *Object) Line() int {
	return o.SLine
}

func (o *Object) Value() interface{} {
	return o.InternalValue
}

func (o *Object) FindMethod(name string) (BaseObject, error) {
	v, ok := o.Methods[name]
	if !ok {
		return nil, fmt.Errorf("[Line %d] method not found: %s", o.Line(), name)
	}
	return v, nil
}

func (o *Object) SetMethod(name string, value BaseObject) {
	o.Methods[name] = value
}

func NewObject(
	TT string,
	InternalValue interface{},
	InspectValue func(self BaseObject) string,
	SLine int,
	Methods map[string]BaseObject,
) *Object {
	return &Object{
		TT:            TT,
		InternalValue: InternalValue,
		InspectValue:  InspectValue,
		SLine:         SLine,
		Methods:       Methods,
	}
}

func NewInteger(
	value int,
	line int,
) *Object {
	return NewObject(
		INTEGER,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("%d", self.Value().(int))
		},
		line,
		map[string]BaseObject{},
	)
}
