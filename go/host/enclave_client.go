package host

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"google.golang.org/grpc/connectivity"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EnclaveRPCClient implements enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type EnclaveRPCClient struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      config.HostConfig
	nodeShortID uint64
}

// TODO - Avoid panicking and return errors instead where appropriate.

func NewEnclaveRPCClient(config config.HostConfig) *EnclaveRPCClient {
	nodeShortID := common.ShortAddress(config.ID)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		common.PanicWithID(nodeShortID, "Failed to connect to enclave RPC service. Cause: %s", err)
	}

	// We wait for the RPC connection to be ready.
	currentTime := time.Now()
	deadline := currentTime.Add(30 * time.Second)
	currentState := connection.GetState()
	for currentState == connectivity.Idle || currentState == connectivity.Connecting || currentState == connectivity.TransientFailure {
		connection.Connect()
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
		currentState = connection.GetState()
	}

	if currentState != connectivity.Ready {
		common.PanicWithID(nodeShortID, "RPC connection failed to establish. Current state is %s", currentState)
	}

	return &EnclaveRPCClient{
		generated.NewEnclaveProtoClient(connection),
		connection,
		config,
		nodeShortID,
	}
}

func (c *EnclaveRPCClient) StopClient() error {
	return c.connection.Close()
}

func (c *EnclaveRPCClient) IsReady() error {
	if c.connection.GetState() != connectivity.Ready {
		return errors.New("RPC connection is not ready")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.IsReady(timeoutCtx, &generated.IsReadyRequest{})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}

func (c *EnclaveRPCClient) Attestation() *common.AttestationReport {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to retrieve attestation. Cause: %s", err)
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg)
}

func (c *EnclaveRPCClient) GenerateSecret() common.EncryptedSharedEnclaveSecret {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to generate secret. Cause: %s", err)
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveRPCClient) ShareSecret(report *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	attestationReportMsg := rpc.ToAttestationReportMsg(report)
	request := generated.FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.ShareSecret(timeoutCtx, &request)
	if err != nil {
		return nil, err
	}
	if response.GetError() != "" {
		return nil, errors.New(response.GetError())
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *EnclaveRPCClient) InitEnclave(secret common.EncryptedSharedEnclaveSecret) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}

func (c *EnclaveRPCClient) IsInitialised() bool {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &generated.IsInitialisedRequest{})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to establish enclave initialisation status. Cause: %s", err)
	}
	return response.IsInitialised
}

func (c *EnclaveRPCClient) ProduceGenesis(blkHash gethcommon.Hash) common.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()
	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{BlockHash: blkHash.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to produce genesis. Cause: %s", err)
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) IngestBlocks(blocks []*types.Block) []common.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock, err := common.EncodeBlock(block)
		if err != nil {
			common.PanicWithID(c.nodeShortID, "Failed to ingest blocks. Cause: %s", err)
		}
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	response, err := c.protoClient.IngestBlocks(timeoutCtx, &generated.IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to ingest blocks. Cause: %s", err)
	}
	responses := response.GetBlockSubmissionResponses()
	result := make([]common.BlockSubmissionResponse, len(responses))
	for i, r := range responses {
		result[i] = rpc.FromBlockSubmissionResponseMsg(r)
	}
	return result
}

func (c *EnclaveRPCClient) Start(block types.Block) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to encode block. Cause: %s", err)
	}
	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to start enclave. Cause: %s", err)
	}
}

func (c *EnclaveRPCClient) SubmitBlock(block types.Block) common.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to encode block. Cause: %s", err)
	}

	processTime := time.Now()
	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to submit block. Cause: %s", err)
	}
	log.Debug("Block %s processed by the enclave over RPC in %s", block.Hash().Hex(), time.Since(processTime))
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) SubmitRollup(rollup common.ExtRollup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	extRollupMsg := rpc.ToExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to submit rollup. Cause: %s", err)
	}
}

func (c *EnclaveRPCClient) SubmitTx(tx common.EncryptedTx) (common.EncryptedResponseSendRawTx, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return nil, err
	}
	return response.EncryptedHash, err
}

func (c *EnclaveRPCClient) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ExecuteOffChainTransaction(timeoutCtx, &generated.OffChainRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result, nil
}

func (c *EnclaveRPCClient) Nonce(address gethcommon.Address) uint64 {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Nonce(timeoutCtx, &generated.NonceRequest{Address: address.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to retrieve nonce: %s", err)
	}
	return response.Nonce
}

func (c *EnclaveRPCClient) RoundWinner(parent common.L2RootHash) (common.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to determine round winner. Cause: %s", err)
	}

	if response.Winner {
		return rpc.FromExtRollupMsg(response.ExtRollup), true, nil
	}
	return common.ExtRollup{}, false, nil
}

func (c *EnclaveRPCClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("failed to stop enclave: %w", err)
	}
	return nil
}

func (c *EnclaveRPCClient) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (common.EncryptedResponseGetTxByHash, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedTx, nil
}

func (c *EnclaveRPCClient) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (common.EncryptedResponseGetTxReceipt, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}

	return response.EncryptedTxReceipt, nil
}

func (c *EnclaveRPCClient) GetRollup(rollupHash common.L2RootHash) (*common.ExtRollup, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetRollup(timeoutCtx, &generated.GetRollupRequest{RollupHash: rollupHash.Bytes()})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rollup with hash %s. Cause: %w", rollupHash.Hex(), err)
	}

	extRollup := rpc.FromExtRollupMsg(response.ExtRollup)
	return &extRollup, nil
}

func (c *EnclaveRPCClient) GetRollupByHeight(rollupHeight int64) (*common.ExtRollup, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetRollupByHeight(timeoutCtx, &generated.GetRollupByHeightRequest{RollupHeight: rollupHeight})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rollup with height %d. Cause: %w", rollupHeight, err)
	}

	extRollup := rpc.FromExtRollupMsg(response.ExtRollup)
	return &extRollup, nil
}

func (c *EnclaveRPCClient) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.AddViewingKey(timeoutCtx, &generated.AddViewingKeyRequest{
		ViewingKey: viewingKeyBytes,
		Signature:  signature,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *EnclaveRPCClient) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedBalance, nil
}

func (c *EnclaveRPCClient) GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetCode(timeoutCtx, &generated.GetCodeRequest{
		Address:    address.Bytes(),
		RollupHash: rollupHash.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return resp.Code, nil
}

func (c *EnclaveRPCClient) StoreAttestation(report *common.AttestationReport) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	msg := rpc.ToAttestationReportMsg(report)
	resp, err := c.protoClient.StoreAttestation(timeoutCtx, &generated.StoreAttestationRequest{
		AttestationReportMsg: &msg,
	})
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}
	return nil
}
