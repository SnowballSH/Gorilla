package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"../gorilla/eval"
	"../gorilla/lexer"
	"../gorilla/object"
	"../gorilla/parser"
	"../gorilla/repl"

	"github.com/rs/cors"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	u, ok := url.QueryUnescape(r.URL.String()[1:])
	if ok != nil {
		return
	}

	fmt.Println("Code: " + u)

	code := u

	env := object.NewEnvironment().AddBuiltin()

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(w, p.Errors())
		os.Exit(1)
	}

	k, _ := w.(io.Writer)
	eval.Eval(program, env, k)
}

func doNothing(_ http.ResponseWriter, _ *http.Request) {}

func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/favicon.ico", doNothing)
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":10000", handler))
}

func main() {
	handleRequests()
}
