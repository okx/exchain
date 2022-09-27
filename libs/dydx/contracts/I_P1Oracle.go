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

// IP1OracleMetaData contains all meta data concerning the IP1Oracle contract.
var IP1OracleMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IP1OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use IP1OracleMetaData.ABI instead.
var IP1OracleABI = IP1OracleMetaData.ABI

// IP1Oracle is an auto generated Go binding around an Ethereum contract.
type IP1Oracle struct {
	IP1OracleCaller     // Read-only binding to the contract
	IP1OracleTransactor // Write-only binding to the contract
	IP1OracleFilterer   // Log filterer for contract events
}

// IP1OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type IP1OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IP1OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IP1OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IP1OracleSession struct {
	Contract     *IP1Oracle        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IP1OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IP1OracleCallerSession struct {
	Contract *IP1OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IP1OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IP1OracleTransactorSession struct {
	Contract     *IP1OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IP1OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type IP1OracleRaw struct {
	Contract *IP1Oracle // Generic contract binding to access the raw methods on
}

// IP1OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IP1OracleCallerRaw struct {
	Contract *IP1OracleCaller // Generic read-only contract binding to access the raw methods on
}

// IP1OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IP1OracleTransactorRaw struct {
	Contract *IP1OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIP1Oracle creates a new instance of IP1Oracle, bound to a specific deployed contract.
func NewIP1Oracle(address common.Address, backend bind.ContractBackend) (*IP1Oracle, error) {
	contract, err := bindIP1Oracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IP1Oracle{IP1OracleCaller: IP1OracleCaller{contract: contract}, IP1OracleTransactor: IP1OracleTransactor{contract: contract}, IP1OracleFilterer: IP1OracleFilterer{contract: contract}}, nil
}

// NewIP1OracleCaller creates a new read-only instance of IP1Oracle, bound to a specific deployed contract.
func NewIP1OracleCaller(address common.Address, caller bind.ContractCaller) (*IP1OracleCaller, error) {
	contract, err := bindIP1Oracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IP1OracleCaller{contract: contract}, nil
}

// NewIP1OracleTransactor creates a new write-only instance of IP1Oracle, bound to a specific deployed contract.
func NewIP1OracleTransactor(address common.Address, transactor bind.ContractTransactor) (*IP1OracleTransactor, error) {
	contract, err := bindIP1Oracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IP1OracleTransactor{contract: contract}, nil
}

// NewIP1OracleFilterer creates a new log filterer instance of IP1Oracle, bound to a specific deployed contract.
func NewIP1OracleFilterer(address common.Address, filterer bind.ContractFilterer) (*IP1OracleFilterer, error) {
	contract, err := bindIP1Oracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IP1OracleFilterer{contract: contract}, nil
}

// bindIP1Oracle binds a generic wrapper to an already deployed contract.
func bindIP1Oracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IP1OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Oracle *IP1OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Oracle.Contract.IP1OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Oracle *IP1OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Oracle.Contract.IP1OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Oracle *IP1OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Oracle.Contract.IP1OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Oracle *IP1OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Oracle *IP1OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Oracle *IP1OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Oracle.Contract.contract.Transact(opts, method, params...)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_IP1Oracle *IP1OracleCaller) GetPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IP1Oracle.contract.Call(opts, &out, "getPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_IP1Oracle *IP1OracleSession) GetPrice() (*big.Int, error) {
	return _IP1Oracle.Contract.GetPrice(&_IP1Oracle.CallOpts)
}

// GetPrice is a free data retrieval call binding the contract method 0x98d5fdca.
//
// Solidity: function getPrice() view returns(uint256)
func (_IP1Oracle *IP1OracleCallerSession) GetPrice() (*big.Int, error) {
	return _IP1Oracle.Contract.GetPrice(&_IP1Oracle.CallOpts)
}
