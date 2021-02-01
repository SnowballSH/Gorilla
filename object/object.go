package object

import "fmt"

const (
	INTEGER = "Integer"
)

type BaseObject interface {
	Type() string
	Inspect() string
	Line() int
}

type Integer struct {
	Value int
	SLine int
}

func (*Integer) Type() string {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Line() int {
	return i.SLine
}
