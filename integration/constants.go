package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
const (
	StartPortGethNetworkTest          uint64 = 30000
	StartPortWalletExtensionTest      uint64 = 31000
	StartPortNodeRunnerTest           uint64 = 32000
	StartPortSimulationDocker                = 33000
	StartPortSimulationGethInMem             = 34000
	StartPortSimulationInMem                 = 35000
	StartPortSimulationAzureEnclave          = 36000
	StartPortSimulationFullNetwork           = 37000
	StartPortManagementContractTests         = 38000
	StartPortRollupChainContractTests        = 39000
)

const (
	EthereumChainID = 1337
	ObscuroChainID  = 777
)
