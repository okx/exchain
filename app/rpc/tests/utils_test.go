package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/dependence/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/app/rpc/types"
	"github.com/stretchr/testify/require"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"strings"
	"testing"
)

const (
	// keys which are provided on node (from test.sh)
	mnemo1                = "plunge silk glide glass curve cycle snack garbage obscure express decade dirt"
	mnemo2                = "lazy cupboard wealth canoe pumpkin gasp play dash antenna monitor material village"
	defaultPassWd         = "12345678"
	defaultCoinType       = 60
	erc20ContractByteCode = "0x60806040526012600260006101000a81548160ff021916908360ff1602179055503480156200002d57600080fd5b506040516200129738038062001297833981018060405281019080805190602001909291908051820192919060200180518201929190505050600260009054906101000a900460ff1660ff16600a0a8302600381905550600354600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160009080519060200190620000e292919062000105565b508060019080519060200190620000fb92919062000105565b50505050620001b4565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200014857805160ff191683800117855562000179565b8280016001018555821562000179579182015b82811115620001785782518255916020019190600101906200015b565b5b5090506200018891906200018c565b5090565b620001b191905b80821115620001ad57600081600090555060010162000193565b5090565b90565b6110d380620001c46000396000f3006080604052600436106100ba576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100bf578063095ea7b31461014f57806318160ddd146101b457806323b872dd146101df578063313ce5671461026457806342966c681461029557806370a08231146102da57806379cc67901461033157806395d89b4114610396578063a9059cbb14610426578063cae9ca5114610473578063dd62ed3e1461051e575b600080fd5b3480156100cb57600080fd5b506100d4610595565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101145780820151818401526020810190506100f9565b50505050905090810190601f1680156101415780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015b57600080fd5b5061019a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610633565b604051808215151515815260200191505060405180910390f35b3480156101c057600080fd5b506101c96106c0565b6040518082815260200191505060405180910390f35b3480156101eb57600080fd5b5061024a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506106c6565b604051808215151515815260200191505060405180910390f35b34801561027057600080fd5b506102796107f3565b604051808260ff1660ff16815260200191505060405180910390f35b3480156102a157600080fd5b506102c060048036038101908080359060200190929190505050610806565b604051808215151515815260200191505060405180910390f35b3480156102e657600080fd5b5061031b600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061090a565b6040518082815260200191505060405180910390f35b34801561033d57600080fd5b5061037c600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610922565b604051808215151515815260200191505060405180910390f35b3480156103a257600080fd5b506103ab610b3c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103eb5780820151818401526020810190506103d0565b50505050905090810190601f1680156104185780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561043257600080fd5b50610471600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610bda565b005b34801561047f57600080fd5b50610504600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610be9565b604051808215151515815260200191505060405180910390f35b34801561052a57600080fd5b5061057f600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d6c565b6040518082815260200191505060405180910390f35b60008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561062b5780601f106106005761010080835404028352916020019161062b565b820191906000526020600020905b81548152906001019060200180831161060e57829003601f168201915b505050505081565b600081600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506001905092915050565b60035481565b6000600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821115151561075357600080fd5b81600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055506107e8848484610d91565b600190509392505050565b600260009054906101000a900460ff1681565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561085657600080fd5b81600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540392505081905550816003600082825403925050819055503373ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040518082815260200191505060405180910390a260019050919050565b60046020528060005260406000206000915090505481565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561097257600080fd5b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482111515156109fd57600080fd5b81600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540392505081905550816003600082825403925050819055508273ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040518082815260200191505060405180910390a26001905092915050565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610bd25780601f10610ba757610100808354040283529160200191610bd2565b820191906000526020600020905b815481529060010190602001808311610bb557829003601f168201915b505050505081565b610be5338383610d91565b5050565b600080849050610bf98585610633565b15610d63578073ffffffffffffffffffffffffffffffffffffffff16638f4ffcb1338630876040518563ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610cf3578082015181840152602081019050610cd8565b50505050905090810190601f168015610d205780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b158015610d4257600080fd5b505af1158015610d56573d6000803e3d6000fd5b5050505060019150610d64565b5b509392505050565b6005602052816000526040600020602052806000526040600020600091509150505481565b6000808373ffffffffffffffffffffffffffffffffffffffff1614151515610db857600080fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410151515610e0657600080fd5b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205401111515610e9457600080fd5b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205401905081600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a380600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600460008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054011415156110a157fe5b505050505600a165627a7a72305820ed94dd1ff19d5d05f76d2df0d1cb9002bb293a6fbb55f287f36aff57fba1b0420029"
	// initial supply: 1000000000
	// token name: btc
	// token symbol: btc
	erc20ContractDeployedByteCode = "0x60806040526012600260006101000a81548160ff021916908360ff1602179055503480156200002d57600080fd5b506040516200129738038062001297833981018060405281019080805190602001909291908051820192919060200180518201929190505050600260009054906101000a900460ff1660ff16600a0a8302600381905550600354600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160009080519060200190620000e292919062000105565b508060019080519060200190620000fb92919062000105565b50505050620001b4565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200014857805160ff191683800117855562000179565b8280016001018555821562000179579182015b82811115620001785782518255916020019190600101906200015b565b5b5090506200018891906200018c565b5090565b620001b191905b80821115620001ad57600081600090555060010162000193565b5090565b90565b6110d380620001c46000396000f3006080604052600436106100ba576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100bf578063095ea7b31461014f57806318160ddd146101b457806323b872dd146101df578063313ce5671461026457806342966c681461029557806370a08231146102da57806379cc67901461033157806395d89b4114610396578063a9059cbb14610426578063cae9ca5114610473578063dd62ed3e1461051e575b600080fd5b3480156100cb57600080fd5b506100d4610595565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101145780820151818401526020810190506100f9565b50505050905090810190601f1680156101415780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015b57600080fd5b5061019a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610633565b604051808215151515815260200191505060405180910390f35b3480156101c057600080fd5b506101c96106c0565b6040518082815260200191505060405180910390f35b3480156101eb57600080fd5b5061024a600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506106c6565b604051808215151515815260200191505060405180910390f35b34801561027057600080fd5b506102796107f3565b604051808260ff1660ff16815260200191505060405180910390f35b3480156102a157600080fd5b506102c060048036038101908080359060200190929190505050610806565b604051808215151515815260200191505060405180910390f35b3480156102e657600080fd5b5061031b600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061090a565b6040518082815260200191505060405180910390f35b34801561033d57600080fd5b5061037c600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610922565b604051808215151515815260200191505060405180910390f35b3480156103a257600080fd5b506103ab610b3c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103eb5780820151818401526020810190506103d0565b50505050905090810190601f1680156104185780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561043257600080fd5b50610471600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610bda565b005b34801561047f57600080fd5b50610504600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610be9565b604051808215151515815260200191505060405180910390f35b34801561052a57600080fd5b5061057f600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d6c565b6040518082815260200191505060405180910390f35b60008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561062b5780601f106106005761010080835404028352916020019161062b565b820191906000526020600020905b81548152906001019060200180831161060e57829003601f168201915b505050505081565b600081600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506001905092915050565b60035481565b6000600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821115151561075357600080fd5b81600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825403925050819055506107e8848484610d91565b600190509392505050565b600260009054906101000a900460ff1681565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561085657600080fd5b81600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540392505081905550816003600082825403925050819055503373ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040518082815260200191505060405180910390a260019050919050565b60046020528060005260406000206000915090505481565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015151561097257600080fd5b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482111515156109fd57600080fd5b81600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540392505081905550816003600082825403925050819055508273ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040518082815260200191505060405180910390a26001905092915050565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610bd25780601f10610ba757610100808354040283529160200191610bd2565b820191906000526020600020905b815481529060010190602001808311610bb557829003601f168201915b505050505081565b610be5338383610d91565b5050565b600080849050610bf98585610633565b15610d63578073ffffffffffffffffffffffffffffffffffffffff16638f4ffcb1338630876040518563ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610cf3578082015181840152602081019050610cd8565b50505050905090810190601f168015610d205780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b158015610d4257600080fd5b505af1158015610d56573d6000803e3d6000fd5b5050505060019150610d64565b5b509392505050565b6005602052816000526040600020602052806000526040600020600091509150505481565b6000808373ffffffffffffffffffffffffffffffffffffffff1614151515610db857600080fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410151515610e0657600080fd5b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205401111515610e9457600080fd5b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205401905081600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a380600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600460008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054011415156110a157fe5b505050505600a165627a7a72305820ed94dd1ff19d5d05f76d2df0d1cb9002bb293a6fbb55f287f36aff57fba1b0420029000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000003627463000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000036274630000000000000000000000000000000000000000000000000000000000"
	// a test contract that emits an event in the constructor
	// pragma solidity ^0.5.1;

	// contract Test {
	//     event Hello(uint256 indexed world);
	//     event TestEvent(uint256 indexed a, uint256 indexed b);

	//     uint256 myStorage;

	//     constructor() public {
	//         emit Hello(1024);
	//     }

	//     function test(uint256 a, uint256 b) public {
	//         myStorage = a;
	//         emit TestEvent(a, b);
	//     }
	// }

	testContractByteCode = "0x608060405234801561001057600080fd5b506104007f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a260d08061004e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063eb8ac92114602d575b600080fd5b606060048036036040811015604157600080fd5b8101908080359060200190929190803590602001909291905050506062565b005b8160008190555080827ff3ca124a697ba07e8c5e80bebcfcc48991fc16a63170e8a9206e30508960d00360405160405180910390a3505056fea265627a7a723158201f1cebf98f21889702248e15b5424244fd05e58ea06ebd542f2d13cb38c23dee64736f6c63430005110032"

	// hash of Hello event
	helloTopic = "0x775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd738898"

	// world parameter in Hello event: 1024
	worldTopic = "0x0000000000000000000000000000000000000000000000000000000000000400"

	// Test Token contract
	ttokenContractByteCode = "0x608060405234801561001057600080fd5b5060405162000b6438038062000b648339818101604052608081101561003557600080fd5b810190808051604051939291908464010000000082111561005557600080fd5b90830190602082018581111561006a57600080fd5b825164010000000081118282018810171561008457600080fd5b82525081516020918201929091019080838360005b838110156100b1578181015183820152602001610099565b50505050905090810190601f1680156100de5780820380516001836020036101000a031916815260200191505b506040526020018051604051939291908464010000000082111561010157600080fd5b90830190602082018581111561011657600080fd5b825164010000000081118282018810171561013057600080fd5b82525081516020918201929091019080838360005b8381101561015d578181015183820152602001610145565b50505050905090810190601f16801561018a5780820380516001836020036101000a031916815260200191505b5060409081526020828101519290910151865192945092506101b191600091870190610201565b5082516101c5906001906020860190610201565b50600280546001600160a01b0390921661010002610100600160a81b031960ff90941660ff1990931692909217929092161790555061029c9050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061024257805160ff191683800117855561026f565b8280016001018555821561026f579182015b8281111561026f578251825591602001919060010190610254565b5061027b92915061027f565b5090565b61029991905b8082111561027b5760008155600101610285565b90565b6108b880620002ac6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806340c10f191161007157806340c10f19146101d957806342966c681461020557806370a082311461022257806395d89b4114610248578063a9059cbb14610250578063dd62ed3e1461027c576100a9565b806306fdde03146100ae578063095ea7b31461012b57806318160ddd1461016b57806323b872dd14610185578063313ce567146101bb575b600080fd5b6100b66102aa565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f05781810151838201526020016100d8565b50505050905090810190601f16801561011d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101576004803603604081101561014157600080fd5b506001600160a01b038135169060200135610340565b604080519115158252519081900360200190f35b6101736103a7565b60408051918252519081900360200190f35b6101576004803603606081101561019b57600080fd5b506001600160a01b038135811691602081013590911690604001356103ad565b6101c3610519565b6040805160ff9092168252519081900360200190f35b610157600480360360408110156101ef57600080fd5b506001600160a01b038135169060200135610522565b6101576004803603602081101561021b57600080fd5b5035610537565b6101736004803603602081101561023857600080fd5b50356001600160a01b031661060c565b6100b6610627565b6101576004803603604081101561026657600080fd5b506001600160a01b038135169060200135610687565b6101736004803603604081101561029257600080fd5b506001600160a01b0381358116916020013516610694565b60008054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103365780601f1061030b57610100808354040283529160200191610336565b820191906000526020600020905b81548152906001019060200180831161031957829003601f168201915b5050505050905090565b3360008181526005602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a35060015b92915050565b60035490565b6000336001600160a01b03851614806103e957506001600160a01b03841660009081526005602090815260408083203384529091529020548211155b610432576040805162461bcd60e51b815260206004820152601560248201527422a9292fa12a27a5a2a72fa120a22fa1a0a62622a960591b604482015290519081900360640190fd5b61043d8484846106bf565b336001600160a01b0385161480159061047b57506001600160a01b038416600090815260056020908152604080832033845290915290205460001914155b1561050f576001600160a01b03841660009081526005602090815260408083203384529091529020546104ae90836107d1565b6001600160a01b03858116600090815260056020908152604080832033808552908352928190208590558051948552519287169391927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a35b5060019392505050565b60025460ff1690565b600061052e83836107e1565b50600192915050565b30600090815260046020526040812054821115610592576040805162461bcd60e51b815260206004820152601460248201527311549497d25394d551919250d251539517d0905360621b604482015290519081900360640190fd5b306000908152600460205260409020546105ac90836107d1565b306000908152600460205260409020556003546105c990836107d1565b60035560408051838152905160009130917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a3506001919050565b6001600160a01b031660009081526004602052604090205490565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156103365780601f1061030b57610100808354040283529160200191610336565b600061052e3384846106bf565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205490565b6001600160a01b038316600090815260046020526040902054811115610723576040805162461bcd60e51b815260206004820152601460248201527311549497d25394d551919250d251539517d0905360621b604482015290519081900360640190fd5b6001600160a01b03831660009081526004602052604090205461074690826107d1565b6001600160a01b0380851660009081526004602052604080822093909355908416815220546107759082610872565b6001600160a01b0380841660008181526004602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b808203828111156103a157600080fd5b6001600160a01b0382166000908152600460205260409020546108049082610872565b6001600160a01b03831660009081526004602052604090205560035461082a9082610872565b6003556040805182815290516001600160a01b038416916000917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9181900360200190a35050565b808201828110156103a157600080fdfea2646970667358221220a3d106c9f53703c13c9efe61ba3fdde06510efec7cce1c90fb4f376e677c7d5364736f6c63430006000033000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000120000000000000000000000002cf4ea7df75b513509d95946b43062e26bd8803500000000000000000000000000000000000000000000000000000000000000036f6b74000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000036f6b740000000000000000000000000000000000000000000000000000000000"
	approveFuncHash       = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	transferFuncHash      = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	recvAddr1             = "0xc9c9b43322f5e1dc401252076fa4e699c9122cd6"
	recvAddr2             = "0x9ad84c8630e0282f78e5479b46e64e17779e3cfb"
	recvAddr1Hash         = "0x000000000000000000000000c9c9b43322f5e1dc401252076fa4e699c9122cd6"
	recvAddr2Hash         = "0x0000000000000000000000009ad84c8630e0282f78e5479b46e64e17779e3cfb"
)

