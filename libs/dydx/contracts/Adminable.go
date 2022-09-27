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

// AdminableMetaData contains all meta data concerning the Adminable contract.
var AdminableMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AdminableABI is the input ABI used to generate the binding from.
// Deprecated: Use AdminableMetaData.ABI instead.
var AdminableABI = AdminableMetaData.ABI

// Adminable is an auto generated Go binding around an Ethereum contract.
type Adminable struct {
	AdminableCaller     // Read-only binding to the contract
	AdminableTransactor // Write-only binding to the contract
	AdminableFilterer   // Log filterer for contract events
}

// AdminableCaller is an auto generated read-only Go binding around an Ethereum contract.
type AdminableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AdminableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AdminableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdminableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AdminableSession struct {
	Contract     *Adminable        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AdminableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AdminableCallerSession struct {
	Contract *AdminableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// AdminableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AdminableTransactorSession struct {
	Contract     *AdminableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AdminableRaw is an auto generated low-level Go binding around an Ethereum contract.
type AdminableRaw struct {
	Contract *Adminable // Generic contract binding to access the raw methods on
}

// AdminableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AdminableCallerRaw struct {
	Contract *AdminableCaller // Generic read-only contract binding to access the raw methods on
}

// AdminableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AdminableTransactorRaw struct {
	Contract *AdminableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdminable creates a new instance of Adminable, bound to a specific deployed contract.
func NewAdminable(address common.Address, backend bind.ContractBackend) (*Adminable, error) {
	contract, err := bindAdminable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Adminable{AdminableCaller: AdminableCaller{contract: contract}, AdminableTransactor: AdminableTransactor{contract: contract}, AdminableFilterer: AdminableFilterer{contract: contract}}, nil
}

// NewAdminableCaller creates a new read-only instance of Adminable, bound to a specific deployed contract.
func NewAdminableCaller(address common.Address, caller bind.ContractCaller) (*AdminableCaller, error) {
	contract, err := bindAdminable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdminableCaller{contract: contract}, nil
}

// NewAdminableTransactor creates a new write-only instance of Adminable, bound to a specific deployed contract.
func NewAdminableTransactor(address common.Address, transactor bind.ContractTransactor) (*AdminableTransactor, error) {
	contract, err := bindAdminable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdminableTransactor{contract: contract}, nil
}

// NewAdminableFilterer creates a new log filterer instance of Adminable, bound to a specific deployed contract.
func NewAdminableFilterer(address common.Address, filterer bind.ContractFilterer) (*AdminableFilterer, error) {
	contract, err := bindAdminable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdminableFilterer{contract: contract}, nil
}

// bindAdminable binds a generic wrapper to an already deployed contract.
func bindAdminable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AdminableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Adminable *AdminableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Adminable.Contract.AdminableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Adminable *AdminableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Adminable.Contract.AdminableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Adminable *AdminableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Adminable.Contract.AdminableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Adminable *AdminableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Adminable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Adminable *AdminableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Adminable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Adminable *AdminableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Adminable.Contract.contract.Transact(opts, method, params...)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_Adminable *AdminableCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Adminable.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_Adminable *AdminableSession) GetAdmin() (common.Address, error) {
	return _Adminable.Contract.GetAdmin(&_Adminable.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_Adminable *AdminableCallerSession) GetAdmin() (common.Address, error) {
	return _Adminable.Contract.GetAdmin(&_Adminable.CallOpts)
}
