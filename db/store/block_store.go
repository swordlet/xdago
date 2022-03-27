package store

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"xdago/core"
	"xdago/db"
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
	fmt.Println(hex.EncodeToString(keyPrefix))
	for _, bytes := range keys {
		//	// 1 + 8 : prefix + time
		//	hash := bytes[9:41]
		//	block := getBlockByHash(Bytes32.wrap(hash), true)
		//	if block != null {
		//		blocks.add(block)
		//	}
		fmt.Println(hex.EncodeToString(bytes))
	}
	return blocks
}