const (
	testContractKind = iota
	erc20ContractKind
)

var (
	keyInfo1, keyInfo2 keys.Info
	Kb                 = keys.NewInMemory(hd.EthSecp256k1Options()...)
	hexAddr1, hexAddr2 ethcmn.Address
	addrCounter        = 2
	defaultGasPrice    sdk.SysCoin
)

func init() {
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	config := sdk.GetConfig()
	config.SetCoinType(defaultCoinType)

	keyInfo1, _ = createAccountWithMnemo(mnemo1, "alice", defaultPassWd)
	keyInfo2, _ = createAccountWithMnemo(mnemo2, "bob", defaultPassWd)
	hexAddr1 = ethcmn.BytesToAddress(keyInfo1.GetAddress().Bytes())
	hexAddr2 = ethcmn.BytesToAddress(keyInfo2.GetAddress().Bytes())
	defaultGasPrice, _ = sdk.ParseDecCoin(defaultMinGasPrice)
}

func TestGetAddress(t *testing.T) {
	addr, err := GetAddress()
	require.NoError(t, err)
	require.True(t, bytes.Equal(addr, hexAddr1[:]))
}

func createAccountWithMnemo(mnemonic, name, passWd string) (info keys.Info, err error) {
	hdPath := keys.CreateHDPath(0, 0).String()
	info, err = Kb.CreateAccount(name, mnemonic, "", passWd, hdPath, hd.EthSecp256k1)
	if err != nil {
		return info, fmt.Errorf("failed. Kb.CreateAccount err : %s", err.Error())
	}

	return info, err
}

