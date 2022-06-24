package evm

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

// Todo - remove all hardcoded values in the next iteration.
// The Contract addresses are the result of the deploying a smart contract from hardcoded owners.
// The "owners" are keys which are the de-facto "admins" of those erc20s and are able to transfer or mint tokens.
// The contracts and addresses cannot be random for now, because there is hardcoded logic in the core
// to generate synthetic "transfer" transactions for each erc20 deposit on ethereum
// and these transactions need to be signed. Which means the platform needs to "own" ERC20s.

// ERC20 - the supported ERC20 tokens. A list of made-up tokens used for testing.
// Todo - this will be removed together will all the keys and addresses.
type ERC20 int

const (
	BTC ERC20 = iota
	ETH
)

var WBtcOwner, _ = crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")

// WBtcContract X- address of the deployed "btc" erc20
var WBtcContract = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))

var WEthOnwer, _ = crypto.HexToECDSA("4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8")

// WEthContract - address of the deployed "eth" erc20
var WEthContract = common.BytesToAddress(common.Hex2Bytes("9802F661d17c65527D7ABB59DAAD5439cb125a67"))

// BridgeAddress - address of the virtual bridge
var BridgeAddress = common.BytesToAddress(common.Hex2Bytes("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))

// ERC20Mapping - maps an L1 Erc20 to an L2 Erc20 address
type ERC20Mapping struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *common.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *common.Address
}

// Bridge encapsulates all logic around processing the interactions with an L1
type Bridge struct {
	SupportedTokens map[ERC20]*ERC20Mapping
	// BridgeAddress The address the bridge on the L2
	BridgeAddress common.Address
}

func NewBridge(obscuroChainID int64, btcAddress *common.Address, ethAddress *common.Address) *Bridge {
	tokens := make(map[ERC20]*ERC20Mapping, 0)

	tokens[BTC] = &ERC20Mapping{
		Name:      BTC,
		L1Address: btcAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), WBtcOwner),
		L2Address: &WBtcContract,
	}

	tokens[ETH] = &ERC20Mapping{
		Name:      ETH,
		L1Address: ethAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), WEthOnwer),
		L2Address: &WEthContract,
	}

	return &Bridge{
		SupportedTokens: tokens,
		BridgeAddress:   BridgeAddress,
	}
}

func (b *Bridge) IsWithdrawal(address common.Address) bool {
	return bytes.Equal(address.Bytes(), b.BridgeAddress.Bytes())
}

// L1Address - returns the L1 address of a token based on the mapping
func (b *Bridge) L1Address(l2Address *common.Address) *common.Address {
	if l2Address == nil {
		return nil
	}
	for _, t := range b.SupportedTokens {
		if bytes.Equal(l2Address.Bytes(), t.L2Address.Bytes()) {
			return t.L1Address
		}
	}
	return nil
}

// GetMapping - finds the maping based on the address that was called in an L1 transaction
func (b *Bridge) GetMapping(l1ContractAddress *common.Address) *ERC20Mapping {
	for _, t := range b.SupportedTokens {
		if bytes.Equal(t.L1Address.Bytes(), l1ContractAddress.Bytes()) {
			return t
		}
	}
	return nil
}

// Todo - move here the methods to extract deposits and post process withdrawal
