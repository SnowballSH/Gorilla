package object

import (
	"strconv"
	"strings"
)

func MessagesFromString(text string) []Message {
	_commands := strings.Split(text, ";")
	var commands []string
	for _, _v := range _commands {
		if strings.TrimSpace(_v) == "" {
			continue
		}
		commands = append(commands, strings.TrimSpace(_v))
	}

	var messages []Message
	for _, command := range commands {
		val, err := strconv.Atoi(command)
		if err != nil {
			messages = append(messages, NewMessage(command))
		} else {
			messages = append(messages, NewMessage(val))
		}
	}

	return messages
}

func ConstantsFromString(text string) []BaseObject {
	_commands := strings.Split(text, ";")
	var commands []string
	for _, _v := range _commands {
		if strings.TrimSpace(_v) == "" {
			continue
		}
		commands = append(commands, strings.TrimSpace(_v))
	}

	var constants []BaseObject
	for _, command := range commands {
		command = strings.TrimSpace(command)
		val, err := strconv.Atoi(command)
		if err == nil {
			constants = append(constants, NewInteger(val, 0))
		} else {
			panic("Not supported: '" + command + "'")
		}
	}

	return constants
}
