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
const PROMPT = "» "
const PROMPT2 = "* "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment().AddBuiltin()

	_, _ = io.WriteString(out, "Gorilla 0.3\n")
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
			_, _ = io.WriteString(out, obj.Inspect()+"\n")
		}
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
