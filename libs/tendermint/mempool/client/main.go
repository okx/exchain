package main

import (
	"context"
	"flag"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/mempool/dydx"
)

const (
	node             = "wss://exchaintestws.okex.org:8443"
	localNode        = "http://127.0.0.1:8545"
	GasPrice  int64  = 100000000 // 0.1 gwei
	GasLimit  uint64 = 3000000
)

var (
	chainID           = int64(65)
	orderContractAddr = common.HexToAddress("0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619")
)

var (
	privHex string
	amount  int64
	isBuy   bool
)

func main() {
	flag.StringVar(&privHex, "priv", "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17", "")
	flag.Int64Var(&amount, "amount", 1, "")
	flag.BoolVar(&isBuy, "buy", true, "")
	flag.Parse()
	priv, err := crypto.HexToECDSA(privHex)
	if err != nil {
		panic(err)
	}
	addr := crypto.PubkeyToAddress(priv.PublicKey)

	client, err := ethclient.Dial(localNode)
	if err != nil {
		panic(err)
	}

	//TODO orderBytes + signature
	order := newP1Order(amount, isBuy)
	order.Maker = addr
	price, ok := big.NewInt(0).SetString("18200000000000000000000", 10)
	if !ok {
		panic(0)
	}
	order.LimitPrice = price
	sig, err := signOrder(order, privHex, chainID, orderContractAddr.String())
	if err != nil {
		panic(err)
	}

	orderBytes, err := order.Encode()
	if err != nil {
		panic(err)
	}
	data := append(orderBytes, sig...)

	unsignedTx := types.NewTransaction(0, orderContractAddr, big.NewInt(0), GasLimit, big.NewInt(GasPrice), data)

	err = client.SendTransaction(context.Background(), unsignedTx)
	if err != nil {
		panic(err)
	}

}

func signOrder(odr dydx.P1Order, hexPriv string, chainId int64, orderContractaddr string) ([]byte, error) {
	priv, err := crypto.HexToECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	orderHash := odr.Hash2(chainId, orderContractaddr)
	signedHash := crypto.Keccak256Hash([]byte(dydx.PREPEND_DEC), orderHash[:])
	sig, err := crypto.Sign(signedHash[:], priv)
	if err != nil {
		return nil, err
	}

	sig[len(sig)-1] += 27
	sig = append(sig, 1)
	return sig, nil
}

func newP1Order(amount int64, isBuy bool) dydx.P1Order {
	odr := dydx.P1Order{
		CallType: 1,
		P1OrdersOrder: contracts.P1OrdersOrder{
			Amount:       big.NewInt(amount),
			LimitPrice:   big.NewInt(0),
			TriggerPrice: big.NewInt(0),
			LimitFee:     big.NewInt(0),
			Expiration:   big.NewInt(time.Now().Unix() * 2),
		},
	}
	if isBuy {
		odr.Flags[31] = 1
	}
	return odr
}
