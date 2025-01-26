package navm

import (
	"strconv"
)

const AARCH64_MACOS_NONE = "aarch64-macos-none"
const X64_WIN_GNU = "x86_64-windows-gnu"

type Architecture struct {
	TargetTriple         string
	Registers64          []string
	ReturnRegister       string
	StackPointerRegister string
	IntSize              int
	StackAlignmentSize   int
}

var Architectures = map[string]*Architecture{
	AARCH64_MACOS_NONE: MakeAarch64MacArchitecture(),
	X64_WIN_GNU:        MakeX64WinGnuArchitecture(),
}

func (a *Architecture) GetGenerator(ir *IR) Generator {
	switch a.TargetTriple {
	case AARCH64_MACOS_NONE:
		return &MacGenerator{
			arch: a,
			ir:   ir,
		}
	case X64_WIN_GNU:
		return &WinGenerator{
			arch: a,
			ir:   ir,
		}
	default:
		panic("Unknown target triple: " + a.TargetTriple)
	}
}

var aarchMac64Registers = []string{"X9", "X10", "X11", "X12", "X13", "X14", "X15"}
var aarchMacReturnRegister = "X0"
var aarchMacStackPointerRegister = "SP"

// use x86_64 registers, not arm
var x64WinGnuRegisters = []string{"R10", "R11", "R12", "R13", "R14", "R15"}
var x64WinGnuReturnRegister = "RAX"
var x64WinGnuStackPointerRegister = "RSP"

func MakeAarch64MacArchitecture() *Architecture {
	return &Architecture{
		TargetTriple:         AARCH64_MACOS_NONE,
		Registers64:          aarchMac64Registers,
		ReturnRegister:       aarchMacReturnRegister,
		StackPointerRegister: aarchMacStackPointerRegister,
		IntSize:              8,
		StackAlignmentSize:   16,
	}
}

func MakeX64WinGnuArchitecture() *Architecture {
	return &Architecture{
		TargetTriple:         X64_WIN_GNU,
		Registers64:          x64WinGnuRegisters,
		ReturnRegister:       x64WinGnuReturnRegister,
		StackPointerRegister: x64WinGnuStackPointerRegister,
		IntSize:              8,
		StackAlignmentSize:   16,
	}
}

func (a *Architecture) GetPhysicalRegister(register int) string {
	if register == 0 {
		panic("0 register should never be used")
	}
	if register == STACK_POINTER_REGISTER {
		return a.StackPointerRegister
	}
	if register == RETURN_REGISTER {
		return a.ReturnRegister
	}
	if register < 0 {
		panic("Invalid register: " + strconv.Itoa(register))
	}
	return a.Registers64[register-1]
}

func (a *Architecture) GetReturnRegister() string {
	return a.ReturnRegister
}

func (a *Architecture) GetStackPointerRegister() string {
	return a.StackPointerRegister
}
