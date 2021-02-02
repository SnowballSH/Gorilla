package main

import (
	"Gorilla/compiler"
	"Gorilla/object"
	"Gorilla/vm"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello, Gorilla!")

	dat, err := ioutil.ReadFile("./examples/bc1.gobc")
	if err != nil {
		panic(err)
	}
	txt := strings.Split(string(dat), "%")

	v, err := compiler.BytecodeFromString(txt[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ms := object.MessagesFromString(txt[1])
	cs := object.ConstantsFromString(txt[2])

	vm_ := vm.New(v, cs, ms)
	vm_.Run()
	fmt.Println(vm_.LastPopped.Inspect())
}
