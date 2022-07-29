package contractdeployer

// Flag names, defaults and usages.
const (
	nodeHostName  = "nodeHost"
	nodeHostUsage = "The host on which to connect the RPC client"

	nodePortName  = "nodePort"
	nodePortUsage = "The port on which to connect the RPC client"

	isL1DeploymentName  = "l1Deployment"
	isL1DeploymentUsage = "Set to true for L1 contract to deployment (false is deployment to obscuro network, the L2)"

	contractNameName  = "contractName"
	contractNameUsage = "The name of the contract bytecode to be deploy (e.g. `MGMT` or `ERC20`)"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the node account"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the chain (defaults to 777 for L2 deployment and 1337 for L1)"
)
