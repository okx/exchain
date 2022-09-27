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

// BaseAdminUpgradeabilityProxyMetaData contains all meta data concerning the BaseAdminUpgradeabilityProxy contract.
var BaseAdminUpgradeabilityProxyMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":false,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// BaseAdminUpgradeabilityProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use BaseAdminUpgradeabilityProxyMetaData.ABI instead.
var BaseAdminUpgradeabilityProxyABI = BaseAdminUpgradeabilityProxyMetaData.ABI

// BaseAdminUpgradeabilityProxy is an auto generated Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxy struct {
	BaseAdminUpgradeabilityProxyCaller     // Read-only binding to the contract
	BaseAdminUpgradeabilityProxyTransactor // Write-only binding to the contract
	BaseAdminUpgradeabilityProxyFilterer   // Log filterer for contract events
}

// BaseAdminUpgradeabilityProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseAdminUpgradeabilityProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseAdminUpgradeabilityProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BaseAdminUpgradeabilityProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseAdminUpgradeabilityProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BaseAdminUpgradeabilityProxySession struct {
	Contract     *BaseAdminUpgradeabilityProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// BaseAdminUpgradeabilityProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BaseAdminUpgradeabilityProxyCallerSession struct {
	Contract *BaseAdminUpgradeabilityProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// BaseAdminUpgradeabilityProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BaseAdminUpgradeabilityProxyTransactorSession struct {
	Contract     *BaseAdminUpgradeabilityProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// BaseAdminUpgradeabilityProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxyRaw struct {
	Contract *BaseAdminUpgradeabilityProxy // Generic contract binding to access the raw methods on
}

// BaseAdminUpgradeabilityProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxyCallerRaw struct {
	Contract *BaseAdminUpgradeabilityProxyCaller // Generic read-only contract binding to access the raw methods on
}

// BaseAdminUpgradeabilityProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BaseAdminUpgradeabilityProxyTransactorRaw struct {
	Contract *BaseAdminUpgradeabilityProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseAdminUpgradeabilityProxy creates a new instance of BaseAdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseAdminUpgradeabilityProxy(address common.Address, backend bind.ContractBackend) (*BaseAdminUpgradeabilityProxy, error) {
	contract, err := bindBaseAdminUpgradeabilityProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxy{BaseAdminUpgradeabilityProxyCaller: BaseAdminUpgradeabilityProxyCaller{contract: contract}, BaseAdminUpgradeabilityProxyTransactor: BaseAdminUpgradeabilityProxyTransactor{contract: contract}, BaseAdminUpgradeabilityProxyFilterer: BaseAdminUpgradeabilityProxyFilterer{contract: contract}}, nil
}

// NewBaseAdminUpgradeabilityProxyCaller creates a new read-only instance of BaseAdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseAdminUpgradeabilityProxyCaller(address common.Address, caller bind.ContractCaller) (*BaseAdminUpgradeabilityProxyCaller, error) {
	contract, err := bindBaseAdminUpgradeabilityProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxyCaller{contract: contract}, nil
}

// NewBaseAdminUpgradeabilityProxyTransactor creates a new write-only instance of BaseAdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseAdminUpgradeabilityProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BaseAdminUpgradeabilityProxyTransactor, error) {
	contract, err := bindBaseAdminUpgradeabilityProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxyTransactor{contract: contract}, nil
}

// NewBaseAdminUpgradeabilityProxyFilterer creates a new log filterer instance of BaseAdminUpgradeabilityProxy, bound to a specific deployed contract.
func NewBaseAdminUpgradeabilityProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BaseAdminUpgradeabilityProxyFilterer, error) {
	contract, err := bindBaseAdminUpgradeabilityProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxyFilterer{contract: contract}, nil
}

// bindBaseAdminUpgradeabilityProxy binds a generic wrapper to an already deployed contract.
func bindBaseAdminUpgradeabilityProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BaseAdminUpgradeabilityProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseAdminUpgradeabilityProxy.Contract.BaseAdminUpgradeabilityProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.BaseAdminUpgradeabilityProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.BaseAdminUpgradeabilityProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseAdminUpgradeabilityProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.contract.Transact(opts, method, params...)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) Admin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.Transact(opts, "admin")
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) Admin() (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Admin(&_BaseAdminUpgradeabilityProxy.TransactOpts)
}

