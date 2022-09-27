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

// P1CurrencyConverterProxyMetaData contains all meta data concerning the P1CurrencyConverterProxy contract.
var P1CurrencyConverterProxyMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchangeWrapper\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenFrom\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenFromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenToAmount\",\"type\":\"uint256\"}],\"name\":\"LogConvertedDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchangeWrapper\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenFrom\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenFromAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenToAmount\",\"type\":\"uint256\"}],\"name\":\"LogConvertedWithdrawal\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"}],\"name\":\"approveMaximumOnPerpetual\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"exchangeWrapper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenFrom\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenFromAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"exchangeWrapper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenFromAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1CurrencyConverterProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use P1CurrencyConverterProxyMetaData.ABI instead.
var P1CurrencyConverterProxyABI = P1CurrencyConverterProxyMetaData.ABI

// P1CurrencyConverterProxy is an auto generated Go binding around an Ethereum contract.
type P1CurrencyConverterProxy struct {
	P1CurrencyConverterProxyCaller     // Read-only binding to the contract
	P1CurrencyConverterProxyTransactor // Write-only binding to the contract
	P1CurrencyConverterProxyFilterer   // Log filterer for contract events
}

// P1CurrencyConverterProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1CurrencyConverterProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1CurrencyConverterProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1CurrencyConverterProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1CurrencyConverterProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1CurrencyConverterProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1CurrencyConverterProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1CurrencyConverterProxySession struct {
	Contract     *P1CurrencyConverterProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// P1CurrencyConverterProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1CurrencyConverterProxyCallerSession struct {
	Contract *P1CurrencyConverterProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// P1CurrencyConverterProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1CurrencyConverterProxyTransactorSession struct {
	Contract     *P1CurrencyConverterProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// P1CurrencyConverterProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1CurrencyConverterProxyRaw struct {
	Contract *P1CurrencyConverterProxy // Generic contract binding to access the raw methods on
}

// P1CurrencyConverterProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1CurrencyConverterProxyCallerRaw struct {
	Contract *P1CurrencyConverterProxyCaller // Generic read-only contract binding to access the raw methods on
}

// P1CurrencyConverterProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1CurrencyConverterProxyTransactorRaw struct {
	Contract *P1CurrencyConverterProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1CurrencyConverterProxy creates a new instance of P1CurrencyConverterProxy, bound to a specific deployed contract.
func NewP1CurrencyConverterProxy(address common.Address, backend bind.ContractBackend) (*P1CurrencyConverterProxy, error) {
	contract, err := bindP1CurrencyConverterProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxy{P1CurrencyConverterProxyCaller: P1CurrencyConverterProxyCaller{contract: contract}, P1CurrencyConverterProxyTransactor: P1CurrencyConverterProxyTransactor{contract: contract}, P1CurrencyConverterProxyFilterer: P1CurrencyConverterProxyFilterer{contract: contract}}, nil
}

// NewP1CurrencyConverterProxyCaller creates a new read-only instance of P1CurrencyConverterProxy, bound to a specific deployed contract.
func NewP1CurrencyConverterProxyCaller(address common.Address, caller bind.ContractCaller) (*P1CurrencyConverterProxyCaller, error) {
	contract, err := bindP1CurrencyConverterProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxyCaller{contract: contract}, nil
}

// NewP1CurrencyConverterProxyTransactor creates a new write-only instance of P1CurrencyConverterProxy, bound to a specific deployed contract.
func NewP1CurrencyConverterProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*P1CurrencyConverterProxyTransactor, error) {
	contract, err := bindP1CurrencyConverterProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxyTransactor{contract: contract}, nil
}

// NewP1CurrencyConverterProxyFilterer creates a new log filterer instance of P1CurrencyConverterProxy, bound to a specific deployed contract.
func NewP1CurrencyConverterProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*P1CurrencyConverterProxyFilterer, error) {
	contract, err := bindP1CurrencyConverterProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxyFilterer{contract: contract}, nil
}

// bindP1CurrencyConverterProxy binds a generic wrapper to an already deployed contract.
func bindP1CurrencyConverterProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1CurrencyConverterProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1CurrencyConverterProxy.Contract.P1CurrencyConverterProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.P1CurrencyConverterProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.P1CurrencyConverterProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1CurrencyConverterProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.contract.Transact(opts, method, params...)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactor) ApproveMaximumOnPerpetual(opts *bind.TransactOpts, perpetual common.Address) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.contract.Transact(opts, "approveMaximumOnPerpetual", perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxySession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.ApproveMaximumOnPerpetual(&_P1CurrencyConverterProxy.TransactOpts, perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactorSession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.ApproveMaximumOnPerpetual(&_P1CurrencyConverterProxy.TransactOpts, perpetual)
}

// Deposit is a paid mutator transaction binding the contract method 0xb13b8d0a.
//
// Solidity: function deposit(address account, address perpetual, address exchangeWrapper, address tokenFrom, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactor) Deposit(opts *bind.TransactOpts, account common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenFrom common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.contract.Transact(opts, "deposit", account, perpetual, exchangeWrapper, tokenFrom, tokenFromAmount, data)
}

