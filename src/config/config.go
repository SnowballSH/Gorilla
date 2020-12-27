package config

import (
	"runtime"
	"strings"
)

var OSNEWLINES = map[string]string{
	"windows": "\r\n",
	"linux":   "\n",
	"js":      "\n",
	"darwin":  "\n",
	"android": "\n",
}

func GetOSNewline(sys string) string {
	if sys == "" || strings.ToLower(sys) == "default" {
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
