package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	watcher "github.com/okex/exchain/x/evm/watcher"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockNumber represents decoding hex string to block values
type BlockNumber int64

const (
	// LatestBlockNumber mapping from "latest" to 0 for tm query
	LatestBlockNumber = BlockNumber(0)

	// EarliestBlockNumber mapping from "earliest" to 1 for tm query (earliest query not supported)
	EarliestBlockNumber = BlockNumber(1)

	// PendingBlockNumber mapping from "pending" to -1 for tm query
	PendingBlockNumber = BlockNumber(-1)
)

var ErrResourceNotFound = errors.New("resource not found")

// NewBlockNumber creates a new BlockNumber instance.
func NewBlockNumber(n *big.Int) BlockNumber {
	return BlockNumber(n.Int64())
}

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*bn = EarliestBlockNumber
		return nil
	case "latest":
		*bn = LatestBlockNumber
		return nil
	case "pending":
		*bn = PendingBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return fmt.Errorf("blocknumber too high")
	}

	*bn = BlockNumber(blckNum)
	return nil
}

// Int64 converts block number to primitive type
func (bn BlockNumber) Int64() int64 {
	return int64(bn)
}

// TmHeight is a util function used for the Tendermint RPC client. It returns
// nil if the block number is "latest". Otherwise, it returns the pointer of the
// int64 value of the height.
func (bn BlockNumber) TmHeight() *int64 {
	if bn == LatestBlockNumber {
		return nil
	}
	height := bn.Int64()
	return &height
}

type BlockNumberOrHash struct {
	BlockNumber      *BlockNumber `json:"blockNumber,omitempty"`
	BlockHash        *common.Hash `json:"blockHash,omitempty"`
	RequireCanonical bool         `json:"requireCanonical,omitempty"`
}

func (bnh *BlockNumberOrHash) UnmarshalJSON(data []byte) error {
	type erased BlockNumberOrHash
	e := erased{}
	err := json.Unmarshal(data, &e)
	if err == nil {
		if e.BlockNumber != nil && e.BlockHash != nil {
			return fmt.Errorf("cannot specify both BlockHash and BlockNumber, choose one or the other")
		}
		bnh.BlockNumber = e.BlockNumber
		bnh.BlockHash = e.BlockHash
		bnh.RequireCanonical = e.RequireCanonical
		return nil
	}
	var input string
	err = json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	switch input {
	case "earliest":
		bn := EarliestBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "latest":
		bn := LatestBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "pending":
		bn := PendingBlockNumber
		bnh.BlockNumber = &bn
		return nil
	default:
		if len(input) == 66 {
			hash := common.Hash{}
			err := hash.UnmarshalText([]byte(input))
			if err != nil {
				return err
			}
			bnh.BlockHash = &hash
			return nil
		} else {
			blckNum, err := hexutil.DecodeUint64(input)
			if err != nil {
				return err
			}
			if blckNum > math.MaxInt64 {
				return fmt.Errorf("blocknumber too high")
			}
			bn := BlockNumber(blckNum)
			bnh.BlockNumber = &bn
			return nil
		}
	}
}

func (bnh *BlockNumberOrHash) Number() (BlockNumber, bool) {
	if bnh.BlockNumber != nil {
		return *bnh.BlockNumber, true
	}
	return BlockNumber(0), false
}

func (bnh *BlockNumberOrHash) Hash() (common.Hash, bool) {
	if bnh.BlockHash != nil {
		return *bnh.BlockHash, true
	}
	return common.Hash{}, false
}

func BlockNumberOrHashWithNumber(blockNr BlockNumber) BlockNumberOrHash {
	return BlockNumberOrHash{
		BlockNumber:      &blockNr,
		BlockHash:        nil,
		RequireCanonical: false,
	}
}

func BlockNumberOrHashWithHash(hash common.Hash, canonical bool) BlockNumberOrHash {
	return BlockNumberOrHash{
		BlockNumber:      nil,
		BlockHash:        &hash,
		RequireCanonical: canonical,
	}
}

// Block represents a transaction returned to RPC clients.
type Block struct {
	Number           hexutil.Uint64     `json:"number"`
	Hash             common.Hash        `json:"hash"`
	ParentHash       common.Hash        `json:"parentHash"`
	Nonce            watcher.BlockNonce `json:"nonce"`
	UncleHash        common.Hash        `json:"sha3Uncles"`
	LogsBloom        ethtypes.Bloom     `json:"logsBloom"`
	TransactionsRoot common.Hash        `json:"transactionsRoot"`
	StateRoot        common.Hash        `json:"stateRoot"`
	Miner            common.Address     `json:"miner"`
	MixHash          common.Hash        `json:"mixHash"`
	Difficulty       hexutil.Uint64     `json:"difficulty"`
	TotalDifficulty  hexutil.Uint64     `json:"totalDifficulty"`
	ExtraData        hexutil.Bytes      `json:"extraData"`
	Size             hexutil.Uint64     `json:"size"`
	GasLimit         hexutil.Uint64     `json:"gasLimit"`
	GasUsed          *hexutil.Big       `json:"gasUsed"`
	Timestamp        hexutil.Uint64     `json:"timestamp"`
	Uncles           []common.Hash      `json:"uncles"`
	ReceiptsRoot     common.Hash        `json:"receiptsRoot"`
	Transactions     interface{}        `json:"transactions"`
}

func RpcBlockFromWatcherBlock(watcherBlock *watcher.FullTxBlock, fullTx bool) *Block {
	b := Block{
		Number:           watcherBlock.Number,
		Hash:             watcherBlock.Hash,
		ParentHash:       watcherBlock.ParentHash,
		Nonce:            watcherBlock.Nonce,
		UncleHash:        watcherBlock.UncleHash,
		LogsBloom:        watcherBlock.LogsBloom,
		TransactionsRoot: watcherBlock.TransactionsRoot,
		StateRoot:        watcherBlock.StateRoot,
		Miner:            watcherBlock.Miner,
		MixHash:          watcherBlock.MixHash,
		Difficulty:       watcherBlock.Difficulty,
		TotalDifficulty:  watcherBlock.TotalDifficulty,
		ExtraData:        watcherBlock.ExtraData,
		Size:             watcherBlock.Size,
		GasLimit:         watcherBlock.GasLimit,
		GasUsed:          watcherBlock.GasUsed,
		Timestamp:        watcherBlock.Timestamp,
		Uncles:           watcherBlock.Uncles,
		ReceiptsRoot:     watcherBlock.ReceiptsRoot,
	}
	if fullTx {
		b.Transactions = watcherBlock.FullTransactions
	} else {
		b.Transactions = watcherBlock.Transactions
	}
	return &b
}
