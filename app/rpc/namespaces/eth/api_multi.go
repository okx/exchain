package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/okx/okbchain/app/rpc/monitor"
	rpctypes "github.com/okx/okbchain/app/rpc/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/evm/watcher"
	"github.com/okx/okbchain/x/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

const (
	FlagEnableMultiCall = "rpc.enable-multi-call"
)

// GetBalanceBatch returns the provided account's balance up to the provided block number.
func (api *PublicEthereumAPI) GetBalanceBatch(addresses []string, blockNrOrHash rpctypes.BlockNumberOrHash) (interface{}, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getBalanceBatch", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("addresses", addresses, "block number", blockNrOrHash)

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	clientCtx := api.clientCtx

	useWatchBackend := api.useWatchBackend(blockNum)
	if !(blockNum == rpctypes.PendingBlockNumber || blockNum == rpctypes.LatestBlockNumber) && !useWatchBackend {
		clientCtx = api.clientCtx.WithHeight(blockNum.Int64())
	}

	type accBalance struct {
		Type    token.AccType `json:"type"`
		Balance *hexutil.Big  `json:"balance"`
	}
	balances := make(map[string]accBalance)
	for _, addr := range addresses {
		address, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, fmt.Errorf("addr:%s,err:%s", addr, err)
		}
		if acc, err := api.wrappedBackend.MustGetAccount(address); err == nil {
			balance := acc.GetCoins().AmountOf(sdk.DefaultBondDenom).BigInt()
			if balance == nil {
				balances[addr] = accBalance{accountType(acc), (*hexutil.Big)(sdk.ZeroInt().BigInt())}
			} else {
				balances[addr] = accBalance{accountType(acc), (*hexutil.Big)(balance)}
			}
			continue
		}

		bs, err := api.clientCtx.Codec.MarshalJSON(auth.NewQueryAccountParams(address))
		if err != nil {
			return nil, err
		}
		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
		if err != nil {
			continue
		}

		var account authexported.Account
		if err := api.clientCtx.Codec.UnmarshalJSON(res, &account); err != nil {
			return nil, err
		}

		val := account.GetCoins().AmountOf(sdk.DefaultBondDenom).BigInt()
		accType := accountType(account)
		if accType == token.UserAccount || accType == token.ContractAccount {
			api.watcherBackend.CommitAccountToRpcDb(account)
			if blockNum != rpctypes.PendingBlockNumber {
				balances[addr] = accBalance{accType, (*hexutil.Big)(val)}
				continue
			}

			// update the address balance with the pending transactions value (if applicable)
			pendingTxs, err := api.backend.UserPendingTransactions(addr, -1)
			if err != nil {
				return nil, err
			}

			for _, tx := range pendingTxs {
				if tx == nil {
					continue
				}

				if tx.From.String() == addr {
					val = new(big.Int).Sub(val, tx.Value.ToInt())
				}
				if tx.To.String() == addr {
					val = new(big.Int).Add(val, tx.Value.ToInt())
				}
			}
		}
		balances[addr] = accBalance{accType, (*hexutil.Big)(val)}
	}
	return balances, nil
}

