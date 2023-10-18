package instructions

type InstAdd struct{}

func (i InstAdd) Emit() []byte {
	return []byte{IHeaderAdd}
}

type InstSub struct{}

func (i InstSub) Emit() []byte {
	return []byte{IHeaderSub}
}

type InstMul struct{}

func (i InstMul) Emit() []byte {
	return []byte{IHeaderMul}
}

type InstDiv struct{}

func (i InstDiv) Emit() []byte {
	return []byte{IHeaderDiv}
}

type InstMod struct{}

func (i InstMod) Emit() []byte {
	return []byte{IHeaderMod}
}
