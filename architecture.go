package navm

import (
	"strconv"
)

const AARCH64_MACOS_NONE = "aarch64-macos-none"

type Architecture struct {
	Registers64          []string
	ReturnRegister       string
	StackPointerRegister string
}

var Architectures = map[string]*Architecture{
	AARCH64_MACOS_NONE: MakeAarch64MacArchitecture(),
}

var aarchMac64Registers = []string{"X9", "X10", "X11", "X12", "X13", "X14", "X15"}
var aarchMacReturnRegister = "X0"

func MakeAarch64MacArchitecture() *Architecture {
	return &Architecture{
		Registers64:          aarchMac64Registers,
		ReturnRegister:       aarchMacReturnRegister,
		StackPointerRegister: "SP",
	}
}

func (a *Architecture) GetPhysicalRegister(register int) string {
	if register == 0 {
		panic("0 register should never be used")
	}
	if register == STACK_POINTER_REGISTER {
		return a.StackPointerRegister
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
