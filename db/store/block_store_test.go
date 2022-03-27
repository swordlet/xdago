//go:build pebble && !rocksdb

////go:build rocksdb && !pebble

//conditional build switch for KV store

package store

import (
	"encoding/hex"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/config"
	"xdago/crypto"
	"xdago/db"
	"xdago/db/factory"
	"xdago/log"
)

var c *config.Config

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	c = config.DevNetConfig()
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
}

func TestGetBlocksUsedTime(t *testing.T) {
	kvFactory := factory.NewKvStoreFactory(c)

	indexSource := kvFactory.GetDB(db.INDEX)
	timeSource := kvFactory.GetDB(db.TIME)
	blockSource := kvFactory.GetDB(db.BLOCK)

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
