package navm

import (
	"testing"
)

func init() {
}

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
