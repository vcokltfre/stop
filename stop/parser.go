package stop

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vcokltfre/stop/stop/instructions"
)

func isIdent(val string) bool {
	for _, r := range val {
		if !(r >= 'a' && r <= 'z') {
			return false
		}
	}

	return true
}

func isReg(val string) (bool, int) {
	if len(val) < 2 {
		return false, 0
	}

	if val[0] != 'r' {
		return false, 0
	}

	if len(val) == 2 {
		return val[1] >= '0' && val[1] <= '9', int(val[1] - '0')
	}

	if len(val) == 3 {
		n := val[1:2]
		v, err := strconv.Atoi(n)
		if err != nil {
			return false, 0
		}

		return v <= 15, v
	}

	return false, 0
}

func isLiteral(val string) (bool, int64) {
	v, err := strconv.ParseInt(val, 10, 64)
	return err == nil, v
}

func cleanArray(arr []string) []string {
	out := []string{}

	for _, val := range arr {
		val = strings.TrimSpace(val)

		if len(val) == 0 {
			continue
		}

		out = append(out, val)
	}

	return out
}

func Parse(code string) ([]instructions.Instruction, error) {
	lines := strings.Split(code, "\n")
	jumps := map[string]int{}

	insts := []instructions.Instruction{}

	for i, line := range lines {
		clean := strings.TrimSpace(line)

		if len(clean) == 0 {
			continue
		}

		if clean[0] == ':' {
			err := func(msg string) error {
				return fmt.Errorf("error on line %d: %s", i+1, msg)
			}

			if len(jumps) >= (1 << 16) {
				return nil, err("too many labels")
			}

			label := strings.TrimSpace(clean[1:])
			if len(label) == 0 {
				return nil, err("label must have a name")
			}

			if !isIdent(label) {
				return nil, err("label must be a valid identifier ([a-z]+)")
			}

			if _, ok := jumps[label]; ok {
				return nil, err("label already defined")
			}

			labelId := len(jumps)
			jumps[label] = labelId

			continue
		}
	}

	for i, line := range lines {
		err := func(msg string) error {
			return fmt.Errorf("error on line %d: %s", i+1, msg)
		}

		isJump := func(loc string) (bool, uint16) {
			v, ok := jumps[loc]
			return ok, uint16(v)
		}

		clean := strings.TrimSpace(line)

		if len(clean) == 0 {
			continue
		}

		if clean[0] == ';' {
			continue
		}

		if clean[0] == ':' {
			label := strings.TrimSpace(clean[1:])
			insts = append(insts, instructions.InstLabel{Label: uint16(jumps[label])})
			continue
		}

		parts := cleanArray(strings.Split(clean, " "))

		switch parts[0] {
		case "hlt":
			insts = append(insts, instructions.InstHlt{})
		case "dbg":
			insts = append(insts, instructions.InstDbg{})
		case "mov":
			// Handle movlit and movreg
			if len(parts) != 3 {
				return nil, err("mov must have two arguments")
			}

			r1ok, r1 := isReg(parts[1])
			if !r1ok {
				return nil, err("mov first argument must be a register")
			}

			r2ok, r2 := isReg(parts[2])
			if r2ok {
				insts = append(insts, instructions.InstMovRegister{Register: uint8(r1), Source: uint8(r2)})
				break
			}

			v, e := strconv.ParseInt(parts[2], 10, 64)
			if e != nil {
				return nil, err("mov second argument must be a register or a number")
			}

			insts = append(insts, instructions.InstMovLiteral{Register: uint8(r1), Value: v})
		case "ld":
			if len(parts) != 2 {
				return nil, err("ld must have one argument")
			}

			r1ok, r1 := isReg(parts[1])
			if !r1ok {
				return nil, err("ld argument must be a register")
			}

			insts = append(insts, instructions.InstLd{Register: uint8(r1)})
		case "st":
			if len(parts) != 2 {
				return nil, err("st must have one argument")
			}

			r1ok, r1 := isReg(parts[1])
			if !r1ok {
				return nil, err("st argument must be a register")
			}

			insts = append(insts, instructions.InstSt{Register: uint8(r1)})
		case "push":
			if len(parts) != 2 {
				return nil, err("push must have one argument")
			}

			vOk, v := isLiteral(parts[1])
			if !vOk {
				return nil, err("push argument must be a number")
			}

			insts = append(insts, instructions.InstPush{Value: v})
		case "dup":
			if len(parts) != 1 {
				return nil, err("dup must have no arguments")
			}

			insts = append(insts, instructions.InstDup{})
		case "drop":
			if len(parts) != 1 {
				return nil, err("drop must have no arguments")
			}

			insts = append(insts, instructions.InstDrop{})
		case "swap":
			if len(parts) != 1 {
				return nil, err("swap must have no arguments")
			}

			insts = append(insts, instructions.InstSwap{})
		case "add":
			if len(parts) != 1 {
				return nil, err("add must have no arguments")
			}

			insts = append(insts, instructions.InstAdd{})
		case "sub":
			if len(parts) != 1 {
				return nil, err("sub must have no arguments")
			}

			insts = append(insts, instructions.InstSub{})
		case "mul":
			if len(parts) != 1 {
				return nil, err("mul must have no arguments")
			}

			insts = append(insts, instructions.InstMul{})
		case "div":
			if len(parts) != 1 {
				return nil, err("div must have no arguments")
			}

			insts = append(insts, instructions.InstDiv{})
		case "mod":
			if len(parts) != 1 {
				return nil, err("mod must have no arguments")
			}

			insts = append(insts, instructions.InstMod{})
		case "call":
			if len(parts) != 2 {
				return nil, err("call must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("call argument must be a label")
			}

			insts = append(insts, instructions.InstCall{Label: jLoc})
		case "jmp":
			if len(parts) != 2 {
				return nil, err("jmp must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("jmp argument must be a label")
			}

			insts = append(insts, instructions.InstJmp{Label: jLoc})
		case "jmpz":
			if len(parts) != 2 {
				return nil, err("jmpz must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("jmpz argument must be a label")
			}

			insts = append(insts, instructions.InstJmpZ{Label: jLoc})
		case "jmpnz":
			if len(parts) != 2 {
				return nil, err("jmpnz must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("jmpnz argument must be a label")
			}

			insts = append(insts, instructions.InstJmpNZ{Label: jLoc})
		case "jmpp":
			if len(parts) != 2 {
				return nil, err("jmpp must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("jmpp argument must be a label")
			}

			insts = append(insts, instructions.InstJmpP{Label: jLoc})
		case "jmpn":
			if len(parts) != 2 {
				return nil, err("jmpn must have one argument")
			}

			jOk, jLoc := isJump(parts[1])
			if !jOk {
				return nil, err("jmpn argument must be a label")
			}

			insts = append(insts, instructions.InstJmpN{Label: jLoc})
		case "putn":
			if len(parts) != 1 {
				return nil, err("putn must have no arguments")
			}

			insts = append(insts, instructions.InstPutN{})
		case "putc":
			if len(parts) != 1 {
				return nil, err("putc must have no arguments")
			}

			insts = append(insts, instructions.InstPutC{})
		case "ret":
			if len(parts) != 1 {
				return nil, err("ret must have no arguments")
			}

			insts = append(insts, instructions.InstRet{})
		default:
			return nil, err("unknown instruction")
		}
	}

	return insts, nil
}
