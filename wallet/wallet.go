package wallet

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"xdago/common"
	"xdago/config"
	"xdago/crypto"
	"xdago/log"
	"xdago/secp256k1"
	"xdago/utils"
)

const (
	VERSION              = 4
	SALT_LENGTH          = 16
	BCRYPT_COST          = 12
	MNEMONIC_PASS_PHRASE = ""
)

type Wallet struct {
	sync.RWMutex
	file             string
	config           *config.Config
	accountsHash     []common.Hash160
	accountsKey      []*secp256k1.PrivateKey
	password         string
	mnemonicPhrase   string
	nextAccountIndex uint32
}

func NewWallet(config *config.Config) Wallet {
	return Wallet{
		file:         config.WalletFilePath(),
		config:       config,
		accountsHash: make([]common.Hash160, 0),
		accountsKey:  make([]*secp256k1.PrivateKey, 0),
	}
}

func (w *Wallet) Exists() bool {
	_, err := os.Stat(w.file)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

func (w *Wallet) Delete() error {
	w.Lock()
	defer w.Unlock()
	return os.Remove(w.file)
}

func (w *Wallet) GetFile() string {
	w.RLock()
	defer w.RUnlock()
	return w.file
}

func (w *Wallet) LockWallet() {
	w.Lock()
	defer w.Unlock()
	w.password = ""
	w.accountsHash = make([]common.Hash160, 0)
	w.accountsKey = make([]*secp256k1.PrivateKey, 0)
}

func (w *Wallet) GetDefKey() *secp256k1.PrivateKey {
	w.RLock()
	defer w.RUnlock()
	if len(w.accountsKey) > 0 {
		return w.accountsKey[0]
	}
	return nil
}

func (w *Wallet) UnlockWallet(password string) bool {
	w.Lock()
	defer w.Unlock()

	if len(password) == 0 {
		log.Crit("password can not be null")
	}

	if w.Exists() {
		data, err := ioutil.ReadFile(w.file)
		if err != nil {
			log.Crit("read wallet file failed", log.Ctx{"err": err.Error()})
		}
		r := utils.NewSimpleReader(data)
		var version uint32
		r.ReadInt(binary.BigEndian, &version)
		switch version {
		// only version 4
		case 4:
			salt := readBytes(r)
			if r.Error() != nil {
				log.Crit("parse wallet salt failed", log.Ctx{"err": r.Error().Error()})
			}
			key, keyErr := crypto.GenerateFromPassword(salt, []byte(password), BCRYPT_COST)
			if keyErr != nil {
				log.Crit("generate wallet Decrypt key failed", log.Ctx{"err": err.Error()})
			}
			newAccounts, readErr := readAccounts(key, r, true, version)
			if readErr != nil {
				return false
			}
			success := w.readHdSeed(key, r)
			if !success {
				return false
			}
			if r.Error() != nil {
				log.Crit("parse wallet private keys failed", log.Ctx{"err": r.Error().Error()})
			}

			w.accountsHash = make([]common.Hash160, 0)
			w.accountsKey = make([]*secp256k1.PrivateKey, 0)
			for _, account := range newAccounts {
				w.accountsKey = append(w.accountsKey, account)
				w.accountsHash = append(w.accountsHash, crypto.ToBytesAddress(account))
			}
		default:
			log.Crit("wallet version error")
		}
	}
	w.password = password
	return true
}

func readAccounts(key []byte, r *utils.SimpleReader, vlq bool, version uint32) ([]*secp256k1.PrivateKey, error) {
	var keys []*secp256k1.PrivateKey
	var total uint32
	r.ReadInt(binary.BigEndian, &total)
	for i := 0; i < int(total); i++ {
		iv := readBytes(r)
		pvKeyBytes, err := crypto.AesDecrypt(readBytes(r), key, iv)
		if err != nil {
			log.Error("decrypt wallet private key failed", log.Ctx{"err": err.Error()})
			return nil, err
		}
		privateKey := secp256k1.PrivKeyFromBytes(pvKeyBytes)
		keys = append(keys, privateKey)
	}
	return keys, nil
}

func (w *Wallet) writeAccounts(key []byte, wr *utils.SimpleWriter) {
	wr.WriteInt(binary.BigEndian, uint32(len(w.accountsKey)))
	for _, pk := range w.accountsKey {
		iv := make([]byte, 16)
		rand.Read(iv)
		writeBytes(iv, wr)
		encKey, err := crypto.AesEncrypt(pk.Serialize(), key, iv)
		if err != nil {
			log.Crit("encrypt wallet private key failed", log.Ctx{"err": err.Error()})
		}
		writeBytes(encKey, wr)
	}
}

func (w *Wallet) readHdSeed(key []byte, r *utils.SimpleReader) bool {
	iv := readBytes(r)
	decryptBites, err := crypto.AesDecrypt(readBytes(r), key, iv)
	if err != nil {
		log.Error("decrypt wallet hd seed failed", log.Ctx{"err": err.Error()})
		return false
	}

	r2 := utils.NewSimpleReader(decryptBites)
	size := bytes2Size(r2)
	w.mnemonicPhrase = string(r2.ReadCString(int(size)))
	r2.ReadInt(binary.BigEndian, &w.nextAccountIndex)
	if r2.Error() != nil {
		log.Crit("parse wallet mnemonic failed", log.Ctx{"err": r2.Error().Error()})
	}
	return true
}

func (w *Wallet) writeHdSeed(key []byte, wr *utils.SimpleWriter) {
	size := 4 + len(w.mnemonicPhrase) + 4
	wr2 := utils.NewSimpleWriter(size)
	wr2.WriteBytes(size2bytes(uint32(len(w.mnemonicPhrase))))
	wr2.WriteBytes([]byte(w.mnemonicPhrase))
	wr2.WriteInt(binary.BigEndian, w.nextAccountIndex)
	if wr2.Error() != nil {
		log.Crit("write wallet hd seed failed", log.Ctx{"err": wr2.Error().Error()})
	}
	iv := make([]byte, 16)
	rand.Read(iv)
	writeBytes(iv, wr)
	enc, err := crypto.AesEncrypt(wr2.BytesUncheck(), key, iv)
	if err != nil {
		log.Crit("encrypt wallet private key failed", log.Ctx{"err": err.Error()})
	}
	writeBytes(enc, wr)
}

func (w *Wallet) Flush() bool {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	wr := utils.NewSimpleWriter(2048)
	wr.WriteInt(binary.BigEndian, uint32(VERSION))
	salt := make([]byte, SALT_LENGTH)
	rand.Read(salt)
	writeBytes(salt, wr)
	key, err := crypto.GenerateFromPassword(salt, []byte(w.password), BCRYPT_COST)
	if err != nil {
		log.Crit("generate wallet encrypt key failed", log.Ctx{"err": err.Error()})
	}
	w.writeAccounts(key, wr)
	w.writeHdSeed(key, wr)
	if wr.Error() != nil {
		log.Crit("write wallet to bytes failed", log.Ctx{"err": wr.Error().Error()})
	}

	dir := path.Dir(w.file)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0666); err != nil {
			log.Crit("create wallet dir failed", log.Ctx{"err": err.Error()})
		}
	}
	err = ioutil.WriteFile(w.file, wr.BytesUncheck(), 0666)
	if err != nil {
		log.Crit("flush wallet data failed", log.Ctx{"err": err.Error()})
	}

	return true
}

