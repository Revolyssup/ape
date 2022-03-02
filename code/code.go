// This package consists ape's bytecode format.Ape uses big-endian

package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte

const (
	/*
		Though this is a stack based virtual machine, we will not have PUSH operation as that limits up to pushing literal values.
		Instead we will use Opconstant which refers to constant expressions whose value can be determined at compiled time and put in a "constants pool"
	*/
	Opconstant Opcode = iota
)

//For debugging purposes
type Definition struct {
	Name          string
	OperandWidths []int //OperandWidths contains the number of bytes each operand takes up
}

var definitions = map[Opcode]*Definition{
	Opconstant: {"Opconstant", []int{2}},
}

func LookupOpcode(op Opcode) (*Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("invalid opcode of type %v", op)
	}
	return def, nil
}

func MakeByteCodeFromOpcode(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1 //for opcode
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	offset := 1 //first byte is for opcode
	instruction[0] = byte(op)
	for i, opr := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(opr))
		}
		offset += width
	}
	return instruction
}
