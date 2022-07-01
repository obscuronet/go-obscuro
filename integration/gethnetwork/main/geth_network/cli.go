package main

import (
	"flag"
	"strings"
)

const (
	// Flag names, defaults and usages.
	numNodesName  = "numNodes"
	numNodesUsage = "The number of nodes on the network"

	startPortName  = "startPort"
	startPortUsage = "The initial port to start allocating ports from"

	websocketStartPortName  = "websocketStartPort"
	websocketStartPortUsage = "The initial port to start allocating websocket ports from"

	prefundedAddrsName  = "prefundedAddrs"
	prefundedAddrsUsage = "The addresses to prefund as a comma-separated list"

	blockTimeSecsName  = "blockTimeSecs"
	blockTimeSecsUsage = "The block time in seconds"
)

type gethConfig struct {
	numNodes           int
	startPort          int
	websocketStartPort int
	prefundedAddrs     []string
	blockTimeSecs      int
}

func defaultHostConfig() gethConfig {
	return gethConfig{
		numNodes:           1,
		startPort:          12000,
		websocketStartPort: 12100,
		prefundedAddrs:     []string{},
		blockTimeSecs:      6,
	}
}

func parseCLIArgs() gethConfig {
	defaultConfig := defaultHostConfig()

	numNodes := flag.Int(numNodesName, defaultConfig.numNodes, numNodesUsage)
	startPort := flag.Int(startPortName, defaultConfig.startPort, startPortUsage)
	websocketStartPort := flag.Int(websocketStartPortName, defaultConfig.websocketStartPort, websocketStartPortUsage)
	prefundedAddrs := flag.String(prefundedAddrsName, "", prefundedAddrsUsage)
	blockTimeSecs := flag.Int(blockTimeSecsName, defaultConfig.blockTimeSecs, blockTimeSecsUsage)

	flag.Parse()

	parsedPrefundedAddrs := strings.Split(*prefundedAddrs, ",")
	if *prefundedAddrs == "" {
		// We handle the special case of an empty list.
		parsedPrefundedAddrs = []string{}
	}

	return gethConfig{
		numNodes:           *numNodes,
		startPort:          *startPort,
		websocketStartPort: *websocketStartPort,
		prefundedAddrs:     parsedPrefundedAddrs,
		blockTimeSecs:      *blockTimeSecs,
	}
}
