package config

import (
	"xdago/common"
	"xdago/crypto"
)

type Config struct {
	configName string

	// Admin spec
	telnetIp       string
	telnetPort     int
	telnetPassword string

	// Mining Pool spec
	poolIp       string
	poolPort     int
	poolTag      string
	poolRation   float64
	rewardRation float64
	fundRation   float64
	directRation float64

	globalMinerLimit        int
	globalMinerChannelLimit int
	maxConnectPerIp         int
	maxMinerPerAccount      int

	maxShareCountPerChannel int
	awardEpoch              int
	waitEpoch               int

	// Node spec
	nodeIp                     string
	nodePort                   int
	maxInboundConnectionsPerIp int
	maxConnections             int
	connectionTimeout          int
	connectionReadTimeout      int

	rootDir        string
	storeDir       string
	storeBackupDir string
	whiteListDir   string
	netDBDir       string

	storeMaxOpenFiles int
	storeMaxThreads   int
	storeFromBackup   bool
	originStoreDir    string

	whitelistUrl  string
	enableRefresh bool
	dnetKeyFile   string
	walletKeyFile string

	ttl          int
	dnetKeyBytes [2048]byte
	xkeys        *crypto.DnetKeys
	whiteIPList  []string

	// LibP2P spec
	libp2pPort    int
	libp2pPrivkey string
	isBootNode    bool
	bootNodes     []string

	// Wallet spec
	walletFilePath string

	// Xdag spec
	xdagEra          uint64
	xdagFieldHeader  common.FieldType
	mainStartAmount  uint64
	apolloForkHeight uint64
	apolloForkAmount uint64

	// Xdag RPC modules
	rpcEnabled  bool
	rpcHost     string
	rpcPortHttp int
	rpcPortWs   int
	//moduleDescriptions []ModuleDescription

	// Xdag Snapshot
	snapshotEnabled bool
	snapshotHeight  uint64
	snapshotTime    uint64
	isSnapshotJ     bool
}

func (c *Config) Xkeys() *crypto.DnetKeys {
	return c.xkeys
}

func (c *Config) SetXkeys(xkeys *crypto.DnetKeys) {
	c.xkeys = xkeys
}

// getter and setter

func (c *Config) ConfigName() string {
	return c.configName
}

func (c *Config) SetConfigName(configName string) {
	c.configName = configName
}

func (c *Config) TelnetIp() string {
	return c.telnetIp
}

func (c *Config) SetTelnetIp(telnetIp string) {
	c.telnetIp = telnetIp
}

func (c *Config) TelnetPort() int {
	return c.telnetPort
}

func (c *Config) SetTelnetPort(telnetPort int) {
	c.telnetPort = telnetPort
}

func (c *Config) TelnetPassword() string {
	return c.telnetPassword
}

func (c *Config) SetTelnetPassword(telnetPassword string) {
	c.telnetPassword = telnetPassword
}

func (c *Config) PoolIp() string {
	return c.poolIp
}

func (c *Config) SetPoolIp(poolIp string) {
	c.poolIp = poolIp
}

func (c *Config) PoolPort() int {
	return c.poolPort
}

func (c *Config) SetPoolPort(poolPort int) {
	c.poolPort = poolPort
}

func (c *Config) PoolTag() string {
	return c.poolTag
}

func (c *Config) SetPoolTag(tag string) {
	c.poolTag = tag
}

func (c *Config) PoolRation() float64 {
	return c.poolRation
}

func (c *Config) SetPoolRation(poolRation float64) {
	c.poolRation = poolRation
}

func (c *Config) RewardRation() float64 {
	return c.rewardRation
}

func (c *Config) SetRewardRation(rewardRation float64) {
	c.rewardRation = rewardRation
}

func (c *Config) FundRation() float64 {
	return c.fundRation
}

func (c *Config) SetFundRation(fundRation float64) {
	c.fundRation = fundRation
}

func (c *Config) DirectRation() float64 {
	return c.directRation
}

func (c *Config) SetDirectRation(directRation float64) {
	c.directRation = directRation
}

func (c *Config) GlobalMinerLimit() int {
	return c.globalMinerLimit
}

func (c *Config) SetGlobalMinerLimit(globalMinerLimit int) {
	c.globalMinerLimit = globalMinerLimit
}

func (c *Config) GlobalMinerChannelLimit() int {
	return c.globalMinerChannelLimit
}

func (c *Config) SetGlobalMinerChannelLimit(globalMinerChannelLimit int) {
	c.globalMinerChannelLimit = globalMinerChannelLimit
}

func (c *Config) MaxConnectPerIp() int {
	return c.maxConnectPerIp
}

