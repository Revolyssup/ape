// This package consists ape's bytecode format.Ape uses big-endian

package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandWidth := len(def.OperandWidths)
	if len(operands) != operandWidth {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandWidth)
	}
	switch operandWidth {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandWidth for %s\n", def.Name)
}
func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		def, err := LookupOpcode(Opcode(ins[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, n := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + n
	}
	return out.String()
}

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
	Opconstant: {"OpConstant", []int{2}}, //The single operand takes 2 bytes(16 bits) which means we can have 65536 unique constants in our constant pool at a time.
}

func LookupOpcode(op Opcode) (*Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("invalid opcode of type %v", op)
	}
	return def, nil
}

func MakeByteCodeFromOpcodeAndOperands(op Opcode, operands ...int) []byte {
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
		case 2: //If we have an operand of 16 bits, then we convert that to bigendian 8-8 bits
			binary.BigEndian.PutUint16(instruction[offset:], uint16(opr))
		}
		offset += width
	}
	return instruction
}

//Opposite of Make. It takes bytecode and spits out operands it read
func ReadOperands(def *Definition, ins Instructions) ([]int, int) { //Decode operands from bytecode instructions
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
