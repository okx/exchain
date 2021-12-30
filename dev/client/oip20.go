// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// Oip20MetaData contains all meta data concerning the Oip20 contract.
var Oip20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"feeReceiver\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"_mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260405162001e6938038062001e6983398181016040528101906200002991906200064e565b856001908051906020019062000041929190620002de565b5084600090805190602001906200005a929190620002de565b5083600260006101000a81548160ff021916908360ff160217905550826003819055506200008f8284620000e360201b60201c565b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f19350505050158015620000d6573d6000803e3d6000fd5b505050505050506200093c565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141562000156576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200014d9062000789565b60405180910390fd5b62000172816003546200028060201b620008651790919060201c565b600381905550620001d181600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546200028060201b620008651790919060201c565b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620002749190620007bc565b60405180910390a35050565b600082828462000291919062000808565b9150811015620002d8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620002cf90620008b5565b60405180910390fd5b92915050565b828054620002ec9062000906565b90600052602060002090601f0160209004810192826200031057600085556200035c565b82601f106200032b57805160ff19168380011785556200035c565b828001600101855582156200035c579182015b828111156200035b5782518255916020019190600101906200033e565b5b5090506200036b91906200036f565b5090565b5b808211156200038a57600081600090555060010162000370565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620003f782620003ac565b810181811067ffffffffffffffff82111715620004195762000418620003bd565b5b80604052505050565b60006200042e6200038e565b90506200043c8282620003ec565b919050565b600067ffffffffffffffff8211156200045f576200045e620003bd565b5b6200046a82620003ac565b9050602081019050919050565b60005b83811015620004975780820151818401526020810190506200047a565b83811115620004a7576000848401525b50505050565b6000620004c4620004be8462000441565b62000422565b905082815260208101848484011115620004e357620004e2620003a7565b5b620004f084828562000477565b509392505050565b600082601f83011262000510576200050f620003a2565b5b815162000522848260208601620004ad565b91505092915050565b600060ff82169050919050565b62000543816200052b565b81146200054f57600080fd5b50565b600081519050620005638162000538565b92915050565b6000819050919050565b6200057e8162000569565b81146200058a57600080fd5b50565b6000815190506200059e8162000573565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620005d182620005a4565b9050919050565b620005e381620005c4565b8114620005ef57600080fd5b50565b6000815190506200060381620005d8565b92915050565b60006200061682620005a4565b9050919050565b620006288162000609565b81146200063457600080fd5b50565b60008151905062000648816200061d565b92915050565b60008060008060008060c087890312156200066e576200066d62000398565b5b600087015167ffffffffffffffff8111156200068f576200068e6200039d565b5b6200069d89828a01620004f8565b965050602087015167ffffffffffffffff811115620006c157620006c06200039d565b5b620006cf89828a01620004f8565b9550506040620006e289828a0162000552565b9450506060620006f589828a016200058d565b93505060806200070889828a01620005f2565b92505060a06200071b89828a0162000637565b9150509295509295509295565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000771601f8362000728565b91506200077e8262000739565b602082019050919050565b60006020820190508181036000830152620007a48162000762565b9050919050565b620007b68162000569565b82525050565b6000602082019050620007d36000830184620007ab565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620008158262000569565b9150620008228362000569565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156200085a5762000859620007d9565b5b828201905092915050565b7f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000600082015250565b60006200089d60148362000728565b9150620008aa8262000865565b602082019050919050565b60006020820190508181036000830152620008d0816200088e565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200091f57607f821691505b60208210811415620009365762000935620008d7565b5b50919050565b61151d806200094c6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80634e6ec247116100715780634e6ec247146101a357806370a08231146101bf57806395d89b41146101ef578063a457c2d71461020d578063a9059cbb1461023d578063dd62ed3e1461026d576100b4565b806306fdde03146100b9578063095ea7b3146100d757806318160ddd1461010757806323b872dd14610125578063313ce567146101555780633950935114610173575b600080fd5b6100c161029d565b6040516100ce9190610def565b60405180910390f35b6100f160048036038101906100ec9190610eaa565b61032f565b6040516100fe9190610f05565b60405180910390f35b61010f610346565b60405161011c9190610f2f565b60405180910390f35b61013f600480360381019061013a9190610f4a565b610350565b60405161014c9190610f05565b60405180910390f35b61015d610401565b60405161016a9190610fb9565b60405180910390f35b61018d60048036038101906101889190610eaa565b610418565b60405161019a9190610f05565b60405180910390f35b6101bd60048036038101906101b89190610eaa565b6104bd565b005b6101d960048036038101906101d49190610fd4565b610647565b6040516101e69190610f2f565b60405180910390f35b6101f7610690565b6040516102049190610def565b60405180910390f35b61022760048036038101906102229190610eaa565b610722565b6040516102349190610f05565b60405180910390f35b61025760048036038101906102529190610eaa565b6107c7565b6040516102649190610f05565b60405180910390f35b61028760048036038101906102829190611001565b6107de565b6040516102949190610f2f565b60405180910390f35b6060600080546102ac90611070565b80601f01602080910402602001604051908101604052809291908181526020018280546102d890611070565b80156103255780601f106102fa57610100808354040283529160200191610325565b820191906000526020600020905b81548152906001019060200180831161030857829003601f168201915b5050505050905090565b600061033c3384846108be565b6001905092915050565b6000600354905090565b600061035d848484610a89565b6103f684336103f185600560008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfd90919063ffffffff16565b6108be565b600190509392505050565b6000600260009054906101000a900460ff16905090565b60006104b333846104ae85600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461086590919063ffffffff16565b6108be565b6001905092915050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561052d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610524906110ee565b60405180910390fd5b6105428160035461086590919063ffffffff16565b60038190555061059a81600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461086590919063ffffffff16565b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161063b9190610f2f565b60405180910390a35050565b6000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60606001805461069f90611070565b80601f01602080910402602001604051908101604052809291908181526020018280546106cb90611070565b80156107185780601f106106ed57610100808354040283529160200191610718565b820191906000526020600020905b8154815290600101906020018083116106fb57829003601f168201915b5050505050905090565b60006107bd33846107b885600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfd90919063ffffffff16565b6108be565b6001905092915050565b60006107d4338484610a89565b6001905092915050565b6000600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828284610874919061113d565b91508110156108b8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108af906111df565b60405180910390fd5b92915050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561092e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161092590611271565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561099e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161099590611303565b60405180910390fd5b80600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610a7c9190610f2f565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610af9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610af090611395565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610b69576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b6090611427565b60405180910390fd5b610bbb81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfd90919063ffffffff16565b600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610c5081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461086590919063ffffffff16565b600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610cf09190610f2f565b60405180910390a3505050565b6000828284610d0c9190611447565b9150811115610d50576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d47906114c7565b60405180910390fd5b92915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610d90578082015181840152602081019050610d75565b83811115610d9f576000848401525b50505050565b6000601f19601f8301169050919050565b6000610dc182610d56565b610dcb8185610d61565b9350610ddb818560208601610d72565b610de481610da5565b840191505092915050565b60006020820190508181036000830152610e098184610db6565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610e4182610e16565b9050919050565b610e5181610e36565b8114610e5c57600080fd5b50565b600081359050610e6e81610e48565b92915050565b6000819050919050565b610e8781610e74565b8114610e9257600080fd5b50565b600081359050610ea481610e7e565b92915050565b60008060408385031215610ec157610ec0610e11565b5b6000610ecf85828601610e5f565b9250506020610ee085828601610e95565b9150509250929050565b60008115159050919050565b610eff81610eea565b82525050565b6000602082019050610f1a6000830184610ef6565b92915050565b610f2981610e74565b82525050565b6000602082019050610f446000830184610f20565b92915050565b600080600060608486031215610f6357610f62610e11565b5b6000610f7186828701610e5f565b9350506020610f8286828701610e5f565b9250506040610f9386828701610e95565b9150509250925092565b600060ff82169050919050565b610fb381610f9d565b82525050565b6000602082019050610fce6000830184610faa565b92915050565b600060208284031215610fea57610fe9610e11565b5b6000610ff884828501610e5f565b91505092915050565b6000806040838503121561101857611017610e11565b5b600061102685828601610e5f565b925050602061103785828601610e5f565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061108857607f821691505b6020821081141561109c5761109b611041565b5b50919050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b60006110d8601f83610d61565b91506110e3826110a2565b602082019050919050565b60006020820190508181036000830152611107816110cb565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061114882610e74565b915061115383610e74565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156111885761118761110e565b5b828201905092915050565b7f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000600082015250565b60006111c9601483610d61565b91506111d482611193565b602082019050919050565b600060208201905081810360008301526111f8816111bc565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b600061125b602483610d61565b9150611266826111ff565b604082019050919050565b6000602082019050818103600083015261128a8161124e565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006112ed602283610d61565b91506112f882611291565b604082019050919050565b6000602082019050818103600083015261131c816112e0565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b600061137f602583610d61565b915061138a82611323565b604082019050919050565b600060208201905081810360008301526113ae81611372565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611411602383610d61565b915061141c826113b5565b604082019050919050565b6000602082019050818103600083015261144081611404565b9050919050565b600061145282610e74565b915061145d83610e74565b9250828210156114705761146f61110e565b5b828203905092915050565b7f64732d6d6174682d7375622d756e646572666c6f770000000000000000000000600082015250565b60006114b1601583610d61565b91506114bc8261147b565b602082019050919050565b600060208201905081810360008301526114e0816114a4565b905091905056fea2646970667358221220206ee5db59571ba825d922503b87562d8c5f4942b30b7e067e14789d9f769f9e64736f6c634300080b0033",
}

