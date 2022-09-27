// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// IExchangeWrapperMetaData contains all meta data concerning the IExchangeWrapper contract.
var IExchangeWrapperMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tradeOriginator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"orderData\",\"type\":\"bytes\"}],\"name\":\"exchange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"desiredMakerToken\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"orderData\",\"type\":\"bytes\"}],\"name\":\"getExchangeCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IExchangeWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use IExchangeWrapperMetaData.ABI instead.
var IExchangeWrapperABI = IExchangeWrapperMetaData.ABI

// IExchangeWrapper is an auto generated Go binding around an Ethereum contract.
type IExchangeWrapper struct {
	IExchangeWrapperCaller     // Read-only binding to the contract
	IExchangeWrapperTransactor // Write-only binding to the contract
	IExchangeWrapperFilterer   // Log filterer for contract events
}

// IExchangeWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type IExchangeWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IExchangeWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IExchangeWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IExchangeWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IExchangeWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IExchangeWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IExchangeWrapperSession struct {
	Contract     *IExchangeWrapper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IExchangeWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IExchangeWrapperCallerSession struct {
	Contract *IExchangeWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IExchangeWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IExchangeWrapperTransactorSession struct {
	Contract     *IExchangeWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IExchangeWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type IExchangeWrapperRaw struct {
	Contract *IExchangeWrapper // Generic contract binding to access the raw methods on
}

// IExchangeWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IExchangeWrapperCallerRaw struct {
	Contract *IExchangeWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// IExchangeWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IExchangeWrapperTransactorRaw struct {
	Contract *IExchangeWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIExchangeWrapper creates a new instance of IExchangeWrapper, bound to a specific deployed contract.
func NewIExchangeWrapper(address common.Address, backend bind.ContractBackend) (*IExchangeWrapper, error) {
	contract, err := bindIExchangeWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IExchangeWrapper{IExchangeWrapperCaller: IExchangeWrapperCaller{contract: contract}, IExchangeWrapperTransactor: IExchangeWrapperTransactor{contract: contract}, IExchangeWrapperFilterer: IExchangeWrapperFilterer{contract: contract}}, nil
}

// NewIExchangeWrapperCaller creates a new read-only instance of IExchangeWrapper, bound to a specific deployed contract.
func NewIExchangeWrapperCaller(address common.Address, caller bind.ContractCaller) (*IExchangeWrapperCaller, error) {
	contract, err := bindIExchangeWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IExchangeWrapperCaller{contract: contract}, nil
}

// NewIExchangeWrapperTransactor creates a new write-only instance of IExchangeWrapper, bound to a specific deployed contract.
func NewIExchangeWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*IExchangeWrapperTransactor, error) {
	contract, err := bindIExchangeWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IExchangeWrapperTransactor{contract: contract}, nil
}

// NewIExchangeWrapperFilterer creates a new log filterer instance of IExchangeWrapper, bound to a specific deployed contract.
func NewIExchangeWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*IExchangeWrapperFilterer, error) {
	contract, err := bindIExchangeWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IExchangeWrapperFilterer{contract: contract}, nil
}

// bindIExchangeWrapper binds a generic wrapper to an already deployed contract.
func bindIExchangeWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IExchangeWrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IExchangeWrapper *IExchangeWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IExchangeWrapper.Contract.IExchangeWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IExchangeWrapper *IExchangeWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.IExchangeWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IExchangeWrapper *IExchangeWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.IExchangeWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IExchangeWrapper *IExchangeWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IExchangeWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IExchangeWrapper *IExchangeWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IExchangeWrapper *IExchangeWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.contract.Transact(opts, method, params...)
}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address makerToken, address takerToken, uint256 desiredMakerToken, bytes orderData) view returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperCaller) GetExchangeCost(opts *bind.CallOpts, makerToken common.Address, takerToken common.Address, desiredMakerToken *big.Int, orderData []byte) (*big.Int, error) {
	var out []interface{}
	err := _IExchangeWrapper.contract.Call(opts, &out, "getExchangeCost", makerToken, takerToken, desiredMakerToken, orderData)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address makerToken, address takerToken, uint256 desiredMakerToken, bytes orderData) view returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperSession) GetExchangeCost(makerToken common.Address, takerToken common.Address, desiredMakerToken *big.Int, orderData []byte) (*big.Int, error) {
	return _IExchangeWrapper.Contract.GetExchangeCost(&_IExchangeWrapper.CallOpts, makerToken, takerToken, desiredMakerToken, orderData)
}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address makerToken, address takerToken, uint256 desiredMakerToken, bytes orderData) view returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperCallerSession) GetExchangeCost(makerToken common.Address, takerToken common.Address, desiredMakerToken *big.Int, orderData []byte) (*big.Int, error) {
	return _IExchangeWrapper.Contract.GetExchangeCost(&_IExchangeWrapper.CallOpts, makerToken, takerToken, desiredMakerToken, orderData)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address tradeOriginator, address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperTransactor) Exchange(opts *bind.TransactOpts, tradeOriginator common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _IExchangeWrapper.contract.Transact(opts, "exchange", tradeOriginator, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address tradeOriginator, address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperSession) Exchange(tradeOriginator common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.Exchange(&_IExchangeWrapper.TransactOpts, tradeOriginator, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address tradeOriginator, address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_IExchangeWrapper *IExchangeWrapperTransactorSession) Exchange(tradeOriginator common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _IExchangeWrapper.Contract.Exchange(&_IExchangeWrapper.TransactOpts, tradeOriginator, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}
