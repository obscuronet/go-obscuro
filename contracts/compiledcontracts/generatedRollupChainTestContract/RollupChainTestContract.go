// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedRollupChainTestContract

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

// RollupChainTestContractMetaData contains all meta data concerning the RollupChainTestContract contract.
var RollupChainTestContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AppendRollupTest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NoForkDetection\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RevertsNoDoubleInitTest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ScrollTreeTest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"f6da5eb5": "AppendRollupTest()",
		"01ac6e5b": "NoForkDetection()",
		"4762fa7c": "RevertsNoDoubleInitTest()",
		"f2b5185d": "ScrollTreeTest()",
	},
	Bin: "0x608060405234801561001057600080fd5b50611504806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806301ac6e5b146100515780634762fa7c1461005b578063f2b5185d14610063578063f6da5eb51461006b575b600080fd5b610059610073565b005b6100596105ae565b61005961068e565b610059610aff565b601573__$7204b1ba8a254ced74f31676d70e6726eb$__639770ebe882610098610fd2565b6040518363ffffffff1660e01b81526004016100b592919061114d565b60006040518083038186803b1580156100cd57600080fd5b505af41580156100e1573d6000803e3d6000fd5b5050604051630ac8339d60e11b815260048101849052600160248201526000925082915073__$7204b1ba8a254ced74f31676d70e6726eb$__90631590673a9060440160a060405180830381865af4158015610141573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016591906111be565b915091508161018f5760405162461bcd60e51b815260040161018690611270565b60405180910390fd5b60005b60058161ffff16101561032c5760008260400151602001516101b4600a611069565b60408051602081019390935282015260600160408051601f198184030181528282528051602091820120838301835286830151820151845290830181905285519151630f0a79d960e41b815290935073__$7204b1ba8a254ced74f31676d70e6726eb$__9163f0a79d909161024a918a918690600401928352602080840192909252805160408401520151606082015260800190565b60006040518083038186803b15801561026257600080fd5b505af4158015610276573d6000803e3d6000fd5b5050604051630cf39b3f60e21b8152600481018990526024810185905273__$7204b1ba8a254ced74f31676d70e6726eb$__92506333ce6cfc915060440160a060405180830381865af41580156102d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f591906111be565b9095509350846103175760405162461bcd60e51b8152600401610186906112b5565b50508080610324906112ff565b915050610192565b506000610338846110cd565b825181519192501461035c5760405162461bcd60e51b81526004016101869061132f565b8160400151602001518160400151602001511461038b5760405162461bcd60e51b815260040161018690611375565b60005b60058161ffff16101561044957604051635c84e72d60e01b815273__$7204b1ba8a254ced74f31676d70e6726eb$__90635c84e72d906103d490889086906004016113c7565b60a060405180830381865af41580156103f1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041591906111be565b9094509150836104375760405162461bcd60e51b815260040161018690611407565b80610441816112ff565b91505061038e565b50604051635c84e72d60e01b815273__$7204b1ba8a254ced74f31676d70e6726eb$__90635c84e72d9061048390879085906004016113c7565b60a060405180830381865af41580156104a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104c491906111be565b909350905082156104e75760405162461bcd60e51b81526004016101869061144b565b6040516367de06ed60e11b81526004810185905260009073__$7204b1ba8a254ced74f31676d70e6726eb$__9063cfbc0dda90602401602060405180830381865af415801561053a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061055e9190611491565b905080156105a75760405162461bcd60e51b81526020600482015260166024820152751cda1bdd5b19081b9bdd081a185d9948199bdc9ad95960521b6044820152606401610186565b5050505050565b600073__$7204b1ba8a254ced74f31676d70e6726eb$__639770ebe8826105d3610fd2565b6040518363ffffffff1660e01b81526004016105f092919061114d565b60006040518083038186803b15801561060857600080fd5b505af415801561061c573d6000803e3d6000fd5b505050508073__$7204b1ba8a254ced74f31676d70e6726eb$__639770ebe89091610645610fd2565b6040518363ffffffff1660e01b815260040161066292919061114d565b60006040518083038186803b15801561067a57600080fd5b505af41580156105a7573d6000803e3d6000fd5b600e73__$7204b1ba8a254ced74f31676d70e6726eb$__639770ebe8826106b3610fd2565b6040518363ffffffff1660e01b81526004016106d092919061114d565b60006040518083038186803b1580156106e857600080fd5b505af41580156106fc573d6000803e3d6000fd5b5050604051630ac8339d60e11b815260048101849052600160248201526000925082915073__$7204b1ba8a254ced74f31676d70e6726eb$__90631590673a9060440160a060405180830381865af415801561075c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078091906111be565b91509150816107a15760405162461bcd60e51b815260040161018690611270565b60005b60058161ffff16101561093e5760008260400151602001516107c6600a611069565b60408051602081019390935282015260600160408051601f198184030181528282528051602091820120838301835286830151820151845290830181905285519151630f0a79d960e41b815290935073__$7204b1ba8a254ced74f31676d70e6726eb$__9163f0a79d909161085c918a918690600401928352602080840192909252805160408401520151606082015260800190565b60006040518083038186803b15801561087457600080fd5b505af4158015610888573d6000803e3d6000fd5b5050604051630cf39b3f60e21b8152600481018990526024810185905273__$7204b1ba8a254ced74f31676d70e6726eb$__92506333ce6cfc915060440160a060405180830381865af41580156108e3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061090791906111be565b9095509350846109295760405162461bcd60e51b8152600401610186906112b5565b50508080610936906112ff565b9150506107a4565b50600061094a846110cd565b825181519192501461096e5760405162461bcd60e51b81526004016101869061132f565b8160400151602001518160400151602001511461099d5760405162461bcd60e51b815260040161018690611375565b60005b60058161ffff161015610a5b57604051635c84e72d60e01b815273__$7204b1ba8a254ced74f31676d70e6726eb$__90635c84e72d906109e690889086906004016113c7565b60a060405180830381865af4158015610a03573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a2791906111be565b909450915083610a495760405162461bcd60e51b815260040161018690611407565b80610a53816112ff565b9150506109a0565b50604051635c84e72d60e01b815273__$7204b1ba8a254ced74f31676d70e6726eb$__90635c84e72d90610a9590879085906004016113c7565b60a060405180830381865af4158015610ab2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ad691906111be565b90935090508215610af95760405162461bcd60e51b81526004016101869061144b565b50505050565b600773__$7204b1ba8a254ced74f31676d70e6726eb$__639770ebe882610b24610fd2565b6040518363ffffffff1660e01b8152600401610b4192919061114d565b60006040518083038186803b158015610b5957600080fd5b505af4158015610b6d573d6000803e3d6000fd5b5050604051630ac8339d60e11b815260048101849052600160248201526000925082915073__$7204b1ba8a254ced74f31676d70e6726eb$__90631590673a9060440160a060405180830381865af4158015610bcd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bf191906111be565b9150915081610c125760405162461bcd60e51b815260040161018690611270565b604080820151602001519051630cf39b3f60e21b815273__$7204b1ba8a254ced74f31676d70e6726eb$__916333ce6cfc91610c5b918791600401918252602082015260400190565b60a060405180830381865af4158015610c78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c9c91906111be565b909250905081610cfe5760405162461bcd60e51b815260206004820152602760248201527f476574526f6c6c75704279486173683a20666972737420726f6c6c7570206e6f6044820152661d08199bdd5b9960ca1b6064820152608401610186565b6000610d09846110cd565b604081810180516020908101518351808301919091526000818501528351603481830301815260548201808652815191840191909120609483018652935190920151825260740182815284519351630f0a79d960e41b8152600481018a9052602481019490945281516044850152516064840152929350919073__$7204b1ba8a254ced74f31676d70e6726eb$__9063f0a79d909060840160006040518083038186803b158015610db957600080fd5b505af4158015610dcd573d6000803e3d6000fd5b5050604051630ac8339d60e11b815260048101899052600260248201526000925082915073__$7204b1ba8a254ced74f31676d70e6726eb$__90631590673a9060440160a060405180830381865af4158015610e2d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5191906111be565b9150915081610ea25760405162461bcd60e51b815260206004820152601f60248201527f476574526f6c6c7570427949443a20726f6c6c7570206e6f7420666f756e64006044820152606401610186565b604051630cf39b3f60e21b8152600481018990526024810185905273__$7204b1ba8a254ced74f31676d70e6726eb$__906333ce6cfc9060440160a060405180830381865af4158015610ef9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f1d91906111be565b909250905081610f795760405162461bcd60e51b815260206004820152602160248201527f476574526f6c6c75704279486173683a20726f6c6c7570206e6f7420666f756e6044820152601960fa1b6064820152608401610186565b826020015181604001516020015114610fc85760405162461bcd60e51b81526020600482015260116024820152700d0c2e6d0cae640c8dedce840dac2e8c6d607b1b6044820152606401610186565b5050505050505050565b60408051808201909152600080825260208201526040805144602082015242918101919091524360608201526000906080016040516020818303038152906040528051906020012090506040518060400160405280828152602001611038612710611069565b60405160200161104a91815260200190565b6040516020818303038152906040528051906020012081525091505090565b6000814244336040516020016110a493929190928352602083019190915260601b6bffffffffffffffffffffffff1916604082015260540190565b6040516020818303038152906040528051906020012060001c6110c791906114ac565b92915050565b61110060408051606081018252600080825260208083018290528351808501855282815290810191909152909182015290565b506004810154600090815260209182526040908190208151606081018352815481526001820154818501528251808401845260028301548152600390920154938201939093529082015290565b8281526060810161116b602083018480518252602090810151910152565b9392505050565b8051801515811461118257600080fd5b919050565b6040805190810167ffffffffffffffff811182821017156111b857634e487b7160e01b600052604160045260246000fd5b60405290565b60008082840360a08112156111d257600080fd5b6111db84611172565b92506080601f19820112156111ef57600080fd5b6040516060810181811067ffffffffffffffff8211171561122057634e487b7160e01b600052604160045260246000fd5b604090815260208681015183528682015190830152605f198301121561124557600080fd5b61124d611187565b606086015181526080909501516020860152604081019490945250909391925050565b60208082526025908201527f476574526f6c6c7570427949443a20666972737420726f6c6c7570206e6f7420604082015264199bdd5b9960da1b606082015260800190565b6020808252602a908201527f476574526f6c6c75704279486173683a20617070656e64656420726f6c6c7570604082015269081b9bdd08199bdd5b9960b21b606082015260800190565b600061ffff8083168181141561132557634e487b7160e01b600052601160045260246000fd5b6001019392505050565b60208082526026908201527f47657448656164526f6c6c75703a20756e6578706563746564206c61737420656040820152651b195b595b9d60d21b606082015260800190565b60208082526032908201527f47657448656164526f6c6c75703a20756e6578706563746564206c61737420656040820152710d8cadacadce840e4ded8d8eae04090c2e6d60731b606082015260800190565b600060a082019050838252825160208301526020830151604083015260408301516113ff606084018280518252602090810151910152565b509392505050565b60208082526024908201527f476574506172656e74526f6c6c75703a20657870656374656420656e6420726f60408201526306c6c75760e41b606082015260800190565b60208082526026908201527f476574506172656e74526f6c6c75703a20756e657870656374656420656e64206040820152650726f6c6c75760d41b606082015260800190565b6000602082840312156114a357600080fd5b61116b82611172565b6000826114c957634e487b7160e01b600052601260045260246000fd5b50069056fea2646970667358221220220c4dc39696a167482f13f695d3fc5646d244ca2afa6c47daac1fd4ff0e167264736f6c634300080c0033",
}

// RollupChainTestContractABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupChainTestContractMetaData.ABI instead.
var RollupChainTestContractABI = RollupChainTestContractMetaData.ABI

// Deprecated: Use RollupChainTestContractMetaData.Sigs instead.
// RollupChainTestContractFuncSigs maps the 4-byte function signature to its string representation.
var RollupChainTestContractFuncSigs = RollupChainTestContractMetaData.Sigs

// RollupChainTestContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RollupChainTestContractMetaData.Bin instead.
var RollupChainTestContractBin = RollupChainTestContractMetaData.Bin

// DeployRollupChainTestContract deploys a new Ethereum contract, binding an instance of RollupChainTestContract to it.
func DeployRollupChainTestContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RollupChainTestContract, error) {
	parsed, err := RollupChainTestContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	rollupChainAddr, _, _, _ := DeployRollupChain(auth, backend)
	RollupChainTestContractBin = strings.Replace(RollupChainTestContractBin, "__$7204b1ba8a254ced74f31676d70e6726eb$__", rollupChainAddr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RollupChainTestContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RollupChainTestContract{RollupChainTestContractCaller: RollupChainTestContractCaller{contract: contract}, RollupChainTestContractTransactor: RollupChainTestContractTransactor{contract: contract}, RollupChainTestContractFilterer: RollupChainTestContractFilterer{contract: contract}}, nil
}

// RollupChainTestContract is an auto generated Go binding around an Ethereum contract.
type RollupChainTestContract struct {
	RollupChainTestContractCaller     // Read-only binding to the contract
	RollupChainTestContractTransactor // Write-only binding to the contract
	RollupChainTestContractFilterer   // Log filterer for contract events
}

// RollupChainTestContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupChainTestContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainTestContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupChainTestContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainTestContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupChainTestContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupChainTestContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupChainTestContractSession struct {
	Contract     *RollupChainTestContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RollupChainTestContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupChainTestContractCallerSession struct {
	Contract *RollupChainTestContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// RollupChainTestContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupChainTestContractTransactorSession struct {
	Contract     *RollupChainTestContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// RollupChainTestContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupChainTestContractRaw struct {
	Contract *RollupChainTestContract // Generic contract binding to access the raw methods on
}

// RollupChainTestContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupChainTestContractCallerRaw struct {
	Contract *RollupChainTestContractCaller // Generic read-only contract binding to access the raw methods on
}

// RollupChainTestContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupChainTestContractTransactorRaw struct {
	Contract *RollupChainTestContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupChainTestContract creates a new instance of RollupChainTestContract, bound to a specific deployed contract.
func NewRollupChainTestContract(address common.Address, backend bind.ContractBackend) (*RollupChainTestContract, error) {
	contract, err := bindRollupChainTestContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollupChainTestContract{RollupChainTestContractCaller: RollupChainTestContractCaller{contract: contract}, RollupChainTestContractTransactor: RollupChainTestContractTransactor{contract: contract}, RollupChainTestContractFilterer: RollupChainTestContractFilterer{contract: contract}}, nil
}

// NewRollupChainTestContractCaller creates a new read-only instance of RollupChainTestContract, bound to a specific deployed contract.
func NewRollupChainTestContractCaller(address common.Address, caller bind.ContractCaller) (*RollupChainTestContractCaller, error) {
	contract, err := bindRollupChainTestContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupChainTestContractCaller{contract: contract}, nil
}

// NewRollupChainTestContractTransactor creates a new write-only instance of RollupChainTestContract, bound to a specific deployed contract.
func NewRollupChainTestContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupChainTestContractTransactor, error) {
	contract, err := bindRollupChainTestContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupChainTestContractTransactor{contract: contract}, nil
}

