package dydx

import (
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/okex/exchain/libs/dydx/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type tempTradeArg struct {
	Maker  string
	Taker  string
	Trader string
	Data   string
}

type TradeOperation struct {
	contracts *Contracts
	// orders    *Orders

	trades    []tempTradeArg
	committed bool
}

func NewTradeOperation(contracts *Contracts) *TradeOperation {
	return &TradeOperation{
		contracts: contracts,
	}
}

func (op *TradeOperation) addTradeArg(
	maker string,
	taker string,
	trader string,
	data string,
) *TradeOperation {
	if op.committed {
		panic("Operation already committed")
	}
	op.trades = append(op.trades, tempTradeArg{
		Trader: trader,
		Data:   data,
		Maker:  strings.ToLower(maker),
		Taker:  strings.ToLower(taker),
	})
	return op
}

func (op *TradeOperation) FillSignedOrder(
	order *SignedOrder,
	amount *big.Int,
	price Price,
	fee Fee,
) error {
	return op.FillSignedOrderWithTaker(order.Taker, order, amount, price, fee)
}

func (op *TradeOperation) FillSignedOrderWithTaker(
	taker string,
	order *SignedOrder,
	amount *big.Int,
	price Price,
	fee Fee,
) error {
	tradeData, err := FillToTradeData(order, amount, price, fee)
	if err != nil {
		return err
	}
	op.addTradeArg(
		order.Maker,
		taker,
		tradeData,
		op.contracts.P1OrdersAddress.String(),
	)
	return nil
}

func (op *TradeOperation) Commit(ops *bind.TransactOpts) (tx *types.Transaction, err error) {
	if op.committed {
		return nil, fmt.Errorf("operation already committed")
	}
	if len(op.trades) == 0 {
		return nil, fmt.Errorf("no tradeArgs have been added to trade")
	}

	defer func() {
		if rec := recover(); rec != nil {
			op.committed = false
			err = fmt.Errorf("error committing trade: %v", rec)
		}
		op.committed = true
	}()

	accountSet := make(map[string]struct{})
	for _, trade := range op.trades {
		accountSet[trade.Maker] = struct{}{}
		accountSet[trade.Taker] = struct{}{}
	}
	accounts := make([]string, 0, len(accountSet))
	for k := range accountSet {
		accounts = append(accounts, k)
	}
	sort.Strings(accounts)

	var tradeArgs []contracts.P1TradeTradeArg
	for _, trade := range op.trades {
		makerIndex := sort.SearchStrings(accounts, trade.Maker)
		takerIndex := sort.SearchStrings(accounts, trade.Taker)
		tradeArgs = append(tradeArgs, contracts.P1TradeTradeArg{
			MakerIndex: big.NewInt(int64(makerIndex)),
			TakerIndex: big.NewInt(int64(takerIndex)),
			Trader:     common.HexToAddress(trade.Trader),
			Data:       common.Hex2Bytes(trade.Data),
		})
	}

	accounts4Eth := make([]common.Address, len(accounts))
	for i, account := range accounts {
		accounts4Eth[i] = common.HexToAddress(account)
	}

	return op.contracts.PerpetualV1.Trade(combineTxOps(ops, op.contracts.txOps), accounts4Eth, tradeArgs)
}
