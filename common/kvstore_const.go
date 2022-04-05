package common

type DatabaseName int

const (
	DB_INDEX DatabaseName = iota // Block index

	DB_BLOCK // Block raw data.

	DB_TIME // Time related block.

	DB_ORPHANIND // Orphan block index

	DB_SNAPSHOT
)

const (
	SETTING_STATS      byte = 0x10
	TIME_HASH_INFO     byte = 0x20
	HASH_BLOCK_INFO    byte = 0x30
	SUMS_BLOCK_INFO    byte = 0x40
	OURS_BLOCK_INFO    byte = 0x50
	SETTING_TOP_STATUS byte = 0x60
	SNAPSHOT_BOOT      byte = 0x70
	BLOCK_HEIGHT       byte = 0x80
	SNAPSHOT_PRESEED   byte = 0x90
	SUM_FILE_NAME           = "sums.dat"
)
