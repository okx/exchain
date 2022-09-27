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

// TestP1MonolithMetaData contains all meta data concerning the TestP1Monolith contract.
var TestP1MonolithMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"TRADER_FLAG_RESULT_2\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_FUNDING_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_FUNDING_IS_POSITIVE_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_PRICE_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_TRADE_RESULT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_TRADE_RESULT_2_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getFunding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"newFunding\",\"type\":\"uint256\"}],\"name\":\"setFunding\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newPrice\",\"type\":\"uint256\"}],\"name\":\"setPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"setSecondTradeResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"setTradeResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestP1MonolithABI is the input ABI used to generate the binding from.
// Deprecated: Use TestP1MonolithMetaData.ABI instead.
var TestP1MonolithABI = TestP1MonolithMetaData.ABI

// TestP1Monolith is an auto generated Go binding around an Ethereum contract.
type TestP1Monolith struct {
	TestP1MonolithCaller     // Read-only binding to the contract
	TestP1MonolithTransactor // Write-only binding to the contract
	TestP1MonolithFilterer   // Log filterer for contract events
}

// TestP1MonolithCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestP1MonolithCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1MonolithTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestP1MonolithTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1MonolithFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestP1MonolithFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestP1MonolithSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestP1MonolithSession struct {
	Contract     *TestP1Monolith   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestP1MonolithCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestP1MonolithCallerSession struct {
	Contract *TestP1MonolithCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TestP1MonolithTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestP1MonolithTransactorSession struct {
	Contract     *TestP1MonolithTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TestP1MonolithRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestP1MonolithRaw struct {
	Contract *TestP1Monolith // Generic contract binding to access the raw methods on
}

// TestP1MonolithCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestP1MonolithCallerRaw struct {
	Contract *TestP1MonolithCaller // Generic read-only contract binding to access the raw methods on
}

// TestP1MonolithTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestP1MonolithTransactorRaw struct {
	Contract *TestP1MonolithTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestP1Monolith creates a new instance of TestP1Monolith, bound to a specific deployed contract.
func NewTestP1Monolith(address common.Address, backend bind.ContractBackend) (*TestP1Monolith, error) {
	contract, err := bindTestP1Monolith(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestP1Monolith{TestP1MonolithCaller: TestP1MonolithCaller{contract: contract}, TestP1MonolithTransactor: TestP1MonolithTransactor{contract: contract}, TestP1MonolithFilterer: TestP1MonolithFilterer{contract: contract}}, nil
}

// NewTestP1MonolithCaller creates a new read-only instance of TestP1Monolith, bound to a specific deployed contract.
func NewTestP1MonolithCaller(address common.Address, caller bind.ContractCaller) (*TestP1MonolithCaller, error) {
	contract, err := bindTestP1Monolith(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1MonolithCaller{contract: contract}, nil
}

// NewTestP1MonolithTransactor creates a new write-only instance of TestP1Monolith, bound to a specific deployed contract.
func NewTestP1MonolithTransactor(address common.Address, transactor bind.ContractTransactor) (*TestP1MonolithTransactor, error) {
	contract, err := bindTestP1Monolith(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestP1MonolithTransactor{contract: contract}, nil
}

// NewTestP1MonolithFilterer creates a new log filterer instance of TestP1Monolith, bound to a specific deployed contract.
func NewTestP1MonolithFilterer(address common.Address, filterer bind.ContractFilterer) (*TestP1MonolithFilterer, error) {
	contract, err := bindTestP1Monolith(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestP1MonolithFilterer{contract: contract}, nil
}

// bindTestP1Monolith binds a generic wrapper to an already deployed contract.
func bindTestP1Monolith(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestP1MonolithABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Monolith *TestP1MonolithRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Monolith.Contract.TestP1MonolithCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Monolith *TestP1MonolithRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.TestP1MonolithTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Monolith *TestP1MonolithRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.TestP1MonolithTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestP1Monolith *TestP1MonolithCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestP1Monolith.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestP1Monolith *TestP1MonolithTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestP1Monolith *TestP1MonolithTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.contract.Transact(opts, method, params...)
}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Monolith *TestP1MonolithCaller) TRADERFLAGRESULT2(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "TRADER_FLAG_RESULT_2")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Monolith *TestP1MonolithSession) TRADERFLAGRESULT2() ([32]byte, error) {
	return _TestP1Monolith.Contract.TRADERFLAGRESULT2(&_TestP1Monolith.CallOpts)
}

// TRADERFLAGRESULT2 is a free data retrieval call binding the contract method 0x6d2c6021.
//
// Solidity: function TRADER_FLAG_RESULT_2() view returns(bytes32)
func (_TestP1Monolith *TestP1MonolithCallerSession) TRADERFLAGRESULT2() ([32]byte, error) {
	return _TestP1Monolith.Contract.TRADERFLAGRESULT2(&_TestP1Monolith.CallOpts)
}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCaller) FUNDING(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "_FUNDING_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithSession) FUNDING() (*big.Int, error) {
	return _TestP1Monolith.Contract.FUNDING(&_TestP1Monolith.CallOpts)
}

// FUNDING is a free data retrieval call binding the contract method 0x4993cc3b.
//
// Solidity: function _FUNDING_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCallerSession) FUNDING() (*big.Int, error) {
	return _TestP1Monolith.Contract.FUNDING(&_TestP1Monolith.CallOpts)
}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Monolith *TestP1MonolithCaller) FUNDINGISPOSITIVE(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "_FUNDING_IS_POSITIVE_")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Monolith *TestP1MonolithSession) FUNDINGISPOSITIVE() (bool, error) {
	return _TestP1Monolith.Contract.FUNDINGISPOSITIVE(&_TestP1Monolith.CallOpts)
}

