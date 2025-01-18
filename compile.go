package navm

import (
	"strconv"

	q "github.com/drellem2/navm/internal/queue"
)

// TODO: In priority order
// 0. Addresses need to point to registers, not constants (e.g. [X0, #4]) (not [#0, #4])
// 0.1 Then make sure that these are working in the register allocation / liveness analysis
// 1. Implement register spilling
// 2. Get cross-compilation working w/ zig
var aarchMac64Registers = []string{"X9", "X10", "X11", "X12", "X13", "X14", "X15"}
var aarchMacReturnRegister = "X0"

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

func Compile(ir *IR) string {
	placeConstantsInRegisters(ir)
	println("Before allocating")
	for _, instr := range ir.instructions {
		println(instr.Print())
	}
	allocateRegisters(ir)
	println("After allocating")
	// print instructions now that they are allocated
	for _, instr := range ir.instructions {
		println(instr.Print())
	}

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
			result += getInstruction("add", instr, ir)
			lastRegister = instr.ret.value
		case sub:
			result += getInstruction("sub", instr, ir)
			lastRegister = instr.ret.value
		case mult:
			result += getInstruction("mul", instr, ir)
			lastRegister = instr.ret.value
		case div:
			result += getInstruction("sdiv", instr, ir)
			lastRegister = instr.ret.value
		case mov:
			result += getTwoArgInstruction("mov", instr, ir)
			lastRegister = instr.ret.value
		case load:
			result += getTwoArgInstruction("ldr", instr, ir)
			lastRegister = instr.ret.value
		case store:
			result += getTwoArgNoRetInstruction("str", instr, ir)
		default:
			panic("Unknown operation: " + strconv.Itoa(int(instr.op)))
		}
	}
	result += getFooter(lastRegister)
	return result
}

// For now we just assume the last register assigned is the return register
func getFooter(lastRegister int) string {
	str := ""
	if lastRegister != 0 {
		str += "  mov " + aarchMacReturnRegister + ", " + getPhysicalRegister(lastRegister) + "\n"
	} else {
		str += "  mov " + aarchMacReturnRegister + ", #0\n"
	}
	str += "  mov X16, #1\n"
	str += "  svc 0\n"
	return str
}

func getTwoArgInstruction(name string, instr Instruction, ir *IR) string {
	retRegister := getPhysicalRegister(instr.ret.value)
	arg2 := getArg(instr.arg2, ir)
	return "  " + name + " " + retRegister + ", " + arg2 + "\n"
}

func getTwoArgNoRetInstruction(name string, instr Instruction, ir *IR) string {
	println("getTwoArgNoRetInstruction: ", instr.arg1.value)
	arg1Register := getPhysicalRegister(instr.arg1.value)
	arg2 := getArg(instr.arg2, ir)
	return "  " + name + " " + arg1Register + ", " + arg2 + "\n"
}

func getInstruction(name string, instr Instruction, ir *IR) string {
	retRegister := getPhysicalRegister(instr.ret.value)
	arg1 := getPhysicalRegister(instr.arg1.value)
	arg2 := getArg(instr.arg2, ir)
	return "  " + name + " " + retRegister + ", " + arg1 + ", " + arg2 + "\n"
}

func getArg(arg Arg, ir *IR) string {
	switch arg.argType {
	case constant:
		return getConstant(arg.value, ir)
	case registerArg:
		if arg.isVirtualRegister {
			panic("Virtual register not allowed at code generation time")
		}
		return getPhysicalRegister(arg.value)
	case address:
		return getAddress(arg, ir)
	default:
		panic("Unknown argument type")
	}
}

func getAddress(arg Arg, ir *IR) string {
	println("GetAddress: " + strconv.Itoa(arg.value))
	return "[" + getPhysicalRegister(arg.value) + ", #" + strconv.Itoa(ir.constants[arg.offsetConstant]) + "]"
}

func getConstant(i int, ir *IR) string {
	return "#" + strconv.Itoa(ir.constants[i])
}

func getPhysicalRegister(i int) string {
	if i == 0 {
		panic("0 register should never be used")
	}
	return aarchMac64Registers[i-1]
}

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
		println("Assigning register", interval.register.value, " to ", interval.register.value)
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
		if instr.arg2.argType == registerArg && instr.arg2.isVirtualRegister {
			instr.arg2.argType = registerArg
			instr.arg2.isVirtualRegister = false
			instr.arg2.value = allocated[instr.arg2.value]
		}
		ir.instructions[i] = instr
	}
}
