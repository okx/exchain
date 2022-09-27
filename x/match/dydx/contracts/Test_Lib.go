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

// TestLibMetaData contains all meta data concerning the TestLib contract.
var TestLibMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"base\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseValue\",\"type\":\"uint256\"}],\"name\":\"baseMul\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseValue\",\"type\":\"uint256\"}],\"name\":\"baseDivMul\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseValue\",\"type\":\"uint256\"}],\"name\":\"baseMulRoundUp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseValue\",\"type\":\"uint256\"}],\"name\":\"baseDiv\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"baseValue\",\"type\":\"uint256\"}],\"name\":\"baseReciprocal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"getFraction\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numerator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"getFractionRoundUp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"min\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"max\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"must\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"requireReason\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"that\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"toUint128\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"toUint120\",\"outputs\":[{\"internalType\":\"uint120\",\"name\":\"\",\"type\":\"uint120\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"toUint32\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"sint\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"add\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"sint\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"sub\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"augend\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"addend\",\"type\":\"tuple\"}],\"name\":\"signedAdd\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"minuend\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"subtrahend\",\"type\":\"tuple\"}],\"name\":\"signedSub\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"load\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signatureBytes\",\"type\":\"bytes\"}],\"name\":\"recover\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"copy\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addToMargin\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"subFromMargin\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addToPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"subFromPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"getPositiveAndNegativeValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"getMargin\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"getPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"newMargin\",\"type\":\"tuple\"}],\"name\":\"setMargin\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPositive\",\"type\":\"bool\"}],\"internalType\":\"structSignedMath.Int\",\"name\":\"newPosition\",\"type\":\"tuple\"}],\"name\":\"setPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"marginIsPositive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"positionIsPositive\",\"type\":\"bool\"},{\"internalType\":\"uint120\",\"name\":\"margin\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"position\",\"type\":\"uint120\"}],\"internalType\":\"structP1Types.Balance\",\"name\":\"\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"nonReentrant1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"nonReentrant2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestLibABI is the input ABI used to generate the binding from.
// Deprecated: Use TestLibMetaData.ABI instead.
var TestLibABI = TestLibMetaData.ABI

// TestLib is an auto generated Go binding around an Ethereum contract.
type TestLib struct {
	TestLibCaller     // Read-only binding to the contract
	TestLibTransactor // Write-only binding to the contract
	TestLibFilterer   // Log filterer for contract events
}

// TestLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestLibSession struct {
	Contract     *TestLib          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestLibCallerSession struct {
	Contract *TestLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// TestLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestLibTransactorSession struct {
	Contract     *TestLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// TestLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestLibRaw struct {
	Contract *TestLib // Generic contract binding to access the raw methods on
}

// TestLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestLibCallerRaw struct {
	Contract *TestLibCaller // Generic read-only contract binding to access the raw methods on
}

// TestLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestLibTransactorRaw struct {
	Contract *TestLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestLib creates a new instance of TestLib, bound to a specific deployed contract.
func NewTestLib(address common.Address, backend bind.ContractBackend) (*TestLib, error) {
	contract, err := bindTestLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestLib{TestLibCaller: TestLibCaller{contract: contract}, TestLibTransactor: TestLibTransactor{contract: contract}, TestLibFilterer: TestLibFilterer{contract: contract}}, nil
}

// NewTestLibCaller creates a new read-only instance of TestLib, bound to a specific deployed contract.
func NewTestLibCaller(address common.Address, caller bind.ContractCaller) (*TestLibCaller, error) {
	contract, err := bindTestLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestLibCaller{contract: contract}, nil
}

// NewTestLibTransactor creates a new write-only instance of TestLib, bound to a specific deployed contract.
func NewTestLibTransactor(address common.Address, transactor bind.ContractTransactor) (*TestLibTransactor, error) {
	contract, err := bindTestLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestLibTransactor{contract: contract}, nil
}

// NewTestLibFilterer creates a new log filterer instance of TestLib, bound to a specific deployed contract.
func NewTestLibFilterer(address common.Address, filterer bind.ContractFilterer) (*TestLibFilterer, error) {
	contract, err := bindTestLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestLibFilterer{contract: contract}, nil
}

