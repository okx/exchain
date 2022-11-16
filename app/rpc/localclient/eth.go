package localclient

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/types"
)

func toBlockNumberOrHash(number *big.Int) (bnh types.BlockNumberOrHash, err error) {
	var bn types.BlockNumber
	bn, err = toBlockNumber(number)
	if err != nil {
		return
	}
	bnh.BlockNumber = &bn
	return
}

func toBlockNumber(number *big.Int) (bn types.BlockNumber, err error) {
	if number == nil {
		bn = types.LatestBlockNumber
		return
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		bn = types.PendingBlockNumber
		return
	}
	if number.Cmp(big.NewInt(0).SetUint64(math.MaxInt64)) > 0 {
		err = fmt.Errorf("blocknumber too high")
		return
	}
	bn = types.BlockNumber(number.Int64())
	return
}

func toCallArg(msg ethereum.CallMsg) types.CallArgs {
	args := types.CallArgs{
		From: &msg.From,
		To:   msg.To,
	}
	if len(msg.Data) > 0 {
		args.Data = (*hexutil.Bytes)(&msg.Data)
	}
	args.Value = (*hexutil.Big)(msg.Value)
	args.Gas = (*hexutil.Uint64)(&msg.Gas)
	args.GasPrice = (*hexutil.Big)(msg.GasPrice)

	return args
}

type Eth struct {
	api *eth.PublicEthereumAPI
}

func NewLocalEth(api *eth.PublicEthereumAPI) *Eth {
	return &Eth{
		api: api,
	}
}

func (c *Eth) CodeAt(_ context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	var blockArg, err = toBlockNumberOrHash(blockNumber)
	if err != nil {
		return nil, err
	}
	return c.api.GetCode(contract, blockArg)
}

func (c *Eth) CallContract(_ context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var blockArg, err = toBlockNumberOrHash(blockNumber)
	if err != nil {
		return nil, err
	}
	return c.api.Call(toCallArg(call), blockArg, nil)
}

func (c *Eth) HeaderByNumber(_ context.Context, number *big.Int) (*ethtypes.Header, error) {
	blockNum, err := toBlockNumber(number)
	if err != nil {
		return nil, err
	}
	block, err := c.api.GetBlockByNumber(blockNum, false)
	if err == nil && block == nil {
		err = ethereum.NotFound
	}
	var head = &ethtypes.Header{
		ParentHash:  block.ParentHash,
		UncleHash:   block.UncleHash,
		Coinbase:    block.Miner,
		Root:        block.StateRoot,
		TxHash:      block.TransactionsRoot,
		ReceiptHash: block.ReceiptsRoot,
		Bloom:       block.LogsBloom,
		Difficulty:  new(big.Int).SetUint64(uint64(block.Difficulty)),
		Number:      new(big.Int).SetUint64(uint64(block.Number)),
		GasLimit:    uint64(block.GasLimit),
		GasUsed:     block.GasUsed.ToInt().Uint64(),
		Time:        uint64(block.Timestamp),
		Extra:       block.ExtraData,
		MixDigest:   block.MixHash,
		Nonce:       ethtypes.BlockNonce(block.Nonce),
	}
	return head, err
}

func (c *Eth) PendingCodeAt(_ context.Context, account common.Address) ([]byte, error) {
	var blockArg, err = toBlockNumberOrHash(big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return c.api.GetCode(account, blockArg)
}

func (c *Eth) PendingNonceAt(_ context.Context, account common.Address) (uint64, error) {
	var blockArg, err = toBlockNumberOrHash(big.NewInt(0))
	if err != nil {
		return 0, err
	}
	tc, err := c.api.GetTransactionCount(account, blockArg)
	if err != nil || tc == nil {
		return 0, err
	}
	return (uint64)(*tc), nil
}

func (c *Eth) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	gp := c.api.GasPrice()
	return (*big.Int)(gp), nil
}

func (c *Eth) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}

func (c *Eth) EstimateGas(_ context.Context, call ethereum.CallMsg) (uint64, error) {
	gas, err := c.api.EstimateGas(toCallArg(call))
	return uint64(gas), err
}

func (c *Eth) SendTransaction(ctx context.Context, tx *ethtypes.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = c.api.SendRawTransaction(data)
	return err
}

func (c *Eth) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]ethtypes.Log, error) {
	panic("implement me")
}

func (c *Eth) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- ethtypes.Log) (ethereum.Subscription, error) {
	panic("implement me")
}

func (c *Eth) NonceAt(_ context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	var blockArg types.BlockNumberOrHash
	blockArg, err = toBlockNumberOrHash(blockNumber)
	if err != nil {
		return
	}
	var tc *hexutil.Uint64
	tc, err = c.api.GetTransactionCount(account, blockArg)
	if err != nil {
		return
	}
	if tc != nil {
		nonce = uint64(*tc)
	}
	return
}

var _ bind.ContractCaller = (*Eth)(nil)
var _ bind.ContractFilterer = (*Eth)(nil)
var _ bind.ContractTransactor = (*Eth)(nil)
