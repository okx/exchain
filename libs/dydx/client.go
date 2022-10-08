package dydx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/libs/dydx/contracts"
)

type DydxClient struct {
	contracts *Contracts

	privKeyHex string
	privKey    *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	from       common.Address

	chainID *big.Int
	ethCli  *ethclient.Client

	//perpetualV1EventCh  chan types.Log
	//perpetualV1EventErr error

	perpetualV1EventLogTradeCh  chan *contracts.PerpetualV1LogTrade
	perpetualV1EventLogTradeErr chan error

	closeCh chan struct{}
}

func (c *DydxClient) Stop() {
	close(c.closeCh)
}

func (c *DydxClient) Trade(order1, order2 *SignedOrder, amount *big.Int, price Price, fee Fee) (*types.Transaction, error) {
	op := NewTradeOperation(c.contracts)
	err := op.FillSignedOrderWithTaker(c.from.String(), order1, amount, price, fee)
	if err != nil {
		return nil, fmt.Errorf("failed to fill order1, err: %w", err)
	}
	err = op.FillSignedOrderWithTaker(c.from.String(), order2, amount, price, fee)
	if err != nil {
		return nil, fmt.Errorf("failed to fill order2, err: %w", err)
	}
	tx, err := op.Commit(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to commit, err: %w", err)
	}
	return tx, nil
}

func (c *DydxClient) GetLogTrades() <-chan *contracts.PerpetualV1LogTrade {
	return c.perpetualV1EventLogTradeCh
}

//func (c *DydxClient) ethEventRoutine(sub ethereum.Subscription) {
//	for {
//		select {
//		case log := <-c.perpetualV1EventCh:
//			if log.Address != c.contracts.PerpetualV1Address {
//				continue
//			}
//		case err := <-sub.Err():
//			c.perpetualV1EventErr = err
//			return
//		case <-c.closeCh:
//			sub.Unsubscribe()
//			return
//		}
//	}
//}

func (c *DydxClient) ethEventLogTradeRoutine(sub ethereum.Subscription) {
	for {
		select {
		// case log := <-c.perpetualV1EventLogTradeCh:
		case err := <-sub.Err():
			c.perpetualV1EventLogTradeErr <- err
			return
		case <-c.closeCh:
			sub.Unsubscribe()
			return
		}
	}
}

func NewDydxClient(chainID *big.Int, ethRpcUrl string, fromBlockNum *big.Int,
	privKey, perpetualV1ContractAddress, p1OrdersContractAddress string) (*DydxClient, error) {
	var client DydxClient
	var err error

	client.privKeyHex = privKey
	client.privKey, err = crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}
	client.pubKey = &client.privKey.PublicKey
	client.from = crypto.PubkeyToAddress(*client.pubKey)

	if chainID == nil {
		return nil, fmt.Errorf("chainID is nil")
	}
	client.chainID = chainID
	client.ethCli, err = ethclient.Dial(ethRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial eth rpc url: %s, err: %w", ethRpcUrl, err)
	}

	txOps, err := bind.NewKeyedTransactorWithChainID(client.privKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOps, err: %w", err)
	}

	client.contracts, err = NewContracts(
		common.HexToAddress(perpetualV1ContractAddress),
		common.HexToAddress(p1OrdersContractAddress),
		txOps,
		client.ethCli,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dydx contracts, err: %w", err)
	}

	//client.contracts.PerpetualV1.FilterLogTrade()
	//client.contracts.P1Orders.FilterLogOrderFilled()

	start := fromBlockNum.Uint64()
	watchOps := &bind.WatchOpts{Start: &start, Context: context.Background()}
	client.perpetualV1EventLogTradeCh = make(chan *contracts.PerpetualV1LogTrade, 128)
	sub, err := client.contracts.PerpetualV1.WatchLogTrade(watchOps, client.perpetualV1EventLogTradeCh, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to watch LogTrade, err: %w", err)
	}
	go client.ethEventLogTradeRoutine(sub)

	//client.closeCh = make(chan struct{})
	//client.perpetualV1EventCh = make(chan types.Log, 512)
	//query := ethereum.FilterQuery{
	//	Addresses: []common.Address{common.HexToAddress(perpetualV1ContractAddress)},
	//	FromBlock: fromBlockNum,
	//}
	//sub, err := client.ethCli.SubscribeFilterLogs(context.Background(), query, client.perpetualV1EventCh)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to subscribe dydx perpetualV1 event, err: %w", err)
	//}
	//
	//go client.ethEventRoutine(sub)

	return &client, nil
}
