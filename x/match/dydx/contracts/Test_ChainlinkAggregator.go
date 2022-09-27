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

// TestChainlinkAggregatorMetaData contains all meta data concerning the TestChainlinkAggregator contract.
var TestChainlinkAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"_ANSWER_\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"int256\",\"name\":\"newAnswer\",\"type\":\"int256\"}],\"name\":\"setAnswer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TestChainlinkAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use TestChainlinkAggregatorMetaData.ABI instead.
var TestChainlinkAggregatorABI = TestChainlinkAggregatorMetaData.ABI

// TestChainlinkAggregator is an auto generated Go binding around an Ethereum contract.
type TestChainlinkAggregator struct {
	TestChainlinkAggregatorCaller     // Read-only binding to the contract
	TestChainlinkAggregatorTransactor // Write-only binding to the contract
	TestChainlinkAggregatorFilterer   // Log filterer for contract events
}

// TestChainlinkAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestChainlinkAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestChainlinkAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestChainlinkAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestChainlinkAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestChainlinkAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestChainlinkAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestChainlinkAggregatorSession struct {
	Contract     *TestChainlinkAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestChainlinkAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestChainlinkAggregatorCallerSession struct {
	Contract *TestChainlinkAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// TestChainlinkAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestChainlinkAggregatorTransactorSession struct {
	Contract     *TestChainlinkAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// TestChainlinkAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestChainlinkAggregatorRaw struct {
	Contract *TestChainlinkAggregator // Generic contract binding to access the raw methods on
}

// TestChainlinkAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestChainlinkAggregatorCallerRaw struct {
	Contract *TestChainlinkAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// TestChainlinkAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestChainlinkAggregatorTransactorRaw struct {
	Contract *TestChainlinkAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestChainlinkAggregator creates a new instance of TestChainlinkAggregator, bound to a specific deployed contract.
func NewTestChainlinkAggregator(address common.Address, backend bind.ContractBackend) (*TestChainlinkAggregator, error) {
	contract, err := bindTestChainlinkAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestChainlinkAggregator{TestChainlinkAggregatorCaller: TestChainlinkAggregatorCaller{contract: contract}, TestChainlinkAggregatorTransactor: TestChainlinkAggregatorTransactor{contract: contract}, TestChainlinkAggregatorFilterer: TestChainlinkAggregatorFilterer{contract: contract}}, nil
}

// NewTestChainlinkAggregatorCaller creates a new read-only instance of TestChainlinkAggregator, bound to a specific deployed contract.
func NewTestChainlinkAggregatorCaller(address common.Address, caller bind.ContractCaller) (*TestChainlinkAggregatorCaller, error) {
	contract, err := bindTestChainlinkAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestChainlinkAggregatorCaller{contract: contract}, nil
}

// NewTestChainlinkAggregatorTransactor creates a new write-only instance of TestChainlinkAggregator, bound to a specific deployed contract.
func NewTestChainlinkAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestChainlinkAggregatorTransactor, error) {
	contract, err := bindTestChainlinkAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestChainlinkAggregatorTransactor{contract: contract}, nil
}

// NewTestChainlinkAggregatorFilterer creates a new log filterer instance of TestChainlinkAggregator, bound to a specific deployed contract.
func NewTestChainlinkAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestChainlinkAggregatorFilterer, error) {
	contract, err := bindTestChainlinkAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestChainlinkAggregatorFilterer{contract: contract}, nil
}

// bindTestChainlinkAggregator binds a generic wrapper to an already deployed contract.
func bindTestChainlinkAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestChainlinkAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestChainlinkAggregator *TestChainlinkAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestChainlinkAggregator.Contract.TestChainlinkAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestChainlinkAggregator *TestChainlinkAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.TestChainlinkAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestChainlinkAggregator *TestChainlinkAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.TestChainlinkAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestChainlinkAggregator *TestChainlinkAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestChainlinkAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestChainlinkAggregator *TestChainlinkAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestChainlinkAggregator *TestChainlinkAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.contract.Transact(opts, method, params...)
}

// ANSWER is a free data retrieval call binding the contract method 0x3ad12a22.
//
// Solidity: function _ANSWER_() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorCaller) ANSWER(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestChainlinkAggregator.contract.Call(opts, &out, "_ANSWER_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ANSWER is a free data retrieval call binding the contract method 0x3ad12a22.
//
// Solidity: function _ANSWER_() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorSession) ANSWER() (*big.Int, error) {
	return _TestChainlinkAggregator.Contract.ANSWER(&_TestChainlinkAggregator.CallOpts)
}

// ANSWER is a free data retrieval call binding the contract method 0x3ad12a22.
//
// Solidity: function _ANSWER_() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorCallerSession) ANSWER() (*big.Int, error) {
	return _TestChainlinkAggregator.Contract.ANSWER(&_TestChainlinkAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestChainlinkAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _TestChainlinkAggregator.Contract.LatestAnswer(&_TestChainlinkAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_TestChainlinkAggregator *TestChainlinkAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _TestChainlinkAggregator.Contract.LatestAnswer(&_TestChainlinkAggregator.CallOpts)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 newAnswer) returns()
func (_TestChainlinkAggregator *TestChainlinkAggregatorTransactor) SetAnswer(opts *bind.TransactOpts, newAnswer *big.Int) (*types.Transaction, error) {
	return _TestChainlinkAggregator.contract.Transact(opts, "setAnswer", newAnswer)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 newAnswer) returns()
func (_TestChainlinkAggregator *TestChainlinkAggregatorSession) SetAnswer(newAnswer *big.Int) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.SetAnswer(&_TestChainlinkAggregator.TransactOpts, newAnswer)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 newAnswer) returns()
func (_TestChainlinkAggregator *TestChainlinkAggregatorTransactorSession) SetAnswer(newAnswer *big.Int) (*types.Transaction, error) {
	return _TestChainlinkAggregator.Contract.SetAnswer(&_TestChainlinkAggregator.TransactOpts, newAnswer)
}
