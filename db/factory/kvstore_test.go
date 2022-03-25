//go:build pebble && !rocksdb

////go:build rocksdb && !pebble
//switch for KV store

package factory

import (
	"encoding/hex"
	"github.com/magiconair/properties/assert"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/config"
	"xdago/crypto"
	"xdago/db"
	"xdago/db/store"
	"xdago/log"
)

// CGO_LDFLAGS=-lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd
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

func TestKv(t *testing.T) {

	kvFactory := NewKvStoreFactory(c)

	blockSource := kvFactory.GetDB(db.BLOCK)
	indexSource := kvFactory.GetDB(db.INDEX)
	orphanSource := kvFactory.GetDB(db.ORPHANIND)

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
	kvFactory := NewKvStoreFactory(c)
	timeSource := kvFactory.GetDB(db.TIME)
	timeSource.Reset()

	hashlow1 := crypto.HashTwice([]byte("1"))
	hashlow2 := crypto.HashTwice([]byte("2"))

	var time1 uint64 = 1602226304712
	value1, _ := hex.DecodeString("1234")
	value2, _ := hex.DecodeString("2345")

	key1 := store.GetTimeKey(time1, hashlow1[:])
	key2 := store.GetTimeKey(time1, hashlow2[:])

	timeSource.Put(key1, value1)
	timeSource.Put(key2, value2)

	all := timeSource.Keys()
	assert.Equal(t, len(all), 2)

	var search uint64 = 1602226304712
	key := store.GetTimeKey(search, nil)
	keys := timeSource.PrefixKeyLookup(key)
	assert.Equal(t, len(keys), 2)
	values := timeSource.PrefixValueLookup(key)
	assert.Equal(t, len(values), 2)

	timeSource.Close()

}
