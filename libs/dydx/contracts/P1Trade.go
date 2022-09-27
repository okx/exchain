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

// P1TradeTradeArg is an auto generated low-level Go binding around an user-defined struct.
type P1TradeTradeArg struct {
	TakerIndex *big.Int
	MakerIndex *big.Int
	Trader     common.Address
	Data       []byte
}

// P1TradeMetaData contains all meta data concerning the P1Trade contract.
var P1TradeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"makerBalance\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"takerBalance\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdrawFinalSettlement\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"takerIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerIndex\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structP1Trade.TradeArg[]\",\"name\":\"trades\",\"type\":\"tuple[]\"}],\"name\":\"trade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1TradeABI is the input ABI used to generate the binding from.
// Deprecated: Use P1TradeMetaData.ABI instead.
var P1TradeABI = P1TradeMetaData.ABI

// P1Trade is an auto generated Go binding around an Ethereum contract.
type P1Trade struct {
	P1TradeCaller     // Read-only binding to the contract
	P1TradeTransactor // Write-only binding to the contract
	P1TradeFilterer   // Log filterer for contract events
}

// P1TradeCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1TradeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1TradeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1TradeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1TradeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1TradeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1TradeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1TradeSession struct {
	Contract     *P1Trade          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1TradeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1TradeCallerSession struct {
	Contract *P1TradeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// P1TradeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1TradeTransactorSession struct {
	Contract     *P1TradeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1TradeRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1TradeRaw struct {
	Contract *P1Trade // Generic contract binding to access the raw methods on
}

// P1TradeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1TradeCallerRaw struct {
	Contract *P1TradeCaller // Generic read-only contract binding to access the raw methods on
}

// P1TradeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1TradeTransactorRaw struct {
	Contract *P1TradeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Trade creates a new instance of P1Trade, bound to a specific deployed contract.
func NewP1Trade(address common.Address, backend bind.ContractBackend) (*P1Trade, error) {
	contract, err := bindP1Trade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Trade{P1TradeCaller: P1TradeCaller{contract: contract}, P1TradeTransactor: P1TradeTransactor{contract: contract}, P1TradeFilterer: P1TradeFilterer{contract: contract}}, nil
}

// NewP1TradeCaller creates a new read-only instance of P1Trade, bound to a specific deployed contract.
func NewP1TradeCaller(address common.Address, caller bind.ContractCaller) (*P1TradeCaller, error) {
	contract, err := bindP1Trade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1TradeCaller{contract: contract}, nil
}

// NewP1TradeTransactor creates a new write-only instance of P1Trade, bound to a specific deployed contract.
func NewP1TradeTransactor(address common.Address, transactor bind.ContractTransactor) (*P1TradeTransactor, error) {
	contract, err := bindP1Trade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1TradeTransactor{contract: contract}, nil
}

// NewP1TradeFilterer creates a new log filterer instance of P1Trade, bound to a specific deployed contract.
func NewP1TradeFilterer(address common.Address, filterer bind.ContractFilterer) (*P1TradeFilterer, error) {
	contract, err := bindP1Trade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1TradeFilterer{contract: contract}, nil
}

