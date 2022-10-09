package dydx

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	dydxlib "github.com/okex/exchain/libs/dydx"
	"github.com/okex/exchain/libs/dydx/contracts"
)

type MatchEngine struct {
	depthBook *DepthBook

	contracts *dydxlib.Contracts
	privKey   *ecdsa.PrivateKey
	from      common.Address
	chainID   *big.Int
	ethCli    *ethclient.Client

	config DydxConfig
}

type DydxConfig struct {
	PrivKeyHex                 string
	ChainID                    string
	EthWsRpcUrl                string
	PerpetualV1ContractAddress string
	P1OrdersContractAddress    string
}

// chainID *big.Int, ethRpcUrl string, fromBlockNum *big.Int,
//	privKey, perpetualV1ContractAddress, p1OrdersContractAddress string

func NewMatchEngine(depthBook *DepthBook, config DydxConfig) (*MatchEngine, error) {
	var engine = &MatchEngine{
		depthBook: depthBook,
		config:    config,
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
		txOps,
		engine.ethCli,
	)

	return engine, nil
}

type MatchResult struct {
	MatchedRecords []*MatchRecord
	TakerOrder     *WrapOrder
}

func (r *MatchResult) AddMatchedRecord(price *big.Int, amount *big.Int, makerOrder *WrapOrder) {
	r.MatchedRecords = append(r.MatchedRecords, &MatchRecord{
		Price:  price,
		Amount: amount,
		Maker:  makerOrder,
		Taker:  r.TakerOrder,
	})
}

type MatchRecord struct {
	Price  *big.Int
	Amount *big.Int
	Taker  *WrapOrder
	Maker  *WrapOrder
}

func (m *MatchEngine) Match(order *WrapOrder) (*MatchResult, error) {
	if order.Type() == BuyOrderType {
		return processOrder(order, m.depthBook.sellOrders, m.depthBook.buyOrders)
	} else if order.Type() == SellOrderType {
		return processOrder(order, m.depthBook.buyOrders, m.depthBook.sellOrders)
	} else {
		return nil, fmt.Errorf("invalid order type")
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

func processOrder(takerOrder *WrapOrder, makerBook *OrderList, takerBook *OrderList) (*MatchResult, error) {
	var matchResult = &MatchResult{
		TakerOrder: takerOrder,
	}
	for {
		makerOrderElem := makerBook.Front()
		if makerOrderElem == nil {
			break
		}
		makerOrder := makerOrderElem.Value.(*WrapOrder)
		if takerOrder.Type() == BuyOrderType && takerOrder.Price().Cmp(makerOrder.Price()) < 0 {
			break
		}
		if takerOrder.Type() == SellOrderType && takerOrder.Price().Cmp(makerOrder.Price()) > 0 {
			break
		}
		//marketPrice := makerOrder.Price()
		//matchAmount := takerOrder.Amount()
		//if takerOrder.Amount().Cmp(makerOrder.Amount()) > 0 {
		//	matchAmount = makerOrder.Amount()
		//}
		//matchResult.AddMatchedRecord(marketPrice, matchAmount, makerOrder)
		//takerOrder.SubAmount(matchAmount)
		//makerOrder.SubAmount(matchAmount)
		//if makerOrder.Amount().Cmp(big.NewInt(0)) == 0 {
		//	makerBook.Remove(makerOrderElem)
		//}
		//if takerOrder.Amount().Cmp(big.NewInt(0)) == 0 {
		//	break
		//}
	}
	//if takerOrder.Amount.Cmp(big.NewInt(0)) > 0 {
	//	takerBook.Insert(takerOrder)
	//}
	return matchResult, nil
}
