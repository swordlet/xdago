package store

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"xdago/common"
	"xdago/core"
	"xdago/db"
	"xdago/log"
	"xdago/utils"
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

func getFileNames(t uint64) []string {
	var names []string
	names = append(names, common.SUM_FILE_NAME)
	subDir := fmt.Sprintf("%02x", uint8(t>>40)) + "/"
	names = append(names, subDir+common.SUM_FILE_NAME)
	subDir = subDir + fmt.Sprintf("%02x", uint8(t>>32)) + "/"
	names = append(names, subDir+common.SUM_FILE_NAME)
	subDir = subDir + fmt.Sprintf("%02x", uint8(t>>24)) + "/"
	names = append(names, subDir+common.SUM_FILE_NAME)

	return names
}

func GetTimeKey(timestamp uint64, hashLow []byte) []byte {
	t := timestamp >> 16
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], t)
	if hashLow == nil {
		return utils.MergeBytes([]byte{common.TIME_HASH_INFO}, b[:])
	} else {
		if len(hashLow) != common.XDAG_HASH_SIZE {
			log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
		}
		return utils.MergeBytes([]byte{common.TIME_HASH_INFO}, b[:], hashLow)
	}
}

func getOurKey(index int32, hashLow []byte) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(index))
	return utils.MergeBytes([]byte{common.OURS_BLOCK_INFO}, b[:], hashLow)
}

func getHeight(height uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], height)
	return utils.MergeBytes([]byte{common.BLOCK_HEIGHT}, b[:])
}

func getOurIndex(key []byte) int32 {
	if len(key) >= 5 {
		b := key[1:5]
		return int32(binary.BigEndian.Uint32(b))
	}
	return 0
}

func getOurHash(key []byte) []byte {
	if len(key) >= 37 {
		return key[5:37]
	}
	return nil
}

func (bs *BlockStore) SaveXdagStatus(stats *core.XDAGStats) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*stats)
	if err != nil {
		log.Error("serialize stats error", log.Ctx{"err": err.Error()})
	}
	bs.indexSource.Put([]byte{common.SETTING_STATS}, buf.Bytes())
}

func (bs *BlockStore) GetXdagStatus() *core.XDAGStats {
	var s core.XDAGStats
	b := bs.indexSource.Get([]byte{common.SETTING_STATS})
	if b == nil {
		return nil
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(&s)
	if err != nil {
		log.Error("deserialize stats error", log.Ctx{"err": err.Error()})
		return nil
	}
	return &s
}

func (bs *BlockStore) SaveXdagtTopStatus(topStats *core.XDAGTopStatus) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*topStats)
	if err != nil {
		log.Error("serialize topStats error", log.Ctx{"err": err.Error()})
	}
	bs.indexSource.Put([]byte{common.SETTING_TOP_STATUS}, buf.Bytes())
}

func (bs *BlockStore) GetXdagTopStatus() *core.XDAGTopStatus {
	var s core.XDAGTopStatus
	b := bs.indexSource.Get([]byte{common.SETTING_TOP_STATUS})
	if b == nil {
		return nil
	}

	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(&s)
	if err != nil {
		log.Error("deserialize topStats error", log.Ctx{"err": err.Error()})
		return nil
	}
	return &s
}

func (bs *BlockStore) SaveBlock(block *core.Block) {
	t := block.GetTimestamp()
	h := block.GetHashLow()
	/// time中只拿key的后缀（hashLow）就够了，值可以不存
	bs.timeSource.Put(GetTimeKey(t, h[:]), []byte{0})
	b := block.GetXdagBlock().GetData()
	bs.blockSource.Put(h[:], b[:])
	bs.saveBlockSums(block)
	bs.SaveBlockInfo(block.Info())
}

func (bs *BlockStore) SaveOurBlock(index int32, hashLow []byte) {
	bs.indexSource.Put(getOurKey(index, hashLow), []byte{0})
}

func (bs *BlockStore) getOurBlock(index int32) []byte {
	var blockHashLow [32]byte
	bs.FetchOurBlocks(func(u int32, block *core.Block) bool {
		if u == index {
			if block != nil && block.GetHashLow() != common.EmptyHash {
				blockHashLow = block.GetHashLow()
				return true
			} else {
				return false
			}
		}
		return false
	})
	return blockHashLow[:]
}

