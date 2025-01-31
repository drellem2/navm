package navm

import (
	"testing"
)

func init() {
}

// TODO: use property based testing

// we don't have move instructions yet, so these programs are assuming 0-initialization of registers

func TestMakeIntervals(t *testing.T) {
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
	intervals := makeIntervals(&ir)
	for _, i := range intervals {
		t.Logf(i.Print())
	}

	if len(intervals) != 3 {
		t.Errorf("Expected 3, got %d", len(intervals))
	}

	if intervals[0].start != 0 || intervals[0].end != 0 {
		t.Errorf("Expected (0, 0), got (%d, %d)", intervals[0].start, intervals[0].end)
	}
	if intervals[1].start != 1 || intervals[1].end != 2 {
		t.Errorf("Expected (1, 2), got (%d, %d)", intervals[1].start, intervals[1].end)
	}
	if intervals[2].start != 0 || intervals[2].end != 2 {
		t.Errorf("Expected (0, 2), got (%d, %d)", intervals[2].start, intervals[2].end)
	}
}

func TestAllocateRegisters(t *testing.T) {
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
	allocateRegisters(Architectures[AARCH64_MACOS_NONE], &ir)
	// Print instructions
	for _, i := range ir.instructions {
		t.Logf(i.Print())
	}

	// Check all of the registers are physical & none are 0
	for _, i := range ir.instructions {
		if i.arg1.registerType == virtualRegister {
			t.Errorf("Expected physical register, got virtual")
		}
		if i.ret.registerType == virtualRegister {
			t.Errorf("Expected physical register, got virtual")
		}
		if i.arg2.argType == registerArg && i.arg2.isVirtualRegister {
			t.Errorf("Expected physical register, got virtual")
		}
	}

	for _, i := range ir.instructions {
		if i.arg1.registerType == physicalRegister && i.arg1.value == 0 {
			t.Errorf("Expected non-zero register")
		}
		if i.ret.registerType == physicalRegister && i.ret.value == 0 {
			t.Errorf("Expected non-zero register")
		}
		if i.arg2.argType == registerArg && !i.arg2.isVirtualRegister && i.arg2.value == 0 {
			t.Errorf("Expected non-zero register")
		}
	}
	if ir.registersLength != 3 {
		t.Errorf("Expected 3, got %d", ir.registersLength)
	}
}

// TODO: actually test output automatically somehow?
func TestAdd(t *testing.T) {
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

	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Expected non-empty string, got %s", result)
	}
}

func TestMov(t *testing.T) {
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
					argType:           registerArg,
					isVirtualRegister: true,
					value:             2,
				},
			},
		},
		constants: []int{1, 2},
	}

	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Expected non-empty string, got %s", result)
	}
}

func TestSub(t *testing.T) {
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
	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Expected non-empty string, got %s", result)
	}
}

func TestMult(t *testing.T) {
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
	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Expected non-empty string, got %s", result)
	}
}

func TestDiv(t *testing.T) {
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
	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Expected non-empty string, got %s", result)
	}
}

func TestLoadAndStoreCompile(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{ // reserve 16 bytes on the stack
				op:   sub,
				ret:  GetStackPointer(),
				arg1: GetStackPointer(),
				arg2: Arg{
					argType: constant,
					value:   2, // only using 8 bytes but need to align to 16
				},
			},
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
				arg2: GetStackPointer().ToAddress(3), // 0 offset
			},
			Instruction{
				op:   load,
				ret:  MakeVirtualRegister(1),
				arg2: GetStackPointer().ToAddress(3), // 0 offset
			},
			Instruction{ // restore 16 bytes on the stack
				op:   add,
				ret:  GetStackPointer(),
				arg1: GetStackPointer(),
				arg2: Arg{
					argType: constant,
					value:   2,
				},
			},
			Instruction{
				op:   add,
				ret:  MakeVirtualRegister(1),
				arg1: MakeVirtualRegister(1),
				arg2: Arg{
					argType: constant,
					value:   3,
				},
			}, // Silly, we add 0 to register 1 to make it the most recently used so it will return
		},
		constants: []int{2, 1, 16, 0},
	}
	result := Compile(&ir, AARCH64_MACOS_NONE)
	if result == "" {
		t.Errorf("Unexpected empty string, got %s", result)
	}
}

func TestSpill(t *testing.T) {
	ir := NewIR()
	ir.MoveConstant(MakeVirtualRegister(1), 0)
	ir.AddRegisters(
		MakeVirtualRegister(2),
		MakeVirtualRegister(1),
		MakeVirtualRegister(1))
	ir.AddRegisters(
		MakeVirtualRegister(3),
		MakeVirtualRegister(2),
		MakeVirtualRegister(2))
	ir.AddRegisters(
		MakeVirtualRegister(4),
		MakeVirtualRegister(3),
		MakeVirtualRegister(3))
	ir.AddRegisters(
		MakeVirtualRegister(5),
		MakeVirtualRegister(4),
		MakeVirtualRegister(4))
	ir.AddRegisters(
		MakeVirtualRegister(6),
		MakeVirtualRegister(5),
		MakeVirtualRegister(5))
	ir.AddRegisters(
		MakeVirtualRegister(7),
		MakeVirtualRegister(6),
		MakeVirtualRegister(6))
	ir.AddRegisters(
		MakeVirtualRegister(8),
		MakeVirtualRegister(7),
		MakeVirtualRegister(1))
	ir.AddRegisters(
		MakeVirtualRegister(9),
		MakeVirtualRegister(7),
		MakeVirtualRegister(2))
	ir.AddRegisters(
		GetReturnRegister(),
		MakeVirtualRegister(8),
		MakeVirtualRegister(4))
	ir.Return()
	ir.registersLength = 10
	ir.constants = []int{1}
	result2 := Interpret(ir)
	result := Compile(ir, AARCH64_MACOS_NONE)
	println("Interpreted result is: ", result2)
	if result == "" {
		t.Errorf("Unexpected empty string, got %s", result)
	}
}
