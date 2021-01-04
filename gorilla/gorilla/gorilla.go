package gorilla

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"../eval"
	"../lexer"
	"../object"
	"../parser"
	"../repl"
)

func RunFile() {
	fnptr := flag.String("run", "", "Run a file")
	flag.Parse()

	fn := *fnptr

	if fn == "" {
		repl.Start(os.Stdin, os.Stdout)
		os.Exit(0)
	}

	b, err := ioutil.ReadFile(fn) // just pass the file name
	if err != nil {
		_, _ = io.WriteString(os.Stdout, "Error while reading file:\n"+err.Error())
		os.Exit(1)
	}

	code := string(b)

	env := object.NewEnvironment()

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	obj := eval.Eval(program, env)
	if obj != nil {
		_, _ = io.WriteString(os.Stdout, obj.Inspect()+"\n")
	}

	os.Exit(0)
}
