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

// BaseUpgradeabilityProxyMetaData contains all meta data concerning the BaseUpgradeabilityProxy contract.
var BaseUpgradeabilityProxyMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
}

// BaseUpgradeabilityProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use BaseUpgradeabilityProxyMetaData.ABI instead.
var BaseUpgradeabilityProxyABI = BaseUpgradeabilityProxyMetaData.ABI

// BaseUpgradeabilityProxy is an auto generated Go binding around an Ethereum contract.
type BaseUpgradeabilityProxy struct {
	BaseUpgradeabilityProxyCaller     // Read-only binding to the contract
	BaseUpgradeabilityProxyTransactor // Write-only binding to the contract
	BaseUpgradeabilityProxyFilterer   // Log filterer for contract events
}

// BaseUpgradeabilityProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type BaseUpgradeabilityProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseUpgradeabilityProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BaseUpgradeabilityProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseUpgradeabilityProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BaseUpgradeabilityProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseUpgradeabilityProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BaseUpgradeabilityProxySession struct {
	Contract     *BaseUpgradeabilityProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BaseUpgradeabilityProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BaseUpgradeabilityProxyCallerSession struct {
	Contract *BaseUpgradeabilityProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// BaseUpgradeabilityProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BaseUpgradeabilityProxyTransactorSession struct {
	Contract     *BaseUpgradeabilityProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// BaseUpgradeabilityProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type BaseUpgradeabilityProxyRaw struct {
	Contract *BaseUpgradeabilityProxy // Generic contract binding to access the raw methods on
}

// BaseUpgradeabilityProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BaseUpgradeabilityProxyCallerRaw struct {
	Contract *BaseUpgradeabilityProxyCaller // Generic read-only contract binding to access the raw methods on
}

// BaseUpgradeabilityProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BaseUpgradeabilityProxyTransactorRaw struct {
	Contract *BaseUpgradeabilityProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseUpgradeabilityProxy creates a new instance of BaseUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseUpgradeabilityProxy(address common.Address, backend bind.ContractBackend) (*BaseUpgradeabilityProxy, error) {
	contract, err := bindBaseUpgradeabilityProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BaseUpgradeabilityProxy{BaseUpgradeabilityProxyCaller: BaseUpgradeabilityProxyCaller{contract: contract}, BaseUpgradeabilityProxyTransactor: BaseUpgradeabilityProxyTransactor{contract: contract}, BaseUpgradeabilityProxyFilterer: BaseUpgradeabilityProxyFilterer{contract: contract}}, nil
}

// NewBaseUpgradeabilityProxyCaller creates a new read-only instance of BaseUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseUpgradeabilityProxyCaller(address common.Address, caller bind.ContractCaller) (*BaseUpgradeabilityProxyCaller, error) {
	contract, err := bindBaseUpgradeabilityProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BaseUpgradeabilityProxyCaller{contract: contract}, nil
}

// NewBaseUpgradeabilityProxyTransactor creates a new write-only instance of BaseUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseUpgradeabilityProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BaseUpgradeabilityProxyTransactor, error) {
	contract, err := bindBaseUpgradeabilityProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BaseUpgradeabilityProxyTransactor{contract: contract}, nil
}

// NewBaseUpgradeabilityProxyFilterer creates a new log filterer instance of BaseUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseUpgradeabilityProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BaseUpgradeabilityProxyFilterer, error) {
	contract, err := bindBaseUpgradeabilityProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BaseUpgradeabilityProxyFilterer{contract: contract}, nil
}

// bindBaseUpgradeabilityProxy binds a generic wrapper to an already deployed contract.
func bindBaseUpgradeabilityProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BaseUpgradeabilityProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseUpgradeabilityProxy.Contract.BaseUpgradeabilityProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.BaseUpgradeabilityProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.BaseUpgradeabilityProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseUpgradeabilityProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.Fallback(&_BaseUpgradeabilityProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BaseUpgradeabilityProxy.Contract.Fallback(&_BaseUpgradeabilityProxy.TransactOpts, calldata)
}

// BaseUpgradeabilityProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the BaseUpgradeabilityProxy contract.
type BaseUpgradeabilityProxyUpgradedIterator struct {
	Event *BaseUpgradeabilityProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *BaseUpgradeabilityProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseUpgradeabilityProxyUpgraded)
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
		it.Event = new(BaseUpgradeabilityProxyUpgraded)
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
func (it *BaseUpgradeabilityProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseUpgradeabilityProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseUpgradeabilityProxyUpgraded represents a Upgraded event raised by the BaseUpgradeabilityProxy contract.
type BaseUpgradeabilityProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BaseUpgradeabilityProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BaseUpgradeabilityProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BaseUpgradeabilityProxyUpgradedIterator{contract: _BaseUpgradeabilityProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BaseUpgradeabilityProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BaseUpgradeabilityProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseUpgradeabilityProxyUpgraded)
				if err := _BaseUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BaseUpgradeabilityProxy *BaseUpgradeabilityProxyFilterer) ParseUpgraded(log types.Log) (*BaseUpgradeabilityProxyUpgraded, error) {
	event := new(BaseUpgradeabilityProxyUpgraded)
	if err := _BaseUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
