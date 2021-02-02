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

func NewMessage(value interface{}) Message {
	switch value.(type) {
	case int:
		return &IntMessage{Value: value.(int)}
	case string:
		return &StringMessage{Value: value.(string)}
	default:
		panic("Message Type is not string or int")
	}
}