// bindTestLib binds a generic wrapper to an already deployed contract.
func bindTestLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestLib *TestLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestLib.Contract.TestLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestLib *TestLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLib.Contract.TestLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestLib *TestLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestLib.Contract.TestLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestLib *TestLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestLib *TestLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestLib *TestLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestLib.Contract.contract.Transact(opts, method, params...)
}

// Add is a free data retrieval call binding the contract method 0x2953a626.
//
// Solidity: function add((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) Add(opts *bind.CallOpts, sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "add", sint, value)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// Add is a free data retrieval call binding the contract method 0x2953a626.
//
// Solidity: function add((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibSession) Add(sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	return _TestLib.Contract.Add(&_TestLib.CallOpts, sint, value)
}

// Add is a free data retrieval call binding the contract method 0x2953a626.
//
// Solidity: function add((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) Add(sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	return _TestLib.Contract.Add(&_TestLib.CallOpts, sint, value)
}

// AddToMargin is a free data retrieval call binding the contract method 0xe78af4ec.
//
// Solidity: function addToMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) AddToMargin(opts *bind.CallOpts, balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "addToMargin", balance, amount)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// AddToMargin is a free data retrieval call binding the contract method 0xe78af4ec.
//
// Solidity: function addToMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) AddToMargin(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.AddToMargin(&_TestLib.CallOpts, balance, amount)
}

// AddToMargin is a free data retrieval call binding the contract method 0xe78af4ec.
//
// Solidity: function addToMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) AddToMargin(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.AddToMargin(&_TestLib.CallOpts, balance, amount)
}

// AddToPosition is a free data retrieval call binding the contract method 0xcf178408.
//
// Solidity: function addToPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) AddToPosition(opts *bind.CallOpts, balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "addToPosition", balance, amount)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// AddToPosition is a free data retrieval call binding the contract method 0xcf178408.
//
// Solidity: function addToPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) AddToPosition(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.AddToPosition(&_TestLib.CallOpts, balance, amount)
}

// AddToPosition is a free data retrieval call binding the contract method 0xcf178408.
//
// Solidity: function addToPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) AddToPosition(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.AddToPosition(&_TestLib.CallOpts, balance, amount)
}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() pure returns(uint256)
func (_TestLib *TestLibCaller) Base(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "base")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() pure returns(uint256)
func (_TestLib *TestLibSession) Base() (*big.Int, error) {
	return _TestLib.Contract.Base(&_TestLib.CallOpts)
}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() pure returns(uint256)
func (_TestLib *TestLibCallerSession) Base() (*big.Int, error) {
	return _TestLib.Contract.Base(&_TestLib.CallOpts)
}

// BaseDiv is a free data retrieval call binding the contract method 0xded7c1c5.
//
// Solidity: function baseDiv(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCaller) BaseDiv(opts *bind.CallOpts, value *big.Int, baseValue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "baseDiv", value, baseValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseDiv is a free data retrieval call binding the contract method 0xded7c1c5.
//
// Solidity: function baseDiv(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibSession) BaseDiv(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseDiv(&_TestLib.CallOpts, value, baseValue)
}

// BaseDiv is a free data retrieval call binding the contract method 0xded7c1c5.
//
// Solidity: function baseDiv(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCallerSession) BaseDiv(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseDiv(&_TestLib.CallOpts, value, baseValue)
}

// BaseDivMul is a free data retrieval call binding the contract method 0x8f6561af.
//
// Solidity: function baseDivMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCaller) BaseDivMul(opts *bind.CallOpts, value *big.Int, baseValue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "baseDivMul", value, baseValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseDivMul is a free data retrieval call binding the contract method 0x8f6561af.
//
// Solidity: function baseDivMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibSession) BaseDivMul(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseDivMul(&_TestLib.CallOpts, value, baseValue)
}

// BaseDivMul is a free data retrieval call binding the contract method 0x8f6561af.
//
// Solidity: function baseDivMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCallerSession) BaseDivMul(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseDivMul(&_TestLib.CallOpts, value, baseValue)
}