// Oip20ABI is the input ABI used to generate the binding from.
// Deprecated: Use Oip20MetaData.ABI instead.
var Oip20ABI = Oip20MetaData.ABI

// Oip20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Oip20MetaData.Bin instead.
var Oip20Bin = Oip20MetaData.Bin

// DeployOip20 deploys a new Ethereum contract, binding an instance of Oip20 to it.
func DeployOip20(auth *bind.TransactOpts, backend bind.ContractBackend, symbol string, name string, decimals uint8, totalSupply *big.Int, ownerAddress common.Address, feeReceiver common.Address) (common.Address, *types.Transaction, *Oip20, error) {
	parsed, err := Oip20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Oip20Bin), backend, symbol, name, decimals, totalSupply, ownerAddress, feeReceiver)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oip20{Oip20Caller: Oip20Caller{contract: contract}, Oip20Transactor: Oip20Transactor{contract: contract}, Oip20Filterer: Oip20Filterer{contract: contract}}, nil
}

// Oip20 is an auto generated Go binding around an Ethereum contract.
type Oip20 struct {
	Oip20Caller     // Read-only binding to the contract
	Oip20Transactor // Write-only binding to the contract
	Oip20Filterer   // Log filterer for contract events
}

// Oip20Caller is an auto generated read-only Go binding around an Ethereum contract.
type Oip20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oip20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Oip20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oip20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Oip20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Oip20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Oip20Session struct {
	Contract     *Oip20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Oip20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Oip20CallerSession struct {
	Contract *Oip20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Oip20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Oip20TransactorSession struct {
	Contract     *Oip20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Oip20Raw is an auto generated low-level Go binding around an Ethereum contract.
type Oip20Raw struct {
	Contract *Oip20 // Generic contract binding to access the raw methods on
}

// Oip20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Oip20CallerRaw struct {
	Contract *Oip20Caller // Generic read-only contract binding to access the raw methods on
}

// Oip20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Oip20TransactorRaw struct {
	Contract *Oip20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewOip20 creates a new instance of Oip20, bound to a specific deployed contract.
func NewOip20(address common.Address, backend bind.ContractBackend) (*Oip20, error) {
	contract, err := bindOip20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oip20{Oip20Caller: Oip20Caller{contract: contract}, Oip20Transactor: Oip20Transactor{contract: contract}, Oip20Filterer: Oip20Filterer{contract: contract}}, nil
}

// NewOip20Caller creates a new read-only instance of Oip20, bound to a specific deployed contract.
func NewOip20Caller(address common.Address, caller bind.ContractCaller) (*Oip20Caller, error) {
	contract, err := bindOip20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Oip20Caller{contract: contract}, nil
}

// NewOip20Transactor creates a new write-only instance of Oip20, bound to a specific deployed contract.
func NewOip20Transactor(address common.Address, transactor bind.ContractTransactor) (*Oip20Transactor, error) {
	contract, err := bindOip20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Oip20Transactor{contract: contract}, nil
}

// NewOip20Filterer creates a new log filterer instance of Oip20, bound to a specific deployed contract.
func NewOip20Filterer(address common.Address, filterer bind.ContractFilterer) (*Oip20Filterer, error) {
	contract, err := bindOip20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Oip20Filterer{contract: contract}, nil
}

// bindOip20 binds a generic wrapper to an already deployed contract.
func bindOip20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Oip20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oip20 *Oip20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oip20.Contract.Oip20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oip20 *Oip20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oip20.Contract.Oip20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oip20 *Oip20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oip20.Contract.Oip20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oip20 *Oip20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oip20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oip20 *Oip20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oip20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oip20 *Oip20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oip20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Oip20 *Oip20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Oip20 *Oip20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Oip20.Contract.Allowance(&_Oip20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Oip20 *Oip20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Oip20.Contract.Allowance(&_Oip20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Oip20 *Oip20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Oip20 *Oip20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _Oip20.Contract.BalanceOf(&_Oip20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Oip20 *Oip20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Oip20.Contract.BalanceOf(&_Oip20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Oip20 *Oip20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Oip20 *Oip20Session) Decimals() (uint8, error) {
	return _Oip20.Contract.Decimals(&_Oip20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Oip20 *Oip20CallerSession) Decimals() (uint8, error) {
	return _Oip20.Contract.Decimals(&_Oip20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Oip20 *Oip20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Oip20 *Oip20Session) Name() (string, error) {
	return _Oip20.Contract.Name(&_Oip20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Oip20 *Oip20CallerSession) Name() (string, error) {
	return _Oip20.Contract.Name(&_Oip20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Oip20 *Oip20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Oip20 *Oip20Session) Symbol() (string, error) {
	return _Oip20.Contract.Symbol(&_Oip20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Oip20 *Oip20CallerSession) Symbol() (string, error) {
	return _Oip20.Contract.Symbol(&_Oip20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Oip20 *Oip20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oip20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Oip20 *Oip20Session) TotalSupply() (*big.Int, error) {
	return _Oip20.Contract.TotalSupply(&_Oip20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Oip20 *Oip20CallerSession) TotalSupply() (*big.Int, error) {
	return _Oip20.Contract.TotalSupply(&_Oip20.CallOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x4e6ec247.
//
// Solidity: function _mint(address account, uint256 amount) returns()
func (_Oip20 *Oip20Transactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "_mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x4e6ec247.
//
// Solidity: function _mint(address account, uint256 amount) returns()
func (_Oip20 *Oip20Session) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Mint(&_Oip20.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x4e6ec247.
//
// Solidity: function _mint(address account, uint256 amount) returns()
func (_Oip20 *Oip20TransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Mint(&_Oip20.TransactOpts, account, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Oip20 *Oip20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Oip20 *Oip20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Approve(&_Oip20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Oip20 *Oip20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Approve(&_Oip20.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Oip20 *Oip20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Oip20 *Oip20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.DecreaseAllowance(&_Oip20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Oip20 *Oip20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.DecreaseAllowance(&_Oip20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Oip20 *Oip20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Oip20 *Oip20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.IncreaseAllowance(&_Oip20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Oip20 *Oip20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.IncreaseAllowance(&_Oip20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Transfer(&_Oip20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.Transfer(&_Oip20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.TransferFrom(&_Oip20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Oip20 *Oip20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Oip20.Contract.TransferFrom(&_Oip20.TransactOpts, sender, recipient, amount)
}

// Oip20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Oip20 contract.
type Oip20ApprovalIterator struct {
	Event *Oip20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Oip20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oip20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Oip20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Oip20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oip20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oip20Approval represents a Approval event raised by the Oip20 contract.
type Oip20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Oip20 *Oip20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*Oip20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Oip20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &Oip20ApprovalIterator{contract: _Oip20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Oip20 *Oip20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Oip20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Oip20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oip20Approval)
				if err := _Oip20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Oip20 *Oip20Filterer) ParseApproval(log types.Log) (*Oip20Approval, error) {
	event := new(Oip20Approval)
	if err := _Oip20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Oip20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Oip20 contract.
type Oip20TransferIterator struct {
	Event *Oip20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Oip20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Oip20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Oip20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Oip20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Oip20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Oip20Transfer represents a Transfer event raised by the Oip20 contract.
type Oip20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Oip20 *Oip20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Oip20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oip20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Oip20TransferIterator{contract: _Oip20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Oip20 *Oip20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Oip20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oip20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Oip20Transfer)
				if err := _Oip20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Oip20 *Oip20Filterer) ParseTransfer(log types.Log) (*Oip20Transfer, error) {
	event := new(Oip20Transfer)
	if err := _Oip20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
