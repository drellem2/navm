package navm

type Runtime struct {
	registers []int
}

func Interpret(ir *IR) int {
	r := Runtime{registers: make([]int, ir.registersLength)}
	for _, i := range ir.instructions {
		arg2 := 0
		switch i.op {
		case add:
			switch i.ret.registerType {
			case noRegisterType:
				panic("No register type for add op")
			case virtualRegister:
				// do nothing
			case physicalRegister:
				panic("Physical register not legal when interpreting")
			default:
				panic("Unknown register type")
			}
			switch i.arg1.registerType {
			case noRegisterType:
				panic("No register type for add op")
			case virtualRegister:
				// do nothing
			case physicalRegister:
				panic("Physical register not legal when interpreting")
			default:
				panic("Unknown register type")
			}
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
		default:
			panic("Unknown operation")
		}
	}
	return r.registers[1]
}

func makeVirtualRegister(value int) Register {
	return Register{registerType: virtualRegister, value: value}
}
