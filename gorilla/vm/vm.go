package vm

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"../code"
	"../compiler"
	"../config"
	"../eval"
	"../object"
)

const StackSize = 1 << 20
const GlobalSize = 1 << 16
const FrameSize = 1 << 20

type VM struct {
	constants []object.Object

	frames      []*Frame
	framesIndex int

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]

	globals []object.Object
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		frames:      frames,
		framesIndex: 1,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		frames:      frames,
		framesIndex: 1,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: globals,
	}
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {

		case code.Jump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1

		case code.JumpElse:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !eval.IsTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case code.SetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()

		case code.LoadGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.SetLocal:
			localIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			frame := vm.currentFrame()

			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()

		case code.LoadLocal:
			localIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case code.LoadBuiltin:
			builtinIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			definition := object.Builtins[builtinIndex]

			err := vm.push(definition.Builtin)
			if err != nil {
				return err
			}

		case code.LoadConst:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.LoadFree:
			freeIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			currentClosure := vm.currentFrame().cl
			err := vm.push(currentClosure.Free[freeIndex])
			if err != nil {
				return err
			}

		case code.Closure:
			constIndex := code.ReadUint16(ins[ip+1:])
			numFree := code.ReadUint16(ins[ip+3:])
			vm.currentFrame().ip += 4

			err := vm.pushClosure(int(constIndex), int(numFree))
			if err != nil {
				return err
			}

		case code.Add, code.Sub, code.Mul, code.Div, code.Mod, code.LARR:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.Eq, code.Neq, code.Gt, code.Gteq:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.Not:
			err := vm.executeNotOperator()
			if err != nil {
				return err
			}

		case code.Neg:
			err := vm.executeNegOperator()
			if err != nil {
				return err
			}

		case code.Pos:
			err := vm.executePosOperator()
			if err != nil {
				return err
			}

		case code.Call:
			numArgs := code.ReadUint16(ins[ip+1:])

			vm.currentFrame().ip += 2

			err := vm.executeCall(int(numArgs))

			if err != nil {
				return err
			}

		case code.GetAttr:
			err := vm.getAttr()

			if err != nil {
				return err
			}

		case code.Array:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			line := int(code.ReadUint16(ins[ip+3:]))
			vm.currentFrame().ip += 4

			array := vm.buildArray(vm.sp-numElements, vm.sp, line)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.Index:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.LoadTrue:
			err := vm.push(object.TRUE)
			if err != nil {
				return err
			}

		case code.LoadFalse:
			err := vm.push(object.FALSE)
			if err != nil {
				return err
			}

		case code.LoadNull:
			err := vm.push(object.NULL)
			if err != nil {
				return err
			}

		case code.Ret:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			if vm.sp < 0 {
				return fmt.Errorf("[Line %d] Return outside of a function", returnValue.Line())
			}

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.RetNull:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			if vm.sp < 0 {
				return fmt.Errorf("[Line Unknown] Return outside of a function")
			}

			err := vm.push(object.NULL)
			if err != nil {
				return err
			}

		case code.Pop:
			vm.pop()

		default:
			return fmt.Errorf("WTF not supported: %d", op)
		}
	}

	return nil
}

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]
	if callee == nil {
		return fmt.Errorf("NIL\n%s\n%s", vm.currentFrame().Instructions(), vm.stack[vm.sp-1-numArgs])
	}

	switch callee := callee.(type) {
	case *object.Closure:
		return vm.callClosure(callee, numArgs)
	case *object.Builtin:
		return vm.callBuiltin(callee, numArgs)
	default:
		return fmt.Errorf("[Line %d] Type '%s' is not callable", callee.Line()+1, callee.Type())
	}
}

func (vm *VM) buildArray(startIndex, endIndex, line int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return object.NewArray(elements, line)
}

func (vm *VM) getAttr() error {
	name := vm.pop().(*object.String).Value
	callee := vm.pop()
	callee = callee.(object.Object)

	ok := false
	var ele object.Object
	for n, v := range callee.Attributes() {
		if n == name {
			ok = true
			ele = v
			break
		}
	}

	if !ok || ele == nil {
		return fmt.Errorf(
			"[Line %d] Type '%s' does not have attribute '%s'",
			callee.Line()+1,
			callee.Type(),
			name)
	}

	ele.SetParent(callee)

	return vm.push(ele)
}

