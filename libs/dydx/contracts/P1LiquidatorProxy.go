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

// SignedMathInt is an auto generated low-level Go binding around an user-defined struct.
type SignedMathInt struct {
	Value      *big.Int
	IsPositive bool
}

// P1LiquidatorProxyMetaData contains all meta data concerning the P1LiquidatorProxy contract.
var P1LiquidatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetualV1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"insuranceFund\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"insuranceFee\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"insuranceFee\",\"type\":\"uint256\"}],\"name\":\"LogInsuranceFeeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"insuranceFund\",\"type\":\"address\"}],\"name\":\"LogInsuranceFundSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidatee\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"LogLiquidatorProxyUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_INSURANCE_FEE_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_INSURANCE_FUND_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_LIQUIDATION_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_PERPETUAL_V1_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"approveMaximumOnPerpetual\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"liquidatee\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"maxPosition\",\"type\":\"tuple\"}],\"name\":\"liquidate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"insuranceFund\",\"type\":\"address\"}],\"name\":\"setInsuranceFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"insuranceFee\",\"type\":\"uint256\"}],\"name\":\"setInsuranceFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1LiquidatorProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use P1LiquidatorProxyMetaData.ABI instead.
var P1LiquidatorProxyABI = P1LiquidatorProxyMetaData.ABI

// P1LiquidatorProxy is an auto generated Go binding around an Ethereum contract.
type P1LiquidatorProxy struct {
	P1LiquidatorProxyCaller     // Read-only binding to the contract
	P1LiquidatorProxyTransactor // Write-only binding to the contract
	P1LiquidatorProxyFilterer   // Log filterer for contract events
}

// P1LiquidatorProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1LiquidatorProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidatorProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1LiquidatorProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidatorProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1LiquidatorProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidatorProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1LiquidatorProxySession struct {
	Contract     *P1LiquidatorProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1LiquidatorProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1LiquidatorProxyCallerSession struct {
	Contract *P1LiquidatorProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// P1LiquidatorProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1LiquidatorProxyTransactorSession struct {
	Contract     *P1LiquidatorProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// P1LiquidatorProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1LiquidatorProxyRaw struct {
	Contract *P1LiquidatorProxy // Generic contract binding to access the raw methods on
}

// P1LiquidatorProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1LiquidatorProxyCallerRaw struct {
	Contract *P1LiquidatorProxyCaller // Generic read-only contract binding to access the raw methods on
}

// P1LiquidatorProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1LiquidatorProxyTransactorRaw struct {
	Contract *P1LiquidatorProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1LiquidatorProxy creates a new instance of P1LiquidatorProxy, bound to a specific deployed contract.
func NewP1LiquidatorProxy(address common.Address, backend bind.ContractBackend) (*P1LiquidatorProxy, error) {
	contract, err := bindP1LiquidatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxy{P1LiquidatorProxyCaller: P1LiquidatorProxyCaller{contract: contract}, P1LiquidatorProxyTransactor: P1LiquidatorProxyTransactor{contract: contract}, P1LiquidatorProxyFilterer: P1LiquidatorProxyFilterer{contract: contract}}, nil
}

// NewP1LiquidatorProxyCaller creates a new read-only instance of P1LiquidatorProxy, bound to a specific deployed contract.
func NewP1LiquidatorProxyCaller(address common.Address, caller bind.ContractCaller) (*P1LiquidatorProxyCaller, error) {
	contract, err := bindP1LiquidatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyCaller{contract: contract}, nil
}

// NewP1LiquidatorProxyTransactor creates a new write-only instance of P1LiquidatorProxy, bound to a specific deployed contract.
func NewP1LiquidatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*P1LiquidatorProxyTransactor, error) {
	contract, err := bindP1LiquidatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyTransactor{contract: contract}, nil
}

// NewP1LiquidatorProxyFilterer creates a new log filterer instance of P1LiquidatorProxy, bound to a specific deployed contract.
func NewP1LiquidatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*P1LiquidatorProxyFilterer, error) {
	contract, err := bindP1LiquidatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyFilterer{contract: contract}, nil
}

// bindP1LiquidatorProxy binds a generic wrapper to an already deployed contract.
func bindP1LiquidatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1LiquidatorProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1LiquidatorProxy *P1LiquidatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1LiquidatorProxy.Contract.P1LiquidatorProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1LiquidatorProxy *P1LiquidatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.P1LiquidatorProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1LiquidatorProxy *P1LiquidatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.P1LiquidatorProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1LiquidatorProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.contract.Transact(opts, method, params...)
}

