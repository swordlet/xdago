package pebbledb

import "xdago/db"

type PebbleKv struct {
	name  string
	alive bool
}

func (p *PebbleKv) GetName() string {
	return p.name
}

func (p *PebbleKv) SetName(name string) {
	p.name = name
}

func (p *PebbleKv) Close() {

}

func (p *PebbleKv) Init() error {
	return nil
}

func (p *PebbleKv) Reset() {

}

func (p *PebbleKv) IsAlive() bool {
	return p.alive
}

func (p *PebbleKv) Put(key, val []byte) error {
	return nil
}

func (p *PebbleKv) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (p *PebbleKv) Delete(key []byte) error {
	return nil
}

func (p *PebbleKv) Keys() ([][]byte, error) {
	return nil, nil
}

func (p *PebbleKv) PrefixKeyLookup(key []byte) ([][]byte, error) {
	return nil, nil
}

func (p *PebbleKv) FetchPrefix(key []byte, f db.FetchFunc) error {
	return nil
}

func (p *PebbleKv) PrefixValueLookup(key []byte) ([][]byte, error) {
	return nil, nil
}
