package object

const (
	_ byte = iota
	INTMESSAGE
	STRINGMESSAGE
)

type Message interface {
	MessageType() byte
}

type IntMessage struct {
	Value int
}

func (*IntMessage) MessageType() byte {
	return INTMESSAGE
}

type StringMessage struct {
	Value string
}

func (*StringMessage) MessageType() byte {
	return STRINGMESSAGE
}
