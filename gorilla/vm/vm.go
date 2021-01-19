package vm

import (
	"fmt"
	"strings"

	"../code"
	"../compiler"
	"../eval"
	"../object"
)

const StackSize = 1 << 14
const GlobalSize = 1 << 16
const FrameSize = 1 << 12

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int // Always points to the next value

	globals []object.Object

	frames      []*Frame
	framesIndex int
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
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalSize),

		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, FrameSize)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		frames:      frames,
		framesIndex: 1,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: s,
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
			pos := int(code.ReadUint16(vm.currentFrame().Instructions()[ip+1:]))
			vm.currentFrame().ip = pos - 1

		case code.JumpElse:
			pos := int(code.ReadUint16(vm.currentFrame().Instructions()[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !eval.IsTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case code.SetGlobal:
			globalIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()

		case code.LoadGlobal:
			globalIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
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

		case code.LoadConst:
			constIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.Add, code.Sub, code.Mul, code.Div:
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
			f := vm.stack[vm.sp-1]
			fn, ok := f.(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("[Line %d] Type '%s' is not callable", f.Line()+1, f.Type())
			}
			frame := NewFrame(fn, vm.sp)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + fn.NumLocals

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

		}
	}

	return nil
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
	default:
		return fmt.Errorf("[Line %d] Unknown integer operator: %d", left.Line()+1, op)
	}

	return vm.push(eval.NewInt(result, left.(*object.Integer).Line()))
}

func (vm *VM) executeStringAddOperation(
	_ code.Opcode,
	left, right object.Object,
) error {
	return vm.push(eval.NewString(left.(*object.String).Value+right.(*object.String).Value, left.Line()))
}

func (vm *VM) executeStringMulOperation(
	_ code.Opcode,
	left, right object.Object,
) error {
	return vm.push(eval.NewString(strings.Repeat(left.(*object.String).Value, int(right.(*object.Integer).Value)), left.Line()))
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER || right.Type() == object.INTEGER {
		return vm.executeIntegerComparison(op, left, right)
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
		return fmt.Errorf("unknown operator: %d", op)
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
	return vm.push(eval.NewInt(-value, right.Line()))
}

func (vm *VM) executePosOperator() error {
	operand := vm.pop()

	return vm.push(operand)
}