func (c *Config) SetMaxConnectPerIp(maxConnectPerIp int) {
	c.maxConnectPerIp = maxConnectPerIp
}

func (c *Config) MaxMinerPerAccount() int {
	return c.maxMinerPerAccount
}

func (c *Config) SetMaxMinerPerAccount(maxMinerPerAccount int) {
	c.maxMinerPerAccount = maxMinerPerAccount
}

func (c *Config) MaxShareCountPerChannel() int {
	return c.maxShareCountPerChannel
}

func (c *Config) SetMaxShareCountPerChannel(maxShareCountPerChannel int) {
	c.maxShareCountPerChannel = maxShareCountPerChannel
}

func (c *Config) AwardEpoch() int {
	return c.awardEpoch
}

func (c *Config) SetAwardEpoch(awardEpoch int) {
	c.awardEpoch = awardEpoch
}

func (c *Config) WaitEpoch() int {
	return c.waitEpoch
}

func (c *Config) SetWaitEpoch(waitEpoch int) {
	c.waitEpoch = waitEpoch
}

func (c *Config) NodeIp() string {
	return c.nodeIp
}

func (c *Config) SetNodeIp(nodeIp string) {
	c.nodeIp = nodeIp
}

func (c *Config) NodePort() int {
	return c.nodePort
}

func (c *Config) SetNodePort(nodePort int) {
	c.nodePort = nodePort
}

func (c *Config) MaxInboundConnectionsPerIp() int {
	return c.maxInboundConnectionsPerIp
}

func (c *Config) SetMaxInboundConnectionsPerIp(maxInboundConnectionsPerIp int) {
	c.maxInboundConnectionsPerIp = maxInboundConnectionsPerIp
}

func (c *Config) MaxConnections() int {
	return c.maxConnections
}

func (c *Config) SetMaxConnections(maxConnections int) {
	c.maxConnections = maxConnections
}

func (c *Config) ConnectionTimeout() int {
	return c.connectionTimeout
}

func (c *Config) SetConnectionTimeout(connectionTimeout int) {
	c.connectionTimeout = connectionTimeout
}

func (c *Config) ConnectionReadTimeout() int {
	return c.connectionReadTimeout
}

func (c *Config) SetConnectionReadTimeout(connectionReadTimeout int) {
	c.connectionReadTimeout = connectionReadTimeout
}

func (c *Config) RootDir() string {
	return c.rootDir
}

func (c *Config) SetRootDir(rootDir string) {
	c.rootDir = rootDir
}

func (c *Config) StoreDir() string {
	return c.storeDir
}

func (c *Config) SetStoreDir(storeDir string) {
	c.storeDir = storeDir
}

func (c *Config) StoreBackupDir() string {
	return c.storeBackupDir
}

func (c *Config) SetStoreBackupDir(storeBackupDir string) {
	c.storeBackupDir = storeBackupDir
}

func (c *Config) WhiteListDir() string {
	return c.whiteListDir
}

func (c *Config) SetWhiteListDir(whiteListDir string) {
	c.whiteListDir = whiteListDir
}

func (c *Config) NetDBDir() string {
	return c.netDBDir
}

func (c *Config) SetNetDBDir(netDBDir string) {
	c.netDBDir = netDBDir
}

func (c *Config) StoreMaxOpenFiles() int {
	return c.storeMaxOpenFiles
}

func (c *Config) SetStoreMaxOpenFiles(storeMaxOpenFiles int) {
	c.storeMaxOpenFiles = storeMaxOpenFiles
}

func (c *Config) StoreMaxThreads() int {
	return c.storeMaxThreads
}

func (c *Config) SetStoreMaxThreads(storeMaxThreads int) {
	c.storeMaxThreads = storeMaxThreads
}

func (c *Config) StoreFromBackup() bool {
	return c.storeFromBackup
}

func (c *Config) SetStoreFromBackup(storeFromBackup bool) {
	c.storeFromBackup = storeFromBackup
}

func (c *Config) OriginStoreDir() string {
	return c.originStoreDir
}

func (c *Config) SetOriginStoreDir(originStoreDir string) {
	c.originStoreDir = originStoreDir
}

func (c *Config) WhitelistUrl() string {
	return c.whitelistUrl
}

func (c *Config) SetWhitelistUrl(whitelistUrl string) {
	c.whitelistUrl = whitelistUrl
}

func (c *Config) EnableRefresh() bool {
	return c.enableRefresh
}

func (c *Config) SetEnableRefresh(enableRefresh bool) {
	c.enableRefresh = enableRefresh
}

func (c *Config) DnetKeyFile() string {
	return c.dnetKeyFile
}

func (c *Config) SetDnetKeyFile(dnetKeyFile string) {
	c.dnetKeyFile = dnetKeyFile
}

