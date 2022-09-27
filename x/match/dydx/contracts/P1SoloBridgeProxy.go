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

// P1SoloBridgeProxyTransfer is an auto generated low-level Go binding around an user-defined struct.
type P1SoloBridgeProxyTransfer struct {
	Account           common.Address
	Perpetual         common.Address
	SoloAccountNumber *big.Int
	SoloMarketId      *big.Int
	Amount            *big.Int
	Options           [32]byte
}

// TypedSignatureSignature is an auto generated low-level Go binding around an user-defined struct.
type TypedSignatureSignature struct {
	R     [32]byte
	S     [32]byte
	VType [2]byte
}

// P1SoloBridgeProxyMetaData contains all meta data concerning the P1SoloBridgeProxy contract.
var P1SoloBridgeProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"soloMargin\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"transferHash\",\"type\":\"bytes32\"}],\"name\":\"LogSignatureInvalidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"soloAccountNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"soloMarketId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"toPerpetual\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"_EIP712_DOMAIN_HASH_\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_SIGNATURE_USED_\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"_SOLO_MARGIN_\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"}],\"name\":\"approveMaximumOnPerpetual\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"soloMarketId\",\"type\":\"uint256\"}],\"name\":\"approveMaximumOnSolo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"soloAccountNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"soloMarketId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"options\",\"type\":\"bytes32\"}],\"internalType\":\"structP1SoloBridgeProxy.Transfer\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"bytes2\",\"name\":\"vType\",\"type\":\"bytes2\"}],\"internalType\":\"structTypedSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"}],\"name\":\"bridgeTransfer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"perpetual\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"soloAccountNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"soloMarketId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"options\",\"type\":\"bytes32\"}],\"internalType\":\"structP1SoloBridgeProxy.Transfer\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"name\":\"invalidateSignature\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// P1SoloBridgeProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use P1SoloBridgeProxyMetaData.ABI instead.
var P1SoloBridgeProxyABI = P1SoloBridgeProxyMetaData.ABI

// P1SoloBridgeProxy is an auto generated Go binding around an Ethereum contract.
type P1SoloBridgeProxy struct {
	P1SoloBridgeProxyCaller     // Read-only binding to the contract
	P1SoloBridgeProxyTransactor // Write-only binding to the contract
	P1SoloBridgeProxyFilterer   // Log filterer for contract events
}

// P1SoloBridgeProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type P1SoloBridgeProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SoloBridgeProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type P1SoloBridgeProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SoloBridgeProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type P1SoloBridgeProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// P1SoloBridgeProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type P1SoloBridgeProxySession struct {
	Contract     *P1SoloBridgeProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// P1SoloBridgeProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type P1SoloBridgeProxyCallerSession struct {
	Contract *P1SoloBridgeProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// P1SoloBridgeProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type P1SoloBridgeProxyTransactorSession struct {
	Contract     *P1SoloBridgeProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// P1SoloBridgeProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type P1SoloBridgeProxyRaw struct {
	Contract *P1SoloBridgeProxy // Generic contract binding to access the raw methods on
}

// P1SoloBridgeProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type P1SoloBridgeProxyCallerRaw struct {
	Contract *P1SoloBridgeProxyCaller // Generic read-only contract binding to access the raw methods on
}

// P1SoloBridgeProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type P1SoloBridgeProxyTransactorRaw struct {
	Contract *P1SoloBridgeProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewP1SoloBridgeProxy creates a new instance of P1SoloBridgeProxy, bound to a specific deployed contract.
func NewP1SoloBridgeProxy(address common.Address, backend bind.ContractBackend) (*P1SoloBridgeProxy, error) {
	contract, err := bindP1SoloBridgeProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxy{P1SoloBridgeProxyCaller: P1SoloBridgeProxyCaller{contract: contract}, P1SoloBridgeProxyTransactor: P1SoloBridgeProxyTransactor{contract: contract}, P1SoloBridgeProxyFilterer: P1SoloBridgeProxyFilterer{contract: contract}}, nil
}

// NewP1SoloBridgeProxyCaller creates a new read-only instance of P1SoloBridgeProxy, bound to a specific deployed contract.
func NewP1SoloBridgeProxyCaller(address common.Address, caller bind.ContractCaller) (*P1SoloBridgeProxyCaller, error) {
	contract, err := bindP1SoloBridgeProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxyCaller{contract: contract}, nil
}

// NewP1SoloBridgeProxyTransactor creates a new write-only instance of P1SoloBridgeProxy, bound to a specific deployed contract.
func NewP1SoloBridgeProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*P1SoloBridgeProxyTransactor, error) {
	contract, err := bindP1SoloBridgeProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxyTransactor{contract: contract}, nil
}

