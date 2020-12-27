package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"../lexer"
	"../token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	i := 0
	for {
		fmt.Printf("Gorilla " + fmt.Sprintf("[%d]", i) + PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == ":quit" {
			fmt.Println("You quit!")
			break
		}

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
		i++
	}
}