// INSURANCEFEE is a free data retrieval call binding the contract method 0x4ce7c2ca.
//
// Solidity: function _INSURANCE_FEE_() view returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) INSURANCEFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "_INSURANCE_FEE_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INSURANCEFEE is a free data retrieval call binding the contract method 0x4ce7c2ca.
//
// Solidity: function _INSURANCE_FEE_() view returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) INSURANCEFEE() (*big.Int, error) {
	return _P1LiquidatorProxy.Contract.INSURANCEFEE(&_P1LiquidatorProxy.CallOpts)
}

// INSURANCEFEE is a free data retrieval call binding the contract method 0x4ce7c2ca.
//
// Solidity: function _INSURANCE_FEE_() view returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) INSURANCEFEE() (*big.Int, error) {
	return _P1LiquidatorProxy.Contract.INSURANCEFEE(&_P1LiquidatorProxy.CallOpts)
}

// INSURANCEFUND is a free data retrieval call binding the contract method 0x3fa8c92a.
//
// Solidity: function _INSURANCE_FUND_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) INSURANCEFUND(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "_INSURANCE_FUND_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// INSURANCEFUND is a free data retrieval call binding the contract method 0x3fa8c92a.
//
// Solidity: function _INSURANCE_FUND_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) INSURANCEFUND() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.INSURANCEFUND(&_P1LiquidatorProxy.CallOpts)
}

// INSURANCEFUND is a free data retrieval call binding the contract method 0x3fa8c92a.
//
// Solidity: function _INSURANCE_FUND_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) INSURANCEFUND() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.INSURANCEFUND(&_P1LiquidatorProxy.CallOpts)
}

// LIQUIDATION is a free data retrieval call binding the contract method 0x786ed92e.
//
// Solidity: function _LIQUIDATION_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) LIQUIDATION(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "_LIQUIDATION_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LIQUIDATION is a free data retrieval call binding the contract method 0x786ed92e.
//
// Solidity: function _LIQUIDATION_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) LIQUIDATION() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.LIQUIDATION(&_P1LiquidatorProxy.CallOpts)
}

// LIQUIDATION is a free data retrieval call binding the contract method 0x786ed92e.
//
// Solidity: function _LIQUIDATION_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) LIQUIDATION() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.LIQUIDATION(&_P1LiquidatorProxy.CallOpts)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) PERPETUALV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "_PERPETUAL_V1_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) PERPETUALV1() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.PERPETUALV1(&_P1LiquidatorProxy.CallOpts)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) PERPETUALV1() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.PERPETUALV1(&_P1LiquidatorProxy.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) IsOwner() (bool, error) {
	return _P1LiquidatorProxy.Contract.IsOwner(&_P1LiquidatorProxy.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) IsOwner() (bool, error) {
	return _P1LiquidatorProxy.Contract.IsOwner(&_P1LiquidatorProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1LiquidatorProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) Owner() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.Owner(&_P1LiquidatorProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1LiquidatorProxy *P1LiquidatorProxyCallerSession) Owner() (common.Address, error) {
	return _P1LiquidatorProxy.Contract.Owner(&_P1LiquidatorProxy.CallOpts)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0x8ffac733.
//
// Solidity: function approveMaximumOnPerpetual() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) ApproveMaximumOnPerpetual(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "approveMaximumOnPerpetual")
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0x8ffac733.
//
// Solidity: function approveMaximumOnPerpetual() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxySession) ApproveMaximumOnPerpetual() (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.ApproveMaximumOnPerpetual(&_P1LiquidatorProxy.TransactOpts)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0x8ffac733.
//
// Solidity: function approveMaximumOnPerpetual() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) ApproveMaximumOnPerpetual() (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.ApproveMaximumOnPerpetual(&_P1LiquidatorProxy.TransactOpts)
}

// Liquidate is a paid mutator transaction binding the contract method 0xa5d9c62a.
//
// Solidity: function liquidate(address liquidatee, address liquidator, bool isBuy, (uint256,bool) maxPosition) returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) Liquidate(opts *bind.TransactOpts, liquidatee common.Address, liquidator common.Address, isBuy bool, maxPosition SignedMathInt) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "liquidate", liquidatee, liquidator, isBuy, maxPosition)
}

