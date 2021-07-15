package ewasm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/vm"

	evmc "github.com/ethereum/evmc/v8/bindings/go/evmc"
)

// EVMC represents the reference to a common EVMC-based VM instance and
// the current execution context as required by go-ethereum design.
type EVMC struct {
	instance *evmc.VM        // The reference to the EVMC VM instance.
	env      *vm.EVM         // The execution context.
	cap      evmc.Capability // The supported EVMC capability (EVM or Ewasm)
	readOnly bool            // The readOnly flag (TODO: Try to get rid of it).
	depth    int
}

func NewEWasm() *evmc.VM {
	return initEVMC(evmc.CapabilityEWASM, "/Users/oker/work/eth/hera/src/libhera.dylib")
}

func initEVMC(cap evmc.Capability, config string) *evmc.VM {
	options := strings.Split(config, ",")
	path := options[0]

	if path == "" {
		panic("EVMC VM path not provided, set --vm.(evm|ewasm)=/path/to/vm")
	}

	instance, err := evmc.Load(path)

	if err != nil {
		panic(err.Error())
	}

	// Set options before checking capabilities.
	for _, option := range options[1:] {
		if idx := strings.Index(option, "="); idx >= 0 {
			name := option[:idx]
			value := option[idx+1:]
			err := instance.SetOption(name, value)
			if err == nil {
			} else {
			}
		}
	}

	if !instance.HasCapability(cap) {
		panic(fmt.Errorf("The EVMC module %s does not have requested capability %d", path, cap))
	}
	return instance
}

// hostContext implements evmc.HostContext interface.
type hostContext struct {
	env      *vm.EVM      // The reference to the EVM execution context.
	contract *vm.Contract // The reference to the current contract, needed by Call-like methods.
}

func (host *hostContext) AccessAccount(addr evmc.Address) evmc.AccessStatus {
	panic("implement me")
}

func (host *hostContext) AccessStorage(addr evmc.Address, key evmc.Hash) evmc.AccessStatus {
	panic("implement me")
}

func (host *hostContext) AccountExists(addr evmc.Address) bool {
	if host.env.ChainConfig().IsEIP158(host.env.Context.BlockNumber) {
		if !host.env.StateDB.Empty(common.Address(addr)) {
			return true
		}
	} else if host.env.StateDB.Exist(common.Address(addr)) {
		return true
	}
	return false
}

func (host *hostContext) GetStorage(addr evmc.Address, key evmc.Hash) evmc.Hash {
	return evmc.Hash(host.env.StateDB.GetState(common.Address(addr), common.Hash(key)))
}

func (host *hostContext) SetStorage(addr evmc.Address, key evmc.Hash, value evmc.Hash) (status evmc.StorageStatus) {
	oldValue := host.env.StateDB.GetState(common.Address(addr), common.Hash(key))
	if evmc.Hash(oldValue) == value {
		return evmc.StorageUnchanged
	}

	current := host.env.StateDB.GetState(common.Address(addr), common.Hash(key))
	original := host.env.StateDB.GetCommittedState(common.Address(addr), common.Hash(key))

	host.env.StateDB.SetState(common.Address(addr), common.Hash(key), common.Hash(value))

	hasNetStorageCostEIP := host.env.ChainConfig().IsConstantinople(host.env.Context.BlockNumber) &&
		!host.env.ChainConfig().IsPetersburg(host.env.Context.BlockNumber)
	if !hasNetStorageCostEIP {

		zero := evmc.Hash{}
		status = evmc.StorageModified
		if evmc.Hash(oldValue) == zero {
			return evmc.StorageAdded
		} else if value == zero {
			host.env.StateDB.AddRefund(params.SstoreRefundGas)
			return evmc.StorageDeleted
		}
		return evmc.StorageModified
	}

	if original == current {
		if original == (common.Hash{}) { // create slot (2.1.1)
			return evmc.StorageAdded
		}
		if value == (evmc.Hash{}) { // delete slot (2.1.2b)
			host.env.StateDB.AddRefund(params.NetSstoreClearRefund)
			return evmc.StorageDeleted
		}
		return evmc.StorageModified
	}
	if original != (common.Hash{}) {
		if current == (common.Hash{}) { // recreate slot (2.2.1.1)
			host.env.StateDB.SubRefund(params.NetSstoreClearRefund)
		} else if value == (evmc.Hash{}) { // delete slot (2.2.1.2)
			host.env.StateDB.AddRefund(params.NetSstoreClearRefund)
		}
	}
	if original == common.Hash(value) {
		if original == (common.Hash{}) { // reset to original inexistent slot (2.2.2.1)
			host.env.StateDB.AddRefund(params.NetSstoreResetClearRefund)
		} else { // reset to original existing slot (2.2.2.2)
			host.env.StateDB.AddRefund(params.NetSstoreResetRefund)
		}
	}
	return evmc.StorageModifiedAgain
}

