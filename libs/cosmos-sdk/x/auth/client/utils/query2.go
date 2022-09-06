package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	types "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
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