func (w *Wallet) IsLocked() bool {
	w.RLock()
	defer w.RUnlock()

	return w.password == ""
}

func (w *Wallet) IsUnLocked() bool {
	return !w.IsLocked()
}

func (w *Wallet) requireUnlocked() {
	if w.password == "" {
		log.Crit("wallet is locked")
	}
}

func (w *Wallet) GetAccounts() []*secp256k1.PrivateKey {
	w.RLock()
	defer w.RUnlock()
	w.requireUnlocked()

	return w.accountsKey
}
func (w *Wallet) SetAccounts(keys []*secp256k1.PrivateKey) {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	w.accountsKey = keys
}

func (w *Wallet) GetAccount(index int) *secp256k1.PrivateKey {
	w.RLock()
	defer w.RUnlock()
	w.requireUnlocked()

	if index < len(w.accountsKey) {
		return w.accountsKey[index]
	}
	return nil
}

func (w *Wallet) GetAccountByAddress(address common.Hash160) *secp256k1.PrivateKey {
	w.RLock()
	defer w.RUnlock()
	w.requireUnlocked()

	for i, addr := range w.accountsHash {
		if addr == address {
			return w.accountsKey[i]
		}
	}
	return nil
}

func (w *Wallet) AddAccount(newKey *secp256k1.PrivateKey) {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	w.accountsKey = append(w.accountsKey, newKey)
	w.accountsHash = append(w.accountsHash, crypto.ToBytesAddress(newKey))
}