func (host *hostContext) GetBalance(addr evmc.Address) evmc.Hash {
	return evmc.Hash(common.BigToHash(host.env.StateDB.GetBalance(common.Address(addr))))
}

func (host *hostContext) GetCodeSize(addr evmc.Address) int {
	return host.env.StateDB.GetCodeSize(common.Address(addr))
}

func (host *hostContext) GetCodeHash(addr evmc.Address) evmc.Hash {
	if host.env.StateDB.Empty(common.Address(addr)) {
		return evmc.Hash{}
	}
	return evmc.Hash(host.env.StateDB.GetCodeHash(common.Address(addr)))
}

func (host *hostContext) GetCode(addr evmc.Address) []byte {
	return host.env.StateDB.GetCode(common.Address(addr))
}

func (host *hostContext) Selfdestruct(addr evmc.Address, beneficiary evmc.Address) {
	db := host.env.StateDB
	if !db.HasSuicided(common.Address(addr)) {
		db.AddRefund(params.SelfdestructRefundGas)
	}
	db.AddBalance(common.Address(beneficiary), db.GetBalance(common.Address(addr)))
	db.Suicide(common.Address(addr))
}

func SetBytesWithHash(to evmc.Hash, from []byte) {
	if len(from) > len(to) {
		from = from[len(from)-common.HashLength:]
	}

	copy(to[common.HashLength-len(from):], from)
}

func SetBytesWithAddress(to evmc.Address, from []byte) {
	if len(from) > len(to) {
		from = from[len(from)-common.AddressLength:]
	}

	copy(to[common.AddressLength-len(from):], from)
}

func (host *hostContext) GetTxContext() evmc.TxContext {

	return evmc.TxContext{

		GasPrice:   evmc.Hash(common.BigToHash(host.env.GasPrice)),
		Origin:     evmc.Address(host.env.Origin),
		Coinbase:   evmc.Address(host.env.Context.Coinbase),
		Number:     host.env.Context.BlockNumber.Int64(),
		Timestamp:  host.env.Context.Time.Int64(),
		GasLimit:   int64(host.env.Context.GasLimit),
		Difficulty: evmc.Hash(common.BigToHash(host.env.Context.Difficulty))}
}

func (host *hostContext) GetBlockHash(number int64) evmc.Hash {
	b := host.env.Context.BlockNumber.Int64()
	if number >= (b-256) && number < b {
		return evmc.Hash(host.env.Context.GetHash(uint64(number)))
	}
	return evmc.Hash{}
}

func (host *hostContext) EmitLog(addr evmc.Address, topics []evmc.Hash, data []byte) {
	var tmpTopic []common.Hash
	if len(topics) > 0 {
		for _, v := range topics {
			tmpTopic = append(tmpTopic, common.Hash(v))
		}
	}

	host.env.StateDB.AddLog(&types.Log{
		Address:     common.Address(addr),
		Topics:      tmpTopic,
		Data:        data,
		BlockNumber: host.env.Context.BlockNumber.Uint64(),
	})
}

