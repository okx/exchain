// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package placeorder

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
	_ = abi.ConvertType
)

// OrdersOrder is an auto generated low-level Go binding around an user-defined struct.
type OrdersOrder struct {
	Flags        [32]byte
	Amount       *big.Int
	LimitPrice   *big.Int
	TriggerPrice *big.Int
	LimitFee     *big.Int
	Maker        common.Address
	Taker        common.Address
	Expiration   *big.Int
}

// PlaceorderMetaData contains all meta data concerning the Placeorder contract.
var PlaceorderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"_EIP712_DOMAIN_HASH_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"triggerPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structOrders.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"flags\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"triggerPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limitFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structOrders.Order\",\"name\":\"_order\",\"type\":\"tuple\"}],\"name\":\"getOrderMessage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_orderMessage\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"_signature\",\"type\":\"bytes32\"}],\"name\":\"getOrderTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// PlaceorderABI is the input ABI used to generate the binding from.
// Deprecated: Use PlaceorderMetaData.ABI instead.
var PlaceorderABI = PlaceorderMetaData.ABI

// Placeorder is an auto generated Go binding around an Ethereum contract.
type Placeorder struct {
	PlaceorderCaller     // Read-only binding to the contract
	PlaceorderTransactor // Write-only binding to the contract
	PlaceorderFilterer   // Log filterer for contract events
}

