package config

import "C"
import (
	"github.com/spf13/viper"
	"net"
	"xdago/core"
	"xdago/log"
)

var defaultConfig = Config{
	telnetIp:                   "127.0.0.1",
	telnetPort:                 7001,
	maxShareCountPerChannel:    20,
	awardEpoch:                 0xf,
	waitEpoch:                  10,
	maxConnections:             1024,
	maxInboundConnectionsPerIp: 8,
	connectionTimeout:          10000,
	connectionReadTimeout:      10000,
	storeMaxOpenFiles:          1024,
	storeMaxThreads:            1,
	storeFromBackup:            false,
	storeBackupDir:             "./testdate",
	enableRefresh:              false,
	ttl:                        5,
	rpcEnabled:                 false,
	snapshotEnabled:            false,
}

func initConfig(rootDir, configName string) *Config {
	c := &defaultConfig
	c.SetRootDir(rootDir)
	c.SetConfigName(configName)
	c.getSetting()
	c.setDir()
	return c
}

func (c *Config) getSetting() {
	v := viper.New()
	v.AddConfigPath(c.rootDir)
	v.SetConfigName(c.configName)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.SetDefault("admin.telnet.ip", "127.0.0.1")
	c.telnetIp = v.GetString("admin.telnet.ip")
	v.SetDefault("admin.telnet.port", 6001)
	c.telnetPort = v.GetInt("admin.telnet.port")
	c.telnetPassword = v.GetString("admin.telnet.password")

	v.SetDefault("pool.ip", "127.0.0.1")
	c.poolIp = v.GetString("pool.ip")
	v.SetDefault("pool.port", 7001)
	c.poolPort = v.GetInt("pool.port")
	v.SetDefault("pool.tag", "xdago")
	c.poolTag = v.GetString("pool.tag")

	c.poolRation = v.GetFloat64("pool.poolRation")
	c.rewardRation = v.GetFloat64("pool.rewardRation")
	c.fundRation = v.GetFloat64("pool.fundRation")
	c.directRation = v.GetFloat64("pool.directRation")

	v.SetDefault("node.ip", "127.0.0.1")
	c.nodeIp = v.GetString("node.ip")
	v.SetDefault("node.port", 8001)
	c.nodePort = v.GetInt("node.port")
	c.maxInboundConnectionsPerIp = v.GetInt("node.maxInboundConnectionsPerIp")

	whiteIpArray := v.GetStringSlice("node.whiteIPs")
	if len(whiteIpArray) > 0 {
		log.Debug("Accessing IP count", len(whiteIpArray))
		for _, address := range whiteIpArray {
			_, err := net.ResolveTCPAddr("tcp4", address)
			if err != nil {
				log.Crit("Parse config white IPs", log.Ctx{"address": address})
			} else {
				c.whiteIPList = append(c.whiteIPList, address)
			}
		}
	}

	c.libp2pPort = v.GetInt("node.libp2p.port")
	c.libp2pPrivkey = v.GetString("node.libp2p.privkey")
	c.isBootNode = v.GetBool("node.libp2p.isbootnode")

	bootNodeList := v.GetStringSlice("node.libp2p.bootnode")
	if len(bootNodeList) > 0 {
		c.bootNodes = append([]string{}, bootNodeList...)
	}

	c.globalMinerLimit = v.GetInt("miner.globalMinerLimit")
	c.globalMinerChannelLimit = v.GetInt("miner.globalMinerChannelLimit")
	c.maxConnectPerIp = v.GetInt("miner.maxConnectPerIp")
	c.maxMinerPerAccount = v.GetInt("miner.maxMinerPerAccount")

	// rpc
	v.SetDefault("rpc.enabled", false)
	v.SetDefault("rpc.http.host", "127.0.0.1")
	v.SetDefault("rpc.http.port", 10001)
	v.SetDefault("rpc.ws.port", 10002)
	c.rpcEnabled = v.GetBool("rpc.enabled")
	if c.rpcEnabled {
		c.rpcHost = v.GetString("rpc.http.host")
		c.rpcPortHttp = v.GetInt("rpc.http.port")
		c.rpcPortWs = v.GetInt("rpc.ws.port")
	}

}

func MainNetConfig() *Config {
	c := initConfig("mainnet", "mainnet-config.json")

	c.whitelistUrl = "https://raw.githubusercontent.com/XDagger/xdag/master/client/netdb-white.txt"

	c.xdagEra = 0x16940000000
	c.mainStartAmount = 1 << 42

	c.apolloForkHeight = 1017323
	c.apolloForkAmount = 1 << 39
	c.xdagFieldHeader = core.XDAG_FIELD_HEAD

	c.dnetKeyFile = c.rootDir + "/dnet_keys.bin"
	c.walletKeyFile = c.rootDir + "/wallet.dat"

	c.walletFilePath = c.rootDir + "/wallet/" + WALLET_FILE_NAME

	return c
}

func DevNetConfig() *Config {
	c := initConfig("devnet", "devnet-config.json")
	c.whitelistUrl = ""

	c.waitEpoch = 1

	c.xdagEra = 0x16900000000
	c.mainStartAmount = 1 << 42

	c.apolloForkHeight = 1000
	c.apolloForkAmount = 1 << 39
	c.xdagFieldHeader = core.XDAG_FIELD_HEAD_TEST

	c.dnetKeyFile = c.rootDir + "/dnet_keys.bin"
	c.walletKeyFile = c.rootDir + "/wallet-testnet.dat"

	c.walletFilePath = c.rootDir + "/wallet/" + WALLET_FILE_NAME

	return c
}

func TestNetConfig() *Config {
	c := initConfig("testnet", "testnet-config.json")
	c.whitelistUrl = "https://raw.githubusercontent.com/XDagger/xdag/master/client/netdb-white-testnet.txt"

	// testnet wait 1 epoch
	c.waitEpoch = 1

	c.xdagEra = 0x16900000000
	c.mainStartAmount = 1 << 42

	c.apolloForkHeight = 196250

	c.apolloForkAmount = 1 << 39
	c.xdagFieldHeader = core.XDAG_FIELD_HEAD_TEST

	c.dnetKeyFile = c.rootDir + "/dnet_keys.bin"
	c.walletKeyFile = c.rootDir + "/wallet-testnet.dat"

	c.walletFilePath = c.rootDir + "/wallet/" + WALLET_FILE_NAME

	return c
}
