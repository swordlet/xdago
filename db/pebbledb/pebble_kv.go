//go:build pebble && !rocksdb

package pebbledb

import (
	"encoding/hex"
	"errors"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
	"os"
	"path"
	"sync"
	"xdago/config"
	"xdago/db"
	"xdago/log"
	"xdago/utils"
)

type PebbleKv struct {
	sync.RWMutex
	config           *config.Config
	name             string
	alive            bool
	db               *pebble.DB
	writeOpt         *pebble.WriteOptions
	prefixSeekLength int
}

func (p *PebbleKv) SetConfig(config *config.Config) {
	p.config = config
}

func NewPebbleKv(name string, preSeekLen int) *PebbleKv {
	log.Debug("New db", log.Ctx{"dbname": name})
	return &PebbleKv{
		name:             name,
		prefixSeekLength: preSeekLen,
		writeOpt: &pebble.WriteOptions{
			Sync: true,
		},
	}
}
func (p *PebbleKv) GetName() string {
	return p.name
}

func (p *PebbleKv) SetName(name string) {
	p.name = name
}

func (p *PebbleKv) Close() {
	p.Lock()
	defer p.Unlock()
	if !p.alive {
		return
	}
	log.Debug("Close db", log.Ctx{"dbname": p.name})
	err := p.db.Close()
	if err != nil {
		log.Crit("Failed close db", log.Ctx{"dbname": p.name, "err": err})
	}
	p.alive = false
}

func (p *PebbleKv) Init() {
	p.Lock()
	defer p.Unlock()

	log.Debug("~> PebbleKVSource.init()", log.Ctx{"dbname": p.name})
	if p.alive {
		return
	}
	if p.name == "" {
		log.Crit("no name set to db")
	}

	cache := pebble.NewCache(128 << 20)
	defer cache.Unref()
	opts := &pebble.Options{
		Cache:                       cache,
		Comparer:                    mvccComparer,
		DisableWAL:                  false,
		FormatMajorVersion:          pebble.FormatNewest,
		L0CompactionThreshold:       2,
		L0StopWritesThreshold:       1000,
		LBaseMaxBytes:               64 << 20, // 64 MB
		Levels:                      make([]pebble.LevelOptions, 7),
		MaxConcurrentCompactions:    3,
		MaxOpenFiles:                16384,
		MemTableSize:                64 << 20,
		MemTableStopWritesThreshold: 4,
		Merger: &pebble.Merger{
			Name: "cockroach_merge_operator",
		},
	}

	for i := 0; i < len(opts.Levels); i++ {
		l := &opts.Levels[i]
		l.BlockSize = 32 << 10       // 32 KB
		l.IndexBlockSize = 256 << 10 // 256 KB
		l.FilterPolicy = bloom.FilterPolicy(10)
		l.FilterType = pebble.TableFilter
		if i > 0 {
			l.TargetFileSize = opts.Levels[i-1].TargetFileSize * 2
		}
		l.EnsureDefaults()
	}
	opts.Levels[6].FilterPolicy = nil
	opts.FlushSplitBytes = opts.Levels[0].TargetFileSize

	opts.EnsureDefaults()

	//if verbose {
	//	opts.EventListener = pebble.MakeLoggingEventListener(nil)
	//	opts.EventListener.TableDeleted = nil
	//	opts.EventListener.TableIngested = nil
	//	opts.EventListener.WALCreated = nil
	//	opts.EventListener.WALDeleted = nil
	//}
	var err error
	log.Info("Opening db", log.Ctx{"dbname": p.name})
	dbPath := p.getPath()
	dir := path.Dir(dbPath)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0666); err != nil {
			panic(err)
		}
	}
	log.Debug("Open existing db or create new db ", log.Ctx{"dbname": p.name})
	p.db, err = pebble.Open(dbPath, opts)
	if err != nil {
		log.Crit("Failed to open db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}
	p.alive = true

	log.Debug("<~ PebbleKVSource.init()", log.Ctx{"dbname": p.name})
}

func (p *PebbleKv) Reset() {
	p.Close()
	err := os.RemoveAll(p.getPath())
	if err != nil {
		log.Crit("Failed to reset db", log.Ctx{"remove db": err.Error()})
	}
	p.Init()
}

func (p *PebbleKv) IsAlive() bool {
	return p.alive
}

func (p *PebbleKv) Put(key, val []byte) {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> PebbleKVSource.put():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(val)})
	var err error
	if val != nil {
		if p.db == nil {

		} else {
			err = p.db.Set(key, val, p.writeOpt)
		}
	} else {
		err = p.db.Delete(key, p.writeOpt)
	}
	if err != nil {
		log.Crit("Failed to put into db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}
	log.Trace("<~ PebbleKVSource.put():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(val)})
}

func (p *PebbleKv) Get(key []byte) []byte {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> PebbleKVSource.get():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})

	ret, closer, err := p.db.Get(key)
	if closer != nil {
		retCopy := make([]byte, len(ret))
		copy(retCopy, ret)
		ret = retCopy
		closer.Close()
	}
	if errors.Is(err, pebble.ErrNotFound) {
		return nil
	}
	if err != nil {
		log.Crit("Failed to get from db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}

	log.Trace("<~ PebbleKVSource.get():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key), "val": len(ret)})

	return ret
}

func (p *PebbleKv) Delete(key []byte) {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> PebbleKVSource.delete():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})
	err := p.db.Delete(key, p.writeOpt)
	if err != nil {
		log.Crit("Failed to delete from db", log.Ctx{"dbname": p.name, "err": err.Error()})
	}

	log.Trace("<~ PebbleKVSource.delete():",
		log.Ctx{"dbname": p.name, "key": hex.EncodeToString(key)})
}

func (p *PebbleKv) Keys() [][]byte {
	p.RLock()
	defer p.RUnlock()

	log.Trace("~> PebbleKVSource.keys():", log.Ctx{"dbname": p.name})
	var keys [][]byte
	iter := p.db.NewIter(nil)
	for iter.First(); iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
	}
	if err := iter.Close(); err != nil {
		log.Crit("Failed to close iterator", log.Ctx{"dbname": p.name, "err": err.Error()})
	}

	log.Trace("<~ PebbleKVSource.keys():", log.Ctx{"dbname": p.name})
	return keys
}

