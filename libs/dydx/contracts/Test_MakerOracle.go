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

// TestMakerOracleMetaData contains all meta data concerning the TestMakerOracle contract.
var TestMakerOracleMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"_PRICE_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_VALID_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"age\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bud\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"orcl\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"slot\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBar\",\"type\":\"uint256\"}],\"name\":\"setBar\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newAge\",\"type\":\"uint256\"}],\"name\":\"setAge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPrice\",\"type\":\"uint256\"}],\"name\":\"setPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"}],\"name\":\"setValidity\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"read\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"peek\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"name\":\"poke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"lift\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"}],\"name\":\"drop\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestMakerOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use TestMakerOracleMetaData.ABI instead.
var TestMakerOracleABI = TestMakerOracleMetaData.ABI

// TestMakerOracle is an auto generated Go binding around an Ethereum contract.
type TestMakerOracle struct {
	TestMakerOracleCaller     // Read-only binding to the contract
	TestMakerOracleTransactor // Write-only binding to the contract
	TestMakerOracleFilterer   // Log filterer for contract events
}

// TestMakerOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestMakerOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMakerOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestMakerOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMakerOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestMakerOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMakerOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestMakerOracleSession struct {
	Contract     *TestMakerOracle  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestMakerOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestMakerOracleCallerSession struct {
	Contract *TestMakerOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// TestMakerOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestMakerOracleTransactorSession struct {
	Contract     *TestMakerOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// TestMakerOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestMakerOracleRaw struct {
	Contract *TestMakerOracle // Generic contract binding to access the raw methods on
}

// TestMakerOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestMakerOracleCallerRaw struct {
	Contract *TestMakerOracleCaller // Generic read-only contract binding to access the raw methods on
}

// TestMakerOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestMakerOracleTransactorRaw struct {
	Contract *TestMakerOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestMakerOracle creates a new instance of TestMakerOracle, bound to a specific deployed contract.
func NewTestMakerOracle(address common.Address, backend bind.ContractBackend) (*TestMakerOracle, error) {
	contract, err := bindTestMakerOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestMakerOracle{TestMakerOracleCaller: TestMakerOracleCaller{contract: contract}, TestMakerOracleTransactor: TestMakerOracleTransactor{contract: contract}, TestMakerOracleFilterer: TestMakerOracleFilterer{contract: contract}}, nil
}

// NewTestMakerOracleCaller creates a new read-only instance of TestMakerOracle, bound to a specific deployed contract.
func NewTestMakerOracleCaller(address common.Address, caller bind.ContractCaller) (*TestMakerOracleCaller, error) {
	contract, err := bindTestMakerOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestMakerOracleCaller{contract: contract}, nil
}

// NewTestMakerOracleTransactor creates a new write-only instance of TestMakerOracle, bound to a specific deployed contract.
func NewTestMakerOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*TestMakerOracleTransactor, error) {
	contract, err := bindTestMakerOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestMakerOracleTransactor{contract: contract}, nil
}

// NewTestMakerOracleFilterer creates a new log filterer instance of TestMakerOracle, bound to a specific deployed contract.
func NewTestMakerOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*TestMakerOracleFilterer, error) {
	contract, err := bindTestMakerOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestMakerOracleFilterer{contract: contract}, nil
}

// bindTestMakerOracle binds a generic wrapper to an already deployed contract.
func bindTestMakerOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestMakerOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestMakerOracle *TestMakerOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestMakerOracle.Contract.TestMakerOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestMakerOracle *TestMakerOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.TestMakerOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestMakerOracle *TestMakerOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.TestMakerOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestMakerOracle *TestMakerOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestMakerOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestMakerOracle *TestMakerOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestMakerOracle *TestMakerOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.contract.Transact(opts, method, params...)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCaller) PRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "_PRICE_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleSession) PRICE() (*big.Int, error) {
	return _TestMakerOracle.Contract.PRICE(&_TestMakerOracle.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCallerSession) PRICE() (*big.Int, error) {
	return _TestMakerOracle.Contract.PRICE(&_TestMakerOracle.CallOpts)
}

