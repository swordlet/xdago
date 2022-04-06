package core

import "xdago/net/node"

type BlockWrapper struct {
	Block      *Block
	Ttl        int
	RemoteNode node.Node //记录区块接收节点
	Timestamp  uint64    //NO_PARENT waiting time
}

func NewBlockWrapper(block *Block, ttl int) BlockWrapper {
	return BlockWrapper{
		Block: block,
		Ttl:   ttl,
	}
}

func NewBlockWrapperWithNode(block *Block, ttl int, remote node.Node) BlockWrapper {
	return BlockWrapper{
		Block:      block,
		Ttl:        ttl,
		RemoteNode: remote,
	}
}