// BaseMul is a free data retrieval call binding the contract method 0xce18b190.
//
// Solidity: function baseMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCaller) BaseMul(opts *bind.CallOpts, value *big.Int, baseValue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "baseMul", value, baseValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseMul is a free data retrieval call binding the contract method 0xce18b190.
//
// Solidity: function baseMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibSession) BaseMul(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseMul(&_TestLib.CallOpts, value, baseValue)
}

// BaseMul is a free data retrieval call binding the contract method 0xce18b190.
//
// Solidity: function baseMul(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCallerSession) BaseMul(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseMul(&_TestLib.CallOpts, value, baseValue)
}

// BaseMulRoundUp is a free data retrieval call binding the contract method 0x6394400f.
//
// Solidity: function baseMulRoundUp(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCaller) BaseMulRoundUp(opts *bind.CallOpts, value *big.Int, baseValue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "baseMulRoundUp", value, baseValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseMulRoundUp is a free data retrieval call binding the contract method 0x6394400f.
//
// Solidity: function baseMulRoundUp(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibSession) BaseMulRoundUp(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseMulRoundUp(&_TestLib.CallOpts, value, baseValue)
}

// BaseMulRoundUp is a free data retrieval call binding the contract method 0x6394400f.
//
// Solidity: function baseMulRoundUp(uint256 value, uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCallerSession) BaseMulRoundUp(value *big.Int, baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseMulRoundUp(&_TestLib.CallOpts, value, baseValue)
}

// BaseReciprocal is a free data retrieval call binding the contract method 0x7dc1e49a.
//
// Solidity: function baseReciprocal(uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCaller) BaseReciprocal(opts *bind.CallOpts, baseValue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "baseReciprocal", baseValue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseReciprocal is a free data retrieval call binding the contract method 0x7dc1e49a.
//
// Solidity: function baseReciprocal(uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibSession) BaseReciprocal(baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseReciprocal(&_TestLib.CallOpts, baseValue)
}

// BaseReciprocal is a free data retrieval call binding the contract method 0x7dc1e49a.
//
// Solidity: function baseReciprocal(uint256 baseValue) pure returns(uint256)
func (_TestLib *TestLibCallerSession) BaseReciprocal(baseValue *big.Int) (*big.Int, error) {
	return _TestLib.Contract.BaseReciprocal(&_TestLib.CallOpts, baseValue)
}

// Copy is a free data retrieval call binding the contract method 0xd46b4aff.
//
// Solidity: function copy((bool,bool,uint120,uint120) balance) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) Copy(opts *bind.CallOpts, balance P1TypesBalance) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "copy", balance)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// Copy is a free data retrieval call binding the contract method 0xd46b4aff.
//
// Solidity: function copy((bool,bool,uint120,uint120) balance) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) Copy(balance P1TypesBalance) (P1TypesBalance, error) {
	return _TestLib.Contract.Copy(&_TestLib.CallOpts, balance)
}

// Copy is a free data retrieval call binding the contract method 0xd46b4aff.
//
// Solidity: function copy((bool,bool,uint120,uint120) balance) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) Copy(balance P1TypesBalance) (P1TypesBalance, error) {
	return _TestLib.Contract.Copy(&_TestLib.CallOpts, balance)
}

// GetFraction is a free data retrieval call binding the contract method 0x5b827a5d.
//
// Solidity: function getFraction(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibCaller) GetFraction(opts *bind.CallOpts, target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "getFraction", target, numerator, denominator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFraction is a free data retrieval call binding the contract method 0x5b827a5d.
//
// Solidity: function getFraction(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibSession) GetFraction(target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	return _TestLib.Contract.GetFraction(&_TestLib.CallOpts, target, numerator, denominator)
}

// GetFraction is a free data retrieval call binding the contract method 0x5b827a5d.
//
// Solidity: function getFraction(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibCallerSession) GetFraction(target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	return _TestLib.Contract.GetFraction(&_TestLib.CallOpts, target, numerator, denominator)
}

// GetFractionRoundUp is a free data retrieval call binding the contract method 0x97259d26.
//
// Solidity: function getFractionRoundUp(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibCaller) GetFractionRoundUp(opts *bind.CallOpts, target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "getFractionRoundUp", target, numerator, denominator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFractionRoundUp is a free data retrieval call binding the contract method 0x97259d26.
//
// Solidity: function getFractionRoundUp(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibSession) GetFractionRoundUp(target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	return _TestLib.Contract.GetFractionRoundUp(&_TestLib.CallOpts, target, numerator, denominator)
}

