package vm

import (
	"fmt"

	"github.com/Revolyssup/ape/code"
	"github.com/Revolyssup/ape/compiler"
	"github.com/Revolyssup/ape/obj"
)

const StackSize = 2048

type VM struct {
	constants    []obj.Object
	instructions code.Instructions
	stackPointer int
	stack        []obj.Object //Always point to next free slot in the stack
}

func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		instructions: bytecode.Instruction,
		constants:    bytecode.Constants,
		stack:        make([]obj.Object, StackSize),
		stackPointer: 0,
	}
}

func (vm *VM) StackTop() obj.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.Opconstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			first, err := vm.pop()
			if err != nil {
				return err
			}
			second, err := vm.pop()
			if err != nil {
				return err
			}
			ans, err := addTwoObjects(first, second)
			if err != nil {
				return err
			}
			vm.push(ans)
		}
	}
	return nil
}

func addTwoObjects(obj1 obj.Object, obj2 obj.Object) (obj.Object, error) {
	if obj1.DataType() != obj2.DataType() {
		return nil, fmt.Errorf("Cannot add two different types %v and %v", obj1.DataType(), obj2.DataType())
	}
	switch obj1.DataType() {
	case obj.INTEGER_OBJ:
		a := obj1.(*obj.Integer)
		b := obj2.(*obj.Integer)
		return &obj.Integer{Value: a.Value + b.Value}, nil
	}
	return nil, fmt.Errorf("Invalid datatype")
}
func (vm *VM) pop() (obj.Object, error) {
	if vm.stackPointer < 0 {
		return nil, fmt.Errorf("Empty stack")
	}
	obj := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return obj, nil
}
func (vm *VM) push(obj obj.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("Stack overflow")
	}
	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++
	return nil
}
