package db

import "xdago/common"

type FetchFunc func([]byte, []byte) bool

type IKVSource interface {
	GetName() string
	SetName(name string)
	IsAlive() bool
	Init()
	Close()
	Reset()
	Put(K, V []byte)
	Get(K []byte) []byte
	Delete(K []byte)
	Keys() [][]byte
	PrefixKeyLookup(key []byte) [][]byte
	FetchPrefix(key []byte, f FetchFunc)
	PrefixValueLookup(key []byte) [][]byte
}

type IDataFactory interface {
	GetDB(name common.DatabaseName) *IKVSource
	Close()
}
