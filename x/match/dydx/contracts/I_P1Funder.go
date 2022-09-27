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

// IP1FunderMetaData contains all meta data concerning the IP1Funder contract.
var IP1FunderMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeDelta\",\"type\":\"uint256\"}],\"name\":\"getFunding\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IP1FunderABI is the input ABI used to generate the binding from.
// Deprecated: Use IP1FunderMetaData.ABI instead.
var IP1FunderABI = IP1FunderMetaData.ABI

// IP1Funder is an auto generated Go binding around an Ethereum contract.
type IP1Funder struct {
	IP1FunderCaller     // Read-only binding to the contract
	IP1FunderTransactor // Write-only binding to the contract
	IP1FunderFilterer   // Log filterer for contract events
}

// IP1FunderCaller is an auto generated read-only Go binding around an Ethereum contract.
type IP1FunderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1FunderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IP1FunderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1FunderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IP1FunderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IP1FunderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IP1FunderSession struct {
	Contract     *IP1Funder        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IP1FunderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IP1FunderCallerSession struct {
	Contract *IP1FunderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IP1FunderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IP1FunderTransactorSession struct {
	Contract     *IP1FunderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IP1FunderRaw is an auto generated low-level Go binding around an Ethereum contract.
type IP1FunderRaw struct {
	Contract *IP1Funder // Generic contract binding to access the raw methods on
}

// IP1FunderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IP1FunderCallerRaw struct {
	Contract *IP1FunderCaller // Generic read-only contract binding to access the raw methods on
}

// IP1FunderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IP1FunderTransactorRaw struct {
	Contract *IP1FunderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIP1Funder creates a new instance of IP1Funder, bound to a specific deployed contract.
func NewIP1Funder(address common.Address, backend bind.ContractBackend) (*IP1Funder, error) {
	contract, err := bindIP1Funder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IP1Funder{IP1FunderCaller: IP1FunderCaller{contract: contract}, IP1FunderTransactor: IP1FunderTransactor{contract: contract}, IP1FunderFilterer: IP1FunderFilterer{contract: contract}}, nil
}

// NewIP1FunderCaller creates a new read-only instance of IP1Funder, bound to a specific deployed contract.
func NewIP1FunderCaller(address common.Address, caller bind.ContractCaller) (*IP1FunderCaller, error) {
	contract, err := bindIP1Funder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IP1FunderCaller{contract: contract}, nil
}

// NewIP1FunderTransactor creates a new write-only instance of IP1Funder, bound to a specific deployed contract.
func NewIP1FunderTransactor(address common.Address, transactor bind.ContractTransactor) (*IP1FunderTransactor, error) {
	contract, err := bindIP1Funder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IP1FunderTransactor{contract: contract}, nil
}

// NewIP1FunderFilterer creates a new log filterer instance of IP1Funder, bound to a specific deployed contract.
func NewIP1FunderFilterer(address common.Address, filterer bind.ContractFilterer) (*IP1FunderFilterer, error) {
	contract, err := bindIP1Funder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IP1FunderFilterer{contract: contract}, nil
}

// bindIP1Funder binds a generic wrapper to an already deployed contract.
func bindIP1Funder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IP1FunderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Funder *IP1FunderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Funder.Contract.IP1FunderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Funder *IP1FunderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Funder.Contract.IP1FunderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Funder *IP1FunderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Funder.Contract.IP1FunderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IP1Funder *IP1FunderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IP1Funder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IP1Funder *IP1FunderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IP1Funder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IP1Funder *IP1FunderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IP1Funder.Contract.contract.Transact(opts, method, params...)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_IP1Funder *IP1FunderCaller) GetFunding(opts *bind.CallOpts, timeDelta *big.Int) (bool, *big.Int, error) {
	var out []interface{}
	err := _IP1Funder.contract.Call(opts, &out, "getFunding", timeDelta)

	if err != nil {
		return *new(bool), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_IP1Funder *IP1FunderSession) GetFunding(timeDelta *big.Int) (bool, *big.Int, error) {
	return _IP1Funder.Contract.GetFunding(&_IP1Funder.CallOpts, timeDelta)
}

// GetFunding is a free data retrieval call binding the contract method 0xebed4bd4.
//
// Solidity: function getFunding(uint256 timeDelta) view returns(bool, uint256)
func (_IP1Funder *IP1FunderCallerSession) GetFunding(timeDelta *big.Int) (bool, *big.Int, error) {
	return _IP1Funder.Contract.GetFunding(&_IP1Funder.CallOpts, timeDelta)
}
