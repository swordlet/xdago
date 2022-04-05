package core

import (
	"math/big"
	"strconv"
	"xdago/utils"
)

type XDAGStats struct {
	Difficulty       *big.Int
	MaxDifficulty    *big.Int
	NBlocks          uint64
	TotalNBlocks     uint64
	NMain            uint64
	TotalNMain       uint64
	NHosts           int
	TotalNHosts      int
	NWaitSync        uint64
	NnoRef           uint64
	NExtra           uint64
	MainTime         uint64
	Balance          uint64
	GlobalMiner      []byte
	OurLastBlockHash []byte
}

func NewXDAGStats(maxDifficulty *big.Int, totalNBlocks, totalNMain, mainTime uint64,
	totalNHosts int) XDAGStats {

	return XDAGStats{
		MaxDifficulty: maxDifficulty,
		TotalNBlocks:  totalNBlocks,
		TotalNMain:    totalNMain,
		MainTime:      mainTime,
		TotalNHosts:   totalNHosts,
	}
}

func (x *XDAGStats) Update(remote XDAGStats) {
	x.TotalNHosts = utils.MaxInt(x.TotalNHosts, remote.TotalNHosts)
	x.TotalNBlocks = utils.MaxUint64(x.TotalNBlocks, remote.NBlocks)
	x.TotalNMain = utils.MaxUint64(x.TotalNMain, remote.TotalNMain)
	if x.MaxDifficulty != nil && remote.MaxDifficulty != nil &&
		remote.MaxDifficulty.Cmp(x.MaxDifficulty) > 0 {
		x.MaxDifficulty.Set(remote.MaxDifficulty)
	}
}

func (x XDAGStats) ToString() string {
	return "XdagStatus[nmain:" + strconv.FormatUint(x.NMain, 10) +
		",totalmain:" + strconv.FormatUint(x.TotalNMain, 10) +
		",nblocks:" + strconv.FormatUint(x.NBlocks, 10) +
		",totalblocks:" + strconv.FormatUint(x.TotalNBlocks, 10) + "]"
}
