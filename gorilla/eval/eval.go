package eval

import (
	"../ast"
	"../config"
	"../object"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var TRUE = object.TRUE
var FALSE = object.FALSE
var NULL = object.NULL

var currentRec = 0

func FromNativeBoolean(input bool, l int) *object.Boolean {
	if input {
		x := TRUE
		x.SLine = l
		return x
	}
	x := FALSE
	x.SLine = l
	return x
}

func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment, out ...io.Writer) object.Object {
	if len(out) > 0 {
		config.SetOut(out[0])
	}

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return object.NewInt(node.Value, node.Token.Line)

	case *ast.Boolean:
		return FromNativeBoolean(node.Value, node.Token.Line)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val, SLine: node.Token.Line}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.FunctionStmt:
		params := node.Parameters
		body := node.Body
		fn := object.NewFunction(params, body, env, node.Token.Line)
		env.Set(node.Name, fn)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return object.NewFunction(params, body, env, node.Token.Line)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.GetAttr:
		return evalGetAttr(node, env)

	case *ast.StringLiteral:
		return object.NewString(node.Value, node.Token.Line)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return object.NewArray(elements, node.Token.Line)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalGetAttr(node *ast.GetAttr, env *object.Environment) object.Object {
	expr := Eval(node.Expr, env)
	if isError(expr) {
		return expr
	}

	attributes := expr.Attributes()
	obj := attributes[node.Name.String()]
	if obj == nil {
		return NewError(
			"[Line %d] Type '%s' does not have attribute '%s'",
			node.Token.Line+1,
			expr.Type(),
			node.Name.String(),
		)
	}

	obj.SetParent(expr)

	return obj
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if IsTruthy(condition) {
			result = Eval(we.Consequence, env)
		} else {
			break
		}
	}

	if result != nil {
		return result
	} else {
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return right
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right == TRUE {
		return FALSE
	}
	if right == FALSE {
		return TRUE
	}
	if right.Type() != object.INTEGER {
		return NewError("[Line %d] cannot negate type '%s' (When attempting to run '-%s')",
			right.Line()+1, right.Type(), right.Inspect())
	}

	value := right.(*object.Integer).Value
	return object.NewInt(-value, right.Line())
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if left == nil || right == nil {
		return NULL
	}

	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.INTEGER:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.ARRAY && operator == "<-":
		return left.(*object.Array).Push(right)

	case operator == "==":
		return FromNativeBoolean(left == right, left.Line())
	case operator == "!=":
		return FromNativeBoolean(left != right, left.Line())
	case left.Type() != right.Type():
		return NewError("[Line %d] type mismatch: %s %s %s (When attempting to run '%s %s %s')",
			left.Line()+1, left.Type(), operator, right.Type(), left.Inspect(), operator, right.Inspect())
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return object.NewInt(leftVal+rightVal, left.Line())
	case "-":
		return object.NewInt(leftVal-rightVal, left.Line())
	case "*":
		return object.NewInt(leftVal*rightVal, left.Line())
	case "/":
		if rightVal == 0 {
			return NewError("[Line %d] Division by Zero", right.Line()+1)
		}
		return object.NewInt(leftVal/rightVal, left.Line())
	case "%":
		if rightVal == 0 {
			return NewError("[Line %d] Modulo by Zero", right.Line()+1)
		}
		return object.NewInt(leftVal%rightVal, left.Line())
	case "<":
		return FromNativeBoolean(leftVal < rightVal, left.Line())
	case ">":
		return FromNativeBoolean(leftVal > rightVal, left.Line())
	case "<=":
		return FromNativeBoolean(leftVal <= rightVal, left.Line())
	case ">=":
		return FromNativeBoolean(leftVal >= rightVal, left.Line())
	case "==":
		return FromNativeBoolean(leftVal == rightVal, left.Line())
	case "!=":
		return FromNativeBoolean(leftVal != rightVal, left.Line())
	default:
		return NULL
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch operator {
	case "+":
		if right.Type() != "STRING" {
			break
		}
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		if len(leftVal)+len(rightVal) >= config.MAXSTRINGSIZE {
			return NewError("[Line %d] String overflow", left.Line()+1)
		}
		return object.NewString(leftVal+rightVal, left.Line())

	case "*":
		if right.Type() != "INTEGER" {
			break
		}
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.Integer).Value
		if rightVal < 0 {
			return NULL
		}
		if len(leftVal)*int(rightVal) >= config.MAXSTRINGSIZE {
			return NewError("[Line %d] String overflow", left.Line()+1)
		}
		return object.NewString(strings.Repeat(leftVal, int(rightVal)), left.Line())

	case "==":
		if right.Type() != "STRING" {
			return FALSE
		}
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return FromNativeBoolean(leftVal == rightVal, left.Line())

	case "!=":
		if right.Type() != "String" {
			return TRUE
		}
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return FromNativeBoolean(leftVal != rightVal, left.Line())

	default:
		break
	}

	return NewError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	default:
		return NewError("[Line %d] Cannot perform index operation: %s[%s]", left.Line()+1, left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Value) - 1)

	if idx < 0 || idx > max {
		return NewError("[Line %d] Array index out of range", arrayObject.Line()+1)
	}

	return arrayObject.Value[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(utf8.RuneCountInString(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NewError("[Line %d] String index out of range", stringObject.Line()+1)
	}

	retString := object.NewString(string([]rune(stringObject.Value)[idx]), stringObject.Line())
	return retString // stringObject.Value[idx]
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if IsTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}

	return NewError("[Line %d] Variable '%s' is not defined", node.Token.Line+1, node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func IsTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return NewError("[Line %d] Argument mismatch (expected %d, got %d)", fn.Line()+1,
				len(fn.Parameters), len(args))
		}
		if currentRec >= config.RecursionLimit {
			currentRec = 0
			return NewError("[Line %d] Max recursion limit hit", fn.Line()+1)
		}
		currentRec++
		env := extendFunctionEnv(fn, args)
		res := Eval(fn.Body, env)
		if res == nil {
			res = NULL
		}
		currentRec--
		return unwrapReturnValue(res)

	case *object.Builtin:
		return fn.Fn(fn.Parent(), fn.Line(), args...)

	default:
		return NewError("[Line %d] Type '%s' is not callable", fn.Line()+1, fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}
	return obj
}
