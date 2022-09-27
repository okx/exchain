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

// PerpetualV1MetaData contains all meta data concerning the PerpetualV1 contract.
var PerpetualV1MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogAccountSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"settlementPrice\",\"type\":\"uint256\"}],\"name\":\"LogFinalSettlementEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"index\",\"type\":\"bytes32\"}],\"name\":\"LogIndex\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"}],\"name\":\"LogSetFunder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"LogSetGlobalOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"LogSetLocalOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minCollateral\",\"type\":\"uint256\"}],\"name\":\"LogSetMinCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"LogSetOracle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"marginAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"positionAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"makerBalance\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"takerBalance\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"balance\",\"type\":\"bytes32\"}],\"name\":\"LogWithdrawFinalSettlement\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"priceLowerBound\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"priceUpperBound\",\"type\":\"uint256\"}],\"name\":\"enableFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountBalance\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFinalSettlementEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFunderContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGlobalIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOracleContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOraclePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"hasAccountPermissions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"}],\"name\":\"setFunder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setGlobalOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setLocalOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minCollateral\",\"type\":\"uint256\"}],\"name\":\"setMinCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"takerIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerIndex\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structP1Trade.TradeArg[]\",\"name\":\"trades\",\"type\":\"tuple[]\"}],\"name\":\"trade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minCollateral\",\"type\":\"uint256\"}],\"name\":\"initializeV1\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PerpetualV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use PerpetualV1MetaData.ABI instead.
var PerpetualV1ABI = PerpetualV1MetaData.ABI

// PerpetualV1 is an auto generated Go binding around an Ethereum contract.
type PerpetualV1 struct {
	PerpetualV1Caller     // Read-only binding to the contract
	PerpetualV1Transactor // Write-only binding to the contract
	PerpetualV1Filterer   // Log filterer for contract events
}

// PerpetualV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type PerpetualV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type PerpetualV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PerpetualV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpetualV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PerpetualV1Session struct {
	Contract     *PerpetualV1      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PerpetualV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PerpetualV1CallerSession struct {
	Contract *PerpetualV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PerpetualV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PerpetualV1TransactorSession struct {
	Contract     *PerpetualV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PerpetualV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type PerpetualV1Raw struct {
	Contract *PerpetualV1 // Generic contract binding to access the raw methods on
}

// PerpetualV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PerpetualV1CallerRaw struct {
	Contract *PerpetualV1Caller // Generic read-only contract binding to access the raw methods on
}

// PerpetualV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PerpetualV1TransactorRaw struct {
	Contract *PerpetualV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewPerpetualV1 creates a new instance of PerpetualV1, bound to a specific deployed contract.
func NewPerpetualV1(address common.Address, backend bind.ContractBackend) (*PerpetualV1, error) {
	contract, err := bindPerpetualV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1{PerpetualV1Caller: PerpetualV1Caller{contract: contract}, PerpetualV1Transactor: PerpetualV1Transactor{contract: contract}, PerpetualV1Filterer: PerpetualV1Filterer{contract: contract}}, nil
}

// NewPerpetualV1Caller creates a new read-only instance of PerpetualV1, bound to a specific deployed contract.
func NewPerpetualV1Caller(address common.Address, caller bind.ContractCaller) (*PerpetualV1Caller, error) {
	contract, err := bindPerpetualV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1Caller{contract: contract}, nil
}

// NewPerpetualV1Transactor creates a new write-only instance of PerpetualV1, bound to a specific deployed contract.
func NewPerpetualV1Transactor(address common.Address, transactor bind.ContractTransactor) (*PerpetualV1Transactor, error) {
	contract, err := bindPerpetualV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1Transactor{contract: contract}, nil
}

// NewPerpetualV1Filterer creates a new log filterer instance of PerpetualV1, bound to a specific deployed contract.
func NewPerpetualV1Filterer(address common.Address, filterer bind.ContractFilterer) (*PerpetualV1Filterer, error) {
	contract, err := bindPerpetualV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1Filterer{contract: contract}, nil
}

// bindPerpetualV1 binds a generic wrapper to an already deployed contract.
func bindPerpetualV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PerpetualV1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpetualV1 *PerpetualV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpetualV1.Contract.PerpetualV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpetualV1 *PerpetualV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualV1.Contract.PerpetualV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpetualV1 *PerpetualV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpetualV1.Contract.PerpetualV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpetualV1 *PerpetualV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpetualV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpetualV1 *PerpetualV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpetualV1 *PerpetualV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpetualV1.Contract.contract.Transact(opts, method, params...)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_PerpetualV1 *PerpetualV1Caller) GetAccountBalance(opts *bind.CallOpts, account common.Address) (P1TypesBalance, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getAccountBalance", account)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_PerpetualV1 *PerpetualV1Session) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _PerpetualV1.Contract.GetAccountBalance(&_PerpetualV1.CallOpts, account)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_PerpetualV1 *PerpetualV1CallerSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _PerpetualV1.Contract.GetAccountBalance(&_PerpetualV1.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1Caller) GetAccountIndex(opts *bind.CallOpts, account common.Address) (P1TypesIndex, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getAccountIndex", account)

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1Session) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _PerpetualV1.Contract.GetAccountIndex(&_PerpetualV1.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1CallerSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _PerpetualV1.Contract.GetAccountIndex(&_PerpetualV1.CallOpts, account)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_PerpetualV1 *PerpetualV1Caller) GetAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_PerpetualV1 *PerpetualV1Session) GetAdmin() (common.Address, error) {
	return _PerpetualV1.Contract.GetAdmin(&_PerpetualV1.CallOpts)
}