// NewP1SoloBridgeProxyFilterer creates a new log filterer instance of P1SoloBridgeProxy, bound to a specific deployed contract.
func NewP1SoloBridgeProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*P1SoloBridgeProxyFilterer, error) {
	contract, err := bindP1SoloBridgeProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxyFilterer{contract: contract}, nil
}

// bindP1SoloBridgeProxy binds a generic wrapper to an already deployed contract.
func bindP1SoloBridgeProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(P1SoloBridgeProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1SoloBridgeProxy.Contract.P1SoloBridgeProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.P1SoloBridgeProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.P1SoloBridgeProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _P1SoloBridgeProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.contract.Transact(opts, method, params...)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCaller) EIP712DOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _P1SoloBridgeProxy.contract.Call(opts, &out, "_EIP712_DOMAIN_HASH_")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) EIP712DOMAINHASH() ([32]byte, error) {
	return _P1SoloBridgeProxy.Contract.EIP712DOMAINHASH(&_P1SoloBridgeProxy.CallOpts)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xc7dc03f9.
//
// Solidity: function _EIP712_DOMAIN_HASH_() view returns(bytes32)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCallerSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _P1SoloBridgeProxy.Contract.EIP712DOMAINHASH(&_P1SoloBridgeProxy.CallOpts)
}

// SIGNATUREUSED is a free data retrieval call binding the contract method 0x0352ddfc.
//
// Solidity: function _SIGNATURE_USED_(bytes32 ) view returns(bool)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCaller) SIGNATUREUSED(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _P1SoloBridgeProxy.contract.Call(opts, &out, "_SIGNATURE_USED_", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SIGNATUREUSED is a free data retrieval call binding the contract method 0x0352ddfc.
//
// Solidity: function _SIGNATURE_USED_(bytes32 ) view returns(bool)
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) SIGNATUREUSED(arg0 [32]byte) (bool, error) {
	return _P1SoloBridgeProxy.Contract.SIGNATUREUSED(&_P1SoloBridgeProxy.CallOpts, arg0)
}

// SIGNATUREUSED is a free data retrieval call binding the contract method 0x0352ddfc.
//
// Solidity: function _SIGNATURE_USED_(bytes32 ) view returns(bool)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCallerSession) SIGNATUREUSED(arg0 [32]byte) (bool, error) {
	return _P1SoloBridgeProxy.Contract.SIGNATUREUSED(&_P1SoloBridgeProxy.CallOpts, arg0)
}

// SOLOMARGIN is a free data retrieval call binding the contract method 0x3f41499d.
//
// Solidity: function _SOLO_MARGIN_() view returns(address)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCaller) SOLOMARGIN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _P1SoloBridgeProxy.contract.Call(opts, &out, "_SOLO_MARGIN_")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SOLOMARGIN is a free data retrieval call binding the contract method 0x3f41499d.
//
// Solidity: function _SOLO_MARGIN_() view returns(address)
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) SOLOMARGIN() (common.Address, error) {
	return _P1SoloBridgeProxy.Contract.SOLOMARGIN(&_P1SoloBridgeProxy.CallOpts)
}

// SOLOMARGIN is a free data retrieval call binding the contract method 0x3f41499d.
//
// Solidity: function _SOLO_MARGIN_() view returns(address)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyCallerSession) SOLOMARGIN() (common.Address, error) {
	return _P1SoloBridgeProxy.Contract.SOLOMARGIN(&_P1SoloBridgeProxy.CallOpts)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactor) ApproveMaximumOnPerpetual(opts *bind.TransactOpts, perpetual common.Address) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.contract.Transact(opts, "approveMaximumOnPerpetual", perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.ApproveMaximumOnPerpetual(&_P1SoloBridgeProxy.TransactOpts, perpetual)
}

// ApproveMaximumOnPerpetual is a paid mutator transaction binding the contract method 0xfe0f8858.
//
// Solidity: function approveMaximumOnPerpetual(address perpetual) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorSession) ApproveMaximumOnPerpetual(perpetual common.Address) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.ApproveMaximumOnPerpetual(&_P1SoloBridgeProxy.TransactOpts, perpetual)
}

// ApproveMaximumOnSolo is a paid mutator transaction binding the contract method 0xf1beb6d7.
//
// Solidity: function approveMaximumOnSolo(uint256 soloMarketId) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactor) ApproveMaximumOnSolo(opts *bind.TransactOpts, soloMarketId *big.Int) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.contract.Transact(opts, "approveMaximumOnSolo", soloMarketId)
}

// ApproveMaximumOnSolo is a paid mutator transaction binding the contract method 0xf1beb6d7.
//
// Solidity: function approveMaximumOnSolo(uint256 soloMarketId) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) ApproveMaximumOnSolo(soloMarketId *big.Int) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.ApproveMaximumOnSolo(&_P1SoloBridgeProxy.TransactOpts, soloMarketId)
}

