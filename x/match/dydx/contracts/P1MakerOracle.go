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

// P1MakerOracleMetaData contains all meta data concerning the P1MakerOracle contract.
var P1MakerOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"adjustment\",\"type\":\"uint256\"}],\"name\":\"LogAdjustmentSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"LogRouteSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_ADJUSTMENTS_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_ROUTER_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setRoute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"adjustment\",\"type\":\"uint256\"}],\"name\":\"setAdjustment\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1MakerOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use P1MakerOracleMetaData.ABI instead.
var P1MakerOracleABI = P1MakerOracleMetaData.ABI

// P1MakerOracle is an auto generated Go binding around an Ethereum contract.
type P1MakerOracle struct {
	P1MakerOracleCaller     // Read-only binding to the contract
	P1MakerOracleTransactor // Write-only binding to the contract
	P1MakerOracleFilterer   // Log filterer for contract events
}

// P1MakerOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1MakerOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MakerOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1MakerOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MakerOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1MakerOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MakerOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1MakerOracleSession struct {
	Contract     *P1MakerOracle    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1MakerOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1MakerOracleCallerSession struct {
	Contract *P1MakerOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// P1MakerOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1MakerOracleTransactorSession struct {
	Contract     *P1MakerOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// P1MakerOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1MakerOracleRaw struct {
	Contract *P1MakerOracle // Generic contract binding to access the raw methods on
}

// P1MakerOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1MakerOracleCallerRaw struct {
	Contract *P1MakerOracleCaller // Generic read-only contract binding to access the raw methods on
}

// P1MakerOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1MakerOracleTransactorRaw struct {
	Contract *P1MakerOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1MakerOracle creates a new instance of P1MakerOracle, bound to a specific deployed contract.
func NewP1MakerOracle(address common.Address, backend bind.ContractBackend) (*P1MakerOracle, error) {
	contract, err := bindP1MakerOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracle{P1MakerOracleCaller: P1MakerOracleCaller{contract: contract}, P1MakerOracleTransactor: P1MakerOracleTransactor{contract: contract}, P1MakerOracleFilterer: P1MakerOracleFilterer{contract: contract}}, nil
}

// NewP1MakerOracleCaller creates a new read-only instance of P1MakerOracle, bound to a specific deployed contract.
func NewP1MakerOracleCaller(address common.Address, caller bind.ContractCaller) (*P1MakerOracleCaller, error) {
	contract, err := bindP1MakerOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleCaller{contract: contract}, nil
}

// NewP1MakerOracleTransactor creates a new write-only instance of P1MakerOracle, bound to a specific deployed contract.
func NewP1MakerOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*P1MakerOracleTransactor, error) {
	contract, err := bindP1MakerOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleTransactor{contract: contract}, nil
}

// NewP1MakerOracleFilterer creates a new log filterer instance of P1MakerOracle, bound to a specific deployed contract.
func NewP1MakerOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*P1MakerOracleFilterer, error) {
	contract, err := bindP1MakerOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleFilterer{contract: contract}, nil
}

// bindP1MakerOracle binds a generic wrapper to an already deployed contract.
func bindP1MakerOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1MakerOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MakerOracle *P1MakerOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MakerOracle.Contract.P1MakerOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MakerOracle *P1MakerOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.P1MakerOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MakerOracle *P1MakerOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.P1MakerOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1MakerOracle *P1MakerOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1MakerOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1MakerOracle *P1MakerOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1MakerOracle *P1MakerOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.contract.Transact(opts, method, params...)
}

