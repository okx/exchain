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

// P1MarginMetaData contains all meta data concerning the P1Margin contract.
var P1MarginMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdrawFinalSettlement\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountBalance\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFinalSettlementEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFunderContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGlobalIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOracleContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOraclePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"hasAccountPermissions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1MarginABI is the input ABI used to generate the binding from.
// Deprecated: Use P1MarginMetaData.ABI instead.
var P1MarginABI = P1MarginMetaData.ABI

// P1Margin is an auto generated Go binding around an Ethereum contract.
type P1Margin struct {
	P1MarginCaller     // Read-only binding to the contract
	P1MarginTransactor // Write-only binding to the contract
	P1MarginFilterer   // Log filterer for contract events
}

// P1MarginCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1MarginCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MarginTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1MarginTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MarginFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1MarginFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1MarginSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1MarginSession struct {
	Contract     *P1Margin         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1MarginCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1MarginCallerSession struct {
	Contract *P1MarginCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// P1MarginTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1MarginTransactorSession struct {
	Contract     *P1MarginTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// P1MarginRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1MarginRaw struct {
	Contract *P1Margin // Generic contract binding to access the raw methods on
}

// P1MarginCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1MarginCallerRaw struct {
	Contract *P1MarginCaller // Generic read-only contract binding to access the raw methods on
}

// P1MarginTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1MarginTransactorRaw struct {
	Contract *P1MarginTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Margin creates a new instance of P1Margin, bound to a specific deployed contract.
func NewP1Margin(address common.Address, backend bind.ContractBackend) (*P1Margin, error) {
	contract, err := bindP1Margin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Margin{P1MarginCaller: P1MarginCaller{contract: contract}, P1MarginTransactor: P1MarginTransactor{contract: contract}, P1MarginFilterer: P1MarginFilterer{contract: contract}}, nil
}

// NewP1MarginCaller creates a new read-only instance of P1Margin, bound to a specific deployed contract.
func NewP1MarginCaller(address common.Address, caller bind.ContractCaller) (*P1MarginCaller, error) {
	contract, err := bindP1Margin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1MarginCaller{contract: contract}, nil
}

// NewP1MarginTransactor creates a new write-only instance of P1Margin, bound to a specific deployed contract.
func NewP1MarginTransactor(address common.Address, transactor bind.ContractTransactor) (*P1MarginTransactor, error) {
	contract, err := bindP1Margin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1MarginTransactor{contract: contract}, nil
}

// NewP1MarginFilterer creates a new log filterer instance of P1Margin, bound to a specific deployed contract.
func NewP1MarginFilterer(address common.Address, filterer bind.ContractFilterer) (*P1MarginFilterer, error) {
	contract, err := bindP1Margin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1MarginFilterer{contract: contract}, nil
}

