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

// P1MirrorOracleMetaData contains all meta data concerning the P1MirrorOracle contract.
var P1MirrorOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"age\",\"type\":\"uint256\"}],\"name\":\"LogMedianPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bar\",\"type\":\"uint256\"}],\"name\":\"LogSetBar\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"LogSetReader\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"LogSetSigner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_AGE_\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_BAR_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ORACLE_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_ORCL_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"_SLOT_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"peek\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"read\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"age\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"orcl\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"bud\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"signerId\",\"type\":\"uint8\"}],\"name\":\"slot\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkSynced\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"val_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"age_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"poke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"lift\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"drop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"setBar\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1MirrorOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use P1MirrorOracleMetaData.ABI instead.
var P1MirrorOracleABI = P1MirrorOracleMetaData.ABI

// P1MirrorOracle is an auto generated Go binding around an Ethereum contract.
type P1MirrorOracle struct {
	P1MirrorOracleCaller     // Read-only binding to the contract
	P1MirrorOracleTransactor // Write-only binding to the contract
	P1MirrorOracleFilterer   // Log filterer for contract events
}

// P1MirrorOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1MirrorOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1MirrorOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1MirrorOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1MirrorOracleSession struct {
	Contract     *P1MirrorOracle   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1MirrorOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1MirrorOracleCallerSession struct {
	Contract *P1MirrorOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// P1MirrorOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1MirrorOracleTransactorSession struct {
	Contract     *P1MirrorOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// P1MirrorOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1MirrorOracleRaw struct {
	Contract *P1MirrorOracle // Generic contract binding to access the raw methods on
}

// P1MirrorOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1MirrorOracleCallerRaw struct {
	Contract *P1MirrorOracleCaller // Generic read-only contract binding to access the raw methods on
}

// P1MirrorOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1MirrorOracleTransactorRaw struct {
	Contract *P1MirrorOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1MirrorOracle creates a new instance of P1MirrorOracle, bound to a specific deployed contract.
func NewP1MirrorOracle(address common.Address, backend bind.ContractBackend) (*P1MirrorOracle, error) {
	contract, err := bindP1MirrorOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracle{P1MirrorOracleCaller: P1MirrorOracleCaller{contract: contract}, P1MirrorOracleTransactor: P1MirrorOracleTransactor{contract: contract}, P1MirrorOracleFilterer: P1MirrorOracleFilterer{contract: contract}}, nil
}

// NewP1MirrorOracleCaller creates a new read-only instance of P1MirrorOracle, bound to a specific deployed contract.
func NewP1MirrorOracleCaller(address common.Address, caller bind.ContractCaller) (*P1MirrorOracleCaller, error) {
	contract, err := bindP1MirrorOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleCaller{contract: contract}, nil
}

// NewP1MirrorOracleTransactor creates a new write-only instance of P1MirrorOracle, bound to a specific deployed contract.
func NewP1MirrorOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*P1MirrorOracleTransactor, error) {
	contract, err := bindP1MirrorOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleTransactor{contract: contract}, nil
}

// NewP1MirrorOracleFilterer creates a new log filterer instance of P1MirrorOracle, bound to a specific deployed contract.
func NewP1MirrorOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*P1MirrorOracleFilterer, error) {
	contract, err := bindP1MirrorOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleFilterer{contract: contract}, nil
}

