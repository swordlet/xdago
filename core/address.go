package core

import (
	"encoding/binary"
	"encoding/hex"
	"xdago/log"
)

type Address struct {
	Data    [XDAG_FIELD_SIZE]byte
	HashLow [XDAG_HASH_SIZE]byte
	Type    FieldType
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
func AddressFromHashLow(hashLow []byte) Address {
	if len(hashLow) != XDAG_HASH_SIZE {
		log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
	}
	adr := Address{
		Parsed: true,
	}
	copy(adr.HashLow[:], hashLow)
	return adr
}

// AddressFromBlock 只用于ref 跟 maxdifflink
func AddressFromBlock(block Block) Address {
	adr := Address{
		Parsed: true,
	}
	copy(adr.HashLow[:], block.GetHashLow())
	return adr
}
func AddressFromType(data []byte, typ FieldType) Address {
	if len(data) != XDAG_FIELD_SIZE {
		log.Crit("address from type, data size error", log.Ctx{"len": len(data)})
	}
	adr := Address{
		Type: typ,
	}
	copy(adr.Data[:], data)
	adr.Parse()
	return adr
}

func AddressFromAmount(hashLow []byte, typ FieldType, amount uint64) Address {
	if len(hashLow) != XDAG_HASH_SIZE {
		log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
	}
	adr := Address{
		Type:   typ,
		Amount: amount,
		Parsed: true,
	}
	copy(adr.HashLow[:], hashLow)
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
	if adr.Data == EmptyHashOrField {
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
