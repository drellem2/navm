package navm

import (
	"testing"
)

func init() {
}

func TestAddConstants(t *testing.T) {
	ir := IR{
		registersLength: 2,
		instructions: []Instruction{
			Instruction{
				op:   add,
				ret:  makeVirtualRegister(1),
				arg1: makeVirtualRegister(1),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
		},
		constants: []int{1, 2},
	}
	result := Interpret(&ir)
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestAddRegisters(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:   add,
				ret:  makeVirtualRegister(2),
				arg1: makeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   add,
				ret:  makeVirtualRegister(1),
				arg1: makeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
		},
		constants: []int{1, 2},
	}
	result := Interpret(&ir)
	if result != 4 {
		t.Errorf("Expected 4, got %d", result)
	}
}

func TestMove(t *testing.T) {
	ir := IR{
		registersLength: 4,
		instructions: []Instruction{
			Instruction{
				op:  mov,
				ret: makeVirtualRegister(3),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   add,
				ret:  makeVirtualRegister(2),
				arg1: makeVirtualRegister(3),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
			Instruction{
				op:  mov,
				ret: makeVirtualRegister(1),
				arg2: Arg{
					argType: virtualRegisterArg,
					value:   2,
				},
			},
		},
		constants: []int{1, 2},
	}
	result := Interpret(&ir)
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
}
