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

// P1FinalSettlementMetaData contains all meta data concerning the P1FinalSettlement contract.
var P1FinalSettlementMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdrawFinalSettlement\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1FinalSettlementABI is the input ABI used to generate the binding from.
// Deprecated: Use P1FinalSettlementMetaData.ABI instead.
var P1FinalSettlementABI = P1FinalSettlementMetaData.ABI

// P1FinalSettlement is an auto generated Go binding around an Ethereum contract.
type P1FinalSettlement struct {
	P1FinalSettlementCaller     // Read-only binding to the contract
	P1FinalSettlementTransactor // Write-only binding to the contract
	P1FinalSettlementFilterer   // Log filterer for contract events
}

// P1FinalSettlementCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1FinalSettlementCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FinalSettlementTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1FinalSettlementTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FinalSettlementFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1FinalSettlementFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1FinalSettlementSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1FinalSettlementSession struct {
	Contract     *P1FinalSettlement // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1FinalSettlementCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1FinalSettlementCallerSession struct {
	Contract *P1FinalSettlementCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// P1FinalSettlementTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1FinalSettlementTransactorSession struct {
	Contract     *P1FinalSettlementTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// P1FinalSettlementRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1FinalSettlementRaw struct {
	Contract *P1FinalSettlement // Generic contract binding to access the raw methods on
}

// P1FinalSettlementCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1FinalSettlementCallerRaw struct {
	Contract *P1FinalSettlementCaller // Generic read-only contract binding to access the raw methods on
}

// P1FinalSettlementTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1FinalSettlementTransactorRaw struct {
	Contract *P1FinalSettlementTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1FinalSettlement creates a new instance of P1FinalSettlement, bound to a specific deployed contract.
