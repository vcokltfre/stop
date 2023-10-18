package instructions

type InstDbg struct{}

func (i InstDbg) Emit() []byte {
	return []byte{IHeaderDbg}
}
