package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"unicode"
)

func Amount2xdag(amount uint64) float64 {
	integer := amount >> 32
	temp := amount - (integer << 32)
	decimal := float64(temp) / math.Pow(2, 32)

	return math.Round((float64(integer)+decimal)*100) / 100
}

func IsAsciiPrintable(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII || !unicode.IsPrint(c) {
			return false
		}
	}
	return true
}

func Hash2String(h [32]byte) string {
	return fmt.Sprintf("%016x%016x%016x%016x",
		binary.LittleEndian.Uint64(h[24:]),
		binary.LittleEndian.Uint64(h[16:24]),
		binary.LittleEndian.Uint64(h[8:16]),
		binary.LittleEndian.Uint64(h[:8]))
}

func Type2String(i uint64) string {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], i)
	var s string
	for _, k := range b {
		s += fmt.Sprintf("%x", k&0x0f)
		s += fmt.Sprintf("%x", (k>>4)&0x0f)
	}
	return s
}
