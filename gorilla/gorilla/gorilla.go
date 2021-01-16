package gorilla

import (
	"../eval"
	"../lexer"
	"../object"
	"../parser"
	"../repl"
	"flag"
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
