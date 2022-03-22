package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"os"
)

const (
	nodeAddressFlag  = "nodeAddress"
	nodeAddressUsage = "The 20 bytes of the node's address"
	portFlag         = "port"
	portUsage        = "The port on which to serve the Obscuro enclave service"
	usage            = `CLI application for the ◠.bscuro enclave service.

Usage:

    <executable> [flags]

The flags are:

-%s   string   %s
-%s   int   %s`
)

func main() {
	if len(os.Args) == 1 {
		usageFmt := fmt.Sprintf(usage, nodeAddressFlag, nodeAddressUsage, portFlag, portUsage)
		fmt.Println(usageFmt)
		return
	}

	nodeAddressBytes, port := parseCLIArgs()

	nodeAddress := common.BytesToAddress([]byte(*nodeAddressBytes))
	if _, err := enclave.StartServer(*port, nodeAddress, nil); err != nil {
		panic(err)
	}

	fmt.Printf("Enclave server listening on port %d.\n", *port)
	select {}
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, *uint64) {
	var nodeAddressBytes = flag.String(nodeAddressFlag, "", nodeAddressUsage)
	var port = flag.Uint64(portFlag, 10000, portUsage)
	flag.Parse()

	return nodeAddressBytes, port
}
