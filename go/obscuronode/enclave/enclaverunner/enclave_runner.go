package enclaverunner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
	"math/big"
	"os"
)

const logPath = "enclave_logs.txt"

// RunEnclave runs an Obscuro enclave as a standalone process.
func RunEnclave(config EnclaveConfig) {
	setLogs(*config.WriteToLogs)

	nodeAddress := common.BigToAddress(big.NewInt(*config.NodeID))
	// TODO - For now, genesisJSON is nil. This means that incoming L1 blocks are not validated by the enclave. In the
	//  future, we should allow the genesisJSON to be passed in somehow, with a default of the default genesis.
	if err := enclave.StartServer(*config.Address, nodeAddress, ethereum_mock.NewMockTxHandler(), false, nil, nil); err != nil {
		panic(err)
	}

	select {}
}

// Sets the log file, defaulting to stdout if writeToLogs is false.
func setLogs(writeToLogs bool) {
	var logFile *os.File
	var err error

	if writeToLogs {
		logFile, err = os.Create(logPath)
		if err != nil {
			panic(err)
		}
	} else {
		logFile = os.Stdout
	}

	log.SetLog(logFile)
}
