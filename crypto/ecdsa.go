package crypto

import (
	"xdago/common"
	"xdago/log"
	"xdago/secp256k1"
	"xdago/secp256k1/ecdsa"
)

func EcdsaSign(key *secp256k1.PrivateKey, hash []byte) (r, s [common.XDAG_FIELD_SIZE]byte) {
	signature := ecdsa.Sign(key, hash)
	serial := signature.Serialize()
	rLen := int(serial[3])
	serial = serial[4:]
	if rLen >= common.XDAG_FIELD_SIZE {
		copy(r[:], serial[rLen-common.XDAG_FIELD_SIZE:rLen])
	} else {
		copy(r[:rLen], serial[:rLen])
	}

	sLen := int(serial[rLen+1])
	serial = serial[rLen+2:]
	if sLen >= common.XDAG_FIELD_SIZE {
		copy(s[:], serial[sLen-common.XDAG_FIELD_SIZE:sLen])
	} else {
		copy(s[:sLen], serial[:sLen])
	}
	log.Debug("Sign")
	return
}

func EcdsaVerify(key *secp256k1.PublicKey, hash, r, s []byte) bool {
	var scalarR, scalarS secp256k1.ModNScalar
	if overflow := scalarR.SetByteSlice(r); overflow {
		log.Crit("ecdsa verify error", log.Ctx{"err": "set scalar R overflow"})
	}
	if overflow := scalarS.SetByteSlice(s); overflow {
		log.Crit("ecdsa verify error", log.Ctx{"err": "set scalar S overflow"})
	}
	signature := ecdsa.NewSignature(&scalarR, &scalarS)

	return signature.Verify(hash, key)
}
