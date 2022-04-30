package core

import (
	"xdago/common"
	"xdago/secp256k1"
)

type IBlockchain interface {
	GetPreSeed() common.Hash
	TryToConnect(block *Block) ImportResult
	CreateNewBlock(pairs map[Address]*secp256k1.PrivateKey, to []Address, mining bool, remark string) *Block
	GetBlockByHash(hash common.Hash, isRaw bool) *Block
	GetBlockByHeight(height uint64) *Block
	CheckNewMain()
	ListMainBlock(count int) []*Block
	ListMinedBlock(count int) []*Block
	GetMemOurBlocks() map[common.Hash]int
	GetXDAGStats() *XDAGStats
	GetXDAGTopStatus() *XDAGTopStatus
	GetSupply(nMain uint64) uint64
	GetBlockByTime(startTime, endTime uint64) []*Block

	//TODO:补充单元测试

	StartCheckMain()   // 启动检查主块链线程
	StopCheckMain()    // 关闭检查主块链线程
	RegisterListener() // 注册监听器
	GetXdagExtStats() XdagExtStats
}
