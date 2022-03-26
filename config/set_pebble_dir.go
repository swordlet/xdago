//go:build pebble && !rocksdb

//conditional build

package config

func (c *Config) SetDir() {
	c.storeDir = c.rootDir + "/pebble/xdagdb"
	c.storeBackupDir = c.rootDir + "/pebble/xdagdb/backupdata"
	c.whiteListDir = c.rootDir + "/netdb-white.txt"
	c.netDBDir = c.rootDir + "/netdb.txt"
}
