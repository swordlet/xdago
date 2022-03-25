package crypto

import "crypto/sha256"

func HashTwice(input []byte) [32]byte {
	h := sha256.Sum256(input)
	return sha256.Sum256(h[:])
}
