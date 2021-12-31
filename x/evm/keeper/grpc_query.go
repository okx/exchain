package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/eth/tracers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/okex/exchain/x/evm/types"
)

const (
	defaultTraceTimeout = 5 * time.Second
)

// TraceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (k Keeper) TraceTx(c context.Context, req *types.QueryTraceTxRequest) (*types.QueryTraceTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.TraceConfig != nil && req.TraceConfig.Limit < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "output limit cannot be negative, got %d", req.TraceConfig.Limit)
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = ctx.WithBlockHeight(req.BlockNumber)
	ctx = ctx.WithBlockTime(req.BlockTime)
	ctx = ctx.WithHeaderHash(common.Hex2Bytes(req.BlockHash))
	k.WithContext(ctx)

	params := k.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(k.eip155ChainID)
	signer := ethtypes.MakeSigner(ethCfg, big.NewInt(ctx.BlockHeight()))
	baseFee := k.feeMarketKeeper.GetBaseFee(ctx)

	for i, tx := range req.Predecessors {
		ethTx := tx.AsTransaction()
		msg, err := ethTx.AsMessage(signer, baseFee)
		if err != nil {
			continue
		}
		k.SetTxHashTransient(ethTx.Hash())
		k.SetTxIndexTransient(uint64(i))

		if _, err := k.ApplyMessage(msg, types.NewNoOpTracer(), true); err != nil {
			continue
		}
	}

	tx := req.Msg.AsTransaction()
	result, err := k.traceTx(ctx, signer, req.TxIndex, ethCfg, tx, baseFee, req.TraceConfig, false)
	if err != nil {
		// error will be returned with detail status from traceTx
		return nil, err
	}

	resultData, err := json.Marshal(result)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTraceTxResponse{
		Data: resultData,
	}, nil
}

func (k *Keeper) traceTx(
	ctx sdk.Context,
	signer ethtypes.Signer,
	txIndex uint64,
	ethCfg *ethparams.ChainConfig,
	tx *ethtypes.Transaction,
	baseFee *big.Int,
	traceConfig *types.TraceConfig,
	commitMessage bool,
) (*interface{}, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer    vm.Tracer
		overrides *ethparams.ChainConfig
		err       error
	)

	msg, err := tx.AsMessage(signer, baseFee)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	txHash := tx.Hash()

	if traceConfig != nil && traceConfig.Overrides != nil {
		overrides = traceConfig.Overrides.EthereumConfig(ethCfg.ChainID)
	}

	switch {
	case traceConfig != nil && traceConfig.Tracer != "":
		timeout := defaultTraceTimeout
		// TODO: change timeout to time.duration
		// Used string to comply with go ethereum
		if traceConfig.Timeout != "" {
			timeout, err = time.ParseDuration(traceConfig.Timeout)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "timeout value: %s", err.Error())
			}
		}

		tCtx := &tracers.Context{
			BlockHash: k.GetHashFn()(uint64(ctx.BlockHeight())),
			TxIndex:   int(txIndex),
			TxHash:    txHash,
		}

		// Construct the JavaScript tracer to execute with
		if tracer, err = tracers.New(traceConfig.Tracer, tCtx); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Handle timeouts and RPC cancellations
		deadlineCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
		defer cancel()

		go func() {
			<-deadlineCtx.Done()
			if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
				tracer.(*tracers.Tracer).Stop(errors.New("execution timeout"))
			}
		}()

	case traceConfig != nil:
		logConfig := vm.LogConfig{
			EnableMemory:     traceConfig.EnableMemory,
			DisableStorage:   traceConfig.DisableStorage,
			DisableStack:     traceConfig.DisableStack,
			EnableReturnData: traceConfig.EnableReturnData,
			Debug:            traceConfig.Debug,
			Limit:            int(traceConfig.Limit),
			Overrides:        overrides,
		}
		tracer = vm.NewStructLogger(&logConfig)
	default:
		tracer = types.NewTracer(types.TracerStruct, msg, ethCfg, ctx.BlockHeight())
	}

	k.SetTxHashTransient(txHash)
	k.SetTxIndexTransient(txIndex)

	res, err := k.ApplyMessage(msg, tracer, commitMessage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result interface{}

	// Depending on the tracer type, format and return the trace result data.
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		// TODO: Return proper returnValue
		result = types.ExecutionResult{
			Gas:         res.GasUsed,
			Failed:      res.Failed(),
			ReturnValue: "",
			StructLogs:  types.FormatLogs(tracer.StructLogs()),
		}
	case *tracers.Tracer:
		result, err = tracer.GetResult()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid tracer type %T", tracer)
	}

	return &result, nil
}
