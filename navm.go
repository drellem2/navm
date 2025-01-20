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
	noOp  Op = iota
	add   Op = iota
	mov   Op = iota
	sub   Op = iota
	mult  Op = iota
	div   Op = iota
	load  Op = iota
	store Op = iota
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

func (r Register) Value() int {
	return r.value
}

const STACK_POINTER_REGISTER = -1

// 0 register is not used, will indicate "no register"
// 1 register will indicate the return value

type ArgType int

const (
	noArgType   ArgType = iota
	registerArg ArgType = iota
	constant    ArgType = iota
	address     ArgType = iota
)

// Basically a union
type Arg struct {
	argType           ArgType
	isVirtualRegister bool
	value             int
	offsetConstant    int
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

func NewIR() *IR {
	return &IR{registersLength: 1}
}

func (ir *IR) GetConstant(c int) int {
	for idx, i := range ir.constants {
		if i == c {
			return idx
		}
	}
	ir.constants = append(ir.constants, c)
	return len(ir.constants) - 1
}

func (ir *IR) MoveConstant(r Register, c int) {
	xrn := Instruction{op: mov, ret: r, arg2: Arg{argType: constant, value: ir.GetConstant(c)}}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) AddRegisters(ret Register, r1 Register, r2 Register) {
	xrn := Instruction{op: add, ret: ret, arg1: r1, arg2: Arg{
		argType:           registerArg,
		isVirtualRegister: true,
		value:             r2.value}}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) SubRegisters(ret Register, r1 Register, r2 Register) {
	xrn := Instruction{op: sub, ret: ret, arg1: r1, arg2: Arg{
		argType:           registerArg,
		isVirtualRegister: true,
		value:             r2.value}}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) MultRegisters(ret Register, r1 Register, r2 Register) {
	xrn := Instruction{op: mult, ret: ret, arg1: r1, arg2: Arg{
		argType:           registerArg,
		isVirtualRegister: true,
		value:             r2.value}}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) DivRegisters(ret Register, r1 Register, r2 Register) {
	xrn := Instruction{op: div, ret: ret, arg1: r1, arg2: Arg{
		argType:           registerArg,
		isVirtualRegister: true,
		value:             r2.value}}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) Load(ret Register, addr Arg) {
	if addr.argType != address {
		panic("Argument addr must be an address")
	}
	xrn := Instruction{op: load, ret: ret, arg2: addr}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) Store(reg Register, addr Arg) {
	if addr.argType != address {
		panic("Argument addr must be an address")
	}
	xrn := Instruction{op: store, ret: reg, arg2: addr}
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) AddInstruction(xrn Instruction) {
	ir.instructions = append(ir.instructions, xrn)
}

func (ir *IR) NewVirtualRegister() Register {
	ret := Register{registerType: virtualRegister, value: ir.registersLength}
	ir.registersLength = ir.registersLength + 1
	return ret
}

func MakeVirtualRegister(value int) Register {
	return Register{registerType: virtualRegister, value: value}
}

func GetStackPointer() Register {
	return Register{registerType: physicalRegister, value: STACK_POINTER_REGISTER}
}

func (r Register) ToArg() Arg {
	return Arg{argType: registerArg, isVirtualRegister: r.registerType == virtualRegister, value: r.value}
}

func (r Register) ToAddress(offset int) Arg {
	return Arg{argType: address, isVirtualRegister: r.registerType == virtualRegister, value: r.value, offsetConstant: offset}
}

func (ir *IR) Print() string {
	ret := ""
	for _, i := range ir.instructions {
		ret += i.Print() + "\n"
	}
	// Add constants
	ret += "Constants: "
	for _, c := range ir.constants {
		ret += strconv.Itoa(c) + " "
	}
	return ret
}

func (i *Instruction) Print() string {
	return strconv.Itoa(int(i.op)) + " " + strconv.Itoa(i.ret.value) + " " + strconv.Itoa(i.arg1.value) + " " + strconv.Itoa(i.arg2.value)
}
