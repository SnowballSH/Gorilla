package config

import (
	"io"
	"os"
	"runtime"
	"strings"
)

const VERSION = "0.3"

var OUT io.Writer = os.Stdout

func SetOut(w io.Writer) {
	OUT = w
}

var RecursionLimit = 1 << 11

func SetRecLimit(i int) {
	RecursionLimit = i
}

var OSNEWLINES = map[string]string{
	"windows": "\r\n",
	"linux":   "\n",
	"js":      "\n",
	"darwin":  "\n",
	"android": "\n",
}

func GetOSNewline(sys string) string {
	sys = strings.ToLower(sys)
	if sys == "" || sys == "default" {
		sys = runtime.GOOS
	}

	var nl string
	if val, ok := OSNEWLINES[sys]; ok {
		nl = val
	} else {
		nl = "\n"
	}
	return nl
}

const MAXSTRINGSIZE = 1 << 12

type Void struct{}
