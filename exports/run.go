package exports

import (
	"bytes"
	"fmt"
	"github.com/SnowballSH/Gorilla/compiler"
	"github.com/SnowballSH/Gorilla/config"
	"github.com/SnowballSH/Gorilla/lexer"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/SnowballSH/Gorilla/repl"
	"github.com/SnowballSH/Gorilla/vm"
	"io"
)

func RunCodeFromString(code string) string {
	w := bytes.Buffer{}

	config.SetOut(&w)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(&w, p.Errors())
		return w.String()
	}

	comp := compiler.NewBytecodeCompiler()
	err := comp.Compile(program)
	if err != nil {
		_, _ = io.WriteString(&w, fmt.Sprintf(" Compilation failed:\n\t%s\n", err))
		return w.String()
	}

	code_ := comp.Bytecodes
	constants := comp.Constants
	messages := comp.Messages

	machine := vm.New(code_, constants, messages)
	e := machine.Run()
	if e != nil {
		_, _ = io.WriteString(&w, fmt.Sprintf(" Runtime Error:\n\t%s\n", e.Inspect()))
		return w.String()
	}

	return w.String()
}
