package config

import (
	"io"
	"os"
	"runtime"
	"strings"
)

var OUT io.Writer = os.Stdout

func SetOut(w io.Writer) {
	OUT = w
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
