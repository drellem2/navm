package navm

import (
	"strconv"
)

type MacGenerator struct {
	arch *Architecture
	ir   *IR
}

func (g *MacGenerator) Init(a *Architecture, ir *IR) {
	g.arch = a
	g.ir = ir
}

func (g *MacGenerator) GetHeader() string {
	return ".global _main\n.align 2\n\n_main:\n"
}

// For now we just assume the last register assigned is the return register
func (g *MacGenerator) GetReturn() string {
	return "  ret\n"
}

func (g *MacGenerator) GetTwoArgInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	retRegister := g.arch.GetPhysicalRegister(instr.ret.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + retRegister + ", " + arg2 + "\n"
}

func (g *MacGenerator) GetTwoArgNoRetInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	arg1Register := g.arch.GetPhysicalRegister(instr.arg1.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + arg1Register + ", " + arg2 + "\n"
}

func (g *MacGenerator) GetInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	retRegister := g.arch.GetPhysicalRegister(instr.ret.value)
	arg1 := g.arch.GetPhysicalRegister(instr.arg1.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + retRegister + ", " + arg1 + ", " + arg2 + "\n"
}

func (g *MacGenerator) GetAddress(arg Arg) string {
	return "[" + g.arch.GetPhysicalRegister(arg.value) + ", #" + strconv.Itoa(g.ir.constants[arg.offsetConstant]) + "]"
}

func (g *MacGenerator) GetConstant(i int) string {
	return "#" + strconv.Itoa(g.ir.constants[i])
}

func (g *MacGenerator) GetArg(arg Arg) string {
	switch arg.argType {
	case constant:
		return g.GetConstant(arg.value)
	case registerArg:
		if arg.isVirtualRegister {
			panic("Virtual register not allowed at code generation time")
		}
		return g.arch.GetPhysicalRegister(arg.value)
	case address:
		return g.GetAddress(arg)
	default:
		panic("Unknown argument type")
	}
}

func (g *MacGenerator) GetTargetInstruction(op GenOp) string {
	switch op {
	case addGenOp:
		return "add"
	case subGenOp:
		return "sub"
	case multGenOp:
		return "mul"
	case divGenOp:
		return "sdiv"
	case movGenOp:
		return "mov"
	case loadGenOp:
		return "ldr"
	case storeGenOp:
		return "str"
	default:
		panic("Unknown operation: " + strconv.Itoa(int(op)))
	}
}
