package compiler

import (
	"ekyu.moe/leb128"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

type Compiler struct {
	nodes []ast.Statement

	channels []chan []byte

	Result []byte
}

func NewCompiler(nodes []ast.Statement) *Compiler {
	return &Compiler{
		nodes:    nodes,
		channels: []chan []byte{},
		Result:   []byte{0x69},
	}
}

func (c *Compiler) compile() (ok bool) {
	for _, x := range c.nodes {
		channel := make(chan []byte)
		c.channels = append(c.channels, channel)
		go func(w ast.Statement) {
			compileNode(channel, w)
		}(x)
	}
	for _, y := range c.channels {
		val := <-y
		c.Result = append(c.Result, val...)
	}
	return true
}

func compileNode(channel chan []byte, node ast.Statement) {
	switch v := node.(type) {
	case *ast.ExpressionStatement:
		switch e := v.Es.(type) {
		case *ast.Integer:
			w := []byte{0x00}
			l := leb128.AppendSleb128(nil, e.Value)
			w = append(w, byte(len(l)))
			w = append(w, l...)
			channel <- w
		}
	}
}
