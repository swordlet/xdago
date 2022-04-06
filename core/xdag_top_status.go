package core

import "math/big"

type XDAGTopStatus struct {
	Top        []byte
	TopDiff    *big.Int
	PreTop     []byte
	PreTopDiff *big.Int
}

func NewXDAGTopStatus() XDAGTopStatus {
	return XDAGTopStatus{
		TopDiff:    big.NewInt(0),
		PreTopDiff: big.NewInt(0),
	}
}
