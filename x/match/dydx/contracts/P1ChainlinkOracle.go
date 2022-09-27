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

// P1ChainlinkOracleMetaData contains all meta data concerning the P1ChainlinkOracle contract.
var P1ChainlinkOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"adjustmentExponent\",\"type\":\"uint96\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ADJUSTMENT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_MAPPING_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_ORACLE_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_READER_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1ChainlinkOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use P1ChainlinkOracleMetaData.ABI instead.
var P1ChainlinkOracleABI = P1ChainlinkOracleMetaData.ABI

// P1ChainlinkOracle is an auto generated Go binding around an Ethereum contract.
type P1ChainlinkOracle struct {
	P1ChainlinkOracleCaller     // Read-only binding to the contract
	P1ChainlinkOracleTransactor // Write-only binding to the contract
	P1ChainlinkOracleFilterer   // Log filterer for contract events
}

// P1ChainlinkOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1ChainlinkOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ChainlinkOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1ChainlinkOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ChainlinkOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1ChainlinkOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1ChainlinkOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1ChainlinkOracleSession struct {
	Contract     *P1ChainlinkOracle // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1ChainlinkOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1ChainlinkOracleCallerSession struct {
	Contract *P1ChainlinkOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// P1ChainlinkOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1ChainlinkOracleTransactorSession struct {
	Contract     *P1ChainlinkOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// P1ChainlinkOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1ChainlinkOracleRaw struct {
	Contract *P1ChainlinkOracle // Generic contract binding to access the raw methods on
}

// P1ChainlinkOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1ChainlinkOracleCallerRaw struct {
	Contract *P1ChainlinkOracleCaller // Generic read-only contract binding to access the raw methods on
}

// P1ChainlinkOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1ChainlinkOracleTransactorRaw struct {
	Contract *P1ChainlinkOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1ChainlinkOracle creates a new instance of P1ChainlinkOracle, bound to a specific deployed contract.
func NewP1ChainlinkOracle(address common.Address, backend bind.ContractBackend) (*P1ChainlinkOracle, error) {
	contract, err := bindP1ChainlinkOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1ChainlinkOracle{P1ChainlinkOracleCaller: P1ChainlinkOracleCaller{contract: contract}, P1ChainlinkOracleTransactor: P1ChainlinkOracleTransactor{contract: contract}, P1ChainlinkOracleFilterer: P1ChainlinkOracleFilterer{contract: contract}}, nil
}

// NewP1ChainlinkOracleCaller creates a new read-only instance of P1ChainlinkOracle, bound to a specific deployed contract.
func NewP1ChainlinkOracleCaller(address common.Address, caller bind.ContractCaller) (*P1ChainlinkOracleCaller, error) {
	contract, err := bindP1ChainlinkOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1ChainlinkOracleCaller{contract: contract}, nil
}

// NewP1ChainlinkOracleTransactor creates a new write-only instance of P1ChainlinkOracle, bound to a specific deployed contract.
func NewP1ChainlinkOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*P1ChainlinkOracleTransactor, error) {
	contract, err := bindP1ChainlinkOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1ChainlinkOracleTransactor{contract: contract}, nil
}

// NewP1ChainlinkOracleFilterer creates a new log filterer instance of P1ChainlinkOracle, bound to a specific deployed contract.
func NewP1ChainlinkOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*P1ChainlinkOracleFilterer, error) {
	contract, err := bindP1ChainlinkOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1ChainlinkOracleFilterer{contract: contract}, nil
}

// bindP1ChainlinkOracle binds a generic wrapper to an already deployed contract.
func bindP1ChainlinkOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1ChainlinkOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1ChainlinkOracle *P1ChainlinkOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1ChainlinkOracle.Contract.P1ChainlinkOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1ChainlinkOracle *P1ChainlinkOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1ChainlinkOracle.Contract.P1ChainlinkOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1ChainlinkOracle *P1ChainlinkOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1ChainlinkOracle.Contract.P1ChainlinkOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1ChainlinkOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1ChainlinkOracle *P1ChainlinkOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1ChainlinkOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1ChainlinkOracle *P1ChainlinkOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1ChainlinkOracle.Contract.contract.Transact(opts, method, params...)
}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleCaller) ADJUSTMENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1ChainlinkOracle.contract.Call(opts, &out, "_ADJUSTMENT_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleSession) ADJUSTMENT() (*big.Int, error) {
	return _P1ChainlinkOracle.Contract.ADJUSTMENT(&_P1ChainlinkOracle.CallOpts)
}

// ADJUSTMENT is a free data retrieval call binding the contract method 0x939a5439.
//
// Solidity: function _ADJUSTMENT_() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerSession) ADJUSTMENT() (*big.Int, error) {
	return _P1ChainlinkOracle.Contract.ADJUSTMENT(&_P1ChainlinkOracle.CallOpts)
}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1ChainlinkOracle *P1ChainlinkOracleCaller) MAPPING(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _P1ChainlinkOracle.contract.Call(opts, &out, "_MAPPING_", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1ChainlinkOracle *P1ChainlinkOracleSession) MAPPING(arg0 common.Address) ([32]byte, error) {
	return _P1ChainlinkOracle.Contract.MAPPING(&_P1ChainlinkOracle.CallOpts, arg0)
}

// MAPPING is a free data retrieval call binding the contract method 0xc1f75961.
//
// Solidity: function _MAPPING_(address ) view returns(bytes32)
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerSession) MAPPING(arg0 common.Address) ([32]byte, error) {
	return _P1ChainlinkOracle.Contract.MAPPING(&_P1ChainlinkOracle.CallOpts, arg0)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleCaller) ORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1ChainlinkOracle.contract.Call(opts, &out, "_ORACLE_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleSession) ORACLE() (common.Address, error) {
	return _P1ChainlinkOracle.Contract.ORACLE(&_P1ChainlinkOracle.CallOpts)
}

// ORACLE is a free data retrieval call binding the contract method 0x73a2ab7c.
//
// Solidity: function _ORACLE_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerSession) ORACLE() (common.Address, error) {
	return _P1ChainlinkOracle.Contract.ORACLE(&_P1ChainlinkOracle.CallOpts)
}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleCaller) READER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1ChainlinkOracle.contract.Call(opts, &out, "_READER_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleSession) READER() (common.Address, error) {
	return _P1ChainlinkOracle.Contract.READER(&_P1ChainlinkOracle.CallOpts)
}

// READER is a free data retrieval call binding the contract method 0x67b141ee.
//
// Solidity: function _READER_() view returns(address)
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerSession) READER() (common.Address, error) {
	return _P1ChainlinkOracle.Contract.READER(&_P1ChainlinkOracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1ChainlinkOracle.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleSession) GetPrice() (*big.Int, error) {
	return _P1ChainlinkOracle.Contract.GetPrice(&_P1ChainlinkOracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1ChainlinkOracle *P1ChainlinkOracleCallerSession) GetPrice() (*big.Int, error) {
	return _P1ChainlinkOracle.Contract.GetPrice(&_P1ChainlinkOracle.CallOpts)
}
