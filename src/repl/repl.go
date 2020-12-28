package repl

import (
	"bufio"
	"fmt"
	"io"

	"../lexer"
	"../parser"
)

// PROMPT is the REPL prompt displayed for each input
const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		_, _ = io.WriteString(out, program.String()+"\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
