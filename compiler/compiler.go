package compiler

import (
	"github.com/Revolyssup/ape/ast"
	"github.com/Revolyssup/ape/code"
	"github.com/Revolyssup/ape/obj"
)

type Compiler struct { //Grouping instructions and constant pool at any time during compilation by a single compiler instance
	instruction code.Instructions
	constants   []obj.Object
}
type ByteCode struct { //Will be extracted from compiler instance at the end of compilation mostly. This is what we will pass to VM
	Instruction code.Instructions
	Constants   []obj.Object
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instruction: c.instruction,
		Constants:   c.constants,
	}
}

func New() *Compiler {
	return &Compiler{
		instruction: code.Instructions{},
		constants:   []obj.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.Compile(node.LeftExpression)
		if err != nil {
			return err
		}
		err = c.Compile(node.RightExpression)
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		integer := &obj.Integer{Value: node.Value}
		c.constants = append(c.constants, integer)
		c.instruction = append(c.instruction, code.MakeByteCodeFromOpcodeAndOperands(code.Opconstant, len(c.constants)-1)...)
	}
	return nil
}
