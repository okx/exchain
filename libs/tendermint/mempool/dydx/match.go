package dydx

import (
	"container/list"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	dydxlib "github.com/okex/exchain/libs/dydx"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

type MatchEngine struct {
	depthBook *DepthBook

	marketPrice *big.Int

	contracts *dydxlib.Contracts
	privKey   *ecdsa.PrivateKey
	from      common.Address
	chainID   *big.Int
	ethCli    *ethclient.Client

	config DydxConfig

	sub ethereum.Subscription

	logger log.Logger
}

type DydxConfig struct {
	PrivKeyHex                 string
	ChainID                    string
	EthWsRpcUrl                string
	PerpetualV1ContractAddress string
	P1OrdersContractAddress    string
	P1MakerOracleAddress       string
	P1MarginAddress            string
}

type LogHandler interface {
	HandleOrderFilled(*contracts.P1OrdersLogOrderFilled)
	SubErr(error)
}

func NewMatchEngine(depthBook *DepthBook, config DydxConfig, handler LogHandler, logger log.Logger) (*MatchEngine, error) {
	var engine = &MatchEngine{
		depthBook: depthBook,
		config:    config,
		logger:    logger,
	}
	if engine.logger == nil {
		engine.logger = log.NewNopLogger()
	}

	var err error
	engine.privKey, err = crypto.HexToECDSA(config.PrivKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	engine.from = crypto.PubkeyToAddress(engine.privKey.PublicKey)
	engine.chainID, _ = new(big.Int).SetString(config.ChainID, 10)
	if engine.chainID == nil {
		return nil, fmt.Errorf("invalid chain id")
	}
	engine.ethCli, err = ethclient.Dial(config.EthWsRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial eth rpc url: %s, err: %w", config.EthWsRpcUrl, err)
	}
	txOps, err := bind.NewKeyedTransactorWithChainID(engine.privKey, engine.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOps, err: %w", err)
	}
	engine.contracts, err = dydxlib.NewContracts(
		common.HexToAddress(config.PerpetualV1ContractAddress),
		common.HexToAddress(config.P1OrdersContractAddress),
		common.HexToAddress(config.P1MakerOracleAddress),
		common.HexToAddress(config.P1MarginAddress),
		txOps,
		engine.ethCli,
	)

	if handler != nil {
		ordersAbi, err := contracts.P1OrdersMetaData.GetAbi()
		if err != nil {
			return nil, fmt.Errorf("failed to get orders abi, err: %w", err)
		}

		var query = ethereum.FilterQuery{
			Addresses: []common.Address{
				common.HexToAddress(config.P1OrdersContractAddress),
			},
			Topics: [][]common.Hash{
				{ordersAbi.Events["LogOrderFilled"].ID},
			},
		}
		ch := make(chan ethtypes.Log, 32)
		engine.sub, err = engine.ethCli.SubscribeFilterLogs(context.Background(), query, ch)
		if err != nil {
			return nil, fmt.Errorf("failed to subscribe filter logs, err: %w", err)
		}

		go func() {
			for {
				select {
				case err := <-engine.sub.Err():
					handler.SubErr(err)
				case log := <-ch:
					filledLog, err := engine.contracts.P1Orders.ParseLogOrderFilled(log)
					if err == nil {
						handler.HandleOrderFilled(filledLog)
					}
				}
			}
		}()
	}

	return engine, nil
}

type MatchResult struct {
	MatchedRecords []*MatchRecord
	TakerOrder     *WrapOrder

	OnChain chan bool
	Tx      *ethtypes.Transaction
}

func (r *MatchResult) AddMatchedRecord(fill *contracts.P1OrdersFill, makerOrder *WrapOrder) {
	r.MatchedRecords = append(r.MatchedRecords, &MatchRecord{
		Fill:  fill,
		Maker: makerOrder,
		Taker: r.TakerOrder,
	})
}

type MatchRecord struct {
	Fill  *contracts.P1OrdersFill
	Taker *WrapOrder
	Maker *WrapOrder
}

func (m *MatchEngine) Stop() {
	if m.sub != nil {
		m.sub.Unsubscribe()
	}
}

func (m *MatchEngine) Match(order *WrapOrder, maketPrice *big.Int) (*MatchResult, error) {
	m.logger.Debug("start match", "order", order.P1Order, "marketPrice", maketPrice)

	if order.Type() == BuyOrderType {
		return processOrder(order, m.depthBook.sellOrders, m.depthBook.buyOrders, maketPrice), nil
	} else if order.Type() == SellOrderType {
		return processOrder(order, m.depthBook.buyOrders, m.depthBook.sellOrders, maketPrice), nil
	} else {
		return nil, fmt.Errorf("invalid order type")
	}
}

func (m *MatchEngine) MatchAndTrade(order *WrapOrder) (*MatchResult, error) {
	marketPrice, err := m.contracts.P1MakerOracle.GetPrice(&bind.CallOpts{
		From: m.contracts.PerpetualV1Address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get market price, err: %w", err)
	}
	matched, err := m.Match(order, marketPrice)
	if err != nil {
		return nil, err
	}

	if len(matched.MatchedRecords) == 0 {
		return nil, nil
	}

	m.logger.Debug("match result", "matched", matched.MatchedRecords)

	op := dydxlib.NewTradeOperation(m.contracts)

	for _, record := range matched.MatchedRecords {
		var solOrder1 = WrapOrderToSignedSolOrder(record.Maker)
		var solOrder2 = WrapOrderToSignedSolOrder(record.Taker)
		if record.Maker.Type() == BuyOrderType {
			var tmp = solOrder1
			solOrder1 = solOrder2
			solOrder2 = tmp
		}

		fill := *record.Fill
		fill.Fee = solOrder1.LimitFee
		err = op.FillSignedSolOrderWithTaker(m.from, solOrder1, &fill)
		if err != nil {
			return matched, fmt.Errorf("failed to fill order, err: %w", err)
		}
		fill.Fee = solOrder2.LimitFee
		err = op.FillSignedSolOrderWithTaker(m.from, solOrder2, &fill)
		if err != nil {
			return matched, fmt.Errorf("failed to fill order, err: %w", err)
		}
	}
	matched.Tx, err = op.Commit(nil)
	if err != nil {
		return matched, fmt.Errorf("failed to commit, err: %w", err)
	}
	matched.OnChain = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-time.After(6 * time.Second):
				receipt, err := m.ethCli.TransactionReceipt(context.Background(), matched.Tx.Hash())
				if err == nil {
					m.logger.Debug("tx receipt received", "hash", matched.Tx.Hash(), "status", receipt.Status)
					if receipt.Status == 1 {
						matched.OnChain <- true
					} else {
						matched.OnChain <- false
					}
				} else {
					m.logger.Error("failed to get receipt", "hash", matched.Tx.Hash(), "err", err)
				}
			}
		}
	}()
	return matched, nil
}

