//go:build pebble && !rocksdb

////go:build rocksdb && !pebble
//conditional build switch for KV store

package core

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"xdago/common"
	"xdago/config"
	"xdago/crypto"
	"xdago/log"
	"xdago/secp256k1"
	"xdago/utils"
)

func testInit() *config.Config {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	c := config.DevNetConfig()
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
	return c
}
func printBlockInfo(block *Block) {
	fmt.Printf("timestamp:%016x\n", block.GetTimestamp())
	fmt.Println("blockhash:", utils.Hash2String(block.GetHash()))
	fmt.Println("blockhashlow:", utils.Hash2String(block.GetHashLow()))
	fmt.Println("type:", utils.Type2String(block.GetType()))
	fmt.Println("inputs:", len(block.Inputs))
	if len(block.Inputs) > 0 {
		for _, address := range block.Inputs {
			fmt.Println("address data:", hex.EncodeToString(address.Data[:]))
			fmt.Println("address hashlow:", utils.Hash2String(address.HashLow))
			fmt.Println("address amount:", utils.Amount2xdag(address.Amount))
		}
	}
	fmt.Println("outputs:", len(block.Outputs))
	if len(block.Inputs) > 0 {
		for _, address := range block.Outputs {
			fmt.Println("address data:", hex.EncodeToString(address.Data[:]))
			fmt.Println("address hashlow:", utils.Hash2String(address.HashLow))
			fmt.Println("address amount:", utils.Amount2xdag(address.Amount))
		}
	}
	fmt.Println("keys size:", len(block.PubKeys))
	fmt.Println("verified keys size:", len(block.VerifiedKeys()))
	if len(block.PubKeys) > 0 {
		for _, key := range block.PubKeys {
			fmt.Println("pub key:", hex.EncodeToString(key.SerializeCompressed()))
		}
	}
	fmt.Println("out sign index:", block.GetOutsigIndex())
	fmt.Println("out sign r:", hex.EncodeToString(block.OutSig[:32]))
	fmt.Println("out sign s:", hex.EncodeToString(block.OutSig[32:]))
	fmt.Println("in sign size", len(block.InSigs))
	for _, sig := range block.InSigs {
		fmt.Println("in sign index", sig[64])
		fmt.Println("in sign r:", hex.EncodeToString(sig[:32]))
		fmt.Println("in sign s:", hex.EncodeToString(sig[32:64]))
	}
	if block.Nonce != common.EmptyField {
		fmt.Println("nonce:", hex.EncodeToString(block.Nonce[:]))
	}
	fmt.Println("xdag block:")
	for _, field := range block.GetXdagBlock().Fields {
		fmt.Println(hex.EncodeToString(field.Data[:]))
	}
}
func TestGenerateBlock(t *testing.T) {
	cfg := testInit()

	blockStrData := "000000000000000038324654050000004d3782fa780100000000000000000000" +
		"c86357a2f57bb9df4f8b43b7a60e24d1ccc547c606f2d7980000000000000000" +
		"afa5fec4f56f7935125806e235d5280d7092c6840f35b397000000000a000000" +
		"a08202c3f60123df5e3a973e21a2dd0418b9926a2eb7c4fc000000000a000000" +
		"08b65d2e2816c0dea73bf1b226c95c2ae3bc683574f559bbc5dd484864b1dbeb" +
		"f02a041d5f7ff83a69c0e35e7eeeb64496f76f69958485787d2c50fd8d9614e6" +
		"7c2b69c79eddeff5d05b2bfc1ee487b9c691979d315586e9928c04ab3ace15bb" +
		"3866f1a25ed00aa18dde715d2a4fc05147d16300c31fefc0f3ebe4d77c63fcbb" +
		"ec6ece350f6be4c84b8705d3b49866a83986578a3a20e876eefe74de0c094bac" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	blockRawData, _ := hex.DecodeString(blockStrData)
	block01 := NewBlockFromXdag(NewXdagBlock(blockRawData))
	printBlockInfo(&block01)

	fmt.Println(
		"=====================================first block use key1========================================")

	time := utils.GetMainTime()
	var pending []Address
	pending = append(pending, AddressFromHashLow(block01.GetHashLow()))

	tx01 := NewBlock(cfg, time, nil, pending, false, nil, "", -1)
	privKey01, _ := secp256k1.GeneratePrivateKey()
	tx01.SignOut(privKey01)

	printBlockInfo(&tx01)

	fmt.Println(
		"=====================================second block use key2========================================")

	tx02 := NewBlock(cfg, time, nil, pending, false, nil, "", -1)
	privKey02, _ := secp256k1.GeneratePrivateKey()
	tx02.SignOut(privKey02)

	printBlockInfo(&tx02)

	fmt.Println(
		"=====================================main block use key2========================================")
	var pendingMain []Address
	pendingMain = append(pendingMain, AddressFromHashLow(tx01.GetHashLow()))
	pendingMain = append(pendingMain, AddressFromHashLow(tx02.GetHashLow()))
	main := NewBlock(cfg, time, nil, pending, true, nil, "", -1)
	main.SignOut(privKey02)
	var minShare [32]byte
	rand.Read(minShare[:])
	main.Nonce = minShare
	printBlockInfo(&main)

	fmt.Println(
		"=====================================transaction1 block use key1========================================")
	var links01 []Address
	links01 = append(links01, AddressFromAmount(tx01.GetHashLow(), common.XDAG_FIELD_IN, 10<<24)) // key1
	links01 = append(links01, AddressFromAmount(tx02.GetHashLow(), common.XDAG_FIELD_OUT, 10<<24))
	var keys01 []secp256k1.PublicKey
	keys01 = append(keys01, *privKey01.PubKey())
	transaction01 := NewBlock(cfg, time, links01, nil, false, keys01, "", 0)
	// 跟输入用的同一把密钥
	transaction01.SignOut(privKey01)
	printBlockInfo(&transaction01)

	fmt.Println(
		"=====================================transaction2 block use key3========================================")
	var links02 []Address
	links02 = append(links02, AddressFromAmount(tx01.GetHashLow(), common.XDAG_FIELD_IN, 10<<24)) // key1
	links02 = append(links02, AddressFromAmount(tx02.GetHashLow(), common.XDAG_FIELD_OUT, 10<<24))
	var keys02 []secp256k1.PublicKey
	keys02 = append(keys02, *privKey01.PubKey())
	transaction02 := NewBlock(cfg, time, links02, nil, false, keys02, "", -1)
	// 跟输入用的不是同一把密钥
	privKey03, _ := secp256k1.GeneratePrivateKey()
	transaction02.SignIn(privKey01)
	transaction02.SignOut(privKey03)
	printBlockInfo(&transaction02)

	fmt.Println(
		"=====================================transaction3 block use key3========================================")
	var links03 []Address
	links03 = append(links03, AddressFromAmount(tx01.GetHashLow(), common.XDAG_FIELD_IN, 10<<24)) // key1
	links03 = append(links03, AddressFromAmount(tx02.GetHashLow(), common.XDAG_FIELD_IN, 10<<24))
	links03 = append(links03, AddressFromAmount(main.GetHashLow(), common.XDAG_FIELD_IN, 20<<24))
	var keys03 []secp256k1.PublicKey
	keys03 = append(keys03, *privKey01.PubKey())
	keys03 = append(keys03, *privKey02.PubKey())
	transaction03 := NewBlock(cfg, time, links03, nil, false, keys03, "", -1)
	// 跟输入用的不是同一把密钥
	transaction03.SignIn(privKey01)
	transaction03.SignIn(privKey02)
	transaction03.SignOut(privKey03)
	printBlockInfo(&transaction03)

	fmt.Println(
		"=====================================verify transaction01 sig========================================")

	input := []*Block{&tx01}
	fmt.Println("can use input?:", canUseInput(&transaction01, input))

	fmt.Println(
		"=====================================verify transaction02 sig========================================")
	fmt.Println("can use input?:", canUseInput(&transaction02, input))

	fmt.Println(
		"=====================================verify transaction03 sig========================================")
	input = append(input, &tx02)
	fmt.Println("can use input?:", canUseInput(&transaction03, input))

}

func canUseInput(transaction *Block, input []*Block) bool {
	keys := transaction.VerifiedKeys()
	for _, inBlock := range input {
		// 获取签名与hash
		subData := inBlock.GetSubRawData(inBlock.GetOutsigIndex())
		//fmt.Println(hex.EncodeToString(subData[:]))

		for _, key := range keys {
			hash := crypto.HashTwice(utils.MergeBytes(subData[:], key.SerializeCompressed()))
			if !crypto.EcdsaVerify(&key, hash[:], inBlock.OutSig[:32], inBlock.OutSig[32:]) {
				return false
			}
		}

	}
	return true
}
