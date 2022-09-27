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

// P1AdminMetaData contains all meta data concerning the P1Admin contract.
var P1AdminMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"settlementPrice\",\"type\":\"uint256\"}],\"name\":\"LogFinalSettlementEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"}],\"name\":\"LogSetFunder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"LogSetGlobalOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minCollateral\",\"type\":\"uint256\"}],\"name\":\"LogSetMinCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"LogSetOracle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdrawFinalSettlement\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setGlobalOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"}],\"name\":\"setFunder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minCollateral\",\"type\":\"uint256\"}],\"name\":\"setMinCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"priceLowerBound\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceUpperBound\",\"type\":\"uint256\"}],\"name\":\"enableFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1AdminABI is the input ABI used to generate the binding from.
// Deprecated: Use P1AdminMetaData.ABI instead.
var P1AdminABI = P1AdminMetaData.ABI

// P1Admin is an auto generated Go binding around an Ethereum contract.
type P1Admin struct {
	P1AdminCaller     // Read-only binding to the contract
	P1AdminTransactor // Write-only binding to the contract
	P1AdminFilterer   // Log filterer for contract events
}

// P1AdminCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1AdminCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1AdminTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1AdminTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1AdminFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1AdminFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1AdminSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1AdminSession struct {
	Contract     *P1Admin          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1AdminCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1AdminCallerSession struct {
	Contract *P1AdminCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// P1AdminTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1AdminTransactorSession struct {
	Contract     *P1AdminTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1AdminRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1AdminRaw struct {
	Contract *P1Admin // Generic contract binding to access the raw methods on
}

// P1AdminCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1AdminCallerRaw struct {
	Contract *P1AdminCaller // Generic read-only contract binding to access the raw methods on
}

// P1AdminTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1AdminTransactorRaw struct {
	Contract *P1AdminTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Admin creates a new instance of P1Admin, bound to a specific deployed contract.
func NewP1Admin(address common.Address, backend bind.ContractBackend) (*P1Admin, error) {
	contract, err := bindP1Admin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Admin{P1AdminCaller: P1AdminCaller{contract: contract}, P1AdminTransactor: P1AdminTransactor{contract: contract}, P1AdminFilterer: P1AdminFilterer{contract: contract}}, nil
}

// NewP1AdminCaller creates a new read-only instance of P1Admin, bound to a specific deployed contract.
func NewP1AdminCaller(address common.Address, caller bind.ContractCaller) (*P1AdminCaller, error) {
	contract, err := bindP1Admin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1AdminCaller{contract: contract}, nil
}

// NewP1AdminTransactor creates a new write-only instance of P1Admin, bound to a specific deployed contract.
func NewP1AdminTransactor(address common.Address, transactor bind.ContractTransactor) (*P1AdminTransactor, error) {
	contract, err := bindP1Admin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1AdminTransactor{contract: contract}, nil
}

// NewP1AdminFilterer creates a new log filterer instance of P1Admin, bound to a specific deployed contract.
func NewP1AdminFilterer(address common.Address, filterer bind.ContractFilterer) (*P1AdminFilterer, error) {
	contract, err := bindP1Admin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1AdminFilterer{contract: contract}, nil
}

