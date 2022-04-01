package snapshot

import (
	"encoding/hex"
	"strconv"
)

type Info struct {
	Type bool   // true PUBKEY false BLOCK_DATA
	Data []byte // 区块数据 or pubkey
}

func (si Info) ToString() string {
	return "SnapshotInfo{" +
		"type=" + strconv.FormatBool(si.Type) +
		", data=" + hex.EncodeToString(si.Data[:]) +
		"}"
}
