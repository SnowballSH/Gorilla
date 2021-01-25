package compiler

import (
	"fmt"
	"testing"

	"../ast"
	"../code"
	"../lexer"
	"../object"
	"../parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

func testConstants(
	_ *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}

		case []code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T",
					i, actual[i])
			}
			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s",
					i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Add),
				code.Make(code.Pop),
			},
		},

		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Sub),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Mul),
				code.Make(code.Pop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Div),
				code.Make(code.Pop),
			},
		},

		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.Pop),
				code.Make(code.LoadConst, 1),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.Pop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadFalse),
				code.Make(code.Pop),
			},
		},

		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Gt),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Gt),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Gteq),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 <= 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Gteq),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Eq),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Neq),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.LoadFalse),
				code.Make(code.Eq),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.LoadFalse),
				code.Make(code.Neq),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            if (true) { 10 }; 3333;
            `,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.LoadTrue),
				// 0001
				code.Make(code.JumpElse, 10),
				// 0004
				code.Make(code.LoadConst, 0),
				// 0007
				code.Make(code.Jump, 11),
				// 0010
				code.Make(code.LoadNull),
				// 0011
				code.Make(code.Pop),
				// 0012
				code.Make(code.LoadConst, 1),
				// 0015
				code.Make(code.Pop),
			},
		},
		{
			input: `
            if (true) { 10 } else { 20 }; 3333;
            `,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.LoadTrue),
				// 0001
				code.Make(code.JumpElse, 10),
				// 0004
				code.Make(code.LoadConst, 0),
				// 0007
				code.Make(code.Jump, 13),
				// 0010
				code.Make(code.LoadConst, 1),
				// 0013
				code.Make(code.Pop),
				// 0014
				code.Make(code.LoadConst, 2),
				// 0017
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
let one = 1
let two = 2
`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.SetGlobal, 1),
			},
		},
		{
			input: `
let one = 1
one
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
let one = 1
let two = one
two
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.SetGlobal, 1),
				code.Make(code.LoadGlobal, 1),
				code.Make(code.Pop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.LoadConst, 0),
					code.Make(code.LoadConst, 1),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 2, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `fn() { 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.LoadConst, 0),
					code.Make(code.LoadConst, 1),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 2, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `fn() { 1; 2 }`,
			expectedConstants: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.LoadConst, 0),
					code.Make(code.Pop),
					code.Make(code.LoadConst, 1),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 2, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionsWithoutRet(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.RetNull),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 0, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.Mul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.Sub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.Sub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.Sub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d",
			compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}
	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.Add)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.Add {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.Add)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != code.Mul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
			previous.Opcode, code.Mul)
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.LoadConst, 0), // The literal "24"
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 1, 0), // The compiled function
				code.Make(code.Call, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            let noArg = fn() { 24 };
            noArg();
            `,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.LoadConst, 0), // The literal "24"
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 1, 0), // The compiled function
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.Call, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            let oneArg = fn(a) { a };
            oneArg(24);
            `,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadLocal, 0),
					code.Make(code.Ret),
				},
				24,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 0, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Call, 1),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            let manyArg = fn(a, b, c) { a; b; c };
            manyArg(24, 25, 26);
            `,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadLocal, 0),
					code.Make(code.Pop),
					code.Make(code.LoadLocal, 1),
					code.Make(code.Pop),
					code.Make(code.LoadLocal, 2),
					code.Make(code.Ret),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 0, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.LoadConst, 2),
				code.Make(code.LoadConst, 3),
				code.Make(code.Call, 3),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            let num = 55;
            fn() { num }
            `,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.LoadGlobal, 0),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.Closure, 1, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            fn() {
                let num = 55;
                num
            }
            `,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.LoadConst, 0),
					code.Make(code.SetLocal, 0),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 1, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            fn() {
                let a = 55;
                let b = 77;
                a + b
            }
            `,
			expectedConstants: []interface{}{
				55,
				77,
				[]code.Instructions{
					code.Make(code.LoadConst, 0),
					code.Make(code.SetLocal, 0),
					code.Make(code.LoadConst, 1),
					code.Make(code.SetLocal, 1),
					code.Make(code.LoadLocal, 0),
					code.Make(code.LoadLocal, 1),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 2, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            fn(a) {
                fn(b) {
                    a + b
                }
            }
            `,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadFree, 0),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
				[]code.Instructions{
					code.Make(code.LoadLocal, 0),
					code.Make(code.Closure, 0, 1),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 1, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            fn(a) {
                fn(b) {
                    fn(c) {
                        a + b + c
                    }
                }
            };
            `,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadFree, 0),
					code.Make(code.LoadFree, 1),
					code.Make(code.Add),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
				[]code.Instructions{
					code.Make(code.LoadFree, 0),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Closure, 0, 2),
					code.Make(code.Ret),
				},
				[]code.Instructions{
					code.Make(code.LoadLocal, 0),
					code.Make(code.Closure, 1, 1),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.Closure, 2, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            let global = 55;
            fn() {
                let a = 66;
                fn() {
                    let b = 77;
                    fn() {
                        let c = 88;
                        global + a + b + c;
                    }
                }
            }
            `,
			expectedConstants: []interface{}{
				55,
				66,
				77,
				88,
				[]code.Instructions{
					code.Make(code.LoadConst, 3),
					code.Make(code.SetLocal, 0),
					code.Make(code.LoadGlobal, 0),
					code.Make(code.LoadFree, 0),
					code.Make(code.Add),
					code.Make(code.LoadFree, 1),
					code.Make(code.Add),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Add),
					code.Make(code.Ret),
				},
				[]code.Instructions{
					code.Make(code.LoadConst, 2),
					code.Make(code.SetLocal, 0),
					code.Make(code.LoadFree, 0),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Closure, 4, 2),
					code.Make(code.Ret),
				},
				[]code.Instructions{
					code.Make(code.LoadConst, 1),
					code.Make(code.SetLocal, 0),
					code.Make(code.LoadLocal, 0),
					code.Make(code.Closure, 5, 1),
					code.Make(code.Ret),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.SetGlobal, 0),
				code.Make(code.Closure, 6, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestAttr(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            "".strip
            `,
			expectedConstants: []interface{}{
				"",
				"strip",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.GetAttr),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.Array, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.LoadConst, 2),
				code.Make(code.Array, 3),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConst, 0),
				code.Make(code.LoadConst, 1),
				code.Make(code.Add),
				code.Make(code.LoadConst, 2),
				code.Make(code.LoadConst, 3),
				code.Make(code.Sub),
				code.Make(code.LoadConst, 4),
				code.Make(code.LoadConst, 5),
				code.Make(code.Mul),
				code.Make(code.Array, 3),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}