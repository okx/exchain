package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

// StateTransition defines data to transitionDB in evm
type StateTransition struct {
	// TxData fields
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	Recipient    *common.Address
	Amount       *big.Int
	Payload      []byte

	ChainID  *big.Int
	Csdb     *CommitStateDB
	TxHash   *common.Hash
	Sender   common.Address
	Simulate bool // i.e CheckTx execution

	CoinDenom string
	GasReturn uint64
}

// GasInfo returns the gas limit, gas consumed and gas refunded from the EVM transition
// execution
type GasInfo struct {
	GasLimit    uint64
	GasConsumed uint64
	GasRefunded uint64
}

// ExecutionResult represents what's returned from a transition
type ExecutionResult struct {
	Logs    []*ethtypes.Log
	Bloom   *big.Int
	Result  *sdk.Result
	GasInfo GasInfo
}

// GetHashFn implements vm.GetHashFunc for Ethermint. It handles 3 cases:
//  1. The requested height matches the current height (and thus same epoch number)
//  2. The requested height is from an previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func GetHashFn(ctx sdk.Context, csdb *CommitStateDB) vm.GetHashFunc {
	return func(height uint64) common.Hash {
		switch {
		case ctx.BlockHeight() == int64(height):
			// Case 1: The requested height matches the one from the context so we can retrieve the header
			// hash directly from the context.
			return csdb.bhash

		case ctx.BlockHeight() > int64(height):
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			return csdb.WithContext(ctx).GetHeightHash(height)

		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return common.Hash{}
		}
	}
}

func (st StateTransition) newEVM(
	ctx sdk.Context,
	csdb *CommitStateDB,
	gasLimit uint64,
	gasPrice *big.Int,
	config ChainConfig,
	vmConfig vm.Config,
) *vm.EVM {
	// Create context for evm
	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     GetHashFn(ctx, csdb),
		Coinbase:    common.BytesToAddress(ctx.BlockHeader().ProposerAddress),
		BlockNumber: big.NewInt(ctx.BlockHeight()),
		Time:        big.NewInt(ctx.BlockHeader().Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
		GasLimit:    gasLimit,
	}

	txCtx := vm.TxContext{
		Origin:   st.Sender,
		GasPrice: gasPrice,
	}

	return vm.NewEVM(blockCtx, txCtx, csdb, config.EthereumConfig(st.ChainID), vmConfig)
}