// ApproveMaximumOnSolo is a paid mutator transaction binding the contract method 0xf1beb6d7.
//
// Solidity: function approveMaximumOnSolo(uint256 soloMarketId) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorSession) ApproveMaximumOnSolo(soloMarketId *big.Int) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.ApproveMaximumOnSolo(&_P1SoloBridgeProxy.TransactOpts, soloMarketId)
}

// BridgeTransfer is a paid mutator transaction binding the contract method 0xb9fed3f6.
//
// Solidity: function bridgeTransfer((address,address,uint256,uint256,uint256,bytes32) transfer, (bytes32,bytes32,bytes2) signature) returns(uint256)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactor) BridgeTransfer(opts *bind.TransactOpts, transfer P1SoloBridgeProxyTransfer, signature TypedSignatureSignature) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.contract.Transact(opts, "bridgeTransfer", transfer, signature)
}

// BridgeTransfer is a paid mutator transaction binding the contract method 0xb9fed3f6.
//
// Solidity: function bridgeTransfer((address,address,uint256,uint256,uint256,bytes32) transfer, (bytes32,bytes32,bytes2) signature) returns(uint256)
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) BridgeTransfer(transfer P1SoloBridgeProxyTransfer, signature TypedSignatureSignature) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.BridgeTransfer(&_P1SoloBridgeProxy.TransactOpts, transfer, signature)
}

// BridgeTransfer is a paid mutator transaction binding the contract method 0xb9fed3f6.
//
// Solidity: function bridgeTransfer((address,address,uint256,uint256,uint256,bytes32) transfer, (bytes32,bytes32,bytes2) signature) returns(uint256)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorSession) BridgeTransfer(transfer P1SoloBridgeProxyTransfer, signature TypedSignatureSignature) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.BridgeTransfer(&_P1SoloBridgeProxy.TransactOpts, transfer, signature)
}

// InvalidateSignature is a paid mutator transaction binding the contract method 0x5249243c.
//
// Solidity: function invalidateSignature((address,address,uint256,uint256,uint256,bytes32) transfer) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactor) InvalidateSignature(opts *bind.TransactOpts, transfer P1SoloBridgeProxyTransfer) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.contract.Transact(opts, "invalidateSignature", transfer)
}

// InvalidateSignature is a paid mutator transaction binding the contract method 0x5249243c.
//
// Solidity: function invalidateSignature((address,address,uint256,uint256,uint256,bytes32) transfer) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxySession) InvalidateSignature(transfer P1SoloBridgeProxyTransfer) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.InvalidateSignature(&_P1SoloBridgeProxy.TransactOpts, transfer)
}

// InvalidateSignature is a paid mutator transaction binding the contract method 0x5249243c.
//
// Solidity: function invalidateSignature((address,address,uint256,uint256,uint256,bytes32) transfer) returns()
func (_P1SoloBridgeProxy *P1SoloBridgeProxyTransactorSession) InvalidateSignature(transfer P1SoloBridgeProxyTransfer) (*types.Transaction, error) {
	return _P1SoloBridgeProxy.Contract.InvalidateSignature(&_P1SoloBridgeProxy.TransactOpts, transfer)
}

