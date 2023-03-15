package types

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	clientcontext "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/tendermint/crypto/merkle"
	tmbytes "github.com/okx/okbchain/libs/tendermint/libs/bytes"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/evm/watcher"
)

var (
	// static gas limit for all blocks
	defaultGasLimit   = hexutil.Uint64(int64(^uint32(0)))
	defaultGasUsed    = hexutil.Uint64(0)
	defaultDifficulty = (*hexutil.Big)(big.NewInt(0))
)

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx clientcontext.CLIContext, bz []byte, height int64) (*evmtypes.MsgEthereumTx, error) {
	tx, err := evmtypes.TxDecoder(clientCtx.Codec)(bz, height)
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
func RpcBlockFromTendermint(clientCtx clientcontext.CLIContext, block *tmtypes.Block, fullTx bool) (*watcher.Block, error) {
	gasLimit, err := BlockMaxGasFromConsensusParams(context.Background(), clientCtx)
	if err != nil {
		return nil, err
	}

	gasUsed, ethTxs, err := EthTransactionsFromTendermint(clientCtx, block.Txs, common.BytesToHash(block.Hash()), uint64(block.Height))
	if err != nil {
		return nil, err
	}

	var bloom ethtypes.Bloom
	clientCtx = clientCtx.WithHeight(block.Height)
	res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryBloom, block.Height))
	if err == nil {
		var bloomRes evmtypes.QueryBloomFilter
		clientCtx.Codec.MustUnmarshalJSON(res, &bloomRes)
		bloom = bloomRes.Bloom
	}

	return FormatBlock(block.Header, block.Size(), block.Hash(), gasLimit, gasUsed, ethTxs, bloom, fullTx), nil
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
		ethTx, err := RawTxToEthTx(clientCtx, tx, int64(blockNumber))
		if err != nil {
			// continue to next transaction in case it's not a MsgEthereumTx
			continue
		}
		// TODO: Remove gas usage calculation if saving gasUsed per block
		gasUsed.Add(gasUsed, big.NewInt(int64(ethTx.GetGas())))
		tx, err := watcher.NewTransaction(ethTx, common.BytesToHash(ethTx.Hash), blockHash, blockNumber, index)
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
	gasUsed *big.Int, transactions []*watcher.Transaction, bloom ethtypes.Bloom, fullTx bool,
) *watcher.Block {
	transactionsRoot := ethtypes.EmptyRootHash
	if len(transactions) > 0 {
		txBzs := make([][]byte, len(transactions))
		for i := 0; i < len(transactions); i++ {
			txBzs[i] = transactions[i].Hash.Bytes()
		}
		transactionsRoot = common.BytesToHash(merkle.SimpleHashFromByteSlices(txBzs))
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
		TransactionsRoot: transactionsRoot,
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

	if fullTx {
		// return empty slice instead of nil for compatibility with Ethereum
		if len(transactions) == 0 {
			ret.Transactions = []*watcher.Transaction{}
		} else {
			ret.Transactions = transactions
		}
	} else {
		txHashes := make([]common.Hash, len(transactions))
		for i, tx := range transactions {
			txHashes[i] = tx.Hash
		}
		ret.Transactions = txHashes
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
		txi, err := txDecoder(block.Txs[i], block.Height)
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

func RawTxToRealTx(clientCtx clientcontext.CLIContext, bz tmtypes.Tx,
	blockHash common.Hash, blockNumber, index uint64) (sdk.Tx, error) {
	realTx, err := evmtypes.TxDecoder(clientCtx.CodecProy)(bz, int64(blockNumber))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	return realTx, nil
}

func RawTxResultToEthReceipt(chainID *big.Int, tr *ctypes.ResultTx, realTx sdk.Tx,
	blockHash common.Hash) (*watcher.TransactionResult, error) {
	// Convert tx bytes to eth transaction
	ethTx, ok := realTx.(*evmtypes.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type %T, expected %T", realTx, evmtypes.MsgEthereumTx{})
	}

	// try to get from event
	if from, err := GetEthSender(tr); err == nil {
		ethTx.BaseTx.From = from
	} else {
		// try to get from sig
		err := ethTx.VerifySig(chainID, tr.Height)
		if err != nil {
			return nil, err
		}
	}

	// Set status codes based on tx result
	var status = hexutil.Uint64(0)
	if tr.TxResult.IsOK() {
		status = hexutil.Uint64(1)
	}

	txData := tr.TxResult.GetData()
	data, err := evmtypes.DecodeResultData(txData)
	if err != nil {
		status = 0 // transaction failed
	}

	if len(data.Logs) == 0 {
		data.Logs = []*ethtypes.Log{}
	}
	contractAddr := &data.ContractAddress
	if data.ContractAddress == common.HexToAddress("0x00000000000000000000") {
		contractAddr = nil
	}

	// fix gasUsed when deliverTx ante handler check sequence invalid
	gasUsed := tr.TxResult.GasUsed
	if tr.TxResult.Code == sdkerrors.ErrInvalidSequence.ABCICode() {
		gasUsed = 0
	}

	receipt := watcher.TransactionReceipt{
		Status: status,
		//CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		LogsBloom:        data.Bloom,
		Logs:             data.Logs,
		TransactionHash:  common.BytesToHash(tr.Hash.Bytes()).String(),
		ContractAddress:  contractAddr,
		GasUsed:          hexutil.Uint64(gasUsed),
		BlockHash:        blockHash.String(),
		BlockNumber:      hexutil.Uint64(tr.Height),
		TransactionIndex: hexutil.Uint64(tr.Index),
		From:             ethTx.GetFrom(),
		To:               ethTx.To(),
	}

	rpcTx, err := watcher.NewTransaction(ethTx, common.BytesToHash(tr.Hash),
		blockHash, uint64(tr.Height), uint64(tr.Index))
	if err != nil {
		return nil, err
	}

	return &watcher.TransactionResult{TxType: hexutil.Uint64(watcher.EthReceipt),
		Receipt: &receipt, EthTx: rpcTx, EthTxLog: tr.TxResult.Log}, nil
}

func GetEthSender(tr *ctypes.ResultTx) (string, error) {
	for _, ev := range tr.TxResult.Events {
		if ev.Type == sdk.EventTypeMessage {
			fromAddr := ""
			realEvmTx := false
			for _, attr := range ev.Attributes {
				if string(attr.Key) == sdk.AttributeKeySender {
					fromAddr = string(attr.Value)
				}
				if string(attr.Key) == sdk.AttributeKeyModule &&
					string(attr.Value) == evmtypes.AttributeValueCategory { // to avoid the evm to cm tx enter
					realEvmTx = true
				}
				// find the sender
				if fromAddr != "" && realEvmTx {
					return fromAddr, nil
				}
			}
		}
	}
	return "", errors.New("No sender in Event")
}
