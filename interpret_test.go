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
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(1),
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
				ret:  MakeVirtualRegister(2),
				arg1: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   add,
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(2),
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
				ret: MakeVirtualRegister(3),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   add,
				ret:  MakeVirtualRegister(2),
				arg1: MakeVirtualRegister(3),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
			Instruction{
				op:  mov,
				ret: MakeVirtualRegister(1),
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

func TestSubRegisters(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  mov,
				ret: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   sub,
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
		},
		constants: []int{1, 2},
	}
	result := Interpret(&ir)
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestMultRegisters(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  mov,
				ret: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   mult,
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
		},
		constants: []int{3, 2},
	}
	result := Interpret(&ir)
	if result != 6 {
		t.Errorf("Expected 6, got %d", result)
	}
}

func TestDivRegisters(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  mov,
				ret: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:   div,
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
		},
		constants: []int{2, 4},
	}
	result := Interpret(&ir)
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestLoadAndStore(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  mov,
				ret: MakeVirtualRegister(2),
				arg2: Arg{
					argType: constant,
					value:   0,
				},
			},
			Instruction{
				op:   store,
				arg1: MakeVirtualRegister(2),
				arg2: Arg{
					argType:        address,
					value:          0,
					offsetConstant: 1,
				},
			},
			Instruction{
				op:  load,
				ret: MakeVirtualRegister(1),
				arg2: Arg{
					argType:        address,
					value:          0,
					offsetConstant: 1,
				},
			},
		},
		constants: []int{2, 1},
	}
	result := Interpret(&ir)
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

// func TestSimpleExpr(t *testing.T) {
// 	// Representing the simple postfix expression 1 2 3 * +
// 	ir := IR{
// 		registersLength: 5,
// 		instructions: []Isntruction{
// 			Instruction{
// 				op: mov,
// 				ret: MakeVirtualRegister(1),
// 				arg2: Arg{
// 					argType: constant,
// 					value: 0
// 				},
// 			},
// 		},
// 		constants: []int{2, 3}
// 	}

// }
