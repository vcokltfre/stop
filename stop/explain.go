package stop

import (
	"fmt"
	"strings"

	"github.com/vcokltfre/stop/stop/instructions"
)

func getReg(data []byte) uint8 {
	return data[0]
}

func getU16(data []byte) uint16 {
	return uint16(data[1])<<8 | uint16(data[0])
}

func getI64(data []byte) int64 {
	var i int64
	for j := 0; j < 8; j++ {
		i |= int64(data[j]) << uint64(j*8)
	}
	return i
}

func Explain(code []byte) {
	index := 0

	for index < len(code) {
		explain := func(name, message string) {
			padding := strings.Repeat(" ", 8-len(name))
			instr := code[index]
			if message == "" {
				fmt.Printf("[%4x] %2x %s%s\n", index, instr, name, padding)
			} else {
				fmt.Printf("[%4x] %2x %s%s %s\n", index, instr, name, padding, message)
			}
		}

		curr := code[index]
		switch curr {
		case instructions.IHeaderHlt:
			explain("HLT", "")
			index += instructions.ISizeHlt
		case instructions.IHeaderDbg:
			explain("DBG", "")
			index += instructions.ISizeDbg
		case instructions.IHeaderMovLiteral:
			explain("MOV", fmt.Sprintf("(literal %d -> register %d)", getI64(code[index+2:index+10]), getReg(code[index+1:index+2])))
			index += instructions.ISizeMovLiteral
		case instructions.IHeaderMovRegister:
			explain("MOV", fmt.Sprintf("(register %d -> register %d)", getReg(code[index+2:index+3]), getReg(code[index+1:index+2])))
			index += instructions.ISizeMovRegister
		case instructions.IHeaderPush:
			explain("PUSH", fmt.Sprintf("(literal %d)", getI64(code[index+1:index+9])))
			index += instructions.ISizePush
		case instructions.IHeaderDup:
			explain("DUP", "")
			index += instructions.ISizeDup
		case instructions.IHeaderDrop:
			explain("DROP", "")
			index += instructions.ISizeDrop
		case instructions.IHeaderSwap:
			explain("SWAP", "")
			index += instructions.ISizeSwap
		case instructions.IHeaderLd:
			explain("LD", fmt.Sprintf("(register %d)", getReg(code[index+1:index+2])))
			index += instructions.ISizeLd
		case instructions.IHeaderSt:
			explain("ST", fmt.Sprintf("(register %d)", getReg(code[index+1:index+2])))
			index += instructions.ISizeSt
		case instructions.IHeaderAdd:
			explain("ADD", "")
			index += instructions.ISizeAdd
		case instructions.IHeaderSub:
			explain("SUB", "")
			index += instructions.ISizeSub
		case instructions.IHeaderMul:
			explain("MUL", "")
			index += instructions.ISizeMul
		case instructions.IHeaderDiv:
			explain("DIV", "")
			index += instructions.ISizeDiv
		case instructions.IHeaderMod:
			explain("MOD", "")
			index += instructions.ISizeMod
		case instructions.IHeaderLabel:
			explain("LABEL", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeLabel
		case instructions.IHeaderCall:
			explain("CALL", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeCall
		case instructions.IHeaderJmp:
			explain("JMP", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpZ:
			explain("JMPZ", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpNZ:
			explain("JMPNZ", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpP:
			explain("JMPP", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpN:
			explain("JMPN", fmt.Sprintf("(label %d)", getU16(code[index+1:index+3])))
			index += instructions.ISizeJmp
		case instructions.IHeaderRet:
			explain("RET", "")
			index += instructions.ISizeRet
		case instructions.IHeaderPutN:
			explain("PUTN", "")
			index += instructions.ISizePutN
		case instructions.IHeaderPutC:
			explain("PUTC", "")
			index += instructions.ISizePutC
		default:
			explain("INVALID", fmt.Sprintf("%x", curr))
			index++
		}
	}
}
