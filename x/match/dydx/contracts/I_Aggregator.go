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

// IAggregatorMetaData contains all meta data concerning the IAggregator contract.
var IAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use IAggregatorMetaData.ABI instead.
var IAggregatorABI = IAggregatorMetaData.ABI

// IAggregator is an auto generated Go binding around an Ethereum contract.
type IAggregator struct {
	IAggregatorCaller     // Read-only binding to the contract
	IAggregatorTransactor // Write-only binding to the contract
	IAggregatorFilterer   // Log filterer for contract events
}

// IAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type IAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAggregatorSession struct {
	Contract     *IAggregator      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAggregatorCallerSession struct {
	Contract *IAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAggregatorTransactorSession struct {
	Contract     *IAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAggregatorRaw struct {
	Contract *IAggregator // Generic contract binding to access the raw methods on
}

// IAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAggregatorCallerRaw struct {
	Contract *IAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// IAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAggregatorTransactorRaw struct {
	Contract *IAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAggregator creates a new instance of IAggregator, bound to a specific deployed contract.
func NewIAggregator(address common.Address, backend bind.ContractBackend) (*IAggregator, error) {
	contract, err := bindIAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAggregator{IAggregatorCaller: IAggregatorCaller{contract: contract}, IAggregatorTransactor: IAggregatorTransactor{contract: contract}, IAggregatorFilterer: IAggregatorFilterer{contract: contract}}, nil
}

// NewIAggregatorCaller creates a new read-only instance of IAggregator, bound to a specific deployed contract.
func NewIAggregatorCaller(address common.Address, caller bind.ContractCaller) (*IAggregatorCaller, error) {
	contract, err := bindIAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorCaller{contract: contract}, nil
}

// NewIAggregatorTransactor creates a new write-only instance of IAggregator, bound to a specific deployed contract.
func NewIAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*IAggregatorTransactor, error) {
	contract, err := bindIAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAggregatorTransactor{contract: contract}, nil
}

// NewIAggregatorFilterer creates a new log filterer instance of IAggregator, bound to a specific deployed contract.
func NewIAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*IAggregatorFilterer, error) {
	contract, err := bindIAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAggregatorFilterer{contract: contract}, nil
}

// bindIAggregator binds a generic wrapper to an already deployed contract.
func bindIAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAggregator *IAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregator.Contract.IAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAggregator *IAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregator.Contract.IAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAggregator *IAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregator.Contract.IAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAggregator *IAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAggregator *IAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAggregator *IAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAggregator.Contract.contract.Transact(opts, method, params...)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_IAggregator *IAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_IAggregator *IAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _IAggregator.Contract.LatestAnswer(&_IAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_IAggregator *IAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _IAggregator.Contract.LatestAnswer(&_IAggregator.CallOpts)
}
