package hostrunner

// Flag names, defaults and usages.
const (
	configName  = "config"
	configUsage = "The path to the node's config file. Overrides all other flags"

	nodeIDName  = "id"
	nodeIDUsage = "The 20 bytes of the node's address"

	isGenesisName  = "isGenesis"
	isGenesisUsage = "Whether the node is the first node to join the network"

	gossipRoundNanosName  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"

	clientRPCAddressName  = "clientRPCAddress"
	clientRPCAddressUsage = "The address on which to listen for client application RPC requests"

	clientRPCTimeoutSecsName  = "clientRPCTimeoutSecs"
	clientRPCTimeoutSecsUsage = "The timeout for client <-> host RPC communication"

	enclaveRPCAddressName  = "enclaveRPCAddress"
	enclaveRPCAddressUsage = "The address to use to connect to the Obscuro enclave service"

	enclaveRPCTimeoutSecsName  = "enclaveRPCTimeoutSecs"
	enclaveRPCTimeoutSecsUsage = "The timeout for host <-> enclave RPC communication"

	p2pAddressName  = "p2pAddress"
	p2pAddressUsage = "The P2P address for our node"

	peerP2PAddressesName = "peerP2PAddresses"
	peerP2PAddrsUsage    = "The P2P addresses of our peer nodes as a comma-separated list"

	l1NodeHostName  = "l1NodeHost"
	l1NodeHostUsage = "The host on which to connect to the Ethereum client"

	l1NodePortName  = "l1NodePort"
	l1NodePortUsage = "The port on which to connect to the Ethereum client"

	l1ConnectionTimeoutSecsName  = "l1ConnectionTimeoutSecs"
	l1ConnectionTimeoutSecsUsage = "The timeout for connecting to the Ethereum client"

	rollupContractAddrName  = "rollupContractAddress"
	rollupContractAddrUsage = "The management contract address on the L1"

	logPathName  = "logPath"
	logPathUsage = "The path to use for the host's log file"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	contractMgmtBlkHashName  = "contractMgmtBlkHash"
	contractMgmtBlkHashUsage = "The block hash where the contract management was deployed"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the L1 chain"
)