// GetAdmin is a free data retrieval call binding the contract method 0x6e9960c3.
//
// Solidity: function getAdmin() view returns(address)
func (_PerpetualV1 *PerpetualV1CallerSession) GetAdmin() (common.Address, error) {
	return _PerpetualV1.Contract.GetAdmin(&_PerpetualV1.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_PerpetualV1 *PerpetualV1Caller) GetFinalSettlementEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getFinalSettlementEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_PerpetualV1 *PerpetualV1Session) GetFinalSettlementEnabled() (bool, error) {
	return _PerpetualV1.Contract.GetFinalSettlementEnabled(&_PerpetualV1.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_PerpetualV1 *PerpetualV1CallerSession) GetFinalSettlementEnabled() (bool, error) {
	return _PerpetualV1.Contract.GetFinalSettlementEnabled(&_PerpetualV1.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Caller) GetFunderContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getFunderContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Session) GetFunderContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetFunderContract(&_PerpetualV1.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_PerpetualV1 *PerpetualV1CallerSession) GetFunderContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetFunderContract(&_PerpetualV1.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1Caller) GetGlobalIndex(opts *bind.CallOpts) (P1TypesIndex, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getGlobalIndex")

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1Session) GetGlobalIndex() (P1TypesIndex, error) {
	return _PerpetualV1.Contract.GetGlobalIndex(&_PerpetualV1.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_PerpetualV1 *PerpetualV1CallerSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _PerpetualV1.Contract.GetGlobalIndex(&_PerpetualV1.CallOpts)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Caller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Session) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.GetIsGlobalOperator(&_PerpetualV1.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1CallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.GetIsGlobalOperator(&_PerpetualV1.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Caller) GetIsLocalOperator(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getIsLocalOperator", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Session) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.GetIsLocalOperator(&_PerpetualV1.CallOpts, account, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1CallerSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.GetIsLocalOperator(&_PerpetualV1.CallOpts, account, operator)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_PerpetualV1 *PerpetualV1Caller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_PerpetualV1 *PerpetualV1Session) GetMinCollateral() (*big.Int, error) {
	return _PerpetualV1.Contract.GetMinCollateral(&_PerpetualV1.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_PerpetualV1 *PerpetualV1CallerSession) GetMinCollateral() (*big.Int, error) {
	return _PerpetualV1.Contract.GetMinCollateral(&_PerpetualV1.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Caller) GetOracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getOracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Session) GetOracleContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetOracleContract(&_PerpetualV1.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_PerpetualV1 *PerpetualV1CallerSession) GetOracleContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetOracleContract(&_PerpetualV1.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_PerpetualV1 *PerpetualV1Caller) GetOraclePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getOraclePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_PerpetualV1 *PerpetualV1Session) GetOraclePrice() (*big.Int, error) {
	return _PerpetualV1.Contract.GetOraclePrice(&_PerpetualV1.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_PerpetualV1 *PerpetualV1CallerSession) GetOraclePrice() (*big.Int, error) {
	return _PerpetualV1.Contract.GetOraclePrice(&_PerpetualV1.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Caller) GetTokenContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "getTokenContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_PerpetualV1 *PerpetualV1Session) GetTokenContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetTokenContract(&_PerpetualV1.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_PerpetualV1 *PerpetualV1CallerSession) GetTokenContract() (common.Address, error) {
	return _PerpetualV1.Contract.GetTokenContract(&_PerpetualV1.CallOpts)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Caller) HasAccountPermissions(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _PerpetualV1.contract.Call(opts, &out, "hasAccountPermissions", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1Session) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.HasAccountPermissions(&_PerpetualV1.CallOpts, account, operator)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_PerpetualV1 *PerpetualV1CallerSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _PerpetualV1.Contract.HasAccountPermissions(&_PerpetualV1.CallOpts, account, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1Transactor) Deposit(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "deposit", account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1Session) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Deposit(&_PerpetualV1.TransactOpts, account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Deposit(&_PerpetualV1.TransactOpts, account, amount)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_PerpetualV1 *PerpetualV1Transactor) EnableFinalSettlement(opts *bind.TransactOpts, priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "enableFinalSettlement", priceLowerBound, priceUpperBound)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_PerpetualV1 *PerpetualV1Session) EnableFinalSettlement(priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.EnableFinalSettlement(&_PerpetualV1.TransactOpts, priceLowerBound, priceUpperBound)
}

// EnableFinalSettlement is a paid mutator transaction binding the contract method 0xf40c3699.
//
// Solidity: function enableFinalSettlement(uint256 priceLowerBound, uint256 priceUpperBound) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) EnableFinalSettlement(priceLowerBound *big.Int, priceUpperBound *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.EnableFinalSettlement(&_PerpetualV1.TransactOpts, priceLowerBound, priceUpperBound)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0xa895155b.
//
// Solidity: function initializeV1(address token, address oracle, address funder, uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1Transactor) InitializeV1(opts *bind.TransactOpts, token common.Address, oracle common.Address, funder common.Address, minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "initializeV1", token, oracle, funder, minCollateral)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0xa895155b.
//
// Solidity: function initializeV1(address token, address oracle, address funder, uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1Session) InitializeV1(token common.Address, oracle common.Address, funder common.Address, minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.InitializeV1(&_PerpetualV1.TransactOpts, token, oracle, funder, minCollateral)
}

// InitializeV1 is a paid mutator transaction binding the contract method 0xa895155b.
//
// Solidity: function initializeV1(address token, address oracle, address funder, uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) InitializeV1(token common.Address, oracle common.Address, funder common.Address, minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.InitializeV1(&_PerpetualV1.TransactOpts, token, oracle, funder, minCollateral)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_PerpetualV1 *PerpetualV1Transactor) SetFunder(opts *bind.TransactOpts, funder common.Address) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "setFunder", funder)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_PerpetualV1 *PerpetualV1Session) SetFunder(funder common.Address) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetFunder(&_PerpetualV1.TransactOpts, funder)
}

// SetFunder is a paid mutator transaction binding the contract method 0x0acc8cd1.
//
// Solidity: function setFunder(address funder) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) SetFunder(funder common.Address) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetFunder(&_PerpetualV1.TransactOpts, funder)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1Transactor) SetGlobalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "setGlobalOperator", operator, approved)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1Session) SetGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetGlobalOperator(&_PerpetualV1.TransactOpts, operator, approved)
}

// SetGlobalOperator is a paid mutator transaction binding the contract method 0x46d256c5.
//
// Solidity: function setGlobalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) SetGlobalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetGlobalOperator(&_PerpetualV1.TransactOpts, operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1Transactor) SetLocalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "setLocalOperator", operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1Session) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetLocalOperator(&_PerpetualV1.TransactOpts, operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetLocalOperator(&_PerpetualV1.TransactOpts, operator, approved)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1Transactor) SetMinCollateral(opts *bind.TransactOpts, minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "setMinCollateral", minCollateral)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1Session) SetMinCollateral(minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetMinCollateral(&_PerpetualV1.TransactOpts, minCollateral)
}

