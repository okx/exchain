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

// PerpetualProxyMetaData contains all meta data concerning the PerpetualProxy contract.
var PerpetualProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":false,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// PerpetualProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use PerpetualProxyMetaData.ABI instead.
var PerpetualProxyABI = PerpetualProxyMetaData.ABI

// PerpetualProxy is an auto generated Go binding around an Ethereum contract.
type PerpetualProxy struct {
	PerpetualProxyCaller     // Read-only binding to the contract
	PerpetualProxyTransactor // Write-only binding to the contract
	PerpetualProxyFilterer   // Log filterer for contract events
}

// PerpetualProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type PerpetualProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PerpetualProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PerpetualProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PerpetualProxySession struct {
	Contract     *PerpetualProxy   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PerpetualProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PerpetualProxyCallerSession struct {
	Contract *PerpetualProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PerpetualProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PerpetualProxyTransactorSession struct {
	Contract     *PerpetualProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PerpetualProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type PerpetualProxyRaw struct {
	Contract *PerpetualProxy // Generic contract binding to access the raw methods on
}

// PerpetualProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PerpetualProxyCallerRaw struct {
	Contract *PerpetualProxyCaller // Generic read-only contract binding to access the raw methods on
}

// PerpetualProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PerpetualProxyTransactorRaw struct {
	Contract *PerpetualProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPerpetualProxy creates a new instance of PerpetualProxy, bound to a specific deployed contract.
func NewPerpetualProxy(address common.Address, backend bind.ContractBackend) (*PerpetualProxy, error) {
	contract, err := bindPerpetualProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PerpetualProxy{PerpetualProxyCaller: PerpetualProxyCaller{contract: contract}, PerpetualProxyTransactor: PerpetualProxyTransactor{contract: contract}, PerpetualProxyFilterer: PerpetualProxyFilterer{contract: contract}}, nil
}

// NewPerpetualProxyCaller creates a new read-only instance of PerpetualProxy, bound to a specific deployed contract.
func NewPerpetualProxyCaller(address common.Address, caller bind.ContractCaller) (*PerpetualProxyCaller, error) {
	contract, err := bindPerpetualProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PerpetualProxyCaller{contract: contract}, nil
}

// NewPerpetualProxyTransactor creates a new write-only instance of PerpetualProxy, bound to a specific deployed contract.
func NewPerpetualProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*PerpetualProxyTransactor, error) {
	contract, err := bindPerpetualProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PerpetualProxyTransactor{contract: contract}, nil
}

// NewPerpetualProxyFilterer creates a new log filterer instance of PerpetualProxy, bound to a specific deployed contract.
func NewPerpetualProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*PerpetualProxyFilterer, error) {
	contract, err := bindPerpetualProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PerpetualProxyFilterer{contract: contract}, nil
}

// bindPerpetualProxy binds a generic wrapper to an already deployed contract.
func bindPerpetualProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PerpetualProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpetualProxy *PerpetualProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpetualProxy.Contract.PerpetualProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpetualProxy *PerpetualProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.PerpetualProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpetualProxy *PerpetualProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.PerpetualProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpetualProxy *PerpetualProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpetualProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpetualProxy *PerpetualProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpetualProxy *PerpetualProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_PerpetualProxy *PerpetualProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualProxy.contract.Transact(opts, "admin")
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_PerpetualProxy *PerpetualProxySession) Admin() (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Admin(&_PerpetualProxy.TransactOpts)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_PerpetualProxy *PerpetualProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Admin(&_PerpetualProxy.TransactOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_PerpetualProxy *PerpetualProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_PerpetualProxy *PerpetualProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.ChangeAdmin(&_PerpetualProxy.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_PerpetualProxy *PerpetualProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.ChangeAdmin(&_PerpetualProxy.TransactOpts, newAdmin)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_PerpetualProxy *PerpetualProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualProxy.contract.Transact(opts, "implementation")
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_PerpetualProxy *PerpetualProxySession) Implementation() (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Implementation(&_PerpetualProxy.TransactOpts)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_PerpetualProxy *PerpetualProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Implementation(&_PerpetualProxy.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PerpetualProxy *PerpetualProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PerpetualProxy *PerpetualProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.UpgradeTo(&_PerpetualProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PerpetualProxy *PerpetualProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.UpgradeTo(&_PerpetualProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PerpetualProxy *PerpetualProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PerpetualProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PerpetualProxy *PerpetualProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.UpgradeToAndCall(&_PerpetualProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PerpetualProxy *PerpetualProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.UpgradeToAndCall(&_PerpetualProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_PerpetualProxy *PerpetualProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _PerpetualProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_PerpetualProxy *PerpetualProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Fallback(&_PerpetualProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_PerpetualProxy *PerpetualProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _PerpetualProxy.Contract.Fallback(&_PerpetualProxy.TransactOpts, calldata)
}

// PerpetualProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the PerpetualProxy contract.
type PerpetualProxyAdminChangedIterator struct {
	Event *PerpetualProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *PerpetualProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualProxyAdminChanged)
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
		it.Event = new(PerpetualProxyAdminChanged)
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
func (it *PerpetualProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualProxyAdminChanged represents a AdminChanged event raised by the PerpetualProxy contract.
type PerpetualProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PerpetualProxy *PerpetualProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*PerpetualProxyAdminChangedIterator, error) {

	logs, sub, err := _PerpetualProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &PerpetualProxyAdminChangedIterator{contract: _PerpetualProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PerpetualProxy *PerpetualProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *PerpetualProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _PerpetualProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualProxyAdminChanged)
				if err := _PerpetualProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PerpetualProxy *PerpetualProxyFilterer) ParseAdminChanged(log types.Log) (*PerpetualProxyAdminChanged, error) {
	event := new(PerpetualProxyAdminChanged)
	if err := _PerpetualProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the PerpetualProxy contract.
type PerpetualProxyUpgradedIterator struct {
	Event *PerpetualProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *PerpetualProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualProxyUpgraded)
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
		it.Event = new(PerpetualProxyUpgraded)
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
func (it *PerpetualProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualProxyUpgraded represents a Upgraded event raised by the PerpetualProxy contract.
type PerpetualProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PerpetualProxy *PerpetualProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PerpetualProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PerpetualProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualProxyUpgradedIterator{contract: _PerpetualProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PerpetualProxy *PerpetualProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PerpetualProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PerpetualProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualProxyUpgraded)
				if err := _PerpetualProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_PerpetualProxy *PerpetualProxyFilterer) ParseUpgraded(log types.Log) (*PerpetualProxyUpgraded, error) {
	event := new(PerpetualProxyUpgraded)
	if err := _PerpetualProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
