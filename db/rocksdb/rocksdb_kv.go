//go:build rocksdb && !pebble

package rocksdb

import "xdago/db"

type RocksKv struct {
	name  string
	alive bool
}

func (p *RocksKv) GetName() string {
	return p.name
}

func (p *RocksKv) SetName(name string) {
	p.name = name
}

func (p *RocksKv) Close() {

}

func (p *RocksKv) Init() error {
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