func (c *Config) WalletKeyFile() string {
	return c.walletKeyFile
}

func (c *Config) SetWalletKeyFile(walletKeyFile string) {
	c.walletKeyFile = walletKeyFile
}

func (c *Config) Ttl() int {
	return c.ttl
}

func (c *Config) SetTtl(ttl int) {
	c.ttl = ttl
}

func (c *Config) DnetKeyBytes() [2048]byte {
	return c.dnetKeyBytes
}

func (c *Config) SetDnetKeyBytes(dnetKeyBytes [2048]byte) {
	c.dnetKeyBytes = dnetKeyBytes
}

func (c *Config) WhiteIPList() []string {
	return c.whiteIPList
}

func (c *Config) SetWhiteIPList(whiteIPList []string) {
	c.whiteIPList = whiteIPList
}

func (c *Config) Libp2pPort() int {
	return c.libp2pPort
}

func (c *Config) SetLibp2pPort(libp2pPort int) {
	c.libp2pPort = libp2pPort
}

func (c *Config) Libp2pPrivkey() string {
	return c.libp2pPrivkey
}

func (c *Config) SetLibp2pPrivkey(libp2pPrivkey string) {
	c.libp2pPrivkey = libp2pPrivkey
}

func (c *Config) IsBootNode() bool {
	return c.isBootNode
}

func (c *Config) SetIsBootNode(isBootNode bool) {
	c.isBootNode = isBootNode
}

func (c *Config) BootNodes() []string {
	return c.bootNodes
}

func (c *Config) SetBootNodes(bootNodes []string) {
	c.bootNodes = bootNodes
}

func (c *Config) WalletFilePath() string {
	return c.walletFilePath
}

func (c *Config) SetWalletFilePath(walletFilePath string) {
	c.walletFilePath = walletFilePath
}

func (c *Config) XdagEra() uint64 {
	return c.xdagEra
}

func (c *Config) SetXdagEra(xdagEra uint64) {
	c.xdagEra = xdagEra
}

func (c *Config) XdagFieldHeader() common.FieldType {
	return c.xdagFieldHeader
}

func (c *Config) SetXdagFieldHeader(xdagFieldHeader common.FieldType) {
	c.xdagFieldHeader = xdagFieldHeader
}

func (c *Config) MainStartAmount() uint64 {
	return c.mainStartAmount
}

func (c *Config) SetMainStartAmount(mainStartAmount uint64) {
	c.mainStartAmount = mainStartAmount
}

func (c *Config) ApolloForkHeight() uint64 {
	return c.apolloForkHeight
}

func (c *Config) SetApolloForkHeight(apolloForkHeight uint64) {
	c.apolloForkHeight = apolloForkHeight
}

func (c *Config) ApolloForkAmount() uint64 {
	return c.apolloForkAmount
}

func (c *Config) SetApolloForkAmount(apolloForkAmount uint64) {
	c.apolloForkAmount = apolloForkAmount
}

func (c *Config) RpcEnabled() bool {
	return c.rpcEnabled
}

func (c *Config) SetRpcEnabled(rpcEnabled bool) {
	c.rpcEnabled = rpcEnabled
}

func (c *Config) RpcHost() string {
	return c.rpcHost
}

func (c *Config) SetRpcHost(rpcHost string) {
	c.rpcHost = rpcHost
}

func (c *Config) RpcPortHttp() int {
	return c.rpcPortHttp
}

func (c *Config) SetRpcPortHttp(rpcPortHttp int) {
	c.rpcPortHttp = rpcPortHttp
}

func (c *Config) RpcPortWs() int {
	return c.rpcPortWs
}

func (c *Config) SetRpcPortWs(rpcPortWs int) {
	c.rpcPortWs = rpcPortWs
}

func (c *Config) SnapshotEnabled() bool {
	return c.snapshotEnabled
}

func (c *Config) SetSnapshotEnabled(snapshotEnabled bool) {
	c.snapshotEnabled = snapshotEnabled
}

func (c *Config) SnapshotHeight() uint64 {
	return c.snapshotHeight
}

func (c *Config) SetSnapshotHeight(snapshotHeight uint64) {
	c.snapshotHeight = snapshotHeight
}

func (c *Config) SnapshotTime() uint64 {
	return c.snapshotTime
}

func (c *Config) SetSnapshotTime(snapshotTime uint64) {
	c.snapshotTime = snapshotTime
}

func (c *Config) IsSnapshotJ() bool {
	return c.isSnapshotJ
}

func (c *Config) SetIsSnapshotJ(isSnapshotJ bool) {
	c.isSnapshotJ = isSnapshotJ
}
