package repl

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/runtime"
	"github.com/c-bata/go-prompt"
	"strings"
)

func completer(d prompt.Document, env *runtime.Environment) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: ":quit", Description: "Quit the repl"},
	}
	for n, k := range env.Store {
		s = append(s, prompt.Suggest{
			Text:        n,
			Description: fmt.Sprintf("Scope Variable, type '%s'", k.Class().Name),
		})
	}
	for n, k := range runtime.Global.Store {
		w := n
		if k.Class() == runtime.GoFuncClass {
			w += "()"
		}
		s = append(s, prompt.Suggest{
			Text:        w,
			Description: fmt.Sprintf("Global Variable, type '%s'", k.Class().Name),
		})
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Start() {
	fmt.Println("Welcome to Gorilla repl. Type :quit to quit.")

	env := runtime.NewEnvironment()

	for {
		text := prompt.Input("> ", func(document prompt.Document) []prompt.Suggest {
			return completer(document, env)
		})
		text = strings.TrimSpace(text)
		if text == ":quit" {
			return
		}

		res, err := exports.CompileGorilla(text)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		vm, lastPopped, err := exports.ExecuteGorillaBytecodeFromSourceAndEnv(res, text, env)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		env = vm.Environment

		if lastPopped != nil {
			fmt.Println(fmt.Sprintf("#=> %s", vm.LastPopped.Inspect()))
		}
	}
}
