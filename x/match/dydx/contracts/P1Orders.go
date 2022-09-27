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

// P1OrdersFill is an auto generated low-level Go binding around an user-defined struct.
type P1OrdersFill struct {
	Amount        *big.Int
	Price         *big.Int
	Fee           *big.Int
	IsNegativeFee bool
}

// P1OrdersOrder is an auto generated low-level Go binding around an user-defined struct.
type P1OrdersOrder struct {
	Flags        [32]byte
	Amount       *big.Int
	LimitPrice   *big.Int
	TriggerPrice *big.Int
	LimitFee     *big.Int
	Maker        common.Address
	Taker        common.Address
	Expiration   *big.Int
}

// P1OrdersOrderQueryOutput is an auto generated low-level Go binding around an user-defined struct.
type P1OrdersOrderQueryOutput struct {
	Status       uint8
	FilledAmount *big.Int
}

// P1OrdersMetaData contains all meta data concerning the P1Orders contract.
var P1OrdersMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetualV1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogOrderApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogOrderCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"triggerPrice\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isNegativeFee\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structP1Orders.Fill\",\"name\":\"fill\",\"type\":\"tuple\"}],\"name\":\"LogOrderFilled\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_EIP712_DOMAIN_HASH_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_FILLED_AMOUNT_\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_PERPETUAL_V1_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_STATUS_\",\"outputs\":[{\"internalType\":\"enumP1Orders.OrderStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"trade\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"traderFlags\",\"type\":\"bytes32\"}],\"internalType\":\"structP1Types.TradeResult\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"triggerPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structP1Orders.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"approveOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"triggerPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structP1Orders.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"orderHashes\",\"type\":\"bytes32[]\"}],\"name\":\"getOrdersStatus\",\"outputs\":[{\"components\":[{\"internalType\":\"enumP1Orders.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"filledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structP1Orders.OrderQueryOutput[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1OrdersABI is the input ABI used to generate the binding from.
// Deprecated: Use P1OrdersMetaData.ABI instead.
var P1OrdersABI = P1OrdersMetaData.ABI

// P1Orders is an auto generated Go binding around an Ethereum contract.
type P1Orders struct {
	P1OrdersCaller     // Read-only binding to the contract
	P1OrdersTransactor // Write-only binding to the contract
	P1OrdersFilterer   // Log filterer for contract events
}

// P1OrdersCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1OrdersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OrdersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1OrdersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OrdersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1OrdersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1OrdersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1OrdersSession struct {
	Contract     *P1Orders         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1OrdersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1OrdersCallerSession struct {
	Contract *P1OrdersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// P1OrdersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1OrdersTransactorSession struct {
	Contract     *P1OrdersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// P1OrdersRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1OrdersRaw struct {
	Contract *P1Orders // Generic contract binding to access the raw methods on
}

// P1OrdersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1OrdersCallerRaw struct {
	Contract *P1OrdersCaller // Generic read-only contract binding to access the raw methods on
}

// P1OrdersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1OrdersTransactorRaw struct {
	Contract *P1OrdersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Orders creates a new instance of P1Orders, bound to a specific deployed contract.
