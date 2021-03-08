package compiler

import (
	"ekyu.moe/leb128"
	"github.com/SnowballSH/Gorilla/grammar"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

func Compile(nodes []ast.Statement) (res []byte, ok bool) {
	result := []byte{grammar.Magic}
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

func compileNode(channel chan []byte, node ast.Node) {
	switch v := node.(type) {
	case *ast.ExpressionStatement:
		compileExpr(channel, v.Es)
	}
}

func compileExpr(channel chan []byte, v ast.Expression) {
	switch e := v.(type) {
	case *ast.Integer:
		w := []byte{grammar.Integer}
		l := leb128.AppendSleb128(nil, e.Value)
		w = append(w, byte(len(l)))
		w = append(w, l...)
		channel <- w

	case *ast.GetVar:
		w := []byte{grammar.GetVar}
		l := []byte(e.Name)
		w = append(w, byte(len(l)))
		w = append(w, l...)
		channel <- w

	case *ast.SetVar:
		w := []byte{grammar.SetVar}
		l := []byte(e.Name)
		w = append(w, byte(len(l)))
		w = append(w, l...)

		ch := make(chan []byte)
		go func() {
			compileExpr(ch, e.Value)
		}()
		w = append(w, <-ch...)
		channel <- w
	}
}
