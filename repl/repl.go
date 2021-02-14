package repl

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/compiler"
	"github.com/SnowballSH/Gorilla/config"
	"github.com/SnowballSH/Gorilla/lexer"
	"github.com/SnowballSH/Gorilla/object"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/SnowballSH/Gorilla/vm"
	"io"
	"strings"
)

const (
	PROMPT1 = "#>> "
	PROMPT2 = "*>> "
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	_, _ = io.WriteString(out, fmt.Sprintf("Gorilla %s\n", config.VERSION))

	var env = object.NewEnvironment()

	var prompt = PROMPT1

	var line = 0

	var indent = 0

	for {
		text := ""

		for {
			line++

			if len(text) > 0 {
				ts := strings.TrimSpace(text)
				st := ts[len(ts)-1]
				if st == '{' ||
					st == '[' ||
					st == '(' ||
					st == ',' {
					indent++
				}
				if st == '}' ||
					st == ']' ||
					st == ')' {
					indent--
				}
			}

			if indent < 0 {
				indent = 0
			}
			_, _ = io.WriteString(out, fmt.Sprintf("%03d%s%s", line, prompt, strings.Repeat("   ", indent)))

			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line := scanner.Text()

			text += line

			if len(strings.TrimSpace(line)) == 0 {
				break
			}

			l := lexer.New(text)
			p := parser.New(l)

			p.ParseProgram()
			if len(p.Errors()) != 0 {
				prompt = PROMPT2
				continue
			}
			break
		}

		indent = 0

		prompt = PROMPT1

		l := lexer.New(text)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewBytecodeCompiler()
		err := comp.Compile(program)
		if err != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Compilation failed:\n\t%s\n", err))
			continue
		}

		code := comp.Bytecodes
		constants := comp.Constants
		messages := comp.Messages

		machine := vm.New(code, constants, messages)
		machine.Frame.Env = env
		e := machine.Run()
		if e != nil {
			_, _ = io.WriteString(out, fmt.Sprintf(" Runtime Error:\n\t%s\n", e.Inspect()))
			continue
		}

		env = machine.Frame.Env

		stackTop := machine.Frame.LastPopped
		if stackTop == nil || stackTop == object.NULLOBJ {
			continue
		}
		_, _ = io.WriteString(out, stackTop.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, " Parser Errors:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
