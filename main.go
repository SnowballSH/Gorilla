package main

import (
	"flag"
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/repl"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"

	"rsc.io/getopt"
)

func executeFile(fn string, compile bool, out string) {
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

	if !compile {
		_, _, err = exports.ExecuteGorillaBytecodeFromSource(bc, code)

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	} else {
		if strings.TrimSpace(out) != "" {
			err = ioutil.WriteFile(out, bc, 0644)
		} else {
			_, err = os.Stdout.Write(bc)
		}

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	}
}

func executeBytecode(fn string) {
	content, err := ioutil.ReadFile(fn)

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	ch := make(chan bool)

	go func() {
		defer func() {
			if a := recover(); a != nil {
				color.Red(`Invalid bytecode: Error while performing execution
If this bytecode is generated by Gorilla without modification, please create an issue on https://github.com/SnowballSH/Gorilla.
Error:
%s`, a)
				ch <- false
			}
		}()

		_, _, err = exports.ExecuteGorillaBytecode(content)

		if err != nil {
			color.Red(err.Error())
			ch <- false
			return
		}

		ch <- true
	}()

	if !<-ch {
		os.Exit(1)
	}
}

func main() {
	help := flag.Bool("help", false, "Shows help text (default: false)")
	file := flag.String("file", "", "Executes a .gr file (default: repl)")
	out := flag.String("out", "", "Output file of compiler (default: stdout)")
	compile := flag.Bool("compile", false, "Generates bytecode (default: false)")
	fromBytecode := flag.Bool("bytecode", false, "Execute from bytecode .grx file (default: false)")

	getopt.Aliases(
		"f", "file",
		"o", "out",
		"h", "help",
		"c", "compile",
		"b", "bytecode",
	)

	getopt.Parse()

	if *compile && *fromBytecode {
		color.Red("ERROR: --compile/-c conflicts with --bytecode/-b")
		os.Exit(1)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if strings.TrimSpace(*file) != "" {
		if *fromBytecode {
			executeBytecode(*file)
		} else {
			executeFile(*file, *compile, *out)
		}
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		repl.Start()
	} else {
		fn := flag.Arg(0)

		if *fromBytecode {
			executeBytecode(fn)
		} else {
			executeFile(fn, *compile, *out)
		}
	}
	os.Exit(0)
}
