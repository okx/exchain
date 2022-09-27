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

// P1MirrorOracleETHUSDMetaData contains all meta data concerning the P1MirrorOracleETHUSD contract.
var P1MirrorOracleETHUSDMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"age\",\"type\":\"uint256\"}],\"name\":\"LogMedianPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bar\",\"type\":\"uint256\"}],\"name\":\"LogSetBar\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"LogSetReader\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"LogSetSigner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"WAT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_AGE_\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_BAR_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ORACLE_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_ORCL_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"_SLOT_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"age\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"bud\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkSynced\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"drop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"lift\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"orcl\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"peek\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"val_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"age_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"poke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"read\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"setBar\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"signerId\",\"type\":\"uint8\"}],\"name\":\"slot\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1MirrorOracleETHUSDABI is the input ABI used to generate the binding from.
// Deprecated: Use P1MirrorOracleETHUSDMetaData.ABI instead.
var P1MirrorOracleETHUSDABI = P1MirrorOracleETHUSDMetaData.ABI

// P1MirrorOracleETHUSD is an auto generated Go binding around an Ethereum contract.
type P1MirrorOracleETHUSD struct {
	P1MirrorOracleETHUSDCaller     // Read-only binding to the contract
	P1MirrorOracleETHUSDTransactor // Write-only binding to the contract
	P1MirrorOracleETHUSDFilterer   // Log filterer for contract events
}

// P1MirrorOracleETHUSDCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1MirrorOracleETHUSDCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleETHUSDTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1MirrorOracleETHUSDTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleETHUSDFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1MirrorOracleETHUSDFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MirrorOracleETHUSDSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1MirrorOracleETHUSDSession struct {
	Contract     *P1MirrorOracleETHUSD // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// P1MirrorOracleETHUSDCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1MirrorOracleETHUSDCallerSession struct {
	Contract *P1MirrorOracleETHUSDCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// P1MirrorOracleETHUSDTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1MirrorOracleETHUSDTransactorSession struct {
	Contract     *P1MirrorOracleETHUSDTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// P1MirrorOracleETHUSDRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1MirrorOracleETHUSDRaw struct {
	Contract *P1MirrorOracleETHUSD // Generic contract binding to access the raw methods on
}

// P1MirrorOracleETHUSDCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1MirrorOracleETHUSDCallerRaw struct {
	Contract *P1MirrorOracleETHUSDCaller // Generic read-only contract binding to access the raw methods on
}

// P1MirrorOracleETHUSDTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1MirrorOracleETHUSDTransactorRaw struct {
	Contract *P1MirrorOracleETHUSDTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1MirrorOracleETHUSD creates a new instance of P1MirrorOracleETHUSD, bound to a specific deployed contract.
func NewP1MirrorOracleETHUSD(address common.Address, backend bind.ContractBackend) (*P1MirrorOracleETHUSD, error) {
	contract, err := bindP1MirrorOracleETHUSD(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSD{P1MirrorOracleETHUSDCaller: P1MirrorOracleETHUSDCaller{contract: contract}, P1MirrorOracleETHUSDTransactor: P1MirrorOracleETHUSDTransactor{contract: contract}, P1MirrorOracleETHUSDFilterer: P1MirrorOracleETHUSDFilterer{contract: contract}}, nil
}

// NewP1MirrorOracleETHUSDCaller creates a new read-only instance of P1MirrorOracleETHUSD, bound to a specific deployed contract.
func NewP1MirrorOracleETHUSDCaller(address common.Address, caller bind.ContractCaller) (*P1MirrorOracleETHUSDCaller, error) {
	contract, err := bindP1MirrorOracleETHUSD(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDCaller{contract: contract}, nil
}

// NewP1MirrorOracleETHUSDTransactor creates a new write-only instance of P1MirrorOracleETHUSD, bound to a specific deployed contract.
func NewP1MirrorOracleETHUSDTransactor(address common.Address, transactor bind.ContractTransactor) (*P1MirrorOracleETHUSDTransactor, error) {
	contract, err := bindP1MirrorOracleETHUSD(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDTransactor{contract: contract}, nil
}

// NewP1MirrorOracleETHUSDFilterer creates a new log filterer instance of P1MirrorOracleETHUSD, bound to a specific deployed contract.
func NewP1MirrorOracleETHUSDFilterer(address common.Address, filterer bind.ContractFilterer) (*P1MirrorOracleETHUSDFilterer, error) {
	contract, err := bindP1MirrorOracleETHUSD(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDFilterer{contract: contract}, nil
}

// bindP1MirrorOracleETHUSD binds a generic wrapper to an already deployed contract.
func bindP1MirrorOracleETHUSD(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1MirrorOracleETHUSDABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MirrorOracleETHUSD.Contract.P1MirrorOracleETHUSDCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.P1MirrorOracleETHUSDTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.P1MirrorOracleETHUSDTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MirrorOracleETHUSD.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.contract.Transact(opts, method, params...)
}

// WAT is a free data retrieval call binding the contract method 0x4e7d3422.
//
// Solidity: function WAT() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) WAT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "WAT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WAT is a free data retrieval call binding the contract method 0x4e7d3422.
//
// Solidity: function WAT() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) WAT() ([32]byte, error) {
	return _P1MirrorOracleETHUSD.Contract.WAT(&_P1MirrorOracleETHUSD.CallOpts)
}

// WAT is a free data retrieval call binding the contract method 0x4e7d3422.
//
// Solidity: function WAT() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) WAT() ([32]byte, error) {
	return _P1MirrorOracleETHUSD.Contract.WAT(&_P1MirrorOracleETHUSD.CallOpts)
}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) AGE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "_AGE_")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) AGE() (uint32, error) {
	return _P1MirrorOracleETHUSD.Contract.AGE(&_P1MirrorOracleETHUSD.CallOpts)
}