func WrapOrderToSignedSolOrder(order *WrapOrder) *dydxlib.SignedSolOrder {
	return &dydxlib.SignedSolOrder{
		order.P1OrdersOrder, order.Sig,
	}
}

func (m *MatchEngine) trade(order1, order2 *dydxlib.SignedSolOrder, fill *contracts.P1OrdersFill) (*ethtypes.Transaction, error) {
	op := dydxlib.NewTradeOperation(m.contracts)
	err := op.FillSignedSolOrderWithTaker(m.from, order1, fill)
	if err != nil {
		return nil, fmt.Errorf("failed to fill order1, err: %w", err)
	}
	err = op.FillSignedSolOrderWithTaker(m.from, order2, fill)
	if err != nil {
		return nil, fmt.Errorf("failed to fill order2, err: %w", err)
	}
	tx, err := op.Commit(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to commit, err: %w", err)
	}
	return tx, nil
}

func IsIntNilOrZero(i *big.Int) bool {
	return i == nil || i.Cmp(zero) == 0
}

func isValidTriggerPrice(order *WrapOrder, marketPrice *big.Int) bool {
	if !IsIntNilOrZero(order.TriggerPrice) && !IsIntNilOrZero(marketPrice) {
		if order.Type() == BuyOrderType {
			if marketPrice.Cmp(order.TriggerPrice) < 0 {
				return false
			}
		} else {
			if marketPrice.Cmp(order.TriggerPrice) > 0 {
				return false
			}
		}
	}
	return true
}

var zero = big.NewInt(0)

func processOrder(takerOrder *WrapOrder, makerBook *OrderList, takerBook *OrderList, marketPrice *big.Int) *MatchResult {
	var matchResult = &MatchResult{
		TakerOrder: takerOrder,
	}

	if takerOrder.LeftAmount.Cmp(zero) <= 0 || !isValidTriggerPrice(takerOrder, marketPrice) {
		takerBook.Insert(takerOrder)
		return matchResult
	}

	var makerOrderElem *list.Element

	for {
		if makerOrderElem == nil {
			makerOrderElem = makerBook.Front()
		} else {
			makerOrderElem = makerOrderElem.Next()
		}
		if makerOrderElem == nil {
			break
		}
		makerOrder := makerOrderElem.Value.(*WrapOrder)

		if makerOrder.LeftAmount.Cmp(zero) <= 0 {
			continue
		}

		if takerOrder.Type() == BuyOrderType && takerOrder.Price().Cmp(makerOrder.Price()) < 0 {
			break
		}
		if takerOrder.Type() == SellOrderType && takerOrder.Price().Cmp(makerOrder.Price()) > 0 {
			break
		}

		if !isValidTriggerPrice(makerOrder, marketPrice) {
			continue
		}

		matchPrice := makerOrder.Price()
		if !IsIntNilOrZero(marketPrice) {
			if takerOrder.Type() == BuyOrderType {
				if takerOrder.Price().Cmp(marketPrice) >= 0 && makerOrder.Price().Cmp(marketPrice) <= 0 {
					matchPrice = marketPrice
				}
			} else {
				if makerOrder.Price().Cmp(marketPrice) >= 0 && takerOrder.Price().Cmp(marketPrice) <= 0 {
					matchPrice = marketPrice
				}
			}
		}

		matchAmount := big.NewInt(0).Set(takerOrder.LeftAmount)
		if matchAmount.Cmp(makerOrder.LeftAmount) > 0 {
			matchAmount.Set(makerOrder.LeftAmount)
		}
		matchResult.AddMatchedRecord(&contracts.P1OrdersFill{
			Amount: matchAmount,
			Price:  matchPrice,
		}, makerOrder)

		takerOrder.LeftAmount.Sub(takerOrder.LeftAmount, matchAmount)
		makerOrder.LeftAmount.Sub(makerOrder.LeftAmount, matchAmount)

		takerOrder.FrozenAmount.Add(takerOrder.FrozenAmount, matchAmount)
		makerOrder.FrozenAmount.Add(makerOrder.FrozenAmount, matchAmount)

		if takerOrder.LeftAmount.Cmp(zero) == 0 {
			break
		}
	}
	takerBook.Insert(takerOrder)
	return matchResult
}
