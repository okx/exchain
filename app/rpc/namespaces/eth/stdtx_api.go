package eth

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"

	"github.com/okex/exchain/app/rpc/monitor"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// GetTransactionReceiptsByBlock returns the transaction receipt identified by block hash or number.
func (api *PublicEthereumAPI) GetAllTransactionResultsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.TransactionResult, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getAllTransactionResultsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	var results []*watcher.TransactionResult

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	// try to get from watch db
	results, err = api.wrappedBackend.GetTxResultByBlock(api.clientCtx, uint64(blockNum), uint64(offset), uint64(limit))
	if err == nil && results != nil && len(results) != 0 {
		return results, nil
	}

	// try to get from node
	height := blockNum.Int64()
	if blockNum == rpctypes.LatestBlockNumber {
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
	}

	resBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}
	blockHash := common.BytesToHash(resBlock.Block.Hash())
	for idx := offset; idx < offset+limit && int(idx) < len(resBlock.Block.Txs); idx++ {
		realTx, err := rpctypes.RawTxToRealTx(api.clientCtx, resBlock.Block.Txs[idx],
			blockHash, uint64(resBlock.Block.Height), uint64(idx))
		if err != nil {
			return nil, err
		}

		if realTx != nil {
			txHash := resBlock.Block.Txs[idx].Hash(resBlock.Block.Height)
			queryTx, err := api.clientCtx.Client.Tx(txHash, false)
			if err != nil {
				// Return nil for transaction when not found
				return nil, err
			}

			var res *watcher.TransactionResult
			switch realTx.GetType() {
			case sdk.EvmTxType:
				res, err = rpctypes.RawTxResultToEthReceipt(api.clientCtx, queryTx, blockHash)
			case sdk.StdTxType:
				res, err = watcher.RawTxResultToStdResponse(api.clientCtx, queryTx, resBlock.Block.Time)
			}

			results = append(results, res)
		}
	}

	return results, nil
}
