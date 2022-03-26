//go:build rocksdb && !pebble

//conditional build

package config

func (c *Config) SetDir() {
	c.storeDir = c.rootDir + "/rocksdb/xdagdb"
	c.storeBackupDir = c.rootDir + "/rocksdb/xdagdb/backupdata"
	c.whiteListDir = c.rootDir + "/netdb-white.txt"
	c.netDBDir = c.rootDir + "/netdb.txt"
}
