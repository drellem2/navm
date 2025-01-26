////////////////////////////////////////////////////////////////////////////////
// Code generation. Should be abstract enough that the same compiler can ///////
// be used for every backend. //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

package navm

type Generator interface {
	Init(a *Architecture, ir *IR)
	GetReturn() string
	GetTwoArgInstruction(name string, instr Instruction) string
	GetTwoArgNoRetInstruction(name string, instr Instruction) string
	GetInstruction(name string, instr Instruction) string
	GetArg(arg Arg) string
}
