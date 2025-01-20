package main

import (
	"os"
	"strconv"
	"unicode"

	navm "github.com/drellem2/navm"
	q "github.com/drellem2/navm/internal/queue"
)

type token struct {
	kind  string
	value string
}

type tokenstack struct {
	tokens []token
}

func (ts *tokenstack) push(t token) {
	ts.tokens = append(ts.tokens, t)
}

func (ts *tokenstack) pop() token {
	t := ts.tokens[len(ts.tokens)-1]
	ts.tokens = ts.tokens[:len(ts.tokens)-1]
	return t
}

func (ts *tokenstack) takeFirst() token {
	t := ts.tokens[0]
	ts.tokens = ts.tokens[1:]
	return t
}

func (ts *tokenstack) peek() token {
	return ts.tokens[len(ts.tokens)-1]
}

func (ts *tokenstack) empty() bool {
	return len(ts.tokens) == 0
}

func (ts *tokenstack) len() int {
	return len(ts.tokens)
}

func (ts *tokenstack) print() string {
	var s string
	for _, t := range ts.tokens {
		s += "(" + t.kind + ", " + t.value + ") "
	}
	return s
}

func tokenize(expr string) tokenstack {
	tokens := tokenstack{}
	currentNum := 0
	currentNumSet := false

	for _, c := range expr {
		if c == ' ' {
			if currentNumSet {
				tokens.push(token{kind: "num", value: strconv.Itoa(currentNum)})
				currentNum = 0
				currentNumSet = false
			}
			continue
		}
		if unicode.IsDigit(c) {
			digit, _ := strconv.Atoi(string(c))
			currentNum = currentNum*10 + digit
			currentNumSet = true
		} else {
			if currentNumSet {
				currentNumSet = false
				currentNum = 0
				tokens.push(token{kind: "num", value: strconv.Itoa(currentNum)})
			}
			switch c {
			case '+':
				tokens.push(token{kind: "op", value: "+"})
			case '-':
				tokens.push(token{kind: "op", value: "-"})
			case '*':
				tokens.push(token{kind: "op", value: "*"})
			case '/':
				tokens.push(token{kind: "op", value: "/"})
			default:
				panic("Unknown token: " + string(c))
			}
		}
	}
	if currentNumSet {
		tokens.push(token{kind: "num", value: strconv.Itoa(currentNum)})
	}
	return tokens
}

// Compiles a postfix expression to navm IR
func toIR(expr string) *navm.IR {
	tokens := tokenize(expr)
	operands := &q.Queue{}
	ir := navm.NewIR()

	// We are going to walk through the expression as if we would interpret it
	// but instead at each step will create IR, get the registers and store their numbers
	// in the operands queue

	for !tokens.empty() {
		t := tokens.takeFirst()
		switch t.kind {
		case "num":
			val, _ := strconv.Atoi(t.value)
			vreg := ir.NewVirtualRegister()
			ir.MoveConstant(vreg, val)
			operands.Push(vreg.Value())

		case "op":
			// pop two numbers off stack, apply operation, push result back onto stack
			if operands.Len() < 2 {
				panic("Not enough operands")
			}
			op2 := operands.PopLast()
			op1 := operands.PopLast()
			switch t.value {
			case "+":
				vreg := ir.NewVirtualRegister()
				ir.AddRegisters(vreg, navm.MakeVirtualRegister(op1), navm.MakeVirtualRegister(op2))
				operands.Push(vreg.Value())
			case "-":
				vreg := ir.NewVirtualRegister()
				ir.SubRegisters(vreg, navm.MakeVirtualRegister(op1), navm.MakeVirtualRegister(op2))
				operands.Push(vreg.Value())
			case "*":
				vreg := ir.NewVirtualRegister()
				ir.MultRegisters(vreg, navm.MakeVirtualRegister(op1), navm.MakeVirtualRegister(op2))
				operands.Push(vreg.Value())
			case "/":
				vreg := ir.NewVirtualRegister()
				ir.DivRegisters(vreg, navm.MakeVirtualRegister(op1), navm.MakeVirtualRegister(op2))
				operands.Push(vreg.Value())
			default:
				panic("Unknown operator: " + t.value)
			}
		default:
			panic("Unknown token kind: " + t.kind)
		}
	}
	return ir
}

func interpret(expr string) int {
	ir := toIR(expr)
	return navm.Interpret(ir)
}

func compile(expr string) string {
	ir := toIR(expr)
	return navm.Compile(ir, navm.AARCH64_MACOS_NONE)
}

// parse cli args
func main() {
	args := os.Args
	if len(args) < 2 {
		panic("Usage: navm [compile/interpret] <postfix-expression>")
	}
	switch args[1] {
	case "compile":
		ir := toIR(args[2])
		result := navm.Compile(ir, navm.AARCH64_MACOS_NONE)
		println(result)
	case "interpret":
		ir := toIR(args[2])
		result := navm.Interpret(ir)
		println(result)
	}
}
