package eval

import (
	"../ast"
	"../object"
	"fmt"
	"strings"
)

var (
	TRUE = &object.Boolean{Value: true}

	FALSE = &object.Boolean{Value: false}

	NULL = &object.Null{}
)

func fromNativeBoolean(input bool, l int) *object.Boolean {
	if input {
		x := TRUE
		x.SLine = l
		return x
	}
	x := FALSE
	x.SLine = l
	return x
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value, SLine: node.Token.Line}

	case *ast.Boolean:
		return fromNativeBoolean(node.Value, node.Token.Line)

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

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body, SLine: node.Token.Line}

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

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
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
		return newError("[Line %d] cannot negate type '%s' (When attempting to run '-%s')",
			right.Line()+1, right.Type(), right.Inspect())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value, SLine: right.Line()}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.Return); ok {
			return returnValue.Value
		}
	}

	return result
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

	case operator == "==":
		return fromNativeBoolean(left == right, left.Line())
	case operator == "!=":
		return fromNativeBoolean(left != right, left.Line())
	case left.Type() != right.Type():
		return newError("[Line %d] type mismatch: %s %s %s (When attempting to run '%s %s %s')",
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
		return &object.Integer{Value: leftVal + rightVal, SLine: left.Line()}
	case "-":
		return &object.Integer{Value: leftVal - rightVal, SLine: left.Line()}
	case "*":
		return &object.Integer{Value: leftVal * rightVal, SLine: left.Line()}
	case "/":
		return &object.Integer{Value: leftVal / rightVal, SLine: left.Line()}
	case "<":
		return fromNativeBoolean(leftVal < rightVal, left.Line())
	case ">":
		return fromNativeBoolean(leftVal > rightVal, left.Line())
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal, left.Line())
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal, left.Line())
	case "==":
		return fromNativeBoolean(leftVal == rightVal, left.Line())
	case "!=":
		return fromNativeBoolean(leftVal != rightVal, left.Line())
	default:
		return NULL
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator == "+" {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	}
	if operator == "*" {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.Integer).Value
		if rightVal < 0 {
			return NULL
		}
		return &object.String{Value: strings.Repeat(leftVal, int(rightVal))}
	}

	return newError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
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
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("[Line %d] Variable '%s' is not defined", node.Token.Line+1, node.Value)
	}

	return val
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

func isTruthy(obj object.Object) bool {
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
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("[Line %d] Type '%s' is not callable", fn.Line(), fn.Type())
	}
	if len(function.Parameters) != len(args) {
		return newError("[Line %d] Argument mismatch (expected %d, got %d)", fn.Line(),
			len(function.Parameters), len(args))
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
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
