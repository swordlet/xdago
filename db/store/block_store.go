package store

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"xdago/core"
	"xdago/db"
	"xdago/log"
	"xdago/utils"
)

const (
	SETTING_STATS      byte = 0x10
	TIME_HASH_INFO     byte = 0x20
	HASH_BLOCK_INFO    byte = 0x30
	SUMS_BLOCK_INFO    byte = 0x40
	OURS_BLOCK_INFO    byte = 0x50
	SETTING_TOP_STATUS byte = 0x60
	SNAPSHOT_BOOT      byte = 0x70
	BLOCK_HEIGHT       byte = 0x80
	SNAPSHOT_PRESEED   byte = 0x90
	SUM_FILE_NAME           = "sums.dat"
)

type BlockStore struct {
	timeSource  db.IKVSource
	indexSource db.IKVSource
	blockSource db.IKVSource
}

func NewBlockStore(indexSource, timeSource, blockSource db.IKVSource) *BlockStore {
	return &BlockStore{
		indexSource: indexSource,
		timeSource:  timeSource,
		blockSource: blockSource,
	}

}

func (bs *BlockStore) Init() {
	bs.indexSource.Init()
	bs.timeSource.Init()
	bs.blockSource.Init()
}

func (bs *BlockStore) Reset() {
	bs.indexSource.Reset()
	bs.timeSource.Reset()
	bs.blockSource.Reset()
}

func (bs *BlockStore) Close() {
	bs.indexSource.Close()
	bs.timeSource.Close()
	bs.blockSource.Close()
}

func GetTimeKey(timestamp uint64, hashLow []byte) []byte {
	t := timestamp >> 16
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], t)
	if hashLow == nil {
		return utils.MergeBytes([]byte{TIME_HASH_INFO}, b[:])
	} else {
		if len(hashLow) != core.XDAG_HASH_SIZE {
			log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
		}
		return utils.MergeBytes([]byte{TIME_HASH_INFO}, b[:], hashLow)
	}
}

func (bs *BlockStore) GetBlocksUsedTime(startTime, endTime uint64) []core.Block {
	var res []core.Block
	time := startTime
	for time < endTime {
		blocks := bs.getBlocksByTime(time)
		time += 0x10000
		if blocks == nil {
			continue
		}
		res = append(res, blocks...)
	}
	return res
}

func (bs *BlockStore) getBlocksByTime(startTime uint64) []core.Block {
	var blocks []core.Block
	keyPrefix := GetTimeKey(startTime, nil)
	keys := bs.timeSource.PrefixKeyLookup(keyPrefix)
	//fmt.Println(hex.EncodeToString(keyPrefix))
	for _, h := range keys {
		// 1 + 8 : prefix + time
		hash := h[9:41]
		block := bs.GetBlockByHash(hash, true)
		if !block.IsEmpty() {
			blocks = append(blocks, block)
		}
		//fmt.Println(hex.EncodeToString(bytes))
	}
	return blocks
}

func (bs *BlockStore) GetBlockByHash(hashLow []byte, isRaw bool) core.Block {
	if isRaw {
		return bs.GetRawBlockByHash(hashLow)
	}
	return bs.GetBlockInfoByHash(hashLow)
}

func (bs *BlockStore) GetRawBlockByHash(hashLow []byte) core.Block {
	block := bs.GetBlockInfoByHash(hashLow)
	if block.IsEmpty() {
		return block
	}
	rawData := bs.blockSource.Get(hashLow)
	if rawData == nil {
		log.Error("No block origin data", log.Ctx{"hash": hex.EncodeToString(hashLow)})
		return core.Block{}
	}
	block.SetXdagBlock(core.NewXdagBlock(rawData))
	block.Parsed = false
	block.Parse()
	return block

}

func (bs *BlockStore) GetBlockInfoByHash(hashLow []byte) core.Block {
	if len(hashLow) != core.XDAG_HASH_SIZE {
		log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
	}
	if !bs.HasBlockInfo(hashLow) {
		return core.Block{}
	}
	var info core.BlockInfo
	value := bs.indexSource.Get(utils.MergeBytes([]byte{HASH_BLOCK_INFO}, hashLow))
	if value == nil {
		return core.Block{}
	}
	err := binary.Read(bytes.NewBuffer(value), binary.LittleEndian, &info)
	if err != nil {
		log.Error("can't deserialize block info data", log.Ctx{"info": hex.EncodeToString(value), "err": err.Error()})
		return core.Block{}
	}
	return core.NewBlockFromInfo(&info)
}

func (bs *BlockStore) HasBlock(hashLow []byte) bool {
	return nil != bs.blockSource.Get(hashLow)
}

func (bs *BlockStore) HasBlockInfo(hashLow []byte) bool {
	return nil != bs.indexSource.Get(utils.MergeBytes([]byte{HASH_BLOCK_INFO}, hashLow))
}
