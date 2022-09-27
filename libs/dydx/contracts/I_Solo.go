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

// I_SoloAccountInfo is an auto generated low-level Go binding around an user-defined struct.
type I_SoloAccountInfo struct {
	Owner  common.Address
	Number *big.Int
}

// I_SoloActionArgs is an auto generated low-level Go binding around an user-defined struct.
type I_SoloActionArgs struct {
	ActionType        uint8
	AccountId         *big.Int
	Amount            I_SoloAssetAmount
	PrimaryMarketId   *big.Int
	SecondaryMarketId *big.Int
	OtherAddress      common.Address
	OtherAccountId    *big.Int
	Data              []byte
}

// I_SoloAssetAmount is an auto generated low-level Go binding around an user-defined struct.
type I_SoloAssetAmount struct {
	Sign         bool
	Denomination uint8
	Ref          uint8
	Value        *big.Int
}

// ISoloMetaData contains all meta data concerning the ISolo contract.
var ISoloMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsLocalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getIsGlobalOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"marketId\",\"type\":\"uint256\"}],\"name\":\"getMarketTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"internalType\":\"structI_Solo.AccountInfo[]\",\"name\":\"accounts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumI_Solo.ActionType\",\"name\":\"actionType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"accountId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"sign\",\"type\":\"bool\"},{\"internalType\":\"enumI_Solo.AssetDenomination\",\"name\":\"denomination\",\"type\":\"uint8\"},{\"internalType\":\"enumI_Solo.AssetReference\",\"name\":\"ref\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structI_Solo.AssetAmount\",\"name\":\"amount\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"primaryMarketId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"secondaryMarketId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"otherAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"otherAccountId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structI_Solo.ActionArgs[]\",\"name\":\"actions\",\"type\":\"tuple[]\"}],\"name\":\"operate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ISoloABI is the input ABI used to generate the binding from.
// Deprecated: Use ISoloMetaData.ABI instead.
var ISoloABI = ISoloMetaData.ABI

// ISolo is an auto generated Go binding around an Ethereum contract.
type ISolo struct {
	ISoloCaller     // Read-only binding to the contract
	ISoloTransactor // Write-only binding to the contract
	ISoloFilterer   // Log filterer for contract events
}

// ISoloCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISoloCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISoloTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISoloTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISoloFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISoloFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISoloSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISoloSession struct {
	Contract     *ISolo            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISoloCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISoloCallerSession struct {
	Contract *ISoloCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ISoloTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISoloTransactorSession struct {
	Contract     *ISoloTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISoloRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISoloRaw struct {
	Contract *ISolo // Generic contract binding to access the raw methods on
}

// ISoloCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISoloCallerRaw struct {
	Contract *ISoloCaller // Generic read-only contract binding to access the raw methods on
}

// ISoloTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISoloTransactorRaw struct {
	Contract *ISoloTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISolo creates a new instance of ISolo, bound to a specific deployed contract.
func NewISolo(address common.Address, backend bind.ContractBackend) (*ISolo, error) {
	contract, err := bindISolo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISolo{ISoloCaller: ISoloCaller{contract: contract}, ISoloTransactor: ISoloTransactor{contract: contract}, ISoloFilterer: ISoloFilterer{contract: contract}}, nil
}

// NewISoloCaller creates a new read-only instance of ISolo, bound to a specific deployed contract.
func NewISoloCaller(address common.Address, caller bind.ContractCaller) (*ISoloCaller, error) {
	contract, err := bindISolo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISoloCaller{contract: contract}, nil
}

// NewISoloTransactor creates a new write-only instance of ISolo, bound to a specific deployed contract.
func NewISoloTransactor(address common.Address, transactor bind.ContractTransactor) (*ISoloTransactor, error) {
	contract, err := bindISolo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISoloTransactor{contract: contract}, nil
}

// NewISoloFilterer creates a new log filterer instance of ISolo, bound to a specific deployed contract.
func NewISoloFilterer(address common.Address, filterer bind.ContractFilterer) (*ISoloFilterer, error) {
	contract, err := bindISolo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISoloFilterer{contract: contract}, nil
}

// bindISolo binds a generic wrapper to an already deployed contract.
func bindISolo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISoloABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISolo *ISoloRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISolo.Contract.ISoloCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISolo *ISoloRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISolo.Contract.ISoloTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISolo *ISoloRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISolo.Contract.ISoloTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISolo *ISoloCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISolo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISolo *ISoloTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISolo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISolo *ISoloTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISolo.Contract.contract.Transact(opts, method, params...)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_ISolo *ISoloCaller) GetIsGlobalOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ISolo.contract.Call(opts, &out, "getIsGlobalOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_ISolo *ISoloSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _ISolo.Contract.GetIsGlobalOperator(&_ISolo.CallOpts, operator)
}

// GetIsGlobalOperator is a free data retrieval call binding the contract method 0x052f72d7.
//
// Solidity: function getIsGlobalOperator(address operator) view returns(bool)
func (_ISolo *ISoloCallerSession) GetIsGlobalOperator(operator common.Address) (bool, error) {
	return _ISolo.Contract.GetIsGlobalOperator(&_ISolo.CallOpts, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_ISolo *ISoloCaller) GetIsLocalOperator(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ISolo.contract.Call(opts, &out, "getIsLocalOperator", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_ISolo *ISoloSession) GetIsLocalOperator(owner common.Address, operator common.Address) (bool, error) {
	return _ISolo.Contract.GetIsLocalOperator(&_ISolo.CallOpts, owner, operator)
}

// GetIsLocalOperator is a free data retrieval call binding the contract method 0x3a031bf0.
//
// Solidity: function getIsLocalOperator(address owner, address operator) view returns(bool)
func (_ISolo *ISoloCallerSession) GetIsLocalOperator(owner common.Address, operator common.Address) (bool, error) {
	return _ISolo.Contract.GetIsLocalOperator(&_ISolo.CallOpts, owner, operator)
}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_ISolo *ISoloCaller) GetMarketTokenAddress(opts *bind.CallOpts, marketId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ISolo.contract.Call(opts, &out, "getMarketTokenAddress", marketId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_ISolo *ISoloSession) GetMarketTokenAddress(marketId *big.Int) (common.Address, error) {
	return _ISolo.Contract.GetMarketTokenAddress(&_ISolo.CallOpts, marketId)
}

// GetMarketTokenAddress is a free data retrieval call binding the contract method 0x062bd3e9.
//
// Solidity: function getMarketTokenAddress(uint256 marketId) view returns(address)
func (_ISolo *ISoloCallerSession) GetMarketTokenAddress(marketId *big.Int) (common.Address, error) {
	return _ISolo.Contract.GetMarketTokenAddress(&_ISolo.CallOpts, marketId)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_ISolo *ISoloTransactor) Operate(opts *bind.TransactOpts, accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _ISolo.contract.Transact(opts, "operate", accounts, actions)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_ISolo *ISoloSession) Operate(accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _ISolo.Contract.Operate(&_ISolo.TransactOpts, accounts, actions)
}

// Operate is a paid mutator transaction binding the contract method 0xa67a6a45.
//
// Solidity: function operate((address,uint256)[] accounts, (uint8,uint256,(bool,uint8,uint8,uint256),uint256,uint256,address,uint256,bytes)[] actions) returns()
func (_ISolo *ISoloTransactorSession) Operate(accounts []I_SoloAccountInfo, actions []I_SoloActionArgs) (*types.Transaction, error) {
	return _ISolo.Contract.Operate(&_ISolo.TransactOpts, accounts, actions)
}
