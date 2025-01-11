package navm

type Runtime struct {
	registers []int
}

func validateRegister(r Register) {
	switch r.registerType {
	case noRegisterType:
		panic("No register type")
	case virtualRegister:
		// do nothing
	case physicalRegister:
		panic("Physical register not legal when interpreting")
	default:
		panic("Unknown register type")
	}
}

func Interpret(ir *IR) int {
	r := Runtime{registers: make([]int, ir.registersLength)}
	lastAssignedRegister := 0
	for _, i := range ir.instructions {
		switch i.op {
		case add:
			lastAssignedRegister = i.ret.value
			runAdd(i, &r, ir)
		case mov:
			lastAssignedRegister = i.ret.value
			runMov(i, &r, ir)

		default:
			panic("Unknown operation")
		}
	}
	if lastAssignedRegister == 0 {
		return 0
	}
	return r.registers[1]
}

func runMov(i Instruction, r *Runtime, ir *IR) {
	arg2 := 0
	validateRegister(i.ret)
	if i.arg1.registerType != noRegisterType {
		panic("arg1 should not be set for MOV instructions")
	}
	switch i.arg2.argType {
	case noArgType:
		panic("No argument type for mov op")
	case constant:
		arg2 = ir.constants[i.arg2.value]
	case physicalRegisterArg:
		panic("Physical register not legal when interpreting")
	case virtualRegisterArg:
		arg2 = r.registers[i.arg2.value]
	default:
		panic("Unknown argument type")
	}
	r.registers[i.ret.value] = arg2
}

func runAdd(i Instruction, r *Runtime, ir *IR) {
	arg2 := 0
	validateRegister(i.ret)
	validateRegister(i.arg1)
	switch i.arg2.argType {
	case noArgType:
		panic("No argument type for add op")
	case constant:
		arg2 = ir.constants[i.arg2.value]
	case virtualRegisterArg:
		arg2 = r.registers[i.arg2.value]
	case physicalRegisterArg:
		panic("Physical register not legal when interpreting")
	default:
		panic("Unknown argument type")
	}
	r.registers[i.ret.value] = r.registers[i.arg1.value] + arg2
}