func (bs *BlockStore) GetKeyIndexByHash(hashLow common.Hash) int32 {
	var keyIndex int32 = -1
	bs.FetchOurBlocks(func(i int32, block *core.Block) bool {
		if hashLow == block.GetHashLow() {
			atomic.StoreInt32(&keyIndex, i)
			return true
		}
		return false
	})
	return atomic.LoadInt32(&keyIndex)
}

func (bs *BlockStore) removeOurBlock(hashLow common.Hash) {
	bs.FetchOurBlocks(func(i int32, block *core.Block) bool {
		if hashLow == block.GetHashLow() {
			bs.indexSource.Delete(getOurKey(i, hashLow[:]))
			return true
		}
		return false
	})
}

func (bs *BlockStore) FetchOurBlocks(f func(int32, *core.Block) bool) {
	bs.indexSource.FetchPrefix([]byte{common.OURS_BLOCK_INFO}, func(k, v []byte) bool {
		index := getOurIndex(k)
		block := bs.GetBlockInfoByHash(getOurHash(k))
		return f(index, block)
	})
}

func (bs *BlockStore) saveBlockSums(block *core.Block) {
	var size uint64 = 512
	sum := block.GetXdagBlock().Sum
	t := block.GetTimestamp()
	fileNames := getFileNames(t)
	for i, file := range fileNames {
		bs.updateSum(file, sum, size, (t>>(40-8*i))&0xff)
	}
}

func (bs *BlockStore) getSums(key string) []byte {
	return bs.indexSource.Get(utils.MergeBytes([]byte{common.SUMS_BLOCK_INFO}, []byte(key)))
}

func (bs *BlockStore) putSums(key string, sums []byte) {
	bs.indexSource.Put(utils.MergeBytes([]byte{common.SUMS_BLOCK_INFO}, []byte(key)), sums)
}

func (bs *BlockStore) updateSum(key string, sum, size, index uint64) {
	sums := bs.getSums(key)
	if sums == nil {
		sums = make([]byte, 4096)
	} else {
		data := sums[16*int(index) : 16*int(index+1)]
		sum += binary.LittleEndian.Uint64(data[:8])
		size += binary.LittleEndian.Uint64(data[8:])
	}
	offset := int(16 * index)
	binary.LittleEndian.PutUint64(sums[offset:offset+8], sum)
	binary.LittleEndian.PutUint64(sums[offset+8:offset+16], size)

	bs.putSums(key, sums)
}

func (bs *BlockStore) LoadSum(startTime, endTime uint64) ([]byte, int) {
	endTime -= startTime
	if endTime == 0 || endTime&(endTime-1) != 0 {
		return nil, -1
	}

	var level int
	for level = -6; endTime != 0; level++ {
		endTime >>= 4
	}

	files := getFileNames(startTime & 0xffffff000000)
	var key string

	if level < 2 {
		key = files[3]
	} else if level < 4 {
		key = files[2]
	} else if level < 6 {
		key = files[1]
	} else {
		key = files[0]
	}

	buf := bs.getSums(key)
	sums := make([]byte, 256, 256)
	if buf == nil {
		return sums, 1
	}

	var sum, size uint64
	if level&1 != 0 {
		for i := 0; i < 256; i++ {
			totalSum := binary.LittleEndian.Uint64(buf[i*16 : i*16+8])
			sum += totalSum
			totalSize := binary.LittleEndian.Uint64(buf[i*16+8 : (i+1)*16])
			size += totalSize
			if i%16 == 0 && i != 0 {
				binary.LittleEndian.PutUint64(sums[i-16:i-8], sum)
				binary.LittleEndian.PutUint64(sums[i-8:i], size)
				sum = 0
				size = 0
			}
		}
	} else {
		index := int((startTime >> (level + 4) * 4) & 0xf0)
		copy(sums[:256], buf[index*16:index*16+256])
	}
	return sums, 1
}

