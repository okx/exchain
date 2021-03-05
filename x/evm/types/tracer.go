package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/spf13/viper"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	tracesDir = "traces"
)

var (
	tracesDB dbm.DB
)

func init() {
	openDB()
}

func openDB() {
	dataDir := filepath.Join(viper.GetString("home"), "data")
	var err error
	tracesDB, err = sdk.NewLevelDB(tracesDir, dataDir)
	if err != nil {
		panic(err)
	}
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

	txHash := hexutil.Encode(tmtypes.Tx(ctx.TxBytes()).Hash())

	saveToDB(txHash, res)
}

func saveToDB(txHash string, res json.RawMessage) {
	if tracesDB == nil {
		panic("traces db is nil")
	}
	err := tracesDB.Set([]byte(txHash), res)
	if err != nil {
		panic(err)
	}
}

func GetTracesFromDB(txHash string) json.RawMessage{
	if tracesDB == nil {
		return []byte{}
	}
	res, err := tracesDB.Get([]byte(txHash))
	if err != nil {
		return []byte{}
	}
	return res
}

func DeleteTracesFromDB(txHash string) {
	if tracesDB == nil {
		return
	}
	tracesDB.Delete([]byte(txHash))
}
