package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/spf13/viper"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	tracesDir          = "traces"
	FlagEnableTraces   = "enable-evm-traces"
	FlagTraceSegment   = "evm-trace-segment"
	FlagTraceFromAddrs = "evm-trace-from-addrs"
	FlagTraceToAddrs   = "evm-trace-to-addrs"
)

var (
	tracesDB     dbm.DB
	enableTraces bool

	// trace from/to addr
	traceFromAddrs, traceToAddrs map[string]struct{}

	// trace segment
	step, total, num int64
)

func init() {
	server.TrapSignal(func() {
		if tracesDB != nil {
			tracesDB.Close()
		}
	})
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

func saveTraceResult(ctx sdk.Context, tracer vm.Tracer, result *core.ExecutionResult) {
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

		res, err = json.Marshal(&TraceExecutionResult{
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

	if err != nil {
		res = []byte(err.Error())
	}

	saveToDB(tmtypes.Tx(ctx.TxBytes()).Hash(), res)
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