// P1SoloBridgeProxyLogSignatureInvalidatedIterator is returned from FilterLogSignatureInvalidated and is used to iterate over the raw logs and unpacked data for LogSignatureInvalidated events raised by the P1SoloBridgeProxy contract.
type P1SoloBridgeProxyLogSignatureInvalidatedIterator struct {
	Event *P1SoloBridgeProxyLogSignatureInvalidated // Event containing the contract specifics and raw log

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
func (it *P1SoloBridgeProxyLogSignatureInvalidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1SoloBridgeProxyLogSignatureInvalidated)
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
		it.Event = new(P1SoloBridgeProxyLogSignatureInvalidated)
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
func (it *P1SoloBridgeProxyLogSignatureInvalidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1SoloBridgeProxyLogSignatureInvalidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1SoloBridgeProxyLogSignatureInvalidated represents a LogSignatureInvalidated event raised by the P1SoloBridgeProxy contract.
type P1SoloBridgeProxyLogSignatureInvalidated struct {
	Account      common.Address
	TransferHash [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogSignatureInvalidated is a free log retrieval operation binding the contract event 0xee0a433ca6f431de6cf9e5ae3dacef4ee7ef663bc14bfc1ee7b0dcca438ef5c3.
//
// Solidity: event LogSignatureInvalidated(address indexed account, bytes32 transferHash)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) FilterLogSignatureInvalidated(opts *bind.FilterOpts, account []common.Address) (*P1SoloBridgeProxyLogSignatureInvalidatedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1SoloBridgeProxy.contract.FilterLogs(opts, "LogSignatureInvalidated", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxyLogSignatureInvalidatedIterator{contract: _P1SoloBridgeProxy.contract, event: "LogSignatureInvalidated", logs: logs, sub: sub}, nil
}

// WatchLogSignatureInvalidated is a free log subscription operation binding the contract event 0xee0a433ca6f431de6cf9e5ae3dacef4ee7ef663bc14bfc1ee7b0dcca438ef5c3.
//
// Solidity: event LogSignatureInvalidated(address indexed account, bytes32 transferHash)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) WatchLogSignatureInvalidated(opts *bind.WatchOpts, sink chan<- *P1SoloBridgeProxyLogSignatureInvalidated, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1SoloBridgeProxy.contract.WatchLogs(opts, "LogSignatureInvalidated", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1SoloBridgeProxyLogSignatureInvalidated)
				if err := _P1SoloBridgeProxy.contract.UnpackLog(event, "LogSignatureInvalidated", log); err != nil {
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

// ParseLogSignatureInvalidated is a log parse operation binding the contract event 0xee0a433ca6f431de6cf9e5ae3dacef4ee7ef663bc14bfc1ee7b0dcca438ef5c3.
//
// Solidity: event LogSignatureInvalidated(address indexed account, bytes32 transferHash)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) ParseLogSignatureInvalidated(log types.Log) (*P1SoloBridgeProxyLogSignatureInvalidated, error) {
	event := new(P1SoloBridgeProxyLogSignatureInvalidated)
	if err := _P1SoloBridgeProxy.contract.UnpackLog(event, "LogSignatureInvalidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// P1SoloBridgeProxyLogTransferredIterator is returned from FilterLogTransferred and is used to iterate over the raw logs and unpacked data for LogTransferred events raised by the P1SoloBridgeProxy contract.
type P1SoloBridgeProxyLogTransferredIterator struct {
	Event *P1SoloBridgeProxyLogTransferred // Event containing the contract specifics and raw log

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
func (it *P1SoloBridgeProxyLogTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(P1SoloBridgeProxyLogTransferred)
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
		it.Event = new(P1SoloBridgeProxyLogTransferred)
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
func (it *P1SoloBridgeProxyLogTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *P1SoloBridgeProxyLogTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// P1SoloBridgeProxyLogTransferred represents a LogTransferred event raised by the P1SoloBridgeProxy contract.
type P1SoloBridgeProxyLogTransferred struct {
	Account           common.Address
	Perpetual         common.Address
	SoloAccountNumber *big.Int
	SoloMarketId      *big.Int
	ToPerpetual       bool
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogTransferred is a free log retrieval operation binding the contract event 0x44704b4f0be8f8a46df98e25b1b154fd1305d5c952ca3edee386cb73b34ae241.
//
// Solidity: event LogTransferred(address indexed account, address perpetual, uint256 soloAccountNumber, uint256 soloMarketId, bool toPerpetual, uint256 amount)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) FilterLogTransferred(opts *bind.FilterOpts, account []common.Address) (*P1SoloBridgeProxyLogTransferredIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1SoloBridgeProxy.contract.FilterLogs(opts, "LogTransferred", accountRule)
	if err != nil {
		return nil, err
	}
	return &P1SoloBridgeProxyLogTransferredIterator{contract: _P1SoloBridgeProxy.contract, event: "LogTransferred", logs: logs, sub: sub}, nil
}

// WatchLogTransferred is a free log subscription operation binding the contract event 0x44704b4f0be8f8a46df98e25b1b154fd1305d5c952ca3edee386cb73b34ae241.
//
// Solidity: event LogTransferred(address indexed account, address perpetual, uint256 soloAccountNumber, uint256 soloMarketId, bool toPerpetual, uint256 amount)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) WatchLogTransferred(opts *bind.WatchOpts, sink chan<- *P1SoloBridgeProxyLogTransferred, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _P1SoloBridgeProxy.contract.WatchLogs(opts, "LogTransferred", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(P1SoloBridgeProxyLogTransferred)
				if err := _P1SoloBridgeProxy.contract.UnpackLog(event, "LogTransferred", log); err != nil {
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

// ParseLogTransferred is a log parse operation binding the contract event 0x44704b4f0be8f8a46df98e25b1b154fd1305d5c952ca3edee386cb73b34ae241.
//
// Solidity: event LogTransferred(address indexed account, address perpetual, uint256 soloAccountNumber, uint256 soloMarketId, bool toPerpetual, uint256 amount)
func (_P1SoloBridgeProxy *P1SoloBridgeProxyFilterer) ParseLogTransferred(log types.Log) (*P1SoloBridgeProxyLogTransferred, error) {
	event := new(P1SoloBridgeProxyLogTransferred)
	if err := _P1SoloBridgeProxy.contract.UnpackLog(event, "LogTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
