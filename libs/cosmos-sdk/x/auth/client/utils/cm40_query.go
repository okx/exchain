package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	codectypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	types "github.com/okx/okbchain/libs/cosmos-sdk/types/ibc-adapter"
	ibc_tx "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
)

func Query40Tx(cliCtx context.CLIContext, hashHexStr string) (*types.TxResponse, error) {
	// strip 0x prefix
	if strings.HasPrefix(hashHexStr, "0x") {
		hashHexStr = hashHexStr[2:]
	}

	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resTx, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
	if err != nil {
		return nil, err
	}

	out, err := mk40TxResult(cliCtx, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, err
	}

	return out, nil
}

func Query40TxsByEvents(cliCtx context.CLIContext, events []string, page, limit int) (*types.SearchTxsResult, error) {
	if len(events) == 0 {
		return nil, errors.New("must declare at least one event to search")
	}

	if page <= 0 {
		return nil, errors.New("page must greater than 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit must greater than 0")
	}

	// XXX: implement ANY
	query := strings.Join(events, " AND ")

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	prove := !cliCtx.TrustNode

	resTxs, err := node.TxSearch(query, prove, page, limit, "")
	if err != nil {
		return nil, err
	}

	if prove {
		for _, tx := range resTxs.Txs {
			err := ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, resTxs.Txs)
	if err != nil {
		return nil, err
	}

	txs, err := format40TxResults(cliCtx, resTxs.Txs, resBlocks)
	if err != nil {
		return nil, err
	}

	result := types.NewSearchTxsResult(uint64(resTxs.TotalCount), uint64(len(txs)), uint64(page), uint64(limit), txs)

	return result, nil
}

// formatTxResults parses the indexed txs into a slice of TxResponse objects.
func format40TxResults(cliCtx context.CLIContext, resTxs []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]*types.TxResponse, error) {
	var err error
	out := make([]*types.TxResponse, len(resTxs))
	for i := range resTxs {
		out[i], err = mk40TxResult(cliCtx, resTxs[i], resBlocks[resTxs[i].Height])
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func mk40TxResult(cliCtx context.CLIContext, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (*types.TxResponse, error) {
	txb, err := ibc_tx.CM40TxDecoder(cliCtx.CodecProy.GetProtocMarshal())(resTx.Tx)
	if nil != err {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	any := p.AsAny()
	return types.NewResponseResultTx(resTx, any, resBlock.Block.Time.Format(time.RFC3339)), nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}
