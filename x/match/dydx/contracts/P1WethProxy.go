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

// P1WethProxyMetaData contains all meta data concerning the P1WethProxy contract.
var P1WethProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"weth\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"_WETH_\",\"outputs\":[{\"internalType\":\"contractWETH9\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"}],\"name\":\"approveMaximumOnPerpetual\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"depositEth\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawEth\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1WethProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use P1WethProxyMetaData.ABI instead.
var P1WethProxyABI = P1WethProxyMetaData.ABI

// P1WethProxy is an auto generated Go binding around an Ethereum contract.
type P1WethProxy struct {
	P1WethProxyCaller     // Read-only binding to the contract
	P1WethProxyTransactor // Write-only binding to the contract
	P1WethProxyFilterer   // Log filterer for contract events
}

// P1WethProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1WethProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1WethProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1WethProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1WethProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1WethProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1WethProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1WethProxySession struct {
	Contract     *P1WethProxy      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1WethProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1WethProxyCallerSession struct {
	Contract *P1WethProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// P1WethProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1WethProxyTransactorSession struct {
	Contract     *P1WethProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// P1WethProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1WethProxyRaw struct {
	Contract *P1WethProxy // Generic contract binding to access the raw methods on
}

// P1WethProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1WethProxyCallerRaw struct {
	Contract *P1WethProxyCaller // Generic read-only contract binding to access the raw methods on
}

// P1WethProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1WethProxyTransactorRaw struct {
	Contract *P1WethProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1WethProxy creates a new instance of P1WethProxy, bound to a specific deployed contract.
func NewP1WethProxy(address common.Address, backend bind.ContractBackend) (*P1WethProxy, error) {
	contract, err := bindP1WethProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1WethProxy{P1WethProxyCaller: P1WethProxyCaller{contract: contract}, P1WethProxyTransactor: P1WethProxyTransactor{contract: contract}, P1WethProxyFilterer: P1WethProxyFilterer{contract: contract}}, nil
}

// NewP1WethProxyCaller creates a new read-only instance of P1WethProxy, bound to a specific deployed contract.
func NewP1WethProxyCaller(address common.Address, caller bind.ContractCaller) (*P1WethProxyCaller, error) {
	contract, err := bindP1WethProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1WethProxyCaller{contract: contract}, nil
}

// NewP1WethProxyTransactor creates a new write-only instance of P1WethProxy, bound to a specific deployed contract.
func NewP1WethProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*P1WethProxyTransactor, error) {
	contract, err := bindP1WethProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1WethProxyTransactor{contract: contract}, nil
}

// NewP1WethProxyFilterer creates a new log filterer instance of P1WethProxy, bound to a specific deployed contract.
func NewP1WethProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*P1WethProxyFilterer, error) {
	contract, err := bindP1WethProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1WethProxyFilterer{contract: contract}, nil
}

// bindP1WethProxy binds a generic wrapper to an already deployed contract.
func bindP1WethProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1WethProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1WethProxy *P1WethProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1WethProxy.Contract.P1WethProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1WethProxy *P1WethProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1WethProxy.Contract.P1WethProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1WethProxy *P1WethProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1WethProxy.Contract.P1WethProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1WethProxy *P1WethProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1WethProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1WethProxy *P1WethProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1WethProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1WethProxy *P1WethProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1WethProxy.Contract.contract.Transact(opts, method, params...)
}

// WETH is a free data retrieval call binding the contract method 0x0d4eec8f.
//
// Solidity: function _WETH_() view returns(address)
func (_P1WethProxy *P1WethProxyCaller) WETH(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1WethProxy.contract.Call(opts, &out, "_WETH_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH is a free data retrieval call binding the contract method 0x0d4eec8f.
//
// Solidity: function _WETH_() view returns(address)
func (_P1WethProxy *P1WethProxySession) WETH() (common.Address, error) {
	return _P1WethProxy.Contract.WETH(&_P1WethProxy.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0x0d4eec8f.
//
// Solidity: function _WETH_() view returns(address)
func (_P1WethProxy *P1WethProxyCallerSession) WETH() (common.Address, error) {
	return _P1WethProxy.Contract.WETH(&_P1WethProxy.CallOpts)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1WethProxy *P1WethProxyTransactor) ApproveMaximumOnPerpetual(opts *bind.TransactOpts, perpetual common.Address) (*types.Transaction, error) {
	return _P1WethProxy.contract.Transact(opts, "approveMaximumOnPerpetual", perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1WethProxy *P1WethProxySession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1WethProxy.Contract.ApproveMaximumOnPerpetual(&_P1WethProxy.TransactOpts, perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1WethProxy *P1WethProxyTransactorSession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1WethProxy.Contract.ApproveMaximumOnPerpetual(&_P1WethProxy.TransactOpts, perpetual)
}

// DepositEth is a paid mutator transaction binding the contract method 0x1e04e856.
//
// Solidity: function depositEth(address perpetual, address account) payable returns()
func (_P1WethProxy *P1WethProxyTransactor) DepositEth(opts *bind.TransactOpts, perpetual common.Address, account common.Address) (*types.Transaction, error) {
	return _P1WethProxy.contract.Transact(opts, "depositEth", perpetual, account)
}

// DepositEth is a paid mutator transaction binding the contract method 0x1e04e856.
//
// Solidity: function depositEth(address perpetual, address account) payable returns()
func (_P1WethProxy *P1WethProxySession) DepositEth(perpetual common.Address, account common.Address) (*types.Transaction, error) {
	return _P1WethProxy.Contract.DepositEth(&_P1WethProxy.TransactOpts, perpetual, account)
}

// DepositEth is a paid mutator transaction binding the contract method 0x1e04e856.
//
// Solidity: function depositEth(address perpetual, address account) payable returns()
func (_P1WethProxy *P1WethProxyTransactorSession) DepositEth(perpetual common.Address, account common.Address) (*types.Transaction, error) {
	return _P1WethProxy.Contract.DepositEth(&_P1WethProxy.TransactOpts, perpetual, account)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0xf2f938dc.
//
// Solidity: function withdrawEth(address perpetual, address account, address destination, uint256 amount) returns()
func (_P1WethProxy *P1WethProxyTransactor) WithdrawEth(opts *bind.TransactOpts, perpetual common.Address, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1WethProxy.contract.Transact(opts, "withdrawEth", perpetual, account, destination, amount)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0xf2f938dc.
//
// Solidity: function withdrawEth(address perpetual, address account, address destination, uint256 amount) returns()
func (_P1WethProxy *P1WethProxySession) WithdrawEth(perpetual common.Address, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1WethProxy.Contract.WithdrawEth(&_P1WethProxy.TransactOpts, perpetual, account, destination, amount)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0xf2f938dc.
//
// Solidity: function withdrawEth(address perpetual, address account, address destination, uint256 amount) returns()
func (_P1WethProxy *P1WethProxyTransactorSession) WithdrawEth(perpetual common.Address, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1WethProxy.Contract.WithdrawEth(&_P1WethProxy.TransactOpts, perpetual, account, destination, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_P1WethProxy *P1WethProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _P1WethProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_P1WethProxy *P1WethProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _P1WethProxy.Contract.Fallback(&_P1WethProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_P1WethProxy *P1WethProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _P1WethProxy.Contract.Fallback(&_P1WethProxy.TransactOpts, calldata)
}