func (w *Wallet) AddAccountRandom() {
	w.Lock()
	defer w.Unlock()
	newKey, _ := secp256k1.GeneratePrivateKey()
	w.accountsKey = append(w.accountsKey, newKey)
	w.accountsHash = append(w.accountsHash, crypto.ToBytesAddress(newKey))
}

func (w *Wallet) AddAccounts(newKeys []*secp256k1.PrivateKey) {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	w.accountsKey = append(w.accountsKey, newKeys...)
	for _, newKey := range newKeys {
		w.accountsHash = append(w.accountsHash, crypto.ToBytesAddress(newKey))
	}
}

func (w *Wallet) RemoveAccountByKey(delKey *secp256k1.PrivateKey) bool {
	return w.RemoveAccountByAddress(crypto.ToBytesAddress(delKey))
}

func (w *Wallet) RemoveAccountByAddress(address common.Hash160) bool {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	index := -1
	for i, addr := range w.accountsHash {
		if addr == address {
			index = i
			break
		}
	}
	if index > -1 {
		w.accountsHash = append(w.accountsHash[:index], w.accountsHash[index+1:]...)
		w.accountsKey = append(w.accountsKey[:index], w.accountsKey[index+1:]...)
		return true
	}
	return false
}

func (w *Wallet) ChangePassword(newPassword string) {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()

	w.password = newPassword
}

func (w *Wallet) GetPassword() string {
	w.RLock()
	defer w.RUnlock()
	w.requireUnlocked()

	return w.password
}

func readBytes(r *utils.SimpleReader) []byte {
	size := bytes2Size(r)
	out := make([]byte, size)
	r.ReadBytes(out)
	return out
}

func writeBytes(b []byte, w *utils.SimpleWriter) {
	size := uint32(len(b))
	w.WriteBytes(size2bytes(size))
	w.WriteBytes(b)
}

func bytes2Size(r *utils.SimpleReader) uint32 {
	var size uint32
	for i := 0; i < 4; i++ {
		b := r.ReadOneByte()
		size = (size << 7) | uint32(b&0x7F)
		if (b & 0x80) == 0 {
			break
		}
	}
	return size
}

func size2bytes(size uint32) []byte {
	var b [4]byte
	i := 3
	b[i] = byte(size & 0x7f)
	size = size >> 7

	for size > 0 {
		i -= 1
		b[i] = byte(size & 0x7f)
		size = size >> 7
	}
	c := i
	for i < 4 {
		if i != 3 {
			b[i] = b[i] | 0x80
		}
		i += 1
	}
	return b[c:]
}

// ================
// HD wallet
// ================

func (w *Wallet) IsHdWalletInitialized() bool {
	return w.mnemonicPhrase != ""
}

func (w *Wallet) InitializeHdWallet(mnemonic string) {
	w.Lock()
	defer w.Unlock()
	w.mnemonicPhrase = mnemonic
	w.nextAccountIndex = 0
}

//GetSeed Returns the HD seed.
func (w *Wallet) GetSeed() []byte {
	return bip39.NewSeed(w.mnemonicPhrase, MNEMONIC_PASS_PHRASE)
}

//Derives a key based on the current HD account index, and put it into the wallet
func (w *Wallet) AddAccountWithNextHdKey() *secp256k1.PrivateKey {
	w.Lock()
	defer w.Unlock()
	w.requireUnlocked()
	if !w.IsHdWalletInitialized() {
		log.Crit("HD Seed is not initialized")
	}

	seed := w.GetSeed()
	masterKey, _ := bip32.NewMasterKey(seed)
	bip44Key := generateBip44Key(masterKey, w.nextAccountIndex)
	w.nextAccountIndex += 1
	key := secp256k1.PrivKeyFromBytes(bip44Key.Key)
	address := crypto.ToBytesAddress(key)
	w.accountsKey = append(w.accountsKey, key)
	w.accountsHash = append(w.accountsHash, address)
	return key
}

const HARDENED_BIT = 0x80000000

func generateBip44Key(masterKey *bip32.Key, index uint32) *bip32.Key {
	childIdx := []uint32{44 | HARDENED_BIT, common.XDAG_BIP44_CION_TYPE | HARDENED_BIT, 0 | HARDENED_BIT, 0, index}

	for _, child := range childIdx {
		masterKey, _ = masterKey.NewChildKey(child)
	}
	return masterKey
}