// ADJUSTMENTS is a free data retrieval call binding the contract method 0x46fb89ce.
//
// Solidity: function _ADJUSTMENTS_(address ) view returns(uint256)
func (_P1MakerOracle *P1MakerOracleCaller) ADJUSTMENTS(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _P1MakerOracle.contract.Call(opts, &out, "_ADJUSTMENTS_", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ADJUSTMENTS is a free data retrieval call binding the contract method 0x46fb89ce.
//
// Solidity: function _ADJUSTMENTS_(address ) view returns(uint256)
func (_P1MakerOracle *P1MakerOracleSession) ADJUSTMENTS(arg0 common.Address) (*big.Int, error) {
	return _P1MakerOracle.Contract.ADJUSTMENTS(&_P1MakerOracle.CallOpts, arg0)
}

// ADJUSTMENTS is a free data retrieval call binding the contract method 0x46fb89ce.
//
// Solidity: function _ADJUSTMENTS_(address ) view returns(uint256)
func (_P1MakerOracle *P1MakerOracleCallerSession) ADJUSTMENTS(arg0 common.Address) (*big.Int, error) {
	return _P1MakerOracle.Contract.ADJUSTMENTS(&_P1MakerOracle.CallOpts, arg0)
}

// ROUTER is a free data retrieval call binding the contract method 0xca670a7a.
//
// Solidity: function _ROUTER_(address ) view returns(address)
func (_P1MakerOracle *P1MakerOracleCaller) ROUTER(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _P1MakerOracle.contract.Call(opts, &out, "_ROUTER_", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ROUTER is a free data retrieval call binding the contract method 0xca670a7a.
//
// Solidity: function _ROUTER_(address ) view returns(address)
func (_P1MakerOracle *P1MakerOracleSession) ROUTER(arg0 common.Address) (common.Address, error) {
	return _P1MakerOracle.Contract.ROUTER(&_P1MakerOracle.CallOpts, arg0)
}

// ROUTER is a free data retrieval call binding the contract method 0xca670a7a.
//
// Solidity: function _ROUTER_(address ) view returns(address)
func (_P1MakerOracle *P1MakerOracleCallerSession) ROUTER(arg0 common.Address) (common.Address, error) {
	return _P1MakerOracle.Contract.ROUTER(&_P1MakerOracle.CallOpts, arg0)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1MakerOracle *P1MakerOracleCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1MakerOracle.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1MakerOracle *P1MakerOracleSession) GetPrice() (*big.Int, error) {
	return _P1MakerOracle.Contract.GetPrice(&_P1MakerOracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_P1MakerOracle *P1MakerOracleCallerSession) GetPrice() (*big.Int, error) {
	return _P1MakerOracle.Contract.GetPrice(&_P1MakerOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MakerOracle *P1MakerOracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1MakerOracle.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MakerOracle *P1MakerOracleSession) IsOwner() (bool, error) {
	return _P1MakerOracle.Contract.IsOwner(&_P1MakerOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_P1MakerOracle *P1MakerOracleCallerSession) IsOwner() (bool, error) {
	return _P1MakerOracle.Contract.IsOwner(&_P1MakerOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MakerOracle *P1MakerOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1MakerOracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MakerOracle *P1MakerOracleSession) Owner() (common.Address, error) {
	return _P1MakerOracle.Contract.Owner(&_P1MakerOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_P1MakerOracle *P1MakerOracleCallerSession) Owner() (common.Address, error) {
	return _P1MakerOracle.Contract.Owner(&_P1MakerOracle.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MakerOracle *P1MakerOracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1MakerOracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MakerOracle *P1MakerOracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MakerOracle.Contract.RenounceOwnership(&_P1MakerOracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_P1MakerOracle *P1MakerOracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _P1MakerOracle.Contract.RenounceOwnership(&_P1MakerOracle.TransactOpts)
}

// SetAdjustment is a paid mutator transaction binding the contract method 0xf77b3a17.
//
// Solidity: function setAdjustment(address oracle, uint256 adjustment) returns()
func (_P1MakerOracle *P1MakerOracleTransactor) SetAdjustment(opts *bind.TransactOpts, oracle common.Address, adjustment *big.Int) (*types.Transaction, error) {
	return _P1MakerOracle.contract.Transact(opts, "setAdjustment", oracle, adjustment)
}

// SetAdjustment is a paid mutator transaction binding the contract method 0xf77b3a17.
//
// Solidity: function setAdjustment(address oracle, uint256 adjustment) returns()
func (_P1MakerOracle *P1MakerOracleSession) SetAdjustment(oracle common.Address, adjustment *big.Int) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.SetAdjustment(&_P1MakerOracle.TransactOpts, oracle, adjustment)
}

// SetAdjustment is a paid mutator transaction binding the contract method 0xf77b3a17.
//
// Solidity: function setAdjustment(address oracle, uint256 adjustment) returns()
func (_P1MakerOracle *P1MakerOracleTransactorSession) SetAdjustment(oracle common.Address, adjustment *big.Int) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.SetAdjustment(&_P1MakerOracle.TransactOpts, oracle, adjustment)
}

// SetRoute is a paid mutator transaction binding the contract method 0x0505e94d.
//
// Solidity: function setRoute(address sender, address oracle) returns()
func (_P1MakerOracle *P1MakerOracleTransactor) SetRoute(opts *bind.TransactOpts, sender common.Address, oracle common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.contract.Transact(opts, "setRoute", sender, oracle)
}

// SetRoute is a paid mutator transaction binding the contract method 0x0505e94d.
//
// Solidity: function setRoute(address sender, address oracle) returns()
func (_P1MakerOracle *P1MakerOracleSession) SetRoute(sender common.Address, oracle common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.SetRoute(&_P1MakerOracle.TransactOpts, sender, oracle)
}

// SetRoute is a paid mutator transaction binding the contract method 0x0505e94d.
//
// Solidity: function setRoute(address sender, address oracle) returns()
func (_P1MakerOracle *P1MakerOracleTransactorSession) SetRoute(sender common.Address, oracle common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.SetRoute(&_P1MakerOracle.TransactOpts, sender, oracle)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MakerOracle *P1MakerOracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MakerOracle *P1MakerOracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.TransferOwnership(&_P1MakerOracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_P1MakerOracle *P1MakerOracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _P1MakerOracle.Contract.TransferOwnership(&_P1MakerOracle.TransactOpts, newOwner)
}

// P1MakerOracleLogAdjustmentSetIterator is returned from FilterLogAdjustmentSet and is used to iterate over the raw logs and unpacked data for LogAdjustmentSet events raised by the P1MakerOracle contract.
type P1MakerOracleLogAdjustmentSetIterator struct {
	Event *P1MakerOracleLogAdjustmentSet // Event containing the contract specifics and raw log

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
func (it *P1MakerOracleLogAdjustmentSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MakerOracleLogAdjustmentSet)
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
		it.Event = new(P1MakerOracleLogAdjustmentSet)
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
func (it *P1MakerOracleLogAdjustmentSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MakerOracleLogAdjustmentSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MakerOracleLogAdjustmentSet represents a LogAdjustmentSet event raised by the P1MakerOracle contract.
type P1MakerOracleLogAdjustmentSet struct {
	Oracle     common.Address
	Adjustment *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAdjustmentSet is a free log retrieval operation binding the contract event 0x6a5ac74b5033ae621af6e2a1a99689adfe6563e782629d2b76aad70664178e21.
//
// Solidity: event LogAdjustmentSet(address indexed oracle, uint256 adjustment)
func (_P1MakerOracle *P1MakerOracleFilterer) FilterLogAdjustmentSet(opts *bind.FilterOpts, oracle []common.Address) (*P1MakerOracleLogAdjustmentSetIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _P1MakerOracle.contract.FilterLogs(opts, "LogAdjustmentSet", oracleRule)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleLogAdjustmentSetIterator{contract: _P1MakerOracle.contract, event: "LogAdjustmentSet", logs: logs, sub: sub}, nil
}

// WatchLogAdjustmentSet is a free log subscription operation binding the contract event 0x6a5ac74b5033ae621af6e2a1a99689adfe6563e782629d2b76aad70664178e21.
//
// Solidity: event LogAdjustmentSet(address indexed oracle, uint256 adjustment)
func (_P1MakerOracle *P1MakerOracleFilterer) WatchLogAdjustmentSet(opts *bind.WatchOpts, sink chan<- *P1MakerOracleLogAdjustmentSet, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _P1MakerOracle.contract.WatchLogs(opts, "LogAdjustmentSet", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MakerOracleLogAdjustmentSet)
				if err := _P1MakerOracle.contract.UnpackLog(event, "LogAdjustmentSet", log); err != nil {
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

// ParseLogAdjustmentSet is a log parse operation binding the contract event 0x6a5ac74b5033ae621af6e2a1a99689adfe6563e782629d2b76aad70664178e21.
//
// Solidity: event LogAdjustmentSet(address indexed oracle, uint256 adjustment)
func (_P1MakerOracle *P1MakerOracleFilterer) ParseLogAdjustmentSet(log types.Log) (*P1MakerOracleLogAdjustmentSet, error) {
	event := new(P1MakerOracleLogAdjustmentSet)
	if err := _P1MakerOracle.contract.UnpackLog(event, "LogAdjustmentSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MakerOracleLogRouteSetIterator is returned from FilterLogRouteSet and is used to iterate over the raw logs and unpacked data for LogRouteSet events raised by the P1MakerOracle contract.
type P1MakerOracleLogRouteSetIterator struct {
	Event *P1MakerOracleLogRouteSet // Event containing the contract specifics and raw log

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
func (it *P1MakerOracleLogRouteSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MakerOracleLogRouteSet)
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
		it.Event = new(P1MakerOracleLogRouteSet)
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
func (it *P1MakerOracleLogRouteSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MakerOracleLogRouteSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MakerOracleLogRouteSet represents a LogRouteSet event raised by the P1MakerOracle contract.
type P1MakerOracleLogRouteSet struct {
	Sender common.Address
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogRouteSet is a free log retrieval operation binding the contract event 0x16f50a0fb14e340612b259bd02dedc506aae22cf39d7215c0b8d5e85030e87b3.
//
// Solidity: event LogRouteSet(address indexed sender, address oracle)
func (_P1MakerOracle *P1MakerOracleFilterer) FilterLogRouteSet(opts *bind.FilterOpts, sender []common.Address) (*P1MakerOracleLogRouteSetIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _P1MakerOracle.contract.FilterLogs(opts, "LogRouteSet", senderRule)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleLogRouteSetIterator{contract: _P1MakerOracle.contract, event: "LogRouteSet", logs: logs, sub: sub}, nil
}

// WatchLogRouteSet is a free log subscription operation binding the contract event 0x16f50a0fb14e340612b259bd02dedc506aae22cf39d7215c0b8d5e85030e87b3.
//
// Solidity: event LogRouteSet(address indexed sender, address oracle)
func (_P1MakerOracle *P1MakerOracleFilterer) WatchLogRouteSet(opts *bind.WatchOpts, sink chan<- *P1MakerOracleLogRouteSet, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _P1MakerOracle.contract.WatchLogs(opts, "LogRouteSet", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MakerOracleLogRouteSet)
				if err := _P1MakerOracle.contract.UnpackLog(event, "LogRouteSet", log); err != nil {
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

// ParseLogRouteSet is a log parse operation binding the contract event 0x16f50a0fb14e340612b259bd02dedc506aae22cf39d7215c0b8d5e85030e87b3.
//
// Solidity: event LogRouteSet(address indexed sender, address oracle)
func (_P1MakerOracle *P1MakerOracleFilterer) ParseLogRouteSet(log types.Log) (*P1MakerOracleLogRouteSet, error) {
	event := new(P1MakerOracleLogRouteSet)
	if err := _P1MakerOracle.contract.UnpackLog(event, "LogRouteSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MakerOracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the P1MakerOracle contract.
type P1MakerOracleOwnershipTransferredIterator struct {
	Event *P1MakerOracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *P1MakerOracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MakerOracleOwnershipTransferred)
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
		it.Event = new(P1MakerOracleOwnershipTransferred)
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
func (it *P1MakerOracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MakerOracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MakerOracleOwnershipTransferred represents a OwnershipTransferred event raised by the P1MakerOracle contract.
type P1MakerOracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MakerOracle *P1MakerOracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*P1MakerOracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MakerOracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &P1MakerOracleOwnershipTransferredIterator{contract: _P1MakerOracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MakerOracle *P1MakerOracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *P1MakerOracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _P1MakerOracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MakerOracleOwnershipTransferred)
				if err := _P1MakerOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_P1MakerOracle *P1MakerOracleFilterer) ParseOwnershipTransferred(log types.Log) (*P1MakerOracleOwnershipTransferred, error) {
	event := new(P1MakerOracleOwnershipTransferred)
	if err := _P1MakerOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
