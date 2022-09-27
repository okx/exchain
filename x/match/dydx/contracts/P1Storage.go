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

// P1StorageMetaData contains all meta data concerning the P1Storage contract.
var P1StorageMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1StorageABI is the input ABI used to generate the binding from.
// Deprecated: Use P1StorageMetaData.ABI instead.
var P1StorageABI = P1StorageMetaData.ABI

// P1Storage is an auto generated Go binding around an Ethereum contract.
type P1Storage struct {
	P1StorageCaller     // Read-only binding to the contract
	P1StorageTransactor // Write-only binding to the contract
	P1StorageFilterer   // Log filterer for contract events
}

// P1StorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1StorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1StorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1StorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1StorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1StorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1StorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1StorageSession struct {
	Contract     *P1Storage        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1StorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1StorageCallerSession struct {
	Contract *P1StorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// P1StorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1StorageTransactorSession struct {
	Contract     *P1StorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// P1StorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1StorageRaw struct {
	Contract *P1Storage // Generic contract binding to access the raw methods on
}

// P1StorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1StorageCallerRaw struct {
	Contract *P1StorageCaller // Generic read-only contract binding to access the raw methods on
}

// P1StorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1StorageTransactorRaw struct {
	Contract *P1StorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Storage creates a new instance of P1Storage, bound to a specific deployed contract.
func NewP1Storage(address common.Address, backend bind.ContractBackend) (*P1Storage, error) {
	contract, err := bindP1Storage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Storage{P1StorageCaller: P1StorageCaller{contract: contract}, P1StorageTransactor: P1StorageTransactor{contract: contract}, P1StorageFilterer: P1StorageFilterer{contract: contract}}, nil
}

// NewP1StorageCaller creates a new read-only instance of P1Storage, bound to a specific deployed contract.
func NewP1StorageCaller(address common.Address, caller bind.ContractCaller) (*P1StorageCaller, error) {
	contract, err := bindP1Storage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1StorageCaller{contract: contract}, nil
}

// NewP1StorageTransactor creates a new write-only instance of P1Storage, bound to a specific deployed contract.
func NewP1StorageTransactor(address common.Address, transactor bind.ContractTransactor) (*P1StorageTransactor, error) {
	contract, err := bindP1Storage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1StorageTransactor{contract: contract}, nil
}

// NewP1StorageFilterer creates a new log filterer instance of P1Storage, bound to a specific deployed contract.
func NewP1StorageFilterer(address common.Address, filterer bind.ContractFilterer) (*P1StorageFilterer, error) {
	contract, err := bindP1Storage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1StorageFilterer{contract: contract}, nil
}

// bindP1Storage binds a generic wrapper to an already deployed contract.
func bindP1Storage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1StorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Storage *P1StorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Storage.Contract.P1StorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Storage *P1StorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Storage.Contract.P1StorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Storage *P1StorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Storage.Contract.P1StorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Storage *P1StorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Storage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Storage *P1StorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Storage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Storage *P1StorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Storage.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Storage *P1StorageCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Storage.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Storage *P1StorageSession) GetAdmin() (common.Address, error) {
	return _P1Storage.Contract.GetAdmin(&_P1Storage.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Storage *P1StorageCallerSession) GetAdmin() (common.Address, error) {
	return _P1Storage.Contract.GetAdmin(&_P1Storage.CallOpts)
}
