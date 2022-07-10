package eth

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/viper"

	"github.com/okex/exchain/app/rpc/monitor"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// GetTransactionsByBlock returns some transactions identified by number or hash.
func (api *PublicEthereumAPI) GetAllTransactionsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.Transaction, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getTransactionsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	txs, e := api.wrappedBackend.GetTransactionsWithStdByBlockNumber(uint64(blockNum), uint64(offset), uint64(limit))
	if e == nil && txs != nil {
		return txs, nil
	}

	height := blockNum.Int64()
	switch blockNum {
	case rpctypes.PendingBlockNumber:
		// get all the EVM pending txs
		pendingTxs, err := api.backend.PendingTransactions()
		if err != nil {
			return nil, err
		}
		switch {
		case len(pendingTxs) <= int(offset):
			return nil, nil
		case len(pendingTxs) < int(offset+limit):
			return pendingTxs[offset:], nil
		default:
			return pendingTxs[offset : offset+limit], nil
		}
	case rpctypes.LatestBlockNumber:
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
	}

	resBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}
	for idx := offset; idx < offset+limit && int(idx) < len(resBlock.Block.Txs); idx++ {
		tx, _ := api.getTransactionByBlockAndIndex(resBlock.Block, idx)
		if tx != nil {
			txs = append(txs, tx)
		}
	}
	return txs, nil
}

// GetTransactionReceiptsByBlock returns the transaction receipt identified by block hash or number.
func (api *PublicEthereumAPI) GetAllTransactionResultsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.TransactionResult, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getTransactionReceiptsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	txs, err := api.GetAllTransactionsByBlock(blockNrOrHash, offset, limit)
	if err != nil || len(txs) == 0 {
		return nil, err
	}

	var results []*watcher.TransactionResult
	//var block *ctypes.ResultBlock
	//var blockHash common.Hash
	for _, tx := range txs {
		var res *watcher.TransactionResult
		// std tx
		if tx.R == nil && tx.S == nil && tx.V == nil {
			stdResponse, _ := api.wrappedBackend.GetTransactionResponse(tx.Hash)
			if stdResponse != nil {
				var realTx authtypes.StdTx
				err := api.clientCtx.Codec.UnmarshalBinaryLengthPrefixed(stdResponse.Tx, &realTx)
				if err != nil {
					return nil, err
				}

				response := sdk.NewResponseResultTx(stdResponse.ResultTx, &realTx, stdResponse.Timestamp.Format(time.RFC3339))
				res = &watcher.TransactionResult{TxType: hexutil.Uint64(watcher.StdResponse), Response: &response}
			}
		} else {
			receipt, _ := api.wrappedBackend.GetTransactionReceipt(tx.Hash)
			if receipt != nil {
				res = &watcher.TransactionResult{TxType: hexutil.Uint64(watcher.EthReceipt), Receipt: receipt}
			}
		}

		if res != nil {
			results = append(results, res)
			continue
		}

		//tx, err := api.clientCtx.Client.Tx(tx.Hash.Bytes(), false)
		//if err != nil {
		//	// Return nil for transaction when not found
		//	return nil, nil
		//}
		//
		//if block == nil {
		//	// Query block for consensus hash
		//	block, err = api.clientCtx.Client.Block(&tx.Height)
		//	if err != nil {
		//		return nil, err
		//	}
		//	blockHash = common.BytesToHash(block.Block.Hash())
		//}
		//
		//// Convert tx bytes to eth transaction
		//ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx)
		//if err != nil {
		//	return nil, err
		//}
		//
		//err = ethTx.VerifySig(ethTx.ChainID(), tx.Height)
		//if err != nil {
		//	return nil, err
		//}
		//
		//// Set status codes based on tx result
		//var status = hexutil.Uint64(0)
		//if tx.TxResult.IsOK() {
		//	status = hexutil.Uint64(1)
		//}
		//
		//txData := tx.TxResult.GetData()
		//data, err := evmtypes.DecodeResultData(txData)
		//if err != nil {
		//	status = 0 // transaction failed
		//}
		//
		//if len(data.Logs) == 0 {
		//	data.Logs = []*ethtypes.Log{}
		//}
		//contractAddr := &data.ContractAddress
		//if data.ContractAddress == common.HexToAddress("0x00000000000000000000") {
		//	contractAddr = nil
		//}
		//
		//// fix gasUsed when deliverTx ante handler check sequence invalid
		//gasUsed := tx.TxResult.GasUsed
		//if tx.TxResult.Code == sdkerrors.ErrInvalidSequence.ABCICode() {
		//	gasUsed = 0
		//}
		//
		//receipt := &watcher.TransactionReceipt{
		//	Status: status,
		//	//CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		//	LogsBloom:        data.Bloom,
		//	Logs:             data.Logs,
		//	TransactionHash:  common.BytesToHash(tx.Hash.Bytes()).String(),
		//	ContractAddress:  contractAddr,
		//	GasUsed:          hexutil.Uint64(gasUsed),
		//	BlockHash:        blockHash.String(),
		//	BlockNumber:      hexutil.Uint64(tx.Height),
		//	TransactionIndex: hexutil.Uint64(tx.Index),
		//	From:             ethTx.GetFrom(),
		//	To:               ethTx.To(),
		//}
		//receipts = append(receipts, receipt)
	}

	return results, nil
}
