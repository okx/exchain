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

// TestP1TraderMetaData contains all meta data concerning the TestP1Trader contract.
var TestP1TraderMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"TRADER_FLAG_RESULT_2\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_TRADE_RESULT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_TRADE_RESULT_2_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"setTradeResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"setSecondTradeResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestP1TraderABI is the input ABI used to generate the binding from.
// Deprecated: Use TestP1TraderMetaData.ABI instead.
var TestP1TraderABI = TestP1TraderMetaData.ABI

// TestP1Trader is an auto generated Go binding around an Ethereum contract.
type TestP1Trader struct {
	TestP1TraderCaller     // Read-only binding to the contract
	TestP1TraderTransactor // Write-only binding to the contract
	TestP1TraderFilterer   // Log filterer for contract events
}

// TestP1TraderCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestP1TraderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1TraderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestP1TraderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1TraderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestP1TraderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1TraderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestP1TraderSession struct {
	Contract     *TestP1Trader     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestP1TraderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestP1TraderCallerSession struct {
	Contract *TestP1TraderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TestP1TraderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestP1TraderTransactorSession struct {
	Contract     *TestP1TraderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TestP1TraderRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestP1TraderRaw struct {
	Contract *TestP1Trader // Generic contract binding to access the raw methods on
}

// TestP1TraderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestP1TraderCallerRaw struct {
	Contract *TestP1TraderCaller // Generic read-only contract binding to access the raw methods on
}

// TestP1TraderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestP1TraderTransactorRaw struct {
	Contract *TestP1TraderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestP1Trader creates a new instance of TestP1Trader, bound to a specific deployed contract.
func NewTestP1Trader(address common.Address, backend bind.ContractBackend) (*TestP1Trader, error) {
	contract, err := bindTestP1Trader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestP1Trader{TestP1TraderCaller: TestP1TraderCaller{contract: contract}, TestP1TraderTransactor: TestP1TraderTransactor{contract: contract}, TestP1TraderFilterer: TestP1TraderFilterer{contract: contract}}, nil
}

// NewTestP1TraderCaller creates a new read-only instance of TestP1Trader, bound to a specific deployed contract.
func NewTestP1TraderCaller(address common.Address, caller bind.ContractCaller) (*TestP1TraderCaller, error) {
	contract, err := bindTestP1Trader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1TraderCaller{contract: contract}, nil
}

// NewTestP1TraderTransactor creates a new write-only instance of TestP1Trader, bound to a specific deployed contract.
func NewTestP1TraderTransactor(address common.Address, transactor bind.ContractTransactor) (*TestP1TraderTransactor, error) {
	contract, err := bindTestP1Trader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1TraderTransactor{contract: contract}, nil
}

// NewTestP1TraderFilterer creates a new log filterer instance of TestP1Trader, bound to a specific deployed contract.
func NewTestP1TraderFilterer(address common.Address, filterer bind.ContractFilterer) (*TestP1TraderFilterer, error) {
	contract, err := bindTestP1Trader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestP1TraderFilterer{contract: contract}, nil
}

// bindTestP1Trader binds a generic wrapper to an already deployed contract.
func bindTestP1Trader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestP1TraderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Trader *TestP1TraderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Trader.Contract.TestP1TraderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Trader *TestP1TraderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Trader.Contract.TestP1TraderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Trader *TestP1TraderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Trader.Contract.TestP1TraderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Trader *TestP1TraderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Trader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Trader *TestP1TraderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Trader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Trader *TestP1TraderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Trader.Contract.contract.Transact(opts, method, params...)
}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Trader *TestP1TraderCaller) TRADERFLAGRESULT2(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestP1Trader.contract.Call(opts, &out, "TRADER_FLAG_RESULT_2")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Trader *TestP1TraderSession) TRADERFLAGRESULT2() ([32]byte, error) {
	return _TestP1Trader.Contract.TRADERFLAGRESULT2(&_TestP1Trader.CallOpts)
}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Trader *TestP1TraderCallerSession) TRADERFLAGRESULT2() ([32]byte, error) {
	return _TestP1Trader.Contract.TRADERFLAGRESULT2(&_TestP1Trader.CallOpts)
}

// TRADERESULT is a free data retrieval call binding the contract method 0x495f9bff.
//
// Solidity: function _TRADE_RESULT_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderCaller) TRADERESULT(opts *bind.CallOpts) (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	var out []interface{}
	err := _TestP1Trader.contract.Call(opts, &out, "_TRADE_RESULT_")

	outstruct := new(struct {
		MarginAmount   *big.Int
		PositionAmount *big.Int
		IsBuy          bool
		TraderFlags    [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MarginAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.PositionAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsBuy = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.TraderFlags = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// TRADERESULT is a free data retrieval call binding the contract method 0x495f9bff.
//
// Solidity: function _TRADE_RESULT_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderSession) TRADERESULT() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Trader.Contract.TRADERESULT(&_TestP1Trader.CallOpts)
}

// TRADERESULT is a free data retrieval call binding the contract method 0x495f9bff.
//
// Solidity: function _TRADE_RESULT_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderCallerSession) TRADERESULT() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Trader.Contract.TRADERESULT(&_TestP1Trader.CallOpts)
}

// TRADERESULT2 is a free data retrieval call binding the contract method 0x75092a30.
//
// Solidity: function _TRADE_RESULT_2_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderCaller) TRADERESULT2(opts *bind.CallOpts) (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	var out []interface{}
	err := _TestP1Trader.contract.Call(opts, &out, "_TRADE_RESULT_2_")

	outstruct := new(struct {
		MarginAmount   *big.Int
		PositionAmount *big.Int
		IsBuy          bool
		TraderFlags    [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MarginAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.PositionAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsBuy = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.TraderFlags = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// TRADERESULT2 is a free data retrieval call binding the contract method 0x75092a30.
//
// Solidity: function _TRADE_RESULT_2_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderSession) TRADERESULT2() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Trader.Contract.TRADERESULT2(&_TestP1Trader.CallOpts)
}

// TRADERESULT2 is a free data retrieval call binding the contract method 0x75092a30.
//
// Solidity: function _TRADE_RESULT_2_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Trader *TestP1TraderCallerSession) TRADERESULT2() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Trader.Contract.TRADERESULT2(&_TestP1Trader.CallOpts)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderTransactor) SetSecondTradeResult(opts *bind.TransactOpts, marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.contract.Transact(opts, "setSecondTradeResult", marginAmount, positionAmount, isBuy, traderFlags)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderSession) SetSecondTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.SetSecondTradeResult(&_TestP1Trader.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderTransactorSession) SetSecondTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.SetSecondTradeResult(&_TestP1Trader.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderTransactor) SetTradeResult(opts *bind.TransactOpts, marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.contract.Transact(opts, "setTradeResult", marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderSession) SetTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.SetTradeResult(&_TestP1Trader.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Trader *TestP1TraderTransactorSession) SetTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.SetTradeResult(&_TestP1Trader.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Trader *TestP1TraderTransactor) Trade(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.contract.Transact(opts, "trade", arg0, arg1, arg2, arg3, arg4, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Trader *TestP1TraderSession) Trade(arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.Trade(&_TestP1Trader.TransactOpts, arg0, arg1, arg2, arg3, arg4, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Trader *TestP1TraderTransactorSession) Trade(arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Trader.Contract.Trade(&_TestP1Trader.TransactOpts, arg0, arg1, arg2, arg3, arg4, traderFlags)
}
