package tests

import (
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/runtime"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

type test struct {
	filename string
	expected runtime.BaseObject
}

var tests = []test{
	{
		filename: "arth.gr",
		expected: runtime.NewInteger(-36),
	},
	{
		filename: "condition.gr",
		expected: runtime.GorillaFalse,
	},
}

func TestAll(t *testing.T) {
	for _, x := range tests {
		content, err := ioutil.ReadFile(x.filename)

		if err != nil {
			panic(err)
		}

		code := strings.TrimSpace(string(content))

		bc, err := exports.CompileGorilla(code)
		assert.Nil(t, err)

		_, res, err := exports.ExecuteGorillaBytecodeFromSource(bc, code)
		assert.Nil(t, err)

		assert.True(t, x.expected.EqualTo(res))
	}
}
