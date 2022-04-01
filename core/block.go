package core

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"xdago/config"
	"xdago/crypto"
	"xdago/log"
	"xdago/secp256k1"
	"xdago/utils"
)

const MAX_LINKS = 15

type Block struct {
	info                *BlockInfo
	IsSaved             bool
	Parsed              bool
	xdagBlock           *XdagBlock
	TransportHeader     uint64
	Inputs              []Address
	Outputs             []Address
	PubKeys             []secp256k1.PublicKey
	InSigs              [][XDAG_FIELD_SIZE*2 + 1]byte // last byte is signature field index in block
	OutSig              [XDAG_FIELD_SIZE * 2]byte
	Nonce               [XDAG_FIELD_SIZE]byte
	tempLength          int
	PreTopCandidate     bool
	PreTopCandidateDiff big.Int
	//IsOurs            bool
	//encoded           []byte
}

var EmptyXdagSignature [XDAG_FIELD_SIZE * 2]byte

func (b *Block) SetXdagBlock(xdagBlock *XdagBlock) {
	b.xdagBlock = xdagBlock
}

func (b Block) IsEmpty() bool {
	return b.info.Timestamp == 0
}

func NewBlock(config *config.Config, timestamp uint64, links []Address, pending []Address,
	mining bool, keys []secp256k1.PublicKey, remark string, defKeyIndex int) Block {
	b := Block{
		Parsed: true,
		info: &BlockInfo{
			Timestamp: timestamp,
		},
	}
	var length int
	b.setType(config.XdagFieldHeader(), length)
	length += 1

	for _, address := range links {
		typ := address.Type
		b.setType(typ, length)
		length += 1
		if typ == XDAG_FIELD_OUT {
			b.Outputs = append(b.Outputs, address)
		} else {
			b.Inputs = append(b.Inputs, address)
		}
	}

	for _, address := range pending {
		b.setType(XDAG_FIELD_OUT, length)
		length += 1
		b.Outputs = append(b.Outputs, address)
	}
	remark = strings.TrimSpace(remark)
	if len(remark) > 0 && len(remark) < 33 && utils.IsAsciiPrintable(remark) {
		b.setType(XDAG_FIELD_REMARK, length)
		length += 1
		strBytes := []byte(remark)
		copy(b.info.Remark[:], strBytes)
	}

	for _, key := range keys {
		var typ FieldType
		if key.Y().Bit(0) == 0 {
			typ = XDAG_FIELD_PUBLIC_KEY_0 // even
		} else {
			typ = XDAG_FIELD_PUBLIC_KEY_1 // odd
		}
		b.setType(typ, length)
		length += 1
		b.PubKeys = append(b.PubKeys, key)
	}

	for i := 0; i < len(keys); i++ {
		if i != defKeyIndex {
			b.setType(XDAG_FIELD_SIGN_IN, length)
			length += 1
			b.setType(XDAG_FIELD_SIGN_IN, length)
			length += 1
		} else {
			b.setType(XDAG_FIELD_SIGN_OUT, length)
			length += 1
			b.setType(XDAG_FIELD_SIGN_OUT, length)
			length += 1
		}
	}

	if defKeyIndex < 0 {
		b.setType(XDAG_FIELD_SIGN_OUT, length)
		length += 1
		b.setType(XDAG_FIELD_SIGN_OUT, length)
		length += 1
	}

	if mining {
		b.setType(XDAG_FIELD_SIGN_IN, MAX_LINKS)
	}

	return b
}

func NewBlockFromInfo(info *BlockInfo) Block {
	return Block{
		info:    info,
		IsSaved: true,
		Parsed:  true,
	}
}

func NewBlockFromXdag(b *XdagBlock) Block {
	block := Block{
		xdagBlock: b,
	}
	block.Parse()
	return block
}

func (b *Block) GetHashLow() []byte {
	if b.info.HashLow == EmptyHashOrField {
		h := b.GetHash()
		copy(b.info.HashLow[8:], h[8:])
	}
	return b.info.HashLow[:]
}

func (b Block) GetHash() [XDAG_HASH_SIZE]byte {
	if b.info.Hash == EmptyHashOrField {
		b.info.Hash = b.calcHash()
	}
	return b.info.Hash
}

func (b *Block) calcHash() [XDAG_HASH_SIZE]byte {
	if b.xdagBlock == nil {
		b.xdagBlock = NewXdagBlock(b.ToBytes())
	}
	return crypto.HashTwice(b.xdagBlock.Data[:])
}

