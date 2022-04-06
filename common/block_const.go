package common

type Hash [XDAG_HASH_SIZE]byte
type Field [XDAG_FIELD_SIZE]byte
type RawBlock [XDAG_BLOCK_SIZE]byte
type Signature [XDAG_FIELD_SIZE * 2]byte

const (
	XDAG_BLOCK_FIELDS = 16
	XDAG_BLOCK_SIZE   = 512
	XDAG_FIELD_SIZE   = 32
	XDAG_HASH_SIZE    = 32
)
const MAX_LINKS = 15

type FieldType byte

const (
	// nonce字段
	XDAG_FIELD_NONCE FieldType = iota
	// 头部字段
	XDAG_FIELD_HEAD
	// 输入
	XDAG_FIELD_IN
	// 输入
	XDAG_FIELD_OUT
	// 输入签名
	XDAG_FIELD_SIGN_IN
	// 输出签名
	XDAG_FIELD_SIGN_OUT
	XDAG_FIELD_PUBLIC_KEY_0
	XDAG_FIELD_PUBLIC_KEY_1
	XDAG_FIELD_HEAD_TEST
	XDAG_FIELD_REMARK
	XDAG_FIELD_RESERVE1
	XDAG_FIELD_RESERVE2
	XDAG_FIELD_RESERVE3
	XDAG_FIELD_RESERVE4
	XDAG_FIELD_RESERVE5
	XDAG_FIELD_RESERVE6
)

var EmptyField Field
var EmptyHash Hash
var EmptyXdagSignature Signature
var EmptyXdagBlock RawBlock
