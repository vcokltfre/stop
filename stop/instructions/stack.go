package instructions

type InstPush struct {
	Value int64
}

func (i InstPush) Emit() []byte {
	return append([]byte{IHeaderPush}, i64ToBytes(i.Value)...)
}

type InstDup struct{}

func (i InstDup) Emit() []byte {
	return []byte{IHeaderDup}
}

type InstDrop struct{}

func (i InstDrop) Emit() []byte {
	return []byte{IHeaderDrop}
}

type InstSwap struct{}

func (i InstSwap) Emit() []byte {
	return []byte{IHeaderSwap}
}
