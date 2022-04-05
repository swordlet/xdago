//go:build pebble && !rocksdb

////go:build rocksdb && !pebble

//conditional build switch for KV store

package store

import (
	"encoding/hex"
	"github.com/magiconair/properties/assert"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/common"
	"xdago/config"
	"xdago/crypto"
	"xdago/db/factory"
	"xdago/log"
)

// CGO_LDFLAGS=-lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd

func testInit() *config.Config {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	cfg := config.DevNetConfig()
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
	return cfg
}

func TestKv(t *testing.T) {
	cfg := testInit()
	kvFactory := factory.NewKvStoreFactory(cfg)

	blockSource := kvFactory.GetDB(common.DB_BLOCK)
	indexSource := kvFactory.GetDB(common.DB_INDEX)
	orphanSource := kvFactory.GetDB(common.DB_ORPHANIND)

	blockSource.Reset()
	indexSource.Reset()
	orphanSource.Reset()

	key, _ := hex.DecodeString("EEEE")
	value, _ := hex.DecodeString("1234")

	blockSource.Put(key, value)
	indexSource.Put(key, value)
	orphanSource.Put(key, value)

	assert.Equal(t, "1234", hex.EncodeToString(blockSource.Get(key)))
	assert.Equal(t, "1234", hex.EncodeToString(indexSource.Get(key)))
	assert.Equal(t, "1234", hex.EncodeToString(orphanSource.Get(key)))

	blockSource.Close()
	indexSource.Close()
	orphanSource.Close()
}

func TestPrefix(t *testing.T) {
	cfg := testInit()
	kvFactory := factory.NewKvStoreFactory(cfg)
	timeSource := kvFactory.GetDB(common.DB_TIME)
	timeSource.Reset()

	hashlow1 := crypto.HashTwice([]byte("1"))
	hashlow2 := crypto.HashTwice([]byte("2"))

	var time1 uint64 = 1602226304712
	value1, _ := hex.DecodeString("1234")
	value2, _ := hex.DecodeString("2345")

	key1 := GetTimeKey(time1, hashlow1[:])
	key2 := GetTimeKey(time1, hashlow2[:])

	timeSource.Put(key1, value1)
	timeSource.Put(key2, value2)

	all := timeSource.Keys()
	assert.Equal(t, len(all), 2)

	var search uint64 = 1602226304712
	key := GetTimeKey(search, nil)
	keys := timeSource.PrefixKeyLookup(key)
	assert.Equal(t, len(keys), 2)
	values := timeSource.PrefixValueLookup(key)
	assert.Equal(t, len(values), 2)

	timeSource.Close()
}
