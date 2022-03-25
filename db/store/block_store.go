package store

import (
	"encoding/binary"
	"xdago/utils"
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

func GetTimeKey(timestamp uint64, hashLow []byte) []byte {
	t := timestamp >> 16
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], t)
	if hashLow == nil {
		return utils.MergeBytes([]byte{TIME_HASH_INFO}, b[:])
	} else {
		return utils.MergeBytes([]byte{TIME_HASH_INFO}, b[:], hashLow)
	}
}
