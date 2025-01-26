////////////////////////////////////////////////////////////////////////////////
// Code generation. Should be abstract enough that the same compiler can ///////
// be used for every backend. //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

package navm

type Generator interface {
	Init(a *Architecture, ir *IR)
	GetHeader() string
	GetReturn() string
	GetTwoArgInstruction(op GenOp, instr Instruction) string
	GetTwoArgNoRetInstruction(op GenOp, instr Instruction) string
	GetInstruction(op GenOp, instr Instruction) string
	GetArg(arg Arg) string
	GetTargetInstruction(op GenOp) string
}

type GenOp int

const (
	noGenOp    GenOp = iota
	addGenOp   GenOp = iota
	subGenOp   GenOp = iota
	multGenOp  GenOp = iota
	divGenOp   GenOp = iota
	movGenOp   GenOp = iota
	loadGenOp  GenOp = iota
	storeGenOp GenOp = iota
)
