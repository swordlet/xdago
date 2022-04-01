package core

import (
	"encoding/binary"
	"xdago/log"
)

var EmptyXdagBlock [XDAG_BLOCK_SIZE]byte

type XdagBlock struct {
	Data   [XDAG_BLOCK_SIZE]byte
	Sum    uint64
	Fields [XDAG_BLOCK_FIELDS]XdagField
}

func NewXdagBlock(data []byte) *XdagBlock {
	if len(data) != XDAG_BLOCK_SIZE {
		log.Crit("new xdag block, data size error", log.Ctx{"len": len(data)})
	}
	xb := XdagBlock{}
	copy(xb.Data[:], data)
	for i := 0; i < XDAG_BLOCK_FIELDS; i++ {
		copy(xb.Fields[i].Data[:], data[i*XDAG_FIELD_SIZE:(i+1)*XDAG_FIELD_SIZE])
	}
	for i := 0; i < XDAG_BLOCK_FIELDS; i++ {
		xb.Sum += xb.Fields[i].GetSum()
		xb.Fields[i].Type = xb.getMsgCode(i)
	}
	return &xb
}

func (xb XdagBlock) getMsgCode(n int) FieldType {
	t := binary.LittleEndian.Uint64(xb.Data[8:16])
	return FieldType((t >> (n << 2)) & 0x0f)
}

func (xb *XdagBlock) GetData() [XDAG_BLOCK_SIZE]byte {
	if xb.Data == EmptyXdagBlock {
		xb.Sum = 0
		for i := 0; i < XDAG_BLOCK_FIELDS; i++ {
			xb.Sum += xb.Fields[i].GetSum()
			copy(xb.Data[i*XDAG_FIELD_SIZE:(i+1)*XDAG_FIELD_SIZE], xb.Fields[i].Data[:])
		}
	}
	return xb.Data
}

func (xb *XdagBlock) Parse() {
	if xb.Data == EmptyXdagBlock {
		return
	}
	for i := 0; i < XDAG_BLOCK_FIELDS; i++ {
		copy(xb.Fields[i].Data[:], xb.Data[i*XDAG_FIELD_SIZE:(i+1)*XDAG_FIELD_SIZE])
	}
	for i := 0; i < XDAG_BLOCK_FIELDS; i++ {
		xb.Sum += xb.Fields[i].GetSum()
		xb.Fields[i].Type = xb.getMsgCode(i)
	}
}
