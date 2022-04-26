//go:build pebble && !rocksdb

////go:build rocksdb && !pebble

//conditional build switch for KV store

package store

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
	"xdago/common"
	"xdago/config"
	"xdago/core"
	"xdago/crypto"
	"xdago/db/factory"
	"xdago/secp256k1"
)

func testBsInit() (*config.Config, *BlockStore) {
	cfg := testInit()
	kvFactory := factory.NewKvStoreFactory(cfg)

	indexSource := kvFactory.GetDB(common.DB_INDEX)
	timeSource := kvFactory.GetDB(common.DB_TIME)
	blockSource := kvFactory.GetDB(common.DB_BLOCK)

	return cfg, NewBlockStore(indexSource, timeSource, blockSource)
}

func TestGetBlocksUsedTime(t *testing.T) {
	cfg := testInit()
	kvFactory := factory.NewKvStoreFactory(cfg)

	indexSource := kvFactory.GetDB(common.DB_INDEX)
	timeSource := kvFactory.GetDB(common.DB_TIME)
	blockSource := kvFactory.GetDB(common.DB_BLOCK)

	bs := NewBlockStore(indexSource, timeSource, blockSource)

	bs.Reset()

	hashlow1 := crypto.HashTwice([]byte("1"))
	hashlow2 := crypto.HashTwice([]byte("2"))
	hashlow3 := crypto.HashTwice([]byte("3"))
	hashlow4 := crypto.HashTwice([]byte("4"))
	hashlow5 := crypto.HashTwice([]byte("5"))

	var time1 uint64 = 1602226304712
	var time2 uint64 = 1602226304712 + 0x10000
	var time3 uint64 = 1602226304712 + 0x1000000
	var time4 uint64 = 1602226304712 + 0x100000000

	value1, _ := hex.DecodeString("1234")
	value2, _ := hex.DecodeString("2345")
	value3, _ := hex.DecodeString("3456")
	value4, _ := hex.DecodeString("4567")
	value5, _ := hex.DecodeString("5678")

	key1 := GetTimeKey(time1, hashlow1[:])
	key2 := GetTimeKey(time2, hashlow2[:])
	key3 := GetTimeKey(time3, hashlow3[:])
	key4 := GetTimeKey(time4, hashlow4[:])

	key5 := GetTimeKey(time1, hashlow5[:]) // same prefix with key1

	timeSource.Put(key1, value1)
	timeSource.Put(key2, value2)
	timeSource.Put(key3, value3)
	timeSource.Put(key4, value4)

	timeSource.Put(key5, value5)

	bs.GetBlocksUsedTime(1602226304712, 1602226304712+0x1000000)
	bs.Close()

}

func TestSaveXdagStatus(t *testing.T) {
	_, bs := testBsInit()
	bs.Init()
	stats := core.NewEmptyXDAGStats()
	stats.NMain = 1
	bs.SaveXdagStatus(stats)
	storedStats := bs.GetXdagStatus()
	assert.Equal(t, storedStats.NMain, stats.NMain)
}

func TestSaveBlock(t *testing.T) {
	cfg, bs := testBsInit()
	bs.Init()
	timestamp := uint64(time.Now().UnixMilli())
	privKey, _ := secp256k1.GeneratePrivateKey()

	b := core.GenerateAddressBlock(cfg, privKey, timestamp)
	bs.SaveBlock(b)
	h := b.GetHashLow()
	storedBlock := bs.GetBlockByHash(h[:], true)

	assert.Equal(t, storedBlock.GetXdagBlock().GetData(), b.GetXdagBlock().GetData())
	assert.Equal(t, bytes.Compare(storedBlock.ToBytes(), b.ToBytes()), 0)
}

func TestSaveOurBlock(t *testing.T) {
	cfg, bs := testBsInit()
	bs.Reset()
	bs.Init()
	timestamp := uint64(time.Now().UnixMilli())
	privKey, _ := secp256k1.GeneratePrivateKey()

	b := core.GenerateAddressBlock(cfg, privKey, timestamp)
	bs.SaveBlock(b)
	h := b.GetHashLow()
	bs.SaveOurBlock(1, h[:])
	assert.Equal(t, bytes.Compare(bs.getOurBlock(1), h[:]), 0)

}