//RecalcHash 重计算 避免矿工挖矿发送share时直接更新 hash
func (b *Block) RecalcHash() [XDAG_HASH_SIZE]byte {
	b.xdagBlock = NewXdagBlock(b.ToBytes())
	return crypto.HashTwice(b.xdagBlock.Data[:])
}

func (b *Block) ToBytes() []byte {
	w := b.getEncodedBody()
	if w.Error() != nil {
		log.Crit("encode block body error", log.Ctx{"err": w.Error().Error()})
	}
	for _, sig := range b.InSigs {

		w.WriteBytes(sig[:XDAG_FIELD_SIZE*2])
	}
	if b.OutSig != EmptyXdagSignature {
		w.WriteBytes(b.OutSig[:])
	}

	length := w.Length()
	b.tempLength = length
	if length == XDAG_BLOCK_FIELDS {
		if w.Error() != nil {
			log.Crit("block to bytes error", log.Ctx{"err": w.Error().Error()})
		}
		return w.BytesUncheck()
	}

	res := XDAG_BLOCK_FIELDS - 1 - length
	for i := 0; i < res; i++ {
		w.WriteBytes(EmptyHashOrField[:])
	}
	w.WriteBytes(b.Nonce[:])
	if w.Error() != nil {
		log.Crit("block to bytes error", log.Ctx{"err": w.Error().Error()})
	}
	return w.BytesUncheck()
}

// block bytes without signature
func (b Block) getEncodedBody() *BlockWriter {
	w := NewBlockWriter(XDAG_BLOCK_SIZE)
	w.WriteBytes(b.getEncodedHeader())
	all := append(b.Inputs, b.Outputs...)
	for _, link := range all {
		w.WriteBytes(link.GetData())
	}
	if b.info.Remark != EmptyHashOrField {
		w.WriteBytes(b.info.Remark[:])
	}
	for _, publicKey := range b.PubKeys {
		w.WriteBytes(publicKey.SerializeCompressed()[1:33])
	}
	return w
}

func (b Block) getEncodedHeader() []byte {
	var fee [8]byte
	binary.LittleEndian.PutUint64(fee[:], b.info.Fee)

	var timestamp [8]byte
	binary.LittleEndian.PutUint64(timestamp[:], b.info.Timestamp)

	var typ [8]byte
	binary.LittleEndian.PutUint64(typ[:], b.info.Type)

	var transport [8]byte
	return utils.MergeBytes(transport[:], typ[:], timestamp[:], fee[:])
}

func (b *Block) GetXdagBlock() *XdagBlock {
	if b.xdagBlock == nil {
		b.xdagBlock = NewXdagBlock(b.ToBytes())
	}
	return b.xdagBlock
}

// Parse 解析512字节数据
func (b *Block) Parse() {
	if b.Parsed {
		return
	}
	if b.info == nil {
		b.info = &BlockInfo{}
	}
	b.info.Hash = b.calcHash()

	header := b.xdagBlock.Fields[0].Data
	b.TransportHeader = binary.LittleEndian.Uint64(header[:8])
	b.info.Type = binary.LittleEndian.Uint64(header[8:16])
	b.info.Timestamp = binary.LittleEndian.Uint64(header[18:24])
	b.info.Fee = binary.LittleEndian.Uint64(header[24:])

	firtSigIndex := 0
	for i, field := range b.xdagBlock.Fields {
		switch field.Type {
		case XDAG_FIELD_IN:
			b.Inputs = append(b.Inputs, AddressFromField(field))
			break
		case XDAG_FIELD_OUT:
			b.Outputs = append(b.Outputs, AddressFromField(field))
			break
		case XDAG_FIELD_REMARK:
			b.info.Remark = field.Data
			break
		case XDAG_FIELD_PUBLIC_KEY_0:
		case XDAG_FIELD_PUBLIC_KEY_1:
			var key [33]byte
			copy(key[1:], field.Data[:])
			if field.Type == XDAG_FIELD_PUBLIC_KEY_0 {
				key[0] = 0x02
			} else {
				key[0] = 0x03
			}
			pubKey, err := secp256k1.ParsePubKey(key[:])
			if err != nil {
				log.Crit("parse public key error", log.Ctx{"err": err.Error()})
			}
			b.PubKeys = append(b.PubKeys, *pubKey)
			break
		case XDAG_FIELD_SIGN_IN:
		case XDAG_FIELD_SIGN_OUT:
			if firtSigIndex == 0 {
				firtSigIndex = i
			}
			if (i-firtSigIndex)%2 == 0 && i+1 < XDAG_BLOCK_FIELDS {
				if field.Type == XDAG_FIELD_SIGN_IN {
					var insig [XDAG_FIELD_SIZE*2 + 1]byte
					copy(insig[:32], field.Data[:])
					copy(insig[32:64], b.xdagBlock.Fields[i+1].Data[:])
					insig[64] = byte(i)
					b.InSigs = append(b.InSigs, insig)
				} else {
					copy(b.OutSig[:32], field.Data[:])
					copy(b.OutSig[32:], b.xdagBlock.Fields[i+1].Data[:])
				}
			}
			if i == MAX_LINKS && field.Type == XDAG_FIELD_IN {
				b.Nonce = field.Data
			}
			break
		default:

		}
	}
	b.Parsed = true
}

