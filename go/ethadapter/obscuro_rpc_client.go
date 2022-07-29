package ethadapter

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/viewkey"
)

// obscuroWalletRPCClient implements the EthClient interface, it's bound to a single wallet (viewing key and address)
//
// Note: Why does this exist (why use a geth-focused interface for our obscuro interactions)?
//  - We have code (like deploying contracts) that we want to run against both L1 and L2, useful to be able to treat them the same
//	- Maintaining this interface helps ensure we remain closely compatible with eth usability
//
type obscuroWalletRPCClient struct {
	client rpcclientlib.Client // the underlying obscuro client - the viewing key and address are stored in here todo: although maybe they shouldn't be...?
	wallet wallet.Wallet
}

// NewObscuroRPCClient creates an obscuro RPC client for a given wallet, it will create and register a viewing key for the wallet as part of this setup
func NewObscuroRPCClient(ipaddress string, port uint, wallet wallet.Wallet) (EthClient, error) {
	client := rpcclientlib.NewClient(fmt.Sprintf("%s:%d", ipaddress, port))
	viewkey.GenerateAndRegisterViewingKey(client, wallet)
	return &obscuroWalletRPCClient{
		client: client,
		wallet: wallet,
	}, nil
}

func (c *obscuroWalletRPCClient) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) BlockByNumber(n *big.Int) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) SendTransaction(signedTx *types.Transaction) error {
	err := c.client.Call(nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
	if err != nil {
		return err
	}
	return nil
}

func (c *obscuroWalletRPCClient) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	var r types.Receipt
	err := c.client.Call(&r, rpcclientlib.RPCGetTxReceipt, hash)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *obscuroWalletRPCClient) Nonce(address gethcommon.Address) (uint64, error) {
	var result uint64
	err := c.client.Call(&result, rpcclientlib.RPCNonce, address)
	return result, err
}

func (c *obscuroWalletRPCClient) Info() Info {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) FetchHeadBlock() *types.Block {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) IsBlockAncestor(block *types.Block, proof common.L1RootHash) bool {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) RPCBlockchainFeed() []*types.Block {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) BlockListener() (chan *types.Header, ethereum.Subscription) {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (c *obscuroWalletRPCClient) Stop() {
	c.client.Stop()
}

func (c *obscuroWalletRPCClient) EthClient() *ethclient.Client {
	// TODO implement me
	panic("implement me")
}

// Formats a transaction for sending to the enclave
func encodeTx(tx *common.L2Tx) string {
	txBinary, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// We convert the transaction binary to the form expected for sending transactions via RPC.
	txBinaryHex := gethcommon.Bytes2Hex(txBinary)

	return "0x" + txBinaryHex
}