func TestRemoveOurBlock(t *testing.T) {
	cfg, bs := testBsInit()
	bs.Reset()
	bs.Init()
	timestamp := uint64(time.Now().UnixMilli())
	privKey, _ := secp256k1.GeneratePrivateKey()

	b := core.GenerateAddressBlock(cfg, privKey, timestamp)
	bs.SaveBlock(b)
	h := b.GetHashLow()
	bs.SaveOurBlock(1, h[:])

	if bytes.Compare(bs.getOurBlock(1), common.EmptyHash[:]) == 0 {
		panic("error")
	}
	bs.removeOurBlock(h)

	assert.Equal(t, bytes.Compare(bs.getOurBlock(1), common.EmptyHash[:]), 0)
}

func TestSaveBlockSums(t *testing.T) {
	cfg, bs := testBsInit()
	bs.Reset()
	bs.Init()
	var timestamp uint64 = 1602951025307
	privKey, _ := secp256k1.GeneratePrivateKey()

	b := core.GenerateAddressBlock(cfg, privKey, timestamp)
	bs.SaveBlock(b)
	sums, res := bs.LoadSum(timestamp, timestamp+64*1024)
	assert.Equal(t, res, 1)
	fmt.Println(sums)
}

func TestGetBlockByTime(t *testing.T) {
	_, bs := testBsInit()
	bs.Reset()
	bs.Init()
	b, _ := hex.DecodeString("00000000000000003833333333530540ffff8741810100000000000000000000032dea64ace570d7ae8668c8a4f52265c16497c9dd8cd62b0000000000000000f1f245ea01d304c3be265cad77f5589acdc45a7b3d35972f0000000000000000f23cddd22c17bf0a083e4bbe63c0e224dfc20a583238ef7a0000000000000000b4407441ad9c0372a7f053a3dbaaa4855589228cef7f05b000000000000000004206427aa89b7066b05379bec0e9264a34c55391f12137bb00000000000000009b55f3a7af41e29d8b6b4e4581387c507726437f7aacc7930000000000000000905786241884e7520a8ad2c777871b28548c78b8964107e20000000000000000a2583dc5f6001020e406edb1c6ed52c41bae2ef1dda9439200000000000000009f5c7e9633614d665fe6739fd122cdb0360b2c688d02685d00000000000000005fbc1107fe34e3faeab63e1ef3e24b6c66053103c4868a6600000000000000003a7883fa0ddb348428d72856ff0527e5aff79b2c739fb946b53ce6b29530a07dc821749a7ffa3f6b6e3417d6c0c54457c9909800b7dc5b034b7a1f979032e4cb000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008ed85467b39cc220720472c5f0b116afaccce977c71a655daae7789782c5fae9")
	b1, _ := hex.DecodeString("00000000000000003833333333530540ffff8941810100000000000000000000032dea64ace570d7ae8668c8a4f52265c16497c9dd8cd62b0000000000000000f1f245ea01d304c3be265cad77f5589acdc45a7b3d35972f0000000000000000f23cddd22c17bf0a083e4bbe63c0e224dfc20a583238ef7a0000000000000000b4407441ad9c0372a7f053a3dbaaa4855589228cef7f05b000000000000000004206427aa89b7066b05379bec0e9264a34c55391f12137bb00000000000000009b55f3a7af41e29d8b6b4e4581387c507726437f7aacc7930000000000000000905786241884e7520a8ad2c777871b28548c78b8964107e20000000000000000a2583dc5f6001020e406edb1c6ed52c41bae2ef1dda9439200000000000000009f5c7e9633614d665fe6739fd122cdb0360b2c688d02685d00000000000000005fbc1107fe34e3faeab63e1ef3e24b6c66053103c4868a6600000000000000003a7883fa0ddb348428d72856ff0527e5aff79b2c739fb946b53ce6b29530a07dc821749a7ffa3f6b6e3417d6c0c54457c9909800b7dc5b034b7a1f979032e4cb000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008ed85467b39cc220720472c5f0b116afaccce977c71a655daae7789782c5fae9")
	block := core.NewBlockFromXdag(core.NewXdagBlock(b))
	block1 := core.NewBlockFromXdag(core.NewXdagBlock(b1))

	timestamp := block.GetTimestamp()

	bs.SaveBlock(block)
	bs.SaveBlock(block1)

	blocks := bs.getBlocksByTime(timestamp)

	assert.Equal(t, len(blocks), 1)
	assert.Equal(t, blocks[0].GetHashLow(), block.GetHashLow())
}
