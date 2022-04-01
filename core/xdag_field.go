package core

import "encoding/binary"

type FieldType byte

const (
	// nonce字段
	XDAG_FIELD_NONCE FieldType = iota
	// 头部字段
	XDAG_FIELD_HEAD
	// 输入
	XDAG_FIELD_IN
	// 输入
	XDAG_FIELD_OUT
	// 输入签名
	XDAG_FIELD_SIGN_IN
	// 输出签名
	XDAG_FIELD_SIGN_OUT
	XDAG_FIELD_PUBLIC_KEY_0
	XDAG_FIELD_PUBLIC_KEY_1
	XDAG_FIELD_HEAD_TEST
	XDAG_FIELD_REMARK
	XDAG_FIELD_RESERVE1
	XDAG_FIELD_RESERVE2
	XDAG_FIELD_RESERVE3
	XDAG_FIELD_RESERVE4
	XDAG_FIELD_RESERVE5
	XDAG_FIELD_RESERVE6
)

var EmptyHashOrField [XDAG_FIELD_SIZE]byte

type XdagField struct {
	Data [XDAG_FIELD_SIZE]byte
	Sum  uint64
	Type FieldType
}

func (x XdagField) GetSum() uint64 {
	x.Sum = 0
	for i := 0; i < 4; i++ {
		x.Sum += binary.LittleEndian.Uint64(x.Data[i*8 : (i+1)*8])
	}
	return x.Sum
}
