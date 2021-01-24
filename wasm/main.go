package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"

	"../gorilla/eval"
	"../gorilla/lexer"
	"../gorilla/object"
	"../gorilla/parser"
	"../gorilla/repl"
)

func run(code string) string {
	env := object.NewEnvironment()

	w := bytes.Buffer{}

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(&w, p.Errors())
		return w.String()
	}

	obj := eval.Eval(program, env, &w)

	if obj == nil {
		return w.String()
	}

	if obj.Type() == "ERROR" {
		_, _ = io.WriteString(&w, obj.Inspect()+"\n")
	}

	return w.String()
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