// bindP1Admin binds a generic wrapper to an already deployed contract.
func bindP1Admin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1AdminABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Admin *P1AdminRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Admin.Contract.P1AdminCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Admin *P1AdminRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Admin.Contract.P1AdminTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Admin *P1AdminRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Admin.Contract.P1AdminTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Admin *P1AdminCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Admin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Admin *P1AdminTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Admin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Admin *P1AdminTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Admin.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Admin *P1AdminCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Admin.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Admin *P1AdminSession) GetAdmin() (common.Address, error) {
	return _P1Admin.Contract.GetAdmin(&_P1Admin.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Admin *P1AdminCallerSession) GetAdmin() (common.Address, error) {
	return _P1Admin.Contract.GetAdmin(&_P1Admin.CallOpts)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_P1Admin *P1AdminTransactor) EnableFinalSettlement(opts *bind.TransactOpts, priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "enableFinalSettlement", priceLowerBound, priceUpperBound)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_P1Admin *P1AdminSession) EnableFinalSettlement(priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _P1Admin.Contract.EnableFinalSettlement(&_P1Admin.TransactOpts, priceLowerBound, priceUpperBound)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_P1Admin *P1AdminTransactorSession) EnableFinalSettlement(priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _P1Admin.Contract.EnableFinalSettlement(&_P1Admin.TransactOpts, priceLowerBound, priceUpperBound)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_P1Admin *P1AdminTransactor) SetFunder(opts *bind.TransactOpts, funder common.Address) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "setFunder", funder)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_P1Admin *P1AdminSession) SetFunder(funder common.Address) (*types.Transaction, error) {
	return _P1Admin.Contract.SetFunder(&_P1Admin.TransactOpts, funder)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_P1Admin *P1AdminTransactorSession) SetFunder(funder common.Address) (*types.Transaction, error) {
	return _P1Admin.Contract.SetFunder(&_P1Admin.TransactOpts, funder)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_P1Admin *P1AdminTransactor) SetGlobalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "setGlobalOperator", operator, approved)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_P1Admin *P1AdminSession) SetGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Admin.Contract.SetGlobalOperator(&_P1Admin.TransactOpts, operator, approved)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_P1Admin *P1AdminTransactorSession) SetGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _P1Admin.Contract.SetGlobalOperator(&_P1Admin.TransactOpts, operator, approved)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_P1Admin *P1AdminTransactor) SetMinCollateral(opts *bind.TransactOpts, minCollateral *big.Int) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "setMinCollateral", minCollateral)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_P1Admin *P1AdminSession) SetMinCollateral(minCollateral *big.Int) (*types.Transaction, error) {
	return _P1Admin.Contract.SetMinCollateral(&_P1Admin.TransactOpts, minCollateral)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_P1Admin *P1AdminTransactorSession) SetMinCollateral(minCollateral *big.Int) (*types.Transaction, error) {
	return _P1Admin.Contract.SetMinCollateral(&_P1Admin.TransactOpts, minCollateral)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_P1Admin *P1AdminTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_P1Admin *P1AdminSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _P1Admin.Contract.SetOracle(&_P1Admin.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_P1Admin *P1AdminTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _P1Admin.Contract.SetOracle(&_P1Admin.TransactOpts, oracle)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Admin *P1AdminTransactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Admin.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Admin *P1AdminSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Admin.Contract.WithdrawFinalSettlement(&_P1Admin.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Admin *P1AdminTransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Admin.Contract.WithdrawFinalSettlement(&_P1Admin.TransactOpts)
}

// P1AdminLogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the P1Admin contract.
type P1AdminLogAccountSettledIterator struct {
	Event *P1AdminLogAccountSettled // Event containing the contract specifics and raw log

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
func (it *P1AdminLogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogAccountSettled)
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
		it.Event = new(P1AdminLogAccountSettled)
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
func (it *P1AdminLogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogAccountSettled represents a LogAccountSettled event raised by the P1Admin contract.
type P1AdminLogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Admin *P1AdminFilterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*P1AdminLogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1AdminLogAccountSettledIterator{contract: _P1Admin.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Admin *P1AdminFilterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *P1AdminLogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogAccountSettled)
				if err := _P1Admin.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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
func (_P1Admin *P1AdminFilterer) ParseLogAccountSettled(log types.Log) (*P1AdminLogAccountSettled, error) {
	event := new(P1AdminLogAccountSettled)
	if err := _P1Admin.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogFinalSettlementEnabledIterator is returned from FilterLogFinalSettlementEnabled and is used to iterate over the raw logs and unpacked data for LogFinalSettlementEnabled events raised by the P1Admin contract.
type P1AdminLogFinalSettlementEnabledIterator struct {
	Event *P1AdminLogFinalSettlementEnabled // Event containing the contract specifics and raw log

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
func (it *P1AdminLogFinalSettlementEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogFinalSettlementEnabled)
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
		it.Event = new(P1AdminLogFinalSettlementEnabled)
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
func (it *P1AdminLogFinalSettlementEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogFinalSettlementEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogFinalSettlementEnabled represents a LogFinalSettlementEnabled event raised by the P1Admin contract.
type P1AdminLogFinalSettlementEnabled struct {
	SettlementPrice *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogFinalSettlementEnabled is a free log retrieval operation binding the contract event 0x68e4c41627e835051be46337f1542645a60c7e6d6ea79efc5f20bdadae5f88d2.
//
// Solidity: event LogFinalSettlementEnabled(uint256 settlementPrice)
func (_P1Admin *P1AdminFilterer) FilterLogFinalSettlementEnabled(opts *bind.FilterOpts) (*P1AdminLogFinalSettlementEnabledIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogFinalSettlementEnabled")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogFinalSettlementEnabledIterator{contract: _P1Admin.contract, event: "LogFinalSettlementEnabled", logs: logs, sub: sub}, nil
}

// WatchLogFinalSettlementEnabled is a free log subscription operation binding the contract event 0x68e4c41627e835051be46337f1542645a60c7e6d6ea79efc5f20bdadae5f88d2.
//
// Solidity: event LogFinalSettlementEnabled(uint256 settlementPrice)
func (_P1Admin *P1AdminFilterer) WatchLogFinalSettlementEnabled(opts *bind.WatchOpts, sink chan<- *P1AdminLogFinalSettlementEnabled) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogFinalSettlementEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogFinalSettlementEnabled)
				if err := _P1Admin.contract.UnpackLog(event, "LogFinalSettlementEnabled", log); err != nil {
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

// ParseLogFinalSettlementEnabled is a log parse operation binding the contract event 0x68e4c41627e835051be46337f1542645a60c7e6d6ea79efc5f20bdadae5f88d2.
//
// Solidity: event LogFinalSettlementEnabled(uint256 settlementPrice)
func (_P1Admin *P1AdminFilterer) ParseLogFinalSettlementEnabled(log types.Log) (*P1AdminLogFinalSettlementEnabled, error) {
	event := new(P1AdminLogFinalSettlementEnabled)
	if err := _P1Admin.contract.UnpackLog(event, "LogFinalSettlementEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the P1Admin contract.
type P1AdminLogIndexIterator struct {
	Event *P1AdminLogIndex // Event containing the contract specifics and raw log

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
func (it *P1AdminLogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogIndex)
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
		it.Event = new(P1AdminLogIndex)
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
func (it *P1AdminLogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogIndex represents a LogIndex event raised by the P1Admin contract.
type P1AdminLogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Admin *P1AdminFilterer) FilterLogIndex(opts *bind.FilterOpts) (*P1AdminLogIndexIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogIndexIterator{contract: _P1Admin.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Admin *P1AdminFilterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *P1AdminLogIndex) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogIndex)
				if err := _P1Admin.contract.UnpackLog(event, "LogIndex", log); err != nil {
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
func (_P1Admin *P1AdminFilterer) ParseLogIndex(log types.Log) (*P1AdminLogIndex, error) {
	event := new(P1AdminLogIndex)
	if err := _P1Admin.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogSetFunderIterator is returned from FilterLogSetFunder and is used to iterate over the raw logs and unpacked data for LogSetFunder events raised by the P1Admin contract.
type P1AdminLogSetFunderIterator struct {
	Event *P1AdminLogSetFunder // Event containing the contract specifics and raw log

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
func (it *P1AdminLogSetFunderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogSetFunder)
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
		it.Event = new(P1AdminLogSetFunder)
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
func (it *P1AdminLogSetFunderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogSetFunderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogSetFunder represents a LogSetFunder event raised by the P1Admin contract.
type P1AdminLogSetFunder struct {
	Funder common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogSetFunder is a free log retrieval operation binding the contract event 0x433b5c8c9ff78f62114ee8804a916537fa42009ebac4965bfed953f771789e47.
//
// Solidity: event LogSetFunder(address funder)
func (_P1Admin *P1AdminFilterer) FilterLogSetFunder(opts *bind.FilterOpts) (*P1AdminLogSetFunderIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogSetFunder")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogSetFunderIterator{contract: _P1Admin.contract, event: "LogSetFunder", logs: logs, sub: sub}, nil
}

// WatchLogSetFunder is a free log subscription operation binding the contract event 0x433b5c8c9ff78f62114ee8804a916537fa42009ebac4965bfed953f771789e47.
//
// Solidity: event LogSetFunder(address funder)
func (_P1Admin *P1AdminFilterer) WatchLogSetFunder(opts *bind.WatchOpts, sink chan<- *P1AdminLogSetFunder) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogSetFunder")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogSetFunder)
				if err := _P1Admin.contract.UnpackLog(event, "LogSetFunder", log); err != nil {
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

// ParseLogSetFunder is a log parse operation binding the contract event 0x433b5c8c9ff78f62114ee8804a916537fa42009ebac4965bfed953f771789e47.
//
// Solidity: event LogSetFunder(address funder)
func (_P1Admin *P1AdminFilterer) ParseLogSetFunder(log types.Log) (*P1AdminLogSetFunder, error) {
	event := new(P1AdminLogSetFunder)
	if err := _P1Admin.contract.UnpackLog(event, "LogSetFunder", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogSetGlobalOperatorIterator is returned from FilterLogSetGlobalOperator and is used to iterate over the raw logs and unpacked data for LogSetGlobalOperator events raised by the P1Admin contract.
type P1AdminLogSetGlobalOperatorIterator struct {
	Event *P1AdminLogSetGlobalOperator // Event containing the contract specifics and raw log

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
func (it *P1AdminLogSetGlobalOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogSetGlobalOperator)
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
		it.Event = new(P1AdminLogSetGlobalOperator)
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
func (it *P1AdminLogSetGlobalOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogSetGlobalOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogSetGlobalOperator represents a LogSetGlobalOperator event raised by the P1Admin contract.
type P1AdminLogSetGlobalOperator struct {
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogSetGlobalOperator is a free log retrieval operation binding the contract event 0xeaeee7699e70e6b31ac89ec999ef6936b97ac1e364f0e1fcf584772372caa0d3.
//
// Solidity: event LogSetGlobalOperator(address operator, bool approved)
func (_P1Admin *P1AdminFilterer) FilterLogSetGlobalOperator(opts *bind.FilterOpts) (*P1AdminLogSetGlobalOperatorIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogSetGlobalOperator")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogSetGlobalOperatorIterator{contract: _P1Admin.contract, event: "LogSetGlobalOperator", logs: logs, sub: sub}, nil
}

// WatchLogSetGlobalOperator is a free log subscription operation binding the contract event 0xeaeee7699e70e6b31ac89ec999ef6936b97ac1e364f0e1fcf584772372caa0d3.
//
// Solidity: event LogSetGlobalOperator(address operator, bool approved)
func (_P1Admin *P1AdminFilterer) WatchLogSetGlobalOperator(opts *bind.WatchOpts, sink chan<- *P1AdminLogSetGlobalOperator) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogSetGlobalOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogSetGlobalOperator)
				if err := _P1Admin.contract.UnpackLog(event, "LogSetGlobalOperator", log); err != nil {
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

// ParseLogSetGlobalOperator is a log parse operation binding the contract event 0xeaeee7699e70e6b31ac89ec999ef6936b97ac1e364f0e1fcf584772372caa0d3.
//
// Solidity: event LogSetGlobalOperator(address operator, bool approved)
func (_P1Admin *P1AdminFilterer) ParseLogSetGlobalOperator(log types.Log) (*P1AdminLogSetGlobalOperator, error) {
	event := new(P1AdminLogSetGlobalOperator)
	if err := _P1Admin.contract.UnpackLog(event, "LogSetGlobalOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogSetMinCollateralIterator is returned from FilterLogSetMinCollateral and is used to iterate over the raw logs and unpacked data for LogSetMinCollateral events raised by the P1Admin contract.
type P1AdminLogSetMinCollateralIterator struct {
	Event *P1AdminLogSetMinCollateral // Event containing the contract specifics and raw log

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
func (it *P1AdminLogSetMinCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogSetMinCollateral)
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
		it.Event = new(P1AdminLogSetMinCollateral)
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
func (it *P1AdminLogSetMinCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogSetMinCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogSetMinCollateral represents a LogSetMinCollateral event raised by the P1Admin contract.
type P1AdminLogSetMinCollateral struct {
	MinCollateral *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogSetMinCollateral is a free log retrieval operation binding the contract event 0x248b36ced4662a14c995e0872f00eb61be4e3dea3913226cdeb513d64728cdca.
//
// Solidity: event LogSetMinCollateral(uint256 minCollateral)
func (_P1Admin *P1AdminFilterer) FilterLogSetMinCollateral(opts *bind.FilterOpts) (*P1AdminLogSetMinCollateralIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogSetMinCollateral")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogSetMinCollateralIterator{contract: _P1Admin.contract, event: "LogSetMinCollateral", logs: logs, sub: sub}, nil
}

// WatchLogSetMinCollateral is a free log subscription operation binding the contract event 0x248b36ced4662a14c995e0872f00eb61be4e3dea3913226cdeb513d64728cdca.
//
// Solidity: event LogSetMinCollateral(uint256 minCollateral)
func (_P1Admin *P1AdminFilterer) WatchLogSetMinCollateral(opts *bind.WatchOpts, sink chan<- *P1AdminLogSetMinCollateral) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogSetMinCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogSetMinCollateral)
				if err := _P1Admin.contract.UnpackLog(event, "LogSetMinCollateral", log); err != nil {
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

// ParseLogSetMinCollateral is a log parse operation binding the contract event 0x248b36ced4662a14c995e0872f00eb61be4e3dea3913226cdeb513d64728cdca.
//
// Solidity: event LogSetMinCollateral(uint256 minCollateral)
func (_P1Admin *P1AdminFilterer) ParseLogSetMinCollateral(log types.Log) (*P1AdminLogSetMinCollateral, error) {
	event := new(P1AdminLogSetMinCollateral)
	if err := _P1Admin.contract.UnpackLog(event, "LogSetMinCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogSetOracleIterator is returned from FilterLogSetOracle and is used to iterate over the raw logs and unpacked data for LogSetOracle events raised by the P1Admin contract.
type P1AdminLogSetOracleIterator struct {
	Event *P1AdminLogSetOracle // Event containing the contract specifics and raw log

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
func (it *P1AdminLogSetOracleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogSetOracle)
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
		it.Event = new(P1AdminLogSetOracle)
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
func (it *P1AdminLogSetOracleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogSetOracleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogSetOracle represents a LogSetOracle event raised by the P1Admin contract.
type P1AdminLogSetOracle struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogSetOracle is a free log retrieval operation binding the contract event 0xad675642c3cba5442815383698d42cd28889533d9671a6d32cffea58ef0874da.
//
// Solidity: event LogSetOracle(address oracle)
func (_P1Admin *P1AdminFilterer) FilterLogSetOracle(opts *bind.FilterOpts) (*P1AdminLogSetOracleIterator, error) {

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogSetOracle")
	if err != nil {
		return nil, err
	}
	return &P1AdminLogSetOracleIterator{contract: _P1Admin.contract, event: "LogSetOracle", logs: logs, sub: sub}, nil
}

// WatchLogSetOracle is a free log subscription operation binding the contract event 0xad675642c3cba5442815383698d42cd28889533d9671a6d32cffea58ef0874da.
//
// Solidity: event LogSetOracle(address oracle)
func (_P1Admin *P1AdminFilterer) WatchLogSetOracle(opts *bind.WatchOpts, sink chan<- *P1AdminLogSetOracle) (event.Subscription, error) {

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogSetOracle")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogSetOracle)
				if err := _P1Admin.contract.UnpackLog(event, "LogSetOracle", log); err != nil {
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

// ParseLogSetOracle is a log parse operation binding the contract event 0xad675642c3cba5442815383698d42cd28889533d9671a6d32cffea58ef0874da.
//
// Solidity: event LogSetOracle(address oracle)
func (_P1Admin *P1AdminFilterer) ParseLogSetOracle(log types.Log) (*P1AdminLogSetOracle, error) {
	event := new(P1AdminLogSetOracle)
	if err := _P1Admin.contract.UnpackLog(event, "LogSetOracle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1AdminLogWithdrawFinalSettlementIterator is returned from FilterLogWithdrawFinalSettlement and is used to iterate over the raw logs and unpacked data for LogWithdrawFinalSettlement events raised by the P1Admin contract.
type P1AdminLogWithdrawFinalSettlementIterator struct {
	Event *P1AdminLogWithdrawFinalSettlement // Event containing the contract specifics and raw log

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
func (it *P1AdminLogWithdrawFinalSettlementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1AdminLogWithdrawFinalSettlement)
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
		it.Event = new(P1AdminLogWithdrawFinalSettlement)
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
func (it *P1AdminLogWithdrawFinalSettlementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1AdminLogWithdrawFinalSettlementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1AdminLogWithdrawFinalSettlement represents a LogWithdrawFinalSettlement event raised by the P1Admin contract.
type P1AdminLogWithdrawFinalSettlement struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFinalSettlement is a free log retrieval operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Admin *P1AdminFilterer) FilterLogWithdrawFinalSettlement(opts *bind.FilterOpts, account []common.Address) (*P1AdminLogWithdrawFinalSettlementIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Admin.contract.FilterLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1AdminLogWithdrawFinalSettlementIterator{contract: _P1Admin.contract, event: "LogWithdrawFinalSettlement", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFinalSettlement is a free log subscription operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Admin *P1AdminFilterer) WatchLogWithdrawFinalSettlement(opts *bind.WatchOpts, sink chan<- *P1AdminLogWithdrawFinalSettlement, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Admin.contract.WatchLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1AdminLogWithdrawFinalSettlement)
				if err := _P1Admin.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
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
func (_P1Admin *P1AdminFilterer) ParseLogWithdrawFinalSettlement(log types.Log) (*P1AdminLogWithdrawFinalSettlement, error) {
	event := new(P1AdminLogWithdrawFinalSettlement)
	if err := _P1Admin.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