// SetMinCollateral is a paid mutator transaction binding the contract method 0x846321a4.
//
// Solidity: function setMinCollateral(uint256 minCollateral) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) SetMinCollateral(minCollateral *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetMinCollateral(&_PerpetualV1.TransactOpts, minCollateral)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_PerpetualV1 *PerpetualV1Transactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_PerpetualV1 *PerpetualV1Session) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetOracle(&_PerpetualV1.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _PerpetualV1.Contract.SetOracle(&_PerpetualV1.TransactOpts, oracle)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_PerpetualV1 *PerpetualV1Transactor) Trade(opts *bind.TransactOpts, accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "trade", accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_PerpetualV1 *PerpetualV1Session) Trade(accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Trade(&_PerpetualV1.TransactOpts, accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) Trade(accounts []common.Address, trades []P1TradeTradeArg) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Trade(&_PerpetualV1.TransactOpts, accounts, trades)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1Transactor) Withdraw(opts *bind.TransactOpts, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "withdraw", account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1Session) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Withdraw(&_PerpetualV1.TransactOpts, account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PerpetualV1.Contract.Withdraw(&_PerpetualV1.TransactOpts, account, destination, amount)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_PerpetualV1 *PerpetualV1Transactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpetualV1.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_PerpetualV1 *PerpetualV1Session) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _PerpetualV1.Contract.WithdrawFinalSettlement(&_PerpetualV1.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_PerpetualV1 *PerpetualV1TransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _PerpetualV1.Contract.WithdrawFinalSettlement(&_PerpetualV1.TransactOpts)
}

