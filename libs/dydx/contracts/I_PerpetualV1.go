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

// I_PerpetualV1TradeArg is an auto generated low-level Go binding around an user-defined struct.
type I_PerpetualV1TradeArg struct {
	TakerIndex *big.Int
	MakerIndex *big.Int
	Trader     common.Address
	Data       []byte
}

// P1TypesBalance is an auto generated low-level Go binding around an user-defined struct.
type P1TypesBalance struct {
	MarginIsPositive   bool
	PositionIsPositive bool
	Margin             *big.Int
	Position           *big.Int
}

// P1TypesIndex is an auto generated low-level Go binding around an user-defined struct.
type P1TypesIndex struct {
	Timestamp  uint32
	IsPositive bool
	Value      *big.Int
}

// IPerpetualV1MetaData contains all meta data concerning the IPerpetualV1 contract.
var IPerpetualV1MetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"takerIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerIndex\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structI_PerpetualV1.TradeArg[]\",\"name\":\"trades\",\"type\":\"tuple[]\"}],\"name\":\"trade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFinalSettlement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setLocalOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountBalance\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getAccountIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOracleContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFunderContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGlobalIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"internalType\":\"structP1Types.Index\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFinalSettlementEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"hasAccountPermissions\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOraclePrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IPerpetualV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use IPerpetualV1MetaData.ABI instead.
var IPerpetualV1ABI = IPerpetualV1MetaData.ABI

// IPerpetualV1 is an auto generated Go binding around an Ethereum contract.
type IPerpetualV1 struct {
	IPerpetualV1Caller     // Read-only binding to the contract
	IPerpetualV1Transactor // Write-only binding to the contract
	IPerpetualV1Filterer   // Log filterer for contract events
}

// IPerpetualV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type IPerpetualV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPerpetualV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IPerpetualV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPerpetualV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPerpetualV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPerpetualV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPerpetualV1Session struct {
	Contract     *IPerpetualV1     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPerpetualV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPerpetualV1CallerSession struct {
	Contract *IPerpetualV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IPerpetualV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPerpetualV1TransactorSession struct {
	Contract     *IPerpetualV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IPerpetualV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type IPerpetualV1Raw struct {
	Contract *IPerpetualV1 // Generic contract binding to access the raw methods on
}

// IPerpetualV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPerpetualV1CallerRaw struct {
	Contract *IPerpetualV1Caller // Generic read-only contract binding to access the raw methods on
}

// IPerpetualV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPerpetualV1TransactorRaw struct {
	Contract *IPerpetualV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIPerpetualV1 creates a new instance of IPerpetualV1, bound to a specific deployed contract.
func NewIPerpetualV1(address common.Address, backend bind.ContractBackend) (*IPerpetualV1, error) {
	contract, err := bindIPerpetualV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPerpetualV1{IPerpetualV1Caller: IPerpetualV1Caller{contract: contract}, IPerpetualV1Transactor: IPerpetualV1Transactor{contract: contract}, IPerpetualV1Filterer: IPerpetualV1Filterer{contract: contract}}, nil
}

// NewIPerpetualV1Caller creates a new read-only instance of IPerpetualV1, bound to a specific deployed contract.
func NewIPerpetualV1Caller(address common.Address, caller bind.ContractCaller) (*IPerpetualV1Caller, error) {
	contract, err := bindIPerpetualV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPerpetualV1Caller{contract: contract}, nil
}

// NewIPerpetualV1Transactor creates a new write-only instance of IPerpetualV1, bound to a specific deployed contract.
func NewIPerpetualV1Transactor(address common.Address, transactor bind.ContractTransactor) (*IPerpetualV1Transactor, error) {
	contract, err := bindIPerpetualV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPerpetualV1Transactor{contract: contract}, nil
}

// NewIPerpetualV1Filterer creates a new log filterer instance of IPerpetualV1, bound to a specific deployed contract.
func NewIPerpetualV1Filterer(address common.Address, filterer bind.ContractFilterer) (*IPerpetualV1Filterer, error) {
	contract, err := bindIPerpetualV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPerpetualV1Filterer{contract: contract}, nil
}

// bindIPerpetualV1 binds a generic wrapper to an already deployed contract.
func bindIPerpetualV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IPerpetualV1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPerpetualV1 *IPerpetualV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPerpetualV1.Contract.IPerpetualV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPerpetualV1 *IPerpetualV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.IPerpetualV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPerpetualV1 *IPerpetualV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.IPerpetualV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPerpetualV1 *IPerpetualV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPerpetualV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPerpetualV1 *IPerpetualV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPerpetualV1 *IPerpetualV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.contract.Transact(opts, method, params...)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_IPerpetualV1 *IPerpetualV1Caller) GetAccountBalance(opts *bind.CallOpts, account common.Address) (P1TypesBalance, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getAccountBalance", account)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_IPerpetualV1 *IPerpetualV1Session) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _IPerpetualV1.Contract.GetAccountBalance(&_IPerpetualV1.CallOpts, account)
}

