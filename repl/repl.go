package repl

import (
	"Gorilla/compiler"
	"Gorilla/lexer"
	"Gorilla/parser"
	"Gorilla/vm"
	"bufio"
	"fmt"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	_, _ = io.WriteString(out, "Gorilla \n")

	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text() + "\n"

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewBytecodeCompiler()
		err := comp.Compile(program)
		if err != nil {
			_, _ = io.WriteString(out, fmt.Sprintf("Compilation failed:\n\t%s\n", err))
			continue
		}

		code := comp.Bytecodes
		constants := comp.Constants
		messages := comp.Messages

		machine := vm.New(code, constants, messages)
		e := machine.Run()
		if e != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Runtime Error:\n\t%s\n", e.Inspect()))
			continue
		}

		stackTop := machine.LastPopped
		if stackTop == nil {
			continue
		}
		_, _ = io.WriteString(out, stackTop.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
