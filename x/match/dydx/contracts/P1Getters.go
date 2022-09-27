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

// P1GettersMetaData contains all meta data concerning the P1Getters contract.
var P1GettersMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountBalance\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOracleContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFunderContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGlobalIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFinalSettlementEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOraclePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"hasAccountPermissions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// P1GettersABI is the input ABI used to generate the binding from.
// Deprecated: Use P1GettersMetaData.ABI instead.
var P1GettersABI = P1GettersMetaData.ABI

// P1Getters is an auto generated Go binding around an Ethereum contract.
type P1Getters struct {
	P1GettersCaller     // Read-only binding to the contract
	P1GettersTransactor // Write-only binding to the contract
	P1GettersFilterer   // Log filterer for contract events
}

// P1GettersCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1GettersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1GettersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1GettersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1GettersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1GettersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1GettersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1GettersSession struct {
	Contract     *P1Getters        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// P1GettersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1GettersCallerSession struct {
	Contract *P1GettersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// P1GettersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1GettersTransactorSession struct {
	Contract     *P1GettersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// P1GettersRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1GettersRaw struct {
	Contract *P1Getters // Generic contract binding to access the raw methods on
}

// P1GettersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1GettersCallerRaw struct {
	Contract *P1GettersCaller // Generic read-only contract binding to access the raw methods on
}

// P1GettersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1GettersTransactorRaw struct {
	Contract *P1GettersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1Getters creates a new instance of P1Getters, bound to a specific deployed contract.