func (vm *VM) callClosure(cl *object.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("[Line %d] Argument mismatch (expected %d, got %d)", cl.Fn.Line()+1,
			cl.Fn.NumParameters, numArgs)
	}

	if cl.Fn == vm.currentFrame().cl.Fn {
		nextOp := vm.currentFrame().NextOp()
		if nextOp == code.RetNull {
			for p := 0; p < numArgs; p++ {
				vm.stack[vm.currentFrame().basePointer+p] = vm.stack[vm.sp-numArgs+p]
			}
			vm.sp -= numArgs + 1
			vm.currentFrame().ip = -1 // reset IP to beginning of the frame
			return nil
		}
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}

	vm.sp = vm.sp - numFree
	closure := &object.Closure{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) callBuiltin(builtin *object.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(builtin.Parent(), builtin.Line(), args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		return vm.push(result)
	} else {
		return vm.push(object.NULL)
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER && rightType == object.INTEGER {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	if leftType == object.STRING && rightType == object.STRING && op == code.Add {
		return vm.executeStringAddOperation(op, left, right)
	}
	if leftType == object.STRING && rightType == object.INTEGER && op == code.Mul {
		return vm.executeStringMulOperation(op, left, right)
	}
	if leftType == object.ARRAY && op == code.LARR {
		return vm.push(left.(*object.Array).Push(right))
	}

	return fmt.Errorf("[Line %d] Unsupported types for binary operation: %s, %s",
		left.Line()+1, leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.Add:
		result = leftValue + rightValue
	case code.Sub:
		result = leftValue - rightValue
	case code.Mul:
		result = leftValue * rightValue
	case code.Div:
		if rightValue == 0 {
			return fmt.Errorf("[Line %d] Division by Zero", right.Line()+1)
		}
		result = leftValue / rightValue
	case code.Mod:
		if rightValue == 0 {
			return fmt.Errorf("[Line %d] Modulo by Zero", right.Line()+1)
		}
		result = leftValue % rightValue
	default:
		return fmt.Errorf("[Line %d] Unknown integer operator: %d", left.Line()+1, op)
	}

	return vm.push(object.NewInt(result, left.(*object.Integer).Line()))
}

func (vm *VM) executeStringAddOperation(
	_ code.Opcode,
	left, right object.Object,
) error {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	if len(leftVal)+len(rightVal) >= config.MAXSTRINGSIZE {
		return fmt.Errorf("[Line %d] String overflow", left.Line()+1)
	}
	return vm.push(object.NewString(leftVal+rightVal, left.Line()))
}

func (vm *VM) executeStringMulOperation(
	_ code.Opcode,
	left, right object.Object,
) error {
	leftVal := left.(*object.String).Value
	rightVal := int(right.(*object.Integer).Value)
	if len(leftVal)*rightVal >= config.MAXSTRINGSIZE {
		return fmt.Errorf("[Line %d] String overflow", left.Line()+1)
	}
	return vm.push(object.NewString(strings.Repeat(leftVal, rightVal), left.Line()))
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return vm.executeIntegerComparison(op, left, right)
	}

	if left.Type() == object.STRING && right.Type() == object.STRING {
		return vm.executeStringComparison(op, left, right)
	}

	switch op {
	case code.Eq:
		return vm.push(eval.FromNativeBoolean(right == left, left.Line()))
	case code.Neq:
		return vm.push(eval.FromNativeBoolean(right != left, left.Line()))

	default:
		if left.Type() != right.Type() {
			return fmt.Errorf("[Line %d] type mismatch: %s, %s (When attempting to run operation one %s and %s)",
				left.Line()+1, left.Type(), right.Type(), left.Inspect(), right.Inspect())
		}

		return vm.push(object.NULL)
	}
}

func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.Eq:
		return vm.push(eval.FromNativeBoolean(rightValue == leftValue, left.Line()))
	case code.Neq:
		return vm.push(eval.FromNativeBoolean(rightValue != leftValue, left.Line()))
	case code.Gt:
		return vm.push(eval.FromNativeBoolean(rightValue < leftValue, left.Line()))
	case code.Gteq:
		return vm.push(eval.FromNativeBoolean(rightValue <= leftValue, left.Line()))
	default:
		return fmt.Errorf(
			"unknown operator: %s %d %s",
			left.Type(), op, right.Type(),
		)
	}
}

func (vm *VM) executeStringComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch op {
	case code.Eq:
		return vm.push(eval.FromNativeBoolean(rightValue == leftValue, left.Line()))
	case code.Neq:
		return vm.push(eval.FromNativeBoolean(rightValue != leftValue, left.Line()))
	default:
		return fmt.Errorf(
			"unknown operator: %s %d %s",
			left.Type(), op, right.Type(),
		)
	}
}

func (vm *VM) executeNotOperator() error {
	operand := vm.pop()

	switch operand {
	case object.TRUE:
		return vm.push(object.FALSE)
	case object.FALSE:
		return vm.push(object.TRUE)
	case object.NULL:
		return vm.push(object.TRUE)
	default:
		return vm.push(object.FALSE)
	}
}

func (vm *VM) executeNegOperator() error {
	right := vm.pop()

	if right == object.TRUE {
		return vm.push(object.FALSE)
	}
	if right == object.FALSE {
		return vm.push(object.TRUE)
	}
	if right.Type() != object.INTEGER {
		return fmt.Errorf("[Line %d] cannot negate type '%s' (When attempting to run '-%s')",
			right.Line()+1, right.Type(), right.Inspect())
	}

	value := right.(*object.Integer).Value
	return vm.push(object.NewInt(-value, right.Line()))
}

func (vm *VM) executePosOperator() error {
	operand := vm.pop()

	return vm.push(operand)
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return vm.executeStringIndex(left, index)
	default:
		return fmt.Errorf("[Line %d] Cannot perform index operation: %s[%s]", left.Line()+1, left.Type(), index.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Value) - 1)

	if i < 0 || i > max {
		return fmt.Errorf("[Line %d] Array index out of range", arrayObject.Line()+1)
	}

	return vm.push(arrayObject.Value[i])
}

func (vm *VM) executeStringIndex(stri, index object.Object) error {
	str := stri.(*object.String)
	i := index.(*object.Integer).Value
	max := int64(utf8.RuneCountInString(str.Value) - 1)

	if i < 0 || i > max {
		return fmt.Errorf("[Line %d] String index out of range", str.Line()+1)
	}

	return vm.push(object.NewString(string([]rune(str.Value)[i]), str.Line()))
}
