package navm

import (
	"testing"
)

func init() {
}

// TODO: use property based testing

func TestMakeIntervals(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  add,
				ret: makeVirtualRegister(2),
				arg1: Arg{
					argType: constant,
					value:   0,
				},
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:  add,
				ret: makeVirtualRegister(1),
				arg1: Arg{
					argType: virtualRegisterArg,
					value:   2,
				},
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
				op:  add,
				ret: makeVirtualRegister(2),
				arg1: Arg{
					argType: constant,
					value:   0,
				},
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:  add,
				ret: makeVirtualRegister(1),
				arg1: Arg{
					argType: virtualRegisterArg,
					value:   2,
				},
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
		},
		constants: []int{1, 2},
	}
	allocateRegisters(&ir)
	// Print instructions
	for _, i := range ir.instructions {
		t.Logf(i.Print())
	}

	// Check all of the registers are physical & none are 0
	for _, i := range ir.instructions {
		if i.arg1.argType == virtualRegisterArg {
			t.Errorf("Expected physical register, got virtual")
		}
		if i.ret.registerType == virtualRegister {
			t.Errorf("Expected physical register, got virtual")
		}
		if i.arg2.argType == virtualRegisterArg {
			t.Errorf("Expected physical register, got virtual")
		}
	}

	for _, i := range ir.instructions {
		if i.arg1.argType == physicalRegisterArg && i.arg1.value == 0 {
			t.Errorf("Expected non-zero register")
		}
		if i.ret.registerType == physicalRegister && i.ret.value == 0 {
			t.Errorf("Expected non-zero register")
		}
		if i.arg2.argType == physicalRegisterArg && i.arg2.value == 0 {
			t.Errorf("Expected non-zero register")
		}
	}
	if ir.registersLength != 3 {
		t.Errorf("Expected 3, got %d", ir.registersLength)
	}
}

func TestCompile(t *testing.T) {
	ir := IR{
		registersLength: 3,
		instructions: []Instruction{
			Instruction{
				op:  add,
				ret: makeVirtualRegister(2),
				arg1: Arg{
					argType: constant,
					value:   0,
				},
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
			Instruction{
				op:  add,
				ret: makeVirtualRegister(1),
				arg1: Arg{
					argType: virtualRegisterArg,
					value:   2,
				},
				arg2: Arg{
					argType: constant,
					value:   1,
				},
			},
		},
		constants: []int{1, 2},
	}

	result := compile(&ir)
	println(result)
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

}
