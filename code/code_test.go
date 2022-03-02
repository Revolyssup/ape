package code

import "testing"

func TestMakeByteCodeFromOpcode(t *testing.T) {
	tests := []struct {
		op               Opcode
		operands         []int
		expectedByteCode []byte
	}{{Opconstant, []int{65534}, []byte{byte(Opconstant), 255, 254}}}
	for _, tt := range tests {
		instruction := MakeByteCodeFromOpcode(tt.op, tt.operands...)
		if len(instruction) != len(tt.expectedByteCode) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expectedByteCode), len(instruction))
		}
		for i, b := range tt.expectedByteCode {
			if instruction[i] != tt.expectedByteCode[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}
