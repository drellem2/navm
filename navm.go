package navm

import (
	"strconv"
)

type Type int

const (
	noType Type = iota
	i64    Type = iota
)

type Op int

const (
	noOp Op = iota
	add  Op = iota
	mov  Op = iota
)

// Add concept of virtual vs physical registers

type RegisterType int

const (
	noRegisterType   RegisterType = iota
	virtualRegister  RegisterType = iota
	physicalRegister RegisterType = iota
)

type Register struct {
	registerType RegisterType
	value        int
}

// 0 register is not used, will indicate "no register"
// 1 register will indicate the return value

type ArgType int

const (
	noArgType           ArgType = iota
	virtualRegisterArg  ArgType = iota
	physicalRegisterArg ArgType = iota
	constant            ArgType = iota
)

// Basically a union
type Arg struct {
	argType ArgType
	value   int
}

type Instruction struct {
	op   Op
	ret  Register
	arg1 Register
	arg2 Arg
}

type IR struct {
	registersLength int // maximum register number + 1
	instructions    []Instruction
	constants       []int
}

func (i *Instruction) Print() string {
	return strconv.Itoa(int(i.op)) + " " + strconv.Itoa(i.ret.value) + " " + strconv.Itoa(i.arg1.value) + " " + strconv.Itoa(i.arg2.value)
}