// PerpetualV1LogAccountSettledIterator is returned from FilterLogAccountSettled and is used to iterate over the raw logs and unpacked data for LogAccountSettled events raised by the PerpetualV1 contract.
type PerpetualV1LogAccountSettledIterator struct {
	Event *PerpetualV1LogAccountSettled // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogAccountSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogAccountSettled)
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
		it.Event = new(PerpetualV1LogAccountSettled)
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
func (it *PerpetualV1LogAccountSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogAccountSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogAccountSettled represents a LogAccountSettled event raised by the PerpetualV1 contract.
type PerpetualV1LogAccountSettled struct {
	Account    common.Address
	IsPositive bool
	Amount     *big.Int
	Balance    [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogAccountSettled is a free log retrieval operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogAccountSettled(opts *bind.FilterOpts, account []common.Address) (*PerpetualV1LogAccountSettledIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogAccountSettledIterator{contract: _PerpetualV1.contract, event: "LogAccountSettled", logs: logs, sub: sub}, nil
}

// WatchLogAccountSettled is a free log subscription operation binding the contract event 0x022694ffbbd957d26de6b85c040be68ec582d13d40114b29130581793a1bf31e.
//
// Solidity: event LogAccountSettled(address indexed account, bool isPositive, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogAccountSettled(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogAccountSettled, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogAccountSettled", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogAccountSettled)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogAccountSettled(log types.Log) (*PerpetualV1LogAccountSettled, error) {
	event := new(PerpetualV1LogAccountSettled)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogAccountSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogDepositIterator is returned from FilterLogDeposit and is used to iterate over the raw logs and unpacked data for LogDeposit events raised by the PerpetualV1 contract.
type PerpetualV1LogDepositIterator struct {
	Event *PerpetualV1LogDeposit // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogDeposit)
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
		it.Event = new(PerpetualV1LogDeposit)
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
func (it *PerpetualV1LogDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogDeposit represents a LogDeposit event raised by the PerpetualV1 contract.
type PerpetualV1LogDeposit struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogDeposit is a free log retrieval operation binding the contract event 0x40a9cb3a9707d3a68091d8ef7ffd4158d01d0b2ad92b1e489abe8312dd543023.
//
// Solidity: event LogDeposit(address indexed account, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogDeposit(opts *bind.FilterOpts, account []common.Address) (*PerpetualV1LogDepositIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogDepositIterator{contract: _PerpetualV1.contract, event: "LogDeposit", logs: logs, sub: sub}, nil
}

// WatchLogDeposit is a free log subscription operation binding the contract event 0x40a9cb3a9707d3a68091d8ef7ffd4158d01d0b2ad92b1e489abe8312dd543023.
//
// Solidity: event LogDeposit(address indexed account, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogDeposit(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogDeposit, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogDeposit", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogDeposit)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogDeposit", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogDeposit(log types.Log) (*PerpetualV1LogDeposit, error) {
	event := new(PerpetualV1LogDeposit)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogFinalSettlementEnabledIterator is returned from FilterLogFinalSettlementEnabled and is used to iterate over the raw logs and unpacked data for LogFinalSettlementEnabled events raised by the PerpetualV1 contract.
type PerpetualV1LogFinalSettlementEnabledIterator struct {
	Event *PerpetualV1LogFinalSettlementEnabled // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogFinalSettlementEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogFinalSettlementEnabled)
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
		it.Event = new(PerpetualV1LogFinalSettlementEnabled)
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
func (it *PerpetualV1LogFinalSettlementEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogFinalSettlementEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogFinalSettlementEnabled represents a LogFinalSettlementEnabled event raised by the PerpetualV1 contract.
type PerpetualV1LogFinalSettlementEnabled struct {
	SettlementPrice *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogFinalSettlementEnabled is a free log retrieval operation binding the contract event 0x68e4c41627e835051be46337f1542645a60c7e6d6ea79efc5f20bdadae5f88d2.
//
// Solidity: event LogFinalSettlementEnabled(uint256 settlementPrice)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogFinalSettlementEnabled(opts *bind.FilterOpts) (*PerpetualV1LogFinalSettlementEnabledIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogFinalSettlementEnabled")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogFinalSettlementEnabledIterator{contract: _PerpetualV1.contract, event: "LogFinalSettlementEnabled", logs: logs, sub: sub}, nil
}

// WatchLogFinalSettlementEnabled is a free log subscription operation binding the contract event 0x68e4c41627e835051be46337f1542645a60c7e6d6ea79efc5f20bdadae5f88d2.
//
// Solidity: event LogFinalSettlementEnabled(uint256 settlementPrice)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogFinalSettlementEnabled(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogFinalSettlementEnabled) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogFinalSettlementEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogFinalSettlementEnabled)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogFinalSettlementEnabled", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogFinalSettlementEnabled(log types.Log) (*PerpetualV1LogFinalSettlementEnabled, error) {
	event := new(PerpetualV1LogFinalSettlementEnabled)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogFinalSettlementEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogIndexIterator is returned from FilterLogIndex and is used to iterate over the raw logs and unpacked data for LogIndex events raised by the PerpetualV1 contract.
type PerpetualV1LogIndexIterator struct {
	Event *PerpetualV1LogIndex // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogIndexIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogIndex)
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
		it.Event = new(PerpetualV1LogIndex)
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
func (it *PerpetualV1LogIndexIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogIndexIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogIndex represents a LogIndex event raised by the PerpetualV1 contract.
type PerpetualV1LogIndex struct {
	Index [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogIndex is a free log retrieval operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogIndex(opts *bind.FilterOpts) (*PerpetualV1LogIndexIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogIndexIterator{contract: _PerpetualV1.contract, event: "LogIndex", logs: logs, sub: sub}, nil
}

// WatchLogIndex is a free log subscription operation binding the contract event 0x995e61c355733308eab39a59e1e1ac167274cdd1ad707fe4d13e127a01076428.
//
// Solidity: event LogIndex(bytes32 index)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogIndex(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogIndex) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogIndex")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogIndex)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogIndex", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogIndex(log types.Log) (*PerpetualV1LogIndex, error) {
	event := new(PerpetualV1LogIndex)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogIndex", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogSetFunderIterator is returned from FilterLogSetFunder and is used to iterate over the raw logs and unpacked data for LogSetFunder events raised by the PerpetualV1 contract.
type PerpetualV1LogSetFunderIterator struct {
	Event *PerpetualV1LogSetFunder // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogSetFunderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogSetFunder)
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
		it.Event = new(PerpetualV1LogSetFunder)
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
func (it *PerpetualV1LogSetFunderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogSetFunderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogSetFunder represents a LogSetFunder event raised by the PerpetualV1 contract.
type PerpetualV1LogSetFunder struct {
	Funder common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogSetFunder is a free log retrieval operation binding the contract event 0x433b5c8c9ff78f62114ee8804a916537fa42009ebac4965bfed953f771789e47.
//
// Solidity: event LogSetFunder(address funder)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogSetFunder(opts *bind.FilterOpts) (*PerpetualV1LogSetFunderIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogSetFunder")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogSetFunderIterator{contract: _PerpetualV1.contract, event: "LogSetFunder", logs: logs, sub: sub}, nil
}

// WatchLogSetFunder is a free log subscription operation binding the contract event 0x433b5c8c9ff78f62114ee8804a916537fa42009ebac4965bfed953f771789e47.
//
// Solidity: event LogSetFunder(address funder)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogSetFunder(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogSetFunder) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogSetFunder")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogSetFunder)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogSetFunder", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogSetFunder(log types.Log) (*PerpetualV1LogSetFunder, error) {
	event := new(PerpetualV1LogSetFunder)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogSetFunder", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogSetGlobalOperatorIterator is returned from FilterLogSetGlobalOperator and is used to iterate over the raw logs and unpacked data for LogSetGlobalOperator events raised by the PerpetualV1 contract.
type PerpetualV1LogSetGlobalOperatorIterator struct {
	Event *PerpetualV1LogSetGlobalOperator // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogSetGlobalOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogSetGlobalOperator)
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
		it.Event = new(PerpetualV1LogSetGlobalOperator)
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
func (it *PerpetualV1LogSetGlobalOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogSetGlobalOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogSetGlobalOperator represents a LogSetGlobalOperator event raised by the PerpetualV1 contract.
type PerpetualV1LogSetGlobalOperator struct {
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogSetGlobalOperator is a free log retrieval operation binding the contract event 0xeaeee7699e70e6b31ac89ec999ef6936b97ac1e364f0e1fcf584772372caa0d3.
//
// Solidity: event LogSetGlobalOperator(address operator, bool approved)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogSetGlobalOperator(opts *bind.FilterOpts) (*PerpetualV1LogSetGlobalOperatorIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogSetGlobalOperator")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogSetGlobalOperatorIterator{contract: _PerpetualV1.contract, event: "LogSetGlobalOperator", logs: logs, sub: sub}, nil
}

// WatchLogSetGlobalOperator is a free log subscription operation binding the contract event 0xeaeee7699e70e6b31ac89ec999ef6936b97ac1e364f0e1fcf584772372caa0d3.
//
// Solidity: event LogSetGlobalOperator(address operator, bool approved)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogSetGlobalOperator(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogSetGlobalOperator) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogSetGlobalOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogSetGlobalOperator)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogSetGlobalOperator", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogSetGlobalOperator(log types.Log) (*PerpetualV1LogSetGlobalOperator, error) {
	event := new(PerpetualV1LogSetGlobalOperator)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogSetGlobalOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogSetLocalOperatorIterator is returned from FilterLogSetLocalOperator and is used to iterate over the raw logs and unpacked data for LogSetLocalOperator events raised by the PerpetualV1 contract.
type PerpetualV1LogSetLocalOperatorIterator struct {
	Event *PerpetualV1LogSetLocalOperator // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogSetLocalOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogSetLocalOperator)
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
		it.Event = new(PerpetualV1LogSetLocalOperator)
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
func (it *PerpetualV1LogSetLocalOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogSetLocalOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogSetLocalOperator represents a LogSetLocalOperator event raised by the PerpetualV1 contract.
type PerpetualV1LogSetLocalOperator struct {
	Sender   common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogSetLocalOperator is a free log retrieval operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogSetLocalOperator(opts *bind.FilterOpts, sender []common.Address) (*PerpetualV1LogSetLocalOperatorIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogSetLocalOperator", senderRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogSetLocalOperatorIterator{contract: _PerpetualV1.contract, event: "LogSetLocalOperator", logs: logs, sub: sub}, nil
}

// WatchLogSetLocalOperator is a free log subscription operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogSetLocalOperator(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogSetLocalOperator, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogSetLocalOperator", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogSetLocalOperator)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogSetLocalOperator", log); err != nil {
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

// ParseLogSetLocalOperator is a log parse operation binding the contract event 0xfe9fa8ad7dbd5e50cbcd1a903ea64717cb80b02e6b737e74f7e2f070b3e4d15f.
//
// Solidity: event LogSetLocalOperator(address indexed sender, address operator, bool approved)
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogSetLocalOperator(log types.Log) (*PerpetualV1LogSetLocalOperator, error) {
	event := new(PerpetualV1LogSetLocalOperator)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogSetLocalOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogSetMinCollateralIterator is returned from FilterLogSetMinCollateral and is used to iterate over the raw logs and unpacked data for LogSetMinCollateral events raised by the PerpetualV1 contract.
type PerpetualV1LogSetMinCollateralIterator struct {
	Event *PerpetualV1LogSetMinCollateral // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogSetMinCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogSetMinCollateral)
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
		it.Event = new(PerpetualV1LogSetMinCollateral)
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
func (it *PerpetualV1LogSetMinCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogSetMinCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogSetMinCollateral represents a LogSetMinCollateral event raised by the PerpetualV1 contract.
type PerpetualV1LogSetMinCollateral struct {
	MinCollateral *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogSetMinCollateral is a free log retrieval operation binding the contract event 0x248b36ced4662a14c995e0872f00eb61be4e3dea3913226cdeb513d64728cdca.
//
// Solidity: event LogSetMinCollateral(uint256 minCollateral)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogSetMinCollateral(opts *bind.FilterOpts) (*PerpetualV1LogSetMinCollateralIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogSetMinCollateral")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogSetMinCollateralIterator{contract: _PerpetualV1.contract, event: "LogSetMinCollateral", logs: logs, sub: sub}, nil
}

// WatchLogSetMinCollateral is a free log subscription operation binding the contract event 0x248b36ced4662a14c995e0872f00eb61be4e3dea3913226cdeb513d64728cdca.
//
// Solidity: event LogSetMinCollateral(uint256 minCollateral)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogSetMinCollateral(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogSetMinCollateral) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogSetMinCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogSetMinCollateral)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogSetMinCollateral", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogSetMinCollateral(log types.Log) (*PerpetualV1LogSetMinCollateral, error) {
	event := new(PerpetualV1LogSetMinCollateral)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogSetMinCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogSetOracleIterator is returned from FilterLogSetOracle and is used to iterate over the raw logs and unpacked data for LogSetOracle events raised by the PerpetualV1 contract.
type PerpetualV1LogSetOracleIterator struct {
	Event *PerpetualV1LogSetOracle // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogSetOracleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogSetOracle)
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
		it.Event = new(PerpetualV1LogSetOracle)
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
func (it *PerpetualV1LogSetOracleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogSetOracleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogSetOracle represents a LogSetOracle event raised by the PerpetualV1 contract.
type PerpetualV1LogSetOracle struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogSetOracle is a free log retrieval operation binding the contract event 0xad675642c3cba5442815383698d42cd28889533d9671a6d32cffea58ef0874da.
//
// Solidity: event LogSetOracle(address oracle)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogSetOracle(opts *bind.FilterOpts) (*PerpetualV1LogSetOracleIterator, error) {

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogSetOracle")
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogSetOracleIterator{contract: _PerpetualV1.contract, event: "LogSetOracle", logs: logs, sub: sub}, nil
}

// WatchLogSetOracle is a free log subscription operation binding the contract event 0xad675642c3cba5442815383698d42cd28889533d9671a6d32cffea58ef0874da.
//
// Solidity: event LogSetOracle(address oracle)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogSetOracle(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogSetOracle) (event.Subscription, error) {

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogSetOracle")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogSetOracle)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogSetOracle", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogSetOracle(log types.Log) (*PerpetualV1LogSetOracle, error) {
	event := new(PerpetualV1LogSetOracle)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogSetOracle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogTradeIterator is returned from FilterLogTrade and is used to iterate over the raw logs and unpacked data for LogTrade events raised by the PerpetualV1 contract.
type PerpetualV1LogTradeIterator struct {
	Event *PerpetualV1LogTrade // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogTrade)
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
		it.Event = new(PerpetualV1LogTrade)
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
func (it *PerpetualV1LogTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogTrade represents a LogTrade event raised by the PerpetualV1 contract.
type PerpetualV1LogTrade struct {
	Maker          common.Address
	Taker          common.Address
	Trader         common.Address
	MarginAmount   *big.Int
	PositionAmount *big.Int
	IsBuy          bool
	MakerBalance   [32]byte
	TakerBalance   [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLogTrade is a free log retrieval operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogTrade(opts *bind.FilterOpts, maker []common.Address, taker []common.Address) (*PerpetualV1LogTradeIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogTrade", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogTradeIterator{contract: _PerpetualV1.contract, event: "LogTrade", logs: logs, sub: sub}, nil
}

// WatchLogTrade is a free log subscription operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogTrade(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogTrade, maker []common.Address, taker []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogTrade", makerRule, takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogTrade)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogTrade", log); err != nil {
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

// ParseLogTrade is a log parse operation binding the contract event 0x5171a2ba3550a103fd09ca39b7dcfdf328a5acef18e290c7802d69c8ba73d8d9.
//
// Solidity: event LogTrade(address indexed maker, address indexed taker, address trader, uint256 marginAmount, uint256 positionAmount, bool isBuy, bytes32 makerBalance, bytes32 takerBalance)
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogTrade(log types.Log) (*PerpetualV1LogTrade, error) {
	event := new(PerpetualV1LogTrade)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogTrade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogWithdrawIterator is returned from FilterLogWithdraw and is used to iterate over the raw logs and unpacked data for LogWithdraw events raised by the PerpetualV1 contract.
type PerpetualV1LogWithdrawIterator struct {
	Event *PerpetualV1LogWithdraw // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogWithdraw)
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
		it.Event = new(PerpetualV1LogWithdraw)
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
func (it *PerpetualV1LogWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogWithdraw represents a LogWithdraw event raised by the PerpetualV1 contract.
type PerpetualV1LogWithdraw struct {
	Account     common.Address
	Destination common.Address
	Amount      *big.Int
	Balance     [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLogWithdraw is a free log retrieval operation binding the contract event 0x74348e8cb927b5536fe550310d0cdf05914498fcb04ad61b99c29e3899b0bce9.
//
// Solidity: event LogWithdraw(address indexed account, address destination, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogWithdraw(opts *bind.FilterOpts, account []common.Address) (*PerpetualV1LogWithdrawIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogWithdraw", accountRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogWithdrawIterator{contract: _PerpetualV1.contract, event: "LogWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogWithdraw is a free log subscription operation binding the contract event 0x74348e8cb927b5536fe550310d0cdf05914498fcb04ad61b99c29e3899b0bce9.
//
// Solidity: event LogWithdraw(address indexed account, address destination, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogWithdraw(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogWithdraw, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogWithdraw", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogWithdraw)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogWithdraw(log types.Log) (*PerpetualV1LogWithdraw, error) {
	event := new(PerpetualV1LogWithdraw)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpetualV1LogWithdrawFinalSettlementIterator is returned from FilterLogWithdrawFinalSettlement and is used to iterate over the raw logs and unpacked data for LogWithdrawFinalSettlement events raised by the PerpetualV1 contract.
type PerpetualV1LogWithdrawFinalSettlementIterator struct {
	Event *PerpetualV1LogWithdrawFinalSettlement // Event containing the contract specifics and raw log

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
func (it *PerpetualV1LogWithdrawFinalSettlementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpetualV1LogWithdrawFinalSettlement)
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
		it.Event = new(PerpetualV1LogWithdrawFinalSettlement)
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
func (it *PerpetualV1LogWithdrawFinalSettlementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpetualV1LogWithdrawFinalSettlementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpetualV1LogWithdrawFinalSettlement represents a LogWithdrawFinalSettlement event raised by the PerpetualV1 contract.
type PerpetualV1LogWithdrawFinalSettlement struct {
	Account common.Address
	Amount  *big.Int
	Balance [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdrawFinalSettlement is a free log retrieval operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) FilterLogWithdrawFinalSettlement(opts *bind.FilterOpts, account []common.Address) (*PerpetualV1LogWithdrawFinalSettlementIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.FilterLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return &PerpetualV1LogWithdrawFinalSettlementIterator{contract: _PerpetualV1.contract, event: "LogWithdrawFinalSettlement", logs: logs, sub: sub}, nil
}

// WatchLogWithdrawFinalSettlement is a free log subscription operation binding the contract event 0xc3b34c584e097adcd5d59ecaf4107928698a4f075c7753b5dbe28cd20d7ac1fd.
//
// Solidity: event LogWithdrawFinalSettlement(address indexed account, uint256 amount, bytes32 balance)
func (_PerpetualV1 *PerpetualV1Filterer) WatchLogWithdrawFinalSettlement(opts *bind.WatchOpts, sink chan<- *PerpetualV1LogWithdrawFinalSettlement, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _PerpetualV1.contract.WatchLogs(opts, "LogWithdrawFinalSettlement", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpetualV1LogWithdrawFinalSettlement)
				if err := _PerpetualV1.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
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
func (_PerpetualV1 *PerpetualV1Filterer) ParseLogWithdrawFinalSettlement(log types.Log) (*PerpetualV1LogWithdrawFinalSettlement, error) {
	event := new(PerpetualV1LogWithdrawFinalSettlement)
	if err := _PerpetualV1.contract.UnpackLog(event, "LogWithdrawFinalSettlement", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
