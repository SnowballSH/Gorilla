package main

import (
	"Gorilla/compiler"
	"Gorilla/lexer"
	"Gorilla/parser"
	"Gorilla/repl"
	"Gorilla/vm"
	"bytes"
	"fmt"
	"io"
	"syscall/js"
)

func run(txt string) string {
	writer := bytes.Buffer{}

	l := lexer.New(txt + "\n")
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(&writer, p.Errors())
		return writer.String()
	}

	comp := compiler.NewBytecodeCompiler()
	err := comp.Compile(program)
	if err != nil {
		_, _ = io.WriteString(&writer, fmt.Sprintf(" Compilation failed:\n\t%s\n", err))
		return writer.String()
	}

	code := comp.Bytecodes
	constants := comp.Constants
	messages := comp.Messages

	machine := vm.New(code, constants, messages)
	e := machine.Run()
	if e != nil {
		_, _ = io.WriteString(&writer, fmt.Sprintf(" Runtime Error:\n\t%s\n", e.Inspect()))
		return writer.String()
	}
	return writer.String()
}

func wrapper() js.Func {
	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid: no arguments"
		}
		code := args[0].String()
		fmt.Printf("input %s\n", code)

		res := run(code)
		return res
	})
	return jsFunc
}

func main() {
	fmt.Println("GO")
	js.Global().Set("runGorilla", wrapper())
	<-make(chan bool)
}
