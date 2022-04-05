//go:build rocksdb && !pebble

package factory

import (
	"strconv"
	"xdago/common"
	"xdago/config"
	"xdago/db"
	"xdago/db/rocksdb"
)

type KvStoreFactory struct {
	databases map[common.DatabaseName]interface{}
	config    *config.Config
}

func NewKvStoreFactory(config *config.Config) KvStoreFactory {
	return KvStoreFactory{
		config:    config,
		databases: make(map[common.DatabaseName]interface{}),
	}
}

func (r *KvStoreFactory) GetDB(name common.DatabaseName) db.IKVSource {

	dataSource, ok := r.databases[name]
	if !ok {
		var kv interface{}
		if name == common.DB_TIME {
			kv = rocksdb.NewRocksKv(strconv.Itoa(int(name)), 9)
		} else {
			kv = rocksdb.NewRocksKv(strconv.Itoa(int(name)), 0)
		}
		kv.(*rocksdb.RocksKv).SetConfig(r.config)
		r.databases[name] = kv
		return kv.(db.IKVSource)
	}

	return dataSource.(db.IKVSource)
}

func (r *KvStoreFactory) Close() {
	for key, value := range r.databases {
		kv := value.(*rocksdb.RocksKv)
		kv.Close()
		delete(r.databases, key)
	}
}
