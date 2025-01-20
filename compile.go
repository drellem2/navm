package navm

import (
	"strconv"

	q "github.com/drellem2/navm/internal/queue"
)

const scratch_register_1 = 1
const scratch_register_2 = 2
const scratch_register_count = 2

// TODO: In priority order
// 1. Get cross-compilation working w/ zig

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

	// Now we need to deal with any spilled registers
	// 1. Allocate enough space on the stack for all spilled registers, respecting alignment
	// 1.1 Ensure that there is one additional stack location for swapping?
	// 2. Add instructions to spill registers
	// 3. Add instructions to free the stack

	stackMax := makeStackSpace(a, ir)
	addSpillInstructions(a, ir)
	freeStackSpace(a, ir, stackMax)

	result := ".global _start\n.align 2\n\n_start:\n"
	for _, instr := range ir.instructions {
		switch instr.op {
		case add:
			result += getInstruction(a, "add", instr, ir)
		case sub:
			result += getInstruction(a, "sub", instr, ir)
		case mult:
			result += getInstruction(a, "mul", instr, ir)
		case div:
			result += getInstruction(a, "sdiv", instr, ir)
		case mov:
			result += getTwoArgInstruction(a, "mov", instr, ir)
		case load:
			result += getTwoArgInstruction(a, "ldr", instr, ir)
		case store:
			result += getTwoArgNoRetInstruction(a, "str", instr, ir)
		case ret:
			result += getReturn()
		default:
			panic("Unknown operation: " + strconv.Itoa(int(instr.op)))
		}
	}
	return result
}

func makeStackSpace(a *Architecture, ir *IR) int {
	// Get largest stack position
	var stackMax int
	for _, instr := range ir.instructions {
		if instr.arg2.argType == stackArg {
			if instr.arg2.value > stackMax {
				stackMax = instr.arg2.value
			}
		}
		if instr.ret.registerType == stackRegister {
			if instr.ret.value > stackMax {
				stackMax = instr.ret.value
			}
		}
		if instr.arg1.registerType == stackRegister {
			if instr.arg1.value > stackMax {
				stackMax = instr.arg1.value
			}
		}
	}
	if stackMax == 0 {
		return 0
	}
	stackMaxSize := stackMax * a.IntSize
	alignRemainder := stackMaxSize % a.StackAlignmentSize
	if alignRemainder != 0 {
		stackMaxSize += a.StackAlignmentSize - alignRemainder
	}
	if stackMaxSize%a.IntSize != 0 {
		panic("Stack size not a multiple of int size")
	}
	stackMax = stackMaxSize / a.IntSize
	xrn := Instruction{
		op:   sub,
		ret:  GetStackPointer(),
		arg1: GetStackPointer(),
		arg2: MakeConstant(ir.GetConstant(stackMax * a.IntSize)),
	}
	ir.instructions = append([]Instruction{xrn}, ir.instructions...)
	return stackMax
}

func freeStackSpace(a *Architecture, ir *IR, stackMax int) {
	if stackMax == 0 {
		return
	}
	xrn := Instruction{
		op:   add,
		ret:  GetStackPointer(),
		arg1: GetStackPointer(),
		arg2: MakeConstant(ir.GetConstant(stackMax * a.IntSize)),
	}
	finalInstr := ir.instructions[len(ir.instructions)-1]
	if finalInstr.op == ret {
		ir.instructions = append(ir.instructions[:len(ir.instructions)-1], xrn, finalInstr)
	} else {
		ir.instructions = append(ir.instructions, xrn)
	}
}

func addSpillInstructions(a *Architecture, ir *IR) {
	xns := make([]Instruction, 0)
	for _, instr := range ir.instructions {
		if instr.arg1.registerType == stackRegister {
			tmpReg1 := MakePhysicalRegister(scratch_register_1)
			loadXrn := Instruction{
				op:   load,
				ret:  tmpReg1,
				arg2: GetStackAddress(a, ir, instr.arg1.value),
			}
			xns = append(xns, loadXrn)
			instr.arg1 = tmpReg1
		}
		if instr.arg2.argType == stackArg {
			tmpReg2 := MakePhysicalRegister(scratch_register_2)
			loadXrn := Instruction{
				op:   load,
				ret:  tmpReg2,
				arg2: GetStackAddress(a, ir, instr.arg2.value),
			}
			xns = append(xns, loadXrn)
			instr.arg2 = tmpReg2.ToArg()
		}
		var storeNeeded bool
		var storeStackPos Arg
		if instr.ret.registerType == stackRegister {
			storeStackPos = GetStackAddress(a, ir, instr.ret.value)
			instr.ret = MakePhysicalRegister(scratch_register_1)
			storeNeeded = true
		}
		xns = append(xns, instr)
		if storeNeeded {
			storeXrn := Instruction{
				op:   store,
				arg1: MakePhysicalRegister(scratch_register_1),
				arg2: storeStackPos,
			}
			xns = append(xns, storeXrn)
		}
	}
	ir.instructions = xns
}

