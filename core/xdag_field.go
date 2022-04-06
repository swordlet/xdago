package core

import (
	"encoding/binary"
	"xdago/common"
)

type XdagField struct {
	Data common.Field
	Sum  uint64
	Type common.FieldType
}

func (x XdagField) GetSum() uint64 {
	x.Sum = 0
	for i := 0; i < 4; i++ {
		x.Sum += binary.LittleEndian.Uint64(x.Data[i*8 : (i+1)*8])
	}
	return x.Sum
}
