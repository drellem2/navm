package navm

import (
	q "github.com/drellem2/navm/internal/queue"
)

// TODO: build liveness intervals and perform simple linear scan
// to allocate registers.
// First step: for 64-bit integers, using X9-X15 registers. Panic if not enough
// Second: implement register spilling

var aarchMac64Registers = []string{"X9", "X10", "X11", "X12", "X13", "X14", "X15"}
var aarchMacReturnRegister = "X0"

func allocateRegisters(ir *IR) {
	// Build liveness intervals
	// Perform linear scan register allocation

	activeQueue := LivenessQueue{active: true}
	inactiveQueue := LivenessQueue{active: false}
	finishedQueue := LivenessQueue{active: true}

	// maps vregisters to physical registers
	allocated := make([]int, ir.registersLength)

	// Free physical registers are just a simple queue, not a priority queue
	physicalRegisters := q.Queue{}
	for i := 0; i < len(aarchMac64Registers); i++ {
		physicalRegisters.Push(i + 1)
	}

	println("Physical registers: ", physicalRegisters.Print())

	// First we will make intervals for all virtual registers
	intervals := makeIntervals(ir)

	// Push all intervals to inactive queue
	for _, val := range intervals[1:] {
		inactiveQueue.Push(val)
	}

	// Linear scan, we iterate through inactive queue and try to assign
	// registers

	println("Inactive queue: ", inactiveQueue.Print())

	for !inactiveQueue.Empty() {
		interval := inactiveQueue.Pop()
		// Check if we can assign a register
		if physicalRegisters.Empty() {
			// Spill register
			panic("Too many virtual registers - spilling not implemented")
		}

		// Free all registers that are not live anymore
		for !activeQueue.Empty() && activeQueue.Peek().end <= interval.start {
			finished := activeQueue.Pop()
			physicalRegisters.Push(finished.physicalRegister)
		}

		// assign a register
		interval.physicalRegister = physicalRegisters.Pop()
		println("Assigned: ", interval.register.value, " to ", interval.physicalRegister)
		activeQueue.Push(interval)
	}

	// Add remaining active intervals to finished queue
	for !activeQueue.Empty() {
		finishedQueue.Push(activeQueue.Pop())
	}

	// Iterate over finished
	for !finishedQueue.Empty() {
		finished := finishedQueue.Pop()
		println("Finished: ", finished.register.value, " to ", finished.physicalRegister)
		allocated[finished.register.value] = finished.physicalRegister
	}

	println("Print all allocated values")
	for k, v := range allocated {
		println("Allocated: ", k, " to ", v)
	}

	// Now iterate through instructions and set all virtual registers to physical registers
	for i, instr := range ir.instructions {
		if instr.arg1.argType == virtualRegisterArg {
			instr.arg1.argType = physicalRegisterArg
			instr.arg1.value = allocated[instr.arg1.value]
		}
		if instr.ret.registerType == virtualRegister {
			instr.ret.registerType = physicalRegister
			instr.ret.value = allocated[instr.ret.value]
		}
		if instr.arg2.argType == virtualRegisterArg {
			instr.arg2.argType = physicalRegisterArg
			instr.arg2.value = allocated[instr.arg2.value]
		}
		ir.instructions[i] = instr
	}
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
