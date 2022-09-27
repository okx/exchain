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

// P1LiquidationMetaData contains all meta data concerning the P1Liquidation contract.
var P1LiquidationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetualV1\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oraclePrice\",\"type\":\"uint256\"}],\"name\":\"LogLiquidated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_PERPETUAL_V1_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1LiquidationABI is the input ABI used to generate the binding from.
// Deprecated: Use P1LiquidationMetaData.ABI instead.
var P1LiquidationABI = P1LiquidationMetaData.ABI

// P1Liquidation is an auto generated Go binding around an Ethereum contract.
type P1Liquidation struct {
	P1LiquidationCaller     // Read-only binding to the contract
	P1LiquidationTransactor // Write-only binding to the contract
	P1LiquidationFilterer   // Log filterer for contract events
}

// P1LiquidationCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1LiquidationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1LiquidationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1LiquidationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1LiquidationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1LiquidationSession struct {
	Contract     *P1Liquidation    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1LiquidationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1LiquidationCallerSession struct {
	Contract *P1LiquidationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// P1LiquidationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1LiquidationTransactorSession struct {
	Contract     *P1LiquidationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// P1LiquidationRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1LiquidationRaw struct {
	Contract *P1Liquidation // Generic contract binding to access the raw methods on
}

// P1LiquidationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1LiquidationCallerRaw struct {
	Contract *P1LiquidationCaller // Generic read-only contract binding to access the raw methods on
}

// P1LiquidationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1LiquidationTransactorRaw struct {
	Contract *P1LiquidationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Liquidation creates a new instance of P1Liquidation, bound to a specific deployed contract.
func NewP1Liquidation(address common.Address, backend bind.ContractBackend) (*P1Liquidation, error) {
	contract, err := bindP1Liquidation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Liquidation{P1LiquidationCaller: P1LiquidationCaller{contract: contract}, P1LiquidationTransactor: P1LiquidationTransactor{contract: contract}, P1LiquidationFilterer: P1LiquidationFilterer{contract: contract}}, nil
}

// NewP1LiquidationCaller creates a new read-only instance of P1Liquidation, bound to a specific deployed contract.
func NewP1LiquidationCaller(address common.Address, caller bind.ContractCaller) (*P1LiquidationCaller, error) {
	contract, err := bindP1Liquidation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1LiquidationCaller{contract: contract}, nil
}

// NewP1LiquidationTransactor creates a new write-only instance of P1Liquidation, bound to a specific deployed contract.
func NewP1LiquidationTransactor(address common.Address, transactor bind.ContractTransactor) (*P1LiquidationTransactor, error) {
	contract, err := bindP1Liquidation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1LiquidationTransactor{contract: contract}, nil
}

// NewP1LiquidationFilterer creates a new log filterer instance of P1Liquidation, bound to a specific deployed contract.
func NewP1LiquidationFilterer(address common.Address, filterer bind.ContractFilterer) (*P1LiquidationFilterer, error) {
	contract, err := bindP1Liquidation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1LiquidationFilterer{contract: contract}, nil
}

// bindP1Liquidation binds a generic wrapper to an already deployed contract.
func bindP1Liquidation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1LiquidationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Liquidation *P1LiquidationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Liquidation.Contract.P1LiquidationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Liquidation *P1LiquidationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Liquidation.Contract.P1LiquidationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Liquidation *P1LiquidationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Liquidation.Contract.P1LiquidationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Liquidation *P1LiquidationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Liquidation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Liquidation *P1LiquidationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Liquidation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Liquidation *P1LiquidationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Liquidation.Contract.contract.Transact(opts, method, params...)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Liquidation *P1LiquidationCaller) PERPETUALV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Liquidation.contract.Call(opts, &out, "_PERPETUAL_V1_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Liquidation *P1LiquidationSession) PERPETUALV1() (common.Address, error) {
	return _P1Liquidation.Contract.PERPETUALV1(&_P1Liquidation.CallOpts)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Liquidation *P1LiquidationCallerSession) PERPETUALV1() (common.Address, error) {
	return _P1Liquidation.Contract.PERPETUALV1(&_P1Liquidation.CallOpts)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Liquidation *P1LiquidationTransactor) Trade(opts *bind.TransactOpts, sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Liquidation.contract.Transact(opts, "trade", sender, maker, taker, price, data, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Liquidation *P1LiquidationSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Liquidation.Contract.Trade(&_P1Liquidation.TransactOpts, sender, maker, taker, price, data, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Liquidation *P1LiquidationTransactorSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Liquidation.Contract.Trade(&_P1Liquidation.TransactOpts, sender, maker, taker, price, data, arg5)
}

// P1LiquidationLogLiquidatedIterator is returned from FilterLogLiquidated and is used to iterate over the raw logs and unpacked data for LogLiquidated events raised by the P1Liquidation contract.
type P1LiquidationLogLiquidatedIterator struct {
	Event *P1LiquidationLogLiquidated // Event containing the contract specifics and raw log

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
func (it *P1LiquidationLogLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1LiquidationLogLiquidated)
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
		it.Event = new(P1LiquidationLogLiquidated)
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
func (it *P1LiquidationLogLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1LiquidationLogLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1LiquidationLogLiquidated represents a LogLiquidated event raised by the P1Liquidation contract.
type P1LiquidationLogLiquidated struct {
	Maker       common.Address
	Taker       common.Address
	Amount      *big.Int
	IsBuy       bool
	OraclePrice *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogLiquidated is a free log retrieval operation binding the contract event 0x6a35c9c914a0cd15e436f6ff44611a525491fdeb755e1044e9841d7e74ba4242.
//
// Solidity: event LogLiquidated(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Liquidation *P1LiquidationFilterer) FilterLogLiquidated(opts *bind.FilterOpts, maker []common.Address, taker []common.Address) (*P1LiquidationLogLiquidatedIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Liquidation.contract.FilterLogs(opts, "LogLiquidated", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return &P1LiquidationLogLiquidatedIterator{contract: _P1Liquidation.contract, event: "LogLiquidated", logs: logs, sub: sub}, nil
}

// WatchLogLiquidated is a free log subscription operation binding the contract event 0x6a35c9c914a0cd15e436f6ff44611a525491fdeb755e1044e9841d7e74ba4242.
//
// Solidity: event LogLiquidated(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Liquidation *P1LiquidationFilterer) WatchLogLiquidated(opts *bind.WatchOpts, sink chan<- *P1LiquidationLogLiquidated, maker []common.Address, taker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Liquidation.contract.WatchLogs(opts, "LogLiquidated", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1LiquidationLogLiquidated)
				if err := _P1Liquidation.contract.UnpackLog(event, "LogLiquidated", log); err != nil {
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

// ParseLogLiquidated is a log parse operation binding the contract event 0x6a35c9c914a0cd15e436f6ff44611a525491fdeb755e1044e9841d7e74ba4242.
//
// Solidity: event LogLiquidated(address indexed maker, address indexed taker, uint256 amount, bool isBuy, uint256 oraclePrice)
func (_P1Liquidation *P1LiquidationFilterer) ParseLogLiquidated(log types.Log) (*P1LiquidationLogLiquidated, error) {
	event := new(P1LiquidationLogLiquidated)
	if err := _P1Liquidation.contract.UnpackLog(event, "LogLiquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