// sendTestTransaction sends a dummy transaction
func sendTestTransaction(t *testing.T, senderAddr, receiverAddr ethcmn.Address, value uint) ethcmn.Hash {
	fromAddrStr, toAddrStr := senderAddr.Hex(), receiverAddr.Hex()
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["to"] = toAddrStr
	param[0]["value"] = hexutil.Uint(value).String()
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	t.Logf("%s transfers %d to %s successfully\n", fromAddrStr, value, toAddrStr)
	return hash
}

// deployTestContract deploys a contract that emits an event in the constructor
func deployTestContract(t *testing.T, senderAddr ethcmn.Address, kind int) (ethcmn.Hash, map[string]interface{}) {
	fromAddrStr := senderAddr.Hex()
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["gasprice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	switch kind {
	case testContractKind:
		param[0]["data"] = testContractByteCode
	case erc20ContractKind:
		param[0]["data"] = erc20ContractDeployedByteCode
	default:
		panic("unsupported contract kind")
	}

	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))

	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt, "transaction failed")
	require.Equal(t, "0x1", receipt["status"].(string))
	t.Logf("%s has deployed a contract %s with tx hash %s successfully\n", fromAddrStr, receipt["contractAddress"], hash.Hex())
	return hash, receipt
}

func getBlockHeightFromTxHash(t *testing.T, hash ethcmn.Hash) hexutil.Uint64 {
	rpcRes := Call(t, "eth_getTransactionByHash", []interface{}{hash})
	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))

	if transaction.BlockNumber == nil {
		return hexutil.Uint64(0)
	}

	return hexutil.Uint64(transaction.BlockNumber.ToInt().Uint64())
}

func getBlockHashFromTxHash(t *testing.T, hash ethcmn.Hash) *ethcmn.Hash {
	rpcRes := Call(t, "eth_getTransactionByHash", []interface{}{hash})
	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))

	return transaction.BlockHash
}

func assertNullFromJSONResponse(t *testing.T, jrm json.RawMessage) {
	require.True(t, bytes.Equal([]byte("null"), jrm))
}

func signWithAccNameAndPasswd(accName, passWd string, msg []byte) (sig []byte, err error) {
	privKey, err := Kb.ExportPrivateKeyObject(accName, passWd)
	if err != nil {
		return
	}

	ethPrivKey, ok := privKey.(ethsecp256k1.PrivKey)
	if !ok {
		return sig, fmt.Errorf("invalid private key type %T", privKey)
	}

	sig, err = crypto.Sign(accounts.TextHash(msg), ethPrivKey.ToECDSA())
	if err != nil {
		return nil, err
	}
	sig[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper

	return
}

func hexToBloom(t *testing.T, hexStr string) ethtypes.Bloom {
	bz, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	require.NoError(t, err)

	return ethtypes.BytesToBloom(bz)
}
