package stop

import "github.com/vcokltfre/stop/stop/instructions"

func Compile(instrs []instructions.Instruction) []byte {
	out := []byte{}

	for _, instr := range instrs {
		out = append(out, instr.Emit()...)
	}

	return out
}
