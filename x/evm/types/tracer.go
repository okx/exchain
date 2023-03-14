package types

import (
	json2 "encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	json "github.com/json-iterator/go"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

type TraceConfig struct {
	// custom javascript tracer
	Tracer string `json:"tracer"`
	// disable stack capture
	DisableStack bool `json:"disableStack"`
	// disable storage capture
	DisableStorage bool `json:"disableStorage"`
	// print output during capture end
	Debug bool `json:"debug"`
	// enable memory capture
	DisableMemory bool `json:"disableMemory"`
	// enable return data capture
	DisableReturnData bool `json:"disableReturnData"`
}

func GetTracerResult(tracer tracers.Tracer, result *core.ExecutionResult) ([]byte, error) {
	var (
		res []byte
		err error
	)
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *logger.StructLogger:
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
	default:
		res, err = tracer.GetResult()
	}
	return res, err
}

// NoOpTracer is an empty implementation of vm.Tracer interface
type NoOpTracer struct{}

func (dt NoOpTracer) CaptureTxStart(gasLimit uint64) {
}

func (dt NoOpTracer) CaptureTxEnd(restGas uint64) {
}

func (dt NoOpTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

func (dt NoOpTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (dt NoOpTracer) GetResult() (json2.RawMessage, error) {
	return json2.RawMessage(`{}`), nil
}

func (dt NoOpTracer) Stop(err error) {

}

// NewNoOpTracer creates a no-op vm.Tracer
func NewNoOpTracer() tracers.Tracer {
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

// CaptureEnd implements vm.Tracer interface
func (dt NoOpTracer) CaptureEnd(
	output []byte,
	gasUsed uint64,
	t time.Duration,
	err error,
) {
}
func defaultTracerConfig() *TraceConfig {
	return &TraceConfig{
		Tracer:            "",
		DisableMemory:     false,
		DisableStorage:    false,
		DisableStack:      false,
		DisableReturnData: false,
		Debug:             false,
	}
}
func TestTracerConfig(traceConfig *TraceConfig) error {
	if traceConfig.Tracer != "" {
		_, err := tracers.New(traceConfig.Tracer, &tracers.Context{}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
func newTracer(ctx sdk.Context, txHash *common.Hash) (tracer tracers.Tracer) {
	if ctx.IsTraceTxLog() {
		var err error
		configBytes := ctx.TraceTxLogConfig()
		traceConfig := &TraceConfig{}
		if configBytes == nil {
			traceConfig = defaultTracerConfig()
		} else {
			err = json.Unmarshal(configBytes, traceConfig)
			if err != nil {
				return NewNoOpTracer()
			}
		}
		if traceConfig.Tracer == "" {
			//Basic tracer with config
			logConfig := logger.Config{
				EnableMemory:     !traceConfig.DisableMemory,
				DisableStorage:   traceConfig.DisableStorage,
				DisableStack:     traceConfig.DisableStack,
				EnableReturnData: !traceConfig.DisableReturnData,
				Debug:            traceConfig.Debug,
			}
			return logger.NewStructLogger(&logConfig)
		}
		// Json-based tracer
		tCtx := &tracers.Context{
			TxHash: *txHash,
		}
		tracer, err = tracers.New(traceConfig.Tracer, tCtx, nil)
		if err != nil {
			return NewNoOpTracer()
		}
		return tracer
	} else {
		//no op tracer
		return NewNoOpTracer()
	}
}
