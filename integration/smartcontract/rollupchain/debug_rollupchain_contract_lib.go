package rollupchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedRollupChainTestContract"
)

// debugRollupChainTestContractLib allows the direct use of the generatedManagementContract package
type debugRollupChainTestContractLib struct {
	genContract *generatedRollupChainTestContract.RollupChainTestContract
	address     *common.Address
}

// newdebugRollupChainContractLib creates an instance of the generated contract package and allows the use of the debugRollupChainTestContractLib properties
func newdebugRollupChainContractLib(address common.Address, client *ethclient.Client) *debugRollupChainTestContractLib {
	genContract, err := generatedRollupChainTestContract.NewRollupChainTestContract(address, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connect contract at addr: %s\n", address)

	return &debugRollupChainTestContractLib{
		genContract,
		&address,
	}
}

func (c *debugRollupChainTestContractLib) Address() *common.Address {
	return c.address
}
