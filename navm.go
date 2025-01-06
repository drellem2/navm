package navm

type Type int

const (
	noType Type = iota
	i64    Type = iota
)

type Op int

const (
	noOp Op = iota
	add  Op = iota
)

type Register int

// 0 register is not used, will indicate "no register"
// 1 register will indicate the return value

type ArgType int

const (
	noArgType ArgType = iota
	register  ArgType = iota
	constant  ArgType = iota
)

// Basically a union
type Arg struct {
	argType ArgType
	value   int
}

type Instruction struct {
	op   Op
	ret  Register
	arg1 Arg
	arg2 Arg
}

type IR struct {
	registersLength int //maximum register number + 1
	instructions    []Instruction
	constants       []int
}
