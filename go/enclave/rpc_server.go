package enclave

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/naoina/toml"
	"github.com/obscuronet/obscuro-playground/go/config"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/common/rpc"
	"github.com/obscuronet/obscuro-playground/go/common/rpc/generated"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/mgmtcontractlib"
	"google.golang.org/grpc"
)

// Receives RPC calls to the enclave process and relays them to the enclave.Enclave.
type server struct {
	generated.UnimplementedEnclaveProtoServer
	enclave     common.Enclave
	rpcServer   *grpc.Server
	nodeShortID uint64
}

// StartServer starts a server on the given port on a separate thread. It creates an enclave.Enclave for the provided nodeID,
// and uses it to respond to incoming RPC messages from the host.
func StartServer(
	enclaveConfig config.EnclaveConfig,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	collector StatsCollector,
) (func(), error) {
	lis, err := net.Listen("tcp", enclaveConfig.Address)
	if err != nil {
		return nil, fmt.Errorf("enclave RPC server could not listen on port: %w", err)
	}

	enclaveServer := server{
		enclave:     NewEnclave(enclaveConfig, mgmtContractLib, erc20ContractLib, collector),
		rpcServer:   grpc.NewServer(),
		nodeShortID: common.ShortAddress(enclaveConfig.HostID),
	}
	generated.RegisterEnclaveProtoServer(enclaveServer.rpcServer, &enclaveServer)

	go func(lis net.Listener) {
		common.LogWithID(enclaveServer.nodeShortID, "Enclave server listening on address %s.", enclaveConfig.Address)
		err = enclaveServer.rpcServer.Serve(lis)
		if err != nil {
			common.LogWithID(enclaveServer.nodeShortID, "enclave RPC server could not serve: %s", err)
		}
	}(lis)

	closeHandle := func() {
		go enclaveServer.Stop(context.Background(), nil) //nolint:errcheck
	}

	tomlConfig, err := toml.Marshal(enclaveConfig)
	if err != nil {
		panic("could not print enclave config")
	}
	log.Info("Enclave service started with following config:\n%s", tomlConfig)

	return closeHandle, nil
}

// IsReady returns a nil error to indicate that the server is ready.
func (s *server) IsReady(context.Context, *generated.IsReadyRequest) (*generated.IsReadyResponse, error) {
	errStr := ""
	if err := s.enclave.IsReady(); err != nil {
		errStr = err.Error()
	}
	return &generated.IsReadyResponse{Error: errStr}, nil
}

func (s *server) Attestation(context.Context, *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation := s.enclave.Attestation()
	msg := rpc.ToAttestationReportMsg(attestation)
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *server) GenerateSecret(context.Context, *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret := s.enclave.GenerateSecret()
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) ShareSecret(_ context.Context, request *generated.FetchSecretRequest) (*generated.ShareSecretResponse, error) {
	attestationReport := rpc.FromAttestationReportMsg(request.AttestationReportMsg)
	secret, err := s.enclave.ShareSecret(attestationReport)
	return &generated.ShareSecretResponse{EncryptedSharedEnclaveSecret: secret}, err
}

func (s *server) InitEnclave(_ context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	errStr := ""
	if err := s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret); err != nil {
		errStr = err.Error()
	}
	return &generated.InitEnclaveResponse{Error: errStr}, nil
}

func (s *server) IsInitialised(context.Context, *generated.IsInitialisedRequest) (*generated.IsInitialisedResponse, error) {
	isInitialised := s.enclave.IsInitialised()
	return &generated.IsInitialisedResponse{IsInitialised: isInitialised}, nil
}

func (s *server) ProduceGenesis(_ context.Context, request *generated.ProduceGenesisRequest) (*generated.ProduceGenesisResponse, error) {
	genesisRollup := s.enclave.ProduceGenesis(gethcommon.BytesToHash(request.GetBlockHash()))
	blockSubmissionResponse := rpc.ToBlockSubmissionResponseMsg(genesisRollup)
	return &generated.ProduceGenesisResponse{BlockSubmissionResponse: &blockSubmissionResponse}, nil
}

func (s *server) IngestBlocks(_ context.Context, request *generated.IngestBlocksRequest) (*generated.IngestBlocksResponse, error) {
	blocks := make([]*types.Block, 0)
	for _, encodedBlock := range request.EncodedBlocks {
		bl := s.decodeBlock(encodedBlock)
		blocks = append(blocks, &bl)
	}

	r := s.enclave.IngestBlocks(blocks)
	blockSubmissionResponses := make([]*generated.BlockSubmissionResponseMsg, len(r))
	for i, response := range r {
		b := rpc.ToBlockSubmissionResponseMsg(response)
		blockSubmissionResponses[i] = &b
	}
	return &generated.IngestBlocksResponse{
		BlockSubmissionResponses: blockSubmissionResponses,
	}, nil
}

