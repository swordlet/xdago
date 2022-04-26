package utils

import "encoding/binary"

func KeyStartWith(source, prefix []byte) bool {
	if len(prefix) > len(source) {
		return false
	}

	for i := 0; i < len(prefix); i++ {
		if prefix[i] != source[i] {
			return false
		}
	}
	return true
}

func MergeBytes(array ...[]byte) []byte {
	var total int
	for _, arr := range array {
		total += len(arr)
	}
	res := make([]byte, total)
	var length int
	for _, arr := range array {
		copy(res[length:], arr[:])
		length += len(arr)
	}
	return res
}

func Copy2(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func U64ToBytes(u uint64, order binary.ByteOrder) []byte {
	var b [8]byte
	order.PutUint64(b[:], u)
	return b[:]
}