// GetAccountBalance is a free data retrieval call binding the contract method 0x93423e9c.
//
// Solidity: function getAccountBalance(address account) view returns((bool,bool,uint120,uint120))
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetAccountBalance(account common.Address) (P1TypesBalance, error) {
	return _IPerpetualV1.Contract.GetAccountBalance(&_IPerpetualV1.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1Caller) GetAccountIndex(opts *bind.CallOpts, account common.Address) (P1TypesIndex, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getAccountIndex", account)

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1Session) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _IPerpetualV1.Contract.GetAccountIndex(&_IPerpetualV1.CallOpts, account)
}

// GetAccountIndex is a free data retrieval call binding the contract method 0x9ba63e9e.
//
// Solidity: function getAccountIndex(address account) view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetAccountIndex(account common.Address) (P1TypesIndex, error) {
	return _IPerpetualV1.Contract.GetAccountIndex(&_IPerpetualV1.CallOpts, account)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Caller) GetFinalSettlementEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getFinalSettlementEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Session) GetFinalSettlementEnabled() (bool, error) {
	return _IPerpetualV1.Contract.GetFinalSettlementEnabled(&_IPerpetualV1.CallOpts)
}

// GetFinalSettlementEnabled is a free data retrieval call binding the contract method 0x7099366b.
//
// Solidity: function getFinalSettlementEnabled() view returns(bool)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetFinalSettlementEnabled() (bool, error) {
	return _IPerpetualV1.Contract.GetFinalSettlementEnabled(&_IPerpetualV1.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Caller) GetFunderContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getFunderContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Session) GetFunderContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetFunderContract(&_IPerpetualV1.CallOpts)
}

// GetFunderContract is a free data retrieval call binding the contract method 0xdc4f3a0e.
//
// Solidity: function getFunderContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetFunderContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetFunderContract(&_IPerpetualV1.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1Caller) GetGlobalIndex(opts *bind.CallOpts) (P1TypesIndex, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getGlobalIndex")

	if err != nil {
		return *new(P1TypesIndex), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesIndex)).(*P1TypesIndex)

	return out0, err

}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1Session) GetGlobalIndex() (P1TypesIndex, error) {
	return _IPerpetualV1.Contract.GetGlobalIndex(&_IPerpetualV1.CallOpts)
}

// GetGlobalIndex is a free data retrieval call binding the contract method 0x80d63681.
//
// Solidity: function getGlobalIndex() view returns((uint32,bool,uint128))
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetGlobalIndex() (P1TypesIndex, error) {
	return _IPerpetualV1.Contract.GetGlobalIndex(&_IPerpetualV1.CallOpts)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Caller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Session) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.GetIsGlobalOperator(&_IPerpetualV1.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.GetIsGlobalOperator(&_IPerpetualV1.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Caller) GetIsLocalOperator(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getIsLocalOperator", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Session) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.GetIsLocalOperator(&_IPerpetualV1.CallOpts, account, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetIsLocalOperator(account common.Address, operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.GetIsLocalOperator(&_IPerpetualV1.CallOpts, account, operator)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1Caller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1Session) GetMinCollateral() (*big.Int, error) {
	return _IPerpetualV1.Contract.GetMinCollateral(&_IPerpetualV1.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetMinCollateral() (*big.Int, error) {
	return _IPerpetualV1.Contract.GetMinCollateral(&_IPerpetualV1.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Caller) GetOracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getOracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Session) GetOracleContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetOracleContract(&_IPerpetualV1.CallOpts)
}

// GetOracleContract is a free data retrieval call binding the contract method 0xe3bbb565.
//
// Solidity: function getOracleContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetOracleContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetOracleContract(&_IPerpetualV1.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1Caller) GetOraclePrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getOraclePrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1Session) GetOraclePrice() (*big.Int, error) {
	return _IPerpetualV1.Contract.GetOraclePrice(&_IPerpetualV1.CallOpts)
}

// GetOraclePrice is a free data retrieval call binding the contract method 0x796da7af.
//
// Solidity: function getOraclePrice() view returns(uint256)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetOraclePrice() (*big.Int, error) {
	return _IPerpetualV1.Contract.GetOraclePrice(&_IPerpetualV1.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Caller) GetTokenContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "getTokenContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1Session) GetTokenContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetTokenContract(&_IPerpetualV1.CallOpts)
}