// VALID is a free data retrieval call binding the contract method 0x6f0fb301.
//
// Solidity: function _VALID_() view returns(bool)
func (_TestMakerOracle *TestMakerOracleCaller) VALID(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "_VALID_")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VALID is a free data retrieval call binding the contract method 0x6f0fb301.
//
// Solidity: function _VALID_() view returns(bool)
func (_TestMakerOracle *TestMakerOracleSession) VALID() (bool, error) {
	return _TestMakerOracle.Contract.VALID(&_TestMakerOracle.CallOpts)
}

// VALID is a free data retrieval call binding the contract method 0x6f0fb301.
//
// Solidity: function _VALID_() view returns(bool)
func (_TestMakerOracle *TestMakerOracleCallerSession) VALID() (bool, error) {
	return _TestMakerOracle.Contract.VALID(&_TestMakerOracle.CallOpts)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_TestMakerOracle *TestMakerOracleCaller) Age(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "age")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_TestMakerOracle *TestMakerOracleSession) Age() (uint32, error) {
	return _TestMakerOracle.Contract.Age(&_TestMakerOracle.CallOpts)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_TestMakerOracle *TestMakerOracleCallerSession) Age() (uint32, error) {
	return _TestMakerOracle.Contract.Age(&_TestMakerOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleSession) Bar() (*big.Int, error) {
	return _TestMakerOracle.Contract.Bar(&_TestMakerOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCallerSession) Bar() (*big.Int, error) {
	return _TestMakerOracle.Contract.Bar(&_TestMakerOracle.CallOpts)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCaller) Bud(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "bud", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleSession) Bud(arg0 common.Address) (*big.Int, error) {
	return _TestMakerOracle.Contract.Bud(&_TestMakerOracle.CallOpts, arg0)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCallerSession) Bud(arg0 common.Address) (*big.Int, error) {
	return _TestMakerOracle.Contract.Bud(&_TestMakerOracle.CallOpts, arg0)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCaller) Orcl(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "orcl", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleSession) Orcl(arg0 common.Address) (*big.Int, error) {
	return _TestMakerOracle.Contract.Orcl(&_TestMakerOracle.CallOpts, arg0)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address ) view returns(uint256)
func (_TestMakerOracle *TestMakerOracleCallerSession) Orcl(arg0 common.Address) (*big.Int, error) {
	return _TestMakerOracle.Contract.Orcl(&_TestMakerOracle.CallOpts, arg0)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_TestMakerOracle *TestMakerOracleCaller) Peek(opts *bind.CallOpts) ([32]byte, bool, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "peek")

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
func (_TestMakerOracle *TestMakerOracleSession) Peek() ([32]byte, bool, error) {
	return _TestMakerOracle.Contract.Peek(&_TestMakerOracle.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_TestMakerOracle *TestMakerOracleCallerSession) Peek() ([32]byte, bool, error) {
	return _TestMakerOracle.Contract.Peek(&_TestMakerOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_TestMakerOracle *TestMakerOracleCaller) Read(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "read")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_TestMakerOracle *TestMakerOracleSession) Read() ([32]byte, error) {
	return _TestMakerOracle.Contract.Read(&_TestMakerOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_TestMakerOracle *TestMakerOracleCallerSession) Read() ([32]byte, error) {
	return _TestMakerOracle.Contract.Read(&_TestMakerOracle.CallOpts)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 ) view returns(address)
func (_TestMakerOracle *TestMakerOracleCaller) Slot(opts *bind.CallOpts, arg0 uint8) (common.Address, error) {
	var out []interface{}
	err := _TestMakerOracle.contract.Call(opts, &out, "slot", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 ) view returns(address)
func (_TestMakerOracle *TestMakerOracleSession) Slot(arg0 uint8) (common.Address, error) {
	return _TestMakerOracle.Contract.Slot(&_TestMakerOracle.CallOpts, arg0)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 ) view returns(address)
func (_TestMakerOracle *TestMakerOracleCallerSession) Slot(arg0 uint8) (common.Address, error) {
	return _TestMakerOracle.Contract.Slot(&_TestMakerOracle.CallOpts, arg0)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Diss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "diss", readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Diss(&_TestMakerOracle.TransactOpts, readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Diss(&_TestMakerOracle.TransactOpts, readers)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Diss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "diss0", reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Diss0(&_TestMakerOracle.TransactOpts, reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Diss0(&_TestMakerOracle.TransactOpts, reader)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Drop(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "drop", signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Drop(&_TestMakerOracle.TransactOpts, signers)
}

// Drop is a paid mutator transaction binding the contract method 0x8ef5eaf0.
//
// Solidity: function drop(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Drop(signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Drop(&_TestMakerOracle.TransactOpts, signers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Kiss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "kiss", readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Kiss(&_TestMakerOracle.TransactOpts, readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Kiss(&_TestMakerOracle.TransactOpts, readers)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Kiss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "kiss0", reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Kiss0(&_TestMakerOracle.TransactOpts, reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Kiss0(&_TestMakerOracle.TransactOpts, reader)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Lift(opts *bind.TransactOpts, signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "lift", signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Lift(&_TestMakerOracle.TransactOpts, signers)
}

// Lift is a paid mutator transaction binding the contract method 0x94318106.
//
// Solidity: function lift(address[] signers) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Lift(signers []common.Address) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Lift(&_TestMakerOracle.TransactOpts, signers)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] , uint256[] , uint8[] , bytes32[] , bytes32[] ) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) Poke(opts *bind.TransactOpts, arg0 []*big.Int, arg1 []*big.Int, arg2 []uint8, arg3 [][32]byte, arg4 [][32]byte) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "poke", arg0, arg1, arg2, arg3, arg4)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] , uint256[] , uint8[] , bytes32[] , bytes32[] ) returns()
func (_TestMakerOracle *TestMakerOracleSession) Poke(arg0 []*big.Int, arg1 []*big.Int, arg2 []uint8, arg3 [][32]byte, arg4 [][32]byte) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Poke(&_TestMakerOracle.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] , uint256[] , uint8[] , bytes32[] , bytes32[] ) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) Poke(arg0 []*big.Int, arg1 []*big.Int, arg2 []uint8, arg3 [][32]byte, arg4 [][32]byte) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.Poke(&_TestMakerOracle.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// SetAge is a paid mutator transaction binding the contract method 0xd5dcf127.
//
// Solidity: function setAge(uint256 newAge) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) SetAge(opts *bind.TransactOpts, newAge *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "setAge", newAge)
}

// SetAge is a paid mutator transaction binding the contract method 0xd5dcf127.
//
// Solidity: function setAge(uint256 newAge) returns()
func (_TestMakerOracle *TestMakerOracleSession) SetAge(newAge *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetAge(&_TestMakerOracle.TransactOpts, newAge)
}

// SetAge is a paid mutator transaction binding the contract method 0xd5dcf127.
//
// Solidity: function setAge(uint256 newAge) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) SetAge(newAge *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetAge(&_TestMakerOracle.TransactOpts, newAge)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 newBar) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) SetBar(opts *bind.TransactOpts, newBar *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "setBar", newBar)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 newBar) returns()
func (_TestMakerOracle *TestMakerOracleSession) SetBar(newBar *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetBar(&_TestMakerOracle.TransactOpts, newBar)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 newBar) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) SetBar(newBar *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetBar(&_TestMakerOracle.TransactOpts, newBar)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) SetPrice(opts *bind.TransactOpts, newPrice *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "setPrice", newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestMakerOracle *TestMakerOracleSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetPrice(&_TestMakerOracle.TransactOpts, newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetPrice(&_TestMakerOracle.TransactOpts, newPrice)
}

// SetValidity is a paid mutator transaction binding the contract method 0xb8a35a01.
//
// Solidity: function setValidity(bool valid) returns()
func (_TestMakerOracle *TestMakerOracleTransactor) SetValidity(opts *bind.TransactOpts, valid bool) (*types.Transaction, error) {
	return _TestMakerOracle.contract.Transact(opts, "setValidity", valid)
}

// SetValidity is a paid mutator transaction binding the contract method 0xb8a35a01.
//
// Solidity: function setValidity(bool valid) returns()
func (_TestMakerOracle *TestMakerOracleSession) SetValidity(valid bool) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetValidity(&_TestMakerOracle.TransactOpts, valid)
}

// SetValidity is a paid mutator transaction binding the contract method 0xb8a35a01.
//
// Solidity: function setValidity(bool valid) returns()
func (_TestMakerOracle *TestMakerOracleTransactorSession) SetValidity(valid bool) (*types.Transaction, error) {
	return _TestMakerOracle.Contract.SetValidity(&_TestMakerOracle.TransactOpts, valid)
}
