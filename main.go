package main

import (
	"flag"
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/repl"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		repl.Start()
	} else {
		fn := flag.Arg(0)
		content, err := ioutil.ReadFile(fn)

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		code := strings.TrimSpace(string(content))

		bc, err := exports.CompileGorilla(code)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		_, _, err = exports.ExecuteGorillaBytecodeFromSource(bc, code)

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
	os.Exit(0)
}