// bindP1Trade binds a generic wrapper to an already deployed contract.
func bindP1Trade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1TradeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Trade *P1TradeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Trade.Contract.P1TradeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Trade *P1TradeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Trade.Contract.P1TradeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Trade *P1TradeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Trade.Contract.P1TradeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Trade *P1TradeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Trade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Trade *P1TradeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Trade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Trade *P1TradeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Trade.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Trade *P1TradeCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Trade.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Trade *P1TradeSession) GetAdmin() (common.Address, error) {
	return _P1Trade.Contract.GetAdmin(&_P1Trade.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Trade *P1TradeCallerSession) GetAdmin() (common.Address, error) {
	return _P1Trade.Contract.GetAdmin(&_P1Trade.CallOpts)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_P1Trade *P1TradeTransactor) Trade(opts *bind.TransactOpts, accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _P1Trade.contract.Transact(opts, "trade", accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_P1Trade *P1TradeSession) Trade(accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _P1Trade.Contract.Trade(&_P1Trade.TransactOpts, accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_P1Trade *P1TradeTransactorSession) Trade(accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _P1Trade.Contract.Trade(&_P1Trade.TransactOpts, accounts, trades)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Trade *P1TradeTransactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Trade.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Trade *P1TradeSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Trade.Contract.WithdrawFinalSettlement(&_P1Trade.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Trade *P1TradeTransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Trade.Contract.WithdrawFinalSettlement(&_P1Trade.TransactOpts)
}

// P1TradeLogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the P1Trade contract.
type P1TradeLogAccountSettledIterator struct {
	Event *P1TradeLogAccountSettled // Event containing the contract specifics and raw log

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
func (it *P1TradeLogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1TradeLogAccountSettled)
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
		it.Event = new(P1TradeLogAccountSettled)
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
func (it *P1TradeLogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1TradeLogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1TradeLogAccountSettled represents a LogAccountSettled event raised by the P1Trade contract.
type P1TradeLogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Trade *P1TradeFilterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*P1TradeLogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Trade.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1TradeLogAccountSettledIterator{contract: _P1Trade.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Trade *P1TradeFilterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *P1TradeLogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Trade.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1TradeLogAccountSettled)
				if err := _P1Trade.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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
func (_P1Trade *P1TradeFilterer) ParseLogAccountSettled(log types.Log) (*P1TradeLogAccountSettled, error) {
	event := new(P1TradeLogAccountSettled)
	if err := _P1Trade.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1TradeLogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the P1Trade contract.
type P1TradeLogIndexIterator struct {
	Event *P1TradeLogIndex // Event containing the contract specifics and raw log

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
func (it *P1TradeLogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1TradeLogIndex)
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
		it.Event = new(P1TradeLogIndex)
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
func (it *P1TradeLogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1TradeLogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1TradeLogIndex represents a LogIndex event raised by the P1Trade contract.
type P1TradeLogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Trade *P1TradeFilterer) FilterLogIndex(opts *bind.FilterOpts) (*P1TradeLogIndexIterator, error) {

	logs, sub, err := _P1Trade.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &P1TradeLogIndexIterator{contract: _P1Trade.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Trade *P1TradeFilterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *P1TradeLogIndex) (event.Subscription, error) {

	logs, sub, err := _P1Trade.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1TradeLogIndex)
				if err := _P1Trade.contract.UnpackLog(event, "LogIndex", log); err != nil {
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
func (_P1Trade *P1TradeFilterer) ParseLogIndex(log types.Log) (*P1TradeLogIndex, error) {
	event := new(P1TradeLogIndex)
	if err := _P1Trade.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1TradeLogTradeIterator is returned from FilterLogTrade and is used to iterate over the raw logs and unpacked data for LogTrade events raised by the P1Trade contract.
type P1TradeLogTradeIterator struct {
	Event *P1TradeLogTrade // Event containing the contract specifics and raw log

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
func (it *P1TradeLogTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1TradeLogTrade)
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
		it.Event = new(P1TradeLogTrade)
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
func (it *P1TradeLogTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1TradeLogTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1TradeLogTrade represents a LogTrade event raised by the P1Trade contract.
type P1TradeLogTrade struct {
	Maker          common.Address
	Taker          common.Address
	Trader         common.Address
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	MakerBalance   [32]byte
	TakerBalance   [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLogTrade is a free log retrieval operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_P1Trade *P1TradeFilterer) FilterLogTrade(opts *bind.FilterOpts, maker []common.Address, taker []common.Address) (*P1TradeLogTradeIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Trade.contract.FilterLogs(opts, "LogTrade", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return &P1TradeLogTradeIterator{contract: _P1Trade.contract, event: "LogTrade", logs: logs, sub: sub}, nil
}

// WatchLogTrade is a free log subscription operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_P1Trade *P1TradeFilterer) WatchLogTrade(opts *bind.WatchOpts, sink chan<- *P1TradeLogTrade, maker []common.Address, taker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _P1Trade.contract.WatchLogs(opts, "LogTrade", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1TradeLogTrade)
				if err := _P1Trade.contract.UnpackLog(event, "LogTrade", log); err != nil {
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

// ParseLogTrade is a log parse operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_P1Trade *P1TradeFilterer) ParseLogTrade(log types.Log) (*P1TradeLogTrade, error) {
	event := new(P1TradeLogTrade)
	if err := _P1Trade.contract.UnpackLog(event, "LogTrade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1TradeLogWithdrawFinalSettlementIterator is returned from FilterLogWithdrawFinalSettlement and is used to iterate over the raw logs and unpacked data for LogWithdrawFinalSettlement events raised by the P1Trade contract.
type P1TradeLogWithdrawFinalSettlementIterator struct {
	Event *P1TradeLogWithdrawFinalSettlement // Event containing the contract specifics and raw log

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
func (it *P1TradeLogWithdrawFinalSettlementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1TradeLogWithdrawFinalSettlement)
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
		it.Event = new(P1TradeLogWithdrawFinalSettlement)
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
func (it *P1TradeLogWithdrawFinalSettlementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1TradeLogWithdrawFinalSettlementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1TradeLogWithdrawFinalSettlement represents a LogWithdrawFinalSettlement event raised by the P1Trade contract.
type P1TradeLogWithdrawFinalSettlement struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFinalSettlement is a free log retrieval operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Trade *P1TradeFilterer) FilterLogWithdrawFinalSettlement(opts *bind.FilterOpts, account []common.Address) (*P1TradeLogWithdrawFinalSettlementIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Trade.contract.FilterLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1TradeLogWithdrawFinalSettlementIterator{contract: _P1Trade.contract, event: "LogWithdrawFinalSettlement", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFinalSettlement is a free log subscription operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Trade *P1TradeFilterer) WatchLogWithdrawFinalSettlement(opts *bind.WatchOpts, sink chan<- *P1TradeLogWithdrawFinalSettlement, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Trade.contract.WatchLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1TradeLogWithdrawFinalSettlement)
				if err := _P1Trade.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
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
func (_P1Trade *P1TradeFilterer) ParseLogWithdrawFinalSettlement(log types.Log) (*P1TradeLogWithdrawFinalSettlement, error) {
	event := new(P1TradeLogWithdrawFinalSettlement)
	if err := _P1Trade.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
