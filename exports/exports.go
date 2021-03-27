package exports

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/compiler"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/SnowballSH/Gorilla/runtime"
	"strings"
)

// CompileGorilla returns an array of byte generated from your gorilla code
func CompileGorilla(code string) ([]byte, error) {
	p := parser.NewParser(parser.NewLexer(strings.TrimSpace(code)))
	res := p.Parse()
	if p.Error != nil {
		return nil, fmt.Errorf(*p.Error)
	}
	comp := compiler.NewCompiler()
	comp.Compile(res)

	return comp.Result, nil
}

// ExecuteGorillaBytecodeFromSource executes the bytecode and returns the last popped object (may be nil if no object is processed)
// and the first runtime error. This method accepts a source argument for better error reporting.
func ExecuteGorillaBytecodeFromSource(bytecode []byte, source string) (*runtime.VM, runtime.BaseObject, error) {
	env := runtime.NewEnvironment()

	return ExecuteGorillaBytecodeFromSourceAndEnv(bytecode, source, env)
}

// ExecuteGorillaBytecodeFromSource executes the bytecode and returns the last popped object (may be nil if no object is processed)
// and the first runtime error. This method accepts a source argument and a env argument for better more customization.
func ExecuteGorillaBytecodeFromSourceAndEnv(bytecode []byte, source string, env *runtime.Environment) (*runtime.VM, runtime.BaseObject, error) {
	vm := runtime.NewVMWithStore(bytecode, env)
	vm.Run()
	if vm.Error != nil {
		return vm, nil, fmt.Errorf(
			fmt.Sprintf("Runtime Error in line %d:\n\n| %s\n%s",
				vm.Error.Line+1,
				strings.Split(strings.ReplaceAll(source, "\r", ""), "\n")[vm.Error.Line], vm.Error.Message),
		)
	}

	return vm, vm.LastPopped, nil
}

// ExecuteGorillaBytecode executes the bytecode and returns the last popped object (may be nil if no object is processed)
// and the first runtime error.
func ExecuteGorillaBytecode(bytecode []byte) (*runtime.VM, runtime.BaseObject, error) {
	env := runtime.NewEnvironment()

	vm := runtime.NewVMWithStore(bytecode, env)
	vm.Run()
	if vm.Error != nil {
		return vm, nil, fmt.Errorf(
			fmt.Sprintf("Runtime Error in line %d:\n\n%s",
				vm.Error.Line+1,
				vm.Error.Message),
		)
	}

	return vm, vm.LastPopped, nil
}
