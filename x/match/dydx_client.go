package match

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/libs/dydx"
)

type DydxClient struct {
	contracts *dydx.Contracts

	privKeyHex string
	privKey    *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	from       common.Address

	chainID *big.Int
	ethCli  *ethclient.Client

	perpetualV1EventCh  chan types.Log
	perpetualV1EventErr error

	closeCh chan struct{}
}

func (c *DydxClient) Stop() {
	close(c.closeCh)
}

func (c *DydxClient) Trade() {

}

func (c *DydxClient) ethEventRoutine(sub ethereum.Subscription) {
	for {
		select {
		case log := <-c.perpetualV1EventCh:
			if log.Address != c.contracts.PerpetualV1Address {
				continue
			}
		case err := <-sub.Err():
			c.perpetualV1EventErr = err
			return
		case <-c.closeCh:
			sub.Unsubscribe()
			return
		}
	}
}

func NewDydxClient(chainID *big.Int, ethRpcUrl,
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

	client.contracts, err = dydx.NewContracts(
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

	client.closeCh = make(chan struct{})
	client.perpetualV1EventCh = make(chan types.Log, 512)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(perpetualV1ContractAddress)},
	}
	sub, err := client.ethCli.SubscribeFilterLogs(context.Background(), query, client.perpetualV1EventCh)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe dydx perpetualV1 event, err: %w", err)
	}

	go client.ethEventRoutine(sub)

	return &client, nil
}
