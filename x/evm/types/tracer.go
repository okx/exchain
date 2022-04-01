package types

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	json "github.com/json-iterator/go"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type TraceConfig struct {
	// custom javascript tracer
	Tracer string `json:"tracer,omitempty"`
	// disable stack capture
	DisableStack bool `json:"disableStack"`
	// disable storage capture
	DisableStorage bool `json:"disableStorage"`
	// print output during capture end
	Debug bool `json:"debug,omitempty"`
	// enable memory capture
	DisableMemory bool `json:"disableMemory"`
	// enable return data capture
	DisableReturnData bool `json:"disableReturnData"`
}

func GetTracerResult(tracer vm.Tracer, result *core.ExecutionResult) ([]byte, error) {
	var (
		res []byte
		err error
	)
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		// If the result contains a revert reason, return it.
		returnVal := fmt.Sprintf("%x", result.Return())
		if len(result.Revert()) > 0 {
			returnVal = fmt.Sprintf("%x", result.Revert())
		}
		res, err = json.ConfigFastest.Marshal(&TraceExecutionResult{
			Gas:         result.UsedGas,
			Failed:      result.Failed(),
			ReturnValue: returnVal,
			StructLogs:  FormatLogs(tracer.StructLogs()),
		})
	case *tracers.Tracer:
		res, err = tracer.GetResult()
	default:
		res = []byte(fmt.Sprintf("bad tracer type %T", tracer))
	}
	return res, err
}

// NoOpTracer is an empty implementation of vm.Tracer interface
type NoOpTracer struct{}

// NewNoOpTracer creates a no-op vm.Tracer
func NewNoOpTracer() *NoOpTracer {
	return &NoOpTracer{}
}

// CaptureStart implements vm.Tracer interface
func (dt NoOpTracer) CaptureStart(
	env *vm.EVM,
	from, to common.Address,
	create bool,
	input []byte,
	gas uint64,
	value *big.Int,
) {
}

// CaptureEnter implements vm.Tracer interface
func (dt NoOpTracer) CaptureEnter(
	typ vm.OpCode,
	from common.Address,
	to common.Address,
	input []byte,
	gas uint64,
	value *big.Int,
) {
}

// CaptureExit implements vm.Tracer interface
func (dt NoOpTracer) CaptureExit(output []byte, gasUsed uint64, err error) {}

// CaptureState implements vm.Tracer interface
func (dt NoOpTracer) CaptureState(
	env *vm.EVM,
	pc uint64,
	op vm.OpCode,
	gas, cost uint64,
	scope *vm.ScopeContext,
	rData []byte,
	depth int,
	err error,
) {
}

// CaptureFault implements vm.Tracer interface
func (dt NoOpTracer) CaptureFault(
	env *vm.EVM,
	pc uint64,
	op vm.OpCode,
	gas, cost uint64,
	scope *vm.ScopeContext,
	depth int,
	err error,
) {
}

// CaptureEnd implements vm.Tracer interface
func (dt NoOpTracer) CaptureEnd(
	output []byte,
	gasUsed uint64,
	t time.Duration,
	err error,
) {
}

func newTracer(ctx sdk.Context, txHash *common.Hash) (tracer vm.Tracer, err error) {
	if ctx.IsTraceTxLog() {
		configBytes := ctx.TraceTxLogConfig()
		var traceConfig TraceConfig
		err = json.Unmarshal(configBytes, &traceConfig)
		if err != nil {
			return nil, err
		}
		if traceConfig.Tracer == "" {
			//Basic tracer with config
			logConfig := vm.LogConfig{
				DisableMemory:     traceConfig.DisableMemory,
				DisableStorage:    traceConfig.DisableStorage,
				DisableStack:      traceConfig.DisableStack,
				DisableReturnData: traceConfig.DisableReturnData,
				Debug:             traceConfig.Debug,
			}
			tracer = vm.NewStructLogger(&logConfig)
			return tracer, nil
		}
		// Json-based tracer
		tCtx := &tracers.Context{
			TxHash: *txHash,
		}
		return tracers.New(traceConfig.Tracer, tCtx)
	} else {
		//no op tracer
		tracer = NewNoOpTracer()
		return tracer, nil
	}
}