// GetTokenContract is a free data retrieval call binding the contract method 0x28b7bede.
//
// Solidity: function getTokenContract() view returns(address)
func (_IPerpetualV1 *IPerpetualV1CallerSession) GetTokenContract() (common.Address, error) {
	return _IPerpetualV1.Contract.GetTokenContract(&_IPerpetualV1.CallOpts)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Caller) HasAccountPermissions(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _IPerpetualV1.contract.Call(opts, &out, "hasAccountPermissions", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1Session) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.HasAccountPermissions(&_IPerpetualV1.CallOpts, account, operator)
}

// HasAccountPermissions is a free data retrieval call binding the contract method 0x84ea2862.
//
// Solidity: function hasAccountPermissions(address account, address operator) view returns(bool)
func (_IPerpetualV1 *IPerpetualV1CallerSession) HasAccountPermissions(account common.Address, operator common.Address) (bool, error) {
	return _IPerpetualV1.Contract.HasAccountPermissions(&_IPerpetualV1.CallOpts, account, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1Transactor) Deposit(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.contract.Transact(opts, "deposit", account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1Session) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Deposit(&_IPerpetualV1.TransactOpts, account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address account, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1TransactorSession) Deposit(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Deposit(&_IPerpetualV1.TransactOpts, account, amount)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_IPerpetualV1 *IPerpetualV1Transactor) SetLocalOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _IPerpetualV1.contract.Transact(opts, "setLocalOperator", operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_IPerpetualV1 *IPerpetualV1Session) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.SetLocalOperator(&_IPerpetualV1.TransactOpts, operator, approved)
}

// SetLocalOperator is a paid mutator transaction binding the contract method 0xb4959e72.
//
// Solidity: function setLocalOperator(address operator, bool approved) returns()
func (_IPerpetualV1 *IPerpetualV1TransactorSession) SetLocalOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.SetLocalOperator(&_IPerpetualV1.TransactOpts, operator, approved)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_IPerpetualV1 *IPerpetualV1Transactor) Trade(opts *bind.TransactOpts, accounts []common.Address, trades []I_PerpetualV1TradeArg) (*types.Transaction, error) {
	return _IPerpetualV1.contract.Transact(opts, "trade", accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_IPerpetualV1 *IPerpetualV1Session) Trade(accounts []common.Address, trades []I_PerpetualV1TradeArg) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Trade(&_IPerpetualV1.TransactOpts, accounts, trades)
}

// Trade is a paid mutator transaction binding the contract method 0x68eec3f6.
//
// Solidity: function trade(address[] accounts, (uint256,uint256,address,bytes)[] trades) returns()
func (_IPerpetualV1 *IPerpetualV1TransactorSession) Trade(accounts []common.Address, trades []I_PerpetualV1TradeArg) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Trade(&_IPerpetualV1.TransactOpts, accounts, trades)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1Transactor) Withdraw(opts *bind.TransactOpts, account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.contract.Transact(opts, "withdraw", account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1Session) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Withdraw(&_IPerpetualV1.TransactOpts, account, destination, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address account, address destination, uint256 amount) returns()
func (_IPerpetualV1 *IPerpetualV1TransactorSession) Withdraw(account common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IPerpetualV1.Contract.Withdraw(&_IPerpetualV1.TransactOpts, account, destination, amount)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_IPerpetualV1 *IPerpetualV1Transactor) WithdrawFinalSettlement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPerpetualV1.contract.Transact(opts, "withdrawFinalSettlement")
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_IPerpetualV1 *IPerpetualV1Session) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _IPerpetualV1.Contract.WithdrawFinalSettlement(&_IPerpetualV1.TransactOpts)
}

// WithdrawFinalSettlement is a paid mutator transaction binding the contract method 0x142c69b3.
//
// Solidity: function withdrawFinalSettlement() returns()
func (_IPerpetualV1 *IPerpetualV1TransactorSession) WithdrawFinalSettlement() (*types.Transaction, error) {
	return _IPerpetualV1.Contract.WithdrawFinalSettlement(&_IPerpetualV1.TransactOpts)
}