func NewP1Orders(address common.Address, backend bind.ContractBackend) (*P1Orders, error) {
	contract, err := bindP1Orders(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Orders{P1OrdersCaller: P1OrdersCaller{contract: contract}, P1OrdersTransactor: P1OrdersTransactor{contract: contract}, P1OrdersFilterer: P1OrdersFilterer{contract: contract}}, nil
}

// NewP1OrdersCaller creates a new read-only instance of P1Orders, bound to a specific deployed contract.
func NewP1OrdersCaller(address common.Address, caller bind.ContractCaller) (*P1OrdersCaller, error) {
	contract, err := bindP1Orders(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1OrdersCaller{contract: contract}, nil
}

// NewP1OrdersTransactor creates a new write-only instance of P1Orders, bound to a specific deployed contract.
func NewP1OrdersTransactor(address common.Address, transactor bind.ContractTransactor) (*P1OrdersTransactor, error) {
	contract, err := bindP1Orders(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1OrdersTransactor{contract: contract}, nil
}

// NewP1OrdersFilterer creates a new log filterer instance of P1Orders, bound to a specific deployed contract.
func NewP1OrdersFilterer(address common.Address, filterer bind.ContractFilterer) (*P1OrdersFilterer, error) {
	contract, err := bindP1Orders(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1OrdersFilterer{contract: contract}, nil
}

// bindP1Orders binds a generic wrapper to an already deployed contract.
func bindP1Orders(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1OrdersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Orders *P1OrdersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Orders.Contract.P1OrdersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Orders *P1OrdersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Orders.Contract.P1OrdersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Orders *P1OrdersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Orders.Contract.P1OrdersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Orders *P1OrdersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Orders.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Orders *P1OrdersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Orders.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Orders *P1OrdersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Orders.Contract.contract.Transact(opts, method, params...)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1Orders *P1OrdersCaller) EIP712DOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _P1Orders.contract.Call(opts, &out, "_EIP712_DOMAIN_HASH_")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1Orders *P1OrdersSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _P1Orders.Contract.EIP712DOMAINHASH(&_P1Orders.CallOpts)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1Orders *P1OrdersCallerSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _P1Orders.Contract.EIP712DOMAINHASH(&_P1Orders.CallOpts)
}

// FILLEDAMOUNT is a free data retrieval call binding the contract method 0x5c457f29.
//
// Solidity: function _FILLED_AMOUNT_(bytes32 ) view returns(uint256)
func (_P1Orders *P1OrdersCaller) FILLEDAMOUNT(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _P1Orders.contract.Call(opts, &out, "_FILLED_AMOUNT_", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FILLEDAMOUNT is a free data retrieval call binding the contract method 0x5c457f29.
//
// Solidity: function _FILLED_AMOUNT_(bytes32 ) view returns(uint256)
func (_P1Orders *P1OrdersSession) FILLEDAMOUNT(arg0 [32]byte) (*big.Int, error) {
	return _P1Orders.Contract.FILLEDAMOUNT(&_P1Orders.CallOpts, arg0)
}

// FILLEDAMOUNT is a free data retrieval call binding the contract method 0x5c457f29.
//
// Solidity: function _FILLED_AMOUNT_(bytes32 ) view returns(uint256)
func (_P1Orders *P1OrdersCallerSession) FILLEDAMOUNT(arg0 [32]byte) (*big.Int, error) {
	return _P1Orders.Contract.FILLEDAMOUNT(&_P1Orders.CallOpts, arg0)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Orders *P1OrdersCaller) PERPETUALV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Orders.contract.Call(opts, &out, "_PERPETUAL_V1_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Orders *P1OrdersSession) PERPETUALV1() (common.Address, error) {
	return _P1Orders.Contract.PERPETUALV1(&_P1Orders.CallOpts)
}

// PERPETUALV1 is a free data retrieval call binding the contract method 0xd4bec8eb.
//
// Solidity: function _PERPETUAL_V1_() view returns(address)
func (_P1Orders *P1OrdersCallerSession) PERPETUALV1() (common.Address, error) {
	return _P1Orders.Contract.PERPETUALV1(&_P1Orders.CallOpts)
}

// STATUS is a free data retrieval call binding the contract method 0x9ea07071.
//
// Solidity: function _STATUS_(bytes32 ) view returns(uint8)
func (_P1Orders *P1OrdersCaller) STATUS(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var out []interface{}
	err := _P1Orders.contract.Call(opts, &out, "_STATUS_", arg0)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// STATUS is a free data retrieval call binding the contract method 0x9ea07071.
//
// Solidity: function _STATUS_(bytes32 ) view returns(uint8)
func (_P1Orders *P1OrdersSession) STATUS(arg0 [32]byte) (uint8, error) {
	return _P1Orders.Contract.STATUS(&_P1Orders.CallOpts, arg0)
}

// STATUS is a free data retrieval call binding the contract method 0x9ea07071.
//
// Solidity: function _STATUS_(bytes32 ) view returns(uint8)
func (_P1Orders *P1OrdersCallerSession) STATUS(arg0 [32]byte) (uint8, error) {
	return _P1Orders.Contract.STATUS(&_P1Orders.CallOpts, arg0)
}

// GetOrdersStatus is a free data retrieval call binding the contract method 0xaacc263e.
//
// Solidity: function getOrdersStatus(bytes32[] orderHashes) view returns((uint8,uint256)[])
func (_P1Orders *P1OrdersCaller) GetOrdersStatus(opts *bind.CallOpts, orderHashes [][32]byte) ([]P1OrdersOrderQueryOutput, error) {
	var out []interface{}
	err := _P1Orders.contract.Call(opts, &out, "getOrdersStatus", orderHashes)

	if err != nil {
		return *new([]P1OrdersOrderQueryOutput), err
	}

	out0 := *abi.ConvertType(out[0], new([]P1OrdersOrderQueryOutput)).(*[]P1OrdersOrderQueryOutput)

	return out0, err

}

// GetOrdersStatus is a free data retrieval call binding the contract method 0xaacc263e.
//
// Solidity: function getOrdersStatus(bytes32[] orderHashes) view returns((uint8,uint256)[])
func (_P1Orders *P1OrdersSession) GetOrdersStatus(orderHashes [][32]byte) ([]P1OrdersOrderQueryOutput, error) {
	return _P1Orders.Contract.GetOrdersStatus(&_P1Orders.CallOpts, orderHashes)
}

// GetOrdersStatus is a free data retrieval call binding the contract method 0xaacc263e.
//
// Solidity: function getOrdersStatus(bytes32[] orderHashes) view returns((uint8,uint256)[])
func (_P1Orders *P1OrdersCallerSession) GetOrdersStatus(orderHashes [][32]byte) ([]P1OrdersOrderQueryOutput, error) {
	return _P1Orders.Contract.GetOrdersStatus(&_P1Orders.CallOpts, orderHashes)
}

// ApproveOrder is a paid mutator transaction binding the contract method 0x867f1690.
//
// Solidity: function approveOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersTransactor) ApproveOrder(opts *bind.TransactOpts, order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.contract.Transact(opts, "approveOrder", order)
}

// ApproveOrder is a paid mutator transaction binding the contract method 0x867f1690.
//
// Solidity: function approveOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersSession) ApproveOrder(order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.Contract.ApproveOrder(&_P1Orders.TransactOpts, order)
}

// ApproveOrder is a paid mutator transaction binding the contract method 0x867f1690.
//
// Solidity: function approveOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersTransactorSession) ApproveOrder(order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.Contract.ApproveOrder(&_P1Orders.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7946c890.
//
// Solidity: function cancelOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersTransactor) CancelOrder(opts *bind.TransactOpts, order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7946c890.
//
// Solidity: function cancelOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersSession) CancelOrder(order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.Contract.CancelOrder(&_P1Orders.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x7946c890.
//
// Solidity: function cancelOrder((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) returns()
func (_P1Orders *P1OrdersTransactorSession) CancelOrder(order P1OrdersOrder) (*types.Transaction, error) {
	return _P1Orders.Contract.CancelOrder(&_P1Orders.TransactOpts, order)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Orders *P1OrdersTransactor) Trade(opts *bind.TransactOpts, sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Orders.contract.Transact(opts, "trade", sender, maker, taker, price, data, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Orders *P1OrdersSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Orders.Contract.Trade(&_P1Orders.TransactOpts, sender, maker, taker, price, data, arg5)
}

// Trade is a paid mutator transaction binding the contract method 0x970c2ba1.
//
// Solidity: function trade(address sender, address maker, address taker, uint256 price, bytes data, bytes32 ) returns((uint256,uint256,bool,bytes32))
func (_P1Orders *P1OrdersTransactorSession) Trade(sender common.Address, maker common.Address, taker common.Address, price *big.Int, data []byte, arg5 [32]byte) (*types.Transaction, error) {
	return _P1Orders.Contract.Trade(&_P1Orders.TransactOpts, sender, maker, taker, price, data, arg5)
}

// P1OrdersLogOrderApprovedIterator is returned from FilterLogOrderApproved and is used to iterate over the raw logs and unpacked data for LogOrderApproved events raised by the P1Orders contract.
type P1OrdersLogOrderApprovedIterator struct {
	Event *P1OrdersLogOrderApproved // Event containing the contract specifics and raw log

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
func (it *P1OrdersLogOrderApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1OrdersLogOrderApproved)
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
		it.Event = new(P1OrdersLogOrderApproved)
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
func (it *P1OrdersLogOrderApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1OrdersLogOrderApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1OrdersLogOrderApproved represents a LogOrderApproved event raised by the P1Orders contract.
type P1OrdersLogOrderApproved struct {
	Maker     common.Address
	OrderHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogOrderApproved is a free log retrieval operation binding the contract event 0xbd06df5febc1b0cd2e8ba37a6bb524ae77524c4aa2dc5e0f5ac64f5d11a50b1b.
//
// Solidity: event LogOrderApproved(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) FilterLogOrderApproved(opts *bind.FilterOpts, maker []common.Address) (*P1OrdersLogOrderApprovedIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _P1Orders.contract.FilterLogs(opts, "LogOrderApproved", makerRule)
	if err != nil {
		return nil, err
	}
	return &P1OrdersLogOrderApprovedIterator{contract: _P1Orders.contract, event: "LogOrderApproved", logs: logs, sub: sub}, nil
}

// WatchLogOrderApproved is a free log subscription operation binding the contract event 0xbd06df5febc1b0cd2e8ba37a6bb524ae77524c4aa2dc5e0f5ac64f5d11a50b1b.
//
// Solidity: event LogOrderApproved(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) WatchLogOrderApproved(opts *bind.WatchOpts, sink chan<- *P1OrdersLogOrderApproved, maker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _P1Orders.contract.WatchLogs(opts, "LogOrderApproved", makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1OrdersLogOrderApproved)
				if err := _P1Orders.contract.UnpackLog(event, "LogOrderApproved", log); err != nil {
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

// ParseLogOrderApproved is a log parse operation binding the contract event 0xbd06df5febc1b0cd2e8ba37a6bb524ae77524c4aa2dc5e0f5ac64f5d11a50b1b.
//
// Solidity: event LogOrderApproved(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) ParseLogOrderApproved(log types.Log) (*P1OrdersLogOrderApproved, error) {
	event := new(P1OrdersLogOrderApproved)
	if err := _P1Orders.contract.UnpackLog(event, "LogOrderApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1OrdersLogOrderCanceledIterator is returned from FilterLogOrderCanceled and is used to iterate over the raw logs and unpacked data for LogOrderCanceled events raised by the P1Orders contract.
type P1OrdersLogOrderCanceledIterator struct {
	Event *P1OrdersLogOrderCanceled // Event containing the contract specifics and raw log

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
func (it *P1OrdersLogOrderCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1OrdersLogOrderCanceled)
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
		it.Event = new(P1OrdersLogOrderCanceled)
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
func (it *P1OrdersLogOrderCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1OrdersLogOrderCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1OrdersLogOrderCanceled represents a LogOrderCanceled event raised by the P1Orders contract.
type P1OrdersLogOrderCanceled struct {
	Maker     common.Address
	OrderHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogOrderCanceled is a free log retrieval operation binding the contract event 0x4117a4c82505f7102c183e1fb9daa8f8e06d56d6af04479fc417fa8c04902893.
//
// Solidity: event LogOrderCanceled(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) FilterLogOrderCanceled(opts *bind.FilterOpts, maker []common.Address) (*P1OrdersLogOrderCanceledIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _P1Orders.contract.FilterLogs(opts, "LogOrderCanceled", makerRule)
	if err != nil {
		return nil, err
	}
	return &P1OrdersLogOrderCanceledIterator{contract: _P1Orders.contract, event: "LogOrderCanceled", logs: logs, sub: sub}, nil
}

// WatchLogOrderCanceled is a free log subscription operation binding the contract event 0x4117a4c82505f7102c183e1fb9daa8f8e06d56d6af04479fc417fa8c04902893.
//
// Solidity: event LogOrderCanceled(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) WatchLogOrderCanceled(opts *bind.WatchOpts, sink chan<- *P1OrdersLogOrderCanceled, maker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _P1Orders.contract.WatchLogs(opts, "LogOrderCanceled", makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1OrdersLogOrderCanceled)
				if err := _P1Orders.contract.UnpackLog(event, "LogOrderCanceled", log); err != nil {
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

// ParseLogOrderCanceled is a log parse operation binding the contract event 0x4117a4c82505f7102c183e1fb9daa8f8e06d56d6af04479fc417fa8c04902893.
//
// Solidity: event LogOrderCanceled(address indexed maker, bytes32 orderHash)
func (_P1Orders *P1OrdersFilterer) ParseLogOrderCanceled(log types.Log) (*P1OrdersLogOrderCanceled, error) {
	event := new(P1OrdersLogOrderCanceled)
	if err := _P1Orders.contract.UnpackLog(event, "LogOrderCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1OrdersLogOrderFilledIterator is returned from FilterLogOrderFilled and is used to iterate over the raw logs and unpacked data for LogOrderFilled events raised by the P1Orders contract.
type P1OrdersLogOrderFilledIterator struct {
	Event *P1OrdersLogOrderFilled // Event containing the contract specifics and raw log

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
func (it *P1OrdersLogOrderFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1OrdersLogOrderFilled)
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
		it.Event = new(P1OrdersLogOrderFilled)
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
func (it *P1OrdersLogOrderFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1OrdersLogOrderFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1OrdersLogOrderFilled represents a LogOrderFilled event raised by the P1Orders contract.
type P1OrdersLogOrderFilled struct {
	OrderHash    [32]byte
	Flags        [32]byte
	TriggerPrice *big.Int
	Fill         P1OrdersFill
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogOrderFilled is a free log retrieval operation binding the contract event 0x5760b5a80923536b02524ebe3b1f92cc973195ac25559c60564e8db9e02d15ad.
//
// Solidity: event LogOrderFilled(bytes32 orderHash, bytes32 flags, uint256 triggerPrice, (uint256,uint256,uint256,bool) fill)
func (_P1Orders *P1OrdersFilterer) FilterLogOrderFilled(opts *bind.FilterOpts) (*P1OrdersLogOrderFilledIterator, error) {

	logs, sub, err := _P1Orders.contract.FilterLogs(opts, "LogOrderFilled")
	if err != nil {
		return nil, err
	}
	return &P1OrdersLogOrderFilledIterator{contract: _P1Orders.contract, event: "LogOrderFilled", logs: logs, sub: sub}, nil
}

// WatchLogOrderFilled is a free log subscription operation binding the contract event 0x5760b5a80923536b02524ebe3b1f92cc973195ac25559c60564e8db9e02d15ad.
//
// Solidity: event LogOrderFilled(bytes32 orderHash, bytes32 flags, uint256 triggerPrice, (uint256,uint256,uint256,bool) fill)
func (_P1Orders *P1OrdersFilterer) WatchLogOrderFilled(opts *bind.WatchOpts, sink chan<- *P1OrdersLogOrderFilled) (event.Subscription, error) {

	logs, sub, err := _P1Orders.contract.WatchLogs(opts, "LogOrderFilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1OrdersLogOrderFilled)
				if err := _P1Orders.contract.UnpackLog(event, "LogOrderFilled", log); err != nil {
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

// ParseLogOrderFilled is a log parse operation binding the contract event 0x5760b5a80923536b02524ebe3b1f92cc973195ac25559c60564e8db9e02d15ad.
//
// Solidity: event LogOrderFilled(bytes32 orderHash, bytes32 flags, uint256 triggerPrice, (uint256,uint256,uint256,bool) fill)
func (_P1Orders *P1OrdersFilterer) ParseLogOrderFilled(log types.Log) (*P1OrdersLogOrderFilled, error) {
	event := new(P1OrdersLogOrderFilled)
	if err := _P1Orders.contract.UnpackLog(event, "LogOrderFilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
