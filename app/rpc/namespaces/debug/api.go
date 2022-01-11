package debug

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm"

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
func (api *PublicDebugAPI) TraceTransaction(txHash common.Hash) (hexutil.Bytes, error) {
	// Get transaction by hash
	tx, err := api.clientCtx.Client.Tx(txHash.Bytes(), false)
	if err != nil {
		//to keep consistent with rpc of ethereum, should be return nil
		return nil, nil
	}

	// check if block number is 0
	if tx.Height == 0 {
		return nil, errors.New("tx height must not be zero")
	}

	res, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}
	block := res.Block
	// check tx index is not out of bound
	if uint32(len(block.Txs)) < tx.Index {
		return nil, fmt.Errorf("transaction not included in block %v", block.Height)
	}
	targetTx, err := evm.TxDecoder(api.clientCtx.Codec)(tx.Tx)
	if err != nil {
		return nil, err
	}

	var predecessors []*sdk.Tx
	for _, txBz := range block.Txs[:tx.Index] {
		tx, err := evm.TxDecoder(api.clientCtx.Codec)(txBz)
		if err != nil {
			return nil, err
		}
		predecessors = append(predecessors, &tx)
	}
	traceTxRequest := sdk.QueryTraceParams{
		TraceTx:           &targetTx,
		TxBytes:           tx.Tx,
		Predecessors:      predecessors,
		PredecessorsBytes: block.Txs[:tx.Index],
		Block:             block,
	}

	bs, err := api.clientCtx.Codec.MarshalJSON(traceTxRequest)
	if err != nil {
		return nil, err
	}

	resTrace, _, err := api.clientCtx.QueryWithData("app/trace", bs)
	if err != nil {
		return nil, err
	}
	return resTrace, nil
}