// bindP1MirrorOracle binds a generic wrapper to an already deployed contract.
func bindP1MirrorOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1MirrorOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MirrorOracle *P1MirrorOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MirrorOracle.Contract.P1MirrorOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MirrorOracle *P1MirrorOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.P1MirrorOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MirrorOracle *P1MirrorOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.P1MirrorOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MirrorOracle *P1MirrorOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MirrorOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MirrorOracle *P1MirrorOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MirrorOracle *P1MirrorOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.contract.Transact(opts, method, params...)
}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleCaller) AGE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "_AGE_")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleSession) AGE() (uint32, error) {
	return _P1MirrorOracle.Contract.AGE(&_P1MirrorOracle.CallOpts)
}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) AGE() (uint32, error) {
	return _P1MirrorOracle.Contract.AGE(&_P1MirrorOracle.CallOpts)
}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCaller) BAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "_BAR_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleSession) BAR() (*big.Int, error) {
	return _P1MirrorOracle.Contract.BAR(&_P1MirrorOracle.CallOpts)
}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) BAR() (*big.Int, error) {
	return _P1MirrorOracle.Contract.BAR(&_P1MirrorOracle.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "_ORACLE_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleSession) ORACLE() (common.Address, error) {
	return _P1MirrorOracle.Contract.ORACLE(&_P1MirrorOracle.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) ORACLE() (common.Address, error) {
	return _P1MirrorOracle.Contract.ORACLE(&_P1MirrorOracle.CallOpts)
}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCaller) ORCL(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "_ORCL_", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleSession) ORCL(arg0 common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.ORCL(&_P1MirrorOracle.CallOpts, arg0)
}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) ORCL(arg0 common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.ORCL(&_P1MirrorOracle.CallOpts, arg0)
}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCaller) SLOT(opts *bind.CallOpts, arg0 uint8) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "_SLOT_", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleSession) SLOT(arg0 uint8) (common.Address, error) {
	return _P1MirrorOracle.Contract.SLOT(&_P1MirrorOracle.CallOpts, arg0)
}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) SLOT(arg0 uint8) (common.Address, error) {
	return _P1MirrorOracle.Contract.SLOT(&_P1MirrorOracle.CallOpts, arg0)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleCaller) Age(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "age")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleSession) Age() (uint32, error) {
	return _P1MirrorOracle.Contract.Age(&_P1MirrorOracle.CallOpts)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Age() (uint32, error) {
	return _P1MirrorOracle.Contract.Age(&_P1MirrorOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleSession) Bar() (*big.Int, error) {
	return _P1MirrorOracle.Contract.Bar(&_P1MirrorOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Bar() (*big.Int, error) {
	return _P1MirrorOracle.Contract.Bar(&_P1MirrorOracle.CallOpts)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCaller) Bud(opts *bind.CallOpts, reader common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "bud", reader)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleSession) Bud(reader common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.Bud(&_P1MirrorOracle.CallOpts, reader)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Bud(reader common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.Bud(&_P1MirrorOracle.CallOpts, reader)
}

// CheckSynced is a free data retrieval call binding the contract method 0xaff85a4b.
//
// Solidity: function checkSynced() view returns(uint256, uint256, bool)
func (_P1MirrorOracle *P1MirrorOracleCaller) CheckSynced(opts *bind.CallOpts) (*big.Int, *big.Int, bool, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "checkSynced")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(bool)).(*bool)

	return out0, out1, out2, err

}

// CheckSynced is a free data retrieval call binding the contract method 0xaff85a4b.
//
// Solidity: function checkSynced() view returns(uint256, uint256, bool)
func (_P1MirrorOracle *P1MirrorOracleSession) CheckSynced() (*big.Int, *big.Int, bool, error) {
	return _P1MirrorOracle.Contract.CheckSynced(&_P1MirrorOracle.CallOpts)
}

// CheckSynced is a free data retrieval call binding the contract method 0xaff85a4b.
//
// Solidity: function checkSynced() view returns(uint256, uint256, bool)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) CheckSynced() (*big.Int, *big.Int, bool, error) {
	return _P1MirrorOracle.Contract.CheckSynced(&_P1MirrorOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracle *P1MirrorOracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracle *P1MirrorOracleSession) IsOwner() (bool, error) {
	return _P1MirrorOracle.Contract.IsOwner(&_P1MirrorOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) IsOwner() (bool, error) {
	return _P1MirrorOracle.Contract.IsOwner(&_P1MirrorOracle.CallOpts)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCaller) Orcl(opts *bind.CallOpts, signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "orcl", signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleSession) Orcl(signer common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.Orcl(&_P1MirrorOracle.CallOpts, signer)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Orcl(signer common.Address) (*big.Int, error) {
	return _P1MirrorOracle.Contract.Orcl(&_P1MirrorOracle.CallOpts, signer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleSession) Owner() (common.Address, error) {
	return _P1MirrorOracle.Contract.Owner(&_P1MirrorOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Owner() (common.Address, error) {
	return _P1MirrorOracle.Contract.Owner(&_P1MirrorOracle.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_P1MirrorOracle *P1MirrorOracleCaller) Peek(opts *bind.CallOpts) ([32]byte, bool, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "peek")

	if err != nil {
		return *new([32]byte), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_P1MirrorOracle *P1MirrorOracleSession) Peek() ([32]byte, bool, error) {
	return _P1MirrorOracle.Contract.Peek(&_P1MirrorOracle.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Peek() ([32]byte, bool, error) {
	return _P1MirrorOracle.Contract.Peek(&_P1MirrorOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracle *P1MirrorOracleCaller) Read(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "read")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracle *P1MirrorOracleSession) Read() ([32]byte, error) {
	return _P1MirrorOracle.Contract.Read(&_P1MirrorOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Read() ([32]byte, error) {
	return _P1MirrorOracle.Contract.Read(&_P1MirrorOracle.CallOpts)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCaller) Slot(opts *bind.CallOpts, signerId uint8) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracle.contract.Call(opts, &out, "slot", signerId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleSession) Slot(signerId uint8) (common.Address, error) {
	return _P1MirrorOracle.Contract.Slot(&_P1MirrorOracle.CallOpts, signerId)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracle *P1MirrorOracleCallerSession) Slot(signerId uint8) (common.Address, error) {
	return _P1MirrorOracle.Contract.Slot(&_P1MirrorOracle.CallOpts, signerId)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Diss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "diss", readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Diss(&_P1MirrorOracle.TransactOpts, readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Diss(&_P1MirrorOracle.TransactOpts, readers)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Diss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "diss0", reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Diss0(&_P1MirrorOracle.TransactOpts, reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Diss0(&_P1MirrorOracle.TransactOpts, reader)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Drop(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "drop", signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Drop(&_P1MirrorOracle.TransactOpts, signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Drop(&_P1MirrorOracle.TransactOpts, signers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Kiss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "kiss", readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Kiss(&_P1MirrorOracle.TransactOpts, readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Kiss(&_P1MirrorOracle.TransactOpts, readers)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Kiss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "kiss0", reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Kiss0(&_P1MirrorOracle.TransactOpts, reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Kiss0(&_P1MirrorOracle.TransactOpts, reader)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Lift(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "lift", signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Lift(&_P1MirrorOracle.TransactOpts, signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Lift(&_P1MirrorOracle.TransactOpts, signers)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) Poke(opts *bind.TransactOpts, val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "poke", val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Poke(&_P1MirrorOracle.TransactOpts, val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.Poke(&_P1MirrorOracle.TransactOpts, val_, age_, v, r, s)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracle *P1MirrorOracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.RenounceOwnership(&_P1MirrorOracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.RenounceOwnership(&_P1MirrorOracle.TransactOpts)
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) SetBar(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "setBar")
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracle *P1MirrorOracleSession) SetBar() (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.SetBar(&_P1MirrorOracle.TransactOpts)
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) SetBar() (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.SetBar(&_P1MirrorOracle.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracle *P1MirrorOracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.TransferOwnership(&_P1MirrorOracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracle *P1MirrorOracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracle.Contract.TransferOwnership(&_P1MirrorOracle.TransactOpts, newOwner)
}

// P1MirrorOracleLogMedianPriceIterator is returned from FilterLogMedianPrice and is used to iterate over the raw logs and unpacked data for LogMedianPrice events raised by the P1MirrorOracle contract.
type P1MirrorOracleLogMedianPriceIterator struct {
	Event *P1MirrorOracleLogMedianPrice // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleLogMedianPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleLogMedianPrice)
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
		it.Event = new(P1MirrorOracleLogMedianPrice)
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
func (it *P1MirrorOracleLogMedianPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleLogMedianPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleLogMedianPrice represents a LogMedianPrice event raised by the P1MirrorOracle contract.
type P1MirrorOracleLogMedianPrice struct {
	Val *big.Int
	Age *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogMedianPrice is a free log retrieval operation binding the contract event 0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c.
//
// Solidity: event LogMedianPrice(uint256 val, uint256 age)
func (_P1MirrorOracle *P1MirrorOracleFilterer) FilterLogMedianPrice(opts *bind.FilterOpts) (*P1MirrorOracleLogMedianPriceIterator, error) {

	logs, sub, err := _P1MirrorOracle.contract.FilterLogs(opts, "LogMedianPrice")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleLogMedianPriceIterator{contract: _P1MirrorOracle.contract, event: "LogMedianPrice", logs: logs, sub: sub}, nil
}

// WatchLogMedianPrice is a free log subscription operation binding the contract event 0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c.
//
// Solidity: event LogMedianPrice(uint256 val, uint256 age)
func (_P1MirrorOracle *P1MirrorOracleFilterer) WatchLogMedianPrice(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleLogMedianPrice) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracle.contract.WatchLogs(opts, "LogMedianPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleLogMedianPrice)
				if err := _P1MirrorOracle.contract.UnpackLog(event, "LogMedianPrice", log); err != nil {
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

// ParseLogMedianPrice is a log parse operation binding the contract event 0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c.
//
// Solidity: event LogMedianPrice(uint256 val, uint256 age)
func (_P1MirrorOracle *P1MirrorOracleFilterer) ParseLogMedianPrice(log types.Log) (*P1MirrorOracleLogMedianPrice, error) {
	event := new(P1MirrorOracleLogMedianPrice)
	if err := _P1MirrorOracle.contract.UnpackLog(event, "LogMedianPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleLogSetBarIterator is returned from FilterLogSetBar and is used to iterate over the raw logs and unpacked data for LogSetBar events raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetBarIterator struct {
	Event *P1MirrorOracleLogSetBar // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleLogSetBarIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleLogSetBar)
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
		it.Event = new(P1MirrorOracleLogSetBar)
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
func (it *P1MirrorOracleLogSetBarIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleLogSetBarIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleLogSetBar represents a LogSetBar event raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetBar struct {
	Bar *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogSetBar is a free log retrieval operation binding the contract event 0x48c6ae1362d7627f13b4207e5f5cd2724aaac090cb9602e9e8aefe15eb8f24a6.
//
// Solidity: event LogSetBar(uint256 bar)
func (_P1MirrorOracle *P1MirrorOracleFilterer) FilterLogSetBar(opts *bind.FilterOpts) (*P1MirrorOracleLogSetBarIterator, error) {

	logs, sub, err := _P1MirrorOracle.contract.FilterLogs(opts, "LogSetBar")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleLogSetBarIterator{contract: _P1MirrorOracle.contract, event: "LogSetBar", logs: logs, sub: sub}, nil
}

// WatchLogSetBar is a free log subscription operation binding the contract event 0x48c6ae1362d7627f13b4207e5f5cd2724aaac090cb9602e9e8aefe15eb8f24a6.
//
// Solidity: event LogSetBar(uint256 bar)
func (_P1MirrorOracle *P1MirrorOracleFilterer) WatchLogSetBar(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleLogSetBar) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracle.contract.WatchLogs(opts, "LogSetBar")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleLogSetBar)
				if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetBar", log); err != nil {
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

// ParseLogSetBar is a log parse operation binding the contract event 0x48c6ae1362d7627f13b4207e5f5cd2724aaac090cb9602e9e8aefe15eb8f24a6.
//
// Solidity: event LogSetBar(uint256 bar)
func (_P1MirrorOracle *P1MirrorOracleFilterer) ParseLogSetBar(log types.Log) (*P1MirrorOracleLogSetBar, error) {
	event := new(P1MirrorOracleLogSetBar)
	if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetBar", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleLogSetReaderIterator is returned from FilterLogSetReader and is used to iterate over the raw logs and unpacked data for LogSetReader events raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetReaderIterator struct {
	Event *P1MirrorOracleLogSetReader // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleLogSetReaderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleLogSetReader)
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
		it.Event = new(P1MirrorOracleLogSetReader)
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
func (it *P1MirrorOracleLogSetReaderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleLogSetReaderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleLogSetReader represents a LogSetReader event raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetReader struct {
	Reader     common.Address
	Authorized bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogSetReader is a free log retrieval operation binding the contract event 0xadb3d91f6b7a78ea487b119a89fd644a0e6cf0909aa48faff97d153e0df682c0.
//
// Solidity: event LogSetReader(address reader, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) FilterLogSetReader(opts *bind.FilterOpts) (*P1MirrorOracleLogSetReaderIterator, error) {

	logs, sub, err := _P1MirrorOracle.contract.FilterLogs(opts, "LogSetReader")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleLogSetReaderIterator{contract: _P1MirrorOracle.contract, event: "LogSetReader", logs: logs, sub: sub}, nil
}

// WatchLogSetReader is a free log subscription operation binding the contract event 0xadb3d91f6b7a78ea487b119a89fd644a0e6cf0909aa48faff97d153e0df682c0.
//
// Solidity: event LogSetReader(address reader, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) WatchLogSetReader(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleLogSetReader) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracle.contract.WatchLogs(opts, "LogSetReader")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleLogSetReader)
				if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetReader", log); err != nil {
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

// ParseLogSetReader is a log parse operation binding the contract event 0xadb3d91f6b7a78ea487b119a89fd644a0e6cf0909aa48faff97d153e0df682c0.
//
// Solidity: event LogSetReader(address reader, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) ParseLogSetReader(log types.Log) (*P1MirrorOracleLogSetReader, error) {
	event := new(P1MirrorOracleLogSetReader)
	if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetReader", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleLogSetSignerIterator is returned from FilterLogSetSigner and is used to iterate over the raw logs and unpacked data for LogSetSigner events raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetSignerIterator struct {
	Event *P1MirrorOracleLogSetSigner // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleLogSetSignerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleLogSetSigner)
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
		it.Event = new(P1MirrorOracleLogSetSigner)
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
func (it *P1MirrorOracleLogSetSignerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleLogSetSignerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleLogSetSigner represents a LogSetSigner event raised by the P1MirrorOracle contract.
type P1MirrorOracleLogSetSigner struct {
	Signer     common.Address
	Authorized bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogSetSigner is a free log retrieval operation binding the contract event 0x8700965646f22bb776d5e0cbb11e1559f8143228405160e62250b866e954d912.
//
// Solidity: event LogSetSigner(address signer, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) FilterLogSetSigner(opts *bind.FilterOpts) (*P1MirrorOracleLogSetSignerIterator, error) {

	logs, sub, err := _P1MirrorOracle.contract.FilterLogs(opts, "LogSetSigner")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleLogSetSignerIterator{contract: _P1MirrorOracle.contract, event: "LogSetSigner", logs: logs, sub: sub}, nil
}

// WatchLogSetSigner is a free log subscription operation binding the contract event 0x8700965646f22bb776d5e0cbb11e1559f8143228405160e62250b866e954d912.
//
// Solidity: event LogSetSigner(address signer, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) WatchLogSetSigner(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleLogSetSigner) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracle.contract.WatchLogs(opts, "LogSetSigner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleLogSetSigner)
				if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetSigner", log); err != nil {
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

// ParseLogSetSigner is a log parse operation binding the contract event 0x8700965646f22bb776d5e0cbb11e1559f8143228405160e62250b866e954d912.
//
// Solidity: event LogSetSigner(address signer, bool authorized)
func (_P1MirrorOracle *P1MirrorOracleFilterer) ParseLogSetSigner(log types.Log) (*P1MirrorOracleLogSetSigner, error) {
	event := new(P1MirrorOracleLogSetSigner)
	if err := _P1MirrorOracle.contract.UnpackLog(event, "LogSetSigner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1MirrorOracle contract.
type P1MirrorOracleOwnershipTransferredIterator struct {
	Event *P1MirrorOracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleOwnershipTransferred)
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
		it.Event = new(P1MirrorOracleOwnershipTransferred)
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
func (it *P1MirrorOracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleOwnershipTransferred represents a OwnershipTransferred event raised by the P1MirrorOracle contract.
type P1MirrorOracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MirrorOracle *P1MirrorOracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1MirrorOracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MirrorOracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleOwnershipTransferredIterator{contract: _P1MirrorOracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MirrorOracle *P1MirrorOracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MirrorOracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleOwnershipTransferred)
				if err := _P1MirrorOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_P1MirrorOracle *P1MirrorOracleFilterer) ParseOwnershipTransferred(log types.Log) (*P1MirrorOracleOwnershipTransferred, error) {
	event := new(P1MirrorOracleOwnershipTransferred)
	if err := _P1MirrorOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
