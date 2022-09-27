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

// AdminUpgradeabilityProxyMetaData contains all meta data concerning the AdminUpgradeabilityProxy contract.
var AdminUpgradeabilityProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":false,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// AdminUpgradeabilityProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use AdminUpgradeabilityProxyMetaData.ABI instead.
var AdminUpgradeabilityProxyABI = AdminUpgradeabilityProxyMetaData.ABI

// AdminUpgradeabilityProxy is an auto generated Go binding around an Ethereum contract.
type AdminUpgradeabilityProxy struct {
	AdminUpgradeabilityProxyCaller     // Read-only binding to the contract
	AdminUpgradeabilityProxyTransactor // Write-only binding to the contract
	AdminUpgradeabilityProxyFilterer   // Log filterer for contract events
}

// AdminUpgradeabilityProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdminUpgradeabilityProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminUpgradeabilityProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdminUpgradeabilityProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminUpgradeabilityProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdminUpgradeabilityProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminUpgradeabilityProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdminUpgradeabilityProxySession struct {
	Contract     *AdminUpgradeabilityProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AdminUpgradeabilityProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdminUpgradeabilityProxyCallerSession struct {
	Contract *AdminUpgradeabilityProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// AdminUpgradeabilityProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdminUpgradeabilityProxyTransactorSession struct {
	Contract     *AdminUpgradeabilityProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// AdminUpgradeabilityProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdminUpgradeabilityProxyRaw struct {
	Contract *AdminUpgradeabilityProxy // Generic contract binding to access the raw methods on
}

// AdminUpgradeabilityProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdminUpgradeabilityProxyCallerRaw struct {
	Contract *AdminUpgradeabilityProxyCaller // Generic read-only contract binding to access the raw methods on
}

// AdminUpgradeabilityProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdminUpgradeabilityProxyTransactorRaw struct {
	Contract *AdminUpgradeabilityProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdminUpgradeabilityProxy creates a new instance of AdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewAdminUpgradeabilityProxy(address common.Address, backend bind.ContractBackend) (*AdminUpgradeabilityProxy, error) {
	contract, err := bindAdminUpgradeabilityProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxy{AdminUpgradeabilityProxyCaller: AdminUpgradeabilityProxyCaller{contract: contract}, AdminUpgradeabilityProxyTransactor: AdminUpgradeabilityProxyTransactor{contract: contract}, AdminUpgradeabilityProxyFilterer: AdminUpgradeabilityProxyFilterer{contract: contract}}, nil
}

// NewAdminUpgradeabilityProxyCaller creates a new read-only instance of AdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewAdminUpgradeabilityProxyCaller(address common.Address, caller bind.ContractCaller) (*AdminUpgradeabilityProxyCaller, error) {
	contract, err := bindAdminUpgradeabilityProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxyCaller{contract: contract}, nil
}

// NewAdminUpgradeabilityProxyTransactor creates a new write-only instance of AdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewAdminUpgradeabilityProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*AdminUpgradeabilityProxyTransactor, error) {
	contract, err := bindAdminUpgradeabilityProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxyTransactor{contract: contract}, nil
}

// NewAdminUpgradeabilityProxyFilterer creates a new log filterer instance of AdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewAdminUpgradeabilityProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*AdminUpgradeabilityProxyFilterer, error) {
	contract, err := bindAdminUpgradeabilityProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxyFilterer{contract: contract}, nil
}

// bindAdminUpgradeabilityProxy binds a generic wrapper to an already deployed contract.
func bindAdminUpgradeabilityProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AdminUpgradeabilityProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdminUpgradeabilityProxy.Contract.AdminUpgradeabilityProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.AdminUpgradeabilityProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.AdminUpgradeabilityProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AdminUpgradeabilityProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.Transact(opts, "admin")
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) Admin() (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Admin(&_AdminUpgradeabilityProxy.TransactOpts)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Admin(&_AdminUpgradeabilityProxy.TransactOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.ChangeAdmin(&_AdminUpgradeabilityProxy.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.ChangeAdmin(&_AdminUpgradeabilityProxy.TransactOpts, newAdmin)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.Transact(opts, "implementation")
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) Implementation() (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Implementation(&_AdminUpgradeabilityProxy.TransactOpts)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Implementation(&_AdminUpgradeabilityProxy.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.UpgradeTo(&_AdminUpgradeabilityProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.UpgradeTo(&_AdminUpgradeabilityProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.UpgradeToAndCall(&_AdminUpgradeabilityProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.UpgradeToAndCall(&_AdminUpgradeabilityProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Fallback(&_AdminUpgradeabilityProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AdminUpgradeabilityProxy.Contract.Fallback(&_AdminUpgradeabilityProxy.TransactOpts, calldata)
}

// AdminUpgradeabilityProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the AdminUpgradeabilityProxy contract.
type AdminUpgradeabilityProxyAdminChangedIterator struct {
	Event *AdminUpgradeabilityProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *AdminUpgradeabilityProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdminUpgradeabilityProxyAdminChanged)
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
		it.Event = new(AdminUpgradeabilityProxyAdminChanged)
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
func (it *AdminUpgradeabilityProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdminUpgradeabilityProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdminUpgradeabilityProxyAdminChanged represents a AdminChanged event raised by the AdminUpgradeabilityProxy contract.
type AdminUpgradeabilityProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*AdminUpgradeabilityProxyAdminChangedIterator, error) {

	logs, sub, err := _AdminUpgradeabilityProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxyAdminChangedIterator{contract: _AdminUpgradeabilityProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *AdminUpgradeabilityProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _AdminUpgradeabilityProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdminUpgradeabilityProxyAdminChanged)
				if err := _AdminUpgradeabilityProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) ParseAdminChanged(log types.Log) (*AdminUpgradeabilityProxyAdminChanged, error) {
	event := new(AdminUpgradeabilityProxyAdminChanged)
	if err := _AdminUpgradeabilityProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AdminUpgradeabilityProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AdminUpgradeabilityProxy contract.
type AdminUpgradeabilityProxyUpgradedIterator struct {
	Event *AdminUpgradeabilityProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *AdminUpgradeabilityProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AdminUpgradeabilityProxyUpgraded)
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
		it.Event = new(AdminUpgradeabilityProxyUpgraded)
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
func (it *AdminUpgradeabilityProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AdminUpgradeabilityProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AdminUpgradeabilityProxyUpgraded represents a Upgraded event raised by the AdminUpgradeabilityProxy contract.
type AdminUpgradeabilityProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AdminUpgradeabilityProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AdminUpgradeabilityProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AdminUpgradeabilityProxyUpgradedIterator{contract: _AdminUpgradeabilityProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AdminUpgradeabilityProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AdminUpgradeabilityProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AdminUpgradeabilityProxyUpgraded)
				if err := _AdminUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AdminUpgradeabilityProxy *AdminUpgradeabilityProxyFilterer) ParseUpgraded(log types.Log) (*AdminUpgradeabilityProxyUpgraded, error) {
	event := new(AdminUpgradeabilityProxyUpgraded)
	if err := _AdminUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
