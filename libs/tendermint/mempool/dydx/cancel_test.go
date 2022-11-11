package dydx

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/mempool/placeorder"
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

func TestMint(t *testing.T) {
	cli, err := ethclient.Dial(Config.EthHttpRpcUrl)
	require.NoError(t, err)
	token, err := contracts.NewTestToken(common.HexToAddress(Config.P1MarginAddress), cli)
	require.NoError(t, err)

	privAdmin, err := crypto.HexToECDSA(Config.PrivKeyHex)
	chainID, _ := new(big.Int).SetString(Config.ChainID, 10)
	adminTxOps, _ := bind.NewKeyedTransactorWithChainID(privAdmin, chainID)
	adminTxOps.GasLimit = 1000000

	tx, err := token.Mint(adminTxOps, common.HexToAddress("0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044"), big.NewInt(1000000))
	require.NoError(t, err)
	t.Logf("mint tx: %v", tx.Hash().Hex())
}

func TestPlaceOrder(t *testing.T) {
	cli, err := ethclient.Dial(Config.EthHttpRpcUrl) //"http://3.113.237.222:26659"
	require.NoError(t, err)
	caller, err := placeorder.NewPlaceorderCaller(common.HexToAddress(placeOrderContractAddr), cli)
	require.NoError(t, err)
	maker := common.HexToAddress("0xbbE4733d85bc2b90682147779DA49caB38C0aA1F")
	order := placeorder.OrdersOrder{
		Amount:       big.NewInt(1),
		LimitPrice:   big.NewInt(18200),
		TriggerPrice: big.NewInt(0),
		LimitFee:     big.NewInt(0),
		Maker:        maker,
		Expiration:   big.NewInt(time.Now().Unix() + oneWeekSeconds),
	}
	msg, err := caller.GetOrderMessage(&bind.CallOpts{From: maker}, order)
	require.NoError(t, err)
	t.Log(hex.EncodeToString(msg))

}
