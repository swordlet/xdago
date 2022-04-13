//go:build pebble && !rocksdb

////go:build rocksdb && !pebble
//conditional build switch for KV store

package wallet

import (
	"encoding/hex"
	"github.com/magiconair/properties/assert"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/config"
	"xdago/log"
	"xdago/secp256k1"
)

const (
	PRIVATE_KEY_STRING = "a392604efc2fad9c0b3da43b5f698a2e3f270f170d859912be0d54742275c5f6"
	PUBLIC_KEY_STRING  = "0x506bc1dc099358e5137292f4efdd57e400f29ba5132aa5d12b18dac1c1f6aab" +
		"a645c0b7b58158babbfa6c6cd5a48aa7340a8749176b120e8516216787a13dc76"
	PUBLIC_KEY_COMPRESS_STRING = "02506bc1dc099358e5137292f4efdd57e400f29ba5132aa5d12b18dac1c1f6aaba"
	ADDRESS                    = "b731bf10ed204f4ebc3d32ac88b7aa61b993fd59"
	PASSWORD                   = "Insecure Pa55w0rd"
	MNEMONIC                   = "scatter major grant return flee easy female jungle" +
		" vivid movie bicycle absent weather inspire carry"
)

func testInit() (*config.Config, string, *Wallet) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	c := config.DevNetConfig()
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
	pwd := "password"
	wallet := NewWallet(c)
	wallet.UnlockWallet(pwd)
	keyBytes, _ := hex.DecodeString(PRIVATE_KEY_STRING)
	privKey := secp256k1.PrivKeyFromBytes(keyBytes)
	wallet.SetAccounts([]*secp256k1.PrivateKey{privKey})
	wallet.Flush()
	wallet.LockWallet()
	return c, pwd, &wallet
}

func TestGetPassword(t *testing.T) {
	_, p, w := testInit()
	w.UnlockWallet(p)
	assert.Equal(t, w.GetPassword(), "password")
}
