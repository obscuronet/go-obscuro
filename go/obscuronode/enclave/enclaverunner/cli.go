package enclaverunner

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "A integer representing the 20 bytes of the node's address (default 1)"

	addressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	writeToLogsName  = "writeToLogs"
	writeToLogsUsage = "Whether to redirect the output to the log file."

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"
)

type EnclaveConfig struct {
	NodeID          int64
	Address         string
	WriteToLogs     bool
	ContractAddress string
}

func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		NodeID:          1,
		Address:         "localhost:11000",
		WriteToLogs:     false,
		ContractAddress: "",
	}
}

func ParseCLIArgs() EnclaveConfig {
	defaultConfig := DefaultEnclaveConfig()

	nodeID := flag.Int64(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	port := flag.String(addressName, defaultConfig.Address, addressUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	flag.Parse()

	return EnclaveConfig{NodeID: *nodeID, Address: *port, WriteToLogs: *writeToLogs, ContractAddress: *contractAddress}
}