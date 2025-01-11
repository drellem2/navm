package navm

// TODO: build liveness intervals and perform simple linear scan
// to allocate registers.
// First step: for 64-bit integers, using X9-X15 registers. Panic if not enough
// Second: implement register spilling

var aarchMac64Registers = []string{"X9", "X10", "X11", "X12", "X13", "X14", "X15"}
var aarchMacReturnRegister = "X0"

func compile(ir *IR) {
	// Build liveness intervals
	// Perform linear scan register allocation

	// activeQueue := LivenessQueue{active: true}
	// inactiveQueue := LivenessQueue{active: false}

	// // First we will make intervals for all virtual registers
	// intervals := makeIntervals(ir)

}

func makeIntervals(ir *IR) []Interval {
	intervals := make([]Interval, ir.registersLength)
	// always skip first, because 0th register is unused
	// range from 1 to len(intervals)-1
	for i := 1; i < len(intervals); i++ {
		intervals[i] = Interval{register: Register{registerType: virtualRegister, value: i}}
		// Set start to max
		intervals[i].start = len(ir.instructions)
	}

	for i, instr := range ir.instructions {
		// Get all virtual registers used in this instruction
		// and update their intervals
		if instr.arg1.argType == virtualRegisterArg {
			intervals[instr.arg1.value].start = min(intervals[instr.arg1.value].start, i)
			intervals[instr.arg1.value].end = i + 1
		}
		if instr.ret.registerType == virtualRegister {
			intervals[instr.ret.value].start = min(intervals[instr.ret.value].start, i)
			intervals[instr.ret.value].end = i + 1
		}
		if instr.arg2.argType == virtualRegisterArg {
			intervals[instr.arg2.value].start = min(intervals[instr.arg2.value].start, i)
			intervals[instr.arg2.value].end = i + 1
		}
	}

	return intervals
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