// Deposit is a paid mutator transaction binding the contract method 0xb13b8d0a.
//
// Solidity: function deposit(address account, address perpetual, address exchangeWrapper, address tokenFrom, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxySession) Deposit(account common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenFrom common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.Deposit(&_P1CurrencyConverterProxy.TransactOpts, account, perpetual, exchangeWrapper, tokenFrom, tokenFromAmount, data)
}

// Deposit is a paid mutator transaction binding the contract method 0xb13b8d0a.
//
// Solidity: function deposit(address account, address perpetual, address exchangeWrapper, address tokenFrom, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactorSession) Deposit(account common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenFrom common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.Deposit(&_P1CurrencyConverterProxy.TransactOpts, account, perpetual, exchangeWrapper, tokenFrom, tokenFromAmount, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5b3901f6.
//
// Solidity: function withdraw(address account, address destination, address perpetual, address exchangeWrapper, address tokenTo, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactor) Withdraw(opts *bind.TransactOpts, account common.Address, destination common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenTo common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.contract.Transact(opts, "withdraw", account, destination, perpetual, exchangeWrapper, tokenTo, tokenFromAmount, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5b3901f6.
//
// Solidity: function withdraw(address account, address destination, address perpetual, address exchangeWrapper, address tokenTo, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxySession) Withdraw(account common.Address, destination common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenTo common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.Withdraw(&_P1CurrencyConverterProxy.TransactOpts, account, destination, perpetual, exchangeWrapper, tokenTo, tokenFromAmount, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5b3901f6.
//
// Solidity: function withdraw(address account, address destination, address perpetual, address exchangeWrapper, address tokenTo, uint256 tokenFromAmount, bytes data) returns(uint256)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyTransactorSession) Withdraw(account common.Address, destination common.Address, perpetual common.Address, exchangeWrapper common.Address, tokenTo common.Address, tokenFromAmount *big.Int, data []byte) (*types.Transaction, error) {
	return _P1CurrencyConverterProxy.Contract.Withdraw(&_P1CurrencyConverterProxy.TransactOpts, account, destination, perpetual, exchangeWrapper, tokenTo, tokenFromAmount, data)
}