// Liquidate is a paid mutator transaction binding the contract method 0xa5d9c62a.
//
// Solidity: function liquidate(address liquidatee, address liquidator, bool isBuy, (uint256,bool) maxPosition) returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxySession) Liquidate(liquidatee common.Address, liquidator common.Address, isBuy bool, maxPosition SignedMathInt) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.Liquidate(&_P1LiquidatorProxy.TransactOpts, liquidatee, liquidator, isBuy, maxPosition)
}

// Liquidate is a paid mutator transaction binding the contract method 0xa5d9c62a.
//
// Solidity: function liquidate(address liquidatee, address liquidator, bool isBuy, (uint256,bool) maxPosition) returns(uint256)
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) Liquidate(liquidatee common.Address, liquidator common.Address, isBuy bool, maxPosition SignedMathInt) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.Liquidate(&_P1LiquidatorProxy.TransactOpts, liquidatee, liquidator, isBuy, maxPosition)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxySession) RenounceOwnership() (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.RenounceOwnership(&_P1LiquidatorProxy.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.RenounceOwnership(&_P1LiquidatorProxy.TransactOpts)
}

// SetInsuranceFee is a paid mutator transaction binding the contract method 0xba32681e.
//
// Solidity: function setInsuranceFee(uint256 insuranceFee) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) SetInsuranceFee(opts *bind.TransactOpts, insuranceFee *big.Int) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "setInsuranceFee", insuranceFee)
}

// SetInsuranceFee is a paid mutator transaction binding the contract method 0xba32681e.
//
// Solidity: function setInsuranceFee(uint256 insuranceFee) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxySession) SetInsuranceFee(insuranceFee *big.Int) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.SetInsuranceFee(&_P1LiquidatorProxy.TransactOpts, insuranceFee)
}

// SetInsuranceFee is a paid mutator transaction binding the contract method 0xba32681e.
//
// Solidity: function setInsuranceFee(uint256 insuranceFee) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) SetInsuranceFee(insuranceFee *big.Int) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.SetInsuranceFee(&_P1LiquidatorProxy.TransactOpts, insuranceFee)
}

// SetInsuranceFund is a paid mutator transaction binding the contract method 0xc3c05293.
//
// Solidity: function setInsuranceFund(address insuranceFund) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) SetInsuranceFund(opts *bind.TransactOpts, insuranceFund common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "setInsuranceFund", insuranceFund)
}

// SetInsuranceFund is a paid mutator transaction binding the contract method 0xc3c05293.
//
// Solidity: function setInsuranceFund(address insuranceFund) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxySession) SetInsuranceFund(insuranceFund common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.SetInsuranceFund(&_P1LiquidatorProxy.TransactOpts, insuranceFund)
}

// SetInsuranceFund is a paid mutator transaction binding the contract method 0xc3c05293.
//
// Solidity: function setInsuranceFund(address insuranceFund) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) SetInsuranceFund(insuranceFund common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.SetInsuranceFund(&_P1LiquidatorProxy.TransactOpts, insuranceFund)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.TransferOwnership(&_P1LiquidatorProxy.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1LiquidatorProxy *P1LiquidatorProxyTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1LiquidatorProxy.Contract.TransferOwnership(&_P1LiquidatorProxy.TransactOpts, newOwner)
}

