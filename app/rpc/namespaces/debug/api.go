package debug

import (
	"encoding/json"
	"fmt"
	syslog "log"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"

	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/app/rpc/monitor"
	"github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	rpctypes "github.com/okex/exchain/app/rpc/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicDebugAPI struct {
	clientCtx clientcontext.CLIContext
	logger    log.Logger
	backend   backend.Backend
	Metrics   map[string]*monitor.RpcMetrics
}

// NewPublicTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewAPI(clientCtx clientcontext.CLIContext, log log.Logger, backend backend.Backend) *PublicDebugAPI {
	api := &PublicDebugAPI{
		clientCtx: clientCtx,
		backend:   backend,
		logger:    log.With("module", "json-rpc", "namespace", "debug"),
	}
	return api
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PublicDebugAPI) TraceTransaction(txHash common.Hash, config evmtypes.TraceConfig) (interface{}, error) {
	start := time.Now()
	syslog.Println("TraceTransaction, start", start)
	defer syslog.Println("TraceTransaction, stop", time.Since(start))
	monitor := monitor.GetMonitor("debug_traceTransaction", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	err := evmtypes.TestTracerConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("tracer err : %s", err.Error())
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	queryParam := sdk.QueryTraceTx{
		TxHash:      txHash,
		ConfigBytes: configBytes,
	}
	queryBytes, err := json.Marshal(&queryParam)
	if err != nil {
		return nil, err
	}
	_, err = api.clientCtx.Client.Tx(txHash.Bytes(), false)
	if err != nil {
		return nil, err
	}
	resTrace, _, err := api.clientCtx.QueryWithData("app/trace", queryBytes)
	if err != nil {
		return nil, err
	}

	var res sdk.Result
	if err := api.clientCtx.Codec.UnmarshalBinaryBare(resTrace, &res); err != nil {
		return nil, err
	}
	var decodedResult interface{}
	if err := json.Unmarshal(res.Data, &decodedResult); err != nil {
		return nil, err
	}

	return decodedResult, nil
}

func (api *PublicDebugAPI) traceBlockByNumber(blockNum rpctypes.BlockNumber, config *evmtypes.TraceConfig) (interface{}, error) {
	syslog.Println("traceBlockByNumber, start", time.Now())
	defer syslog.Println("traceBlockByNumber, stop", time.Now())
	err := evmtypes.TestTracerConfig(config)
	if err != nil {
		return nil, fmt.Errorf("tracer err : %s", err.Error())
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	//query block from db
	height := blockNum.Int64()
	_, err = api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}

	queryParam := sdk.QueryTraceBlock{
		Height:      height,
		ConfigBytes: configBytes,
	}
	queryBytes, err := json.Marshal(&queryParam)
	if err != nil {
		return nil, err
	}
	resTrace, _, err := api.clientCtx.QueryWithData("app/traceBlock", queryBytes)
	if err != nil {
		return nil, err
	}

	var results []sdk.QueryTraceTxResult
	if err := json.Unmarshal(resTrace, &results); err != nil {
		return nil, err
	}
	rpcResults := []types.TraceTxResult{}
	for _, res := range results {
		rpcRes := types.TraceTxResult{
			TxIndex: res.TxIndex,
			Txhash:  res.TxHash,
		}
		if res.Error != nil {
			rpcRes.Error = res.Error.Error()
		} else {
			if err := json.Unmarshal(res.Result, &rpcRes.Result); err != nil {
				rpcRes.Error = err.Error()
			}
		}
		rpcResults = append(rpcResults, rpcRes)
	}
	return rpcResults, nil
}

func (api *PublicDebugAPI) TraceBlockByHash(hash common.Hash, config *evmtypes.TraceConfig) (interface{}, error) {
	monitor := monitor.GetMonitor("debug_traceBlockByHash", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	block, err := api.backend.GetBlockByHash(hash, false)
	if err != nil {
		return nil, err
	}
	return api.TraceBlockByNumber(rpctypes.BlockNumber(block.Number), config)
}

func (api *PublicDebugAPI) TraceBlockByNumber(blockNum rpctypes.BlockNumber, config *evmtypes.TraceConfig) (interface{}, error) {
	monitor := monitor.GetMonitor("debug_traceBlockByNumber", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	return api.traceBlockByNumber(blockNum, config)
}
