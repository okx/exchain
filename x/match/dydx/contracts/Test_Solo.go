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

// TestSoloMetaData contains all meta data concerning the TestSolo contract.
var TestSoloMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_TOKEN_ADDRESSES_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marketId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"setTokenAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marketId\",\"type\":\"uint256\"}],\"name\":\"getMarketTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"internalType\":\"structI_Solo.AccountInfo[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumI_Solo.ActionType\",\"name\":\"actionType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"accountId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"sign\",\"type\":\"bool\"},{\"internalType\":\"enumI_Solo.AssetDenomination\",\"name\":\"denomination\",\"type\":\"uint8\"},{\"internalType\":\"enumI_Solo.AssetReference\",\"name\":\"ref\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structI_Solo.AssetAmount\",\"name\":\"amount\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"primaryMarketId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"secondaryMarketId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"otherAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"otherAccountId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structI_Solo.ActionArgs[]\",\"name\":\"actions\",\"type\":\"tuple[]\"}],\"name\":\"operate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestSoloABI is the input ABI used to generate the binding from.
// Deprecated: Use TestSoloMetaData.ABI instead.
var TestSoloABI = TestSoloMetaData.ABI

// TestSolo is an auto generated Go binding around an Ethereum contract.
type TestSolo struct {
	TestSoloCaller     // Read-only binding to the contract
	TestSoloTransactor // Write-only binding to the contract
	TestSoloFilterer   // Log filterer for contract events
}

// TestSoloCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestSoloCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestSoloTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestSoloTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestSoloFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestSoloFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestSoloSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestSoloSession struct {
	Contract     *TestSolo         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestSoloCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestSoloCallerSession struct {
	Contract *TestSoloCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TestSoloTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestSoloTransactorSession struct {
	Contract     *TestSoloTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TestSoloRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestSoloRaw struct {
	Contract *TestSolo // Generic contract binding to access the raw methods on
}

// TestSoloCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestSoloCallerRaw struct {
	Contract *TestSoloCaller // Generic read-only contract binding to access the raw methods on
}

// TestSoloTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestSoloTransactorRaw struct {
	Contract *TestSoloTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestSolo creates a new instance of TestSolo, bound to a specific deployed contract.
func NewTestSolo(address common.Address, backend bind.ContractBackend) (*TestSolo, error) {
	contract, err := bindTestSolo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestSolo{TestSoloCaller: TestSoloCaller{contract: contract}, TestSoloTransactor: TestSoloTransactor{contract: contract}, TestSoloFilterer: TestSoloFilterer{contract: contract}}, nil
}

// NewTestSoloCaller creates a new read-only instance of TestSolo, bound to a specific deployed contract.
func NewTestSoloCaller(address common.Address, caller bind.ContractCaller) (*TestSoloCaller, error) {
	contract, err := bindTestSolo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestSoloCaller{contract: contract}, nil
}

// NewTestSoloTransactor creates a new write-only instance of TestSolo, bound to a specific deployed contract.
func NewTestSoloTransactor(address common.Address, transactor bind.ContractTransactor) (*TestSoloTransactor, error) {
	contract, err := bindTestSolo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestSoloTransactor{contract: contract}, nil
}

// NewTestSoloFilterer creates a new log filterer instance of TestSolo, bound to a specific deployed contract.
func NewTestSoloFilterer(address common.Address, filterer bind.ContractFilterer) (*TestSoloFilterer, error) {
	contract, err := bindTestSolo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestSoloFilterer{contract: contract}, nil
}

// bindTestSolo binds a generic wrapper to an already deployed contract.
func bindTestSolo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestSoloABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestSolo *TestSoloRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestSolo.Contract.TestSoloCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestSolo *TestSoloRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestSolo.Contract.TestSoloTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestSolo *TestSoloRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestSolo.Contract.TestSoloTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestSolo *TestSoloCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestSolo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestSolo *TestSoloTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestSolo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestSolo *TestSoloTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestSolo.Contract.contract.Transact(opts, method, params...)
}

