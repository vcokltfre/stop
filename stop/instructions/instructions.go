package instructions

const (
	IHeaderHlt uint8 = 0x00 // Halt
	IHeaderDbg uint8 = 0x01 // Debug

	IHeaderMovLiteral  uint8 = 0x08 // Move value
	IHeaderMovRegister uint8 = 0x09 // Move register

	IHeaderPush uint8 = 0x10 // Push value
	IHeaderDup  uint8 = 0x11 // Duplicate value
	IHeaderDrop uint8 = 0x12 // Drop value
	IHeaderSwap uint8 = 0x13 // Swap values

	IHeaderLd uint8 = 0x20 // Load value
	IHeaderSt uint8 = 0x21 // Store value

	IHeaderAdd uint8 = 0x30 // Add
	IHeaderSub uint8 = 0x31 // Subtract
	IHeaderMul uint8 = 0x32 // Multiply
	IHeaderDiv uint8 = 0x33 // Divide
	IHeaderMod uint8 = 0x34 // Modulo

	IHeaderLabel uint8 = 0xA0 // Label
	IHeaderCall  uint8 = 0xA1 // Call
	IHeaderJmp   uint8 = 0xA2 // Jump
	IHeaderJmpZ  uint8 = 0xA3 // Jump if zero
	IHeaderJmpNZ uint8 = 0xA4 // Jump if not zero
	IHeaderJmpP  uint8 = 0xA5 // Jump if positive
	IHeaderJmpN  uint8 = 0xA6 // Jump if negative
	IHeaderRet   uint8 = 0xA7 // Return

	IHeaderPutN uint8 = 0xB0 // Put number
	IHeaderPutC uint8 = 0xB1 // Put character
)

type Instruction interface {
	Emit() []byte
}