func NewP1FinalSettlement(address common.Address, backend bind.ContractBackend) (*P1FinalSettlement, error) {
	contract, err := bindP1FinalSettlement(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlement{P1FinalSettlementCaller: P1FinalSettlementCaller{contract: contract}, P1FinalSettlementTransactor: P1FinalSettlementTransactor{contract: contract}, P1FinalSettlementFilterer: P1FinalSettlementFilterer{contract: contract}}, nil
}

// NewP1FinalSettlementCaller creates a new read-only instance of P1FinalSettlement, bound to a specific deployed contract.
func NewP1FinalSettlementCaller(address common.Address, caller bind.ContractCaller) (*P1FinalSettlementCaller, error) {
	contract, err := bindP1FinalSettlement(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementCaller{contract: contract}, nil
}

// NewP1FinalSettlementTransactor creates a new write-only instance of P1FinalSettlement, bound to a specific deployed contract.
func NewP1FinalSettlementTransactor(address common.Address, transactor bind.ContractTransactor) (*P1FinalSettlementTransactor, error) {
	contract, err := bindP1FinalSettlement(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementTransactor{contract: contract}, nil
}

// NewP1FinalSettlementFilterer creates a new log filterer instance of P1FinalSettlement, bound to a specific deployed contract.
func NewP1FinalSettlementFilterer(address common.Address, filterer bind.ContractFilterer) (*P1FinalSettlementFilterer, error) {
	contract, err := bindP1FinalSettlement(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementFilterer{contract: contract}, nil
}

// bindP1FinalSettlement binds a generic wrapper to an already deployed contract.
func bindP1FinalSettlement(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1FinalSettlementABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1FinalSettlement *P1FinalSettlementRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1FinalSettlement.Contract.P1FinalSettlementCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1FinalSettlement *P1FinalSettlementRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.P1FinalSettlementTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1FinalSettlement *P1FinalSettlementRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.P1FinalSettlementTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1FinalSettlement *P1FinalSettlementCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1FinalSettlement.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1FinalSettlement *P1FinalSettlementTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1FinalSettlement *P1FinalSettlementTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1FinalSettlement *P1FinalSettlementCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1FinalSettlement.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1FinalSettlement *P1FinalSettlementSession) GetAdmin() (common.Address, error) {
	return _P1FinalSettlement.Contract.GetAdmin(&_P1FinalSettlement.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1FinalSettlement *P1FinalSettlementCallerSession) GetAdmin() (common.Address, error) {
	return _P1FinalSettlement.Contract.GetAdmin(&_P1FinalSettlement.CallOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1FinalSettlement *P1FinalSettlementTransactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1FinalSettlement.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1FinalSettlement *P1FinalSettlementSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.WithdrawFinalSettlement(&_P1FinalSettlement.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1FinalSettlement *P1FinalSettlementTransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1FinalSettlement.Contract.WithdrawFinalSettlement(&_P1FinalSettlement.TransactOpts)
}

// P1FinalSettlementLogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the P1FinalSettlement contract.
type P1FinalSettlementLogAccountSettledIterator struct {
	Event *P1FinalSettlementLogAccountSettled // Event containing the contract specifics and raw log

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
func (it *P1FinalSettlementLogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FinalSettlementLogAccountSettled)
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
		it.Event = new(P1FinalSettlementLogAccountSettled)
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
func (it *P1FinalSettlementLogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FinalSettlementLogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FinalSettlementLogAccountSettled represents a LogAccountSettled event raised by the P1FinalSettlement contract.
type P1FinalSettlementLogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*P1FinalSettlementLogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1FinalSettlement.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementLogAccountSettledIterator{contract: _P1FinalSettlement.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *P1FinalSettlementLogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1FinalSettlement.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FinalSettlementLogAccountSettled)
				if err := _P1FinalSettlement.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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

// ParseLogAccountSettled is a log parse operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) ParseLogAccountSettled(log types.Log) (*P1FinalSettlementLogAccountSettled, error) {
	event := new(P1FinalSettlementLogAccountSettled)
	if err := _P1FinalSettlement.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1FinalSettlementLogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the P1FinalSettlement contract.
type P1FinalSettlementLogIndexIterator struct {
	Event *P1FinalSettlementLogIndex // Event containing the contract specifics and raw log

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
func (it *P1FinalSettlementLogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FinalSettlementLogIndex)
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
		it.Event = new(P1FinalSettlementLogIndex)
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
func (it *P1FinalSettlementLogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FinalSettlementLogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FinalSettlementLogIndex represents a LogIndex event raised by the P1FinalSettlement contract.
type P1FinalSettlementLogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1FinalSettlement *P1FinalSettlementFilterer) FilterLogIndex(opts *bind.FilterOpts) (*P1FinalSettlementLogIndexIterator, error) {

	logs, sub, err := _P1FinalSettlement.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementLogIndexIterator{contract: _P1FinalSettlement.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1FinalSettlement *P1FinalSettlementFilterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *P1FinalSettlementLogIndex) (event.Subscription, error) {

	logs, sub, err := _P1FinalSettlement.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FinalSettlementLogIndex)
				if err := _P1FinalSettlement.contract.UnpackLog(event, "LogIndex", log); err != nil {
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

// ParseLogIndex is a log parse operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1FinalSettlement *P1FinalSettlementFilterer) ParseLogIndex(log types.Log) (*P1FinalSettlementLogIndex, error) {
	event := new(P1FinalSettlementLogIndex)
	if err := _P1FinalSettlement.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1FinalSettlementLogWithdrawFinalSettlementIterator is returned from FilterLogWithdrawFinalSettlement and is used to iterate over the raw logs and unpacked data for LogWithdrawFinalSettlement events raised by the P1FinalSettlement contract.
type P1FinalSettlementLogWithdrawFinalSettlementIterator struct {
	Event *P1FinalSettlementLogWithdrawFinalSettlement // Event containing the contract specifics and raw log

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
func (it *P1FinalSettlementLogWithdrawFinalSettlementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1FinalSettlementLogWithdrawFinalSettlement)
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
		it.Event = new(P1FinalSettlementLogWithdrawFinalSettlement)
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
func (it *P1FinalSettlementLogWithdrawFinalSettlementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1FinalSettlementLogWithdrawFinalSettlementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1FinalSettlementLogWithdrawFinalSettlement represents a LogWithdrawFinalSettlement event raised by the P1FinalSettlement contract.
type P1FinalSettlementLogWithdrawFinalSettlement struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFinalSettlement is a free log retrieval operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) FilterLogWithdrawFinalSettlement(opts *bind.FilterOpts, account []common.Address) (*P1FinalSettlementLogWithdrawFinalSettlementIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1FinalSettlement.contract.FilterLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1FinalSettlementLogWithdrawFinalSettlementIterator{contract: _P1FinalSettlement.contract, event: "LogWithdrawFinalSettlement", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFinalSettlement is a free log subscription operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) WatchLogWithdrawFinalSettlement(opts *bind.WatchOpts, sink chan<- *P1FinalSettlementLogWithdrawFinalSettlement, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1FinalSettlement.contract.WatchLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1FinalSettlementLogWithdrawFinalSettlement)
				if err := _P1FinalSettlement.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
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

// ParseLogWithdrawFinalSettlement is a log parse operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1FinalSettlement *P1FinalSettlementFilterer) ParseLogWithdrawFinalSettlement(log types.Log) (*P1FinalSettlementLogWithdrawFinalSettlement, error) {
	event := new(P1FinalSettlementLogWithdrawFinalSettlement)
	if err := _P1FinalSettlement.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
