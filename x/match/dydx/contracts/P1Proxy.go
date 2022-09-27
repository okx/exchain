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

// P1ProxyMetaData contains all meta data concerning the P1Proxy contract.
var P1ProxyMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"}],\"name\":\"approveMaximumOnPerpetual\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use P1ProxyMetaData.ABI instead.
var P1ProxyABI = P1ProxyMetaData.ABI

// P1Proxy is an auto generated Go binding around an Ethereum contract.
type P1Proxy struct {
	P1ProxyCaller     // Read-only binding to the contract
	P1ProxyTransactor // Write-only binding to the contract
	P1ProxyFilterer   // Log filterer for contract events
}

// P1ProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1ProxySession struct {
	Contract     *P1Proxy          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1ProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1ProxyCallerSession struct {
	Contract *P1ProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// P1ProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1ProxyTransactorSession struct {
	Contract     *P1ProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1ProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1ProxyRaw struct {
	Contract *P1Proxy // Generic contract binding to access the raw methods on
}

// P1ProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1ProxyCallerRaw struct {
	Contract *P1ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// P1ProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1ProxyTransactorRaw struct {
	Contract *P1ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Proxy creates a new instance of P1Proxy, bound to a specific deployed contract.
func NewP1Proxy(address common.Address, backend bind.ContractBackend) (*P1Proxy, error) {
	contract, err := bindP1Proxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Proxy{P1ProxyCaller: P1ProxyCaller{contract: contract}, P1ProxyTransactor: P1ProxyTransactor{contract: contract}, P1ProxyFilterer: P1ProxyFilterer{contract: contract}}, nil
}

// NewP1ProxyCaller creates a new read-only instance of P1Proxy, bound to a specific deployed contract.
func NewP1ProxyCaller(address common.Address, caller bind.ContractCaller) (*P1ProxyCaller, error) {
	contract, err := bindP1Proxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1ProxyCaller{contract: contract}, nil
}

// NewP1ProxyTransactor creates a new write-only instance of P1Proxy, bound to a specific deployed contract.
func NewP1ProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*P1ProxyTransactor, error) {
	contract, err := bindP1Proxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1ProxyTransactor{contract: contract}, nil
}

// NewP1ProxyFilterer creates a new log filterer instance of P1Proxy, bound to a specific deployed contract.
func NewP1ProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*P1ProxyFilterer, error) {
	contract, err := bindP1Proxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1ProxyFilterer{contract: contract}, nil
}

// bindP1Proxy binds a generic wrapper to an already deployed contract.
func bindP1Proxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1ProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Proxy *P1ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Proxy.Contract.P1ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Proxy *P1ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Proxy.Contract.P1ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Proxy *P1ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Proxy.Contract.P1ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Proxy *P1ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Proxy *P1ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Proxy *P1ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Proxy.Contract.contract.Transact(opts, method, params...)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1Proxy *P1ProxyTransactor) ApproveMaximumOnPerpetual(opts *bind.TransactOpts, perpetual common.Address) (*types.Transaction, error) {
	return _P1Proxy.contract.Transact(opts, "approveMaximumOnPerpetual", perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1Proxy *P1ProxySession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1Proxy.Contract.ApproveMaximumOnPerpetual(&_P1Proxy.TransactOpts, perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1Proxy *P1ProxyTransactorSession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1Proxy.Contract.ApproveMaximumOnPerpetual(&_P1Proxy.TransactOpts, perpetual)
}