// NewRollupChainTestContractFilterer creates a new log filterer instance of RollupChainTestContract, bound to a specific deployed contract.
func NewRollupChainTestContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupChainTestContractFilterer, error) {
	contract, err := bindRollupChainTestContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupChainTestContractFilterer{contract: contract}, nil
}

// bindRollupChainTestContract binds a generic wrapper to an already deployed contract.
func bindRollupChainTestContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RollupChainTestContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupChainTestContract *RollupChainTestContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupChainTestContract.Contract.RollupChainTestContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupChainTestContract *RollupChainTestContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.RollupChainTestContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupChainTestContract *RollupChainTestContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.RollupChainTestContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupChainTestContract *RollupChainTestContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupChainTestContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupChainTestContract *RollupChainTestContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupChainTestContract *RollupChainTestContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.contract.Transact(opts, method, params...)
}

// AppendRollupTest is a paid mutator transaction binding the contract method 0xf6da5eb5.
//
// Solidity: function AppendRollupTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactor) AppendRollupTest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.contract.Transact(opts, "AppendRollupTest")
}

// AppendRollupTest is a paid mutator transaction binding the contract method 0xf6da5eb5.
//
// Solidity: function AppendRollupTest() returns()
func (_RollupChainTestContract *RollupChainTestContractSession) AppendRollupTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.AppendRollupTest(&_RollupChainTestContract.TransactOpts)
}

