//go:build pebble && !rocksdb

package pebbledb

import (
	"strconv"
	"xdago/config"
	"xdago/db"
)

type KvStoreFactory struct {
	databases map[db.DatabaseName]interface{}
	config    *config.Config
}

func NewKvStoreFactory(config *config.Config) KvStoreFactory {
	return KvStoreFactory{
		config:    config,
		databases: make(map[db.DatabaseName]interface{}),
	}
}

func (r *KvStoreFactory) GetDB(name db.DatabaseName) *db.IKVSource {

	dataSource, ok := r.databases[name]
	if !ok {
		var kv interface{}
		if name == db.TIME {
			kv = NewPebbleKv(strconv.Itoa(int(name)), 10)
		} else {
			kv = NewPebbleKv(strconv.Itoa(int(name)), 0)
		}
		kv.(*PebbleKv).SetConfig(r.config)
		r.databases[name] = kv
		return kv.(*db.IKVSource)
	}

	return dataSource.(*db.IKVSource)
}

func (r *KvStoreFactory) Close() {
	for key, value := range r.databases {
		kv := value.(*PebbleKv)
		kv.Close()
		delete(r.databases, key)
	}
}
