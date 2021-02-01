package object

const (
	_ byte = iota
	INTMESSAGE
	STRINGMESSAGE
)

type Message interface {
	Type() string
}
