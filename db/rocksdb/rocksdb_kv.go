//go:build rocksdb && !pebble

package rocksdb

import (
	"errors"
	"github.com/linxGnu/grocksdb"
	"sync"
	"xdago/db"
	"xdago/log"
)

type RocksKv struct {
	sync.RWMutex
	name             string
	alive            bool
	db               *grocksdb.DB
	readOpt          *grocksdb.ReadOptions
	prefixSeekLength int
}

func NewRocksKv(name string, preSeekLen int) *RocksKv {
	log.Debug("New db", log.Ctx{"dbname": name})
	return &RocksKv{
		name:             name,
		prefixSeekLength: preSeekLen,
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
	p.readOpt.Destroy()
	p.alive = false
}

func (p *RocksKv) Init() error {
	p.Lock()
	defer p.Unlock()
	log.Debug("Init db", log.Ctx{"dbname": p.name})
	if p.alive {
		return nil
	}
	if p.name == "" {
		return errors.New("no name set to db")
	}
	options := grocksdb.NewDefaultOptions()
	options.SetCreateIfMissingColumnFamilies(true)
	options.SetCompression(grocksdb.LZ4Compression)
	options.SetBottommostCompression(grocksdb.LZ4HCCompression)
	options.SetLevelCompactionDynamicLevelBytes(true)
	options.SetMaxOpenFiles(10)    // TODO: from config
	options.IncreaseParallelism(4) // TODO: from config
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

	return nil
}

func (p *RocksKv) Reset() {

}

func (p *RocksKv) IsAlive() bool {
	return p.alive
}

func (p *RocksKv) Put(key, val []byte) error {
	return nil
}

func (p *RocksKv) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (p *RocksKv) Delete(key []byte) error {
	return nil
}

func (p *RocksKv) Keys() ([][]byte, error) {
	return nil, nil
}

func (p *RocksKv) PrefixKeyLookup(key []byte) ([][]byte, error) {
	return nil, nil
}

func (p *RocksKv) FetchPrefix(key []byte, f db.FetchFunc) error {
	return nil
}

func (p *RocksKv) PrefixValueLookup(key []byte) ([][]byte, error) {
	return nil, nil
}