// GetFractionRoundUp is a free data retrieval call binding the contract method 0x97259d26.
//
// Solidity: function getFractionRoundUp(uint256 target, uint256 numerator, uint256 denominator) pure returns(uint256)
func (_TestLib *TestLibCallerSession) GetFractionRoundUp(target *big.Int, numerator *big.Int, denominator *big.Int) (*big.Int, error) {
	return _TestLib.Contract.GetFractionRoundUp(&_TestLib.CallOpts, target, numerator, denominator)
}

// GetMargin is a free data retrieval call binding the contract method 0xc1f7ea1e.
//
// Solidity: function getMargin((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) GetMargin(opts *bind.CallOpts, balance P1TypesBalance) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "getMargin", balance)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// GetMargin is a free data retrieval call binding the contract method 0xc1f7ea1e.
//
// Solidity: function getMargin((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibSession) GetMargin(balance P1TypesBalance) (SignedMathInt, error) {
	return _TestLib.Contract.GetMargin(&_TestLib.CallOpts, balance)
}

// GetMargin is a free data retrieval call binding the contract method 0xc1f7ea1e.
//
// Solidity: function getMargin((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) GetMargin(balance P1TypesBalance) (SignedMathInt, error) {
	return _TestLib.Contract.GetMargin(&_TestLib.CallOpts, balance)
}

// GetPosition is a free data retrieval call binding the contract method 0x1491ac01.
//
// Solidity: function getPosition((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) GetPosition(opts *bind.CallOpts, balance P1TypesBalance) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "getPosition", balance)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// GetPosition is a free data retrieval call binding the contract method 0x1491ac01.
//
// Solidity: function getPosition((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibSession) GetPosition(balance P1TypesBalance) (SignedMathInt, error) {
	return _TestLib.Contract.GetPosition(&_TestLib.CallOpts, balance)
}

// GetPosition is a free data retrieval call binding the contract method 0x1491ac01.
//
// Solidity: function getPosition((bool,bool,uint120,uint120) balance) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) GetPosition(balance P1TypesBalance) (SignedMathInt, error) {
	return _TestLib.Contract.GetPosition(&_TestLib.CallOpts, balance)
}

// GetPositiveAndNegativeValue is a free data retrieval call binding the contract method 0xd5ae406c.
//
// Solidity: function getPositiveAndNegativeValue((bool,bool,uint120,uint120) balance, uint256 price) pure returns(uint256, uint256)
func (_TestLib *TestLibCaller) GetPositiveAndNegativeValue(opts *bind.CallOpts, balance P1TypesBalance, price *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "getPositiveAndNegativeValue", balance, price)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetPositiveAndNegativeValue is a free data retrieval call binding the contract method 0xd5ae406c.
//
// Solidity: function getPositiveAndNegativeValue((bool,bool,uint120,uint120) balance, uint256 price) pure returns(uint256, uint256)
func (_TestLib *TestLibSession) GetPositiveAndNegativeValue(balance P1TypesBalance, price *big.Int) (*big.Int, *big.Int, error) {
	return _TestLib.Contract.GetPositiveAndNegativeValue(&_TestLib.CallOpts, balance, price)
}

// GetPositiveAndNegativeValue is a free data retrieval call binding the contract method 0xd5ae406c.
//
// Solidity: function getPositiveAndNegativeValue((bool,bool,uint120,uint120) balance, uint256 price) pure returns(uint256, uint256)
func (_TestLib *TestLibCallerSession) GetPositiveAndNegativeValue(balance P1TypesBalance, price *big.Int) (*big.Int, *big.Int, error) {
	return _TestLib.Contract.GetPositiveAndNegativeValue(&_TestLib.CallOpts, balance, price)
}