// Admin is a paid mutator transaction binding the contract method 0xf851a440.
//
// Solidity: function admin() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) Admin() (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Admin(&_BaseAdminUpgradeabilityProxy.TransactOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.ChangeAdmin(&_BaseAdminUpgradeabilityProxy.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.ChangeAdmin(&_BaseAdminUpgradeabilityProxy.TransactOpts, newAdmin)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) Implementation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.Transact(opts, "implementation")
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) Implementation() (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Implementation(&_BaseAdminUpgradeabilityProxy.TransactOpts)
}

// Implementation is a paid mutator transaction binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() returns(address)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) Implementation() (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Implementation(&_BaseAdminUpgradeabilityProxy.TransactOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.UpgradeTo(&_BaseAdminUpgradeabilityProxy.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.UpgradeTo(&_BaseAdminUpgradeabilityProxy.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.UpgradeToAndCall(&_BaseAdminUpgradeabilityProxy.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.UpgradeToAndCall(&_BaseAdminUpgradeabilityProxy.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Fallback(&_BaseAdminUpgradeabilityProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BaseAdminUpgradeabilityProxy.Contract.Fallback(&_BaseAdminUpgradeabilityProxy.TransactOpts, calldata)
}

// BaseAdminUpgradeabilityProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the BaseAdminUpgradeabilityProxy contract.
type BaseAdminUpgradeabilityProxyAdminChangedIterator struct {
	Event *BaseAdminUpgradeabilityProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *BaseAdminUpgradeabilityProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseAdminUpgradeabilityProxyAdminChanged)
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
		it.Event = new(BaseAdminUpgradeabilityProxyAdminChanged)
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
func (it *BaseAdminUpgradeabilityProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseAdminUpgradeabilityProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseAdminUpgradeabilityProxyAdminChanged represents a AdminChanged event raised by the BaseAdminUpgradeabilityProxy contract.
type BaseAdminUpgradeabilityProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*BaseAdminUpgradeabilityProxyAdminChangedIterator, error) {

	logs, sub, err := _BaseAdminUpgradeabilityProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxyAdminChangedIterator{contract: _BaseAdminUpgradeabilityProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *BaseAdminUpgradeabilityProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _BaseAdminUpgradeabilityProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseAdminUpgradeabilityProxyAdminChanged)
				if err := _BaseAdminUpgradeabilityProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) ParseAdminChanged(log types.Log) (*BaseAdminUpgradeabilityProxyAdminChanged, error) {
	event := new(BaseAdminUpgradeabilityProxyAdminChanged)
	if err := _BaseAdminUpgradeabilityProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseAdminUpgradeabilityProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the BaseAdminUpgradeabilityProxy contract.
type BaseAdminUpgradeabilityProxyUpgradedIterator struct {
	Event *BaseAdminUpgradeabilityProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *BaseAdminUpgradeabilityProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseAdminUpgradeabilityProxyUpgraded)
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
		it.Event = new(BaseAdminUpgradeabilityProxyUpgraded)
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
func (it *BaseAdminUpgradeabilityProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseAdminUpgradeabilityProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseAdminUpgradeabilityProxyUpgraded represents a Upgraded event raised by the BaseAdminUpgradeabilityProxy contract.
type BaseAdminUpgradeabilityProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BaseAdminUpgradeabilityProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BaseAdminUpgradeabilityProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BaseAdminUpgradeabilityProxyUpgradedIterator{contract: _BaseAdminUpgradeabilityProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BaseAdminUpgradeabilityProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BaseAdminUpgradeabilityProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseAdminUpgradeabilityProxyUpgraded)
				if err := _BaseAdminUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_BaseAdminUpgradeabilityProxy *BaseAdminUpgradeabilityProxyFilterer) ParseUpgraded(log types.Log) (*BaseAdminUpgradeabilityProxyUpgraded, error) {
	event := new(BaseAdminUpgradeabilityProxyUpgraded)
	if err := _BaseAdminUpgradeabilityProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
