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

// P1DeleveragingMetaData contains all meta data concerning the P1Deleveraging contract.
var P1DeleveragingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetualV1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"deleveragingOperator\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oraclePrice\",\"type\":\"uint256\"}],\"name\":\"LogDeleveraged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"deleveragingOperator\",\"type\":\"address\"}],\"name\":\"LogDeleveragingOperatorSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"LogMarkedForDeleveraging\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"LogUnmarkedForDeleveraging\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DELEVERAGING_TIMELOCK_S\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_DELEVERAGING_OPERATOR_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_MARKED_TIMESTAMP_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_PERPETUAL_V1_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"mark\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"unmark\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"setDeleveragingOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1DeleveragingABI is the input ABI used to generate the binding from.
// Deprecated: Use P1DeleveragingMetaData.ABI instead.
var P1DeleveragingABI = P1DeleveragingMetaData.ABI

// P1Deleveraging is an auto generated Go binding around an Ethereum contract.
type P1Deleveraging struct {
	P1DeleveragingCaller     // Read-only binding to the contract
	P1DeleveragingTransactor // Write-only binding to the contract
	P1DeleveragingFilterer   // Log filterer for contract events
}

// P1DeleveragingCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1DeleveragingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1DeleveragingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1DeleveragingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1DeleveragingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1DeleveragingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1DeleveragingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1DeleveragingSession struct {
	Contract     *P1Deleveraging   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1DeleveragingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1DeleveragingCallerSession struct {
	Contract *P1DeleveragingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// P1DeleveragingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1DeleveragingTransactorSession struct {
	Contract     *P1DeleveragingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// P1DeleveragingRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1DeleveragingRaw struct {
	Contract *P1Deleveraging // Generic contract binding to access the raw methods on
}

// P1DeleveragingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1DeleveragingCallerRaw struct {
	Contract *P1DeleveragingCaller // Generic read-only contract binding to access the raw methods on
}

// P1DeleveragingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1DeleveragingTransactorRaw struct {
	Contract *P1DeleveragingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Deleveraging creates a new instance of P1Deleveraging, bound to a specific deployed contract.
