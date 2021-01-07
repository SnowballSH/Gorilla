package repl

import (
	"bufio"
	"fmt"
	"io"

	"../eval"
	"../lexer"
	"../object"
	"../parser"
)

// PROMPT is the REPL prompt displayed for each input
const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment().AddBuiltin()

	_, _ = io.WriteString(out, "Gorilla 0.1\n")
	i := 0
	for {
		_, _ = io.WriteString(out, fmt.Sprintf("[%d]%s", i, PROMPT))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		obj := eval.Eval(program, env)
		if obj != nil {
			_, _ = io.WriteString(out, obj.Inspect()+"\n")
		}
		i++
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
