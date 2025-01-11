package navm

type Runtime struct {
	registers []int
}

func Interpret(ir *IR) int {
	r := Runtime{registers: make([]int, ir.registersLength)}
	for _, i := range ir.instructions {
		var arg1, arg2 = 0, 0
		switch i.op {
		case add:
			switch i.arg1.argType {
			case noArgType:
				panic("No argument type for add op")
			case constant:
				arg1 = ir.constants[i.arg1.value]
			case virtualRegisterArg:
				arg1 = r.registers[i.arg1.value]
			case physicalRegisterArg:
				panic("Physical register not legal when interpreting")
			default:
				panic("Unknown argument type")
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
			r.registers[i.ret.value] = arg1 + arg2
		default:
			panic("Unknown operation")
		}
	}
	return r.registers[1]
}

func makeVirtualRegister(value int) Register {
	return Register{registerType: virtualRegister, value: value}
}
