package repl

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/parser"
	"os"
	"strings"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Gorilla repl. Type :quit to quit.")

	for {
		fmt.Print("> ")
		text, e := reader.ReadString('\n')
		if e != nil {
			return
		}
		if strings.TrimSpace(text) == ":quit" {
			return
		}

		par := parser.NewParser(parser.NewLexer(text))
		res := par.Parse()
		if par.Error != nil {
			fmt.Println(*par.Error)
			continue
		}
		for _, item := range res {
			fmt.Println(item.String())
		}
	}
}