// FUNDINGISPOSITIVE is a free data retrieval call binding the contract method 0x910fb073.
//
// Solidity: function _FUNDING_IS_POSITIVE_() view returns(bool)
func (_TestP1Monolith *TestP1MonolithCallerSession) FUNDINGISPOSITIVE() (bool, error) {
	return _TestP1Monolith.Contract.FUNDINGISPOSITIVE(&_TestP1Monolith.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCaller) PRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "_PRICE_")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithSession) PRICE() (*big.Int, error) {
	return _TestP1Monolith.Contract.PRICE(&_TestP1Monolith.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x9f5cf46a.
//
// Solidity: function _PRICE_() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCallerSession) PRICE() (*big.Int, error) {
	return _TestP1Monolith.Contract.PRICE(&_TestP1Monolith.CallOpts)
}

// TRADERESULT is a free data retrieval call binding the contract method 0x495f9bff.
//
// Solidity: function _TRADE_RESULT_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Monolith *TestP1MonolithCaller) TRADERESULT(opts *bind.CallOpts) (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "_TRADE_RESULT_")

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
func (_TestP1Monolith *TestP1MonolithSession) TRADERESULT() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Monolith.Contract.TRADERESULT(&_TestP1Monolith.CallOpts)
}

// TRADERESULT is a free data retrieval call binding the contract method 0x495f9bff.
//
// Solidity: function _TRADE_RESULT_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Monolith *TestP1MonolithCallerSession) TRADERESULT() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Monolith.Contract.TRADERESULT(&_TestP1Monolith.CallOpts)
}

// TRADERESULT2 is a free data retrieval call binding the contract method 0x75092a30.
//
// Solidity: function _TRADE_RESULT_2_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Monolith *TestP1MonolithCaller) TRADERESULT2(opts *bind.CallOpts) (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "_TRADE_RESULT_2_")

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
func (_TestP1Monolith *TestP1MonolithSession) TRADERESULT2() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Monolith.Contract.TRADERESULT2(&_TestP1Monolith.CallOpts)
}

// TRADERESULT2 is a free data retrieval call binding the contract method 0x75092a30.
//
// Solidity: function _TRADE_RESULT_2_() view returns(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags)
func (_TestP1Monolith *TestP1MonolithCallerSession) TRADERESULT2() (struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}, error) {
	return _TestP1Monolith.Contract.TRADERESULT2(&_TestP1Monolith.CallOpts)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Monolith *TestP1MonolithCaller) GetFunding(opts *bind.CallOpts, arg0 *big.Int) (bool, *big.Int, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "getFunding", arg0)

	if err != nil {
		return *new(bool), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Monolith *TestP1MonolithSession) GetFunding(arg0 *big.Int) (bool, *big.Int, error) {
	return _TestP1Monolith.Contract.GetFunding(&_TestP1Monolith.CallOpts, arg0)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 ) view returns(bool, uint256)
