package stop

import (
	"fmt"

	"github.com/vcokltfre/stop/stop/instructions"
)

const (
	STACK_SIZE      = 1024
	CALL_STACK_SIZE = 64
	DEBUG           = false
)

type VM struct {
	stack        []int64
	stackTop     int
	callStack    []int
	callStackTop int
	program      []byte
	registers    []int64

	jumps map[uint16]int
	index int
}

func (v *VM) debug(data ...any) {
	if DEBUG {
		fmt.Println(append([]any{"[DEBUG]", fmt.Sprintf("%x", v.index)}, data...)...)
	}
}

func (v *VM) buildJumps() {
	index := 0

	for index < len(v.program) {
		curr := v.program[index]
		switch curr {
		case instructions.IHeaderHlt:
			index += instructions.ISizeHlt
		case instructions.IHeaderDbg:
			index += instructions.ISizeDbg
		case instructions.IHeaderMovLiteral:
			index += instructions.ISizeMovLiteral
		case instructions.IHeaderMovRegister:
			index += instructions.ISizeMovRegister
		case instructions.IHeaderPush:
			index += instructions.ISizePush
		case instructions.IHeaderDup:
			index += instructions.ISizeDup
		case instructions.IHeaderDrop:
			index += instructions.ISizeDrop
		case instructions.IHeaderSwap:
			index += instructions.ISizeSwap
		case instructions.IHeaderLd:
			index += instructions.ISizeLd
		case instructions.IHeaderSt:
			index += instructions.ISizeSt
		case instructions.IHeaderAdd:
			index += instructions.ISizeAdd
		case instructions.IHeaderSub:
			index += instructions.ISizeSub
		case instructions.IHeaderMul:
			index += instructions.ISizeMul
		case instructions.IHeaderDiv:
			index += instructions.ISizeDiv
		case instructions.IHeaderMod:
			index += instructions.ISizeMod
		case instructions.IHeaderLabel:
			index += instructions.ISizeLabel
			v.jumps[uint16(v.program[index-1])<<8|uint16(v.program[index-2])] = index
		case instructions.IHeaderCall:
			index += instructions.ISizeCall
		case instructions.IHeaderJmp:
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpZ:
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpNZ:
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpP:
			index += instructions.ISizeJmp
		case instructions.IHeaderJmpN:
			index += instructions.ISizeJmp
		case instructions.IHeaderRet:
			index += instructions.ISizeRet
		case instructions.IHeaderPutN:
			index += instructions.ISizePutN
		case instructions.IHeaderPutC:
			index += instructions.ISizePutC
		default:
			panic("invalid instruction: " + fmt.Sprintf("%x", curr))
		}
	}
}

func (v *VM) stackPush(val int64) {
	if v.stackTop >= STACK_SIZE-1 {
		panic("stack overflow")
	}

	v.debug("pushing", v.stackTop, val)

	v.stackTop++
	v.stack[v.stackTop] = val
}

func (v *VM) stackPop() int64 {
	if v.stackTop < 0 {
		panic("stack underflow")
	}

	v.debug("popping", v.stackTop, v.stack[v.stackTop])

	v.stackTop--
	return v.stack[v.stackTop+1]
}

func (v *VM) callStackPush(val int) {
	if v.callStackTop >= CALL_STACK_SIZE-1 {
		panic("call stack overflow")
	}

	v.callStackTop++
	v.callStack[v.callStackTop] = val
}

func (v *VM) callStackPop() int {
	if v.callStackTop < 0 {
		panic("call stack underflow")
	}

	v.callStackTop--
	return v.callStack[v.callStackTop+1]
}

func (v *VM) getReg() uint8 {
	return v.program[v.index]
}

func (v *VM) getU16() uint16 {
	return uint16(v.program[v.index+1])<<8 | uint16(v.program[v.index])
}

func (v *VM) getI64() int64 {
	var i int64
	for j := 0; j < 8; j++ {
		i |= int64(v.program[v.index+j]) << uint64(j*8)
	}
	return i
}

func (v *VM) instDebug() {
	fmt.Println("TODO: implement debug instruction")
}

func (v *VM) instMovLiteral() {
	v.index += 1
	reg := v.getReg()
	v.index += 1
	val := v.getI64()
	v.registers[int(reg)] = val
	v.index += 8
}

func (v *VM) instMovRegister() {
	v.index += 1
	reg := v.getReg()
	v.index += 1
	src := v.getReg()
	v.registers[int(reg)] = v.registers[int(src)]
	v.index += 1
}

func (v *VM) instPush() {
	v.index++
	v.stackPush(v.getI64())
	v.index += 8
}

func (v *VM) instDup() {
	val := v.stackPop()
	v.stackPush(val)
	v.stackPush(val)
	v.index += instructions.ISizeDup
}

func (v *VM) instDrop() {
	v.stackPop()
	v.index += instructions.ISizeDrop
}

func (v *VM) instSwap() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a)
	v.stackPush(b)
	v.index += instructions.ISizeSwap
}

func (v *VM) instLd() {
	v.index += 1
	reg := v.getReg()
	v.stackPush(v.registers[int(reg)])
	v.index += 1
}

