package repl

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/compiler"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/SnowballSH/Gorilla/runtime"
	"os"
	"strings"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Gorilla repl. Type :quit to quit.")

	env := runtime.NewEnvironment()

	for {
		fmt.Print("> ")
		text, e := reader.ReadString('\n')
		if e != nil {
			return
		}
		text = strings.TrimSpace(text)
		if text == ":quit" {
			return
		}

		par := parser.NewParser(parser.NewLexer(text))
		res := par.Parse()
		if par.Error != nil {
			fmt.Println("Syntax Error:\n" + *par.Error)
			continue
		}

		comp := compiler.NewCompiler()
		comp.Compile(res)

		vm := runtime.NewVMWithStore(comp.Result, env)
		vm.Run()
		if vm.Error != nil {
			fmt.Println(
				fmt.Sprintf("Runtime Error in line %d:\n\n| %s\n%s",
					vm.Error.Line+1,
					strings.Split(strings.ReplaceAll(text, "\r", ""), "\n")[vm.Error.Line], vm.Error.Message),
			)
			continue
		}
		env = vm.Environment
		fmt.Println(fmt.Sprintf("#=> %s", vm.LastPopped.Inspect()))
	}
}
