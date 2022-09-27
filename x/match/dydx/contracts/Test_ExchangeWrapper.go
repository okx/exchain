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

// TestExchangeWrapperMetaData contains all meta data concerning the TestExchangeWrapper contract.
var TestExchangeWrapperMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"EXCHANGE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_MAKER_AMOUNT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_TAKER_AMOUNT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"makerAmount\",\"type\":\"uint256\"}],\"name\":\"setMakerAmount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"takerAmount\",\"type\":\"uint256\"}],\"name\":\"setTakerAmount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getExchangeCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestedFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"orderData\",\"type\":\"bytes\"}],\"name\":\"exchange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestExchangeWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use TestExchangeWrapperMetaData.ABI instead.
var TestExchangeWrapperABI = TestExchangeWrapperMetaData.ABI

// TestExchangeWrapper is an auto generated Go binding around an Ethereum contract.
type TestExchangeWrapper struct {
	TestExchangeWrapperCaller     // Read-only binding to the contract
	TestExchangeWrapperTransactor // Write-only binding to the contract
	TestExchangeWrapperFilterer   // Log filterer for contract events
}

// TestExchangeWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestExchangeWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExchangeWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestExchangeWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExchangeWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestExchangeWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExchangeWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestExchangeWrapperSession struct {
	Contract     *TestExchangeWrapper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestExchangeWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestExchangeWrapperCallerSession struct {
	Contract *TestExchangeWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TestExchangeWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestExchangeWrapperTransactorSession struct {
	Contract     *TestExchangeWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TestExchangeWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestExchangeWrapperRaw struct {
	Contract *TestExchangeWrapper // Generic contract binding to access the raw methods on
}

// TestExchangeWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestExchangeWrapperCallerRaw struct {
	Contract *TestExchangeWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// TestExchangeWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestExchangeWrapperTransactorRaw struct {
	Contract *TestExchangeWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestExchangeWrapper creates a new instance of TestExchangeWrapper, bound to a specific deployed contract.
func NewTestExchangeWrapper(address common.Address, backend bind.ContractBackend) (*TestExchangeWrapper, error) {
	contract, err := bindTestExchangeWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestExchangeWrapper{TestExchangeWrapperCaller: TestExchangeWrapperCaller{contract: contract}, TestExchangeWrapperTransactor: TestExchangeWrapperTransactor{contract: contract}, TestExchangeWrapperFilterer: TestExchangeWrapperFilterer{contract: contract}}, nil
}

// NewTestExchangeWrapperCaller creates a new read-only instance of TestExchangeWrapper, bound to a specific deployed contract.
func NewTestExchangeWrapperCaller(address common.Address, caller bind.ContractCaller) (*TestExchangeWrapperCaller, error) {
	contract, err := bindTestExchangeWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestExchangeWrapperCaller{contract: contract}, nil
}

// NewTestExchangeWrapperTransactor creates a new write-only instance of TestExchangeWrapper, bound to a specific deployed contract.
func NewTestExchangeWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*TestExchangeWrapperTransactor, error) {
	contract, err := bindTestExchangeWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestExchangeWrapperTransactor{contract: contract}, nil
}

// NewTestExchangeWrapperFilterer creates a new log filterer instance of TestExchangeWrapper, bound to a specific deployed contract.
func NewTestExchangeWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*TestExchangeWrapperFilterer, error) {
	contract, err := bindTestExchangeWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestExchangeWrapperFilterer{contract: contract}, nil
}

// bindTestExchangeWrapper binds a generic wrapper to an already deployed contract.
func bindTestExchangeWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestExchangeWrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestExchangeWrapper *TestExchangeWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestExchangeWrapper.Contract.TestExchangeWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestExchangeWrapper *TestExchangeWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.TestExchangeWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestExchangeWrapper *TestExchangeWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.TestExchangeWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestExchangeWrapper *TestExchangeWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestExchangeWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestExchangeWrapper *TestExchangeWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestExchangeWrapper *TestExchangeWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.contract.Transact(opts, method, params...)
}

// EXCHANGEADDRESS is a free data retrieval call binding the contract method 0x60aec0f0.
//
// Solidity: function EXCHANGE_ADDRESS() view returns(address)
func (_TestExchangeWrapper *TestExchangeWrapperCaller) EXCHANGEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestExchangeWrapper.contract.Call(opts, &out, "EXCHANGE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EXCHANGEADDRESS is a free data retrieval call binding the contract method 0x60aec0f0.
//
// Solidity: function EXCHANGE_ADDRESS() view returns(address)
func (_TestExchangeWrapper *TestExchangeWrapperSession) EXCHANGEADDRESS() (common.Address, error) {
	return _TestExchangeWrapper.Contract.EXCHANGEADDRESS(&_TestExchangeWrapper.CallOpts)
}

// EXCHANGEADDRESS is a free data retrieval call binding the contract method 0x60aec0f0.
//
// Solidity: function EXCHANGE_ADDRESS() view returns(address)
func (_TestExchangeWrapper *TestExchangeWrapperCallerSession) EXCHANGEADDRESS() (common.Address, error) {
	return _TestExchangeWrapper.Contract.EXCHANGEADDRESS(&_TestExchangeWrapper.CallOpts)
}

// MAKERAMOUNT is a free data retrieval call binding the contract method 0x13dfc516.
//
// Solidity: function _MAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCaller) MAKERAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestExchangeWrapper.contract.Call(opts, &out, "_MAKER_AMOUNT_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAKERAMOUNT is a free data retrieval call binding the contract method 0x13dfc516.
//
// Solidity: function _MAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperSession) MAKERAMOUNT() (*big.Int, error) {
	return _TestExchangeWrapper.Contract.MAKERAMOUNT(&_TestExchangeWrapper.CallOpts)
}