// PlaceorderCaller is an auto generated read-only Go binding around an Ethereum contract.
type PlaceorderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlaceorderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PlaceorderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlaceorderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PlaceorderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlaceorderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PlaceorderSession struct {
	Contract     *Placeorder       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PlaceorderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PlaceorderCallerSession struct {
	Contract *PlaceorderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// PlaceorderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PlaceorderTransactorSession struct {
	Contract     *PlaceorderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PlaceorderRaw is an auto generated low-level Go binding around an Ethereum contract.
type PlaceorderRaw struct {
	Contract *Placeorder // Generic contract binding to access the raw methods on
}

// PlaceorderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PlaceorderCallerRaw struct {
	Contract *PlaceorderCaller // Generic read-only contract binding to access the raw methods on
}

// PlaceorderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PlaceorderTransactorRaw struct {
	Contract *PlaceorderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPlaceorder creates a new instance of Placeorder, bound to a specific deployed contract.
func NewPlaceorder(address common.Address, backend bind.ContractBackend) (*Placeorder, error) {
	contract, err := bindPlaceorder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Placeorder{PlaceorderCaller: PlaceorderCaller{contract: contract}, PlaceorderTransactor: PlaceorderTransactor{contract: contract}, PlaceorderFilterer: PlaceorderFilterer{contract: contract}}, nil
}

// NewPlaceorderCaller creates a new read-only instance of Placeorder, bound to a specific deployed contract.
func NewPlaceorderCaller(address common.Address, caller bind.ContractCaller) (*PlaceorderCaller, error) {
	contract, err := bindPlaceorder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PlaceorderCaller{contract: contract}, nil
}

// NewPlaceorderTransactor creates a new write-only instance of Placeorder, bound to a specific deployed contract.
func NewPlaceorderTransactor(address common.Address, transactor bind.ContractTransactor) (*PlaceorderTransactor, error) {
	contract, err := bindPlaceorder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PlaceorderTransactor{contract: contract}, nil
}

// NewPlaceorderFilterer creates a new log filterer instance of Placeorder, bound to a specific deployed contract.
func NewPlaceorderFilterer(address common.Address, filterer bind.ContractFilterer) (*PlaceorderFilterer, error) {
	contract, err := bindPlaceorder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PlaceorderFilterer{contract: contract}, nil
}

// bindPlaceorder binds a generic wrapper to an already deployed contract.
func bindPlaceorder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PlaceorderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Placeorder *PlaceorderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Placeorder.Contract.PlaceorderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Placeorder *PlaceorderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Placeorder.Contract.PlaceorderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Placeorder *PlaceorderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Placeorder.Contract.PlaceorderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Placeorder *PlaceorderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Placeorder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Placeorder *PlaceorderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Placeorder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Placeorder *PlaceorderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Placeorder.Contract.contract.Transact(opts, method, params...)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_Placeorder *PlaceorderCaller) EIP712DOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Placeorder.contract.Call(opts, &out, "_EIP712_DOMAIN_HASH_")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_Placeorder *PlaceorderSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _Placeorder.Contract.EIP712DOMAINHASH(&_Placeorder.CallOpts)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_Placeorder *PlaceorderCallerSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _Placeorder.Contract.EIP712DOMAINHASH(&_Placeorder.CallOpts)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xd4f2b529.
//
// Solidity: function getOrderHash((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) view returns(bytes32)
func (_Placeorder *PlaceorderCaller) GetOrderHash(opts *bind.CallOpts, order OrdersOrder) ([32]byte, error) {
	var out []interface{}
	err := _Placeorder.contract.Call(opts, &out, "getOrderHash", order)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetOrderHash is a free data retrieval call binding the contract method 0xd4f2b529.
//
// Solidity: function getOrderHash((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) view returns(bytes32)
func (_Placeorder *PlaceorderSession) GetOrderHash(order OrdersOrder) ([32]byte, error) {
	return _Placeorder.Contract.GetOrderHash(&_Placeorder.CallOpts, order)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xd4f2b529.
//
// Solidity: function getOrderHash((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) order) view returns(bytes32)
func (_Placeorder *PlaceorderCallerSession) GetOrderHash(order OrdersOrder) ([32]byte, error) {
	return _Placeorder.Contract.GetOrderHash(&_Placeorder.CallOpts, order)
}

// GetOrderMessage is a free data retrieval call binding the contract method 0x2c6bc25d.
//
// Solidity: function getOrderMessage((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) _order) pure returns(bytes)
func (_Placeorder *PlaceorderCaller) GetOrderMessage(opts *bind.CallOpts, _order OrdersOrder) ([]byte, error) {
	var out []interface{}
	err := _Placeorder.contract.Call(opts, &out, "getOrderMessage", _order)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetOrderMessage is a free data retrieval call binding the contract method 0x2c6bc25d.
//
// Solidity: function getOrderMessage((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) _order) pure returns(bytes)
func (_Placeorder *PlaceorderSession) GetOrderMessage(_order OrdersOrder) ([]byte, error) {
	return _Placeorder.Contract.GetOrderMessage(&_Placeorder.CallOpts, _order)
}

// GetOrderMessage is a free data retrieval call binding the contract method 0x2c6bc25d.
//
// Solidity: function getOrderMessage((bytes32,uint256,uint256,uint256,uint256,address,address,uint256) _order) pure returns(bytes)
func (_Placeorder *PlaceorderCallerSession) GetOrderMessage(_order OrdersOrder) ([]byte, error) {
	return _Placeorder.Contract.GetOrderMessage(&_Placeorder.CallOpts, _order)
}

// GetOrderTransaction is a free data retrieval call binding the contract method 0xdd53ef8c.
//
// Solidity: function getOrderTransaction(bytes _orderMessage, bytes32 _signature) pure returns(bytes)
func (_Placeorder *PlaceorderCaller) GetOrderTransaction(opts *bind.CallOpts, _orderMessage []byte, _signature [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Placeorder.contract.Call(opts, &out, "getOrderTransaction", _orderMessage, _signature)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetOrderTransaction is a free data retrieval call binding the contract method 0xdd53ef8c.
//
// Solidity: function getOrderTransaction(bytes _orderMessage, bytes32 _signature) pure returns(bytes)
func (_Placeorder *PlaceorderSession) GetOrderTransaction(_orderMessage []byte, _signature [32]byte) ([]byte, error) {
	return _Placeorder.Contract.GetOrderTransaction(&_Placeorder.CallOpts, _orderMessage, _signature)
}

// GetOrderTransaction is a free data retrieval call binding the contract method 0xdd53ef8c.
//
// Solidity: function getOrderTransaction(bytes _orderMessage, bytes32 _signature) pure returns(bytes)
func (_Placeorder *PlaceorderCallerSession) GetOrderTransaction(_orderMessage []byte, _signature [32]byte) ([]byte, error) {
	return _Placeorder.Contract.GetOrderTransaction(&_Placeorder.CallOpts, _orderMessage, _signature)
}