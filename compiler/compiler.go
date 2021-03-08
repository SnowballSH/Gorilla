package compiler

import (
	"ekyu.moe/leb128"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

func Compile(nodes []ast.Statement) (res []byte, ok bool) {
	result := []byte{0x69}
	var channels []chan []byte
	for _, x := range nodes {
		channel := make(chan []byte)
		channels = append(channels, channel)
		go func(w ast.Statement) {
			compileNode(channel, w)
		}(x)
	}
	for _, y := range channels {
		val := <-y
		result = append(result, val...)
	}
	return result, true
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