// TOKENADDRESSES is a free data retrieval call binding the contract method 0x81f02d5b.
//
// Solidity: function _TOKEN_ADDRESSES_(uint256 ) view returns(address)
func (_TestSolo *TestSoloCaller) TOKENADDRESSES(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TestSolo.contract.Call(opts, &out, "_TOKEN_ADDRESSES_", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENADDRESSES is a free data retrieval call binding the contract method 0x81f02d5b.
//
// Solidity: function _TOKEN_ADDRESSES_(uint256 ) view returns(address)
func (_TestSolo *TestSoloSession) TOKENADDRESSES(arg0 *big.Int) (common.Address, error) {
	return _TestSolo.Contract.TOKENADDRESSES(&_TestSolo.CallOpts, arg0)
}

// TOKENADDRESSES is a free data retrieval call binding the contract method 0x81f02d5b.
//
// Solidity: function _TOKEN_ADDRESSES_(uint256 ) view returns(address)
func (_TestSolo *TestSoloCallerSession) TOKENADDRESSES(arg0 *big.Int) (common.Address, error) {
	return _TestSolo.Contract.TOKENADDRESSES(&_TestSolo.CallOpts, arg0)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_TestSolo *TestSoloCaller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _TestSolo.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_TestSolo *TestSoloSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _TestSolo.Contract.GetIsGlobalOperator(&_TestSolo.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_TestSolo *TestSoloCallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _TestSolo.Contract.GetIsGlobalOperator(&_TestSolo.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_TestSolo *TestSoloCaller) GetIsLocalOperator(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _TestSolo.contract.Call(opts, &out, "getIsLocalOperator", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_TestSolo *TestSoloSession) GetIsLocalOperator(owner common.Address, operator common.Address) (bool, error) {
	return _TestSolo.Contract.GetIsLocalOperator(&_TestSolo.CallOpts, owner, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_TestSolo *TestSoloCallerSession) GetIsLocalOperator(owner common.Address, operator common.Address) (bool, error) {
	return _TestSolo.Contract.GetIsLocalOperator(&_TestSolo.CallOpts, owner, operator)
}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_TestSolo *TestSoloCaller) GetMarketTokenAddress(opts *bind.CallOpts, marketId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TestSolo.contract.Call(opts, &out, "getMarketTokenAddress", marketId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_TestSolo *TestSoloSession) GetMarketTokenAddress(marketId *big.Int) (common.Address, error) {
	return _TestSolo.Contract.GetMarketTokenAddress(&_TestSolo.CallOpts, marketId)
}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_TestSolo *TestSoloCallerSession) GetMarketTokenAddress(marketId *big.Int) (common.Address, error) {
	return _TestSolo.Contract.GetMarketTokenAddress(&_TestSolo.CallOpts, marketId)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_TestSolo *TestSoloTransactor) Operate(opts *bind.TransactOpts, accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _TestSolo.contract.Transact(opts, "operate", accounts, actions)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_TestSolo *TestSoloSession) Operate(accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _TestSolo.Contract.Operate(&_TestSolo.TransactOpts, accounts, actions)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_TestSolo *TestSoloTransactorSession) Operate(accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _TestSolo.Contract.Operate(&_TestSolo.TransactOpts, accounts, actions)
}

// SetIsGlobalOperator is a paid mutator transaction binding the contract method 0x4d499acc.
//
// Solidity: function setIsGlobalOperator(address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloTransactor) SetIsGlobalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.contract.Transact(opts, "setIsGlobalOperator", operator, approved)
}

// SetIsGlobalOperator is a paid mutator transaction binding the contract method 0x4d499acc.
//
// Solidity: function setIsGlobalOperator(address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloSession) SetIsGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.Contract.SetIsGlobalOperator(&_TestSolo.TransactOpts, operator, approved)
}

// SetIsGlobalOperator is a paid mutator transaction binding the contract method 0x4d499acc.
//
// Solidity: function setIsGlobalOperator(address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloTransactorSession) SetIsGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.Contract.SetIsGlobalOperator(&_TestSolo.TransactOpts, operator, approved)
}

// SetIsLocalOperator is a paid mutator transaction binding the contract method 0x1e6b2f5c.
//
// Solidity: function setIsLocalOperator(address owner, address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloTransactor) SetIsLocalOperator(opts *bind.TransactOpts, owner common.Address, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.contract.Transact(opts, "setIsLocalOperator", owner, operator, approved)
}

// SetIsLocalOperator is a paid mutator transaction binding the contract method 0x1e6b2f5c.
//
// Solidity: function setIsLocalOperator(address owner, address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloSession) SetIsLocalOperator(owner common.Address, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.Contract.SetIsLocalOperator(&_TestSolo.TransactOpts, owner, operator, approved)
}

// SetIsLocalOperator is a paid mutator transaction binding the contract method 0x1e6b2f5c.
//
// Solidity: function setIsLocalOperator(address owner, address operator, bool approved) returns(bool)
func (_TestSolo *TestSoloTransactorSession) SetIsLocalOperator(owner common.Address, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestSolo.Contract.SetIsLocalOperator(&_TestSolo.TransactOpts, owner, operator, approved)
}

// SetTokenAddress is a paid mutator transaction binding the contract method 0x287e96c1.
//
// Solidity: function setTokenAddress(uint256 marketId, address tokenAddress) returns()
func (_TestSolo *TestSoloTransactor) SetTokenAddress(opts *bind.TransactOpts, marketId *big.Int, tokenAddress common.Address) (*types.Transaction, error) {
	return _TestSolo.contract.Transact(opts, "setTokenAddress", marketId, tokenAddress)
}

// SetTokenAddress is a paid mutator transaction binding the contract method 0x287e96c1.
//
// Solidity: function setTokenAddress(uint256 marketId, address tokenAddress) returns()
func (_TestSolo *TestSoloSession) SetTokenAddress(marketId *big.Int, tokenAddress common.Address) (*types.Transaction, error) {
	return _TestSolo.Contract.SetTokenAddress(&_TestSolo.TransactOpts, marketId, tokenAddress)
}

// SetTokenAddress is a paid mutator transaction binding the contract method 0x287e96c1.
//
// Solidity: function setTokenAddress(uint256 marketId, address tokenAddress) returns()
func (_TestSolo *TestSoloTransactorSession) SetTokenAddress(marketId *big.Int, tokenAddress common.Address) (*types.Transaction, error) {
	return _TestSolo.Contract.SetTokenAddress(&_TestSolo.TransactOpts, marketId, tokenAddress)
}
