package instructions

type InstLabel struct {
	Label uint16
}

func (i InstLabel) Emit() []byte {
	return append([]byte{IHeaderLabel}, u16ToBytes(i.Label)...)
}

type InstCall struct {
	Label uint16
}

func (i InstCall) Emit() []byte {
	return append([]byte{IHeaderCall}, u16ToBytes(i.Label)...)
}

type InstJmp struct {
	Label uint16
}

func (i InstJmp) Emit() []byte {
	return append([]byte{IHeaderJmp}, u16ToBytes(i.Label)...)
}

type InstJmpZ struct {
	Label uint16
}

func (i InstJmpZ) Emit() []byte {
	return append([]byte{IHeaderJmpZ}, u16ToBytes(i.Label)...)
}

type InstJmpNZ struct {
	Label uint16
}

func (i InstJmpNZ) Emit() []byte {
	return append([]byte{IHeaderJmpNZ}, u16ToBytes(i.Label)...)
}

type InstJmpP struct {
	Label uint16
}

func (i InstJmpP) Emit() []byte {
	return append([]byte{IHeaderJmpP}, u16ToBytes(i.Label)...)
}

type InstJmpN struct {
	Label uint16
}

func (i InstJmpN) Emit() []byte {
	return append([]byte{IHeaderJmpN}, u16ToBytes(i.Label)...)
}

type InstRet struct{}

func (i InstRet) Emit() []byte {
	return []byte{IHeaderRet}
}
