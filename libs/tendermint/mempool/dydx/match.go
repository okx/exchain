package dydx

import (
	"container/list"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rpc"

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

type PubSub interface {
	Unsubscribe(id rpc.ID) bool
	SubscribeLogs(conn chan<- *ethtypes.Log, query ethereum.FilterQuery) (rpc.ID, error)
}

type MatchEngine struct {
	depthBook *DepthBook

	marketPrice *big.Int

	contracts *dydxlib.Contracts
	privKey   *ecdsa.PrivateKey
	from      common.Address
	nonce     uint64
	chainID   *big.Int
	ethCli    *ethclient.Client
	httpCli   *ethclient.Client

	config DydxConfig

	sub ethereum.Subscription

	pubsub   PubSub
	pubsubID string

	logger log.Logger
}

type DydxConfig struct {
	PrivKeyHex                 string
	ChainID                    string
	EthWsRpcUrl                string
	EthHttpRpcUrl              string
	PerpetualV1ContractAddress string
	P1OrdersContractAddress    string
	P1MakerOracleAddress       string
	P1MarginAddress            string
	VMode                      bool
}

type LogHandler interface {
	HandleOrderFilled(*contracts.P1OrdersLogOrderFilled)
	SubErr(error)
}

func NewMatchEngine(api PubSub, depthBook *DepthBook, config DydxConfig, handler LogHandler, logger log.Logger) (*MatchEngine, error) {
	var engine = &MatchEngine{
		depthBook: depthBook,
		config:    config,
		logger:    logger,
		pubsub:    api,
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
	if !config.VMode || api == nil {
		engine.ethCli, err = ethclient.Dial(config.EthWsRpcUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to dial eth rpc url: %s, err: %w", config.EthWsRpcUrl, err)
		}
	}

	engine.httpCli, err = ethclient.Dial(config.EthHttpRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial eth rpc url: %s, err: %w", config.EthHttpRpcUrl, err)
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
		engine.httpCli,
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

		if config.VMode && api != nil {
			ch := make(chan *ethtypes.Log, 32)
			id, err := api.SubscribeLogs(ch, query)
			if err != nil {
				return nil, fmt.Errorf("failed to subscribe local logs, err: %w", err)
			}
			engine.pubsubID = string(id)
			go func() {
				for log := range ch {
					filledLog, err := engine.contracts.P1Orders.ParseLogOrderFilled(*log)
					if err == nil {
						handler.HandleOrderFilled(filledLog)
					}
				}
			}()
		} else {
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
	}

	return engine, nil
}

type MatchResult struct {
	MatchedRecords []*MatchRecord
	TakerOrder     *WrapOrder

	OnChain chan bool
	Tx      *ethtypes.Transaction
	NoSend  bool

	tradeOps *dydxlib.TradeOperation
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
	if m.pubsubID != "" {
		m.pubsub.Unsubscribe(rpc.ID(m.pubsubID))
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

func (m *MatchEngine) Rollback(matchResult *MatchResult) {
	for _, record := range matchResult.MatchedRecords {
		record.Maker.Unfrozen(record.Fill.Amount)
		record.Taker.Unfrozen(record.Fill.Amount)
	}
}

func (m *MatchEngine) matchAndTrade(order *WrapOrder, noSend bool) (*MatchResult, error) {
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
	matched.NoSend = noSend

	m.logger.Debug("match result", "matched", matched.MatchedRecords)

	var needRollback bool
	defer func() {
		if needRollback {
			m.Rollback(matched)
		}
	}()

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
			needRollback = true
			return matched, fmt.Errorf("failed to fill order, err: %w", err)
		}
		fill.Fee = solOrder2.LimitFee
		err = op.FillSignedSolOrderWithTaker(m.from, solOrder2, &fill)
		if err != nil {
			needRollback = true
			return matched, fmt.Errorf("failed to fill order, err: %w", err)
		}
	}
	if !noSend {
		matched.Tx, err = op.Commit(&bind.TransactOpts{NoSend: noSend})
		if err != nil {
			needRollback = true
			return matched, fmt.Errorf("failed to commit, err: %w", err)
		}
		m.logger.Debug("commit tx", "tx", matched.Tx.Hash().Hex())
	} else {
		matched.tradeOps = op
	}

	matched.OnChain = make(chan bool, 1)
	if noSend {
		return matched, nil
	}

	go func(txHash common.Hash) {
		m.logger.Debug("wait tx", "tx", txHash.Hex())
		count := 0
		for {
			if count == 10 {
				m.logger.Error("wait tx timeout", "tx", txHash.Hex())
				matched.OnChain <- false
				return
			}
			select {
			case <-time.After(5 * time.Second):
				receipt, err := m.httpCli.TransactionReceipt(context.Background(), txHash)
				if err == nil {
					m.logger.Debug("tx receipt received", "hash", txHash, "status", receipt.Status)
					if receipt.Status == 1 {
						matched.OnChain <- true
					} else {
						matched.OnChain <- false
					}
					return
				} else {
					m.logger.Error("failed to get receipt", "hash", txHash, "err", err)
				}
			}
			count += 1
		}
	}(matched.Tx.Hash())
	return matched, nil
}

func (m *MatchEngine) MatchAndTrade(order *WrapOrder) (*MatchResult, error) {
	return m.matchAndTrade(order, m.config.VMode)
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

		takerOrder.Frozen(matchAmount)
		makerOrder.Frozen(matchAmount)

		if takerOrder.LeftAmount.Cmp(zero) == 0 {
			break
		}
	}
	takerBook.Insert(takerOrder)
	return matchResult
}
