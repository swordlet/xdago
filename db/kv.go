package db

type DatabaseName int

const (
	INDEX DatabaseName = iota // Block index

	BLOCK // Block raw data.

	TIME // Time related block.

	ORPHANIND // Orphan block index

	SNAPSHOT
)

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
	GetDB(name DatabaseName) IKVSource
	Close()
}