// Load is a free data retrieval call binding the contract method 0xf0350799.
//
// Solidity: function load(bytes32 slot) view returns(bytes32)
func (_TestLib *TestLibCaller) Load(opts *bind.CallOpts, slot [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "load", slot)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Load is a free data retrieval call binding the contract method 0xf0350799.
//
// Solidity: function load(bytes32 slot) view returns(bytes32)
func (_TestLib *TestLibSession) Load(slot [32]byte) ([32]byte, error) {
	return _TestLib.Contract.Load(&_TestLib.CallOpts, slot)
}

// Load is a free data retrieval call binding the contract method 0xf0350799.
//
// Solidity: function load(bytes32 slot) view returns(bytes32)
func (_TestLib *TestLibCallerSession) Load(slot [32]byte) ([32]byte, error) {
	return _TestLib.Contract.Load(&_TestLib.CallOpts, slot)
}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibCaller) Max(opts *bind.CallOpts, a *big.Int, b *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "max", a, b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibSession) Max(a *big.Int, b *big.Int) (*big.Int, error) {
	return _TestLib.Contract.Max(&_TestLib.CallOpts, a, b)
}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibCallerSession) Max(a *big.Int, b *big.Int) (*big.Int, error) {
	return _TestLib.Contract.Max(&_TestLib.CallOpts, a, b)
}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibCaller) Min(opts *bind.CallOpts, a *big.Int, b *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "min", a, b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibSession) Min(a *big.Int, b *big.Int) (*big.Int, error) {
	return _TestLib.Contract.Min(&_TestLib.CallOpts, a, b)
}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_TestLib *TestLibCallerSession) Min(a *big.Int, b *big.Int) (*big.Int, error) {
	return _TestLib.Contract.Min(&_TestLib.CallOpts, a, b)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 hash, bytes signatureBytes) pure returns(address)
func (_TestLib *TestLibCaller) Recover(opts *bind.CallOpts, hash [32]byte, signatureBytes []byte) (common.Address, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "recover", hash, signatureBytes)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 hash, bytes signatureBytes) pure returns(address)
func (_TestLib *TestLibSession) Recover(hash [32]byte, signatureBytes []byte) (common.Address, error) {
	return _TestLib.Contract.Recover(&_TestLib.CallOpts, hash, signatureBytes)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 hash, bytes signatureBytes) pure returns(address)
func (_TestLib *TestLibCallerSession) Recover(hash [32]byte, signatureBytes []byte) (common.Address, error) {
	return _TestLib.Contract.Recover(&_TestLib.CallOpts, hash, signatureBytes)
}

// SetMargin is a free data retrieval call binding the contract method 0x80778a4f.
//
// Solidity: function setMargin((bool,bool,uint120,uint120) balance, (uint256,bool) newMargin) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) SetMargin(opts *bind.CallOpts, balance P1TypesBalance, newMargin SignedMathInt) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "setMargin", balance, newMargin)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// SetMargin is a free data retrieval call binding the contract method 0x80778a4f.
//
// Solidity: function setMargin((bool,bool,uint120,uint120) balance, (uint256,bool) newMargin) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) SetMargin(balance P1TypesBalance, newMargin SignedMathInt) (P1TypesBalance, error) {
	return _TestLib.Contract.SetMargin(&_TestLib.CallOpts, balance, newMargin)
}

// SetMargin is a free data retrieval call binding the contract method 0x80778a4f.
//
// Solidity: function setMargin((bool,bool,uint120,uint120) balance, (uint256,bool) newMargin) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) SetMargin(balance P1TypesBalance, newMargin SignedMathInt) (P1TypesBalance, error) {
	return _TestLib.Contract.SetMargin(&_TestLib.CallOpts, balance, newMargin)
}

// SetPosition is a free data retrieval call binding the contract method 0x0d261665.
//
// Solidity: function setPosition((bool,bool,uint120,uint120) balance, (uint256,bool) newPosition) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) SetPosition(opts *bind.CallOpts, balance P1TypesBalance, newPosition SignedMathInt) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "setPosition", balance, newPosition)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// SetPosition is a free data retrieval call binding the contract method 0x0d261665.
//
// Solidity: function setPosition((bool,bool,uint120,uint120) balance, (uint256,bool) newPosition) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) SetPosition(balance P1TypesBalance, newPosition SignedMathInt) (P1TypesBalance, error) {
	return _TestLib.Contract.SetPosition(&_TestLib.CallOpts, balance, newPosition)
}

