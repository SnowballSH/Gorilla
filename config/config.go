package config

import (
	"io"
	"os"
)

var OUT io.Writer = os.Stdout

func SetOut(out io.Writer) {
	OUT = out
}

const VERSION = "0.4.0"
