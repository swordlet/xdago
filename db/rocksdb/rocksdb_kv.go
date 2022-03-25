//go:build rocksdb && !pebble

package rocksdb

import (
	"encoding/hex"
	"github.com/linxGnu/grocksdb"
	"os"
	"path"
	"sync"
	"xdago/config"
	"xdago/db"
	"xdago/log"
	"xdago/utils"
)

type RocksKv struct {
	sync.RWMutex
	config           *config.Config
	name             string
	alive            bool
	db               *grocksdb.DB
	readOpt          *grocksdb.ReadOptions
	writeOpt         *grocksdb.WriteOptions
	prefixSeekLength int
}

func (p *RocksKv) SetConfig(config *config.Config) {
	p.config = config
}

func NewRocksKv(name string, preSeekLen int) *RocksKv {
	log.Debug("New db", log.Ctx{"dbname": name})
	return &RocksKv{
		name:             name,
		prefixSeekLength: preSeekLen,
		writeOpt:         grocksdb.NewDefaultWriteOptions(),
	}
}
func (p *RocksKv) GetName() string {
	return p.name
}

func (p *RocksKv) SetName(name string) {
	p.name = name
}

func (p *RocksKv) Close() {
	p.Lock()
	defer p.Unlock()
	if !p.alive {
		return
	}
	log.Debug("Close db", log.Ctx{"dbname": p.name})
	p.db.Close()
	p.writeOpt.Destroy()
	p.readOpt.Destroy()
	p.alive = false
}

func (p *RocksKv) Init() {
	p.Lock()
	defer p.Unlock()

	log.Debug("~> RocksdbKVSource.init()", log.Ctx{"dbname": p.name})
	if p.alive {
		return
	}
	if p.name == "" {
		log.Crit("no name set to db")
	}
	options := grocksdb.NewDefaultOptions()
	options.SetCreateIfMissingColumnFamilies(true)
	options.SetCompression(grocksdb.LZ4Compression)
	options.SetBottommostCompression(grocksdb.LZ4HCCompression)
	options.SetLevelCompactionDynamicLevelBytes(true)
	options.SetMaxOpenFiles(p.config.StoreMaxOpenFiles())
	options.IncreaseParallelism(p.config.StoreMaxThreads())
	options.SetPrefixExtractor(grocksdb.NewFixedPrefixTransform(p.prefixSeekLength))

	tableConfig := grocksdb.NewDefaultBlockBasedTableOptions()
	tableConfig.SetBlockSize(16 * 1024)
	tableConfig.SetBlockCache(grocksdb.NewLRUCache(32 * 1024 * 1024))
	tableConfig.SetCacheIndexAndFilterBlocks(true)
	tableConfig.SetPinL0FilterAndIndexBlocksInCache(true)
	tableConfig.SetFilterPolicy(grocksdb.NewBloomFilter(10))
	options.SetBlockBasedTableFactory(tableConfig)

	p.readOpt = grocksdb.NewDefaultReadOptions()
	p.readOpt.SetPrefixSameAsStart(true)
	p.readOpt.SetVerifyChecksums(false)

	log.Info("Opening db", log.Ctx{"dbname": p.name})
	dbPath := p.getPath()
	dir := path.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0666); err != nil {
			panic(err)
		}
	}
	if p.config.StoreFromBackup() {
		if stat, err := os.Stat(p.backupPath()); err != nil || stat.Mode().Perm() != 0666 {
			log.Crit("db backup file permission error", log.Ctx{"dbname": p.name})
		}
		// TODO: open backup db
	}
	log.Debug("Open existing db or create new db ", log.Ctx{"dbname": p.name})
	var err error
	p.db, err = grocksdb.OpenDb(options, dbPath)
	if err != nil {
		log.Crit("Failed to open db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}
	p.alive = true

	log.Debug("<~ RocksdbKVSource.init()", log.Ctx{"dbname": p.name})

}

func (p *RocksKv) backup() {
	// TODO: back up db to backup path
}

func (p *RocksKv) Reset() {
	p.Close()
	err := os.RemoveAll(p.getPath())
	if err != nil {
		log.Crit("Failed to reset db", log.Ctx{"remove db": err.Error()})
	}
	p.Init()
}

func (p *RocksKv) IsAlive() bool {
	return p.alive
}

func (p *RocksKv) Put(key, val []byte) {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> RocksdbKVSource.put():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(val)})
	var err error
	if val != nil {
		if p.db == nil {

		} else {
			err = p.db.Put(p.writeOpt, key, val)
		}
	} else {
		err = p.db.Delete(p.writeOpt, key)
	}
	if err != nil {
		log.Crit("Failed to put into db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}
	log.Trace("<~ RocksdbKVSource.put():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(val)})

}

func (p *RocksKv) Get(key []byte) []byte {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> RocksdbKVSource.get():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})

	val, err := p.db.GetBytes(p.readOpt, key)
	if err != nil {
		log.Crit("Failed to get from db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}

	log.Trace("<~ RocksdbKVSource.get():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(val)})

	return val
}

func (p *RocksKv) Delete(key []byte) {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> RocksdbKVSource.delete():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})
	err := p.db.Delete(p.writeOpt, key)
	if err != nil {
		log.Crit("Failed to delete from db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}

	log.Trace("<~ RocksdbKVSource.delete():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})
}

func (p *RocksKv) Keys() [][]byte {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> RocksdbKVSource.keys():", log.Ctx{"dbname": p.name})
	var keys [][]byte
	iter := p.db.NewIterator(p.readOpt)
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key().Data())
	}

	log.Trace("<~ RocksdbKVSource.keys():", log.Ctx{"dbname": p.name})
	return keys
}

func (p *RocksKv) FetchPrefix(key []byte, f db.FetchFunc) {
	p.RLock()
	defer p.RUnlock()

	iter := p.db.NewIterator(p.readOpt)
	for iter.Seek(key); iter.Valid(); iter.Next() {
		if utils.KeyStartWith(iter.Key().Data(), key) {
			if f(iter.Key().Data(), iter.Value().Data()) {
				return
			}
		} else {
			return
		}
	}
}

func (p *RocksKv) PrefixKeyLookup(key []byte) [][]byte {
	var keyList [][]byte
	p.FetchPrefix(key, func(k, v []byte) bool {
		keyList = append(keyList, k)
		return false
	})
	return keyList
}

func (p *RocksKv) PrefixValueLookup(key []byte) [][]byte {
	var valueList [][]byte
	p.FetchPrefix(key, func(k, v []byte) bool {
		valueList = append(valueList, v)
		return false
	})
	return valueList
}

func (p *RocksKv) getPath() string {
	return path.Join(p.config.StoreDir(), p.name)
}

func (p *RocksKv) backupPath() string {
	return path.Join(p.config.StoreDir(), "backup", p.name)
}
