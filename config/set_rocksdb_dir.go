//go:build rocksdb && !pebble

package config

func (c *Config) setDir() {
	c.storeDir = c.rootDir + "/rocksdb/xdagdb"
	c.storeBackupDir = c.rootDir + "/rocksdb/xdagdb/backupdata"
	c.whiteListDir = c.rootDir + "/netdb-white.txt"
	c.netDBDir = c.rootDir + "/netdb.txt"
}
