package instructions

type InstMovLiteral struct {
	Register uint8
	Value    int64
}

func (i InstMovLiteral) Emit() []byte {
	return append([]byte{IHeaderMovLiteral, i.Register}, i64ToBytes(i.Value)...)
}

type InstMovRegister struct {
	Register uint8
	Source   uint8
}

func (i InstMovRegister) Emit() []byte {
	return []byte{IHeaderMovRegister, i.Register, i.Source}
}

type InstLd struct {
	Register uint8
}

func (i InstLd) Emit() []byte {
	return []byte{IHeaderLd, i.Register}
}

type InstSt struct {
	Register uint8
}

func (i InstSt) Emit() []byte {
	return []byte{IHeaderSt, i.Register}
}