func NewP1Getters(address common.Address, backend bind.ContractBackend) (*P1Getters, error) {
	contract, err := bindP1Getters(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1Getters{P1GettersCaller: P1GettersCaller{contract: contract}, P1GettersTransactor: P1GettersTransactor{contract: contract}, P1GettersFilterer: P1GettersFilterer{contract: contract}}, nil
}

// NewP1GettersCaller creates a new read-only instance of P1Getters, bound to a specific deployed contract.
func NewP1GettersCaller(address common.Address, caller bind.ContractCaller) (*P1GettersCaller, error) {
	contract, err := bindP1Getters(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1GettersCaller{contract: contract}, nil
}

// NewP1GettersTransactor creates a new write-only instance of P1Getters, bound to a specific deployed contract.
func NewP1GettersTransactor(address common.Address, transactor bind.ContractTransactor) (*P1GettersTransactor, error) {
	contract, err := bindP1Getters(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1GettersTransactor{contract: contract}, nil
}

// NewP1GettersFilterer creates a new log filterer instance of P1Getters, bound to a specific deployed contract.
func NewP1GettersFilterer(address common.Address, filterer bind.ContractFilterer) (*P1GettersFilterer, error) {
	contract, err := bindP1Getters(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1GettersFilterer{contract: contract}, nil
}

// bindP1Getters binds a generic wrapper to an already deployed contract.
func bindP1Getters(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1GettersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Getters *P1GettersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Getters.Contract.P1GettersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Getters *P1GettersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Getters.Contract.P1GettersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Getters *P1GettersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Getters.Contract.P1GettersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1Getters *P1GettersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1Getters.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1Getters *P1GettersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1Getters.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1Getters *P1GettersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1Getters.Contract.contract.Transact(opts, method, params...)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Getters *P1GettersCaller) GetAccountBalance(opts *bind.CallOpts, account common.Address) (P1TypesBalance, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getAccountBalance", account)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Getters *P1GettersSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _P1Getters.Contract.GetAccountBalance(&_P1Getters.CallOpts, account)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_P1Getters *P1GettersCallerSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _P1Getters.Contract.GetAccountBalance(&_P1Getters.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersCaller) GetAccountIndex(opts *bind.CallOpts, account common.Address) (P1TypesIndex, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getAccountIndex", account)

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _P1Getters.Contract.GetAccountIndex(&_P1Getters.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersCallerSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _P1Getters.Contract.GetAccountIndex(&_P1Getters.CallOpts, account)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Getters *P1GettersCaller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Getters *P1GettersSession) GetAdmin() (common.Address, error) {
	return _P1Getters.Contract.GetAdmin(&_P1Getters.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_P1Getters *P1GettersCallerSession) GetAdmin() (common.Address, error) {
	return _P1Getters.Contract.GetAdmin(&_P1Getters.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Getters *P1GettersCaller) GetFinalSettlementEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getFinalSettlementEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Getters *P1GettersSession) GetFinalSettlementEnabled() (bool, error) {
	return _P1Getters.Contract.GetFinalSettlementEnabled(&_P1Getters.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_P1Getters *P1GettersCallerSession) GetFinalSettlementEnabled() (bool, error) {
	return _P1Getters.Contract.GetFinalSettlementEnabled(&_P1Getters.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Getters *P1GettersCaller) GetFunderContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getFunderContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Getters *P1GettersSession) GetFunderContract() (common.Address, error) {
	return _P1Getters.Contract.GetFunderContract(&_P1Getters.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_P1Getters *P1GettersCallerSession) GetFunderContract() (common.Address, error) {
	return _P1Getters.Contract.GetFunderContract(&_P1Getters.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersCaller) GetGlobalIndex(opts *bind.CallOpts) (P1TypesIndex, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getGlobalIndex")

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _P1Getters.Contract.GetGlobalIndex(&_P1Getters.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_P1Getters *P1GettersCallerSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _P1Getters.Contract.GetGlobalIndex(&_P1Getters.CallOpts)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Getters *P1GettersCaller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Getters *P1GettersSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _P1Getters.Contract.GetIsGlobalOperator(&_P1Getters.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_P1Getters *P1GettersCallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _P1Getters.Contract.GetIsGlobalOperator(&_P1Getters.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersCaller) GetIsLocalOperator(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getIsLocalOperator", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _P1Getters.Contract.GetIsLocalOperator(&_P1Getters.CallOpts, account, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersCallerSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _P1Getters.Contract.GetIsLocalOperator(&_P1Getters.CallOpts, account, operator)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Getters *P1GettersCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Getters *P1GettersSession) GetMinCollateral() (*big.Int, error) {
	return _P1Getters.Contract.GetMinCollateral(&_P1Getters.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_P1Getters *P1GettersCallerSession) GetMinCollateral() (*big.Int, error) {
	return _P1Getters.Contract.GetMinCollateral(&_P1Getters.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Getters *P1GettersCaller) GetOracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getOracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Getters *P1GettersSession) GetOracleContract() (common.Address, error) {
	return _P1Getters.Contract.GetOracleContract(&_P1Getters.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_P1Getters *P1GettersCallerSession) GetOracleContract() (common.Address, error) {
	return _P1Getters.Contract.GetOracleContract(&_P1Getters.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Getters *P1GettersCaller) GetOraclePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getOraclePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Getters *P1GettersSession) GetOraclePrice() (*big.Int, error) {
	return _P1Getters.Contract.GetOraclePrice(&_P1Getters.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_P1Getters *P1GettersCallerSession) GetOraclePrice() (*big.Int, error) {
	return _P1Getters.Contract.GetOraclePrice(&_P1Getters.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Getters *P1GettersCaller) GetTokenContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "getTokenContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Getters *P1GettersSession) GetTokenContract() (common.Address, error) {
	return _P1Getters.Contract.GetTokenContract(&_P1Getters.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_P1Getters *P1GettersCallerSession) GetTokenContract() (common.Address, error) {
	return _P1Getters.Contract.GetTokenContract(&_P1Getters.CallOpts)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersCaller) HasAccountPermissions(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _P1Getters.contract.Call(opts, &out, "hasAccountPermissions", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _P1Getters.Contract.HasAccountPermissions(&_P1Getters.CallOpts, account, operator)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_P1Getters *P1GettersCallerSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _P1Getters.Contract.HasAccountPermissions(&_P1Getters.CallOpts, account, operator)
}
