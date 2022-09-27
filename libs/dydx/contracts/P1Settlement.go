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

// P1SettlementMetaData contains all meta data concerning the P1Settlement contract.
var P1SettlementMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1SettlementABI is the input ABI used to generate the binding from.
// Deprecated: Use P1SettlementMetaData.ABI instead.
var P1SettlementABI = P1SettlementMetaData.ABI

// P1Settlement is an auto generated Go binding around an Ethereum contract.
type P1Settlement struct {
	P1SettlementCaller     // Read-only binding to the contract
	P1SettlementTransactor // Write-only binding to the contract
	P1SettlementFilterer   // Log filterer for contract events
}

// P1SettlementCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1SettlementCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SettlementTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1SettlementTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SettlementFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1SettlementFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SettlementSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1SettlementSession struct {
	Contract     *P1Settlement     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1SettlementCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1SettlementCallerSession struct {
	Contract *P1SettlementCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// P1SettlementTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1SettlementTransactorSession struct {
	Contract     *P1SettlementTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// P1SettlementRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1SettlementRaw struct {
	Contract *P1Settlement // Generic contract binding to access the raw methods on
}

// P1SettlementCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1SettlementCallerRaw struct {
	Contract *P1SettlementCaller // Generic read-only contract binding to access the raw methods on
}

// P1SettlementTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1SettlementTransactorRaw struct {
	Contract *P1SettlementTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Settlement creates a new instance of P1Settlement, bound to a specific deployed contract.
func NewP1Settlement(address common.Address, backend bind.ContractBackend) (*P1Settlement, error) {
	contract, err := bindP1Settlement(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Settlement{P1SettlementCaller: P1SettlementCaller{contract: contract}, P1SettlementTransactor: P1SettlementTransactor{contract: contract}, P1SettlementFilterer: P1SettlementFilterer{contract: contract}}, nil
}

// NewP1SettlementCaller creates a new read-only instance of P1Settlement, bound to a specific deployed contract.
func NewP1SettlementCaller(address common.Address, caller bind.ContractCaller) (*P1SettlementCaller, error) {
	contract, err := bindP1Settlement(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1SettlementCaller{contract: contract}, nil
}

// NewP1SettlementTransactor creates a new write-only instance of P1Settlement, bound to a specific deployed contract.
func NewP1SettlementTransactor(address common.Address, transactor bind.ContractTransactor) (*P1SettlementTransactor, error) {
	contract, err := bindP1Settlement(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1SettlementTransactor{contract: contract}, nil
}

// NewP1SettlementFilterer creates a new log filterer instance of P1Settlement, bound to a specific deployed contract.
func NewP1SettlementFilterer(address common.Address, filterer bind.ContractFilterer) (*P1SettlementFilterer, error) {
	contract, err := bindP1Settlement(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1SettlementFilterer{contract: contract}, nil
}

// bindP1Settlement binds a generic wrapper to an already deployed contract.
func bindP1Settlement(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1SettlementABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Settlement *P1SettlementRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Settlement.Contract.P1SettlementCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Settlement *P1SettlementRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Settlement.Contract.P1SettlementTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Settlement *P1SettlementRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Settlement.Contract.P1SettlementTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Settlement *P1SettlementCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Settlement.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Settlement *P1SettlementTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Settlement.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Settlement *P1SettlementTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Settlement.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Settlement *P1SettlementCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Settlement.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Settlement *P1SettlementSession) GetAdmin() (common.Address, error) {
	return _P1Settlement.Contract.GetAdmin(&_P1Settlement.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Settlement *P1SettlementCallerSession) GetAdmin() (common.Address, error) {
	return _P1Settlement.Contract.GetAdmin(&_P1Settlement.CallOpts)
}

// P1SettlementLogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the P1Settlement contract.
type P1SettlementLogAccountSettledIterator struct {
	Event *P1SettlementLogAccountSettled // Event containing the contract specifics and raw log

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
func (it *P1SettlementLogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1SettlementLogAccountSettled)
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
		it.Event = new(P1SettlementLogAccountSettled)
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
func (it *P1SettlementLogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1SettlementLogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1SettlementLogAccountSettled represents a LogAccountSettled event raised by the P1Settlement contract.
type P1SettlementLogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Settlement *P1SettlementFilterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*P1SettlementLogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Settlement.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1SettlementLogAccountSettledIterator{contract: _P1Settlement.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Settlement *P1SettlementFilterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *P1SettlementLogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Settlement.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1SettlementLogAccountSettled)
				if err := _P1Settlement.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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
func (_P1Settlement *P1SettlementFilterer) ParseLogAccountSettled(log types.Log) (*P1SettlementLogAccountSettled, error) {
	event := new(P1SettlementLogAccountSettled)
	if err := _P1Settlement.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1SettlementLogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the P1Settlement contract.
type P1SettlementLogIndexIterator struct {
	Event *P1SettlementLogIndex // Event containing the contract specifics and raw log

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
func (it *P1SettlementLogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1SettlementLogIndex)
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
		it.Event = new(P1SettlementLogIndex)
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
func (it *P1SettlementLogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1SettlementLogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1SettlementLogIndex represents a LogIndex event raised by the P1Settlement contract.
type P1SettlementLogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Settlement *P1SettlementFilterer) FilterLogIndex(opts *bind.FilterOpts) (*P1SettlementLogIndexIterator, error) {

	logs, sub, err := _P1Settlement.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &P1SettlementLogIndexIterator{contract: _P1Settlement.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Settlement *P1SettlementFilterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *P1SettlementLogIndex) (event.Subscription, error) {

	logs, sub, err := _P1Settlement.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1SettlementLogIndex)
				if err := _P1Settlement.contract.UnpackLog(event, "LogIndex", log); err != nil {
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
func (_P1Settlement *P1SettlementFilterer) ParseLogIndex(log types.Log) (*P1SettlementLogIndex, error) {
	event := new(P1SettlementLogIndex)
	if err := _P1Settlement.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
