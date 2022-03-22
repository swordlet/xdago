package config

const (
	MAIN_CHAIN_PERIOD uint64 = 64 << 10

	//setmain设置区块为主块时标志该位
	BI_MAIN byte = 0x01

	//跟BI_MAIN差不多 不过BI_MAIN是确定的 BI_MAIN_CHAIN是还未确定的
	BI_MAIN_CHAIN byte = 0x02

	//区块被应用apply后可能会标志该标识位（因为有可能区块存在问题不过还是被指向了 但是会标示为拒绝状态）
	BI_APPLIED byte = 0x04

	//区块应用apply过后会置该标识位
	BI_MAIN_REF byte = 0x08

	//从孤块链中移除 即有区块链接孤块的时候 将孤块置为BI_REF
	BI_REF byte = 0x10

	//添加区块时如果该区块的签名可以用自身的公钥解 则说明该区块是自己的区块
	BI_OURS byte = 0x20

	//候补主块未持久化
	BI_EXTRA                byte   = 0x40
	BI_REMARK               byte   = 0x80
	SEND_PERIOD             uint64 = 2
	DNET_PKT_XDAG           int    = 0x8B
	BLOCK_HEAD_WORD         int    = 0x3fca9e2b
	REQUEST_BLOCKS_MAX_TIME uint64 = 1 << 20
	REQUEST_WAIT            uint64 = 64
	MAX_ALLOWED_EXTRA       uint64 = 65536
	FUND_ADDRESS            string = "FQglVQtb60vQv2DOWEUL7yh3smtj7g1s"
	//每一轮的确认数是16
	CONFIRMATIONS_COUNT int = 16
	MAIN_BIG_PERIOD_LOG int = 21

	WALLET_FILE_NAME string = "wallet.data"

	CLIENT_VERSION string = "0.4.6"

	//同步问题 分叉高度
	SYNC_FIX_HEIGHT uint64 = 0
)

type MessageType int

const (
	PRE_TOP MessageType = iota
	NEW_LINK
)