// SetPosition is a free data retrieval call binding the contract method 0x0d261665.
//
// Solidity: function setPosition((bool,bool,uint120,uint120) balance, (uint256,bool) newPosition) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) SetPosition(balance P1TypesBalance, newPosition SignedMathInt) (P1TypesBalance, error) {
	return _TestLib.Contract.SetPosition(&_TestLib.CallOpts, balance, newPosition)
}

// SignedAdd is a free data retrieval call binding the contract method 0xb63f6580.
//
// Solidity: function signedAdd((uint256,bool) augend, (uint256,bool) addend) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) SignedAdd(opts *bind.CallOpts, augend SignedMathInt, addend SignedMathInt) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "signedAdd", augend, addend)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// SignedAdd is a free data retrieval call binding the contract method 0xb63f6580.
//
// Solidity: function signedAdd((uint256,bool) augend, (uint256,bool) addend) pure returns((uint256,bool))
func (_TestLib *TestLibSession) SignedAdd(augend SignedMathInt, addend SignedMathInt) (SignedMathInt, error) {
	return _TestLib.Contract.SignedAdd(&_TestLib.CallOpts, augend, addend)
}

// SignedAdd is a free data retrieval call binding the contract method 0xb63f6580.
//
// Solidity: function signedAdd((uint256,bool) augend, (uint256,bool) addend) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) SignedAdd(augend SignedMathInt, addend SignedMathInt) (SignedMathInt, error) {
	return _TestLib.Contract.SignedAdd(&_TestLib.CallOpts, augend, addend)
}

// SignedSub is a free data retrieval call binding the contract method 0x51894bf6.
//
// Solidity: function signedSub((uint256,bool) minuend, (uint256,bool) subtrahend) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) SignedSub(opts *bind.CallOpts, minuend SignedMathInt, subtrahend SignedMathInt) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "signedSub", minuend, subtrahend)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// SignedSub is a free data retrieval call binding the contract method 0x51894bf6.
//
// Solidity: function signedSub((uint256,bool) minuend, (uint256,bool) subtrahend) pure returns((uint256,bool))
func (_TestLib *TestLibSession) SignedSub(minuend SignedMathInt, subtrahend SignedMathInt) (SignedMathInt, error) {
	return _TestLib.Contract.SignedSub(&_TestLib.CallOpts, minuend, subtrahend)
}

// SignedSub is a free data retrieval call binding the contract method 0x51894bf6.
//
// Solidity: function signedSub((uint256,bool) minuend, (uint256,bool) subtrahend) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) SignedSub(minuend SignedMathInt, subtrahend SignedMathInt) (SignedMathInt, error) {
	return _TestLib.Contract.SignedSub(&_TestLib.CallOpts, minuend, subtrahend)
}

// Sub is a free data retrieval call binding the contract method 0xd165c800.
//
// Solidity: function sub((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibCaller) Sub(opts *bind.CallOpts, sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "sub", sint, value)

	if err != nil {
		return *new(SignedMathInt), err
	}

	out0 := *abi.ConvertType(out[0], new(SignedMathInt)).(*SignedMathInt)

	return out0, err

}

// Sub is a free data retrieval call binding the contract method 0xd165c800.
//
// Solidity: function sub((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibSession) Sub(sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	return _TestLib.Contract.Sub(&_TestLib.CallOpts, sint, value)
}

// Sub is a free data retrieval call binding the contract method 0xd165c800.
//
// Solidity: function sub((uint256,bool) sint, uint256 value) pure returns((uint256,bool))
func (_TestLib *TestLibCallerSession) Sub(sint SignedMathInt, value *big.Int) (SignedMathInt, error) {
	return _TestLib.Contract.Sub(&_TestLib.CallOpts, sint, value)
}

// SubFromMargin is a free data retrieval call binding the contract method 0x7dd1c962.
//
// Solidity: function subFromMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) SubFromMargin(opts *bind.CallOpts, balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "subFromMargin", balance, amount)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// SubFromMargin is a free data retrieval call binding the contract method 0x7dd1c962.
//
// Solidity: function subFromMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) SubFromMargin(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.SubFromMargin(&_TestLib.CallOpts, balance, amount)
}

