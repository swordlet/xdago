//go:build pebble && !rocksdb

package pebbledb

import (
	"strconv"
	"sync"
	"xdago/config"
	"xdago/db"
)

type PebbleFactory struct {
	databases sync.Map
	config    *config.Config
}

func NewPebbleFactory(config *config.Config) PebbleFactory {
	return PebbleFactory{
		config: config,
	}
}

func (r *PebbleFactory) GetDB(name db.DatabaseName) *db.IKVSource {

	dataSource, _ := r.databases.LoadOrStore(name, func() *PebbleKv {
		var kv *PebbleKv
		if name == db.TIME {
			kv = NewPebbleKv(strconv.Itoa(int(name)), 10)
		} else {
			kv = NewPebbleKv(strconv.Itoa(int(name)), 0)
		}
		kv.SetConfig(r.config)
		return kv
	}())
	return dataSource.(*db.IKVSource)
}

func (r *PebbleFactory) Close() {
	r.databases.Range(func(key, value interface{}) bool {
		kv := value.(*PebbleKv)
		kv.Close()
		r.databases.Delete(key)
		return true
	})
}