// MultiCall performs multiple raw contract call.
func (api *PublicEthereumAPI) MultiCall(args []rpctypes.CallArgs, blockNr rpctypes.BlockNumber, _ *[]evmtypes.StateOverrides) ([]hexutil.Bytes, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_multiCall", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("args", args, "block number", blockNr)

	blockNrOrHash := rpctypes.BlockNumberOrHashWithNumber(blockNr)
	rets := make([]hexutil.Bytes, 0, len(args))
	for _, arg := range args {
		ret, err := api.Call(arg, blockNrOrHash, nil)
		if err != nil {
			return rets, err
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

// GetTransactionsByBlock returns some transactions identified by number or hash.
func (api *PublicEthereumAPI) GetTransactionsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.Transaction, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getTransactionsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	txs, e := api.wrappedBackend.GetTransactionsByBlockNumber(uint64(blockNum), uint64(offset), uint64(limit))
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

	resBlock, err := api.backend.Block(&height)
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
func (api *PublicEthereumAPI) GetTransactionReceiptsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.TransactionReceipt, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_getTransactionReceiptsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	txs, err := api.GetTransactionsByBlock(blockNrOrHash, offset, limit)
	if err != nil || len(txs) == 0 {
		return nil, err
	}

	var receipts []*watcher.TransactionReceipt
	var block *ctypes.ResultBlock
	var blockHash common.Hash
	for _, tx := range txs {
		res, _ := api.wrappedBackend.GetTransactionReceipt(tx.Hash)
		if res != nil {
			receipts = append(receipts, res)
			continue
		}

		tx, err := api.clientCtx.Client.Tx(tx.Hash.Bytes(), false)
		if err != nil {
			// Return nil for transaction when not found
			return nil, nil
		}

		if block == nil {
			// Query block for consensus hash
			block, err = api.backend.Block(&tx.Height)
			if err != nil {
				return nil, err
			}
			blockHash = common.BytesToHash(block.Block.Hash())
		}

		// Convert tx bytes to eth transaction
		ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx, tx.Height)
		if err != nil {
			return nil, err
		}

		err = ethTx.VerifySig(ethTx.ChainID(), tx.Height)
		if err != nil {
			return nil, err
		}

		// Set status codes based on tx result
		var status = hexutil.Uint64(0)
		if tx.TxResult.IsOK() {
			status = hexutil.Uint64(1)
		}

		txData := tx.TxResult.GetData()
		data, err := evmtypes.DecodeResultData(txData)
		if err != nil {
			status = 0 // transaction failed
		}

		if len(data.Logs) == 0 || status == 0 {
			data.Logs = []*ethtypes.Log{}
			data.Bloom = ethtypes.BytesToBloom(make([]byte, 256))
		}

		contractAddr := &data.ContractAddress
		if data.ContractAddress == common.HexToAddress("0x00000000000000000000") {
			contractAddr = nil
		}

		// fix gasUsed when deliverTx ante handler check sequence invalid
		gasUsed := tx.TxResult.GasUsed
		if tx.TxResult.Code == sdkerrors.ErrInvalidSequence.ABCICode() {
			gasUsed = 0
		}

		receipt := &watcher.TransactionReceipt{
			Status: status,
			//CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
			LogsBloom:        data.Bloom,
			Logs:             data.Logs,
			TransactionHash:  common.BytesToHash(tx.Hash.Bytes()).String(),
			ContractAddress:  contractAddr,
			GasUsed:          hexutil.Uint64(gasUsed),
			BlockHash:        blockHash.String(),
			BlockNumber:      hexutil.Uint64(tx.Height),
			TransactionIndex: hexutil.Uint64(tx.Index),
			From:             ethTx.GetFrom(),
			To:               ethTx.To(),
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

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

	height := blockNum.Int64()
	if blockNum == rpctypes.LatestBlockNumber {
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
	}

	// try to get from watch db
	results, err = api.wrappedBackend.GetTxResultByBlock(api.clientCtx, uint64(height), uint64(offset), uint64(limit))
	if err == nil {
		return results, nil
	}

	// try to get from node
	resBlock, err := api.backend.Block(&height)
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
			txHash := resBlock.Block.Txs[idx].Hash()
			queryTx, err := api.clientCtx.Client.Tx(txHash, false)
			if err != nil {
				// Return nil for transaction when not found
				return nil, err
			}

			var res *watcher.TransactionResult
			switch realTx.GetType() {
			case sdk.EvmTxType:
				res, err = rpctypes.RawTxResultToEthReceipt(api.chainIDEpoch, queryTx, realTx, blockHash)
			case sdk.StdTxType:
				res, err = watcher.RawTxResultToStdResponse(api.clientCtx, queryTx, realTx, resBlock.Block.Time)
			}

			if err != nil {
				// Return nil for transaction when not found
				return nil, err
			}

			results = append(results, res)
		}
	}

	return results, nil
}
