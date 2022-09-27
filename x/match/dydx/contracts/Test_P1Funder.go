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

// TestP1FunderMetaData contains all meta data concerning the TestP1Funder contract.
var TestP1FunderMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"_FUNDING_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_FUNDING_IS_POSITIVE_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getFunding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"newFunding\",\"type\":\"uint256\"}],\"name\":\"setFunding\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestP1FunderABI is the input ABI used to generate the binding from.
// Deprecated: Use TestP1FunderMetaData.ABI instead.
var TestP1FunderABI = TestP1FunderMetaData.ABI

// TestP1Funder is an auto generated Go binding around an Ethereum contract.
type TestP1Funder struct {
	TestP1FunderCaller     // Read-only binding to the contract
	TestP1FunderTransactor // Write-only binding to the contract
	TestP1FunderFilterer   // Log filterer for contract events
}

// TestP1FunderCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestP1FunderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1FunderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestP1FunderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1FunderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestP1FunderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1FunderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestP1FunderSession struct {
	Contract     *TestP1Funder     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestP1FunderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestP1FunderCallerSession struct {
	Contract *TestP1FunderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TestP1FunderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestP1FunderTransactorSession struct {
	Contract     *TestP1FunderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TestP1FunderRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestP1FunderRaw struct {
	Contract *TestP1Funder // Generic contract binding to access the raw methods on
}

// TestP1FunderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestP1FunderCallerRaw struct {
	Contract *TestP1FunderCaller // Generic read-only contract binding to access the raw methods on
}

// TestP1FunderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestP1FunderTransactorRaw struct {
	Contract *TestP1FunderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestP1Funder creates a new instance of TestP1Funder, bound to a specific deployed contract.
func NewTestP1Funder(address common.Address, backend bind.ContractBackend) (*TestP1Funder, error) {
	contract, err := bindTestP1Funder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestP1Funder{TestP1FunderCaller: TestP1FunderCaller{contract: contract}, TestP1FunderTransactor: TestP1FunderTransactor{contract: contract}, TestP1FunderFilterer: TestP1FunderFilterer{contract: contract}}, nil
}

// NewTestP1FunderCaller creates a new read-only instance of TestP1Funder, bound to a specific deployed contract.
func NewTestP1FunderCaller(address common.Address, caller bind.ContractCaller) (*TestP1FunderCaller, error) {
	contract, err := bindTestP1Funder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1FunderCaller{contract: contract}, nil
}

// NewTestP1FunderTransactor creates a new write-only instance of TestP1Funder, bound to a specific deployed contract.
func NewTestP1FunderTransactor(address common.Address, transactor bind.ContractTransactor) (*TestP1FunderTransactor, error) {
	contract, err := bindTestP1Funder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1FunderTransactor{contract: contract}, nil
}

// NewTestP1FunderFilterer creates a new log filterer instance of TestP1Funder, bound to a specific deployed contract.
func NewTestP1FunderFilterer(address common.Address, filterer bind.ContractFilterer) (*TestP1FunderFilterer, error) {
	contract, err := bindTestP1Funder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestP1FunderFilterer{contract: contract}, nil
}

// bindTestP1Funder binds a generic wrapper to an already deployed contract.
func bindTestP1Funder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestP1FunderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Funder *TestP1FunderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Funder.Contract.TestP1FunderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Funder *TestP1FunderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Funder.Contract.TestP1FunderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Funder *TestP1FunderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Funder.Contract.TestP1FunderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Funder *TestP1FunderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Funder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Funder *TestP1FunderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Funder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Funder *TestP1FunderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Funder.Contract.contract.Transact(opts, method, params...)
}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Funder *TestP1FunderCaller) FUNDING(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Funder.contract.Call(opts, &out, "_FUNDING_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Funder *TestP1FunderSession) FUNDING() (*big.Int, error) {
	return _TestP1Funder.Contract.FUNDING(&_TestP1Funder.CallOpts)
}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Funder *TestP1FunderCallerSession) FUNDING() (*big.Int, error) {
	return _TestP1Funder.Contract.FUNDING(&_TestP1Funder.CallOpts)
}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Funder *TestP1FunderCaller) FUNDINGISPOSITIVE(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestP1Funder.contract.Call(opts, &out, "_FUNDING_IS_POSITIVE_")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Funder *TestP1FunderSession) FUNDINGISPOSITIVE() (bool, error) {
	return _TestP1Funder.Contract.FUNDINGISPOSITIVE(&_TestP1Funder.CallOpts)
}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Funder *TestP1FunderCallerSession) FUNDINGISPOSITIVE() (bool, error) {
	return _TestP1Funder.Contract.FUNDINGISPOSITIVE(&_TestP1Funder.CallOpts)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Funder *TestP1FunderCaller) GetFunding(opts *bind.CallOpts, arg0 *big.Int) (bool, *big.Int, error) {
	var out []interface{}
	err := _TestP1Funder.contract.Call(opts, &out, "getFunding", arg0)

	if err != nil {
		return *new(bool), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Funder *TestP1FunderSession) GetFunding(arg0 *big.Int) (bool, *big.Int, error) {
	return _TestP1Funder.Contract.GetFunding(&_TestP1Funder.CallOpts, arg0)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Funder *TestP1FunderCallerSession) GetFunding(arg0 *big.Int) (bool, *big.Int, error) {
	return _TestP1Funder.Contract.GetFunding(&_TestP1Funder.CallOpts, arg0)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Funder *TestP1FunderTransactor) SetFunding(opts *bind.TransactOpts, isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Funder.contract.Transact(opts, "setFunding", isPositive, newFunding)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Funder *TestP1FunderSession) SetFunding(isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Funder.Contract.SetFunding(&_TestP1Funder.TransactOpts, isPositive, newFunding)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Funder *TestP1FunderTransactorSession) SetFunding(isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Funder.Contract.SetFunding(&_TestP1Funder.TransactOpts, isPositive, newFunding)
}
