//go:build pebble && !rocksdb

////go:build rocksdb && !pebble
//switch for KV store

package config

import (
	"github.com/magiconair/properties/assert"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/log"
	"xdago/utils"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
}

func TestTestNetConfig(t *testing.T) {
	c := TestNetConfig()

	assert.Equal(t, c.RootDir(), "testnet")
	assert.Equal(t, c.TelnetIp(), "127.0.0.1")
	assert.Equal(t, c.PoolRation(), float64(5))
	assert.Equal(t, c.MaxInboundConnectionsPerIp(), 16)
	assert.Equal(t, c.Libp2pPort(), 12345)
	assert.Equal(t, c.IsBootNode(), false)
	assert.Equal(t, c.GlobalMinerLimit(), 8192)
	assert.Equal(t, c.MaxMinerPerAccount(), 512)
	assert.Equal(t, c.RpcEnabled(), false)
	if c.RpcEnabled() {
		assert.Equal(t, c.RpcPortHttp(), 10001)
	} else {
		assert.Equal(t, c.RpcPortHttp(), 0)
	}
	assert.Equal(t, c.ApolloForkHeight(), uint64(196250))
	assert.Equal(t, c.XdagEra(), uint64(0x16900000000))
	assert.Equal(t, c.StoreDir(), "testnet/pebble/xdagdb")
	assert.Equal(t, utils.Amount2xdag(c.MainStartAmount()), 1024.0)
	assert.Equal(t, utils.Amount2xdag(c.ApolloForkAmount()), 128.0)

}
