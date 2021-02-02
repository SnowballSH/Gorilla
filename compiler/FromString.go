package compiler

import (
	"Gorilla/code"
	"fmt"
	"strings"
)

func BytecodeFromString(text string) ([]code.Opcode, error) {
	_commands := strings.Split(text, ";")
	var commands []string
	for _, _v := range _commands {
		if strings.TrimSpace(_v) == "" {
			continue
		}
		commands = append(commands, strings.TrimSpace(_v))
	}

	var ops []code.Opcode
	for _, command := range commands {
		v, ok := code.StringToCode[command]
		if !ok {
			return nil, fmt.Errorf("command not found: %s", command)
		}
		ops = append(ops, v)
	}

	return ops, nil
}
