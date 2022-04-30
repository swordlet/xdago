package core

import (
	"math/big"
	"strconv"
	"xdago/utils"
)

type XDAGStats struct {
	Difficulty       *big.Int
	maxDifficulty    *big.Int
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

func (x *XDAGStats) MaxDifficulty() *big.Int {
	return x.maxDifficulty
}

func (x *XDAGStats) SetMaxDifficulty(maxDifficulty *big.Int) {
	if x.maxDifficulty.Cmp(maxDifficulty) < 0 {
		x.maxDifficulty = maxDifficulty
	}

}

func NewEmptyXDAGStats() *XDAGStats {
	return &XDAGStats{
		Difficulty:    big.NewInt(0),
		maxDifficulty: big.NewInt(0),
	}
}

func NewXDAGStats(maxDifficulty *big.Int, totalNBlocks, totalNMain, mainTime uint64,
	totalNHosts int) *XDAGStats {

	return &XDAGStats{
		maxDifficulty: maxDifficulty,
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
	if x.maxDifficulty != nil && remote.maxDifficulty != nil &&
		remote.maxDifficulty.Cmp(x.maxDifficulty) > 0 {
		x.maxDifficulty.Set(remote.maxDifficulty)
	}
}

func (x XDAGStats) ToString() string {
	return "XdagStatus[nmain:" + strconv.FormatUint(x.NMain, 10) +
		",totalmain:" + strconv.FormatUint(x.TotalNMain, 10) +
		",nblocks:" + strconv.FormatUint(x.NBlocks, 10) +
		",totalblocks:" + strconv.FormatUint(x.TotalNBlocks, 10) + "]"
}
