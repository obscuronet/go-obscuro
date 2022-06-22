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
	Bin: "0x608060405234801561001057600080fd5b50612729806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063d4c8066411610066578063d4c8066414610132578063e0643dfc14610162578063e0fd84bd14610195578063e34fbfc8146101b1578063f1846d0c146101cd57610093565b8063324ff8661461009857806365a293c2146100b657806368e10383146100e65780638ef74f8914610102575b600080fd5b6100a06101e9565b6040516100ad91906111a4565b60405180910390f35b6100d060048036038101906100cb9190611210565b6102c2565b6040516100dd9190611287565b60405180910390f35b61010060048036038101906100fb919061149c565b61036e565b005b61011c6004803603810190610117919061152c565b610436565b6040516101299190611287565b60405180910390f35b61014c6004803603810190610147919061152c565b6104d6565b6040516101599190611574565b60405180910390f35b61017c6004803603810190610177919061158f565b6104f6565b60405161018c9493929190611606565b60405180910390f35b6101af60048036038101906101aa91906116cd565b610563565b005b6101cb60048036038101906101c69190611767565b6106a0565b005b6101e760048036038101906101e29190611855565b6106f3565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156102b957838290600052602060002001805461022c90611953565b80601f016020809104026020016040519081016040528092919081815260200182805461025890611953565b80156102a55780601f1061027a576101008083540402835291602001916102a5565b820191906000526020600020905b81548152906001019060200180831161028857829003601f168201915b50505050508152602001906001019061020d565b50505050905090565b600381815481106102d257600080fd5b9060005260206000200160009150905080546102ed90611953565b80601f016020809104026020016040519081016040528092919081815260200182805461031990611953565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b600460009054906101000a900460ff161561038857600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555060038190806001815401808255809150506001900390600052602060002001600090919091909150908161042f9190611b30565b5050505050565b6001602052806000526040600020600091509050805461045590611953565b80601f016020809104026020016040519081016040528092919081815260200182805461048190611953565b80156104ce5780601f106104a3576101008083540402835291602001916104ce565b820191906000526020600020905b8154815290600101906020018083116104b157829003601f168201915b505050505081565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061051257600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166105b957600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826106ee929190611c0d565b505050565b6000600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061074e57600080fd5b600061077c87878660405160200161076893929190611d6c565b6040516020818303038152906040526108c6565b9050600061078a8287610901565b90508773ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146107c489610928565b6107cd83610928565b6040516020016107de929190611f37565b6040516020818303038152906040529061082e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108259190611287565b60405180910390fd5b506001600260008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003849080600181540180825580915050600190039060005260206000200160009091909190915090816108bb9190611b30565b505050505050505050565b60006108d28251610aeb565b826040516020016108e4929190611fc8565b604051602081830303815290604052805190602001209050919050565b60008060006109108585610c4b565b9150915061091d81610ccc565b819250505092915050565b60606000602867ffffffffffffffff81111561094757610946611371565b5b6040519080825280601f01601f1916602001820160405280156109795781602001600182028036833780820191505090505b50905060005b6014811015610ae15760008160136109979190612026565b60086109a3919061205a565b60026109af91906121e7565b8573ffffffffffffffffffffffffffffffffffffffff166109d09190612261565b60f81b9050600060108260f81c6109e7919061229f565b60f81b905060008160f81c60106109fe91906122d0565b8360f81c610a0c919061230b565b60f81b9050610a1a82610e98565b85856002610a28919061205a565b81518110610a3957610a3861233f565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350610a7181610e98565b856001866002610a81919061205a565b610a8b919061236e565b81518110610a9c57610a9b61233f565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505050508080610ad9906123c4565b91505061097f565b5080915050919050565b606060008203610b32576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610c46565b600082905060005b60008214610b64578080610b4d906123c4565b915050600a82610b5d9190612261565b9150610b3a565b60008167ffffffffffffffff811115610b8057610b7f611371565b5b6040519080825280601f01601f191660200182016040528015610bb25781602001600182028036833780820191505090505b5090505b60008514610c3f57600182610bcb9190612026565b9150600a85610bda919061240c565b6030610be6919061236e565b60f81b818381518110610bfc57610bfb61233f565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610c389190612261565b9450610bb6565b8093505050505b919050565b6000806041835103610c8c5760008060006020860151925060408601519150606086015160001a9050610c8087828585610ede565b94509450505050610cc5565b6040835103610cbc576000806020850151915060408501519050610cb1868383610fea565b935093505050610cc5565b60006002915091505b9250929050565b60006004811115610ce057610cdf61243d565b5b816004811115610cf357610cf261243d565b5b0315610e955760016004811115610d0d57610d0c61243d565b5b816004811115610d2057610d1f61243d565b5b03610d60576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d57906124b8565b60405180910390fd5b60026004811115610d7457610d7361243d565b5b816004811115610d8757610d8661243d565b5b03610dc7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dbe90612524565b60405180910390fd5b60036004811115610ddb57610dda61243d565b5b816004811115610dee57610ded61243d565b5b03610e2e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e25906125b6565b60405180910390fd5b600480811115610e4157610e4061243d565b5b816004811115610e5457610e5361243d565b5b03610e94576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e8b90612648565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610ec35760308260f81c610eb99190612668565b60f81b9050610ed9565b60578260f81c610ed39190612668565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610f19576000600391509150610fe1565b601b8560ff1614158015610f315750601c8560ff1614155b15610f43576000600491509150610fe1565b600060018787878760405160008152602001604052604051610f6894939291906126ae565b6020604051602081039080840390855afa158015610f8a573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610fd857600060019250925050610fe1565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c61102d919061236e565b905061103b87828885610ede565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b838110156110af578082015181840152602081019050611094565b838111156110be576000848401525b50505050565b6000601f19601f8301169050919050565b60006110e082611075565b6110ea8185611080565b93506110fa818560208601611091565b611103816110c4565b840191505092915050565b600061111a83836110d5565b905092915050565b6000602082019050919050565b600061113a82611049565b6111448185611054565b93508360208202850161115685611065565b8060005b858110156111925784840389528151611173858261110e565b945061117e83611122565b925060208a0199505060018101905061115a565b50829750879550505050505092915050565b600060208201905081810360008301526111be818461112f565b905092915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6111ed816111da565b81146111f857600080fd5b50565b60008135905061120a816111e4565b92915050565b600060208284031215611226576112256111d0565b5b6000611234848285016111fb565b91505092915050565b600082825260208201905092915050565b600061125982611075565b611263818561123d565b9350611273818560208601611091565b61127c816110c4565b840191505092915050565b600060208201905081810360008301526112a1818461124e565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006112d4826112a9565b9050919050565b6112e4816112c9565b81146112ef57600080fd5b50565b600081359050611301816112db565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261132c5761132b611307565b5b8235905067ffffffffffffffff8111156113495761134861130c565b5b60208301915083600182028301111561136557611364611311565b5b9250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6113a9826110c4565b810181811067ffffffffffffffff821117156113c8576113c7611371565b5b80604052505050565b60006113db6111c6565b90506113e782826113a0565b919050565b600067ffffffffffffffff82111561140757611406611371565b5b611410826110c4565b9050602081019050919050565b82818337600083830152505050565b600061143f61143a846113ec565b6113d1565b90508281526020810184848401111561145b5761145a61136c565b5b61146684828561141d565b509392505050565b600082601f83011261148357611482611307565b5b813561149384826020860161142c565b91505092915050565b600080600080606085870312156114b6576114b56111d0565b5b60006114c4878288016112f2565b945050602085013567ffffffffffffffff8111156114e5576114e46111d5565b5b6114f187828801611316565b9350935050604085013567ffffffffffffffff811115611514576115136111d5565b5b6115208782880161146e565b91505092959194509250565b600060208284031215611542576115416111d0565b5b6000611550848285016112f2565b91505092915050565b60008115159050919050565b61156e81611559565b82525050565b60006020820190506115896000830184611565565b92915050565b600080604083850312156115a6576115a56111d0565b5b60006115b4858286016111fb565b92505060206115c5858286016111fb565b9150509250929050565b6000819050919050565b6115e2816115cf565b82525050565b6115f1816112c9565b82525050565b611600816111da565b82525050565b600060808201905061161b60008301876115d9565b61162860208301866115e8565b61163560408301856115d9565b61164260608301846115f7565b95945050505050565b611654816115cf565b811461165f57600080fd5b50565b6000813590506116718161164b565b92915050565b60008083601f84011261168d5761168c611307565b5b8235905067ffffffffffffffff8111156116aa576116a961130c565b5b6020830191508360018202830111156116c6576116c5611311565b5b9250929050565b60008060008060008060a087890312156116ea576116e96111d0565b5b60006116f889828a01611662565b965050602061170989828a016112f2565b955050604061171a89828a01611662565b945050606061172b89828a016111fb565b935050608087013567ffffffffffffffff81111561174c5761174b6111d5565b5b61175889828a01611677565b92509250509295509295509295565b6000806020838503121561177e5761177d6111d0565b5b600083013567ffffffffffffffff81111561179c5761179b6111d5565b5b6117a885828601611677565b92509250509250929050565b600067ffffffffffffffff8211156117cf576117ce611371565b5b6117d8826110c4565b9050602081019050919050565b60006117f86117f3846117b4565b6113d1565b9050828152602081018484840111156118145761181361136c565b5b61181f84828561141d565b509392505050565b600082601f83011261183c5761183b611307565b5b813561184c8482602086016117e5565b91505092915050565b600080600080600060a08688031215611871576118706111d0565b5b600061187f888289016112f2565b9550506020611890888289016112f2565b945050604086013567ffffffffffffffff8111156118b1576118b06111d5565b5b6118bd88828901611827565b935050606086013567ffffffffffffffff8111156118de576118dd6111d5565b5b6118ea88828901611827565b925050608086013567ffffffffffffffff81111561190b5761190a6111d5565b5b6119178882890161146e565b9150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061196b57607f821691505b60208210810361197e5761197d611924565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026119e67fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826119a9565b6119f086836119a9565b95508019841693508086168417925050509392505050565b6000819050919050565b6000611a2d611a28611a23846111da565b611a08565b6111da565b9050919050565b6000819050919050565b611a4783611a12565b611a5b611a5382611a34565b8484546119b6565b825550505050565b600090565b611a70611a63565b611a7b818484611a3e565b505050565b5b81811015611a9f57611a94600082611a68565b600181019050611a81565b5050565b601f821115611ae457611ab581611984565b611abe84611999565b81016020851015611acd578190505b611ae1611ad985611999565b830182611a80565b50505b505050565b600082821c905092915050565b6000611b0760001984600802611ae9565b1980831691505092915050565b6000611b208383611af6565b9150826002028217905092915050565b611b3982611075565b67ffffffffffffffff811115611b5257611b51611371565b5b611b5c8254611953565b611b67828285611aa3565b600060209050601f831160018114611b9a5760008415611b88578287015190505b611b928582611b14565b865550611bfa565b601f198416611ba886611984565b60005b82811015611bd057848901518255600182019150602085019450602081019050611bab565b86831015611bed5784890151611be9601f891682611af6565b8355505b6001600288020188555050505b505050505050565b600082905092915050565b611c178383611c02565b67ffffffffffffffff811115611c3057611c2f611371565b5b611c3a8254611953565b611c45828285611aa3565b6000601f831160018114611c745760008415611c62578287013590505b611c6c8582611b14565b865550611cd4565b601f198416611c8286611984565b60005b82811015611caa57848901358255600182019150602085019450602081019050611c85565b86831015611cc75784890135611cc3601f891682611af6565b8355505b6001600288020188555050505b50505050505050565b60008160601b9050919050565b6000611cf582611cdd565b9050919050565b6000611d0782611cea565b9050919050565b611d1f611d1a826112c9565b611cfc565b82525050565b600081519050919050565b600081905092915050565b6000611d4682611d25565b611d508185611d30565b9350611d60818560208601611091565b80840191505092915050565b6000611d788286611d0e565b601482019150611d888285611d0e565b601482019150611d988284611d3b565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611e0c602d83611da5565b9150611e1782611db0565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611e7e602483611da5565b9150611e8982611e22565b602482019050919050565b6000611e9f82611075565b611ea98185611da5565b9350611eb9818560208601611091565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611f21602483611da5565b9150611f2c82611ec5565b602482019050919050565b6000611f4282611dff565b9150611f4d82611e71565b9150611f598285611e94565b9150611f6482611f14565b9150611f708284611e94565b91508190509392505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611fb2601a83611da5565b9150611fbd82611f7c565b601a82019050919050565b6000611fd382611fa5565b9150611fdf8285611e94565b9150611feb8284611d3b565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000612031826111da565b915061203c836111da565b92508282101561204f5761204e611ff7565b5b828203905092915050565b6000612065826111da565b9150612070836111da565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156120a9576120a8611ff7565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b600185111561210b578086048111156120e7576120e6611ff7565b5b60018516156120f65780820291505b8081029050612104856120b4565b94506120cb565b94509492505050565b60008261212457600190506121e0565b8161213257600090506121e0565b8160018114612148576002811461215257612181565b60019150506121e0565b60ff84111561216457612163611ff7565b5b8360020a91508482111561217b5761217a611ff7565b5b506121e0565b5060208310610133831016604e8410600b84101617156121b65782820a9050838111156121b1576121b0611ff7565b5b6121e0565b6121c384848460016120c1565b925090508184048111156121da576121d9611ff7565b5b81810290505b9392505050565b60006121f2826111da565b91506121fd836111da565b925061222a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484612114565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061226c826111da565b9150612277836111da565b92508261228757612286612232565b5b828204905092915050565b600060ff82169050919050565b60006122aa82612292565b91506122b583612292565b9250826122c5576122c4612232565b5b828204905092915050565b60006122db82612292565b91506122e683612292565b92508160ff0483118215151615612300576122ff611ff7565b5b828202905092915050565b600061231682612292565b915061232183612292565b92508282101561233457612333611ff7565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000612379826111da565b9150612384836111da565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156123b9576123b8611ff7565b5b828201905092915050565b60006123cf826111da565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361240157612400611ff7565b5b600182019050919050565b6000612417826111da565b9150612422836111da565b92508261243257612431612232565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b60006124a260188361123d565b91506124ad8261246c565b602082019050919050565b600060208201905081810360008301526124d181612495565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b600061250e601f8361123d565b9150612519826124d8565b602082019050919050565b6000602082019050818103600083015261253d81612501565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006125a060228361123d565b91506125ab82612544565b604082019050919050565b600060208201905081810360008301526125cf81612593565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061263260228361123d565b915061263d826125d6565b604082019050919050565b6000602082019050818103600083015261266181612625565b9050919050565b600061267382612292565b915061267e83612292565b92508260ff0382111561269457612693611ff7565b5b828201905092915050565b6126a881612292565b82525050565b60006080820190506126c360008301876115d9565b6126d0602083018661269f565b6126dd60408301856115d9565b6126ea60608301846115d9565b9594505050505056fea2646970667358221220716d0b769db7d0d075a3f51ca881378d1f6177c20fb1356aa0818fb78f819a9164736f6c634300080f0033",
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
