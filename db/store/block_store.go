package store

import "xdago/db"

type pebbleKv struct {
	name  string
	alive bool
}

func (p *pebbleKv) GetName() string {
	return p.name
}

func (p *pebbleKv) SetName(name string) {
	p.name = name
}

func (p *pebbleKv) Close() {

}

func (p *pebbleKv) Init() error {
	return nil
}

func (p *pebbleKv) Reset() {

}

func (p *pebbleKv) IsAlive() bool {
	return p.alive
}

func (p *pebbleKv) Put(key, val []byte) error {
	return nil
}

func (p *pebbleKv) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (p *pebbleKv) Delete(key []byte) error {
	return nil
}

func (p *pebbleKv) Keys() ([][]byte, error) {
	return nil, nil
}

func (p *pebbleKv) PrefixKeyLookup(key []byte) ([][]byte, error) {
	return nil, nil
}

func (p *pebbleKv) FetchPrefix(key []byte, f db.FetchFunc) error {
	return nil
}

func (p *pebbleKv) PrefixValueLookup(key []byte) ([][]byte, error) {
	return nil, nil
}