// TransitionDb will transition the state by applying the current transaction and
// returning the evm execution result.
// NOTE: State transition checks are run during AnteHandler execution.
func (st StateTransition) TransitionDb(ctx sdk.Context, config ChainConfig) (*ExecutionResult, error) {
	contractCreation := st.Recipient == nil

	cost, err := core.IntrinsicGas(st.Payload, contractCreation, config.IsHomestead(), config.IsIstanbul())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid intrinsic gas for transaction")
	}

	// This gas limit the the transaction gas limit with intrinsic gas subtracted
	gasLimit := st.GasLimit - ctx.GasMeter().GasConsumed()

	csdb := st.Csdb.WithContext(ctx)
	if st.Simulate {
		// gasLimit is set here because stdTxs incur gaskv charges in the ante handler, but for eth_call
		// the cost needs to be the same as an Ethereum transaction sent through the web3 API
		consumedGas := ctx.GasMeter().GasConsumed()
		gasLimit = st.GasLimit - cost
		if consumedGas < cost {
			// If Cosmos standard tx ante handler cost is less than EVM intrinsic cost
			// gas must be consumed to match to accurately simulate an Ethereum transaction
			ctx.GasMeter().ConsumeGas(cost-consumedGas, "Intrinsic gas match")
		}
	}

	// This gas meter is set up to consume gas from gaskv during evm execution and be ignored
	currentGasMeter := ctx.GasMeter()
	evmGasMeter := sdk.NewInfiniteGasMeter()
	csdb.WithContext(ctx.WithGasMeter(evmGasMeter))

	params := csdb.GetParams()

	var tracer vm.Tracer
	tracer = vm.NewStructLogger(nil)

	vmConfig := vm.Config{
		ExtraEips: params.ExtraEIPs,
		Debug:     true,
		Tracer:    tracer,
	}

	evm := st.newEVM(ctx, csdb, gasLimit, st.Price, config, vmConfig)

	var (
		ret             []byte
		leftOverGas     uint64
		contractAddress common.Address
		recipientLog    string
		senderRef       = vm.AccountRef(st.Sender)
	)

	// Get nonce of account outside of the EVM
	currentNonce := csdb.GetNonce(st.Sender)
	// Set nonce of sender account before evm state transition for usage in generating Create address
	csdb.SetNonce(st.Sender, st.AccountNonce)

	// create contract or execute call
	switch contractCreation {
	case true:
		if !params.EnableCreate {
			return nil, ErrCreateDisabled
		}

		ret, contractAddress, leftOverGas, err = evm.Create(senderRef, st.Payload, gasLimit, st.Amount)
		recipientLog = fmt.Sprintf("contract address %s", contractAddress.String())
	default:
		if !params.EnableCall {
			return nil, ErrCallDisabled
		}

		// Increment the nonce for the next transaction	(just for evm state transition)
		csdb.SetNonce(st.Sender, csdb.GetNonce(st.Sender)+1)
		ret, leftOverGas, err = evm.Call(senderRef, *st.Recipient, st.Payload, gasLimit, st.Amount)
		recipientLog = fmt.Sprintf("recipient address %s", st.Recipient.String())
	}

	gasConsumed := gasLimit - leftOverGas

	result := &core.ExecutionResult{
		UsedGas:    gasConsumed,
		Err:        err,
		ReturnData: ret,
	}

	// The maximum refund should not be more than half of the consumption
	gasReturn := gasConsumed / 2
	if gasReturn > csdb.refund {
		gasReturn = csdb.refund
	}
	st.GasReturn = gasReturn

	if err != nil {
		// Consume gas before returning
		ctx.GasMeter().ConsumeGas(gasConsumed, "evm execution consumption")
		if !st.Simulate {
			saveTraceResult(ctx, tracer, result)
		}
		return nil, newRevertError(ret, err)
	}

	// Resets nonce to value pre state transition
	csdb.SetNonce(st.Sender, currentNonce)

	// Generate bloom filter to be saved in tx receipt data
	bloomInt := big.NewInt(0)

	var (
		bloomFilter ethtypes.Bloom
		logs        []*ethtypes.Log
	)

	if st.TxHash != nil && !st.Simulate {
		logs, err = csdb.GetLogs(*st.TxHash)
		if err != nil {
			return nil, err
		}

		bloomInt = big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs))
		bloomFilter = ethtypes.BytesToBloom(bloomInt.Bytes())
	}

	if !st.Simulate {
		// Finalise state if not a simulated transaction
		// TODO: change to depend on config
		if err := csdb.Finalise(true); err != nil {
			return nil, err
		}

		if _, err = csdb.Commit(true); err != nil {
			return nil, err
		}
		saveTraceResult(ctx, tracer, result)
	}

	// Encode all necessary data into slice of bytes to return in sdk result
	resultData := ResultData{
		Bloom:  bloomFilter,
		Logs:   logs,
		Ret:    ret,
		TxHash: *st.TxHash,
	}

	if contractCreation {
		resultData.ContractAddress = contractAddress
	}

	resBz, err := EncodeResultData(resultData)
	if err != nil {
		return nil, err
	}

	resultLog := fmt.Sprintf(
		"executed EVM state transition; sender address %s; %s", st.Sender.String(), recipientLog,
	)

	executionResult := &ExecutionResult{
		Logs:  logs,
		Bloom: bloomInt,
		Result: &sdk.Result{
			Data: resBz,
			Log:  resultLog,
		},
		GasInfo: GasInfo{
			GasConsumed: gasConsumed,
			GasLimit:    gasLimit,
			GasRefunded: leftOverGas,
		},
	}

	// Consume gas from evm execution
	// Out of gas check does not need to be done here since it is done within the EVM execution
	ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(gasConsumed, "EVM execution consumption")

	return executionResult, nil
}

func (st StateTransition) RefundGas(ctx sdk.Context) error {

	gasRemaining := st.GasLimit - ctx.GasMeter().GasConsumed() + st.GasReturn
	feeReturn := big.NewInt(1).Mul(st.Price, big.NewInt(1).SetUint64(gasRemaining))

	if feeReturn.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	senderAddress, err := sdk.AccAddressFromHex(strings.TrimPrefix(st.Sender.Hex(), "0x"))
	if err != nil {
		return err
	}

	if err = st.Csdb.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx.WithGasMeter(sdk.NewInfiniteGasMeter()),
		types.FeeCollectorName,
		senderAddress,
		sdk.NewCoins(sdk.NewCoin(st.CoinDenom, sdk.NewDecFromBigIntWithPrec(feeReturn, sdk.Precision))),
	); err != nil {
		return err
	}

	return nil
}

func newRevertError(data []byte, e error) error {
	var resultError []string
	if data == nil || e.Error() != vm.ErrExecutionReverted.Error() {
		return e
	}
	resultError = append(resultError, e.Error())
	reason, errUnpack := abi.UnpackRevert(data)
	if errUnpack == nil {
		resultError = append(resultError, vm.ErrExecutionReverted.Error()+":"+reason)
	} else {
		resultError = append(resultError, hexutil.Encode(data))
	}
	resultError = append(resultError, ErrorHexData)
	resultError = append(resultError, hexutil.Encode(data))
	ret, error := json.Marshal(resultError)

	//failed to marshal, return original data in error
	if error != nil {
		return fmt.Errorf(e.Error()+"[%v]", hexutil.Encode(data))
	}
	return errors.New(string(ret))
}
