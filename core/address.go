package core

import (
	"encoding/binary"
	"encoding/hex"
	"xdago/common"
	"xdago/log"
)

type Address struct {
	Data    common.Field
	HashLow common.Hash
	Type    common.FieldType
	Amount  uint64
	Parsed  bool
}

func AddressFromField(field XdagField) Address {
	adr := Address{
		Type: field.Type,
		Data: field.Data,
	}
	adr.Parse()
	return adr
}

// AddressFromHashLow 只用于ref 跟 maxdifflink
func AddressFromHashLow(hashLow [32]byte) Address {
	if hashLow == common.EmptyHash {
		log.Crit("address from zero hash low")
	}
	adr := Address{
		Parsed: true,
	}
	adr.HashLow = hashLow
	return adr
}

// AddressFromBlock 只用于ref 跟 maxdifflink
func AddressFromBlock(block Block) Address {
	adr := Address{
		Parsed: true,
	}
	adr.HashLow = block.GetHashLow()
	return adr
}
func AddressFromType(data [32]byte, typ common.FieldType) Address {
	if data == common.EmptyField {
		log.Crit("address from type, zero hash low")
	}
	adr := Address{
		Type: typ,
	}
	adr.Data = data
	adr.Parse()
	return adr
}

func AddressFromAmount(hashLow [32]byte, typ common.FieldType, amount uint64) Address {
	if hashLow == common.EmptyHash {
		log.Crit("address from amount, zero hash low")
	}
	adr := Address{
		Type:   typ,
		Amount: amount,
		Parsed: true,
	}
	adr.HashLow = hashLow
	return adr
}

func (adr *Address) Parse() {
	if !adr.Parsed {
		copy(adr.HashLow[8:32], adr.Data[0:24])
		adr.Amount = binary.LittleEndian.Uint64(adr.Data[24:])
		adr.Parsed = true
	}
}

func (adr *Address) GetData() []byte {
	if adr.Data == common.EmptyField {
		binary.LittleEndian.PutUint64(adr.Data[24:], adr.Amount)
		copy(adr.Data[0:24], adr.HashLow[8:])
	}
	return adr.Data[:]
}

func (adr *Address) GetAmount() uint64 {
	adr.Parse()
	return adr.Amount
}

func (adr *Address) GetHashLow() []byte {
	adr.Parse()
	return adr.HashLow[:]
}

func (adr Address) ToString() string {
	return "Block(A) Hash[" + hex.EncodeToString(adr.HashLow[:]) + "]"
}