// AGE is a free data retrieval call binding the contract method 0xe2f1028e.
//
// Solidity: function _AGE_() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) AGE() (uint32, error) {
	return _P1MirrorOracleETHUSD.Contract.AGE(&_P1MirrorOracleETHUSD.CallOpts)
}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) BAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "_BAR_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) BAR() (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.BAR(&_P1MirrorOracleETHUSD.CallOpts)
}

// BAR is a free data retrieval call binding the contract method 0x82bdfc35.
//
// Solidity: function _BAR_() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) BAR() (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.BAR(&_P1MirrorOracleETHUSD.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "_ORACLE_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) ORACLE() (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.ORACLE(&_P1MirrorOracleETHUSD.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) ORACLE() (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.ORACLE(&_P1MirrorOracleETHUSD.CallOpts)
}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) ORCL(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "_ORCL_", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) ORCL(arg0 common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.ORCL(&_P1MirrorOracleETHUSD.CallOpts, arg0)
}

// ORCL is a free data retrieval call binding the contract method 0x8f8d10bb.
//
// Solidity: function _ORCL_(address ) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) ORCL(arg0 common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.ORCL(&_P1MirrorOracleETHUSD.CallOpts, arg0)
}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) SLOT(opts *bind.CallOpts, arg0 uint8) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "_SLOT_", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) SLOT(arg0 uint8) (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.SLOT(&_P1MirrorOracleETHUSD.CallOpts, arg0)
}

