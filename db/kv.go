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
	Init() error
	Close()
	Reset()
	Put(K, V []byte) error
	Get(K []byte) ([]byte, error)
	Delete(K []byte) error
	Keys() ([][]byte, error)
	PrefixKeyLookup(key []byte) ([][]byte, error)
	FetchPrefix(key []byte, f FetchFunc) error
	PrefixValueLookup(key []byte) ([][]byte, error)
}

type IDataFactory interface {
	GetDB(name DatabaseName) IKVSource
	Close()
}