// AppendRollupTest is a paid mutator transaction binding the contract method 0xf6da5eb5.
//
// Solidity: function AppendRollupTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactorSession) AppendRollupTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.AppendRollupTest(&_RollupChainTestContract.TransactOpts)
}

// NoForkDetection is a paid mutator transaction binding the contract method 0x01ac6e5b.
//
// Solidity: function NoForkDetection() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactor) NoForkDetection(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.contract.Transact(opts, "NoForkDetection")
}

// NoForkDetection is a paid mutator transaction binding the contract method 0x01ac6e5b.
//
// Solidity: function NoForkDetection() returns()
func (_RollupChainTestContract *RollupChainTestContractSession) NoForkDetection() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.NoForkDetection(&_RollupChainTestContract.TransactOpts)
}

// NoForkDetection is a paid mutator transaction binding the contract method 0x01ac6e5b.
//
// Solidity: function NoForkDetection() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactorSession) NoForkDetection() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.NoForkDetection(&_RollupChainTestContract.TransactOpts)
}

// RevertsNoDoubleInitTest is a paid mutator transaction binding the contract method 0x4762fa7c.
//
// Solidity: function RevertsNoDoubleInitTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactor) RevertsNoDoubleInitTest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.contract.Transact(opts, "RevertsNoDoubleInitTest")
}

// RevertsNoDoubleInitTest is a paid mutator transaction binding the contract method 0x4762fa7c.
//
// Solidity: function RevertsNoDoubleInitTest() returns()
func (_RollupChainTestContract *RollupChainTestContractSession) RevertsNoDoubleInitTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.RevertsNoDoubleInitTest(&_RollupChainTestContract.TransactOpts)
}

// RevertsNoDoubleInitTest is a paid mutator transaction binding the contract method 0x4762fa7c.
//
// Solidity: function RevertsNoDoubleInitTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactorSession) RevertsNoDoubleInitTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.RevertsNoDoubleInitTest(&_RollupChainTestContract.TransactOpts)
}

// ScrollTreeTest is a paid mutator transaction binding the contract method 0xf2b5185d.
//
// Solidity: function ScrollTreeTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactor) ScrollTreeTest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupChainTestContract.contract.Transact(opts, "ScrollTreeTest")
}

// ScrollTreeTest is a paid mutator transaction binding the contract method 0xf2b5185d.
//
// Solidity: function ScrollTreeTest() returns()
func (_RollupChainTestContract *RollupChainTestContractSession) ScrollTreeTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.ScrollTreeTest(&_RollupChainTestContract.TransactOpts)
}

// ScrollTreeTest is a paid mutator transaction binding the contract method 0xf2b5185d.
//
// Solidity: function ScrollTreeTest() returns()
func (_RollupChainTestContract *RollupChainTestContractTransactorSession) ScrollTreeTest() (*types.Transaction, error) {
	return _RollupChainTestContract.Contract.ScrollTreeTest(&_RollupChainTestContract.TransactOpts)
}
