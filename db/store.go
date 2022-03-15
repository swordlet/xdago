package db

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
	Keys() [][]byte
	PrefixKeyLookup(key []byte) ([][]byte, error)
	FetchPrefix(key []byte, f FetchFunc) error
	PrefixValueLookup(key []byte) ([][]byte, error)
}