// SLOT is a free data retrieval call binding the contract method 0x1006b5d7.
//
// Solidity: function _SLOT_(uint8 ) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) SLOT(arg0 uint8) (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.SLOT(&_P1MirrorOracleETHUSD.CallOpts, arg0)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Age(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "age")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Age() (uint32, error) {
	return _P1MirrorOracleETHUSD.Contract.Age(&_P1MirrorOracleETHUSD.CallOpts)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Age() (uint32, error) {
	return _P1MirrorOracleETHUSD.Contract.Age(&_P1MirrorOracleETHUSD.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Bar() (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Bar(&_P1MirrorOracleETHUSD.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Bar() (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Bar(&_P1MirrorOracleETHUSD.CallOpts)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Bud(opts *bind.CallOpts, reader common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "bud", reader)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Bud(reader common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Bud(&_P1MirrorOracleETHUSD.CallOpts, reader)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Bud(reader common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Bud(&_P1MirrorOracleETHUSD.CallOpts, reader)
}

// CheckSynced is a free data retrieval call binding the contract method 0xaff85a4b.
//
// Solidity: function checkSynced() view returns(uint256, uint256, bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) CheckSynced(opts *bind.CallOpts) (*big.Int, *big.Int, bool, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "checkSynced")

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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) CheckSynced() (*big.Int, *big.Int, bool, error) {
	return _P1MirrorOracleETHUSD.Contract.CheckSynced(&_P1MirrorOracleETHUSD.CallOpts)
}

// CheckSynced is a free data retrieval call binding the contract method 0xaff85a4b.
//
// Solidity: function checkSynced() view returns(uint256, uint256, bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) CheckSynced() (*big.Int, *big.Int, bool, error) {
	return _P1MirrorOracleETHUSD.Contract.CheckSynced(&_P1MirrorOracleETHUSD.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) IsOwner() (bool, error) {
	return _P1MirrorOracleETHUSD.Contract.IsOwner(&_P1MirrorOracleETHUSD.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) IsOwner() (bool, error) {
	return _P1MirrorOracleETHUSD.Contract.IsOwner(&_P1MirrorOracleETHUSD.CallOpts)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Orcl(opts *bind.CallOpts, signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "orcl", signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Orcl(signer common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Orcl(&_P1MirrorOracleETHUSD.CallOpts, signer)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Orcl(signer common.Address) (*big.Int, error) {
	return _P1MirrorOracleETHUSD.Contract.Orcl(&_P1MirrorOracleETHUSD.CallOpts, signer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Owner() (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.Owner(&_P1MirrorOracleETHUSD.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Owner() (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.Owner(&_P1MirrorOracleETHUSD.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Peek(opts *bind.CallOpts) ([32]byte, bool, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "peek")

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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Peek() ([32]byte, bool, error) {
	return _P1MirrorOracleETHUSD.Contract.Peek(&_P1MirrorOracleETHUSD.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Peek() ([32]byte, bool, error) {
	return _P1MirrorOracleETHUSD.Contract.Peek(&_P1MirrorOracleETHUSD.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Read(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "read")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Read() ([32]byte, error) {
	return _P1MirrorOracleETHUSD.Contract.Read(&_P1MirrorOracleETHUSD.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Read() ([32]byte, error) {
	return _P1MirrorOracleETHUSD.Contract.Read(&_P1MirrorOracleETHUSD.CallOpts)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCaller) Slot(opts *bind.CallOpts, signerId uint8) (common.Address, error) {
	var out []interface{}
	err := _P1MirrorOracleETHUSD.contract.Call(opts, &out, "slot", signerId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Slot(signerId uint8) (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.Slot(&_P1MirrorOracleETHUSD.CallOpts, signerId)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDCallerSession) Slot(signerId uint8) (common.Address, error) {
	return _P1MirrorOracleETHUSD.Contract.Slot(&_P1MirrorOracleETHUSD.CallOpts, signerId)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Diss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "diss", readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Diss(&_P1MirrorOracleETHUSD.TransactOpts, readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Diss(&_P1MirrorOracleETHUSD.TransactOpts, readers)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Diss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "diss0", reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Diss0(&_P1MirrorOracleETHUSD.TransactOpts, reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Diss0(&_P1MirrorOracleETHUSD.TransactOpts, reader)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Drop(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "drop", signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Drop(&_P1MirrorOracleETHUSD.TransactOpts, signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Drop(&_P1MirrorOracleETHUSD.TransactOpts, signers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Kiss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "kiss", readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Kiss(&_P1MirrorOracleETHUSD.TransactOpts, readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Kiss(&_P1MirrorOracleETHUSD.TransactOpts, readers)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Kiss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "kiss0", reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Kiss0(&_P1MirrorOracleETHUSD.TransactOpts, reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Kiss0(&_P1MirrorOracleETHUSD.TransactOpts, reader)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Lift(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "lift", signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Lift(&_P1MirrorOracleETHUSD.TransactOpts, signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Lift(&_P1MirrorOracleETHUSD.TransactOpts, signers)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) Poke(opts *bind.TransactOpts, val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "poke", val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Poke(&_P1MirrorOracleETHUSD.TransactOpts, val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.Poke(&_P1MirrorOracleETHUSD.TransactOpts, val_, age_, v, r, s)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.RenounceOwnership(&_P1MirrorOracleETHUSD.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.RenounceOwnership(&_P1MirrorOracleETHUSD.TransactOpts)
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) SetBar(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "setBar")
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) SetBar() (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.SetBar(&_P1MirrorOracleETHUSD.TransactOpts)
}

// SetBar is a paid mutator transaction binding the contract method 0x24a904b5.
//
// Solidity: function setBar() returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) SetBar() (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.SetBar(&_P1MirrorOracleETHUSD.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.TransferOwnership(&_P1MirrorOracleETHUSD.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MirrorOracleETHUSD.Contract.TransferOwnership(&_P1MirrorOracleETHUSD.TransactOpts, newOwner)
}

// P1MirrorOracleETHUSDLogMedianPriceIterator is returned from FilterLogMedianPrice and is used to iterate over the raw logs and unpacked data for LogMedianPrice events raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogMedianPriceIterator struct {
	Event *P1MirrorOracleETHUSDLogMedianPrice // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleETHUSDLogMedianPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleETHUSDLogMedianPrice)
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
		it.Event = new(P1MirrorOracleETHUSDLogMedianPrice)
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
func (it *P1MirrorOracleETHUSDLogMedianPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleETHUSDLogMedianPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleETHUSDLogMedianPrice represents a LogMedianPrice event raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogMedianPrice struct {
	Val *big.Int
	Age *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogMedianPrice is a free log retrieval operation binding the contract event 0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c.
//
// Solidity: event LogMedianPrice(uint256 val, uint256 age)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) FilterLogMedianPrice(opts *bind.FilterOpts) (*P1MirrorOracleETHUSDLogMedianPriceIterator, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.FilterLogs(opts, "LogMedianPrice")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDLogMedianPriceIterator{contract: _P1MirrorOracleETHUSD.contract, event: "LogMedianPrice", logs: logs, sub: sub}, nil
}

// WatchLogMedianPrice is a free log subscription operation binding the contract event 0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c.
//
// Solidity: event LogMedianPrice(uint256 val, uint256 age)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) WatchLogMedianPrice(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleETHUSDLogMedianPrice) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.WatchLogs(opts, "LogMedianPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleETHUSDLogMedianPrice)
				if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogMedianPrice", log); err != nil {
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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) ParseLogMedianPrice(log types.Log) (*P1MirrorOracleETHUSDLogMedianPrice, error) {
	event := new(P1MirrorOracleETHUSDLogMedianPrice)
	if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogMedianPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleETHUSDLogSetBarIterator is returned from FilterLogSetBar and is used to iterate over the raw logs and unpacked data for LogSetBar events raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetBarIterator struct {
	Event *P1MirrorOracleETHUSDLogSetBar // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleETHUSDLogSetBarIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleETHUSDLogSetBar)
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
		it.Event = new(P1MirrorOracleETHUSDLogSetBar)
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
func (it *P1MirrorOracleETHUSDLogSetBarIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleETHUSDLogSetBarIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleETHUSDLogSetBar represents a LogSetBar event raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetBar struct {
	Bar *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogSetBar is a free log retrieval operation binding the contract event 0x48c6ae1362d7627f13b4207e5f5cd2724aaac090cb9602e9e8aefe15eb8f24a6.
//
// Solidity: event LogSetBar(uint256 bar)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) FilterLogSetBar(opts *bind.FilterOpts) (*P1MirrorOracleETHUSDLogSetBarIterator, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.FilterLogs(opts, "LogSetBar")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDLogSetBarIterator{contract: _P1MirrorOracleETHUSD.contract, event: "LogSetBar", logs: logs, sub: sub}, nil
}

// WatchLogSetBar is a free log subscription operation binding the contract event 0x48c6ae1362d7627f13b4207e5f5cd2724aaac090cb9602e9e8aefe15eb8f24a6.
//
// Solidity: event LogSetBar(uint256 bar)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) WatchLogSetBar(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleETHUSDLogSetBar) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.WatchLogs(opts, "LogSetBar")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleETHUSDLogSetBar)
				if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetBar", log); err != nil {
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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) ParseLogSetBar(log types.Log) (*P1MirrorOracleETHUSDLogSetBar, error) {
	event := new(P1MirrorOracleETHUSDLogSetBar)
	if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetBar", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleETHUSDLogSetReaderIterator is returned from FilterLogSetReader and is used to iterate over the raw logs and unpacked data for LogSetReader events raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetReaderIterator struct {
	Event *P1MirrorOracleETHUSDLogSetReader // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleETHUSDLogSetReaderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleETHUSDLogSetReader)
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
		it.Event = new(P1MirrorOracleETHUSDLogSetReader)
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
func (it *P1MirrorOracleETHUSDLogSetReaderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleETHUSDLogSetReaderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleETHUSDLogSetReader represents a LogSetReader event raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetReader struct {
	Reader     common.Address
	Authorized bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogSetReader is a free log retrieval operation binding the contract event 0xadb3d91f6b7a78ea487b119a89fd644a0e6cf0909aa48faff97d153e0df682c0.
//
// Solidity: event LogSetReader(address reader, bool authorized)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) FilterLogSetReader(opts *bind.FilterOpts) (*P1MirrorOracleETHUSDLogSetReaderIterator, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.FilterLogs(opts, "LogSetReader")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDLogSetReaderIterator{contract: _P1MirrorOracleETHUSD.contract, event: "LogSetReader", logs: logs, sub: sub}, nil
}

// WatchLogSetReader is a free log subscription operation binding the contract event 0xadb3d91f6b7a78ea487b119a89fd644a0e6cf0909aa48faff97d153e0df682c0.
//
// Solidity: event LogSetReader(address reader, bool authorized)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) WatchLogSetReader(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleETHUSDLogSetReader) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.WatchLogs(opts, "LogSetReader")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleETHUSDLogSetReader)
				if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetReader", log); err != nil {
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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) ParseLogSetReader(log types.Log) (*P1MirrorOracleETHUSDLogSetReader, error) {
	event := new(P1MirrorOracleETHUSDLogSetReader)
	if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetReader", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleETHUSDLogSetSignerIterator is returned from FilterLogSetSigner and is used to iterate over the raw logs and unpacked data for LogSetSigner events raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetSignerIterator struct {
	Event *P1MirrorOracleETHUSDLogSetSigner // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleETHUSDLogSetSignerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleETHUSDLogSetSigner)
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
		it.Event = new(P1MirrorOracleETHUSDLogSetSigner)
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
func (it *P1MirrorOracleETHUSDLogSetSignerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleETHUSDLogSetSignerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleETHUSDLogSetSigner represents a LogSetSigner event raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDLogSetSigner struct {
	Signer     common.Address
	Authorized bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogSetSigner is a free log retrieval operation binding the contract event 0x8700965646f22bb776d5e0cbb11e1559f8143228405160e62250b866e954d912.
//
// Solidity: event LogSetSigner(address signer, bool authorized)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) FilterLogSetSigner(opts *bind.FilterOpts) (*P1MirrorOracleETHUSDLogSetSignerIterator, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.FilterLogs(opts, "LogSetSigner")
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDLogSetSignerIterator{contract: _P1MirrorOracleETHUSD.contract, event: "LogSetSigner", logs: logs, sub: sub}, nil
}

// WatchLogSetSigner is a free log subscription operation binding the contract event 0x8700965646f22bb776d5e0cbb11e1559f8143228405160e62250b866e954d912.
//
// Solidity: event LogSetSigner(address signer, bool authorized)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) WatchLogSetSigner(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleETHUSDLogSetSigner) (event.Subscription, error) {

	logs, sub, err := _P1MirrorOracleETHUSD.contract.WatchLogs(opts, "LogSetSigner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleETHUSDLogSetSigner)
				if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetSigner", log); err != nil {
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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) ParseLogSetSigner(log types.Log) (*P1MirrorOracleETHUSDLogSetSigner, error) {
	event := new(P1MirrorOracleETHUSDLogSetSigner)
	if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "LogSetSigner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MirrorOracleETHUSDOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDOwnershipTransferredIterator struct {
	Event *P1MirrorOracleETHUSDOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1MirrorOracleETHUSDOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MirrorOracleETHUSDOwnershipTransferred)
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
		it.Event = new(P1MirrorOracleETHUSDOwnershipTransferred)
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
func (it *P1MirrorOracleETHUSDOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MirrorOracleETHUSDOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MirrorOracleETHUSDOwnershipTransferred represents a OwnershipTransferred event raised by the P1MirrorOracleETHUSD contract.
type P1MirrorOracleETHUSDOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1MirrorOracleETHUSDOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MirrorOracleETHUSD.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1MirrorOracleETHUSDOwnershipTransferredIterator{contract: _P1MirrorOracleETHUSD.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1MirrorOracleETHUSDOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MirrorOracleETHUSD.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MirrorOracleETHUSDOwnershipTransferred)
				if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_P1MirrorOracleETHUSD *P1MirrorOracleETHUSDFilterer) ParseOwnershipTransferred(log types.Log) (*P1MirrorOracleETHUSDOwnershipTransferred, error) {
	event := new(P1MirrorOracleETHUSDOwnershipTransferred)
	if err := _P1MirrorOracleETHUSD.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
