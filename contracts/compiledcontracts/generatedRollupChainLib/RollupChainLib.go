// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedRollupChainLib

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RollupChainMetaData contains all meta data concerning the RollupChain contract.
var RollupChainMetaData = &bind.MetaData{
	ABI: "[]",
	Sigs: map[string]string{
		"f0a79d90": "AppendRollup(RollupChain.List storage,uint256,RollupChain.Rollup)",
		"5c84e72d": "GetParentRollup(RollupChain.List storage,RollupChain.RollupElement)",
		"33ce6cfc": "GetRollupByHash(RollupChain.List storage,bytes32)",
		"1590673a": "GetRollupByID(RollupChain.List storage,uint256)",
		"cfbc0dda": "HasSecondCousinFork(RollupChain.List storage)",
		"9770ebe8": "Initialize(RollupChain.List storage,RollupChain.Rollup)",
	},
	Bin: "0x6107f761003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061006c5760003560e01c80631590673a1461007157806333ce6cfc146100be5780635c84e72d146100d15780639770ebe8146100e4578063cfbc0dda14610106578063f0a79d9014610129575b600080fd5b61008461007f3660046105f8565b610149565b60408051921515835281516020808501919091528083015184830152910151805160608401520151608082015260a0015b60405180910390f35b6100846100cc3660046105f8565b6101a0565b6100846100df366004610677565b6101d1565b8180156100f057600080fd5b506101046100ff366004610706565b6101e9565b005b610119610114366004610733565b6102bf565b60405190151581526020016100b5565b81801561013557600080fd5b5061010461014436600461074c565b610472565b60006101536105c8565b505060009081526020918252604090819020815160608101835281548152600182015481850152825180840184526002830154815260039092015493820193909352908201528051151591565b60006101aa6105c8565b60008381526001850160205260409020546101c6908590610149565b915091509250929050565b60006101db6105c8565b6101c6848460200151610149565b600682015460ff16156102435760405162461bcd60e51b815260206004820152601b60248201527f63616e6e6f7420626520696e697469616c697a656420616761696e000000000060448201526064015b60405180910390fd5b60068201805460ff191660019081179091556040805160608101825282815260006020808301828152838501878152868452888352858420945185559051848701555180516002808601919091559082015160039094019390935560048701859055600587019290925593810151845293820190935291902055565b6000806102cb83610573565b90506000806102da85846101d1565b91509150816103175760405162461bcd60e51b81526020600482015260096024820152681b9bc81c185c995b9d60ba1b604482015260640161023a565b60008061032487846101d1565b91509150816103675760405162461bcd60e51b815260206004820152600f60248201526e1b9bc819dc985b99081c185c995b9d608a1b604482015260640161023a565b805160009081526002880160209081526040808320805482518185028101850190935280835291929091908301828280156103c157602002820191906000526020600020905b8154815260200190600101908083116103ad575b5050505050905060005b8151811015610463576000806103fa8b8585815181106103ed576103ed610782565b6020026020010151610149565b9150915081610414575060009a9950505050505050505050565b865181511415610425575050610451565b8051600090815260028c0160205260409020541561044e575060019a9950505050505050505050565b50505b8061045b81610798565b9150506103cb565b50600098975050505050505050565b600583018054908190600061048683610798565b91905055506000806104988686610149565b91509150816104dc5760405162461bcd60e51b815260206004820152601060248201526f1c185c995b9d081b9bdd08199bdd5b9960821b604482015260640161023a565b60408051606081018252848152602080820188815282840188815260008881528b84528581209451855591516001808601919091559051805160028087019190915590840151600390950194909455898252928a018252838120805480850182559082528282200187905587820151815291890190522083905560048601548151141561056b57600486018390555b505050505050565b61057b6105c8565b506004810154600090815260209182526040908190208151606081018352815481526001820154818501528251808401845260028301548152600390920154938201939093529082015290565b60408051606081018252600080825260208083018290528351808501855282815290810191909152909182015290565b6000806040838503121561060b57600080fd5b50508035926020909101359150565b60006040828403121561062c57600080fd5b6040516040810181811067ffffffffffffffff8211171561065d57634e487b7160e01b600052604160045260246000fd5b604052823581526020928301359281019290925250919050565b60008082840360a081121561068b57600080fd5b833592506080601f19820112156106a157600080fd5b506040516060810181811067ffffffffffffffff821117156106d357634e487b7160e01b600052604160045260246000fd5b806040525060208401358152604084013560208201526106f6856060860161061a565b6040820152809150509250929050565b6000806060838503121561071957600080fd5b8235915061072a846020850161061a565b90509250929050565b60006020828403121561074557600080fd5b5035919050565b60008060006080848603121561076157600080fd5b8335925060208401359150610779856040860161061a565b90509250925092565b634e487b7160e01b600052603260045260246000fd5b60006000198214156107ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220d6022b65ab5fa55482e28e79903bcf3e8adb93810ee440236c36709de2c9883364736f6c634300080c0033",
}

// RollupChainABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupChainMetaData.ABI instead.
var RollupChainABI = RollupChainMetaData.ABI

// Deprecated: Use RollupChainMetaData.Sigs instead.
// RollupChainFuncSigs maps the 4-byte function signature to its string representation.
var RollupChainFuncSigs = RollupChainMetaData.Sigs

// RollupChainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RollupChainMetaData.Bin instead.
var RollupChainBin = RollupChainMetaData.Bin

// DeployRollupChain deploys a new Ethereum contract, binding an instance of RollupChain to it.
func DeployRollupChain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RollupChain, error) {
	parsed, err := RollupChainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RollupChainBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RollupChain{RollupChainCaller: RollupChainCaller{contract: contract}, RollupChainTransactor: RollupChainTransactor{contract: contract}, RollupChainFilterer: RollupChainFilterer{contract: contract}}, nil
}

// RollupChain is an auto generated Go binding around an Ethereum contract.
type RollupChain struct {
	RollupChainCaller     // Read-only binding to the contract
	RollupChainTransactor // Write-only binding to the contract
	RollupChainFilterer   // Log filterer for contract events
}

// RollupChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupChainSession struct {
	Contract     *RollupChain      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupChainCallerSession struct {
	Contract *RollupChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RollupChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupChainTransactorSession struct {
	Contract     *RollupChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RollupChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupChainRaw struct {
	Contract *RollupChain // Generic contract binding to access the raw methods on
}

// RollupChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupChainCallerRaw struct {
	Contract *RollupChainCaller // Generic read-only contract binding to access the raw methods on
}

// RollupChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupChainTransactorRaw struct {
	Contract *RollupChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupChain creates a new instance of RollupChain, bound to a specific deployed contract.
func NewRollupChain(address common.Address, backend bind.ContractBackend) (*RollupChain, error) {
	contract, err := bindRollupChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollupChain{RollupChainCaller: RollupChainCaller{contract: contract}, RollupChainTransactor: RollupChainTransactor{contract: contract}, RollupChainFilterer: RollupChainFilterer{contract: contract}}, nil
}

// NewRollupChainCaller creates a new read-only instance of RollupChain, bound to a specific deployed contract.
func NewRollupChainCaller(address common.Address, caller bind.ContractCaller) (*RollupChainCaller, error) {
	contract, err := bindRollupChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupChainCaller{contract: contract}, nil
}

// NewRollupChainTransactor creates a new write-only instance of RollupChain, bound to a specific deployed contract.
func NewRollupChainTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupChainTransactor, error) {
	contract, err := bindRollupChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupChainTransactor{contract: contract}, nil
}

// NewRollupChainFilterer creates a new log filterer instance of RollupChain, bound to a specific deployed contract.
func NewRollupChainFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupChainFilterer, error) {
	contract, err := bindRollupChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupChainFilterer{contract: contract}, nil
}

// bindRollupChain binds a generic wrapper to an already deployed contract.
func bindRollupChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RollupChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupChain *RollupChainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupChain.Contract.RollupChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupChain *RollupChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChain.Contract.RollupChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupChain *RollupChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupChain.Contract.RollupChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupChain *RollupChainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupChain *RollupChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupChain *RollupChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupChain.Contract.contract.Transact(opts, method, params...)
}
