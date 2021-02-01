package object

const (
	_ uint16 = iota
	INTEGER
)

type Object interface {
	Type() string
	Inspect() string
	Line() int
}