// P1CurrencyConverterProxyLogConvertedDepositIterator is returned from FilterLogConvertedDeposit and is used to iterate over the raw logs and unpacked data for LogConvertedDeposit events raised by the P1CurrencyConverterProxy contract.
type P1CurrencyConverterProxyLogConvertedDepositIterator struct {
	Event *P1CurrencyConverterProxyLogConvertedDeposit // Event containing the contract specifics and raw log

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
func (it *P1CurrencyConverterProxyLogConvertedDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1CurrencyConverterProxyLogConvertedDeposit)
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
		it.Event = new(P1CurrencyConverterProxyLogConvertedDeposit)
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
func (it *P1CurrencyConverterProxyLogConvertedDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1CurrencyConverterProxyLogConvertedDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1CurrencyConverterProxyLogConvertedDeposit represents a LogConvertedDeposit event raised by the P1CurrencyConverterProxy contract.
type P1CurrencyConverterProxyLogConvertedDeposit struct {
	Account         common.Address
	Source          common.Address
	Perpetual       common.Address
	ExchangeWrapper common.Address
	TokenFrom       common.Address
	TokenTo         common.Address
	TokenFromAmount *big.Int
	TokenToAmount   *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogConvertedDeposit is a free log retrieval operation binding the contract event 0xb82979dec0b27d2050fc2ec2e499e291b29e8ce7d4cfd1f711b2557d5e609c08.
//
// Solidity: event LogConvertedDeposit(address indexed account, address source, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) FilterLogConvertedDeposit(opts *bind.FilterOpts, account []common.Address) (*P1CurrencyConverterProxyLogConvertedDepositIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1CurrencyConverterProxy.contract.FilterLogs(opts, "LogConvertedDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxyLogConvertedDepositIterator{contract: _P1CurrencyConverterProxy.contract, event: "LogConvertedDeposit", logs: logs, sub: sub}, nil
}

// WatchLogConvertedDeposit is a free log subscription operation binding the contract event 0xb82979dec0b27d2050fc2ec2e499e291b29e8ce7d4cfd1f711b2557d5e609c08.
//
// Solidity: event LogConvertedDeposit(address indexed account, address source, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) WatchLogConvertedDeposit(opts *bind.WatchOpts, sink chan<- *P1CurrencyConverterProxyLogConvertedDeposit, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1CurrencyConverterProxy.contract.WatchLogs(opts, "LogConvertedDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1CurrencyConverterProxyLogConvertedDeposit)
				if err := _P1CurrencyConverterProxy.contract.UnpackLog(event, "LogConvertedDeposit", log); err != nil {
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

// ParseLogConvertedDeposit is a log parse operation binding the contract event 0xb82979dec0b27d2050fc2ec2e499e291b29e8ce7d4cfd1f711b2557d5e609c08.
//
// Solidity: event LogConvertedDeposit(address indexed account, address source, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) ParseLogConvertedDeposit(log types.Log) (*P1CurrencyConverterProxyLogConvertedDeposit, error) {
	event := new(P1CurrencyConverterProxyLogConvertedDeposit)
	if err := _P1CurrencyConverterProxy.contract.UnpackLog(event, "LogConvertedDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1CurrencyConverterProxyLogConvertedWithdrawalIterator is returned from FilterLogConvertedWithdrawal and is used to iterate over the raw logs and unpacked data for LogConvertedWithdrawal events raised by the P1CurrencyConverterProxy contract.
type P1CurrencyConverterProxyLogConvertedWithdrawalIterator struct {
	Event *P1CurrencyConverterProxyLogConvertedWithdrawal // Event containing the contract specifics and raw log

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
func (it *P1CurrencyConverterProxyLogConvertedWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1CurrencyConverterProxyLogConvertedWithdrawal)
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
		it.Event = new(P1CurrencyConverterProxyLogConvertedWithdrawal)
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
func (it *P1CurrencyConverterProxyLogConvertedWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1CurrencyConverterProxyLogConvertedWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1CurrencyConverterProxyLogConvertedWithdrawal represents a LogConvertedWithdrawal event raised by the P1CurrencyConverterProxy contract.
type P1CurrencyConverterProxyLogConvertedWithdrawal struct {
	Account         common.Address
	Destination     common.Address
	Perpetual       common.Address
	ExchangeWrapper common.Address
	TokenFrom       common.Address
	TokenTo         common.Address
	TokenFromAmount *big.Int
	TokenToAmount   *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogConvertedWithdrawal is a free log retrieval operation binding the contract event 0x80ff1ccb23c5efa0821893d4924d2f4abe9ddf65316ec0f50606c0642f1d1367.
//
// Solidity: event LogConvertedWithdrawal(address indexed account, address destination, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) FilterLogConvertedWithdrawal(opts *bind.FilterOpts, account []common.Address) (*P1CurrencyConverterProxyLogConvertedWithdrawalIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1CurrencyConverterProxy.contract.FilterLogs(opts, "LogConvertedWithdrawal", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1CurrencyConverterProxyLogConvertedWithdrawalIterator{contract: _P1CurrencyConverterProxy.contract, event: "LogConvertedWithdrawal", logs: logs, sub: sub}, nil
}

// WatchLogConvertedWithdrawal is a free log subscription operation binding the contract event 0x80ff1ccb23c5efa0821893d4924d2f4abe9ddf65316ec0f50606c0642f1d1367.
//
// Solidity: event LogConvertedWithdrawal(address indexed account, address destination, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) WatchLogConvertedWithdrawal(opts *bind.WatchOpts, sink chan<- *P1CurrencyConverterProxyLogConvertedWithdrawal, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1CurrencyConverterProxy.contract.WatchLogs(opts, "LogConvertedWithdrawal", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1CurrencyConverterProxyLogConvertedWithdrawal)
				if err := _P1CurrencyConverterProxy.contract.UnpackLog(event, "LogConvertedWithdrawal", log); err != nil {
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

// ParseLogConvertedWithdrawal is a log parse operation binding the contract event 0x80ff1ccb23c5efa0821893d4924d2f4abe9ddf65316ec0f50606c0642f1d1367.
//
// Solidity: event LogConvertedWithdrawal(address indexed account, address destination, address perpetual, address exchangeWrapper, address tokenFrom, address tokenTo, uint256 tokenFromAmount, uint256 tokenToAmount)
func (_P1CurrencyConverterProxy *P1CurrencyConverterProxyFilterer) ParseLogConvertedWithdrawal(log types.Log) (*P1CurrencyConverterProxyLogConvertedWithdrawal, error) {
	event := new(P1CurrencyConverterProxyLogConvertedWithdrawal)
	if err := _P1CurrencyConverterProxy.contract.UnpackLog(event, "LogConvertedWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