func (v *VM) instSt() {
	v.index += 1
	reg := v.getReg()
	v.registers[int(reg)] = v.stackPop()
	v.index += 1
}

func (v *VM) instAdd() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a + b)
	v.index += instructions.ISizeAdd
}

func (v *VM) instSub() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a - b)
	v.index += instructions.ISizeSub
}

func (v *VM) instMul() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a * b)
	v.index += instructions.ISizeMul
}

func (v *VM) instDiv() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a / b)
	v.index += instructions.ISizeDiv
}

func (v *VM) instMod() {
	a := v.stackPop()
	b := v.stackPop()
	v.stackPush(a % b)
	v.index += instructions.ISizeMod
}

func (v *VM) instJmp() {
	v.index += 1
	addr := v.getU16()
	v.index = v.jumps[addr]
}

func (v *VM) instJmpZ() {
	v.index += 1
	addr := v.getU16()
	val := v.stackPop()
	if val == 0 {
		v.index = v.jumps[addr]
	} else {
		v.index += 2
	}
}

func (v *VM) instJmpNZ() {
	v.index += 1
	addr := v.getU16()
	val := v.stackPop()
	if val != 0 {
		v.index = v.jumps[addr]
	} else {
		v.index += 2
	}
}

func (v *VM) instJmpP() {
	v.index += 1
	addr := v.getU16()
	val := v.stackPop()
	if val > 0 {
		v.index = v.jumps[addr]
	} else {
		v.index += 2
	}
}

func (v *VM) instJmpN() {
	v.index += 1
	addr := v.getU16()
	val := v.stackPop()
	if val < 0 {
		v.index = v.jumps[addr]
	} else {
		v.index += 2
	}
}

func (v *VM) instRet() {
	v.index = v.callStackPop()
}

func (v *VM) instPutN() {
	v.index += 1
	val := v.stackPop()
	fmt.Println(val)
}

func (v *VM) instPutC() {
	v.index += 1
	val := v.stackPop()
	fmt.Printf("%c", val)
}

func (v *VM) step() bool {
	if v.index >= len(v.program) {
		return true
	}

	curr := v.program[v.index]
	switch curr {
	case instructions.IHeaderHlt:
		return true
	case instructions.IHeaderDbg:
		v.instDebug()
		v.index += instructions.ISizeDbg
	case instructions.IHeaderMovLiteral:
		v.debug("mov literal")
		v.instMovLiteral()
	case instructions.IHeaderMovRegister:
		v.debug("mov register")
		v.instMovRegister()
	case instructions.IHeaderPush:
		v.debug("push")
		v.instPush()
	case instructions.IHeaderDup:
		v.debug("dup")
		v.instDup()
	case instructions.IHeaderDrop:
		v.debug("drop")
		v.instDrop()
	case instructions.IHeaderSwap:
		v.debug("swap")
		v.instSwap()
	case instructions.IHeaderLd:
		v.debug("ld")
		v.instLd()
	case instructions.IHeaderSt:
		v.debug("st")
		v.instSt()
	case instructions.IHeaderAdd:
		v.debug("add")
		v.instAdd()
	case instructions.IHeaderSub:
		v.debug("sub")
		v.instSub()
	case instructions.IHeaderMul:
		v.debug("mul")
		v.instMul()
	case instructions.IHeaderDiv:
		v.debug("div")
		v.instDiv()
	case instructions.IHeaderMod:
		v.debug("mod")
		v.instMod()
	case instructions.IHeaderLabel:
		v.debug("label")
		v.index += instructions.ISizeLabel
	case instructions.IHeaderCall:
		v.debug("call")
		v.index += 1
		v.callStackPush(v.index + 2)
		v.index = v.jumps[v.getU16()]
	case instructions.IHeaderJmp:
		v.debug("jmp")
		v.instJmp()
	case instructions.IHeaderJmpZ:
		v.debug("jmpz")
		v.instJmpZ()
	case instructions.IHeaderJmpNZ:
		v.debug("jmpnz")
		v.instJmpNZ()
	case instructions.IHeaderJmpP:
		v.debug("jmpp")
		v.instJmpP()
	case instructions.IHeaderJmpN:
		v.debug("jmpn")
		v.instJmpN()
	case instructions.IHeaderRet:
		v.debug("ret")
		v.instRet()
	case instructions.IHeaderPutN:
		v.debug("putn")
		v.instPutN()
	case instructions.IHeaderPutC:
		v.debug("putc")
		v.instPutC()
	default:
		panic("invalid instruction: " + fmt.Sprintf("%x", curr) + " at " + fmt.Sprintf("%x (%d)", v.index, v.index))
	}
	return false
}

func (v *VM) Run(code []byte) {
	v.stack = make([]int64, STACK_SIZE)
	v.callStack = make([]int, CALL_STACK_SIZE)
	v.program = code
	v.jumps = make(map[uint16]int)
	v.registers = make([]int64, 16)

	v.buildJumps()

	for {
		stop := v.step()
		if stop {
			break
		}
	}
}