func (p *PebbleKv) FetchPrefix(key []byte, f db.FetchFunc) {
	p.RLock()
	defer p.RUnlock()

	iter := p.db.NewIter(prefixIterOptions(key))
	for iter.First(); iter.Valid(); iter.Next() {
		if f(utils.Copy2(iter.Key()), utils.Copy2(iter.Value())) {
			break
		}
	}
	if err := iter.Close(); err != nil {
		log.Crit("Failed to close prefix iterator", log.Ctx{"dbname": p.name, "err": err.Error()})
	}
}

func (p *PebbleKv) PrefixKeyLookup(key []byte) [][]byte {
	var keyList [][]byte
	p.FetchPrefix(key, func(k, v []byte) bool {
		keyList = append(keyList, k)
		return false
	})
	return keyList
}

func (p *PebbleKv) PrefixValueLookup(key []byte) [][]byte {
	var valueList [][]byte
	p.FetchPrefix(key, func(k, v []byte) bool {
		valueList = append(valueList, v)
		return false
	})
	return valueList
}

func (p *PebbleKv) getPath() string {
	return path.Join(p.config.StoreDir(), p.name)
}

func (p *PebbleKv) backupPath() string {
	return path.Join(p.config.StoreDir(), "backup", p.name)
}

var keyUpperBound = func(b []byte) []byte {
	end := make([]byte, len(b))
	copy(end, b)
	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil // no upper-bound
}

var prefixIterOptions = func(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: keyUpperBound(prefix),
	}
}
