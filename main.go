package main

import (
	"Gorilla/compiler"
	"Gorilla/lexer"
	"Gorilla/parser"
	"Gorilla/repl"
	"Gorilla/vm"
	"fmt"
	"github.com/alecthomas/kong"
	"io"
	"io/ioutil"
	"os"
)

var CLI struct {
	Run struct {
		Path string `arg name:"path" help:"Paths to remove." type:"path"`
	} `cmd help:"Run file"`
	Repl struct {
	} `cmd help:"Start Gorilla Repl"`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "repl":
		repl.Start(os.Stdin, ctx.Stdout)
	case "run <path>":
		out := ctx.Stdout

		txt, err := ioutil.ReadFile(CLI.Run.Path)
		if err != nil {
			_, _ = fmt.Fprintln(ctx.Stderr, err.Error())
			os.Exit(1)
		}

		l := lexer.New(string(txt))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(out, p.Errors())
			os.Exit(1)
		}

		comp := compiler.NewBytecodeCompiler()
		err = comp.Compile(program)
		if err != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Compilation failed:\n\t%s\n", err))
			os.Exit(1)
		}

		code := comp.Bytecodes
		constants := comp.Constants
		messages := comp.Messages

		machine := vm.New(code, constants, messages)
		e := machine.Run()
		if e != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Runtime Error:\n\t%s\n", e.Inspect()))
			os.Exit(1)
		}
	}
}