// Converts our virtual stack pointer into a real address
func GetStackAddress(a *Architecture, ir *IR, stackPos int) Arg {
	return GetStackPointer().ToAddress(ir.GetConstant((stackPos - 1) * a.IntSize))
}

// For now we just assume the last register assigned is the return register
func getReturn() string {
	return "  ret\n"
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

type allocType int

const (
	noAllocType   allocType = iota
	registerAlloc allocType = iota
	stackAlloc    allocType = iota
)

type allocation struct {
	allocTyp allocType
	value    int
}

func makeRegisterAlloc(value int) allocation {
	return allocation{registerAlloc, value}
}

func makeStackAlloc(value int) allocation {
	return allocation{stackAlloc, value}
}

func allocateRegisters(a *Architecture, ir *IR) {
	// Build liveness intervals
	// Perform linear scan register allocation

	activeQueue := LivenessQueue{active: true}
	inactiveQueue := LivenessQueue{active: false}
	finishedQueue := LivenessQueue{active: true}

	var virtualStackPointer int

	// maps vregisters to physical registers
	allocated := make([]allocation, ir.registersLength)

	// Free physical registers are just a simple queue, not a priority queue
	physicalRegisters := q.Queue{}
	// We skip the first two registers so we can use them later as scratch registers
	// when restoring spills
	for i := scratch_register_count; i < len(a.Registers64); i++ {
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
			// panic("Too many virtual registers - spilling not implemented")
			// Here we will choose the active range ending last and spill its register
			// TODO - do comparison against current interval and decide if we need to spill
			spill := activeQueue.PopLast()
			physReg := spill.physicalRegister
			spill.physicalRegister = 0
			virtualStackPointer = virtualStackPointer + 1
			spill.stackPosition = virtualStackPointer
			finishedQueue.Push(spill)

			interval.physicalRegister = physReg
			activeQueue.Push(interval)
		} else {
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
	}

	// Add remaining active intervals to finished queue
	for !activeQueue.Empty() {
		finishedQueue.Push(activeQueue.Pop())
	}

	// Iterate over finished
	for !finishedQueue.Empty() {
		finished := finishedQueue.Pop()
		if finished.stackPosition != 0 {
			allocated[finished.register.value] = makeStackAlloc(finished.stackPosition)
		} else {
			allocated[finished.register.value] = makeRegisterAlloc(finished.physicalRegister)
		}
	}

	// Now iterate through instructions and set all virtual registers to physical registers
	for i, instr := range ir.instructions {
		ir.instructions[i] = allocateInstruction(instr, allocated)
	}
}

func allocateInstruction(instr Instruction, allocated []allocation) Instruction {
	instr.ret = allocateRegister(instr.ret, allocated)
	instr.arg1 = allocateRegister(instr.arg1, allocated)
	instr.arg2 = allocateArg(instr.arg2, allocated)
	return instr
}

func allocateArg(arg Arg, allocated []allocation) Arg {
	if arg.isVirtualRegister && (arg.argType == registerArg || arg.argType == address) {
		arg.isVirtualRegister = false
		if arg.value < 0 { // special registers are not allocated
			return arg
		}
		if allocated[arg.value].allocTyp == stackAlloc {
			arg.argType = stackArg
		}
		arg.value = allocated[arg.value].value
	}
	return arg
}

func allocateRegister(register Register, allocated []allocation) Register {
	if register.registerType == virtualRegister {
		if register.value < 0 { // special registers are not allocated
			register.registerType = physicalRegister
			return register
		}
		if allocated[register.value].allocTyp == stackAlloc {
			register.registerType = stackRegister
		} else {
			register.registerType = physicalRegister
		}
		register.value = allocated[register.value].value
	}
	return register
}