func NewP1Deleveraging(address common.Address, backend bind.ContractBackend) (*P1Deleveraging, error) {
	contract, err := bindP1Deleveraging(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Deleveraging{P1DeleveragingCaller: P1DeleveragingCaller{contract: contract}, P1DeleveragingTransactor: P1DeleveragingTransactor{contract: contract}, P1DeleveragingFilterer: P1DeleveragingFilterer{contract: contract}}, nil
}

// NewP1DeleveragingCaller creates a new read-only instance of P1Deleveraging, bound to a specific deployed contract.
func NewP1DeleveragingCaller(address common.Address, caller bind.ContractCaller) (*P1DeleveragingCaller, error) {
	contract, err := bindP1Deleveraging(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingCaller{contract: contract}, nil
}

// NewP1DeleveragingTransactor creates a new write-only instance of P1Deleveraging, bound to a specific deployed contract.
func NewP1DeleveragingTransactor(address common.Address, transactor bind.ContractTransactor) (*P1DeleveragingTransactor, error) {
	contract, err := bindP1Deleveraging(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingTransactor{contract: contract}, nil
}

// NewP1DeleveragingFilterer creates a new log filterer instance of P1Deleveraging, bound to a specific deployed contract.
func NewP1DeleveragingFilterer(address common.Address, filterer bind.ContractFilterer) (*P1DeleveragingFilterer, error) {
	contract, err := bindP1Deleveraging(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingFilterer{contract: contract}, nil
}

// bindP1Deleveraging binds a generic wrapper to an already deployed contract.
func bindP1Deleveraging(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1DeleveragingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Deleveraging *P1DeleveragingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Deleveraging.Contract.P1DeleveragingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Deleveraging *P1DeleveragingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.P1DeleveragingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Deleveraging *P1DeleveragingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.P1DeleveragingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Deleveraging *P1DeleveragingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Deleveraging.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Deleveraging *P1DeleveragingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Deleveraging *P1DeleveragingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.contract.Transact(opts, method, params...)
}

// DELEVERAGINGTIMELOCKS is a free data retrieval call binding the contract method 0x741c1195.
//
// Solidity: function DELEVERAGING_TIMELOCK_S() view returns(uint256)
func (_P1Deleveraging *P1DeleveragingCaller) DELEVERAGINGTIMELOCKS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "DELEVERAGING_TIMELOCK_S")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DELEVERAGINGTIMELOCKS is a free data retrieval call binding the contract method 0x741c1195.
//
// Solidity: function DELEVERAGING_TIMELOCK_S() view returns(uint256)
func (_P1Deleveraging *P1DeleveragingSession) DELEVERAGINGTIMELOCKS() (*big.Int, error) {
	return _P1Deleveraging.Contract.DELEVERAGINGTIMELOCKS(&_P1Deleveraging.CallOpts)
}

// DELEVERAGINGTIMELOCKS is a free data retrieval call binding the contract method 0x741c1195.
//
// Solidity: function DELEVERAGING_TIMELOCK_S() view returns(uint256)
func (_P1Deleveraging *P1DeleveragingCallerSession) DELEVERAGINGTIMELOCKS() (*big.Int, error) {
	return _P1Deleveraging.Contract.DELEVERAGINGTIMELOCKS(&_P1Deleveraging.CallOpts)
}

// DELEVERAGINGOPERATOR is a free data retrieval call binding the contract method 0x8e085c98.
//
// Solidity: function _DELEVERAGING_OPERATOR_() view returns(address)
func (_P1Deleveraging *P1DeleveragingCaller) DELEVERAGINGOPERATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "_DELEVERAGING_OPERATOR_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DELEVERAGINGOPERATOR is a free data retrieval call binding the contract method 0x8e085c98.
//
// Solidity: function _DELEVERAGING_OPERATOR_() view returns(address)
func (_P1Deleveraging *P1DeleveragingSession) DELEVERAGINGOPERATOR() (common.Address, error) {
	return _P1Deleveraging.Contract.DELEVERAGINGOPERATOR(&_P1Deleveraging.CallOpts)
}

// DELEVERAGINGOPERATOR is a free data retrieval call binding the contract method 0x8e085c98.
//
// Solidity: function _DELEVERAGING_OPERATOR_() view returns(address)
func (_P1Deleveraging *P1DeleveragingCallerSession) DELEVERAGINGOPERATOR() (common.Address, error) {
	return _P1Deleveraging.Contract.DELEVERAGINGOPERATOR(&_P1Deleveraging.CallOpts)
}

// MARKEDTIMESTAMP is a free data retrieval call binding the contract method 0x5ce999b6.
//
// Solidity: function _MARKED_TIMESTAMP_(address ) view returns(uint256)
func (_P1Deleveraging *P1DeleveragingCaller) MARKEDTIMESTAMP(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "_MARKED_TIMESTAMP_", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MARKEDTIMESTAMP is a free data retrieval call binding the contract method 0x5ce999b6.
//
// Solidity: function _MARKED_TIMESTAMP_(address ) view returns(uint256)
func (_P1Deleveraging *P1DeleveragingSession) MARKEDTIMESTAMP(arg0 common.Address) (*big.Int, error) {
	return _P1Deleveraging.Contract.MARKEDTIMESTAMP(&_P1Deleveraging.CallOpts, arg0)
}

// MARKEDTIMESTAMP is a free data retrieval call binding the contract method 0x5ce999b6.
//
// Solidity: function _MARKED_TIMESTAMP_(address ) view returns(uint256)
func (_P1Deleveraging *P1DeleveragingCallerSession) MARKEDTIMESTAMP(arg0 common.Address) (*big.Int, error) {
	return _P1Deleveraging.Contract.MARKEDTIMESTAMP(&_P1Deleveraging.CallOpts, arg0)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Deleveraging *P1DeleveragingCaller) PERPETUALV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "_PERPETUAL_V1_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Deleveraging *P1DeleveragingSession) PERPETUALV1() (common.Address, error) {
	return _P1Deleveraging.Contract.PERPETUALV1(&_P1Deleveraging.CallOpts)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Deleveraging *P1DeleveragingCallerSession) PERPETUALV1() (common.Address, error) {
	return _P1Deleveraging.Contract.PERPETUALV1(&_P1Deleveraging.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1Deleveraging *P1DeleveragingCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1Deleveraging *P1DeleveragingSession) IsOwner() (bool, error) {
	return _P1Deleveraging.Contract.IsOwner(&_P1Deleveraging.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1Deleveraging *P1DeleveragingCallerSession) IsOwner() (bool, error) {
	return _P1Deleveraging.Contract.IsOwner(&_P1Deleveraging.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1Deleveraging *P1DeleveragingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Deleveraging.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1Deleveraging *P1DeleveragingSession) Owner() (common.Address, error) {
	return _P1Deleveraging.Contract.Owner(&_P1Deleveraging.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1Deleveraging *P1DeleveragingCallerSession) Owner() (common.Address, error) {
	return _P1Deleveraging.Contract.Owner(&_P1Deleveraging.CallOpts)
}

// Mark is a paid mutator transaction binding the contract method 0x7ceeb880.
//
// Solidity: function mark(address account) returns()
func (_P1Deleveraging *P1DeleveragingTransactor) Mark(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "mark", account)
}

// Mark is a paid mutator transaction binding the contract method 0x7ceeb880.
//
// Solidity: function mark(address account) returns()
func (_P1Deleveraging *P1DeleveragingSession) Mark(account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Mark(&_P1Deleveraging.TransactOpts, account)
}

// Mark is a paid mutator transaction binding the contract method 0x7ceeb880.
//
// Solidity: function mark(address account) returns()
func (_P1Deleveraging *P1DeleveragingTransactorSession) Mark(account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Mark(&_P1Deleveraging.TransactOpts, account)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1Deleveraging *P1DeleveragingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1Deleveraging *P1DeleveragingSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1Deleveraging.Contract.RenounceOwnership(&_P1Deleveraging.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1Deleveraging *P1DeleveragingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1Deleveraging.Contract.RenounceOwnership(&_P1Deleveraging.TransactOpts)
}

// SetDeleveragingOperator is a paid mutator transaction binding the contract method 0x4ca5dcb2.
//
// Solidity: function setDeleveragingOperator(address newOperator) returns()
func (_P1Deleveraging *P1DeleveragingTransactor) SetDeleveragingOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "setDeleveragingOperator", newOperator)
}

// SetDeleveragingOperator is a paid mutator transaction binding the contract method 0x4ca5dcb2.
//
// Solidity: function setDeleveragingOperator(address newOperator) returns()
func (_P1Deleveraging *P1DeleveragingSession) SetDeleveragingOperator(newOperator common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.SetDeleveragingOperator(&_P1Deleveraging.TransactOpts, newOperator)
}

// SetDeleveragingOperator is a paid mutator transaction binding the contract method 0x4ca5dcb2.
//
// Solidity: function setDeleveragingOperator(address newOperator) returns()
func (_P1Deleveraging *P1DeleveragingTransactorSession) SetDeleveragingOperator(newOperator common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.SetDeleveragingOperator(&_P1Deleveraging.TransactOpts, newOperator)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_P1Deleveraging *P1DeleveragingTransactor) Trade(opts *bind.TransactOpts, sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "trade", sender, maker, taker, price, data, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_P1Deleveraging *P1DeleveragingSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Trade(&_P1Deleveraging.TransactOpts, sender, maker, taker, price, data, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_P1Deleveraging *P1DeleveragingTransactorSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Trade(&_P1Deleveraging.TransactOpts, sender, maker, taker, price, data, traderFlags)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1Deleveraging *P1DeleveragingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1Deleveraging *P1DeleveragingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.TransferOwnership(&_P1Deleveraging.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1Deleveraging *P1DeleveragingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.TransferOwnership(&_P1Deleveraging.TransactOpts, newOwner)
}

// Unmark is a paid mutator transaction binding the contract method 0xa5db6198.
//
// Solidity: function unmark(address account) returns()
func (_P1Deleveraging *P1DeleveragingTransactor) Unmark(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.contract.Transact(opts, "unmark", account)
}

// Unmark is a paid mutator transaction binding the contract method 0xa5db6198.
//
// Solidity: function unmark(address account) returns()
func (_P1Deleveraging *P1DeleveragingSession) Unmark(account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Unmark(&_P1Deleveraging.TransactOpts, account)
}

// Unmark is a paid mutator transaction binding the contract method 0xa5db6198.
//
// Solidity: function unmark(address account) returns()
func (_P1Deleveraging *P1DeleveragingTransactorSession) Unmark(account common.Address) (*types.Transaction, error) {
	return _P1Deleveraging.Contract.Unmark(&_P1Deleveraging.TransactOpts, account)
}

// P1DeleveragingLogDeleveragedIterator is returned from FilterLogDeleveraged and is used to iterate over the raw logs and unpacked data for LogDeleveraged events raised by the P1Deleveraging contract.
type P1DeleveragingLogDeleveragedIterator struct {
	Event *P1DeleveragingLogDeleveraged // Event containing the contract specifics and raw log

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
func (it *P1DeleveragingLogDeleveragedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1DeleveragingLogDeleveraged)
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
		it.Event = new(P1DeleveragingLogDeleveraged)
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
func (it *P1DeleveragingLogDeleveragedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1DeleveragingLogDeleveragedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1DeleveragingLogDeleveraged represents a LogDeleveraged event raised by the P1Deleveraging contract.
type P1DeleveragingLogDeleveraged struct {
	Maker       common.Address
	Taker       common.Address
	Amount      *big.Int
	IsBuy       bool
	OraclePrice *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogDeleveraged is a free log retrieval operation binding the contract event 0x750eb1daab9f5a06890a5126e981abeeb7d50b590d48e4a9e523016de22985bb.
//
// Solidity: event LogDeleveraged(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Deleveraging *P1DeleveragingFilterer) FilterLogDeleveraged(opts *bind.FilterOpts, maker []common.Address, taker []common.Address) (*P1DeleveragingLogDeleveragedIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Deleveraging.contract.FilterLogs(opts, "LogDeleveraged", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingLogDeleveragedIterator{contract: _P1Deleveraging.contract, event: "LogDeleveraged", logs: logs, sub: sub}, nil
}

// WatchLogDeleveraged is a free log subscription operation binding the contract event 0x750eb1daab9f5a06890a5126e981abeeb7d50b590d48e4a9e523016de22985bb.
//
// Solidity: event LogDeleveraged(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Deleveraging *P1DeleveragingFilterer) WatchLogDeleveraged(opts *bind.WatchOpts, sink chan<- *P1DeleveragingLogDeleveraged, maker []common.Address, taker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Deleveraging.contract.WatchLogs(opts, "LogDeleveraged", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1DeleveragingLogDeleveraged)
				if err := _P1Deleveraging.contract.UnpackLog(event, "LogDeleveraged", log); err != nil {
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

// ParseLogDeleveraged is a log parse operation binding the contract event 0x750eb1daab9f5a06890a5126e981abeeb7d50b590d48e4a9e523016de22985bb.
//
// Solidity: event LogDeleveraged(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Deleveraging *P1DeleveragingFilterer) ParseLogDeleveraged(log types.Log) (*P1DeleveragingLogDeleveraged, error) {
	event := new(P1DeleveragingLogDeleveraged)
	if err := _P1Deleveraging.contract.UnpackLog(event, "LogDeleveraged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1DeleveragingLogDeleveragingOperatorSetIterator is returned from FilterLogDeleveragingOperatorSet and is used to iterate over the raw logs and unpacked data for LogDeleveragingOperatorSet events raised by the P1Deleveraging contract.
type P1DeleveragingLogDeleveragingOperatorSetIterator struct {
	Event *P1DeleveragingLogDeleveragingOperatorSet // Event containing the contract specifics and raw log

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
func (it *P1DeleveragingLogDeleveragingOperatorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1DeleveragingLogDeleveragingOperatorSet)
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
		it.Event = new(P1DeleveragingLogDeleveragingOperatorSet)
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
func (it *P1DeleveragingLogDeleveragingOperatorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1DeleveragingLogDeleveragingOperatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1DeleveragingLogDeleveragingOperatorSet represents a LogDeleveragingOperatorSet event raised by the P1Deleveraging contract.
type P1DeleveragingLogDeleveragingOperatorSet struct {
	DeleveragingOperator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogDeleveragingOperatorSet is a free log retrieval operation binding the contract event 0x40f472379384d27e4916dabc7440b7ce7f7282a3e665176fe8fbcaf4c12a90a3.
//
// Solidity: event LogDeleveragingOperatorSet(address deleveragingOperator)
func (_P1Deleveraging *P1DeleveragingFilterer) FilterLogDeleveragingOperatorSet(opts *bind.FilterOpts) (*P1DeleveragingLogDeleveragingOperatorSetIterator, error) {

	logs, sub, err := _P1Deleveraging.contract.FilterLogs(opts, "LogDeleveragingOperatorSet")
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingLogDeleveragingOperatorSetIterator{contract: _P1Deleveraging.contract, event: "LogDeleveragingOperatorSet", logs: logs, sub: sub}, nil
}

// WatchLogDeleveragingOperatorSet is a free log subscription operation binding the contract event 0x40f472379384d27e4916dabc7440b7ce7f7282a3e665176fe8fbcaf4c12a90a3.
//
// Solidity: event LogDeleveragingOperatorSet(address deleveragingOperator)
func (_P1Deleveraging *P1DeleveragingFilterer) WatchLogDeleveragingOperatorSet(opts *bind.WatchOpts, sink chan<- *P1DeleveragingLogDeleveragingOperatorSet) (event.Subscription, error) {

	logs, sub, err := _P1Deleveraging.contract.WatchLogs(opts, "LogDeleveragingOperatorSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1DeleveragingLogDeleveragingOperatorSet)
				if err := _P1Deleveraging.contract.UnpackLog(event, "LogDeleveragingOperatorSet", log); err != nil {
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

// ParseLogDeleveragingOperatorSet is a log parse operation binding the contract event 0x40f472379384d27e4916dabc7440b7ce7f7282a3e665176fe8fbcaf4c12a90a3.
//
// Solidity: event LogDeleveragingOperatorSet(address deleveragingOperator)
func (_P1Deleveraging *P1DeleveragingFilterer) ParseLogDeleveragingOperatorSet(log types.Log) (*P1DeleveragingLogDeleveragingOperatorSet, error) {
	event := new(P1DeleveragingLogDeleveragingOperatorSet)
	if err := _P1Deleveraging.contract.UnpackLog(event, "LogDeleveragingOperatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1DeleveragingLogMarkedForDeleveragingIterator is returned from FilterLogMarkedForDeleveraging and is used to iterate over the raw logs and unpacked data for LogMarkedForDeleveraging events raised by the P1Deleveraging contract.
type P1DeleveragingLogMarkedForDeleveragingIterator struct {
	Event *P1DeleveragingLogMarkedForDeleveraging // Event containing the contract specifics and raw log

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
func (it *P1DeleveragingLogMarkedForDeleveragingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1DeleveragingLogMarkedForDeleveraging)
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
		it.Event = new(P1DeleveragingLogMarkedForDeleveraging)
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
func (it *P1DeleveragingLogMarkedForDeleveragingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1DeleveragingLogMarkedForDeleveragingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1DeleveragingLogMarkedForDeleveraging represents a LogMarkedForDeleveraging event raised by the P1Deleveraging contract.
type P1DeleveragingLogMarkedForDeleveraging struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogMarkedForDeleveraging is a free log retrieval operation binding the contract event 0xc5c00dddf28309e25cd647751e97c5c3b8be3425d81f33ba7835ac8e710aa8b9.
//
// Solidity: event LogMarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) FilterLogMarkedForDeleveraging(opts *bind.FilterOpts, account []common.Address) (*P1DeleveragingLogMarkedForDeleveragingIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Deleveraging.contract.FilterLogs(opts, "LogMarkedForDeleveraging", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingLogMarkedForDeleveragingIterator{contract: _P1Deleveraging.contract, event: "LogMarkedForDeleveraging", logs: logs, sub: sub}, nil
}

// WatchLogMarkedForDeleveraging is a free log subscription operation binding the contract event 0xc5c00dddf28309e25cd647751e97c5c3b8be3425d81f33ba7835ac8e710aa8b9.
//
// Solidity: event LogMarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) WatchLogMarkedForDeleveraging(opts *bind.WatchOpts, sink chan<- *P1DeleveragingLogMarkedForDeleveraging, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Deleveraging.contract.WatchLogs(opts, "LogMarkedForDeleveraging", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1DeleveragingLogMarkedForDeleveraging)
				if err := _P1Deleveraging.contract.UnpackLog(event, "LogMarkedForDeleveraging", log); err != nil {
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

// ParseLogMarkedForDeleveraging is a log parse operation binding the contract event 0xc5c00dddf28309e25cd647751e97c5c3b8be3425d81f33ba7835ac8e710aa8b9.
//
// Solidity: event LogMarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) ParseLogMarkedForDeleveraging(log types.Log) (*P1DeleveragingLogMarkedForDeleveraging, error) {
	event := new(P1DeleveragingLogMarkedForDeleveraging)
	if err := _P1Deleveraging.contract.UnpackLog(event, "LogMarkedForDeleveraging", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1DeleveragingLogUnmarkedForDeleveragingIterator is returned from FilterLogUnmarkedForDeleveraging and is used to iterate over the raw logs and unpacked data for LogUnmarkedForDeleveraging events raised by the P1Deleveraging contract.
type P1DeleveragingLogUnmarkedForDeleveragingIterator struct {
	Event *P1DeleveragingLogUnmarkedForDeleveraging // Event containing the contract specifics and raw log

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
func (it *P1DeleveragingLogUnmarkedForDeleveragingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1DeleveragingLogUnmarkedForDeleveraging)
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
		it.Event = new(P1DeleveragingLogUnmarkedForDeleveraging)
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
func (it *P1DeleveragingLogUnmarkedForDeleveragingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1DeleveragingLogUnmarkedForDeleveragingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1DeleveragingLogUnmarkedForDeleveraging represents a LogUnmarkedForDeleveraging event raised by the P1Deleveraging contract.
type P1DeleveragingLogUnmarkedForDeleveraging struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogUnmarkedForDeleveraging is a free log retrieval operation binding the contract event 0x9771b33aeb53e1bf116f495e5f185f99fd84b523e3da950a84ca36ad10b8160e.
//
// Solidity: event LogUnmarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) FilterLogUnmarkedForDeleveraging(opts *bind.FilterOpts, account []common.Address) (*P1DeleveragingLogUnmarkedForDeleveragingIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Deleveraging.contract.FilterLogs(opts, "LogUnmarkedForDeleveraging", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingLogUnmarkedForDeleveragingIterator{contract: _P1Deleveraging.contract, event: "LogUnmarkedForDeleveraging", logs: logs, sub: sub}, nil
}

// WatchLogUnmarkedForDeleveraging is a free log subscription operation binding the contract event 0x9771b33aeb53e1bf116f495e5f185f99fd84b523e3da950a84ca36ad10b8160e.
//
// Solidity: event LogUnmarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) WatchLogUnmarkedForDeleveraging(opts *bind.WatchOpts, sink chan<- *P1DeleveragingLogUnmarkedForDeleveraging, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Deleveraging.contract.WatchLogs(opts, "LogUnmarkedForDeleveraging", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1DeleveragingLogUnmarkedForDeleveraging)
				if err := _P1Deleveraging.contract.UnpackLog(event, "LogUnmarkedForDeleveraging", log); err != nil {
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

// ParseLogUnmarkedForDeleveraging is a log parse operation binding the contract event 0x9771b33aeb53e1bf116f495e5f185f99fd84b523e3da950a84ca36ad10b8160e.
//
// Solidity: event LogUnmarkedForDeleveraging(address indexed account)
func (_P1Deleveraging *P1DeleveragingFilterer) ParseLogUnmarkedForDeleveraging(log types.Log) (*P1DeleveragingLogUnmarkedForDeleveraging, error) {
	event := new(P1DeleveragingLogUnmarkedForDeleveraging)
	if err := _P1Deleveraging.contract.UnpackLog(event, "LogUnmarkedForDeleveraging", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1DeleveragingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1Deleveraging contract.
type P1DeleveragingOwnershipTransferredIterator struct {
	Event *P1DeleveragingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1DeleveragingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1DeleveragingOwnershipTransferred)
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
		it.Event = new(P1DeleveragingOwnershipTransferred)
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
func (it *P1DeleveragingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1DeleveragingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1DeleveragingOwnershipTransferred represents a OwnershipTransferred event raised by the P1Deleveraging contract.
type P1DeleveragingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1Deleveraging *P1DeleveragingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1DeleveragingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1Deleveraging.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1DeleveragingOwnershipTransferredIterator{contract: _P1Deleveraging.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1Deleveraging *P1DeleveragingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1DeleveragingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1Deleveraging.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1DeleveragingOwnershipTransferred)
				if err := _P1Deleveraging.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_P1Deleveraging *P1DeleveragingFilterer) ParseOwnershipTransferred(log types.Log) (*P1DeleveragingOwnershipTransferred, error) {
	event := new(P1DeleveragingOwnershipTransferred)
	if err := _P1Deleveraging.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
