//go:build rocksdb && !pebble

package rocksdb

import (
	"strconv"
	"sync"
	"xdago/config"
	"xdago/db"
)

type RocksdbFactory struct {
	databases sync.Map
	config    *config.Config
}

func NewRocksdbFactory(config *config.Config) RocksdbFactory {
	return RocksdbFactory{
		config: config,
	}
}

func (r *RocksdbFactory) GetDB(name db.DatabaseName) *db.IKVSource {

	dataSource, _ := r.databases.LoadOrStore(name, func() *RocksKv {
		var kv *RocksKv
		if name == db.TIME {
			kv = NewRocksKv(strconv.Itoa(int(name)), 10)
		} else {
			kv = NewRocksKv(strconv.Itoa(int(name)), 0)
		}
		kv.SetConfig(r.config)
		return kv
	}())
	return dataSource.(*db.IKVSource)
}

func (r *RocksdbFactory) Close() {
	r.databases.Range(func(key, value interface{}) bool {
		kv := value.(*RocksKv)
		kv.Close()
		r.databases.Delete(key)
		return true
	})
}
