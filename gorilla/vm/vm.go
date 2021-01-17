package vm

import (
	"fmt"

	"../code"
	"../compiler"
	"../eval"
	"../object"
)

const StackSize = 1 << 14

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
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
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {

		case code.LoadConst:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

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

	return fmt.Errorf("[Line %d] Unsupported types for binary operation: %s, %s",
		left.Line(), leftType, rightType)
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
		result = leftValue / rightValue
	default:
		return fmt.Errorf("[Line %d] Unknown integer operator: %d", left.(*object.Integer).Line(), op)
	}

	return vm.push(eval.NewInt(result, left.(*object.Integer).Line()))
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