func (bs *BlockStore) SaveBlockInfo(info *core.BlockInfo) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*info)
	if err != nil {
		log.Error("serialize stats error", log.Ctx{"err": err.Error()})
	}
	bs.indexSource.Put(utils.MergeBytes([]byte{common.HASH_BLOCK_INFO}, info.HashLow[:]), buf.Bytes())
	bs.indexSource.Put(getHeight(info.Height), info.HashLow[:])
}

func (bs *BlockStore) HasBlock(hashLow []byte) bool {
	return nil != bs.blockSource.Get(hashLow)
}

func (bs *BlockStore) HasBlockInfo(hashLow []byte) bool {
	return nil != bs.indexSource.Get(utils.MergeBytes([]byte{common.HASH_BLOCK_INFO}, hashLow))
}

func (bs *BlockStore) GetBlocksUsedTime(startTime, endTime uint64) []*core.Block {
	var res []*core.Block
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

func (bs *BlockStore) getBlocksByTime(startTime uint64) []*core.Block {
	var blocks []*core.Block
	keyPrefix := GetTimeKey(startTime, nil)
	keys := bs.timeSource.PrefixKeyLookup(keyPrefix)
	//fmt.Println(hex.EncodeToString(keyPrefix))
	for _, h := range keys {
		// 1 + 8 : prefix + time
		hash := h[9:41]
		block := bs.GetBlockByHash(hash, true)
		if block != nil {
			blocks = append(blocks, block)
		}
		//fmt.Println(hex.EncodeToString(bytes))
	}
	return blocks
}

//GetBlockByHeight 通过高度获取区块
func (bs *BlockStore) GetBlockByHeight(height uint64) *core.Block {
	hashLow := bs.indexSource.Get(getHeight(height))
	if hashLow == nil {
		return nil
	}
	return bs.GetBlockByHash(hashLow, false)
}

func (bs *BlockStore) GetBlockByHash(hashLow []byte, isRaw bool) *core.Block {
	if isRaw {
		return bs.GetRawBlockByHash(hashLow)
	}
	return bs.GetBlockInfoByHash(hashLow)
}

func (bs *BlockStore) GetRawBlockByHash(hashLow []byte) *core.Block {
	block := bs.GetBlockInfoByHash(hashLow)
	if block == nil {
		return nil
	}
	rawData := bs.blockSource.Get(hashLow)
	if rawData == nil {
		log.Error("No block origin data", log.Ctx{"hash": hex.EncodeToString(hashLow)})
		return nil
	}
	block.SetXdagBlock(core.NewXdagBlock(rawData))
	block.Parsed = false
	block.Parse()
	return block

}

func (bs *BlockStore) GetBlockInfoByHash(hashLow []byte) *core.Block {
	if len(hashLow) != common.XDAG_HASH_SIZE {
		log.Crit("hashlow size error", log.Ctx{"len": len(hashLow)})
	}
	if !bs.HasBlockInfo(hashLow) {
		return nil
	}
	var info core.BlockInfo
	value := bs.indexSource.Get(utils.MergeBytes([]byte{common.HASH_BLOCK_INFO}, hashLow))
	if value == nil {
		return nil
	}
	dec := gob.NewDecoder(bytes.NewBuffer(value))
	err := dec.Decode(&info)
	if err != nil {
		log.Error("can't deserialize block info data", log.Ctx{"info": hex.EncodeToString(value), "err": err.Error()})
		return nil
	}
	return core.NewBlockFromInfo(&info)
}

func (bs *BlockStore) IsSnapshotBoot() bool {
	data := bs.indexSource.Get([]byte{common.SNAPSHOT_BOOT})
	if data == nil {
		return false
	} else {
		res := binary.BigEndian.Uint32(data)
		return res == 1
	}
}

func (bs *BlockStore) SetSnapshotBoot() {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], 1)
	bs.indexSource.Put([]byte{common.SNAPSHOT_BOOT}, b[:])
}

func (bs *BlockStore) SavePreSeed(preSeed []byte) {
	bs.indexSource.Put([]byte{common.SNAPSHOT_PRESEED}, preSeed)
}

func (bs *BlockStore) GetPreSeed() []byte {
	return bs.indexSource.Get([]byte{common.SNAPSHOT_PRESEED})
}
