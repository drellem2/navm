package navm

type Runtime struct {
	returnRegister int
	registers      []int
	memory         []byte
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
	r := Runtime{
		registers: make([]int, ir.registersLength),
		memory:    make([]byte, 1024)}
	for _, i := range ir.instructions {
		(&r).print()
		switch i.op {
		case add:
			runAdd(i, &r, ir)
		case sub:
			runSub(i, &r, ir)
		case mult:
			runMult(i, &r, ir)
		case div:
			runDiv(i, &r, ir)
		case mov:
			runMov(i, &r, ir)
		case load:
			runLoad(i, &r, ir)
		case store:
			runStore(i, &r, ir)
		case ret:
			return r.returnRegister

		default:
			panic("Unknown operation")
		}
	}
	return r.returnRegister
}

func (r *Runtime) getRegister(i int) int {
	if i == STACK_POINTER_REGISTER {
		panic("Stack pointer not implemented in interpreter")
	}
	if i == RETURN_REGISTER {
		return r.returnRegister
	}
	return r.registers[i]
}

func (r *Runtime) setRegister(i int, value int) {
	if i == STACK_POINTER_REGISTER {
		panic("Stack pointer not implemented in interpreter")
	}
	if i == RETURN_REGISTER {
		r.returnRegister = value
		return
	}
	r.registers[i] = value
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
	case registerArg:
		if !i.arg2.isVirtualRegister {
			panic("Physical register not legal when interpreting")
		}
		arg2 = r.getRegister(i.arg2.value)
	default:
		panic("Unknown argument type")
	}
	r.setRegister(i.ret.value, arg2)
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
	case registerArg:
		if !i.arg2.isVirtualRegister {
			panic("Physical register not legal when interpreting")
		}
		arg2 = r.getRegister(i.arg2.value)
	default:
		panic("Unknown argument type")
	}
	r.setRegister(i.ret.value, r.getRegister(i.arg1.value)+arg2)
}

func runSub(i Instruction, r *Runtime, ir *IR) {
	arg2 := 0
	validateRegister(i.ret)
	validateRegister(i.arg1)
	switch i.arg2.argType {
	case noArgType:
		panic("No argument type for add op")
	case constant:
		arg2 = ir.constants[i.arg2.value]
	case registerArg:
		if !i.arg2.isVirtualRegister {
			panic("Physical register not legal when interpreting")
		}
		arg2 = r.getRegister(i.arg2.value)
	default:
		panic("Unknown argument type")
	}
	r.setRegister(i.ret.value, r.getRegister(i.arg1.value)-arg2)
}

func runMult(i Instruction, r *Runtime, ir *IR) {
	arg2 := 0
	validateRegister(i.ret)
	validateRegister(i.arg1)
	switch i.arg2.argType {
	case noArgType:
		panic("No argument type for add op")
	case constant:
		arg2 = ir.constants[i.arg2.value]
	case registerArg:
		if !i.arg2.isVirtualRegister {
			panic("Physical register not legal when interpreting")
		}
		arg2 = r.getRegister(i.arg2.value)
	default:
		panic("Unknown argument type")
	}
	r.setRegister(i.ret.value, r.getRegister(i.arg1.value)*arg2)
}

func runDiv(i Instruction, r *Runtime, ir *IR) {
	arg2 := 0
	validateRegister(i.ret)
	validateRegister(i.arg1)
	switch i.arg2.argType {
	case noArgType:
		panic("No argument type for add op")
	case constant:
		arg2 = ir.constants[i.arg2.value]
	case registerArg:
		if !i.arg2.isVirtualRegister {
			panic("Physical register not legal when interpreting")
		}
		arg2 = r.getRegister(i.arg2.value)
	default:
		panic("Unknown argument type")
	}
	r.setRegister(i.ret.value, r.getRegister(i.arg1.value)/arg2)
}

func runLoad(i Instruction, r *Runtime, ir *IR) {
	validateRegister(i.ret)
	if i.arg2.argType != address {
		panic("Load arg2 should be an address")
	}
	r.setRegister(i.ret.value, 0)
	for t := 0; t < 8; t++ {
		r.setRegister(i.ret.value, r.getRegister(i.ret.value)<<8)
		r.setRegister(i.ret.value, r.getRegister(i.ret.value)|int(r.memory[r.getRegister(i.arg2.value)+ir.constants[i.arg2.offsetConstant]+t]))
	}
}

func runStore(i Instruction, r *Runtime, ir *IR) {
	validateRegister(i.arg1)
	// TODO validateAddress(i.arg2)
	if i.arg2.argType != address {
		panic("Store arg2 should be an address")
	}
	// Now we do the opposite and store 8 bytes
	for t := 0; t < 8; t++ {
		r.memory[r.getRegister(i.arg2.value)+ir.constants[i.arg2.offsetConstant]+t] = byte(r.getRegister(i.arg1.value) >> uint(8*(7-t)))
	}
}

func (r *Runtime) print() {
	for idx, i := range r.registers {
		if i != 0 {
			println("Register ", idx, " = ", i)
		}
	}
	for idx, i := range r.memory {
		if i != 0 {
			println("Memory ", idx, " = ", i)
		}
	}
}
