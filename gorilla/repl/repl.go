package repl

import (
	"bufio"
	"fmt"
	"io"

	"../compiler"
	"../config"
	"../eval"
	"../lexer"
	"../object"
	"../parser"
	"../vm"
)

// PROMPT is the REPL prompt displayed for each input
const PROMPT = "» "
const PROMPT2 = "* "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	_, _ = io.WriteString(out, "Gorilla "+config.VERSION+"\n")
	i := 0
	status := 0
	txt := ""
	for {
		i++

		prompt := PROMPT
		if status == 1 {
			prompt = PROMPT2
		}
		_, _ = io.WriteString(out, fmt.Sprintf("[%d]%s", i, prompt))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		txt += line + "\n"
		l := lexer.New(txt)
		p := parser.New(l)

		p.ParseProgram()

		if status == 0 {
			if len(p.Errors()) != 0 {
				status = 1
				//PrintParserErrors(out, p.Errors())
				continue
			}
		} else {
			if line == "" {

			} else if len(p.Errors()) != 0 {
				//status = 1
				//PrintParserErrors(out, p.Errors())
				continue
			}
		}

		status = 0

		l = lexer.New(txt)
		txt = ""
		p = parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}
		obj := eval.Eval(program, env)
		if obj != nil {
			if obj == object.NULL {
				continue
			}
			_, _ = io.WriteString(out, obj.Inspect()+"\n")
		}
	}
}

func StartCompile(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	var constants []object.Object
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	_, _ = io.WriteString(out, "Gorilla "+config.VERSION+"\n")
	i := 0
	status := 0
	txt := ""
	for {
		i++

		prompt := PROMPT
		if status == 1 {
			prompt = PROMPT2
		}
		_, _ = io.WriteString(out, fmt.Sprintf("[%d]%s", i, prompt))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		txt += line + "\n"
		l := lexer.New(txt)
		p := parser.New(l)

		p.ParseProgram()

		if status == 0 {
			if len(p.Errors()) != 0 {
				status = 1
				continue
			}
		} else {
			if line == "" {

			} else if len(p.Errors()) != 0 {
				continue
			}
		}

		status = 0

		l = lexer.New(txt)
		txt = ""
		p = parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			_, _ = io.WriteString(out, fmt.Sprintf("Compilation failed:\n\t%s\n", err))
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)

		err = machine.Run()
		if err != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Runtime Error:\n\t%s\n", err))
			continue
		}

		stackTop := machine.LastPopped()
		if stackTop == object.NULL || stackTop == nil {
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
