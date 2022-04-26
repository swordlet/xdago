package store

import (
	"encoding/binary"
	"encoding/hex"
	"xdago/common"
	"xdago/core"
	"xdago/db"
	"xdago/log"
	"xdago/utils"
)

const (
	ORPHAN_PREFEX uint8 = 0x00
)

var ORPHAN_SIZE, _ = hex.DecodeString("FFFFFFFFFFFFFFFF")

type OrphanPool struct {
	orphanSource db.IKVSource
}

func NewOrphanPool(orphan db.IKVSource) *OrphanPool {
	return &OrphanPool{
		orphanSource: orphan,
	}
}

func (p *OrphanPool) Init() {
	p.orphanSource.Init()
	if p.orphanSource.Get(ORPHAN_SIZE) == nil {
		p.orphanSource.Put(ORPHAN_SIZE, utils.U64ToBytes(0, binary.BigEndian))
	}
}

func (p *OrphanPool) Reset() {
	p.orphanSource.Reset()
	p.orphanSource.Put(ORPHAN_SIZE, utils.U64ToBytes(0, binary.BigEndian))
}

func (p *OrphanPool) GetOrphan(num, sendTime uint64) []core.Address {
	var res []core.Address
	if p.orphanSource.Get(ORPHAN_SIZE) == nil || p.getOrphanSize() == 0 {
		return nil
	}

	orphanSize := p.getOrphanSize()
	addNum := utils.MinUint64(orphanSize, num)
	key := []byte{ORPHAN_PREFEX}
	ans := p.orphanSource.PrefixValueLookup(key)

	for _, an := range ans {
		if addNum == 0 {
			break
		}
		// TODO:判断时间，这里出现过orphanSource获取key时为空的情况
		if p.orphanSource.Get(an) == nil {
			continue
		}
		timestamp := binary.LittleEndian.Uint64(p.orphanSource.Get(an))
		if timestamp < sendTime {
			addNum--
			var field common.Field
			copy(field[:], an[1:33])
			res = append(res, core.AddressFromType(field, common.XDAG_FIELD_OUT))
		}
	}
	return res
}

func (p *OrphanPool) DeleteByHash(hashLow []byte) {
	log.Debug("orphan delete by hash", log.Ctx{"hash": hex.EncodeToString(hashLow)})
	p.orphanSource.Delete(utils.MergeBytes([]byte{ORPHAN_PREFEX}, hashLow))
	curSize := binary.BigEndian.Uint64(p.orphanSource.Get(ORPHAN_SIZE))
	p.orphanSource.Put(ORPHAN_SIZE, utils.U64ToBytes(curSize-1, binary.BigEndian))
}

func (p *OrphanPool) AddOrphan(block *core.Block) {
	hash := block.GetHashLow()
	p.orphanSource.Put(utils.MergeBytes([]byte{ORPHAN_PREFEX}, hash[:]),
		utils.U64ToBytes(block.GetTimestamp(), binary.BigEndian))

	curSize := binary.BigEndian.Uint64(p.orphanSource.Get(ORPHAN_SIZE))
	p.orphanSource.Put(ORPHAN_SIZE, utils.U64ToBytes(curSize+1, binary.BigEndian))
}

func (p *OrphanPool) getOrphanSize() uint64 {
	curSize := binary.BigEndian.Uint64(p.orphanSource.Get(ORPHAN_SIZE))
	log.Debug("current orphan size",
		log.Ctx{"size": curSize, "hex": hex.EncodeToString(p.orphanSource.Get(ORPHAN_SIZE))})
	return curSize
}

func (p *OrphanPool) ContainsKey(hashLow []byte) bool {
	return p.orphanSource.Get(utils.MergeBytes([]byte{ORPHAN_PREFEX}, hashLow)) != nil
}
