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

// P1FundingOracleMetaData contains all meta data concerning the P1FundingOracle contract.
var P1FundingOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fundingRateProvider\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fundingRateProvider\",\"type\":\"address\"}],\"name\":\"LogFundingRateProviderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fundingRate\",\"type\":\"bytes32\"}],\"name\":\"LogFundingRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_ABS_DIFF_PER_SECOND\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_ABS_VALUE\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_FUNDING_RATE_PROVIDER_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"newRate\",\"type\":\"tuple\"}],\"name\":\"setFundingRate\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newProvider\",\"type\":\"address\"}],\"name\":\"setFundingRateProvider\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeDelta\",\"type\":\"uint256\"}],\"name\":\"getFunding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1FundingOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use P1FundingOracleMetaData.ABI instead.
var P1FundingOracleABI = P1FundingOracleMetaData.ABI

// P1FundingOracle is an auto generated Go binding around an Ethereum contract.
type P1FundingOracle struct {
	P1FundingOracleCaller     // Read-only binding to the contract
	P1FundingOracleTransactor // Write-only binding to the contract
	P1FundingOracleFilterer   // Log filterer for contract events
}

// P1FundingOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1FundingOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FundingOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1FundingOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FundingOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1FundingOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FundingOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1FundingOracleSession struct {
	Contract     *P1FundingOracle  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1FundingOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1FundingOracleCallerSession struct {
	Contract *P1FundingOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// P1FundingOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1FundingOracleTransactorSession struct {
	Contract     *P1FundingOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// P1FundingOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1FundingOracleRaw struct {
	Contract *P1FundingOracle // Generic contract binding to access the raw methods on
}

// P1FundingOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1FundingOracleCallerRaw struct {
	Contract *P1FundingOracleCaller // Generic read-only contract binding to access the raw methods on
}

// P1FundingOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1FundingOracleTransactorRaw struct {
	Contract *P1FundingOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1FundingOracle creates a new instance of P1FundingOracle, bound to a specific deployed contract.
func NewP1FundingOracle(address common.Address, backend bind.ContractBackend) (*P1FundingOracle, error) {
	contract, err := bindP1FundingOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1FundingOracle{P1FundingOracleCaller: P1FundingOracleCaller{contract: contract}, P1FundingOracleTransactor: P1FundingOracleTransactor{contract: contract}, P1FundingOracleFilterer: P1FundingOracleFilterer{contract: contract}}, nil
}

// NewP1FundingOracleCaller creates a new read-only instance of P1FundingOracle, bound to a specific deployed contract.
func NewP1FundingOracleCaller(address common.Address, caller bind.ContractCaller) (*P1FundingOracleCaller, error) {
	contract, err := bindP1FundingOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleCaller{contract: contract}, nil
}

// NewP1FundingOracleTransactor creates a new write-only instance of P1FundingOracle, bound to a specific deployed contract.
func NewP1FundingOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*P1FundingOracleTransactor, error) {
	contract, err := bindP1FundingOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleTransactor{contract: contract}, nil
}

// NewP1FundingOracleFilterer creates a new log filterer instance of P1FundingOracle, bound to a specific deployed contract.
func NewP1FundingOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*P1FundingOracleFilterer, error) {
	contract, err := bindP1FundingOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleFilterer{contract: contract}, nil
}

// bindP1FundingOracle binds a generic wrapper to an already deployed contract.
func bindP1FundingOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1FundingOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1FundingOracle *P1FundingOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1FundingOracle.Contract.P1FundingOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1FundingOracle *P1FundingOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.P1FundingOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1FundingOracle *P1FundingOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.P1FundingOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1FundingOracle *P1FundingOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1FundingOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1FundingOracle *P1FundingOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1FundingOracle *P1FundingOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.contract.Transact(opts, method, params...)
}

// MAXABSDIFFPERSECOND is a free data retrieval call binding the contract method 0x56a281ea.
//
// Solidity: function MAX_ABS_DIFF_PER_SECOND() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleCaller) MAXABSDIFFPERSECOND(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "MAX_ABS_DIFF_PER_SECOND")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXABSDIFFPERSECOND is a free data retrieval call binding the contract method 0x56a281ea.
//
// Solidity: function MAX_ABS_DIFF_PER_SECOND() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleSession) MAXABSDIFFPERSECOND() (*big.Int, error) {
	return _P1FundingOracle.Contract.MAXABSDIFFPERSECOND(&_P1FundingOracle.CallOpts)
}

// MAXABSDIFFPERSECOND is a free data retrieval call binding the contract method 0x56a281ea.
//
// Solidity: function MAX_ABS_DIFF_PER_SECOND() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleCallerSession) MAXABSDIFFPERSECOND() (*big.Int, error) {
	return _P1FundingOracle.Contract.MAXABSDIFFPERSECOND(&_P1FundingOracle.CallOpts)
}

