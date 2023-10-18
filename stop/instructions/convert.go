package instructions

func i64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	for j := 0; j < 8; j++ {
		b[j] = byte(i >> uint64(j*8))
	}
	return b
}

func u16ToBytes(i uint16) []byte {
	b := make([]byte, 2)
	for j := 0; j < 2; j++ {
		b[j] = byte(i >> uint16(j*8))
	}
	return b
}

// func bytesToI64(b []byte) int64 {
// 	var i int64
// 	for j := 0; j < 8; j++ {
// 		i |= int64(b[j]) << uint64(j*8)
// 	}
// 	return i
// }