func (_TestP1Monolith *TestP1MonolithCallerSession) GetFunding(arg0 *big.Int) (bool, *big.Int, error) {
	return _TestP1Monolith.Contract.GetFunding(&_TestP1Monolith.CallOpts, arg0)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestP1Monolith.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithSession) GetPrice() (*big.Int, error) {
	return _TestP1Monolith.Contract.GetPrice(&_TestP1Monolith.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_TestP1Monolith *TestP1MonolithCallerSession) GetPrice() (*big.Int, error) {
	return _TestP1Monolith.Contract.GetPrice(&_TestP1Monolith.CallOpts)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Monolith *TestP1MonolithTransactor) SetFunding(opts *bind.TransactOpts, isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.contract.Transact(opts, "setFunding", isPositive, newFunding)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Monolith *TestP1MonolithSession) SetFunding(isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetFunding(&_TestP1Monolith.TransactOpts, isPositive, newFunding)
}

// SetFunding is a paid mutator transaction binding the contract method 0xe41a054f.
//
// Solidity: function setFunding(bool isPositive, uint256 newFunding) returns()
func (_TestP1Monolith *TestP1MonolithTransactorSession) SetFunding(isPositive bool, newFunding *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetFunding(&_TestP1Monolith.TransactOpts, isPositive, newFunding)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Monolith *TestP1MonolithTransactor) SetPrice(opts *bind.TransactOpts, newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.contract.Transact(opts, "setPrice", newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Monolith *TestP1MonolithSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetPrice(&_TestP1Monolith.TransactOpts, newPrice)
}

// SetPrice is a paid mutator transaction binding the contract method 0x91b7f5ed.
//
// Solidity: function setPrice(uint256 newPrice) returns()
func (_TestP1Monolith *TestP1MonolithTransactorSession) SetPrice(newPrice *big.Int) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetPrice(&_TestP1Monolith.TransactOpts, newPrice)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithTransactor) SetSecondTradeResult(opts *bind.TransactOpts, marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.contract.Transact(opts, "setSecondTradeResult", marginAmount, positionAmount, isBuy, traderFlags)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithSession) SetSecondTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetSecondTradeResult(&_TestP1Monolith.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetSecondTradeResult is a paid mutator transaction binding the contract method 0x63a3d85f.
//
// Solidity: function setSecondTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithTransactorSession) SetSecondTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetSecondTradeResult(&_TestP1Monolith.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithTransactor) SetTradeResult(opts *bind.TransactOpts, marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.contract.Transact(opts, "setTradeResult", marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithSession) SetTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetTradeResult(&_TestP1Monolith.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// SetTradeResult is a paid mutator transaction binding the contract method 0xe53adbb2.
//
// Solidity: function setTradeResult(uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 traderFlags) returns()
func (_TestP1Monolith *TestP1MonolithTransactorSession) SetTradeResult(marginAmount *big.Int, positionAmount *big.Int, isBuy bool, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.SetTradeResult(&_TestP1Monolith.TransactOpts, marginAmount, positionAmount, isBuy, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Monolith *TestP1MonolithTransactor) Trade(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.contract.Transact(opts, "trade", arg0, arg1, arg2, arg3, arg4, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Monolith *TestP1MonolithSession) Trade(arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.Trade(&_TestP1Monolith.TransactOpts, arg0, arg1, arg2, arg3, arg4, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address , address , address , uint256 , bytes , bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_TestP1Monolith *TestP1MonolithTransactorSession) Trade(arg0 common.Address, arg1 common.Address, arg2 common.Address, arg3 *big.Int, arg4 []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _TestP1Monolith.Contract.Trade(&_TestP1Monolith.TransactOpts, arg0, arg1, arg2, arg3, arg4, traderFlags)
}