// bindP1Margin binds a generic wrapper to an already deployed contract.
func bindP1Margin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1MarginABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Margin *P1MarginRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Margin.Contract.P1MarginCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Margin *P1MarginRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Margin.Contract.P1MarginTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Margin *P1MarginRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Margin.Contract.P1MarginTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Margin *P1MarginCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Margin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Margin *P1MarginTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Margin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Margin *P1MarginTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Margin.Contract.contract.Transact(opts, method, params...)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Margin *P1MarginCaller) GetAccountBalance(opts *bind.CallOpts, account common.Address) (P1TypesBalance, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getAccountBalance", account)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Margin *P1MarginSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _P1Margin.Contract.GetAccountBalance(&_P1Margin.CallOpts, account)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Margin *P1MarginCallerSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _P1Margin.Contract.GetAccountBalance(&_P1Margin.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginCaller) GetAccountIndex(opts *bind.CallOpts, account common.Address) (P1TypesIndex, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getAccountIndex", account)

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _P1Margin.Contract.GetAccountIndex(&_P1Margin.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginCallerSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _P1Margin.Contract.GetAccountIndex(&_P1Margin.CallOpts, account)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Margin *P1MarginCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Margin *P1MarginSession) GetAdmin() (common.Address, error) {
	return _P1Margin.Contract.GetAdmin(&_P1Margin.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Margin *P1MarginCallerSession) GetAdmin() (common.Address, error) {
	return _P1Margin.Contract.GetAdmin(&_P1Margin.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Margin *P1MarginCaller) GetFinalSettlementEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getFinalSettlementEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Margin *P1MarginSession) GetFinalSettlementEnabled() (bool, error) {
	return _P1Margin.Contract.GetFinalSettlementEnabled(&_P1Margin.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Margin *P1MarginCallerSession) GetFinalSettlementEnabled() (bool, error) {
	return _P1Margin.Contract.GetFinalSettlementEnabled(&_P1Margin.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Margin *P1MarginCaller) GetFunderContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getFunderContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Margin *P1MarginSession) GetFunderContract() (common.Address, error) {
	return _P1Margin.Contract.GetFunderContract(&_P1Margin.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Margin *P1MarginCallerSession) GetFunderContract() (common.Address, error) {
	return _P1Margin.Contract.GetFunderContract(&_P1Margin.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginCaller) GetGlobalIndex(opts *bind.CallOpts) (P1TypesIndex, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getGlobalIndex")

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _P1Margin.Contract.GetGlobalIndex(&_P1Margin.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Margin *P1MarginCallerSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _P1Margin.Contract.GetGlobalIndex(&_P1Margin.CallOpts)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Margin *P1MarginCaller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Margin *P1MarginSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _P1Margin.Contract.GetIsGlobalOperator(&_P1Margin.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Margin *P1MarginCallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _P1Margin.Contract.GetIsGlobalOperator(&_P1Margin.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginCaller) GetIsLocalOperator(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getIsLocalOperator", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _P1Margin.Contract.GetIsLocalOperator(&_P1Margin.CallOpts, account, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginCallerSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _P1Margin.Contract.GetIsLocalOperator(&_P1Margin.CallOpts, account, operator)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Margin *P1MarginCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Margin *P1MarginSession) GetMinCollateral() (*big.Int, error) {
	return _P1Margin.Contract.GetMinCollateral(&_P1Margin.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Margin *P1MarginCallerSession) GetMinCollateral() (*big.Int, error) {
	return _P1Margin.Contract.GetMinCollateral(&_P1Margin.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Margin *P1MarginCaller) GetOracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getOracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Margin *P1MarginSession) GetOracleContract() (common.Address, error) {
	return _P1Margin.Contract.GetOracleContract(&_P1Margin.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Margin *P1MarginCallerSession) GetOracleContract() (common.Address, error) {
	return _P1Margin.Contract.GetOracleContract(&_P1Margin.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Margin *P1MarginCaller) GetOraclePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getOraclePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Margin *P1MarginSession) GetOraclePrice() (*big.Int, error) {
	return _P1Margin.Contract.GetOraclePrice(&_P1Margin.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Margin *P1MarginCallerSession) GetOraclePrice() (*big.Int, error) {
	return _P1Margin.Contract.GetOraclePrice(&_P1Margin.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Margin *P1MarginCaller) GetTokenContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "getTokenContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Margin *P1MarginSession) GetTokenContract() (common.Address, error) {
	return _P1Margin.Contract.GetTokenContract(&_P1Margin.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Margin *P1MarginCallerSession) GetTokenContract() (common.Address, error) {
	return _P1Margin.Contract.GetTokenContract(&_P1Margin.CallOpts)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginCaller) HasAccountPermissions(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Margin.contract.Call(opts, &out, "hasAccountPermissions", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _P1Margin.Contract.HasAccountPermissions(&_P1Margin.CallOpts, account, operator)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Margin *P1MarginCallerSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _P1Margin.Contract.HasAccountPermissions(&_P1Margin.CallOpts, account, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_P1Margin *P1MarginTransactor) Deposit(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.contract.Transact(opts, "deposit", account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_P1Margin *P1MarginSession) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.Contract.Deposit(&_P1Margin.TransactOpts, account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_P1Margin *P1MarginTransactorSession) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.Contract.Deposit(&_P1Margin.TransactOpts, account, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_P1Margin *P1MarginTransactor) Withdraw(opts *bind.TransactOpts, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.contract.Transact(opts, "withdraw", account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_P1Margin *P1MarginSession) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.Contract.Withdraw(&_P1Margin.TransactOpts, account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_P1Margin *P1MarginTransactorSession) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _P1Margin.Contract.Withdraw(&_P1Margin.TransactOpts, account, destination, amount)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Margin *P1MarginTransactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Margin.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Margin *P1MarginSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Margin.Contract.WithdrawFinalSettlement(&_P1Margin.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_P1Margin *P1MarginTransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _P1Margin.Contract.WithdrawFinalSettlement(&_P1Margin.TransactOpts)
}

// P1MarginLogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the P1Margin contract.
type P1MarginLogAccountSettledIterator struct {
	Event *P1MarginLogAccountSettled // Event containing the contract specifics and raw log

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
func (it *P1MarginLogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MarginLogAccountSettled)
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
		it.Event = new(P1MarginLogAccountSettled)
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
func (it *P1MarginLogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MarginLogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MarginLogAccountSettled represents a LogAccountSettled event raised by the P1Margin contract.
type P1MarginLogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*P1MarginLogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1MarginLogAccountSettledIterator{contract: _P1Margin.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *P1MarginLogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MarginLogAccountSettled)
				if err := _P1Margin.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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
func (_P1Margin *P1MarginFilterer) ParseLogAccountSettled(log types.Log) (*P1MarginLogAccountSettled, error) {
	event := new(P1MarginLogAccountSettled)
	if err := _P1Margin.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MarginLogDepositIterator is returned from FilterLogDeposit and is used to iterate over the raw logs and unpacked data for LogDeposit events raised by the P1Margin contract.
type P1MarginLogDepositIterator struct {
	Event *P1MarginLogDeposit // Event containing the contract specifics and raw log

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
func (it *P1MarginLogDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MarginLogDeposit)
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
		it.Event = new(P1MarginLogDeposit)
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
func (it *P1MarginLogDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MarginLogDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MarginLogDeposit represents a LogDeposit event raised by the P1Margin contract.
type P1MarginLogDeposit struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogDeposit is a free log retrieval operation binding the contract event 0x40a9cb3a9707d3a68091d8ef7ffd4158d01d0b2ad92b1e489abe8312dd543023.
//
// Solidity: event LogDeposit(address indexed account, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) FilterLogDeposit(opts *bind.FilterOpts, account []common.Address) (*P1MarginLogDepositIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.FilterLogs(opts, "LogDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1MarginLogDepositIterator{contract: _P1Margin.contract, event: "LogDeposit", logs: logs, sub: sub}, nil
}

// WatchLogDeposit is a free log subscription operation binding the contract event 0x40a9cb3a9707d3a68091d8ef7ffd4158d01d0b2ad92b1e489abe8312dd543023.
//
// Solidity: event LogDeposit(address indexed account, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) WatchLogDeposit(opts *bind.WatchOpts, sink chan<- *P1MarginLogDeposit, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.WatchLogs(opts, "LogDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MarginLogDeposit)
				if err := _P1Margin.contract.UnpackLog(event, "LogDeposit", log); err != nil {
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

// ParseLogDeposit is a log parse operation binding the contract event 0x40a9cb3a9707d3a68091d8ef7ffd4158d01d0b2ad92b1e489abe8312dd543023.
//
// Solidity: event LogDeposit(address indexed account, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) ParseLogDeposit(log types.Log) (*P1MarginLogDeposit, error) {
	event := new(P1MarginLogDeposit)
	if err := _P1Margin.contract.UnpackLog(event, "LogDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MarginLogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the P1Margin contract.
type P1MarginLogIndexIterator struct {
	Event *P1MarginLogIndex // Event containing the contract specifics and raw log

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
func (it *P1MarginLogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MarginLogIndex)
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
		it.Event = new(P1MarginLogIndex)
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
func (it *P1MarginLogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MarginLogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MarginLogIndex represents a LogIndex event raised by the P1Margin contract.
type P1MarginLogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Margin *P1MarginFilterer) FilterLogIndex(opts *bind.FilterOpts) (*P1MarginLogIndexIterator, error) {

	logs, sub, err := _P1Margin.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &P1MarginLogIndexIterator{contract: _P1Margin.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_P1Margin *P1MarginFilterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *P1MarginLogIndex) (event.Subscription, error) {

	logs, sub, err := _P1Margin.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MarginLogIndex)
				if err := _P1Margin.contract.UnpackLog(event, "LogIndex", log); err != nil {
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
func (_P1Margin *P1MarginFilterer) ParseLogIndex(log types.Log) (*P1MarginLogIndex, error) {
	event := new(P1MarginLogIndex)
	if err := _P1Margin.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MarginLogWithdrawIterator is returned from FilterLogWithdraw and is used to iterate over the raw logs and unpacked data for LogWithdraw events raised by the P1Margin contract.
type P1MarginLogWithdrawIterator struct {
	Event *P1MarginLogWithdraw // Event containing the contract specifics and raw log

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
func (it *P1MarginLogWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MarginLogWithdraw)
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
		it.Event = new(P1MarginLogWithdraw)
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
func (it *P1MarginLogWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MarginLogWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MarginLogWithdraw represents a LogWithdraw event raised by the P1Margin contract.
type P1MarginLogWithdraw struct {
	Account     common.Address
	Destination common.Address
	Amount      *big.Int
	Balance     [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogWithdraw is a free log retrieval operation binding the contract event 0x74348e8cb927b5536fe550310d0cdf05914498fcb04ad61b99c29e3899b0bce9.
//
// Solidity: event LogWithdraw(address indexed account, address destination, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) FilterLogWithdraw(opts *bind.FilterOpts, account []common.Address) (*P1MarginLogWithdrawIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.FilterLogs(opts, "LogWithdraw", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1MarginLogWithdrawIterator{contract: _P1Margin.contract, event: "LogWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogWithdraw is a free log subscription operation binding the contract event 0x74348e8cb927b5536fe550310d0cdf05914498fcb04ad61b99c29e3899b0bce9.
//
// Solidity: event LogWithdraw(address indexed account, address destination, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) WatchLogWithdraw(opts *bind.WatchOpts, sink chan<- *P1MarginLogWithdraw, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.WatchLogs(opts, "LogWithdraw", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MarginLogWithdraw)
				if err := _P1Margin.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
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

// ParseLogWithdraw is a log parse operation binding the contract event 0x74348e8cb927b5536fe550310d0cdf05914498fcb04ad61b99c29e3899b0bce9.
//
// Solidity: event LogWithdraw(address indexed account, address destination, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) ParseLogWithdraw(log types.Log) (*P1MarginLogWithdraw, error) {
	event := new(P1MarginLogWithdraw)
	if err := _P1Margin.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1MarginLogWithdrawFinalSettlementIterator is returned from FilterLogWithdrawFinalSettlement and is used to iterate over the raw logs and unpacked data for LogWithdrawFinalSettlement events raised by the P1Margin contract.
type P1MarginLogWithdrawFinalSettlementIterator struct {
	Event *P1MarginLogWithdrawFinalSettlement // Event containing the contract specifics and raw log

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
func (it *P1MarginLogWithdrawFinalSettlementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1MarginLogWithdrawFinalSettlement)
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
		it.Event = new(P1MarginLogWithdrawFinalSettlement)
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
func (it *P1MarginLogWithdrawFinalSettlementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1MarginLogWithdrawFinalSettlementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1MarginLogWithdrawFinalSettlement represents a LogWithdrawFinalSettlement event raised by the P1Margin contract.
type P1MarginLogWithdrawFinalSettlement struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFinalSettlement is a free log retrieval operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) FilterLogWithdrawFinalSettlement(opts *bind.FilterOpts, account []common.Address) (*P1MarginLogWithdrawFinalSettlementIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.FilterLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1MarginLogWithdrawFinalSettlementIterator{contract: _P1Margin.contract, event: "LogWithdrawFinalSettlement", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFinalSettlement is a free log subscription operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_P1Margin *P1MarginFilterer) WatchLogWithdrawFinalSettlement(opts *bind.WatchOpts, sink chan<- *P1MarginLogWithdrawFinalSettlement, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1Margin.contract.WatchLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1MarginLogWithdrawFinalSettlement)
				if err := _P1Margin.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
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
func (_P1Margin *P1MarginFilterer) ParseLogWithdrawFinalSettlement(log types.Log) (*P1MarginLogWithdrawFinalSettlement, error) {
	event := new(P1MarginLogWithdrawFinalSettlement)
	if err := _P1Margin.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
