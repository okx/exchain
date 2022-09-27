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

// TestP1OracleMetaData contains all meta data concerning the TestP1Oracle contract.
var TestP1OracleMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"_PRICE_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPrice\",\"type\":\"uint256\"}],\"name\":\"setPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestP1OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use TestP1OracleMetaData.ABI instead.
var TestP1OracleABI = TestP1OracleMetaData.ABI

// TestP1Oracle is an auto generated Go binding around an Ethereum contract.
type TestP1Oracle struct {
	TestP1OracleCaller     // Read-only binding to the contract
	TestP1OracleTransactor // Write-only binding to the contract
	TestP1OracleFilterer   // Log filterer for contract events
}

// TestP1OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestP1OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestP1OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestP1OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestP1OracleSession struct {
	Contract     *TestP1Oracle     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestP1OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestP1OracleCallerSession struct {
	Contract *TestP1OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TestP1OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestP1OracleTransactorSession struct {
	Contract     *TestP1OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TestP1OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestP1OracleRaw struct {
	Contract *TestP1Oracle // Generic contract binding to access the raw methods on
}

// TestP1OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestP1OracleCallerRaw struct {
	Contract *TestP1OracleCaller // Generic read-only contract binding to access the raw methods on
}

// TestP1OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestP1OracleTransactorRaw struct {
	Contract *TestP1OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestP1Oracle creates a new instance of TestP1Oracle, bound to a specific deployed contract.
func NewTestP1Oracle(address common.Address, backend bind.ContractBackend) (*TestP1Oracle, error) {
	contract, err := bindTestP1Oracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestP1Oracle{TestP1OracleCaller: TestP1OracleCaller{contract: contract}, TestP1OracleTransactor: TestP1OracleTransactor{contract: contract}, TestP1OracleFilterer: TestP1OracleFilterer{contract: contract}}, nil
}

// NewTestP1OracleCaller creates a new read-only instance of TestP1Oracle, bound to a specific deployed contract.
func NewTestP1OracleCaller(address common.Address, caller bind.ContractCaller) (*TestP1OracleCaller, error) {
	contract, err := bindTestP1Oracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1OracleCaller{contract: contract}, nil
}

// NewTestP1OracleTransactor creates a new write-only instance of TestP1Oracle, bound to a specific deployed contract.
func NewTestP1OracleTransactor(address common.Address, transactor bind.ContractTransactor) (*TestP1OracleTransactor, error) {
	contract, err := bindTestP1Oracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1OracleTransactor{contract: contract}, nil
}

// NewTestP1OracleFilterer creates a new log filterer instance of TestP1Oracle, bound to a specific deployed contract.
func NewTestP1OracleFilterer(address common.Address, filterer bind.ContractFilterer) (*TestP1OracleFilterer, error) {
	contract, err := bindTestP1Oracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestP1OracleFilterer{contract: contract}, nil
}

// bindTestP1Oracle binds a generic wrapper to an already deployed contract.
func bindTestP1Oracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestP1OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Oracle *TestP1OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Oracle.Contract.TestP1OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Oracle *TestP1OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.TestP1OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Oracle *TestP1OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.TestP1OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Oracle *TestP1OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Oracle *TestP1OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Oracle *TestP1OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.contract.Transact(opts, method, params...)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Oracle *TestP1OracleCaller) PRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Oracle.contract.Call(opts, &out, "_PRICE_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Oracle *TestP1OracleSession) PRICE() (*big.Int, error) {
	return _TestP1Oracle.Contract.PRICE(&_TestP1Oracle.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Oracle *TestP1OracleCallerSession) PRICE() (*big.Int, error) {
	return _TestP1Oracle.Contract.PRICE(&_TestP1Oracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Oracle *TestP1OracleCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Oracle.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Oracle *TestP1OracleSession) GetPrice() (*big.Int, error) {
	return _TestP1Oracle.Contract.GetPrice(&_TestP1Oracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Oracle *TestP1OracleCallerSession) GetPrice() (*big.Int, error) {
	return _TestP1Oracle.Contract.GetPrice(&_TestP1Oracle.CallOpts)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Oracle *TestP1OracleTransactor) SetPrice(opts *bind.TransactOpts, newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Oracle.contract.Transact(opts, "setPrice", newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Oracle *TestP1OracleSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.SetPrice(&_TestP1Oracle.TransactOpts, newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Oracle *TestP1OracleTransactorSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Oracle.Contract.SetPrice(&_TestP1Oracle.TransactOpts, newPrice)
}