// MAXABSVALUE is a free data retrieval call binding the contract method 0x499c9c6d.
//
// Solidity: function MAX_ABS_VALUE() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleCaller) MAXABSVALUE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "MAX_ABS_VALUE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXABSVALUE is a free data retrieval call binding the contract method 0x499c9c6d.
//
// Solidity: function MAX_ABS_VALUE() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleSession) MAXABSVALUE() (*big.Int, error) {
	return _P1FundingOracle.Contract.MAXABSVALUE(&_P1FundingOracle.CallOpts)
}

// MAXABSVALUE is a free data retrieval call binding the contract method 0x499c9c6d.
//
// Solidity: function MAX_ABS_VALUE() view returns(uint128)
func (_P1FundingOracle *P1FundingOracleCallerSession) MAXABSVALUE() (*big.Int, error) {
	return _P1FundingOracle.Contract.MAXABSVALUE(&_P1FundingOracle.CallOpts)
}

// FUNDINGRATEPROVIDER is a free data retrieval call binding the contract method 0x0b8781ee.
//
// Solidity: function _FUNDING_RATE_PROVIDER_() view returns(address)
func (_P1FundingOracle *P1FundingOracleCaller) FUNDINGRATEPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "_FUNDING_RATE_PROVIDER_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNDINGRATEPROVIDER is a free data retrieval call binding the contract method 0x0b8781ee.
//
// Solidity: function _FUNDING_RATE_PROVIDER_() view returns(address)
func (_P1FundingOracle *P1FundingOracleSession) FUNDINGRATEPROVIDER() (common.Address, error) {
	return _P1FundingOracle.Contract.FUNDINGRATEPROVIDER(&_P1FundingOracle.CallOpts)
}

