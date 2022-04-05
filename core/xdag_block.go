package core

import (
	"encoding/binary"
	"xdago/common"
	"xdago/log"
)

type XdagBlock struct {
	Data   [common.XDAG_BLOCK_SIZE]byte
	Sum    uint64
	Fields [common.XDAG_BLOCK_FIELDS]XdagField
}

func NewXdagBlock(data []byte) *XdagBlock {
	if len(data) != common.XDAG_BLOCK_SIZE {
		log.Crit("new xdag block, data size error", log.Ctx{"len": len(data)})
	}
	xb := XdagBlock{}
	copy(xb.Data[:], data)
	for i := 0; i < common.XDAG_BLOCK_FIELDS; i++ {
		copy(xb.Fields[i].Data[:], data[i*common.XDAG_FIELD_SIZE:(i+1)*common.XDAG_FIELD_SIZE])
	}
	for i := 0; i < common.XDAG_BLOCK_FIELDS; i++ {
		xb.Sum += xb.Fields[i].GetSum()
		xb.Fields[i].Type = xb.getMsgCode(i)
	}
	return &xb
}

func (xb XdagBlock) getMsgCode(n int) common.FieldType {
	t := binary.LittleEndian.Uint64(xb.Data[8:16])
	return common.FieldType((t >> (n << 2)) & 0x0f)
}

func (xb *XdagBlock) GetData() [common.XDAG_BLOCK_SIZE]byte {
	if xb.Data == common.EmptyXdagBlock {
		xb.Sum = 0
		for i := 0; i < common.XDAG_BLOCK_FIELDS; i++ {
			xb.Sum += xb.Fields[i].GetSum()
			copy(xb.Data[i*common.XDAG_FIELD_SIZE:(i+1)*common.XDAG_FIELD_SIZE], xb.Fields[i].Data[:])
		}
	}
	return xb.Data
}

func (xb *XdagBlock) Parse() {
	if xb.Data == common.EmptyXdagBlock {
		return
	}
	for i := 0; i < common.XDAG_BLOCK_FIELDS; i++ {
		copy(xb.Fields[i].Data[:], xb.Data[i*common.XDAG_FIELD_SIZE:(i+1)*common.XDAG_FIELD_SIZE])
	}
	for i := 0; i < common.XDAG_BLOCK_FIELDS; i++ {
		xb.Sum += xb.Fields[i].GetSum()
		xb.Fields[i].Type = xb.getMsgCode(i)
	}
}