// P1LiquidatorProxyLogInsuranceFeeSetIterator is returned from FilterLogInsuranceFeeSet and is used to iterate over the raw logs and unpacked data for LogInsuranceFeeSet events raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogInsuranceFeeSetIterator struct {
	Event *P1LiquidatorProxyLogInsuranceFeeSet // Event containing the contract specifics and raw log

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
func (it *P1LiquidatorProxyLogInsuranceFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1LiquidatorProxyLogInsuranceFeeSet)
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
		it.Event = new(P1LiquidatorProxyLogInsuranceFeeSet)
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
func (it *P1LiquidatorProxyLogInsuranceFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1LiquidatorProxyLogInsuranceFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1LiquidatorProxyLogInsuranceFeeSet represents a LogInsuranceFeeSet event raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogInsuranceFeeSet struct {
	InsuranceFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogInsuranceFeeSet is a free log retrieval operation binding the contract event 0xb66e25e76b9dc01dd4aa8dffb4c93b795262619460b2d05e1c4cc18388f78f48.
//
// Solidity: event LogInsuranceFeeSet(uint256 insuranceFee)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) FilterLogInsuranceFeeSet(opts *bind.FilterOpts) (*P1LiquidatorProxyLogInsuranceFeeSetIterator, error) {

	logs, sub, err := _P1LiquidatorProxy.contract.FilterLogs(opts, "LogInsuranceFeeSet")
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyLogInsuranceFeeSetIterator{contract: _P1LiquidatorProxy.contract, event: "LogInsuranceFeeSet", logs: logs, sub: sub}, nil
}

// WatchLogInsuranceFeeSet is a free log subscription operation binding the contract event 0xb66e25e76b9dc01dd4aa8dffb4c93b795262619460b2d05e1c4cc18388f78f48.
//
// Solidity: event LogInsuranceFeeSet(uint256 insuranceFee)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) WatchLogInsuranceFeeSet(opts *bind.WatchOpts, sink chan<- *P1LiquidatorProxyLogInsuranceFeeSet) (event.Subscription, error) {

	logs, sub, err := _P1LiquidatorProxy.contract.WatchLogs(opts, "LogInsuranceFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1LiquidatorProxyLogInsuranceFeeSet)
				if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogInsuranceFeeSet", log); err != nil {
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

// ParseLogInsuranceFeeSet is a log parse operation binding the contract event 0xb66e25e76b9dc01dd4aa8dffb4c93b795262619460b2d05e1c4cc18388f78f48.
//
// Solidity: event LogInsuranceFeeSet(uint256 insuranceFee)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) ParseLogInsuranceFeeSet(log types.Log) (*P1LiquidatorProxyLogInsuranceFeeSet, error) {
	event := new(P1LiquidatorProxyLogInsuranceFeeSet)
	if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogInsuranceFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1LiquidatorProxyLogInsuranceFundSetIterator is returned from FilterLogInsuranceFundSet and is used to iterate over the raw logs and unpacked data for LogInsuranceFundSet events raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogInsuranceFundSetIterator struct {
	Event *P1LiquidatorProxyLogInsuranceFundSet // Event containing the contract specifics and raw log

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
func (it *P1LiquidatorProxyLogInsuranceFundSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1LiquidatorProxyLogInsuranceFundSet)
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
		it.Event = new(P1LiquidatorProxyLogInsuranceFundSet)
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
func (it *P1LiquidatorProxyLogInsuranceFundSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1LiquidatorProxyLogInsuranceFundSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1LiquidatorProxyLogInsuranceFundSet represents a LogInsuranceFundSet event raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogInsuranceFundSet struct {
	InsuranceFund common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogInsuranceFundSet is a free log retrieval operation binding the contract event 0x02be8aef8c7fb3cfe392924a6868452212a6fb3db277cc1a24af9f7c4af80ebd.
//
// Solidity: event LogInsuranceFundSet(address insuranceFund)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) FilterLogInsuranceFundSet(opts *bind.FilterOpts) (*P1LiquidatorProxyLogInsuranceFundSetIterator, error) {

	logs, sub, err := _P1LiquidatorProxy.contract.FilterLogs(opts, "LogInsuranceFundSet")
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyLogInsuranceFundSetIterator{contract: _P1LiquidatorProxy.contract, event: "LogInsuranceFundSet", logs: logs, sub: sub}, nil
}

// WatchLogInsuranceFundSet is a free log subscription operation binding the contract event 0x02be8aef8c7fb3cfe392924a6868452212a6fb3db277cc1a24af9f7c4af80ebd.
//
// Solidity: event LogInsuranceFundSet(address insuranceFund)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) WatchLogInsuranceFundSet(opts *bind.WatchOpts, sink chan<- *P1LiquidatorProxyLogInsuranceFundSet) (event.Subscription, error) {

	logs, sub, err := _P1LiquidatorProxy.contract.WatchLogs(opts, "LogInsuranceFundSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1LiquidatorProxyLogInsuranceFundSet)
				if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogInsuranceFundSet", log); err != nil {
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

// ParseLogInsuranceFundSet is a log parse operation binding the contract event 0x02be8aef8c7fb3cfe392924a6868452212a6fb3db277cc1a24af9f7c4af80ebd.
//
// Solidity: event LogInsuranceFundSet(address insuranceFund)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) ParseLogInsuranceFundSet(log types.Log) (*P1LiquidatorProxyLogInsuranceFundSet, error) {
	event := new(P1LiquidatorProxyLogInsuranceFundSet)
	if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogInsuranceFundSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1LiquidatorProxyLogLiquidatorProxyUsedIterator is returned from FilterLogLiquidatorProxyUsed and is used to iterate over the raw logs and unpacked data for LogLiquidatorProxyUsed events raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogLiquidatorProxyUsedIterator struct {
	Event *P1LiquidatorProxyLogLiquidatorProxyUsed // Event containing the contract specifics and raw log

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
func (it *P1LiquidatorProxyLogLiquidatorProxyUsedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1LiquidatorProxyLogLiquidatorProxyUsed)
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
		it.Event = new(P1LiquidatorProxyLogLiquidatorProxyUsed)
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
func (it *P1LiquidatorProxyLogLiquidatorProxyUsedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1LiquidatorProxyLogLiquidatorProxyUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1LiquidatorProxyLogLiquidatorProxyUsed represents a LogLiquidatorProxyUsed event raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyLogLiquidatorProxyUsed struct {
	Liquidatee        common.Address
	Liquidator        common.Address
	IsBuy             bool
	LiquidationAmount *big.Int
	FeeAmount         *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogLiquidatorProxyUsed is a free log retrieval operation binding the contract event 0x56f54e5e291f84831023c9ddf34fe42973dae320af11193db2b5f7af27719ba6.
//
// Solidity: event LogLiquidatorProxyUsed(address indexed liquidatee, address indexed liquidator, bool isBuy, uint256 liquidationAmount, uint256 feeAmount)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) FilterLogLiquidatorProxyUsed(opts *bind.FilterOpts, liquidatee []common.Address, liquidator []common.Address) (*P1LiquidatorProxyLogLiquidatorProxyUsedIterator, error) {

	var liquidateeRule []interface{}
	for _, liquidateeItem := range liquidatee {
		liquidateeRule = append(liquidateeRule, liquidateeItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _P1LiquidatorProxy.contract.FilterLogs(opts, "LogLiquidatorProxyUsed", liquidateeRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyLogLiquidatorProxyUsedIterator{contract: _P1LiquidatorProxy.contract, event: "LogLiquidatorProxyUsed", logs: logs, sub: sub}, nil
}

// WatchLogLiquidatorProxyUsed is a free log subscription operation binding the contract event 0x56f54e5e291f84831023c9ddf34fe42973dae320af11193db2b5f7af27719ba6.
//
// Solidity: event LogLiquidatorProxyUsed(address indexed liquidatee, address indexed liquidator, bool isBuy, uint256 liquidationAmount, uint256 feeAmount)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) WatchLogLiquidatorProxyUsed(opts *bind.WatchOpts, sink chan<- *P1LiquidatorProxyLogLiquidatorProxyUsed, liquidatee []common.Address, liquidator []common.Address) (event.Subscription, error) {

	var liquidateeRule []interface{}
	for _, liquidateeItem := range liquidatee {
		liquidateeRule = append(liquidateeRule, liquidateeItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _P1LiquidatorProxy.contract.WatchLogs(opts, "LogLiquidatorProxyUsed", liquidateeRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1LiquidatorProxyLogLiquidatorProxyUsed)
				if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogLiquidatorProxyUsed", log); err != nil {
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

// ParseLogLiquidatorProxyUsed is a log parse operation binding the contract event 0x56f54e5e291f84831023c9ddf34fe42973dae320af11193db2b5f7af27719ba6.
//
// Solidity: event LogLiquidatorProxyUsed(address indexed liquidatee, address indexed liquidator, bool isBuy, uint256 liquidationAmount, uint256 feeAmount)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) ParseLogLiquidatorProxyUsed(log types.Log) (*P1LiquidatorProxyLogLiquidatorProxyUsed, error) {
	event := new(P1LiquidatorProxyLogLiquidatorProxyUsed)
	if err := _P1LiquidatorProxy.contract.UnpackLog(event, "LogLiquidatorProxyUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1LiquidatorProxyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyOwnershipTransferredIterator struct {
	Event *P1LiquidatorProxyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1LiquidatorProxyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1LiquidatorProxyOwnershipTransferred)
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
		it.Event = new(P1LiquidatorProxyOwnershipTransferred)
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
func (it *P1LiquidatorProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1LiquidatorProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1LiquidatorProxyOwnershipTransferred represents a OwnershipTransferred event raised by the P1LiquidatorProxy contract.
type P1LiquidatorProxyOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1LiquidatorProxyOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1LiquidatorProxy.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1LiquidatorProxyOwnershipTransferredIterator{contract: _P1LiquidatorProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1LiquidatorProxyOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1LiquidatorProxy.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1LiquidatorProxyOwnershipTransferred)
				if err := _P1LiquidatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1LiquidatorProxy *P1LiquidatorProxyFilterer) ParseOwnershipTransferred(log types.Log) (*P1LiquidatorProxyOwnershipTransferred, error) {
	event := new(P1LiquidatorProxyOwnershipTransferred)
	if err := _P1LiquidatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
