package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName    = "nodeID"
	nodeIDDefault = 1
	nodeIDUsage   = "A integer representing the 20 bytes of the node's address (default 1)"

	addressName    = "address"
	addressDefault = ":11000"
	addressUsage   = "The address on which to serve the Obscuro enclave service"

	writeToLogsName    = "writeToLogs"
	writeToLogsDefault = false
	writeToLogsUsage   = "Whether to redirect the output to the log file."
)

type enclaveConfig struct {
	nodeID      *int64
	address     *string
	writeToLogs *bool
}

func parseCLIArgs() enclaveConfig {
	nodeID := flag.Int64(nodeIDName, nodeIDDefault, nodeIDUsage)
	port := flag.String(addressName, addressDefault, addressUsage)
	writeToLogs := flag.Bool(writeToLogsName, writeToLogsDefault, writeToLogsUsage)
	flag.Parse()

	return enclaveConfig{nodeID, port, writeToLogs}
}
