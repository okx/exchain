package types

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	watcher "github.com/okex/exchain/x/evm/watcher"
)

var (
	// static gas limit for all blocks
	defaultGasLimit   = hexutil.Uint64(int64(^uint32(0)))
	defaultGasUsed    = hexutil.Uint64(0)
	defaultDifficulty = (*hexutil.Big)(big.NewInt(0))
)

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx clientcontext.CLIContext, bz []byte) (*evmtypes.MsgEthereumTx, error) {
	tx, err := evmtypes.TxDecoder(clientCtx.Codec)(bz, evmtypes.IGNORE_HEIGHT_CHECKING)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ethTx, ok := tx.(*evmtypes.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type %T, expected %T", tx, evmtypes.MsgEthereumTx{})
	}
	return ethTx, nil
}

func ToTransaction(tx *evmtypes.MsgEthereumTx, from *common.Address) *watcher.Transaction {
	rpcTx := &watcher.Transaction{
		From:     *from,
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}
	return rpcTx
}

// RpcBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func RpcBlockFromTendermint(clientCtx clientcontext.CLIContext, block *tmtypes.Block) (*watcher.Block, error) {
	gasLimit, err := BlockMaxGasFromConsensusParams(context.Background(), clientCtx)
	if err != nil {
		return nil, err
	}

	gasUsed, ethTxs, err := EthTransactionsFromTendermint(clientCtx, block.Txs, common.BytesToHash(block.Hash()), uint64(block.Height))
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

	return FormatBlock(block.Header, block.Size(), block.Hash(), gasLimit, gasUsed, ethTxs, bloom), nil
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header) *ethtypes.Header {
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: ethtypes.EmptyRootHash,
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
func EthTransactionsFromTendermint(clientCtx clientcontext.CLIContext, txs []tmtypes.Tx, blockHash common.Hash, blockNumber uint64) (*big.Int, []*watcher.Transaction, error) {
	var transactions []*watcher.Transaction
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
		txHash := tx.Hash(int64(blockNumber))
		tx, err := watcher.NewTransaction(ethTx, common.BytesToHash(txHash), blockHash, blockNumber, index)
		if err == nil {
			transactions = append(transactions, tx)
			index++
		}
	}

	return gasUsed, transactions, nil
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
) *watcher.Block {
	if len(header.DataHash) == 0 {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}
	parentHash := header.LastBlockID.Hash
	if parentHash == nil {
		parentHash = ethtypes.EmptyRootHash.Bytes()
	}
	ret := &watcher.Block{
		Number:           hexutil.Uint64(header.Height),
		Hash:             common.BytesToHash(curBlockHash),
		ParentHash:       common.BytesToHash(parentHash),
		Nonce:            watcher.BlockNonce{},    // PoW specific
		UncleHash:        ethtypes.EmptyUncleHash, // No uncles in Tendermint
		LogsBloom:        bloom,
		TransactionsRoot: common.BytesToHash(header.DataHash),
		StateRoot:        common.BytesToHash(header.AppHash),
		Miner:            common.BytesToAddress(header.ProposerAddress),
		MixHash:          common.Hash{},
		Difficulty:       hexutil.Uint64(0),
		TotalDifficulty:  hexutil.Uint64(0),
		ExtraData:        hexutil.Bytes{},
		Size:             hexutil.Uint64(size),
		GasLimit:         hexutil.Uint64(gasLimit), // Static gas limit
		GasUsed:          (*hexutil.Big)(gasUsed),
		Timestamp:        hexutil.Uint64(header.Time.Unix()),
		Uncles:           []common.Hash{},
		ReceiptsRoot:     ethtypes.EmptyRootHash,
	}
	if !reflect.ValueOf(transactions).IsNil() {
		ret.Transactions = transactions
	} else {
		ret.Transactions = []*watcher.Transaction{}
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
		txi, err := txDecoder(block.Txs[i], evmtypes.IGNORE_HEIGHT_CHECKING)
		if err != nil {
			continue
		}

		gasUsed += txi.GetGas()
	}
	return gasUsed
}

// EthHeaderWithBlockHashFromTendermint gets the eth Header with block hash from Tendermint block inside
func EthHeaderWithBlockHashFromTendermint(tmHeader *tmtypes.Header) (header *EthHeaderWithBlockHash, err error) {
	if tmHeader == nil {
		return header, errors.New("failed. nil tendermint block header")
	}

	header = &EthHeaderWithBlockHash{
		ParentHash:  common.BytesToHash(tmHeader.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(tmHeader.ProposerAddress),
		Root:        common.BytesToHash(tmHeader.AppHash),
		TxHash:      common.BytesToHash(tmHeader.DataHash),
		ReceiptHash: ethtypes.EmptyRootHash,
		Number:      (*hexutil.Big)(big.NewInt(tmHeader.Height)),
		// difficulty is not available for DPOS
		Difficulty: defaultDifficulty,
		GasLimit:   defaultGasLimit,
		GasUsed:    defaultGasUsed,
		Time:       hexutil.Uint64(tmHeader.Time.Unix()),
		Hash:       common.BytesToHash(tmHeader.Hash()),
	}

	return
}
