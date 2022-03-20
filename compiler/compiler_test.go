package compiler

import (
	"fmt"
	"testing"

	"github.com/Revolyssup/ape/ast"
	"github.com/Revolyssup/ape/code"
	"github.com/Revolyssup/ape/lexer"
	"github.com/Revolyssup/ape/obj"
	"github.com/Revolyssup/ape/parser"
)

type testCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{
			input:             "1+2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.MakeByteCodeFromOpcodeAndOperands(code.Opconstant, 0),
				code.MakeByteCodeFromOpcodeAndOperands(code.Opconstant, 1),
			},
		},
	}
	runTests(t, tests)
}
func runTests(t *testing.T, tests []testCase) {
	for _, tt := range tests {
		prog := parse(tt.input)
		c := New()
		err := c.Compile(prog)
		if err != nil {
			t.Fatal("err ", err)
		}
		bc := c.ByteCode()
		err = testInstructions(bc.Instruction, tt.expectedInstructions)
		if err != nil {
			t.Fatal("err ", err)
		}
		err = testConstants(tt.expectedConstants, bc.Constants)
		if err != nil {
			t.Fatal("err ", err)
		}
	}
}
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(actual code.Instructions, expected []code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted.String(), actual.String())
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}
	return nil
}
func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}
func testConstants(expected []interface{}, actual []obj.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(expected int64, actual obj.Object) error {
	actualint, ok := actual.(*obj.Integer)
	if !ok {
		return fmt.Errorf("expected integer got %+v", actualint)
	}
	if actualint.Value != expected {
		return fmt.Errorf("value mismatch. Expected %v got %v", expected, actualint.Value)
	}
	return nil
}
