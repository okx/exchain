package types

import (
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	json "github.com/json-iterator/go"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	tracesDir = "traces"

	FlagEnableTraces           = "evm-trace-enable"
	FlagTraceSegment           = "evm-trace-segment"
	FlagTraceFromAddrs         = "evm-trace-from-addrs"
	FlagTraceToAddrs           = "evm-trace-to-addrs"
	FlagTraceDisableMemory     = "evm-trace-nomemory"
	FlagTraceDisableStack      = "evm-trace-nostack"
	FlagTraceDisableStorage    = "evm-trace-nostorage"
	FlagTraceDisableReturnData = "evm-trace-noreturndata"
	FlagTraceDebug             = "evm-trace-debug"
)

var (
	tracesDB     dbm.DB
	enableTraces bool

	// trace from/to addr
	traceFromAddrs, traceToAddrs map[string]struct{}

	// trace segment
	step, total, num int64

	evmLogConfig *vm.LogConfig
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

func CloseTracer() {
	if tracesDB != nil {
		tracesDB.Close()
	}
}

func InitTxTraces() {
	enableTraces = viper.GetBool(FlagEnableTraces)
	if !enableTraces {
		return
	}

	etp := viper.GetString(FlagTraceSegment)
	segment := strings.Split(etp, "-")
	if len(segment) != 3 {
		panic(fmt.Errorf("invalid evm trace params: %s", etp))
	}

	var err error
	step, err = strconv.ParseInt(segment[0], 10, 64)
	if err != nil || step <= 0 {
		panic(fmt.Errorf("invalid evm trace params: %s", etp))
	}
	total, err = strconv.ParseInt(segment[1], 10, 64)
	if err != nil || total <= 0 {
		panic(fmt.Errorf("invalid evm trace params: %s", etp))
	}
	num, err = strconv.ParseInt(segment[2], 10, 64)
	if err != nil || num < 0 || num >= total {
		panic(fmt.Errorf("invalid evm trace params: %s", etp))
	}

	traceFromAddrs = make(map[string]struct{})
	traceToAddrs = make(map[string]struct{})
	fromAddrsStr := viper.GetString(FlagTraceFromAddrs)
	if fromAddrsStr != "" {
		for _, addr := range strings.Split(fromAddrsStr, ",") {
			traceFromAddrs[common.HexToAddress(addr).String()] = struct{}{}
		}
	}
	toAddrsStr := viper.GetString(FlagTraceToAddrs)
	if toAddrsStr != "" {
		for _, addr := range strings.Split(toAddrsStr, ",") {
			traceToAddrs[common.HexToAddress(addr).String()] = struct{}{}
		}
	}

	evmLogConfig = &vm.LogConfig{
		DisableMemory:     viper.GetBool(FlagTraceDisableMemory),
		DisableStack:      viper.GetBool(FlagTraceDisableStack),
		DisableStorage:    viper.GetBool(FlagTraceDisableStorage),
		DisableReturnData: viper.GetBool(FlagTraceDisableReturnData),
		Debug:             viper.GetBool(FlagTraceDebug),
	}

	dataDir := filepath.Join(viper.GetString("home"), "data")
	tracesDB, err = sdk.NewLevelDB(tracesDir, dataDir)
	if err != nil {
		panic(err)
	}
}

func checkTracesSegment(height int64, from, to string) bool {
	_, fromOk := traceFromAddrs[from]
	_, toOk := traceToAddrs[to]

	return enableTraces &&
		((height-1)/step)%total == num &&
		(len(traceFromAddrs) == 0 || (len(traceFromAddrs) > 0 && fromOk)) &&
		(len(traceToAddrs) == 0 || to == "" || (len(traceToAddrs) > 0 && toOk))
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
func saveTraceResult(ctx sdk.Context, tracer vm.Tracer, result *core.ExecutionResult) {

	res, err := GetTracerResult(tracer, result)
	if err != nil {
		res = []byte(err.Error())
	}
	saveToDB(tmtypes.Tx(ctx.TxBytes()).Hash(ctx.BlockHeight()), res)
}

func saveToDB(txHash []byte, value json.RawMessage) {
	if tracesDB == nil {
		panic("traces db is nil")
	}
	err := tracesDB.SetSync(txHash, value)
	if err != nil {
		panic(err)
	}
}

func GetTracesFromDB(txHash []byte) json.RawMessage {
	if tracesDB == nil {
		return []byte{}
	}
	res, err := tracesDB.Get(txHash)
	if err != nil {
		return []byte{}
	}
	return res
}

func DeleteTracesFromDB(txHash []byte) error {
	if tracesDB == nil {
		return fmt.Errorf("traces db is nil")
	}
	return tracesDB.Delete(txHash)
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