// FUNDINGRATEPROVIDER is a free data retrieval call binding the contract method 0x0b8781ee.
//
// Solidity: function _FUNDING_RATE_PROVIDER_() view returns(address)
func (_P1FundingOracle *P1FundingOracleCallerSession) FUNDINGRATEPROVIDER() (common.Address, error) {
	return _P1FundingOracle.Contract.FUNDINGRATEPROVIDER(&_P1FundingOracle.CallOpts)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_P1FundingOracle *P1FundingOracleCaller) GetFunding(opts *bind.CallOpts, timeDelta *big.Int) (bool, *big.Int, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "getFunding", timeDelta)

	if err != nil {
		return *new(bool), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_P1FundingOracle *P1FundingOracleSession) GetFunding(timeDelta *big.Int) (bool, *big.Int, error) {
	return _P1FundingOracle.Contract.GetFunding(&_P1FundingOracle.CallOpts, timeDelta)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_P1FundingOracle *P1FundingOracleCallerSession) GetFunding(timeDelta *big.Int) (bool, *big.Int, error) {
	return _P1FundingOracle.Contract.GetFunding(&_P1FundingOracle.CallOpts, timeDelta)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1FundingOracle *P1FundingOracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1FundingOracle *P1FundingOracleSession) IsOwner() (bool, error) {
	return _P1FundingOracle.Contract.IsOwner(&_P1FundingOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1FundingOracle *P1FundingOracleCallerSession) IsOwner() (bool, error) {
	return _P1FundingOracle.Contract.IsOwner(&_P1FundingOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1FundingOracle *P1FundingOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1FundingOracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1FundingOracle *P1FundingOracleSession) Owner() (common.Address, error) {
	return _P1FundingOracle.Contract.Owner(&_P1FundingOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1FundingOracle *P1FundingOracleCallerSession) Owner() (common.Address, error) {
	return _P1FundingOracle.Contract.Owner(&_P1FundingOracle.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1FundingOracle *P1FundingOracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FundingOracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1FundingOracle *P1FundingOracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1FundingOracle.Contract.RenounceOwnership(&_P1FundingOracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1FundingOracle *P1FundingOracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1FundingOracle.Contract.RenounceOwnership(&_P1FundingOracle.TransactOpts)
}

// SetFundingRate is a paid mutator transaction binding the contract method 0xef460e36.
//
// Solidity: function setFundingRate((uint256,bool) newRate) returns((uint32,bool,uint128))
func (_P1FundingOracle *P1FundingOracleTransactor) SetFundingRate(opts *bind.TransactOpts, newRate SignedMathInt) (*types.Transaction, error) {
	return _P1FundingOracle.contract.Transact(opts, "setFundingRate", newRate)
}

// SetFundingRate is a paid mutator transaction binding the contract method 0xef460e36.
//
// Solidity: function setFundingRate((uint256,bool) newRate) returns((uint32,bool,uint128))
func (_P1FundingOracle *P1FundingOracleSession) SetFundingRate(newRate SignedMathInt) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.SetFundingRate(&_P1FundingOracle.TransactOpts, newRate)
}

// SetFundingRate is a paid mutator transaction binding the contract method 0xef460e36.
//
// Solidity: function setFundingRate((uint256,bool) newRate) returns((uint32,bool,uint128))
func (_P1FundingOracle *P1FundingOracleTransactorSession) SetFundingRate(newRate SignedMathInt) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.SetFundingRate(&_P1FundingOracle.TransactOpts, newRate)
}

// SetFundingRateProvider is a paid mutator transaction binding the contract method 0x109f60e3.
//
// Solidity: function setFundingRateProvider(address newProvider) returns()
func (_P1FundingOracle *P1FundingOracleTransactor) SetFundingRateProvider(opts *bind.TransactOpts, newProvider common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.contract.Transact(opts, "setFundingRateProvider", newProvider)
}

// SetFundingRateProvider is a paid mutator transaction binding the contract method 0x109f60e3.
//
// Solidity: function setFundingRateProvider(address newProvider) returns()
func (_P1FundingOracle *P1FundingOracleSession) SetFundingRateProvider(newProvider common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.SetFundingRateProvider(&_P1FundingOracle.TransactOpts, newProvider)
}

// SetFundingRateProvider is a paid mutator transaction binding the contract method 0x109f60e3.
//
// Solidity: function setFundingRateProvider(address newProvider) returns()
func (_P1FundingOracle *P1FundingOracleTransactorSession) SetFundingRateProvider(newProvider common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.SetFundingRateProvider(&_P1FundingOracle.TransactOpts, newProvider)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1FundingOracle *P1FundingOracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1FundingOracle *P1FundingOracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.TransferOwnership(&_P1FundingOracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1FundingOracle *P1FundingOracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1FundingOracle.Contract.TransferOwnership(&_P1FundingOracle.TransactOpts, newOwner)
}

// P1FundingOracleLogFundingRateProviderSetIterator is returned from FilterLogFundingRateProviderSet and is used to iterate over the raw logs and unpacked data for LogFundingRateProviderSet events raised by the P1FundingOracle contract.
type P1FundingOracleLogFundingRateProviderSetIterator struct {
	Event *P1FundingOracleLogFundingRateProviderSet // Event containing the contract specifics and raw log

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
func (it *P1FundingOracleLogFundingRateProviderSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FundingOracleLogFundingRateProviderSet)
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
		it.Event = new(P1FundingOracleLogFundingRateProviderSet)
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
func (it *P1FundingOracleLogFundingRateProviderSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FundingOracleLogFundingRateProviderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FundingOracleLogFundingRateProviderSet represents a LogFundingRateProviderSet event raised by the P1FundingOracle contract.
type P1FundingOracleLogFundingRateProviderSet struct {
	FundingRateProvider common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterLogFundingRateProviderSet is a free log retrieval operation binding the contract event 0x232d43841005a98dbf929d234a7a8d2c4c570becee067c9c81bcd4e0acd0ab92.
//
// Solidity: event LogFundingRateProviderSet(address fundingRateProvider)
func (_P1FundingOracle *P1FundingOracleFilterer) FilterLogFundingRateProviderSet(opts *bind.FilterOpts) (*P1FundingOracleLogFundingRateProviderSetIterator, error) {

	logs, sub, err := _P1FundingOracle.contract.FilterLogs(opts, "LogFundingRateProviderSet")
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleLogFundingRateProviderSetIterator{contract: _P1FundingOracle.contract, event: "LogFundingRateProviderSet", logs: logs, sub: sub}, nil
}

// WatchLogFundingRateProviderSet is a free log subscription operation binding the contract event 0x232d43841005a98dbf929d234a7a8d2c4c570becee067c9c81bcd4e0acd0ab92.
//
// Solidity: event LogFundingRateProviderSet(address fundingRateProvider)
func (_P1FundingOracle *P1FundingOracleFilterer) WatchLogFundingRateProviderSet(opts *bind.WatchOpts, sink chan<- *P1FundingOracleLogFundingRateProviderSet) (event.Subscription, error) {

	logs, sub, err := _P1FundingOracle.contract.WatchLogs(opts, "LogFundingRateProviderSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FundingOracleLogFundingRateProviderSet)
				if err := _P1FundingOracle.contract.UnpackLog(event, "LogFundingRateProviderSet", log); err != nil {
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

// ParseLogFundingRateProviderSet is a log parse operation binding the contract event 0x232d43841005a98dbf929d234a7a8d2c4c570becee067c9c81bcd4e0acd0ab92.
//
// Solidity: event LogFundingRateProviderSet(address fundingRateProvider)
func (_P1FundingOracle *P1FundingOracleFilterer) ParseLogFundingRateProviderSet(log types.Log) (*P1FundingOracleLogFundingRateProviderSet, error) {
	event := new(P1FundingOracleLogFundingRateProviderSet)
	if err := _P1FundingOracle.contract.UnpackLog(event, "LogFundingRateProviderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1FundingOracleLogFundingRateUpdatedIterator is returned from FilterLogFundingRateUpdated and is used to iterate over the raw logs and unpacked data for LogFundingRateUpdated events raised by the P1FundingOracle contract.
type P1FundingOracleLogFundingRateUpdatedIterator struct {
	Event *P1FundingOracleLogFundingRateUpdated // Event containing the contract specifics and raw log

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
func (it *P1FundingOracleLogFundingRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FundingOracleLogFundingRateUpdated)
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
		it.Event = new(P1FundingOracleLogFundingRateUpdated)
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
func (it *P1FundingOracleLogFundingRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FundingOracleLogFundingRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FundingOracleLogFundingRateUpdated represents a LogFundingRateUpdated event raised by the P1FundingOracle contract.
type P1FundingOracleLogFundingRateUpdated struct {
	FundingRate [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogFundingRateUpdated is a free log retrieval operation binding the contract event 0x2ebf65220b5046a8d9cff102710ef15de0a0bf3709dcc11c3af50abe472e1c22.
//
// Solidity: event LogFundingRateUpdated(bytes32 fundingRate)
func (_P1FundingOracle *P1FundingOracleFilterer) FilterLogFundingRateUpdated(opts *bind.FilterOpts) (*P1FundingOracleLogFundingRateUpdatedIterator, error) {

	logs, sub, err := _P1FundingOracle.contract.FilterLogs(opts, "LogFundingRateUpdated")
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleLogFundingRateUpdatedIterator{contract: _P1FundingOracle.contract, event: "LogFundingRateUpdated", logs: logs, sub: sub}, nil
}

// WatchLogFundingRateUpdated is a free log subscription operation binding the contract event 0x2ebf65220b5046a8d9cff102710ef15de0a0bf3709dcc11c3af50abe472e1c22.
//
// Solidity: event LogFundingRateUpdated(bytes32 fundingRate)
func (_P1FundingOracle *P1FundingOracleFilterer) WatchLogFundingRateUpdated(opts *bind.WatchOpts, sink chan<- *P1FundingOracleLogFundingRateUpdated) (event.Subscription, error) {

	logs, sub, err := _P1FundingOracle.contract.WatchLogs(opts, "LogFundingRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FundingOracleLogFundingRateUpdated)
				if err := _P1FundingOracle.contract.UnpackLog(event, "LogFundingRateUpdated", log); err != nil {
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

// ParseLogFundingRateUpdated is a log parse operation binding the contract event 0x2ebf65220b5046a8d9cff102710ef15de0a0bf3709dcc11c3af50abe472e1c22.
//
// Solidity: event LogFundingRateUpdated(bytes32 fundingRate)
func (_P1FundingOracle *P1FundingOracleFilterer) ParseLogFundingRateUpdated(log types.Log) (*P1FundingOracleLogFundingRateUpdated, error) {
	event := new(P1FundingOracleLogFundingRateUpdated)
	if err := _P1FundingOracle.contract.UnpackLog(event, "LogFundingRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1FundingOracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1FundingOracle contract.
type P1FundingOracleOwnershipTransferredIterator struct {
	Event *P1FundingOracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1FundingOracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FundingOracleOwnershipTransferred)
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
		it.Event = new(P1FundingOracleOwnershipTransferred)
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
func (it *P1FundingOracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FundingOracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FundingOracleOwnershipTransferred represents a OwnershipTransferred event raised by the P1FundingOracle contract.
type P1FundingOracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1FundingOracle *P1FundingOracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1FundingOracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1FundingOracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1FundingOracleOwnershipTransferredIterator{contract: _P1FundingOracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1FundingOracle *P1FundingOracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1FundingOracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1FundingOracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FundingOracleOwnershipTransferred)
				if err := _P1FundingOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_P1FundingOracle *P1FundingOracleFilterer) ParseOwnershipTransferred(log types.Log) (*P1FundingOracleOwnershipTransferred, error) {
	event := new(P1FundingOracleOwnershipTransferred)
	if err := _P1FundingOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
