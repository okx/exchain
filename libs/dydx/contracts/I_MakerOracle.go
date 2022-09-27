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

// IMakerOracleMetaData contains all meta data concerning the IMakerOracle contract.
var IMakerOracleMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"peek\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"read\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"age\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"orcl\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"bud\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"signerId\",\"type\":\"uint8\"}],\"name\":\"slot\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"val_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"age_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"poke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"kiss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"readers\",\"type\":\"address[]\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"reader\",\"type\":\"address\"}],\"name\":\"diss\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IMakerOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use IMakerOracleMetaData.ABI instead.
var IMakerOracleABI = IMakerOracleMetaData.ABI

// IMakerOracle is an auto generated Go binding around an Ethereum contract.
type IMakerOracle struct {
	IMakerOracleCaller     // Read-only binding to the contract
	IMakerOracleTransactor // Write-only binding to the contract
	IMakerOracleFilterer   // Log filterer for contract events
}

// IMakerOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type IMakerOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IMakerOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IMakerOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IMakerOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IMakerOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IMakerOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IMakerOracleSession struct {
	Contract     *IMakerOracle     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IMakerOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IMakerOracleCallerSession struct {
	Contract *IMakerOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IMakerOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IMakerOracleTransactorSession struct {
	Contract     *IMakerOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IMakerOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type IMakerOracleRaw struct {
	Contract *IMakerOracle // Generic contract binding to access the raw methods on
}

// IMakerOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IMakerOracleCallerRaw struct {
	Contract *IMakerOracleCaller // Generic read-only contract binding to access the raw methods on
}

// IMakerOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IMakerOracleTransactorRaw struct {
	Contract *IMakerOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIMakerOracle creates a new instance of IMakerOracle, bound to a specific deployed contract.
func NewIMakerOracle(address common.Address, backend bind.ContractBackend) (*IMakerOracle, error) {
	contract, err := bindIMakerOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IMakerOracle{IMakerOracleCaller: IMakerOracleCaller{contract: contract}, IMakerOracleTransactor: IMakerOracleTransactor{contract: contract}, IMakerOracleFilterer: IMakerOracleFilterer{contract: contract}}, nil
}

// NewIMakerOracleCaller creates a new read-only instance of IMakerOracle, bound to a specific deployed contract.
func NewIMakerOracleCaller(address common.Address, caller bind.ContractCaller) (*IMakerOracleCaller, error) {
	contract, err := bindIMakerOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IMakerOracleCaller{contract: contract}, nil
}

// NewIMakerOracleTransactor creates a new write-only instance of IMakerOracle, bound to a specific deployed contract.
func NewIMakerOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*IMakerOracleTransactor, error) {
	contract, err := bindIMakerOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IMakerOracleTransactor{contract: contract}, nil
}

// NewIMakerOracleFilterer creates a new log filterer instance of IMakerOracle, bound to a specific deployed contract.
func NewIMakerOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*IMakerOracleFilterer, error) {
	contract, err := bindIMakerOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IMakerOracleFilterer{contract: contract}, nil
}

// bindIMakerOracle binds a generic wrapper to an already deployed contract.
func bindIMakerOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IMakerOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IMakerOracle *IMakerOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IMakerOracle.Contract.IMakerOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IMakerOracle *IMakerOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IMakerOracle.Contract.IMakerOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IMakerOracle *IMakerOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IMakerOracle.Contract.IMakerOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IMakerOracle *IMakerOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IMakerOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IMakerOracle *IMakerOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IMakerOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IMakerOracle *IMakerOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IMakerOracle.Contract.contract.Transact(opts, method, params...)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_IMakerOracle *IMakerOracleCaller) Age(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "age")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_IMakerOracle *IMakerOracleSession) Age() (uint32, error) {
	return _IMakerOracle.Contract.Age(&_IMakerOracle.CallOpts)
}

