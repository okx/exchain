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

// P1TypesTradeResult is an auto generated low-level Go binding around an user-defined struct.
type P1TypesTradeResult struct {
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	TraderFlags    [32]byte
}

// IP1TraderMetaData contains all meta data concerning the IP1Trader contract.
var IP1TraderMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IP1TraderABI is the input ABI used to generate the binding from.
// Deprecated: Use IP1TraderMetaData.ABI instead.
var IP1TraderABI = IP1TraderMetaData.ABI

// IP1Trader is an auto generated Go binding around an Ethereum contract.
type IP1Trader struct {
	IP1TraderCaller     // Read-only binding to the contract
	IP1TraderTransactor // Write-only binding to the contract
	IP1TraderFilterer   // Log filterer for contract events
}

// IP1TraderCaller is an auto generated read-only Go binding around an Ethereum contract.
type IP1TraderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1TraderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IP1TraderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1TraderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IP1TraderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1TraderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IP1TraderSession struct {
	Contract     *IP1Trader        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IP1TraderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IP1TraderCallerSession struct {
	Contract *IP1TraderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IP1TraderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IP1TraderTransactorSession struct {
	Contract     *IP1TraderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IP1TraderRaw is an auto generated low-level Go binding around an Ethereum contract.
type IP1TraderRaw struct {
	Contract *IP1Trader // Generic contract binding to access the raw methods on
}

// IP1TraderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IP1TraderCallerRaw struct {
	Contract *IP1TraderCaller // Generic read-only contract binding to access the raw methods on
}

// IP1TraderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IP1TraderTransactorRaw struct {
	Contract *IP1TraderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIP1Trader creates a new instance of IP1Trader, bound to a specific deployed contract.
func NewIP1Trader(address common.Address, backend bind.ContractBackend) (*IP1Trader, error) {
	contract, err := bindIP1Trader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IP1Trader{IP1TraderCaller: IP1TraderCaller{contract: contract}, IP1TraderTransactor: IP1TraderTransactor{contract: contract}, IP1TraderFilterer: IP1TraderFilterer{contract: contract}}, nil
}

// NewIP1TraderCaller creates a new read-only instance of IP1Trader, bound to a specific deployed contract.
func NewIP1TraderCaller(address common.Address, caller bind.ContractCaller) (*IP1TraderCaller, error) {
	contract, err := bindIP1Trader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IP1TraderCaller{contract: contract}, nil
}

// NewIP1TraderTransactor creates a new write-only instance of IP1Trader, bound to a specific deployed contract.
func NewIP1TraderTransactor(address common.Address, transactor bind.ContractTransactor) (*IP1TraderTransactor, error) {
	contract, err := bindIP1Trader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IP1TraderTransactor{contract: contract}, nil
}

// NewIP1TraderFilterer creates a new log filterer instance of IP1Trader, bound to a specific deployed contract.
func NewIP1TraderFilterer(address common.Address, filterer bind.ContractFilterer) (*IP1TraderFilterer, error) {
	contract, err := bindIP1Trader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IP1TraderFilterer{contract: contract}, nil
}

// bindIP1Trader binds a generic wrapper to an already deployed contract.
func bindIP1Trader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IP1TraderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Trader *IP1TraderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Trader.Contract.IP1TraderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Trader *IP1TraderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Trader.Contract.IP1TraderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Trader *IP1TraderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Trader.Contract.IP1TraderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Trader *IP1TraderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Trader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Trader *IP1TraderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Trader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Trader *IP1TraderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Trader.Contract.contract.Transact(opts, method, params...)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_IP1Trader *IP1TraderTransactor) Trade(opts *bind.TransactOpts, sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _IP1Trader.contract.Transact(opts, "trade", sender, maker, taker, price, data, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_IP1Trader *IP1TraderSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _IP1Trader.Contract.Trade(&_IP1Trader.TransactOpts, sender, maker, taker, price, data, traderFlags)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 traderFlags) returns((uint256,uint256,bool,bytes32))
func (_IP1Trader *IP1TraderTransactorSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, traderFlags [32]byte) (*types.Transaction, error) {
	return _IP1Trader.Contract.Trade(&_IP1Trader.TransactOpts, sender, maker, taker, price, data, traderFlags)
}
