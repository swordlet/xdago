package core

import (
	"math/big"
	"xdago/common"
)

type XdagExtStats struct {
	HashRateTotal    []*big.Int
	HashRateOurs     []*big.Int
	HashRateLastTime uint64
}

func NewXdagExtStats() *XdagExtStats {
	return &XdagExtStats{
		HashRateTotal: make([]*big.Int, common.HASH_RATE_LAST_MAX_TIME, common.HASH_RATE_LAST_MAX_TIME),
		HashRateOurs:  make([]*big.Int, common.HASH_RATE_LAST_MAX_TIME, common.HASH_RATE_LAST_MAX_TIME),
	}
}