func (host *hostContext) Call(kind evmc.CallKind,
	destination evmc.Address, sender evmc.Address, value evmc.Hash, input []byte, gas int64, depth int,
	static bool, salt evmc.Hash) (output []byte, gasLeft int64, createAddr evmc.Address, err error) {

	gasU := uint64(gas)
	var gasLeftU uint64
	valueInBig := common.Hash(value).Big()

	switch kind {
	case evmc.Call:
		if static {
			output, gasLeftU, err = host.env.StaticCall(host.contract, common.Address(destination), input, gasU)
		} else {
			output, gasLeftU, err = host.env.Call(host.contract, common.Address(destination), input, gasU, valueInBig)
		}
	case evmc.DelegateCall:
		output, gasLeftU, err = host.env.DelegateCall(host.contract, common.Address(destination), input, gasU)
	case evmc.CallCode:
		output, gasLeftU, err = host.env.CallCode(host.contract, common.Address(destination), input, gasU, valueInBig)
	case evmc.Create:
		var createOutput []byte
		var tmpAddr common.Address
		createOutput, tmpAddr, gasLeftU, err = host.env.Create(host.contract, input, gasU, valueInBig)
		createAddr = evmc.Address(tmpAddr)
		isHomestead := host.env.ChainConfig().IsHomestead(host.env.Context.BlockNumber)
		if !isHomestead && err == vm.ErrCodeStoreOutOfGas {
			err = nil
		}
		if err == vm.ErrExecutionReverted {
			// Assign return buffer from REVERT.
			// TODO: Bad API design: return data buffer and the code is returned in the same place. In worst case
			//       the code is returned also when there is not enough funds to deploy the code.
			output = createOutput
		}
	case evmc.Create2:
		var createOutput []byte
		var tmpSalt uint256.Int
		tmpSalt.SetFromBig(common.Hash(salt).Big())
		var tmpAddr common.Address
		createOutput, tmpAddr, gasLeftU, err = host.env.Create2(host.contract, input, gasU, valueInBig, &tmpSalt)
		createAddr = evmc.Address(tmpAddr)
		if err == vm.ErrExecutionReverted {
			// Assign return buffer from REVERT.
			// TODO: Bad API design: return data buffer and the code is returned in the same place. In worst case
			//       the code is returned also when there is not enough funds to deploy the code.
			output = createOutput
		}
	default:
		panic(fmt.Errorf("EVMC: Unknown call kind %d", kind))
	}

	// Map errors.
	if err == vm.ErrExecutionReverted {
		err = evmc.Revert
	} else if err != nil {
		err = evmc.Failure
	}

	gasLeft = int64(gasLeftU)
	return output, gasLeft, createAddr, err
}

// getRevision translates ChainConfig's HF block information into EVMC revision.
func getRevision(env *vm.EVM) evmc.Revision {
	n := env.Context.BlockNumber
	conf := env.ChainConfig()
	switch {
	case conf.IsPetersburg(n):
		return evmc.Petersburg
	case conf.IsConstantinople(n):
		return evmc.Constantinople
	case conf.IsByzantium(n):
		return evmc.Byzantium
	case conf.IsEIP158(n):
		return evmc.SpuriousDragon
	case conf.IsEIP150(n):
		return evmc.TangerineWhistle
	case conf.IsHomestead(n):
		return evmc.Homestead
	default:
		return evmc.Frontier
	}
}

// Run implements Interpreter.Run().
func (evm *EVMC) Run(contract *vm.Contract, input []byte, readOnly bool) (ret []byte, err error) {
	evm.depth++
	defer func() { evm.depth-- }()

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}

	kind := evmc.Call
	if evm.env.StateDB.GetCodeSize(contract.Address()) == 0 {
		// Guess if this is a CREATE.
		kind = evmc.Create
	}

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	if readOnly && !evm.readOnly {
		evm.readOnly = true
		defer func() { evm.readOnly = false }()
	}

	output, gasLeft, err := evm.instance.Execute(
		&hostContext{evm.env, contract},
		getRevision(evm.env),
		kind,
		evm.readOnly,
		evm.depth-1,
		int64(contract.Gas),
		evmc.Address(contract.Address()),
		evmc.Address(contract.Caller()),
		input,
		evmc.Hash(common.BigToHash(contract.Value())),
		contract.Code,
		evmc.Hash{})

	contract.Gas = uint64(gasLeft)

	if err == evmc.Revert {
		err = vm.ErrExecutionReverted
	} else if evmcError, ok := err.(evmc.Error); ok && evmcError.IsInternalError() {
		panic(fmt.Sprintf("EVMC VM internal error: %s", evmcError.Error()))
	}

	return output, err
}

// CanRun implements Interpreter.CanRun().
func (evm *EVMC) CanRun(code []byte) bool {
	required := evmc.CapabilityEVM1
	wasmPreamble := []byte("\x00asm")
	if bytes.HasPrefix(code, wasmPreamble) {
		required = evmc.CapabilityEWASM
	}
	return evm.cap == required
}