func (b *Block) SignIn(key *secp256k1.PrivateKey) {
	b.sign(key, XDAG_FIELD_SIGN_IN)
}

func (b *Block) SignOut(key *secp256k1.PrivateKey) {
	b.sign(key, XDAG_FIELD_SIGN_OUT)
}

func (b *Block) sign(key *secp256k1.PrivateKey, typ FieldType) {
	encoded := b.ToBytes()
	digest := utils.MergeBytes(encoded, key.PubKey().SerializeCompressed())
	hash := crypto.HashTwice(digest)
	r, s := crypto.EcdsaSign(key, hash[:])
	if typ == XDAG_FIELD_SIGN_OUT {
		copy(b.OutSig[:32], r[:])
		copy(b.OutSig[32:], s[:])
	} else {
		var sig [XDAG_FIELD_SIZE*2 + 1]byte
		copy(sig[:32], r[:])
		copy(sig[32:64], s[:])
		sig[64] = byte(b.tempLength)
		b.InSigs = append(b.InSigs, sig)
	}
}

func (b Block) VerifiedKeys() (res []secp256k1.PublicKey) {
	for _, sig := range b.InSigs {
		digest := b.GetSubRawData(int(sig[64]))
		for _, pubkey := range b.PubKeys {
			hash := crypto.HashTwice(utils.MergeBytes(digest[:], pubkey.SerializeCompressed()))
			if crypto.EcdsaVerify(&pubkey, hash[:], sig[:32], sig[32:64]) {
				res = append(res, pubkey)
			}
		}
	}
	digest := b.GetSubRawData(b.GetOutsigIndex())
	for _, pubkey := range b.PubKeys {
		hash := crypto.HashTwice(utils.MergeBytes(digest[:], pubkey.SerializeCompressed()))
		if crypto.EcdsaVerify(&pubkey, hash[:], b.OutSig[:32], b.OutSig[32:]) {
			res = append(res, pubkey)
		}
	}
	return
}

func (b *Block) setType(typ FieldType, n int) {
	b.info.Type |= uint64(typ) << (n << 2)
}

// GetSubRawData 根据length获取前length个字段的数据 主要用于签名
func (b Block) GetSubRawData(length int) [XDAG_BLOCK_SIZE]byte {
	var res [XDAG_BLOCK_SIZE]byte
	copy(res[:length*XDAG_FIELD_SIZE], b.xdagBlock.Data[:length*XDAG_FIELD_SIZE])
	for i := length; i < XDAG_BLOCK_FIELDS; i++ {
		typ := binary.LittleEndian.Uint64(b.xdagBlock.Data[8:16])
		typeB := FieldType((typ >> (i << 2)) & 0x0f)
		if typeB == XDAG_FIELD_SIGN_IN || typeB == XDAG_FIELD_SIGN_OUT {
			continue
		}
		copy(res[i*XDAG_FIELD_SIZE:(i+1)*XDAG_FIELD_SIZE], b.xdagBlock.Data[i*XDAG_FIELD_SIZE:(i+1)*XDAG_FIELD_SIZE])
	}
	return res
}

//GetOutsigIndex 取输出签名在字段的索引
func (b Block) GetOutsigIndex() int {
	i := 0
	temp := b.info.Type
	for i < XDAG_BLOCK_FIELDS && FieldType(temp&0x0f) != XDAG_FIELD_SIGN_OUT {
		temp = temp >> 4
		i++
	}
	return i
}

func (b Block) GetTimestamp() uint64 {
	return b.info.Timestamp
}

func (b Block) GetType() uint64 {
	return b.info.Type
}

func (b Block) GetFee() uint64 {
	return b.info.Fee
}
func (b Block) GetLinks() []Address {
	return append(b.Inputs, b.Outputs...)

}
func (b Block) Equals(a Block) bool {
	return b.info.HashLow == a.info.HashLow
}

func (b Block) ToString() string {
	return "Block info:[Hash:{" + hex.EncodeToString(b.info.HashLow[:]) +
		"}][Time:{" + strconv.FormatUint(b.info.Timestamp, 16) +
		"}]"
}