// MAKERAMOUNT is a free data retrieval call binding the contract method 0x13dfc516.
//
// Solidity: function _MAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCallerSession) MAKERAMOUNT() (*big.Int, error) {
	return _TestExchangeWrapper.Contract.MAKERAMOUNT(&_TestExchangeWrapper.CallOpts)
}

// TAKERAMOUNT is a free data retrieval call binding the contract method 0x8c4d443f.
//
// Solidity: function _TAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCaller) TAKERAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestExchangeWrapper.contract.Call(opts, &out, "_TAKER_AMOUNT_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TAKERAMOUNT is a free data retrieval call binding the contract method 0x8c4d443f.
//
// Solidity: function _TAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperSession) TAKERAMOUNT() (*big.Int, error) {
	return _TestExchangeWrapper.Contract.TAKERAMOUNT(&_TestExchangeWrapper.CallOpts)
}

// TAKERAMOUNT is a free data retrieval call binding the contract method 0x8c4d443f.
//
// Solidity: function _TAKER_AMOUNT_() view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCallerSession) TAKERAMOUNT() (*big.Int, error) {
	return _TestExchangeWrapper.Contract.TAKERAMOUNT(&_TestExchangeWrapper.CallOpts)
}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address , address , uint256 , bytes ) view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCaller) GetExchangeCost(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*big.Int, error) {
	var out []interface{}
	err := _TestExchangeWrapper.contract.Call(opts, &out, "getExchangeCost", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address , address , uint256 , bytes ) view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperSession) GetExchangeCost(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*big.Int, error) {
	return _TestExchangeWrapper.Contract.GetExchangeCost(&_TestExchangeWrapper.CallOpts, arg0, arg1, arg2, arg3)
}

// GetExchangeCost is a free data retrieval call binding the contract method 0x3a8fdd7d.
//
// Solidity: function getExchangeCost(address , address , uint256 , bytes ) view returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperCallerSession) GetExchangeCost(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*big.Int, error) {
	return _TestExchangeWrapper.Contract.GetExchangeCost(&_TestExchangeWrapper.CallOpts, arg0, arg1, arg2, arg3)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address , address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperTransactor) Exchange(opts *bind.TransactOpts, arg0 common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _TestExchangeWrapper.contract.Transact(opts, "exchange", arg0, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address , address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperSession) Exchange(arg0 common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.Exchange(&_TestExchangeWrapper.TransactOpts, arg0, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}

// Exchange is a paid mutator transaction binding the contract method 0x7d98ebac.
//
// Solidity: function exchange(address , address receiver, address makerToken, address takerToken, uint256 requestedFillAmount, bytes orderData) returns(uint256)
func (_TestExchangeWrapper *TestExchangeWrapperTransactorSession) Exchange(arg0 common.Address, receiver common.Address, makerToken common.Address, takerToken common.Address, requestedFillAmount *big.Int, orderData []byte) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.Exchange(&_TestExchangeWrapper.TransactOpts, arg0, receiver, makerToken, takerToken, requestedFillAmount, orderData)
}

// SetMakerAmount is a paid mutator transaction binding the contract method 0x3ebddc21.
//
// Solidity: function setMakerAmount(uint256 makerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperTransactor) SetMakerAmount(opts *bind.TransactOpts, makerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.contract.Transact(opts, "setMakerAmount", makerAmount)
}

// SetMakerAmount is a paid mutator transaction binding the contract method 0x3ebddc21.
//
// Solidity: function setMakerAmount(uint256 makerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperSession) SetMakerAmount(makerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.SetMakerAmount(&_TestExchangeWrapper.TransactOpts, makerAmount)
}

// SetMakerAmount is a paid mutator transaction binding the contract method 0x3ebddc21.
//
// Solidity: function setMakerAmount(uint256 makerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperTransactorSession) SetMakerAmount(makerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.SetMakerAmount(&_TestExchangeWrapper.TransactOpts, makerAmount)
}

// SetTakerAmount is a paid mutator transaction binding the contract method 0xcc246d8b.
//
// Solidity: function setTakerAmount(uint256 takerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperTransactor) SetTakerAmount(opts *bind.TransactOpts, takerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.contract.Transact(opts, "setTakerAmount", takerAmount)
}

// SetTakerAmount is a paid mutator transaction binding the contract method 0xcc246d8b.
//
// Solidity: function setTakerAmount(uint256 takerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperSession) SetTakerAmount(takerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.SetTakerAmount(&_TestExchangeWrapper.TransactOpts, takerAmount)
}

// SetTakerAmount is a paid mutator transaction binding the contract method 0xcc246d8b.
//
// Solidity: function setTakerAmount(uint256 takerAmount) returns()
func (_TestExchangeWrapper *TestExchangeWrapperTransactorSession) SetTakerAmount(takerAmount *big.Int) (*types.Transaction, error) {
	return _TestExchangeWrapper.Contract.SetTakerAmount(&_TestExchangeWrapper.TransactOpts, takerAmount)
}
