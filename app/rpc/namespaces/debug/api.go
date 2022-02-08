package debug

import (
	"encoding/json"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"

	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicDebugAPI struct {
	clientCtx clientcontext.CLIContext
	logger    log.Logger
	backend   backend.Backend
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
func (api *PublicDebugAPI) TraceTransaction(txHash common.Hash) (interface{}, error) {
	resTrace, _, err := api.clientCtx.QueryWithData("app/trace", txHash.Bytes())
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
