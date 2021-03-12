package compiler

import (
	"ekyu.moe/leb128"
	"github.com/SnowballSH/Gorilla/grammar"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

func Compile(nodes []ast.Statement) (res []byte, ok bool) {
	result := []byte{grammar.Magic}
	lastLine := 0

	for _, x := range nodes {
		r, l := compileNode(x, lastLine)
		lastLine = l
		result = append(result, r...)
	}
	return result, true
}

func compileNode(node ast.Node, lastLine int) ([]byte, int) {
	var x []byte
	for node.Line() > lastLine {
		lastLine++
		x = append(x, grammar.Advance)
	}

	switch v := node.(type) {
	case *ast.ExpressionStatement:
		k, l := compileExpr(v.Es, lastLine)
		lastLine = l
		x = append(x, k...)
		x = append(x, grammar.Pop)
	}
	return x, lastLine
}

func compileExpr(v ast.Expression, lastLine int) ([]byte, int) {
	var w []byte

	for v.Line() > lastLine {
		lastLine++
		w = append(w, grammar.Advance)
	}

	switch e := v.(type) {
	case *ast.Integer:
		w = append(w, grammar.Integer)

		l := leb128.AppendSleb128(nil, e.Value)
		w = append(w, byte(len(l)))
		w = append(w, l...)

	case *ast.GetVar:
		w = append(w, grammar.GetVar)

		l := []byte(e.Name)
		w = append(w, byte(len(l)))
		w = append(w, l...)

	case *ast.SetVar:
		w = append(w, grammar.SetVar)

		l := []byte(e.Name)
		w = append(w, byte(len(l)))
		w = append(w, l...)

		x, o := compileExpr(e.Value, lastLine)
		lastLine = o
		w = append(w, x...)
	}

	return w, lastLine
}