// SubFromMargin is a free data retrieval call binding the contract method 0x7dd1c962.
//
// Solidity: function subFromMargin((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) SubFromMargin(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.SubFromMargin(&_TestLib.CallOpts, balance, amount)
}

// SubFromPosition is a free data retrieval call binding the contract method 0x437348e7.
//
// Solidity: function subFromPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCaller) SubFromPosition(opts *bind.CallOpts, balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "subFromPosition", balance, amount)

	if err != nil {
		return *new(P1TypesBalance), err
	}

	out0 := *abi.ConvertType(out[0], new(P1TypesBalance)).(*P1TypesBalance)

	return out0, err

}

// SubFromPosition is a free data retrieval call binding the contract method 0x437348e7.
//
// Solidity: function subFromPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibSession) SubFromPosition(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.SubFromPosition(&_TestLib.CallOpts, balance, amount)
}

// SubFromPosition is a free data retrieval call binding the contract method 0x437348e7.
//
// Solidity: function subFromPosition((bool,bool,uint120,uint120) balance, uint256 amount) pure returns((bool,bool,uint120,uint120))
func (_TestLib *TestLibCallerSession) SubFromPosition(balance P1TypesBalance, amount *big.Int) (P1TypesBalance, error) {
	return _TestLib.Contract.SubFromPosition(&_TestLib.CallOpts, balance, amount)
}

// That is a free data retrieval call binding the contract method 0x4c93be61.
//
// Solidity: function that(bool must, string requireReason, address addr) pure returns()
func (_TestLib *TestLibCaller) That(opts *bind.CallOpts, must bool, requireReason string, addr common.Address) error {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "that", must, requireReason, addr)

	if err != nil {
		return err
	}

	return err

}

// That is a free data retrieval call binding the contract method 0x4c93be61.
//
// Solidity: function that(bool must, string requireReason, address addr) pure returns()
func (_TestLib *TestLibSession) That(must bool, requireReason string, addr common.Address) error {
	return _TestLib.Contract.That(&_TestLib.CallOpts, must, requireReason, addr)
}

// That is a free data retrieval call binding the contract method 0x4c93be61.
//
// Solidity: function that(bool must, string requireReason, address addr) pure returns()
func (_TestLib *TestLibCallerSession) That(must bool, requireReason string, addr common.Address) error {
	return _TestLib.Contract.That(&_TestLib.CallOpts, must, requireReason, addr)
}

// ToUint120 is a free data retrieval call binding the contract method 0x1e4e4bad.
//
// Solidity: function toUint120(uint256 value) pure returns(uint120)
func (_TestLib *TestLibCaller) ToUint120(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "toUint120", value)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ToUint120 is a free data retrieval call binding the contract method 0x1e4e4bad.
//
// Solidity: function toUint120(uint256 value) pure returns(uint120)
func (_TestLib *TestLibSession) ToUint120(value *big.Int) (*big.Int, error) {
	return _TestLib.Contract.ToUint120(&_TestLib.CallOpts, value)
}

// ToUint120 is a free data retrieval call binding the contract method 0x1e4e4bad.
//
// Solidity: function toUint120(uint256 value) pure returns(uint120)
func (_TestLib *TestLibCallerSession) ToUint120(value *big.Int) (*big.Int, error) {
	return _TestLib.Contract.ToUint120(&_TestLib.CallOpts, value)
}

// ToUint128 is a free data retrieval call binding the contract method 0x809fdd33.
//
// Solidity: function toUint128(uint256 value) pure returns(uint128)
func (_TestLib *TestLibCaller) ToUint128(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "toUint128", value)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ToUint128 is a free data retrieval call binding the contract method 0x809fdd33.
//
// Solidity: function toUint128(uint256 value) pure returns(uint128)
func (_TestLib *TestLibSession) ToUint128(value *big.Int) (*big.Int, error) {
	return _TestLib.Contract.ToUint128(&_TestLib.CallOpts, value)
}

// ToUint128 is a free data retrieval call binding the contract method 0x809fdd33.
//
// Solidity: function toUint128(uint256 value) pure returns(uint128)
func (_TestLib *TestLibCallerSession) ToUint128(value *big.Int) (*big.Int, error) {
	return _TestLib.Contract.ToUint128(&_TestLib.CallOpts, value)
}

