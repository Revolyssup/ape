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
		}
	}
	return nil
}

func (vm *VM) push(obj obj.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("Stack overflow")
	}
	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++
	return nil
}