// Age is a free data retrieval call binding the contract method 0x262a9dff.
//
// Solidity: function age() view returns(uint32)
func (_IMakerOracle *IMakerOracleCallerSession) Age() (uint32, error) {
	return _IMakerOracle.Contract.Age(&_IMakerOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_IMakerOracle *IMakerOracleCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_IMakerOracle *IMakerOracleSession) Bar() (*big.Int, error) {
	return _IMakerOracle.Contract.Bar(&_IMakerOracle.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_IMakerOracle *IMakerOracleCallerSession) Bar() (*big.Int, error) {
	return _IMakerOracle.Contract.Bar(&_IMakerOracle.CallOpts)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_IMakerOracle *IMakerOracleCaller) Bud(opts *bind.CallOpts, reader common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "bud", reader)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_IMakerOracle *IMakerOracleSession) Bud(reader common.Address) (*big.Int, error) {
	return _IMakerOracle.Contract.Bud(&_IMakerOracle.CallOpts, reader)
}

// Bud is a free data retrieval call binding the contract method 0x4fce7a2a.
//
// Solidity: function bud(address reader) view returns(uint256)
func (_IMakerOracle *IMakerOracleCallerSession) Bud(reader common.Address) (*big.Int, error) {
	return _IMakerOracle.Contract.Bud(&_IMakerOracle.CallOpts, reader)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_IMakerOracle *IMakerOracleCaller) Orcl(opts *bind.CallOpts, signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "orcl", signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_IMakerOracle *IMakerOracleSession) Orcl(signer common.Address) (*big.Int, error) {
	return _IMakerOracle.Contract.Orcl(&_IMakerOracle.CallOpts, signer)
}

// Orcl is a free data retrieval call binding the contract method 0x020b2e32.
//
// Solidity: function orcl(address signer) view returns(uint256)
func (_IMakerOracle *IMakerOracleCallerSession) Orcl(signer common.Address) (*big.Int, error) {
	return _IMakerOracle.Contract.Orcl(&_IMakerOracle.CallOpts, signer)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_IMakerOracle *IMakerOracleCaller) Peek(opts *bind.CallOpts) ([32]byte, bool, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "peek")

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
func (_IMakerOracle *IMakerOracleSession) Peek() ([32]byte, bool, error) {
	return _IMakerOracle.Contract.Peek(&_IMakerOracle.CallOpts)
}

// Peek is a free data retrieval call binding the contract method 0x59e02dd7.
//
// Solidity: function peek() view returns(bytes32, bool)
func (_IMakerOracle *IMakerOracleCallerSession) Peek() ([32]byte, bool, error) {
	return _IMakerOracle.Contract.Peek(&_IMakerOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_IMakerOracle *IMakerOracleCaller) Read(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "read")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_IMakerOracle *IMakerOracleSession) Read() ([32]byte, error) {
	return _IMakerOracle.Contract.Read(&_IMakerOracle.CallOpts)
}

// Read is a free data retrieval call binding the contract method 0x57de26a4.
//
// Solidity: function read() view returns(bytes32)
func (_IMakerOracle *IMakerOracleCallerSession) Read() ([32]byte, error) {
	return _IMakerOracle.Contract.Read(&_IMakerOracle.CallOpts)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_IMakerOracle *IMakerOracleCaller) Slot(opts *bind.CallOpts, signerId uint8) (common.Address, error) {
	var out []interface{}
	err := _IMakerOracle.contract.Call(opts, &out, "slot", signerId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_IMakerOracle *IMakerOracleSession) Slot(signerId uint8) (common.Address, error) {
	return _IMakerOracle.Contract.Slot(&_IMakerOracle.CallOpts, signerId)
}

// Slot is a free data retrieval call binding the contract method 0x8d0e5a9a.
//
// Solidity: function slot(uint8 signerId) view returns(address)
func (_IMakerOracle *IMakerOracleCallerSession) Slot(signerId uint8) (common.Address, error) {
	return _IMakerOracle.Contract.Slot(&_IMakerOracle.CallOpts, signerId)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleTransactor) Diss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.contract.Transact(opts, "diss", readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Diss(&_IMakerOracle.TransactOpts, readers)
}

// Diss is a paid mutator transaction binding the contract method 0x46d4577d.
//
// Solidity: function diss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleTransactorSession) Diss(readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Diss(&_IMakerOracle.TransactOpts, readers)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_IMakerOracle *IMakerOracleTransactor) Diss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.contract.Transact(opts, "diss0", reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_IMakerOracle *IMakerOracleSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Diss0(&_IMakerOracle.TransactOpts, reader)
}

// Diss0 is a paid mutator transaction binding the contract method 0x65c4ce7a.
//
// Solidity: function diss(address reader) returns()
func (_IMakerOracle *IMakerOracleTransactorSession) Diss0(reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Diss0(&_IMakerOracle.TransactOpts, reader)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleTransactor) Kiss(opts *bind.TransactOpts, readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.contract.Transact(opts, "kiss", readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Kiss(&_IMakerOracle.TransactOpts, readers)
}

// Kiss is a paid mutator transaction binding the contract method 0x1b25b65f.
//
// Solidity: function kiss(address[] readers) returns()
func (_IMakerOracle *IMakerOracleTransactorSession) Kiss(readers []common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Kiss(&_IMakerOracle.TransactOpts, readers)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_IMakerOracle *IMakerOracleTransactor) Kiss0(opts *bind.TransactOpts, reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.contract.Transact(opts, "kiss0", reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_IMakerOracle *IMakerOracleSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Kiss0(&_IMakerOracle.TransactOpts, reader)
}

// Kiss0 is a paid mutator transaction binding the contract method 0xf29c29c4.
//
// Solidity: function kiss(address reader) returns()
func (_IMakerOracle *IMakerOracleTransactorSession) Kiss0(reader common.Address) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Kiss0(&_IMakerOracle.TransactOpts, reader)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_IMakerOracle *IMakerOracleTransactor) Poke(opts *bind.TransactOpts, val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _IMakerOracle.contract.Transact(opts, "poke", val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_IMakerOracle *IMakerOracleSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Poke(&_IMakerOracle.TransactOpts, val_, age_, v, r, s)
}

// Poke is a paid mutator transaction binding the contract method 0x89bbb8b2.
//
// Solidity: function poke(uint256[] val_, uint256[] age_, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_IMakerOracle *IMakerOracleTransactorSession) Poke(val_ []*big.Int, age_ []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _IMakerOracle.Contract.Poke(&_IMakerOracle.TransactOpts, val_, age_, v, r, s)
}
