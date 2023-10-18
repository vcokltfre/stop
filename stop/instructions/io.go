package instructions

type InstPutN struct{}

func (i InstPutN) Emit() []byte {
	return []byte{IHeaderPutN}
}

type InstPutC struct{}

func (i InstPutC) Emit() []byte {
	return []byte{IHeaderPutC}
}
