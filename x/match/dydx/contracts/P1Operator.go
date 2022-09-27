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

// P1OperatorMetaData contains all meta data concerning the P1Operator contract.
var P1OperatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"LogSetLocalOperator\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setLocalOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1OperatorABI is the input ABI used to generate the binding from.
// Deprecated: Use P1OperatorMetaData.ABI instead.
var P1OperatorABI = P1OperatorMetaData.ABI

// P1Operator is an auto generated Go binding around an Ethereum contract.
type P1Operator struct {
	P1OperatorCaller     // Read-only binding to the contract
	P1OperatorTransactor // Write-only binding to the contract
	P1OperatorFilterer   // Log filterer for contract events
}

// P1OperatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1OperatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OperatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1OperatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OperatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1OperatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OperatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1OperatorSession struct {
	Contract     *P1Operator       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1OperatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1OperatorCallerSession struct {
	Contract *P1OperatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// P1OperatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1OperatorTransactorSession struct {
	Contract     *P1OperatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// P1OperatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1OperatorRaw struct {
	Contract *P1Operator // Generic contract binding to access the raw methods on
}

// P1OperatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1OperatorCallerRaw struct {
	Contract *P1OperatorCaller // Generic read-only contract binding to access the raw methods on
}

// P1OperatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1OperatorTransactorRaw struct {
	Contract *P1OperatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Operator creates a new instance of P1Operator, bound to a specific deployed contract.
func NewP1Operator(address common.Address, backend bind.ContractBackend) (*P1Operator, error) {
	contract, err := bindP1Operator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Operator{P1OperatorCaller: P1OperatorCaller{contract: contract}, P1OperatorTransactor: P1OperatorTransactor{contract: contract}, P1OperatorFilterer: P1OperatorFilterer{contract: contract}}, nil
}

// NewP1OperatorCaller creates a new read-only instance of P1Operator, bound to a specific deployed contract.
func NewP1OperatorCaller(address common.Address, caller bind.ContractCaller) (*P1OperatorCaller, error) {
	contract, err := bindP1Operator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1OperatorCaller{contract: contract}, nil
}

// NewP1OperatorTransactor creates a new write-only instance of P1Operator, bound to a specific deployed contract.
func NewP1OperatorTransactor(address common.Address, transactor bind.ContractTransactor) (*P1OperatorTransactor, error) {
	contract, err := bindP1Operator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1OperatorTransactor{contract: contract}, nil
}

// NewP1OperatorFilterer creates a new log filterer instance of P1Operator, bound to a specific deployed contract.
func NewP1OperatorFilterer(address common.Address, filterer bind.ContractFilterer) (*P1OperatorFilterer, error) {
	contract, err := bindP1Operator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1OperatorFilterer{contract: contract}, nil
}

// bindP1Operator binds a generic wrapper to an already deployed contract.
func bindP1Operator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1OperatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Operator *P1OperatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Operator.Contract.P1OperatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Operator *P1OperatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Operator.Contract.P1OperatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Operator *P1OperatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Operator.Contract.P1OperatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Operator *P1OperatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Operator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Operator *P1OperatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Operator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Operator *P1OperatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Operator.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Operator *P1OperatorCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Operator.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Operator *P1OperatorSession) GetAdmin() (common.Address, error) {
	return _P1Operator.Contract.GetAdmin(&_P1Operator.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Operator *P1OperatorCallerSession) GetAdmin() (common.Address, error) {
	return _P1Operator.Contract.GetAdmin(&_P1Operator.CallOpts)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_P1Operator *P1OperatorTransactor) SetLocalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Operator.contract.Transact(opts, "setLocalOperator", operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_P1Operator *P1OperatorSession) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Operator.Contract.SetLocalOperator(&_P1Operator.TransactOpts, operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_P1Operator *P1OperatorTransactorSession) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Operator.Contract.SetLocalOperator(&_P1Operator.TransactOpts, operator, approved)
}

// P1OperatorLogSetLocalOperatorIterator is returned from FilterLogSetLocalOperator and is used to iterate over the raw logs and unpacked data for LogSetLocalOperator events raised by the P1Operator contract.
type P1OperatorLogSetLocalOperatorIterator struct {
	Event *P1OperatorLogSetLocalOperator // Event containing the contract specifics and raw log

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
func (it *P1OperatorLogSetLocalOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1OperatorLogSetLocalOperator)
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
		it.Event = new(P1OperatorLogSetLocalOperator)
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
func (it *P1OperatorLogSetLocalOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1OperatorLogSetLocalOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1OperatorLogSetLocalOperator represents a LogSetLocalOperator event raised by the P1Operator contract.
type P1OperatorLogSetLocalOperator struct {
	Sender   common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogSetLocalOperator is a free log retrieval operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_P1Operator *P1OperatorFilterer) FilterLogSetLocalOperator(opts *bind.FilterOpts, sender []common.Address) (*P1OperatorLogSetLocalOperatorIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _P1Operator.contract.FilterLogs(opts, "LogSetLocalOperator", senderRule)
	if err != nil {
		return nil, err
	}
	return &P1OperatorLogSetLocalOperatorIterator{contract: _P1Operator.contract, event: "LogSetLocalOperator", logs: logs, sub: sub}, nil
}

// WatchLogSetLocalOperator is a free log subscription operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_P1Operator *P1OperatorFilterer) WatchLogSetLocalOperator(opts *bind.WatchOpts, sink chan<- *P1OperatorLogSetLocalOperator, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _P1Operator.contract.WatchLogs(opts, "LogSetLocalOperator", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1OperatorLogSetLocalOperator)
				if err := _P1Operator.contract.UnpackLog(event, "LogSetLocalOperator", log); err != nil {
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

// ParseLogSetLocalOperator is a log parse operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_P1Operator *P1OperatorFilterer) ParseLogSetLocalOperator(log types.Log) (*P1OperatorLogSetLocalOperator, error) {
	event := new(P1OperatorLogSetLocalOperator)
	if err := _P1Operator.contract.UnpackLog(event, "LogSetLocalOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
