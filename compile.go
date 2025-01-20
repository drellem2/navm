package navm

import (
	"strconv"

	q "github.com/drellem2/navm/internal/queue"
)

// TODO: In priority order
// 0. Add stack manipulation commands and test add to load/store test
//    so it won't segfault

// 1. Add sp register
/*
   // Example:
   SUB SP, SP, #16        // Reserve 64 bytes on the stack
   STR X0, [SP]           // Store X0 at the top of the stack
   STR X1, [SP, #8]       // Store X1 at SP + 8
   ADD SP, SP, #64        // Free 64 bytes by restoring SP
*/
// 1. Implement register spilling (on the stack)
// 2. Get cross-compilation working w/ zig

// Add an initial pass that forces constants certain constants into registers
// e.g. in arm, both mul operands must be in registers
func placeConstantsInRegisters(ir *IR) {
	// Find all constants that are used in mult/div instructions
	// Place them in registers
	xns := ir.instructions[:0]
	for _, instr := range ir.instructions {
		if instr.op == mult || instr.op == div {
			if instr.arg2.argType == constant {
				vreg := ir.NewVirtualRegister()
				// move instruction
				movInstr := Instruction{
					op:  mov,
					ret: vreg,
					arg2: Arg{
						argType: constant,
						value:   instr.arg2.value,
					},
				}
				xns = append(xns, movInstr)
				instr.arg2 = Arg{
					argType:           registerArg,
					isVirtualRegister: true,
					value:             vreg.value,
				}
				xns = append(xns, instr)
				continue
			}
		}
		xns = append(xns, instr)
	}
	ir.instructions = xns
}

func Compile(ir *IR, architecture string) string {
	a := Architectures[architecture]
	if a == nil {
		panic("Unknown or unsupported architecture: " + architecture)
	}
	placeConstantsInRegisters(ir)
	allocateRegisters(a, ir)

	// Now do really simple code generation
	// Example:
	// .global _start
	// .align 2
	// _start:
	//   mov X9, #1
	//   mov X10, #2
	//   add X0, X9, X10
	//   mov X16, #1
	//   svc 0

	result := ".global _start\n.align 2\n\n_start:\n"
	lastRegister := 0
	for _, instr := range ir.instructions {
		switch instr.op {
		case add:
			result += getInstruction(a, "add", instr, ir)
			lastRegister = instr.ret.value
		case sub:
			result += getInstruction(a, "sub", instr, ir)
			lastRegister = instr.ret.value
		case mult:
			result += getInstruction(a, "mul", instr, ir)
			lastRegister = instr.ret.value
		case div:
			result += getInstruction(a, "sdiv", instr, ir)
			lastRegister = instr.ret.value
		case mov:
			result += getTwoArgInstruction(a, "mov", instr, ir)
			lastRegister = instr.ret.value
		case load:
			result += getTwoArgInstruction(a, "ldr", instr, ir)
			lastRegister = instr.ret.value
		case store:
			result += getTwoArgNoRetInstruction(a, "str", instr, ir)
		default:
			panic("Unknown operation: " + strconv.Itoa(int(instr.op)))
		}
	}
	result += getFooter(a, lastRegister)
	return result
}

// For now we just assume the last register assigned is the return register
func getFooter(a *Architecture, lastRegister int) string {
	str := ""
	if lastRegister != 0 {
		str += "  mov " + a.ReturnRegister + ", " + a.GetPhysicalRegister(lastRegister) + "\n"
	} else {
		str += "  mov " + a.ReturnRegister + ", #0\n"
	}
	str += "  mov X16, #1\n"
	str += "  svc 0\n"
	return str
}

func getTwoArgInstruction(a *Architecture, name string, instr Instruction, ir *IR) string {
	retRegister := a.GetPhysicalRegister(instr.ret.value)
	arg2 := getArg(a, instr.arg2, ir)
	return "  " + name + " " + retRegister + ", " + arg2 + "\n"
}

func getTwoArgNoRetInstruction(a *Architecture, name string, instr Instruction, ir *IR) string {
	arg1Register := a.GetPhysicalRegister(instr.arg1.value)
	arg2 := getArg(a, instr.arg2, ir)
	return "  " + name + " " + arg1Register + ", " + arg2 + "\n"
}

func getInstruction(a *Architecture, name string, instr Instruction, ir *IR) string {
	retRegister := a.GetPhysicalRegister(instr.ret.value)
	arg1 := a.GetPhysicalRegister(instr.arg1.value)
	arg2 := getArg(a, instr.arg2, ir)
	return "  " + name + " " + retRegister + ", " + arg1 + ", " + arg2 + "\n"
}

func getArg(a *Architecture, arg Arg, ir *IR) string {
	switch arg.argType {
	case constant:
		return getConstant(arg.value, ir)
	case registerArg:
		if arg.isVirtualRegister {
			panic("Virtual register not allowed at code generation time")
		}
		return a.GetPhysicalRegister(arg.value)
	case address:
		return getAddress(a, arg, ir)
	default:
		panic("Unknown argument type")
	}
}

func getAddress(a *Architecture, arg Arg, ir *IR) string {
	return "[" + a.GetPhysicalRegister(arg.value) + ", #" + strconv.Itoa(ir.constants[arg.offsetConstant]) + "]"
}

func getConstant(i int, ir *IR) string {
	return "#" + strconv.Itoa(ir.constants[i])
}

func allocateRegisters(a *Architecture, ir *IR) {
	// Build liveness intervals
	// Perform linear scan register allocation

	activeQueue := LivenessQueue{active: true}
	inactiveQueue := LivenessQueue{active: false}
	finishedQueue := LivenessQueue{active: true}

	// maps vregisters to physical registers
	allocated := make([]int, ir.registersLength)

	// Free physical registers are just a simple queue, not a priority queue
	physicalRegisters := q.Queue{}
	for i := 0; i < len(a.Registers64); i++ {
		physicalRegisters.Push(i + 1)
	}

	// First we will make intervals for all virtual registers
	intervals := makeIntervals(ir)

	// Push all intervals to inactive queue
	for _, val := range intervals[1:] {
		inactiveQueue.Push(val)
	}

	// Linear scan, we iterate through inactive queue and try to assign
	// registers

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
			finishedQueue.Push(finished)
			physicalRegisters.Push(finished.physicalRegister)
		}

		// assign a register
		interval.physicalRegister = physicalRegisters.Pop()
		activeQueue.Push(interval)
	}

	// Add remaining active intervals to finished queue
	for !activeQueue.Empty() {
		finishedQueue.Push(activeQueue.Pop())
	}

	// Iterate over finished
	for !finishedQueue.Empty() {
		finished := finishedQueue.Pop()
		allocated[finished.register.value] = finished.physicalRegister
	}

	// Now iterate through instructions and set all virtual registers to physical registers
	for i, instr := range ir.instructions {
		if instr.arg1.registerType == virtualRegister {
			instr.arg1.registerType = physicalRegister
			instr.arg1.value = allocated[instr.arg1.value]
		}
		if instr.ret.registerType == virtualRegister {
			instr.ret.registerType = physicalRegister
			instr.ret.value = allocated[instr.ret.value]
		}
		if instr.arg2.isVirtualRegister && (instr.arg2.argType == registerArg || instr.arg2.argType == address) {
			instr.arg2.isVirtualRegister = false
			instr.arg2.value = allocated[instr.arg2.value]
		}
		ir.instructions[i] = instr
	}
}
