package gorilla

import (
	"../compiler"
	"../eval"
	"../lexer"
	"../object"
	"../parser"
	"../repl"
	"../vm"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func RunFile() {
	compile := flag.Bool("c", false, "Use the compiler")
	flag.Parse()

	if flag.NArg() < 1 {
		var fn func(in io.Reader, out io.Writer)
		if *compile {
			fn = repl.StartCompile
		} else {
			fn = repl.Start
		}
		fn(os.Stdin, os.Stdout)
		os.Exit(0)
	}

	fn := flag.Arg(0)

	b, err := ioutil.ReadFile(fn) // just pass the file name
	if err != nil {
		_, _ = io.WriteString(os.Stdout, "Error while reading file:\n"+err.Error())
		os.Exit(1)
	}

	code := string(b)

	if *compile {
		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(os.Stdout, p.Errors())
			os.Exit(1)
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Compilation failed:\n\t%s\n", err))
			os.Exit(1)
		}

		code := comp.Bytecode()

		machine := vm.New(code)

		err = machine.Run()
		if err != nil {
			_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Runtime Error:\n\t%s\n", err))
			os.Exit(1)
		}

		os.Exit(0)
	}

	env := object.NewEnvironment().AddBuiltin()

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	obj := eval.Eval(program, env)

	if obj != nil && obj.Type() == "ERROR" {
		_, _ = io.WriteString(os.Stdout, obj.Inspect()+"\n")
	}

	os.Exit(0)
}
