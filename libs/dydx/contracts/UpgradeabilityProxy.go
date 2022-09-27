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

// UpgradeabilityProxyMetaData contains all meta data concerning the UpgradeabilityProxy contract.
var UpgradeabilityProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
}

// UpgradeabilityProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeabilityProxyMetaData.ABI instead.
var UpgradeabilityProxyABI = UpgradeabilityProxyMetaData.ABI

// UpgradeabilityProxy is an auto generated Go binding around an Ethereum contract.
type UpgradeabilityProxy struct {
	UpgradeabilityProxyCaller     // Read-only binding to the contract
	UpgradeabilityProxyTransactor // Write-only binding to the contract
	UpgradeabilityProxyFilterer   // Log filterer for contract events
}

// UpgradeabilityProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeabilityProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeabilityProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeabilityProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeabilityProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeabilityProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeabilityProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeabilityProxySession struct {
	Contract     *UpgradeabilityProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// UpgradeabilityProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeabilityProxyCallerSession struct {
	Contract *UpgradeabilityProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// UpgradeabilityProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeabilityProxyTransactorSession struct {
	Contract     *UpgradeabilityProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// UpgradeabilityProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeabilityProxyRaw struct {
	Contract *UpgradeabilityProxy // Generic contract binding to access the raw methods on
}

// UpgradeabilityProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeabilityProxyCallerRaw struct {
	Contract *UpgradeabilityProxyCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeabilityProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeabilityProxyTransactorRaw struct {
	Contract *UpgradeabilityProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeabilityProxy creates a new instance of UpgradeabilityProxy, bound to a specific deployed contract.
func NewUpgradeabilityProxy(address common.Address, backend bind.ContractBackend) (*UpgradeabilityProxy, error) {
	contract, err := bindUpgradeabilityProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradeabilityProxy{UpgradeabilityProxyCaller: UpgradeabilityProxyCaller{contract: contract}, UpgradeabilityProxyTransactor: UpgradeabilityProxyTransactor{contract: contract}, UpgradeabilityProxyFilterer: UpgradeabilityProxyFilterer{contract: contract}}, nil
}

// NewUpgradeabilityProxyCaller creates a new read-only instance of UpgradeabilityProxy, bound to a specific deployed contract.
func NewUpgradeabilityProxyCaller(address common.Address, caller bind.ContractCaller) (*UpgradeabilityProxyCaller, error) {
	contract, err := bindUpgradeabilityProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeabilityProxyCaller{contract: contract}, nil
}

// NewUpgradeabilityProxyTransactor creates a new write-only instance of UpgradeabilityProxy, bound to a specific deployed contract.
func NewUpgradeabilityProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeabilityProxyTransactor, error) {
	contract, err := bindUpgradeabilityProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeabilityProxyTransactor{contract: contract}, nil
}

// NewUpgradeabilityProxyFilterer creates a new log filterer instance of UpgradeabilityProxy, bound to a specific deployed contract.
func NewUpgradeabilityProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeabilityProxyFilterer, error) {
	contract, err := bindUpgradeabilityProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeabilityProxyFilterer{contract: contract}, nil
}

// bindUpgradeabilityProxy binds a generic wrapper to an already deployed contract.
func bindUpgradeabilityProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeabilityProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeabilityProxy *UpgradeabilityProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeabilityProxy.Contract.UpgradeabilityProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeabilityProxy *UpgradeabilityProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.UpgradeabilityProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeabilityProxy *UpgradeabilityProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.UpgradeabilityProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeabilityProxy *UpgradeabilityProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeabilityProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeabilityProxy *UpgradeabilityProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeabilityProxy *UpgradeabilityProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_UpgradeabilityProxy *UpgradeabilityProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _UpgradeabilityProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_UpgradeabilityProxy *UpgradeabilityProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.Fallback(&_UpgradeabilityProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_UpgradeabilityProxy *UpgradeabilityProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _UpgradeabilityProxy.Contract.Fallback(&_UpgradeabilityProxy.TransactOpts, calldata)
}

// UpgradeabilityProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the UpgradeabilityProxy contract.
type UpgradeabilityProxyUpgradedIterator struct {
	Event *UpgradeabilityProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *UpgradeabilityProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeabilityProxyUpgraded)
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
		it.Event = new(UpgradeabilityProxyUpgraded)
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
func (it *UpgradeabilityProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeabilityProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeabilityProxyUpgraded represents a Upgraded event raised by the UpgradeabilityProxy contract.
type UpgradeabilityProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UpgradeabilityProxy *UpgradeabilityProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*UpgradeabilityProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UpgradeabilityProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeabilityProxyUpgradedIterator{contract: _UpgradeabilityProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UpgradeabilityProxy *UpgradeabilityProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UpgradeabilityProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UpgradeabilityProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeabilityProxyUpgraded)
				if err := _UpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_UpgradeabilityProxy *UpgradeabilityProxyFilterer) ParseUpgraded(log types.Log) (*UpgradeabilityProxyUpgraded, error) {
	event := new(UpgradeabilityProxyUpgraded)
	if err := _UpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
