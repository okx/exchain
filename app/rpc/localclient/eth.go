package localclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/types"
)

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func toBlockNumberOrHash(number *big.Int) (types.BlockNumberOrHash, error) {
	var blockArg types.BlockNumberOrHash
	err := blockArg.UnmarshalJSON([]byte(toBlockNumArg(number)))
	return blockArg, err
}

func toBlockNumber(number *big.Int) (types.BlockNumber, error) {
	var blockArg types.BlockNumber
	err := blockArg.UnmarshalJSON([]byte(toBlockNumArg(number)))
	return blockArg, err
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
	var err error
	//blockNum, err := toBlockNumber(number)
	if err != nil {
		return nil, err
	}
	var head *ethtypes.Header
	//block, err := c.api.GetBlockByNumber(blockNum, false)
	//if err == nil && head == nil {
	//	err = ethereum.NotFound
	//}
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

var _ bind.ContractCaller = (*Eth)(nil)
var _ bind.ContractFilterer = (*Eth)(nil)
var _ bind.ContractTransactor = (*Eth)(nil)
