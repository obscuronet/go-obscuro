// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

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

// GeneratedManagementContractMetaData contains all meta data concerning the GeneratedManagementContract contract.
var GeneratedManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hostAddresses\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506119b9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063d4c8066411610066578063d4c8066414610132578063e0643dfc14610162578063e0fd84bd14610195578063e34fbfc8146101b1578063f1846d0c146101cd57610093565b8063324ff8661461009857806365a293c2146100b657806368e10383146100e65780638ef74f8914610102575b600080fd5b6100a06101e9565b6040516100ad9190610b0a565b60405180910390f35b6100d060048036038101906100cb9190610b76565b6102c2565b6040516100dd9190610bed565b60405180910390f35b61010060048036038101906100fb9190610e02565b61036e565b005b61011c60048036038101906101179190610e92565b610436565b6040516101299190610bed565b60405180910390f35b61014c60048036038101906101479190610e92565b6104d6565b6040516101599190610eda565b60405180910390f35b61017c60048036038101906101779190610ef5565b6104f6565b60405161018c9493929190610f6c565b60405180910390f35b6101af60048036038101906101aa9190611033565b610563565b005b6101cb60048036038101906101c691906110cd565b6106a0565b005b6101e760048036038101906101e291906111bb565b6106f3565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156102b957838290600052602060002001805461022c906112b9565b80601f0160208091040260200160405190810160405280929190818152602001828054610258906112b9565b80156102a55780601f1061027a576101008083540402835291602001916102a5565b820191906000526020600020905b81548152906001019060200180831161028857829003601f168201915b50505050508152602001906001019061020d565b50505050905090565b600381815481106102d257600080fd5b9060005260206000200160009150905080546102ed906112b9565b80601f0160208091040260200160405190810160405280929190818152602001828054610319906112b9565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b600460009054906101000a900460ff161561038857600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555060038190806001815401808255809150506001900390600052602060002001600090919091909150908161042f9190611496565b5050505050565b60016020528060005260406000206000915090508054610455906112b9565b80601f0160208091040260200160405190810160405280929190818152602001828054610481906112b9565b80156104ce5780601f106104a3576101008083540402835291602001916104ce565b820191906000526020600020905b8154815290600101906020018083116104b157829003601f168201915b505050505081565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061051257600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166105b957600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826106ee929190611573565b505050565b6000600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061074e57600080fd5b600061077c878786604051602001610768939291906116d2565b604051602081830303815290604052610814565b90506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555060038390806001815401808255809150506001900390600052602060002001600090919091909150908161080a9190611496565b5050505050505050565b6000610820825161084f565b82604051602001610832929190611793565b604051602081830303815290604052805190602001209050919050565b606060008203610896576040518060400160405280600181526020017f300000000000000000000000000000000000000000000000000000000000000081525090506109aa565b600082905060005b600082146108c85780806108b1906117f1565b915050600a826108c19190611868565b915061089e565b60008167ffffffffffffffff8111156108e4576108e3610cd7565b5b6040519080825280601f01601f1916602001820160405280156109165781602001600182028036833780820191505090505b5090505b600085146109a35760018261092f9190611899565b9150600a8561093e91906118cd565b603061094a91906118fe565b60f81b8183815181106109605761095f611954565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a8561099c9190611868565b945061091a565b8093505050505b919050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610a155780820151818401526020810190506109fa565b83811115610a24576000848401525b50505050565b6000601f19601f8301169050919050565b6000610a46826109db565b610a5081856109e6565b9350610a608185602086016109f7565b610a6981610a2a565b840191505092915050565b6000610a808383610a3b565b905092915050565b6000602082019050919050565b6000610aa0826109af565b610aaa81856109ba565b935083602082028501610abc856109cb565b8060005b85811015610af85784840389528151610ad98582610a74565b9450610ae483610a88565b925060208a01995050600181019050610ac0565b50829750879550505050505092915050565b60006020820190508181036000830152610b248184610a95565b905092915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b610b5381610b40565b8114610b5e57600080fd5b50565b600081359050610b7081610b4a565b92915050565b600060208284031215610b8c57610b8b610b36565b5b6000610b9a84828501610b61565b91505092915050565b600082825260208201905092915050565b6000610bbf826109db565b610bc98185610ba3565b9350610bd98185602086016109f7565b610be281610a2a565b840191505092915050565b60006020820190508181036000830152610c078184610bb4565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c3a82610c0f565b9050919050565b610c4a81610c2f565b8114610c5557600080fd5b50565b600081359050610c6781610c41565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610c9257610c91610c6d565b5b8235905067ffffffffffffffff811115610caf57610cae610c72565b5b602083019150836001820283011115610ccb57610cca610c77565b5b9250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610d0f82610a2a565b810181811067ffffffffffffffff82111715610d2e57610d2d610cd7565b5b80604052505050565b6000610d41610b2c565b9050610d4d8282610d06565b919050565b600067ffffffffffffffff821115610d6d57610d6c610cd7565b5b610d7682610a2a565b9050602081019050919050565b82818337600083830152505050565b6000610da5610da084610d52565b610d37565b905082815260208101848484011115610dc157610dc0610cd2565b5b610dcc848285610d83565b509392505050565b600082601f830112610de957610de8610c6d565b5b8135610df9848260208601610d92565b91505092915050565b60008060008060608587031215610e1c57610e1b610b36565b5b6000610e2a87828801610c58565b945050602085013567ffffffffffffffff811115610e4b57610e4a610b3b565b5b610e5787828801610c7c565b9350935050604085013567ffffffffffffffff811115610e7a57610e79610b3b565b5b610e8687828801610dd4565b91505092959194509250565b600060208284031215610ea857610ea7610b36565b5b6000610eb684828501610c58565b91505092915050565b60008115159050919050565b610ed481610ebf565b82525050565b6000602082019050610eef6000830184610ecb565b92915050565b60008060408385031215610f0c57610f0b610b36565b5b6000610f1a85828601610b61565b9250506020610f2b85828601610b61565b9150509250929050565b6000819050919050565b610f4881610f35565b82525050565b610f5781610c2f565b82525050565b610f6681610b40565b82525050565b6000608082019050610f816000830187610f3f565b610f8e6020830186610f4e565b610f9b6040830185610f3f565b610fa86060830184610f5d565b95945050505050565b610fba81610f35565b8114610fc557600080fd5b50565b600081359050610fd781610fb1565b92915050565b60008083601f840112610ff357610ff2610c6d565b5b8235905067ffffffffffffffff8111156110105761100f610c72565b5b60208301915083600182028301111561102c5761102b610c77565b5b9250929050565b60008060008060008060a087890312156110505761104f610b36565b5b600061105e89828a01610fc8565b965050602061106f89828a01610c58565b955050604061108089828a01610fc8565b945050606061109189828a01610b61565b935050608087013567ffffffffffffffff8111156110b2576110b1610b3b565b5b6110be89828a01610fdd565b92509250509295509295509295565b600080602083850312156110e4576110e3610b36565b5b600083013567ffffffffffffffff81111561110257611101610b3b565b5b61110e85828601610fdd565b92509250509250929050565b600067ffffffffffffffff82111561113557611134610cd7565b5b61113e82610a2a565b9050602081019050919050565b600061115e6111598461111a565b610d37565b90508281526020810184848401111561117a57611179610cd2565b5b611185848285610d83565b509392505050565b600082601f8301126111a2576111a1610c6d565b5b81356111b284826020860161114b565b91505092915050565b600080600080600060a086880312156111d7576111d6610b36565b5b60006111e588828901610c58565b95505060206111f688828901610c58565b945050604086013567ffffffffffffffff81111561121757611216610b3b565b5b6112238882890161118d565b935050606086013567ffffffffffffffff81111561124457611243610b3b565b5b6112508882890161118d565b925050608086013567ffffffffffffffff81111561127157611270610b3b565b5b61127d88828901610dd4565b9150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806112d157607f821691505b6020821081036112e4576112e361128a565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261134c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261130f565b611356868361130f565b95508019841693508086168417925050509392505050565b6000819050919050565b600061139361138e61138984610b40565b61136e565b610b40565b9050919050565b6000819050919050565b6113ad83611378565b6113c16113b98261139a565b84845461131c565b825550505050565b600090565b6113d66113c9565b6113e18184846113a4565b505050565b5b81811015611405576113fa6000826113ce565b6001810190506113e7565b5050565b601f82111561144a5761141b816112ea565b611424846112ff565b81016020851015611433578190505b61144761143f856112ff565b8301826113e6565b50505b505050565b600082821c905092915050565b600061146d6000198460080261144f565b1980831691505092915050565b6000611486838361145c565b9150826002028217905092915050565b61149f826109db565b67ffffffffffffffff8111156114b8576114b7610cd7565b5b6114c282546112b9565b6114cd828285611409565b600060209050601f83116001811461150057600084156114ee578287015190505b6114f8858261147a565b865550611560565b601f19841661150e866112ea565b60005b8281101561153657848901518255600182019150602085019450602081019050611511565b86831015611553578489015161154f601f89168261145c565b8355505b6001600288020188555050505b505050505050565b600082905092915050565b61157d8383611568565b67ffffffffffffffff81111561159657611595610cd7565b5b6115a082546112b9565b6115ab828285611409565b6000601f8311600181146115da57600084156115c8578287013590505b6115d2858261147a565b86555061163a565b601f1984166115e8866112ea565b60005b82811015611610578489013582556001820191506020850194506020810190506115eb565b8683101561162d5784890135611629601f89168261145c565b8355505b6001600288020188555050505b50505050505050565b60008160601b9050919050565b600061165b82611643565b9050919050565b600061166d82611650565b9050919050565b61168561168082610c2f565b611662565b82525050565b600081519050919050565b600081905092915050565b60006116ac8261168b565b6116b68185611696565b93506116c68185602086016109f7565b80840191505092915050565b60006116de8286611674565b6014820191506116ee8285611674565b6014820191506116fe82846116a1565b9150819050949350505050565b600081905092915050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b600061174c601a8361170b565b915061175782611716565b601a82019050919050565b600061176d826109db565b611777818561170b565b93506117878185602086016109f7565b80840191505092915050565b600061179e8261173f565b91506117aa8285611762565b91506117b682846116a1565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006117fc82610b40565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361182e5761182d6117c2565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061187382610b40565b915061187e83610b40565b92508261188e5761188d611839565b5b828204905092915050565b60006118a482610b40565b91506118af83610b40565b9250828210156118c2576118c16117c2565b5b828203905092915050565b60006118d882610b40565b91506118e383610b40565b9250826118f3576118f2611839565b5b828206905092915050565b600061190982610b40565b915061191483610b40565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611949576119486117c2565b5b828201905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea26469706673582212209e1f25417d11783742d4683ecde949cb17dc73de8af25258a8ccd7a40fd21ff064736f6c634300080f0033",
}

// GeneratedManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GeneratedManagementContractMetaData.ABI instead.
var GeneratedManagementContractABI = GeneratedManagementContractMetaData.ABI

// GeneratedManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GeneratedManagementContractMetaData.Bin instead.
var GeneratedManagementContractBin = GeneratedManagementContractMetaData.Bin

// DeployGeneratedManagementContract deploys a new Ethereum contract, binding an instance of GeneratedManagementContract to it.
func DeployGeneratedManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GeneratedManagementContract, error) {
	parsed, err := GeneratedManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GeneratedManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// GeneratedManagementContract is an auto generated Go binding around an Ethereum contract.
type GeneratedManagementContract struct {
	GeneratedManagementContractCaller     // Read-only binding to the contract
	GeneratedManagementContractTransactor // Write-only binding to the contract
	GeneratedManagementContractFilterer   // Log filterer for contract events
}

// GeneratedManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedManagementContractSession struct {
	Contract     *GeneratedManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GeneratedManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedManagementContractCallerSession struct {
	Contract *GeneratedManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// GeneratedManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedManagementContractTransactorSession struct {
	Contract     *GeneratedManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// GeneratedManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedManagementContractRaw struct {
	Contract *GeneratedManagementContract // Generic contract binding to access the raw methods on
}

// GeneratedManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCallerRaw struct {
	Contract *GeneratedManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactorRaw struct {
	Contract *GeneratedManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeneratedManagementContract creates a new instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContract(address common.Address, backend bind.ContractBackend) (*GeneratedManagementContract, error) {
	contract, err := bindGeneratedManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// NewGeneratedManagementContractCaller creates a new read-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractCaller(address common.Address, caller bind.ContractCaller) (*GeneratedManagementContractCaller, error) {
	contract, err := bindGeneratedManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractCaller{contract: contract}, nil
}

// NewGeneratedManagementContractTransactor creates a new write-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedManagementContractTransactor, error) {
	contract, err := bindGeneratedManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractTransactor{contract: contract}, nil
}

// NewGeneratedManagementContractFilterer creates a new log filterer instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedManagementContractFilterer, error) {
	contract, err := bindGeneratedManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractFilterer{contract: contract}, nil
}

// bindGeneratedManagementContract binds a generic wrapper to an already deployed contract.
func bindGeneratedManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transact(opts, method, params...)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCaller) GetHostAddresses(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "GetHostAddresses")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) AttestationRequests(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attestationRequests", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Attested(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attested", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) HostAddresses(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "hostAddresses", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) HostAddresses(arg0 *big.Int) (string, error) {
	return _GeneratedManagementContract.Contract.HostAddresses(&_GeneratedManagementContract.CallOpts, arg0)
}

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) HostAddresses(arg0 *big.Int) (string, error) {
	return _GeneratedManagementContract.Contract.HostAddresses(&_GeneratedManagementContract.CallOpts, arg0)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		AggregatorID common.Address
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.L1Block = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret, hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}