func (s *server) Start(_ context.Context, request *generated.StartRequest) (*generated.StartResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	s.enclave.Start(bl)
	return &generated.StartResponse{}, nil
}

func (s *server) SubmitBlock(_ context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	blockSubmissionResponse := s.enclave.SubmitBlock(bl)

	msg := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
}

func (s *server) SubmitRollup(_ context.Context, request *generated.SubmitRollupRequest) (*generated.SubmitRollupResponse, error) {
	extRollup := rpc.FromExtRollupMsg(request.ExtRollup)
	s.enclave.SubmitRollup(extRollup)
	return &generated.SubmitRollupResponse{}, nil
}

func (s *server) SubmitTx(_ context.Context, request *generated.SubmitTxRequest) (*generated.SubmitTxResponse, error) {
	err := s.enclave.SubmitTx(request.EncryptedTx)
	return &generated.SubmitTxResponse{}, err
}

func (s *server) ExecuteOffChainTransaction(_ context.Context, request *generated.OffChainRequest) (*generated.OffChainResponse, error) {
	result, err := s.enclave.ExecuteOffChainTransaction(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.OffChainResponse{Result: result}, nil
}

func (s *server) GetTransactionCount(_ context.Context, request *generated.GetTransactionCountRequest) (*generated.GetTransactionCountResponse, error) {
	result, err := s.enclave.GetTransactionCount(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionCountResponse{Result: result}, nil
}

func (s *server) RoundWinner(_ context.Context, request *generated.RoundWinnerRequest) (*generated.RoundWinnerResponse, error) {
	extRollup, winner, err := s.enclave.RoundWinner(gethcommon.BytesToHash(request.Parent))
	if err != nil {
		return nil, err
	}
	extRollupMsg := rpc.ToExtRollupMsg(&extRollup)
	return &generated.RoundWinnerResponse{Winner: winner, ExtRollup: &extRollupMsg}, nil
}

func (s *server) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	defer s.rpcServer.GracefulStop()
	err := s.enclave.Stop()
	return &generated.StopResponse{}, err
}

func (s *server) GetTransaction(_ context.Context, request *generated.GetTransactionRequest) (*generated.GetTransactionResponse, error) {
	tx := s.enclave.GetTransaction(gethcommon.BytesToHash(request.TxHash))
	if tx == nil {
		return &generated.GetTransactionResponse{Known: false, EncodedTransaction: []byte{}}, nil
	}

	var buffer bytes.Buffer
	if err := tx.EncodeRLP(&buffer); err != nil {
		common.LogWithID(s.nodeShortID, "failed to decode transaction sent to enclave: %v", err)
	}
	return &generated.GetTransactionResponse{Known: true, EncodedTransaction: buffer.Bytes()}, nil
}

func (s *server) GetTransactionReceipt(_ context.Context, request *generated.GetTransactionReceiptRequest) (*generated.GetTransactionReceiptResponse, error) {
	encryptedTxReceipt, err := s.enclave.GetTransactionReceipt(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionReceiptResponse{EncryptedTxReceipt: encryptedTxReceipt}, nil
}

func (s *server) GetRollup(_ context.Context, request *generated.GetRollupRequest) (*generated.GetRollupResponse, error) {
	extRollup := s.enclave.GetRollup(gethcommon.BytesToHash(request.RollupHash))
	if extRollup == nil {
		return &generated.GetRollupResponse{Known: false, ExtRollup: nil}, nil
	}

	extRollupMsg := rpc.ToExtRollupMsg(extRollup)
	return &generated.GetRollupResponse{Known: true, ExtRollup: &extRollupMsg}, nil
}

func (s *server) AddViewingKey(_ context.Context, request *generated.AddViewingKeyRequest) (*generated.AddViewingKeyResponse, error) {
	err := s.enclave.AddViewingKey(request.ViewingKey, request.Signature)
	if err != nil {
		return nil, err
	}
	return &generated.AddViewingKeyResponse{}, nil
}

func (s *server) GetBalance(_ context.Context, request *generated.GetBalanceRequest) (*generated.GetBalanceResponse, error) {
	encryptedBalance, err := s.enclave.GetBalance(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetBalanceResponse{EncryptedBalance: encryptedBalance}, nil
}

func (s *server) decodeBlock(encodedBlock []byte) types.Block {
	block := types.Block{}
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		common.LogWithID(s.nodeShortID, "failed to decode block sent to enclave: %v", err)
	}
	return block
}
