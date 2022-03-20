package vm

import (
	"fmt"
	"testing"

	"github.com/Revolyssup/ape/ast"
	"github.com/Revolyssup/ape/compiler"
	"github.com/Revolyssup/ape/lexer"
	"github.com/Revolyssup/ape/obj"
	"github.com/Revolyssup/ape/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
func testIntegerObject(expected int64, actual obj.Object) error {
	result, ok := actual.(*obj.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}
	return nil
}
func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2}, // FIXME
	}
	runVmTests(t, tests)
}
func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := New(comp.ByteCode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()
		testExpectedObject(t, tt.expected, stackElem)

	}
}
func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual obj.Object,
) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}
