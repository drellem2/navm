package navm

import (
	"strconv"
)

type WinGenerator struct {
	arch *Architecture
	ir   *IR
}

func (g *WinGenerator) Init(a *Architecture, ir *IR) {
	g.arch = a
	g.ir = ir
}

func (g *WinGenerator) GetHeader() string {
	return "section .text\n\tglobal main\n\nmain:\n"
}

// For now we just assume the last register assigned is the return register
func (g *WinGenerator) GetReturn() string {
	return "  ret\n"
}

func (g *WinGenerator) GetTwoArgInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	retRegister := g.arch.GetPhysicalRegister(instr.ret.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + retRegister + ", " + arg2 + "\n"
}

func (g *WinGenerator) GetTwoArgNoRetInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	arg1Register := g.arch.GetPhysicalRegister(instr.arg1.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + arg1Register + ", " + arg2 + "\n"
}

func (g *WinGenerator) GetInstruction(op GenOp, instr Instruction) string {
	name := g.GetTargetInstruction(op)
	retRegister := g.arch.GetPhysicalRegister(instr.ret.value)
	arg1 := g.arch.GetPhysicalRegister(instr.arg1.value)
	arg2 := g.GetArg(instr.arg2)
	return "  " + name + " " + retRegister + ", " + arg1 + ", " + arg2 + "\n"
}

func (g *WinGenerator) GetAddress(arg Arg) string {
	return "[" + g.arch.GetPhysicalRegister(arg.value) + ", " + strconv.Itoa(g.ir.constants[arg.offsetConstant]) + "]"
}

func (g *WinGenerator) GetConstant(i int) string {
	return strconv.Itoa(g.ir.constants[i])
}

func (g *WinGenerator) GetArg(arg Arg) string {
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

func (g *WinGenerator) GetTargetInstruction(op GenOp) string {
	switch op {
	case addGenOp:
		return "add"
	case subGenOp:
		return "sub"
	case multGenOp:
		return "mul"
	case divGenOp:
		return "div"
	case loadGenOp:
		return "ldr"
	case storeGenOp:
		return "str"
	case movGenOp:
		return "mov"
	default:
		panic("Unknown instruction")
	}
}
