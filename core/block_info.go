package core

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"xdago/common"
	"xdago/snapshot"
	"xdago/utils"
)

type BlockInfo struct {
	Type        uint64
	Flags       int
	Height      uint64
	Difficulty  big.Int
	Ref         []byte
	MaxDiffLink []byte
	Fee         uint64
	Remark      [common.XDAG_FIELD_SIZE]byte
	Hash        [common.XDAG_HASH_SIZE]byte
	HashLow     [common.XDAG_HASH_SIZE]byte
	Amount      uint64
	Timestamp   uint64
	IsSnapshot  bool
	SnapInfo    snapshot.Info
}

func (bi BlockInfo) Equals(o BlockInfo) bool {
	return bi.Type == o.Type &&
		bi.Flags == o.Flags &&
		bi.Height == o.Height &&
		bi.Amount == o.Amount &&
		bi.Timestamp == o.Timestamp &&
		bi.Hash == o.Hash
}

func (bi BlockInfo) ToString() string {
	return "BlockInfo{" + "height=" + strconv.FormatUint(bi.Height, 10) +
		", hash=" + hex.EncodeToString(bi.Hash[:]) +
		", hashlow=" + hex.EncodeToString(bi.HashLow[:]) +
		", amount=" + strconv.FormatFloat(utils.Amount2xdag(bi.Amount), 'f', 2, 64) +
		", type=" + strconv.FormatUint(bi.Type, 10) +
		", difficulty=" + bi.Difficulty.String() +
		", ref=" + hex.EncodeToString(bi.Ref) +
		", maxDiffLink=" + hex.EncodeToString(bi.MaxDiffLink) +
		", flags=" + strconv.FormatUint(uint64(bi.Flags), 16) +
		", fee=" + strconv.FormatUint(bi.Fee, 10) +
		", timestamp=" + strconv.FormatUint(bi.Timestamp, 10) +
		", remark=" + string(bi.Remark[:]) +
		", isSnapshot=" + strconv.FormatBool(bi.IsSnapshot) +
		", snapshotInfo=" + bi.SnapInfo.ToString() +
		"}"
}
