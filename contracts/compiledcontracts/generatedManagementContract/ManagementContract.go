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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hostAddresses\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506126c3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063d4c8066411610066578063d4c8066414610132578063e0643dfc14610162578063e0fd84bd14610195578063e34fbfc8146101b1578063f1846d0c146101cd57610093565b8063324ff8661461009857806365a293c2146100b65780638ef74f89146100e6578063c719bf5014610116575b600080fd5b6100a06101e9565b6040516100ad919061116e565b60405180910390f35b6100d060048036038101906100cb91906111da565b6102c2565b6040516100dd9190611251565b60405180910390f35b61010060048036038101906100fb91906112d1565b61036e565b60405161010d9190611251565b60405180910390f35b610130600480360381019061012b9190611363565b61040e565b005b61014c600480360381019061014791906112d1565b6104a0565b60405161015991906113de565b60405180910390f35b61017c600480360381019061017791906113f9565b6104c0565b60405161018c9493929190611470565b60405180910390f35b6101af60048036038101906101aa9190611537565b61052d565b005b6101cb60048036038101906101c691906115d1565b61066a565b005b6101e760048036038101906101e291906117ef565b6106bd565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156102b957838290600052602060002001805461022c906118ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610258906118ed565b80156102a55780601f1061027a576101008083540402835291602001916102a5565b820191906000526020600020905b81548152906001019060200180831161028857829003601f168201915b50505050508152602001906001019061020d565b50505050905090565b600381815481106102d257600080fd5b9060005260206000200160009150905080546102ed906118ed565b80601f0160208091040260200160405190810160405280929190818152602001828054610319906118ed565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b6001602052806000526040600020600091509050805461038d906118ed565b80601f01602080910402602001604051908101604052809291908181526020018280546103b9906118ed565b80156104065780601f106103db57610100808354040283529160200191610406565b820191906000526020600020905b8154815290600101906020018083116103e957829003601f168201915b505050505081565b600460009054906101000a900460ff161561042857600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550505050565b60026020528060005260406000206000915054906101000a900460ff1681565b600060205281600052604060002081815481106104dc57600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661058357600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826106b8929190611ad5565b505050565b6000600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061071857600080fd5b600061074687878660405160200161073293929190611c34565b604051602081830303815290604052610890565b9050600061075482876108cb565b90508773ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461078e896108f2565b610797836108f2565b6040516020016107a8929190611dff565b604051602081830303815290604052906107f8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107ef9190611251565b60405180910390fd5b506001600260008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003849080600181540180825580915050600190039060005260206000200160009091909190915090816108859190611e44565b505050505050505050565b600061089c8251610ab5565b826040516020016108ae929190611f62565b604051602081830303815290604052805190602001209050919050565b60008060006108da8585610c15565b915091506108e781610c96565b819250505092915050565b60606000602867ffffffffffffffff81111561091157610910611623565b5b6040519080825280601f01601f1916602001820160405280156109435781602001600182028036833780820191505090505b50905060005b6014811015610aab5760008160136109619190611fc0565b600861096d9190611ff4565b60026109799190612181565b8573ffffffffffffffffffffffffffffffffffffffff1661099a91906121fb565b60f81b9050600060108260f81c6109b19190612239565b60f81b905060008160f81c60106109c8919061226a565b8360f81c6109d691906122a5565b60f81b90506109e482610e62565b858560026109f29190611ff4565b81518110610a0357610a026122d9565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350610a3b81610e62565b856001866002610a4b9190611ff4565b610a559190612308565b81518110610a6657610a656122d9565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505050508080610aa39061235e565b915050610949565b5080915050919050565b606060008203610afc576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610c10565b600082905060005b60008214610b2e578080610b179061235e565b915050600a82610b2791906121fb565b9150610b04565b60008167ffffffffffffffff811115610b4a57610b49611623565b5b6040519080825280601f01601f191660200182016040528015610b7c5781602001600182028036833780820191505090505b5090505b60008514610c0957600182610b959190611fc0565b9150600a85610ba491906123a6565b6030610bb09190612308565b60f81b818381518110610bc657610bc56122d9565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610c0291906121fb565b9450610b80565b8093505050505b919050565b6000806041835103610c565760008060006020860151925060408601519150606086015160001a9050610c4a87828585610ea8565b94509450505050610c8f565b6040835103610c86576000806020850151915060408501519050610c7b868383610fb4565b935093505050610c8f565b60006002915091505b9250929050565b60006004811115610caa57610ca96123d7565b5b816004811115610cbd57610cbc6123d7565b5b0315610e5f5760016004811115610cd757610cd66123d7565b5b816004811115610cea57610ce96123d7565b5b03610d2a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d2190612452565b60405180910390fd5b60026004811115610d3e57610d3d6123d7565b5b816004811115610d5157610d506123d7565b5b03610d91576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d88906124be565b60405180910390fd5b60036004811115610da557610da46123d7565b5b816004811115610db857610db76123d7565b5b03610df8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610def90612550565b60405180910390fd5b600480811115610e0b57610e0a6123d7565b5b816004811115610e1e57610e1d6123d7565b5b03610e5e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e55906125e2565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610e8d5760308260f81c610e839190612602565b60f81b9050610ea3565b60578260f81c610e9d9190612602565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610ee3576000600391509150610fab565b601b8560ff1614158015610efb5750601c8560ff1614155b15610f0d576000600491509150610fab565b600060018787878760405160008152602001604052604051610f329493929190612648565b6020604051602081039080840390855afa158015610f54573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610fa257600060019250925050610fab565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610ff79190612308565b905061100587828885610ea8565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561107957808201518184015260208101905061105e565b83811115611088576000848401525b50505050565b6000601f19601f8301169050919050565b60006110aa8261103f565b6110b4818561104a565b93506110c481856020860161105b565b6110cd8161108e565b840191505092915050565b60006110e4838361109f565b905092915050565b6000602082019050919050565b600061110482611013565b61110e818561101e565b9350836020820285016111208561102f565b8060005b8581101561115c578484038952815161113d85826110d8565b9450611148836110ec565b925060208a01995050600181019050611124565b50829750879550505050505092915050565b6000602082019050818103600083015261118881846110f9565b905092915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6111b7816111a4565b81146111c257600080fd5b50565b6000813590506111d4816111ae565b92915050565b6000602082840312156111f0576111ef61119a565b5b60006111fe848285016111c5565b91505092915050565b600082825260208201905092915050565b60006112238261103f565b61122d8185611207565b935061123d81856020860161105b565b6112468161108e565b840191505092915050565b6000602082019050818103600083015261126b8184611218565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061129e82611273565b9050919050565b6112ae81611293565b81146112b957600080fd5b50565b6000813590506112cb816112a5565b92915050565b6000602082840312156112e7576112e661119a565b5b60006112f5848285016112bc565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083601f840112611323576113226112fe565b5b8235905067ffffffffffffffff8111156113405761133f611303565b5b60208301915083600182028301111561135c5761135b611308565b5b9250929050565b60008060006040848603121561137c5761137b61119a565b5b600061138a868287016112bc565b935050602084013567ffffffffffffffff8111156113ab576113aa61119f565b5b6113b78682870161130d565b92509250509250925092565b60008115159050919050565b6113d8816113c3565b82525050565b60006020820190506113f360008301846113cf565b92915050565b600080604083850312156114105761140f61119a565b5b600061141e858286016111c5565b925050602061142f858286016111c5565b9150509250929050565b6000819050919050565b61144c81611439565b82525050565b61145b81611293565b82525050565b61146a816111a4565b82525050565b60006080820190506114856000830187611443565b6114926020830186611452565b61149f6040830185611443565b6114ac6060830184611461565b95945050505050565b6114be81611439565b81146114c957600080fd5b50565b6000813590506114db816114b5565b92915050565b60008083601f8401126114f7576114f66112fe565b5b8235905067ffffffffffffffff81111561151457611513611303565b5b6020830191508360018202830111156115305761152f611308565b5b9250929050565b60008060008060008060a087890312156115545761155361119a565b5b600061156289828a016114cc565b965050602061157389828a016112bc565b955050604061158489828a016114cc565b945050606061159589828a016111c5565b935050608087013567ffffffffffffffff8111156115b6576115b561119f565b5b6115c289828a016114e1565b92509250509295509295509295565b600080602083850312156115e8576115e761119a565b5b600083013567ffffffffffffffff8111156116065761160561119f565b5b611612858286016114e1565b92509250509250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61165b8261108e565b810181811067ffffffffffffffff8211171561167a57611679611623565b5b80604052505050565b600061168d611190565b90506116998282611652565b919050565b600067ffffffffffffffff8211156116b9576116b8611623565b5b6116c28261108e565b9050602081019050919050565b82818337600083830152505050565b60006116f16116ec8461169e565b611683565b90508281526020810184848401111561170d5761170c61161e565b5b6117188482856116cf565b509392505050565b600082601f830112611735576117346112fe565b5b81356117458482602086016116de565b91505092915050565b600067ffffffffffffffff82111561176957611768611623565b5b6117728261108e565b9050602081019050919050565b600061179261178d8461174e565b611683565b9050828152602081018484840111156117ae576117ad61161e565b5b6117b98482856116cf565b509392505050565b600082601f8301126117d6576117d56112fe565b5b81356117e684826020860161177f565b91505092915050565b600080600080600060a0868803121561180b5761180a61119a565b5b6000611819888289016112bc565b955050602061182a888289016112bc565b945050604086013567ffffffffffffffff81111561184b5761184a61119f565b5b61185788828901611720565b935050606086013567ffffffffffffffff8111156118785761187761119f565b5b61188488828901611720565b925050608086013567ffffffffffffffff8111156118a5576118a461119f565b5b6118b1888289016117c1565b9150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061190557607f821691505b602082108103611918576119176118be565b5b50919050565b600082905092915050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261198b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261194e565b611995868361194e565b95508019841693508086168417925050509392505050565b6000819050919050565b60006119d26119cd6119c8846111a4565b6119ad565b6111a4565b9050919050565b6000819050919050565b6119ec836119b7565b611a006119f8826119d9565b84845461195b565b825550505050565b600090565b611a15611a08565b611a208184846119e3565b505050565b5b81811015611a4457611a39600082611a0d565b600181019050611a26565b5050565b601f821115611a8957611a5a81611929565b611a638461193e565b81016020851015611a72578190505b611a86611a7e8561193e565b830182611a25565b50505b505050565b600082821c905092915050565b6000611aac60001984600802611a8e565b1980831691505092915050565b6000611ac58383611a9b565b9150826002028217905092915050565b611adf838361191e565b67ffffffffffffffff811115611af857611af7611623565b5b611b0282546118ed565b611b0d828285611a48565b6000601f831160018114611b3c5760008415611b2a578287013590505b611b348582611ab9565b865550611b9c565b601f198416611b4a86611929565b60005b82811015611b7257848901358255600182019150602085019450602081019050611b4d565b86831015611b8f5784890135611b8b601f891682611a9b565b8355505b6001600288020188555050505b50505050505050565b60008160601b9050919050565b6000611bbd82611ba5565b9050919050565b6000611bcf82611bb2565b9050919050565b611be7611be282611293565b611bc4565b82525050565b600081519050919050565b600081905092915050565b6000611c0e82611bed565b611c188185611bf8565b9350611c2881856020860161105b565b80840191505092915050565b6000611c408286611bd6565b601482019150611c508285611bd6565b601482019150611c608284611c03565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611cd4602d83611c6d565b9150611cdf82611c78565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611d46602483611c6d565b9150611d5182611cea565b602482019050919050565b6000611d678261103f565b611d718185611c6d565b9350611d8181856020860161105b565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611de9602483611c6d565b9150611df482611d8d565b602482019050919050565b6000611e0a82611cc7565b9150611e1582611d39565b9150611e218285611d5c565b9150611e2c82611ddc565b9150611e388284611d5c565b91508190509392505050565b611e4d8261103f565b67ffffffffffffffff811115611e6657611e65611623565b5b611e7082546118ed565b611e7b828285611a48565b600060209050601f831160018114611eae5760008415611e9c578287015190505b611ea68582611ab9565b865550611f0e565b601f198416611ebc86611929565b60005b82811015611ee457848901518255600182019150602085019450602081019050611ebf565b86831015611f015784890151611efd601f891682611a9b565b8355505b6001600288020188555050505b505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611f4c601a83611c6d565b9150611f5782611f16565b601a82019050919050565b6000611f6d82611f3f565b9150611f798285611d5c565b9150611f858284611c03565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611fcb826111a4565b9150611fd6836111a4565b925082821015611fe957611fe8611f91565b5b828203905092915050565b6000611fff826111a4565b915061200a836111a4565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561204357612042611f91565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b60018511156120a55780860481111561208157612080611f91565b5b60018516156120905780820291505b808102905061209e8561204e565b9450612065565b94509492505050565b6000826120be576001905061217a565b816120cc576000905061217a565b81600181146120e257600281146120ec5761211b565b600191505061217a565b60ff8411156120fe576120fd611f91565b5b8360020a91508482111561211557612114611f91565b5b5061217a565b5060208310610133831016604e8410600b84101617156121505782820a90508381111561214b5761214a611f91565b5b61217a565b61215d848484600161205b565b9250905081840481111561217457612173611f91565b5b81810290505b9392505050565b600061218c826111a4565b9150612197836111a4565b92506121c47fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846120ae565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000612206826111a4565b9150612211836111a4565b925082612221576122206121cc565b5b828204905092915050565b600060ff82169050919050565b60006122448261222c565b915061224f8361222c565b92508261225f5761225e6121cc565b5b828204905092915050565b60006122758261222c565b91506122808361222c565b92508160ff048311821515161561229a57612299611f91565b5b828202905092915050565b60006122b08261222c565b91506122bb8361222c565b9250828210156122ce576122cd611f91565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000612313826111a4565b915061231e836111a4565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561235357612352611f91565b5b828201905092915050565b6000612369826111a4565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361239b5761239a611f91565b5b600182019050919050565b60006123b1826111a4565b91506123bc836111a4565b9250826123cc576123cb6121cc565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b600061243c601883611207565b915061244782612406565b602082019050919050565b6000602082019050818103600083015261246b8161242f565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b60006124a8601f83611207565b91506124b382612472565b602082019050919050565b600060208201905081810360008301526124d78161249b565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061253a602283611207565b9150612545826124de565b604082019050919050565b600060208201905081810360008301526125698161252d565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006125cc602283611207565b91506125d782612570565b604082019050919050565b600060208201905081810360008301526125fb816125bf565b9050919050565b600061260d8261222c565b91506126188361222c565b92508260ff0382111561262e5761262d611f91565b5b828201905092915050565b6126428161222c565b82525050565b600060808201905061265d6000830187611443565b61266a6020830186612639565b6126776040830185611443565b6126846060830184611443565b9594505050505056fea2646970667358221220cfda5ff98125bb45192158a173e482832639af2d11da1796f8a4a23b7be21f2464736f6c634300080f0033",
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

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
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