// ToUint32 is a free data retrieval call binding the contract method 0xc8193255.
//
// Solidity: function toUint32(uint256 value) pure returns(uint32)
func (_TestLib *TestLibCaller) ToUint32(opts *bind.CallOpts, value *big.Int) (uint32, error) {
	var out []interface{}
	err := _TestLib.contract.Call(opts, &out, "toUint32", value)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ToUint32 is a free data retrieval call binding the contract method 0xc8193255.
//
// Solidity: function toUint32(uint256 value) pure returns(uint32)
func (_TestLib *TestLibSession) ToUint32(value *big.Int) (uint32, error) {
	return _TestLib.Contract.ToUint32(&_TestLib.CallOpts, value)
}

// ToUint32 is a free data retrieval call binding the contract method 0xc8193255.
//
// Solidity: function toUint32(uint256 value) pure returns(uint32)
func (_TestLib *TestLibCallerSession) ToUint32(value *big.Int) (uint32, error) {
	return _TestLib.Contract.ToUint32(&_TestLib.CallOpts, value)
}

// NonReentrant1 is a paid mutator transaction binding the contract method 0x4437e178.
//
// Solidity: function nonReentrant1() returns(uint256)
func (_TestLib *TestLibTransactor) NonReentrant1(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLib.contract.Transact(opts, "nonReentrant1")
}

// NonReentrant1 is a paid mutator transaction binding the contract method 0x4437e178.
//
// Solidity: function nonReentrant1() returns(uint256)
func (_TestLib *TestLibSession) NonReentrant1() (*types.Transaction, error) {
	return _TestLib.Contract.NonReentrant1(&_TestLib.TransactOpts)
}

// NonReentrant1 is a paid mutator transaction binding the contract method 0x4437e178.
//
// Solidity: function nonReentrant1() returns(uint256)
func (_TestLib *TestLibTransactorSession) NonReentrant1() (*types.Transaction, error) {
	return _TestLib.Contract.NonReentrant1(&_TestLib.TransactOpts)
}

// NonReentrant2 is a paid mutator transaction binding the contract method 0x7bdfba4b.
//
// Solidity: function nonReentrant2() returns(uint256)
func (_TestLib *TestLibTransactor) NonReentrant2(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestLib.contract.Transact(opts, "nonReentrant2")
}

// NonReentrant2 is a paid mutator transaction binding the contract method 0x7bdfba4b.
//
// Solidity: function nonReentrant2() returns(uint256)
func (_TestLib *TestLibSession) NonReentrant2() (*types.Transaction, error) {
	return _TestLib.Contract.NonReentrant2(&_TestLib.TransactOpts)
}

// NonReentrant2 is a paid mutator transaction binding the contract method 0x7bdfba4b.
//
// Solidity: function nonReentrant2() returns(uint256)
func (_TestLib *TestLibTransactorSession) NonReentrant2() (*types.Transaction, error) {
	return _TestLib.Contract.NonReentrant2(&_TestLib.TransactOpts)
}

// Store is a paid mutator transaction binding the contract method 0x4000e4f6.
//
// Solidity: function store(bytes32 slot, bytes32 value) returns()
func (_TestLib *TestLibTransactor) Store(opts *bind.TransactOpts, slot [32]byte, value [32]byte) (*types.Transaction, error) {
	return _TestLib.contract.Transact(opts, "store", slot, value)
}

// Store is a paid mutator transaction binding the contract method 0x4000e4f6.
//
// Solidity: function store(bytes32 slot, bytes32 value) returns()
func (_TestLib *TestLibSession) Store(slot [32]byte, value [32]byte) (*types.Transaction, error) {
	return _TestLib.Contract.Store(&_TestLib.TransactOpts, slot, value)
}

// Store is a paid mutator transaction binding the contract method 0x4000e4f6.
//
// Solidity: function store(bytes32 slot, bytes32 value) returns()
func (_TestLib *TestLibTransactorSession) Store(slot [32]byte, value [32]byte) (*types.Transaction, error) {
	return _TestLib.Contract.Store(&_TestLib.TransactOpts, slot, value)
}
