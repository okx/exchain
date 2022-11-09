package dydx

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"testing"
	"time"
)

func TestCancel(t *testing.T) {
	orderBytes, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000004be4e7267b6ae0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bbe4733d85bc2b90682147779da49cab38c0aa1f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063732962")
	require.NoError(t, err)
	var order P1Order
	err = order.DecodeFrom(orderBytes)
	require.NoError(t, err)

	privHex := "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
	priv, err := crypto.HexToECDSA(privHex)
	require.NoError(t, err)
	_ = priv

	addr := crypto.PubkeyToAddress(priv.PublicKey)
	_ = addr

	cli, err := ethclient.Dial(Config.EthHttpRpcUrl)
	require.NoError(t, err)
	p1Orders, err := contracts.NewP1Orders(common.HexToAddress(Config.P1OrdersContractAddress), cli)
	require.NoError(t, err)
	_ = p1Orders

	txOps, err := bind.NewKeyedTransactorWithChainID(priv, big.NewInt(65))
	//txOps.NoSend = true
	tx, err := p1Orders.CancelOrder(txOps, order.P1OrdersOrder)
	require.NoError(t, err)
	fmt.Printf("%#v\n", tx)

	fmt.Println(tx.Hash())
}

func TestWithdraw(t *testing.T) {
	privHex := "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
	priv, err := crypto.HexToECDSA(privHex)
	require.NoError(t, err)
	_ = priv

	addr := crypto.PubkeyToAddress(priv.PublicKey)
	_ = addr

	cli, err := ethclient.Dial(Config.EthHttpRpcUrl)
	require.NoError(t, err)
	perpetualV1, err := contracts.NewPerpetualV1(common.HexToAddress(Config.PerpetualV1ContractAddress), cli)
	require.NoError(t, err)

	tokenAddr, err := perpetualV1.GetTokenContract(nil)
	require.NoError(t, err)
	t.Logf("token contract: %s\n", tokenAddr)

	erc20c, err := contracts.NewTestToken(tokenAddr, cli)
	require.NoError(t, err)

	tokenBalance, err := erc20c.BalanceOf(nil, addr)
	require.NoError(t, err)
	log.Printf("erc20 balance: %s\n", tokenBalance)

	balance, err := perpetualV1.GetAccountBalance(nil, addr)
	require.NoError(t, err)
	log.Printf("balance of %s: %v\n", addr, balance)

	txOps, err := bind.NewKeyedTransactorWithChainID(priv, big.NewInt(65))
	//txOps.NoSend = true
	tx, err := perpetualV1.Withdraw(txOps, addr, addr, big.NewInt(1))
	require.NoError(t, err)
	log.Println("tx hash:", tx.Hash())
	time.Sleep(5 * time.Second)

	tokenBalance, err = erc20c.BalanceOf(nil, addr)
	require.NoError(t, err)
	log.Printf("erc20 balance: %s\n", tokenBalance)

	balance, err = perpetualV1.GetAccountBalance(nil, addr)
	require.NoError(t, err)
	log.Printf("balance of %s: %v\n", addr, balance)
}
