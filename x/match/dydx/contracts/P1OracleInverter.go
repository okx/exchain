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

// P1OracleInverterMetaData contains all meta data concerning the P1OracleInverter contract.
var P1OracleInverterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"adjustmentExponent\",\"type\":\"uint96\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ADJUSTMENT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_MAPPING_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ORACLE_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_READER_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1OracleInverterABI is the input ABI used to generate the binding from.
// Deprecated: Use P1OracleInverterMetaData.ABI instead.
var P1OracleInverterABI = P1OracleInverterMetaData.ABI

// P1OracleInverter is an auto generated Go binding around an Ethereum contract.
type P1OracleInverter struct {
	P1OracleInverterCaller     // Read-only binding to the contract
	P1OracleInverterTransactor // Write-only binding to the contract
	P1OracleInverterFilterer   // Log filterer for contract events
}

// P1OracleInverterCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1OracleInverterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OracleInverterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1OracleInverterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OracleInverterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1OracleInverterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OracleInverterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1OracleInverterSession struct {
	Contract     *P1OracleInverter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1OracleInverterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1OracleInverterCallerSession struct {
	Contract *P1OracleInverterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// P1OracleInverterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1OracleInverterTransactorSession struct {
	Contract     *P1OracleInverterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// P1OracleInverterRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1OracleInverterRaw struct {
	Contract *P1OracleInverter // Generic contract binding to access the raw methods on
}

// P1OracleInverterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1OracleInverterCallerRaw struct {
	Contract *P1OracleInverterCaller // Generic read-only contract binding to access the raw methods on
}

// P1OracleInverterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1OracleInverterTransactorRaw struct {
	Contract *P1OracleInverterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1OracleInverter creates a new instance of P1OracleInverter, bound to a specific deployed contract.
func NewP1OracleInverter(address common.Address, backend bind.ContractBackend) (*P1OracleInverter, error) {
	contract, err := bindP1OracleInverter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1OracleInverter{P1OracleInverterCaller: P1OracleInverterCaller{contract: contract}, P1OracleInverterTransactor: P1OracleInverterTransactor{contract: contract}, P1OracleInverterFilterer: P1OracleInverterFilterer{contract: contract}}, nil
}

// NewP1OracleInverterCaller creates a new read-only instance of P1OracleInverter, bound to a specific deployed contract.
func NewP1OracleInverterCaller(address common.Address, caller bind.ContractCaller) (*P1OracleInverterCaller, error) {
	contract, err := bindP1OracleInverter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1OracleInverterCaller{contract: contract}, nil
}

// NewP1OracleInverterTransactor creates a new write-only instance of P1OracleInverter, bound to a specific deployed contract.
func NewP1OracleInverterTransactor(address common.Address, transactor bind.ContractTransactor) (*P1OracleInverterTransactor, error) {
	contract, err := bindP1OracleInverter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1OracleInverterTransactor{contract: contract}, nil
}

// NewP1OracleInverterFilterer creates a new log filterer instance of P1OracleInverter, bound to a specific deployed contract.
func NewP1OracleInverterFilterer(address common.Address, filterer bind.ContractFilterer) (*P1OracleInverterFilterer, error) {
	contract, err := bindP1OracleInverter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1OracleInverterFilterer{contract: contract}, nil
}

// bindP1OracleInverter binds a generic wrapper to an already deployed contract.
func bindP1OracleInverter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1OracleInverterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1OracleInverter *P1OracleInverterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1OracleInverter.Contract.P1OracleInverterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1OracleInverter *P1OracleInverterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1OracleInverter.Contract.P1OracleInverterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1OracleInverter *P1OracleInverterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1OracleInverter.Contract.P1OracleInverterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1OracleInverter *P1OracleInverterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1OracleInverter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1OracleInverter *P1OracleInverterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1OracleInverter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1OracleInverter *P1OracleInverterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1OracleInverter.Contract.contract.Transact(opts, method, params...)
}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterCaller) ADJUSTMENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1OracleInverter.contract.Call(opts, &out, "_ADJUSTMENT_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterSession) ADJUSTMENT() (*big.Int, error) {
	return _P1OracleInverter.Contract.ADJUSTMENT(&_P1OracleInverter.CallOpts)
}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterCallerSession) ADJUSTMENT() (*big.Int, error) {
	return _P1OracleInverter.Contract.ADJUSTMENT(&_P1OracleInverter.CallOpts)
}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1OracleInverter *P1OracleInverterCaller) MAPPING(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _P1OracleInverter.contract.Call(opts, &out, "_MAPPING_", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1OracleInverter *P1OracleInverterSession) MAPPING(arg0 common.Address) ([32]byte, error) {
	return _P1OracleInverter.Contract.MAPPING(&_P1OracleInverter.CallOpts, arg0)
}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1OracleInverter *P1OracleInverterCallerSession) MAPPING(arg0 common.Address) ([32]byte, error) {
	return _P1OracleInverter.Contract.MAPPING(&_P1OracleInverter.CallOpts, arg0)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1OracleInverter *P1OracleInverterCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1OracleInverter.contract.Call(opts, &out, "_ORACLE_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1OracleInverter *P1OracleInverterSession) ORACLE() (common.Address, error) {
	return _P1OracleInverter.Contract.ORACLE(&_P1OracleInverter.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1OracleInverter *P1OracleInverterCallerSession) ORACLE() (common.Address, error) {
	return _P1OracleInverter.Contract.ORACLE(&_P1OracleInverter.CallOpts)
}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1OracleInverter *P1OracleInverterCaller) READER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1OracleInverter.contract.Call(opts, &out, "_READER_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1OracleInverter *P1OracleInverterSession) READER() (common.Address, error) {
	return _P1OracleInverter.Contract.READER(&_P1OracleInverter.CallOpts)
}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1OracleInverter *P1OracleInverterCallerSession) READER() (common.Address, error) {
	return _P1OracleInverter.Contract.READER(&_P1OracleInverter.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1OracleInverter.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterSession) GetPrice() (*big.Int, error) {
	return _P1OracleInverter.Contract.GetPrice(&_P1OracleInverter.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1OracleInverter *P1OracleInverterCallerSession) GetPrice() (*big.Int, error) {
	return _P1OracleInverter.Contract.GetPrice(&_P1OracleInverter.CallOpts)
}
