package core

import (
	"xdago/config"
	"xdago/secp256k1"
)

func GenerateAddressBlock(config *config.Config, key *secp256k1.PrivateKey, xdagTime uint64) *Block {
	b := NewBlock(config, xdagTime, nil, nil, false, nil, "", -1)
	b.SignOut(key)
	return b
}
