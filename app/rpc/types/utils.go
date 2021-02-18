package types

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	evmtypes "github.com/okex/okexchain/x/evm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx clientcontext.CLIContext, bz []byte) (*evmtypes.MsgEthereumTx, error) {
	tx, err := evmtypes.TxDecoder(clientCtx.Codec)(bz)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ethTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type %T, expected %T", tx, evmtypes.MsgEthereumTx{})
	}
	return &ethTx, nil
}

// NewTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransaction(tx *evmtypes.MsgEthereumTx, txHash, blockHash common.Hash, blockNumber, index uint64) (*Transaction, error) {
	// Verify signature and retrieve sender address
	from, err := tx.VerifySig(tx.ChainID())
	if err != nil {
		return nil, err
	}

	rpcTx := &Transaction{
		From:     from,
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Hash:     txHash,
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}

	if blockHash != (common.Hash{}) {
		rpcTx.BlockHash = &blockHash
		rpcTx.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		rpcTx.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	return rpcTx, nil
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(clientCtx clientcontext.CLIContext, block *tmtypes.Block, fullTx bool) (map[string]interface{}, error) {
	var blockTxs interface{}
	gasLimit, err := BlockMaxGasFromConsensusParams(context.Background(), clientCtx)
	if err != nil {
		return nil, err
	}

	transactions, gasUsed, ethTxs, err := EthTransactionsFromTendermint(clientCtx, block.Txs, common.BytesToHash(block.Hash()), uint64(block.Height))
	if err != nil {
		return nil, err
	}

	res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryBloom, block.Height))
	if err != nil {
		return nil, err
	}

	var bloomRes evmtypes.QueryBloomFilter
	clientCtx.Codec.MustUnmarshalJSON(res, &bloomRes)

	bloom := bloomRes.Bloom
	if fullTx {
		blockTxs = ethTxs
	} else {
		blockTxs = transactions
	}

	return FormatBlock(block.Header, block.Size(), block.Hash(), gasLimit, gasUsed, blockTxs, bloom), nil
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header) *ethtypes.Header {
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   common.Hash{},
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: common.Hash{},
		Difficulty:  nil,
		Number:      big.NewInt(header.Height),
		Time:        uint64(header.Time.Unix()),
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
	}
}

// EthTransactionsFromTendermint returns a slice of ethereum transaction hashes and the total gas usage from a set of
// tendermint block transactions.
func EthTransactionsFromTendermint(clientCtx clientcontext.CLIContext, txs []tmtypes.Tx, blockHash common.Hash, blockNumber uint64) ([]common.Hash, *big.Int, []*Transaction, error) {
	var transactionHashes []common.Hash
	var transactions []*Transaction
	gasUsed := big.NewInt(0)
	index := uint64(0)

	for _, tx := range txs {
		ethTx, err := RawTxToEthTx(clientCtx, tx)
		if err != nil {
			// continue to next transaction in case it's not a MsgEthereumTx
			continue
		}
		// TODO: Remove gas usage calculation if saving gasUsed per block
		gasUsed.Add(gasUsed, big.NewInt(int64(ethTx.GetGas())))
		transactionHashes = append(transactionHashes, common.BytesToHash(tx.Hash()))
		tx, err := NewTransaction(ethTx, common.BytesToHash(tx.Hash()), blockHash, blockNumber, index)
		if err == nil {
			transactions = append(transactions, tx)
			index++
		}
	}

	return transactionHashes, gasUsed, transactions, nil
}

// BlockMaxGasFromConsensusParams returns the gas limit for the latest block from the chain consensus params.
func BlockMaxGasFromConsensusParams(_ context.Context, clientCtx clientcontext.CLIContext) (int64, error) {
	//resConsParams, err := clientCtx.Client.ConsensusParams(nil)
	//if err != nil {
	//	return 0, err
	//}
	//
	//gasLimit := resConsParams.ConsensusParams.Block.MaxGas
	//if gasLimit == -1 {
	//	// Sets gas limit to max uint32 to not error with javascript dev tooling
	//	// This -1 value indicating no block gas limit is set to max uint64 with geth hexutils
	//	// which errors certain javascript dev tooling which only supports up to 53 bits
	//	gasLimit = int64(^uint32(0))
	//}
	//
	//return gasLimit, nil

	return int64(^uint32(0)), nil
}

// FormatBlock creates an ethereum block from a tendermint header and ethereum-formatted
// transactions.
func FormatBlock(
	header tmtypes.Header, size int, curBlockHash tmbytes.HexBytes, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) map[string]interface{} {
	if len(header.DataHash) == 0 {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}

	ret := map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(curBlockHash),
		"parentHash":       hexutil.Bytes(header.LastBlockID.Hash),
		"nonce":            hexutil.Uint64(0), // PoW specific
		"sha3Uncles":       common.Hash{},     // No uncles in Tendermint
		"logsBloom":        bloom,
		"transactionsRoot": hexutil.Bytes(header.DataHash),
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            common.Address{},
		"mixHash":          common.Hash{},
		"difficulty":       hexutil.Uint64(0),
		"totalDifficulty":  hexutil.Uint64(0),
		"extraData":        hexutil.Bytes{},
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"uncles":           []string{},
		"receiptsRoot":     common.Hash{},
	}
	switch transactions.(type) {
	case []common.Hash:
		ret["transactions"] = transactions.([]common.Hash)
	case []*Transaction:
		ret["transactions"] = transactions.([]*Transaction)
	}
	return ret
}

// GetKeyByAddress returns the private key matching the given address. If not found it returns false.
func GetKeyByAddress(keys []ethsecp256k1.PrivKey, address common.Address) (key *ethsecp256k1.PrivKey, exist bool) {
	for _, key := range keys {
		if bytes.Equal(key.PubKey().Address().Bytes(), address.Bytes()) {
			return &key, true
		}
	}
	return nil, false
}

// GetBlockCumulativeGas returns the cumulative gas used on a block up to a given
// transaction index. The returned gas used includes the gas from both the SDK and
// EVM module transactions.
func GetBlockCumulativeGas(cdc *codec.Codec, block *tmtypes.Block, idx int) uint64 {
	var gasUsed uint64
	txDecoder := evmtypes.TxDecoder(cdc)

	for i := 0; i < idx && i < len(block.Txs); i++ {
		txi, err := txDecoder(block.Txs[i])
		if err != nil {
			continue
		}

		switch tx := txi.(type) {
		case authtypes.StdTx:
			gasUsed += tx.GetGas()
		case evmtypes.MsgEthereumTx:
			gasUsed += tx.GetGas()
		}
	}
	return gasUsed
}
