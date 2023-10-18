package instructions

type InstHlt struct{}

func (i InstHlt) Emit() []byte {
	return []byte{IHeaderHlt}
}
