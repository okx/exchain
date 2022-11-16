package dydx

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/stretchr/testify/require"
)

//{
//"65": {
//"PerpetualV1": "0x85574c0114F5387eaE45e83dE515e55f667180F7",
//"PerpetualProxy": "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
//"P1FundingOracle": "0x010579e42d4f9aE141717e180D51d8F22145f515",
//"P1MakerOracle": "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
//"MarginToken": "0xC87EF36830A0D94E42bB2D82a0b2bB939368b10B",
//"P1Orders": "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
//"P1Liquidation": "0xbe10Af4f5f2408E583259EAa3dDc41905ea195A0",
//"P1Deleveraging": "0x4817C9f094158F4134d87cdb815D61902c71330c",
//"P1CurrencyConverterProxy": "0x43F1700F40276B8F7E2B2de2bd99f6A424E13376",
//"P1WethProxy": "0xD9604832D9c966d485FDcF368198F4Ea01B09149",
//"P1LiquidatorProxy": "0xC5Df7ccfCf552c76B1FD5D00086565d77f0BC6Db"
//}
//}

var privKeyCaptain = "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
var privKeyAlice = "e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d"
var privKeyBob = "75dee45fc7b2dd69ec22dc6a825a2d982aee4ca2edd42c53ced0912173c4a788"
var privKeyTuring = "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3"
var privKeyDevnetSuper = "824c346a2b5fa81768c75408202493a9cb0a7f5879ff4988d23da2c6b1afb9cf"

var addrCaptain = privKeyToAddress(privKeyCaptain)
var addrBob = privKeyToAddress(privKeyBob)
var addrTuring = privKeyToAddress(privKeyTuring)
var addrAlice = privKeyToAddress(privKeyAlice)
var addrDevnetSuper = privKeyToAddress(privKeyDevnetSuper)

// operator : 0xfefac29bfa769d8a6c17b685816dadbd30e3f395e997ed955a5461914be75ed5

var config = DydxConfig{
	// PrivKeyHex:                 "fefac29bfa769d8a6c17b685816dadbd30e3f395e997ed955a5461914be75ed5",
	PrivKeyHex:                 privKeyTuring,
	ChainID:                    "65",
	EthHttpRpcUrl:              "https://exchaintestrpc.okex.org",
	PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
	P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
	P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
}

func TestDepositOKT(t *testing.T) {
	fmt.Println(addrBob, addrTuring, addrAlice)
}

func privKeyToAddress(privKeyHex string) common.Address {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		panic(err)
	}
	return crypto.PubkeyToAddress(privKey.PublicKey)
}

type testTool struct {
	*testing.T
}

func TestTransfer(t *testing.T) {
	var config = DydxConfig{
		PrivKeyHex:                 privKeyDevnetSuper,
		ChainID:                    "64",
		EthHttpRpcUrl:              "http://52.199.88.250:26659",
		PerpetualV1ContractAddress: "0xbc0Bf2Bf737344570c02d8D8335ceDc02cECee71",
		P1OrdersContractAddress:    "0x632D131CCCE01206F08390cB66D1AdEf9b264C61",
		P1MakerOracleAddress:       "0xF306F8B7531561d0f92BA965a163B6C6d422ade1",
	}
	book := NewDepthBook()
	me, err := NewMatchEngine(nil, nil, book, config, nil, nil)
	require.NoError(t, err)

	toAddrs := []common.Address{
		addrAlice, addrBob, addrCaptain, addrTuring,
	}

	nonce, err := me.httpCli.NonceAt(context.Background(), me.from, nil)
	require.NoError(t, err)

	gp, err := me.httpCli.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	for _, to := range toAddrs {
		tx := evmtypes.NewTransaction(nonce, to, new(big.Int).Mul(big.NewInt(1000000000000000000), big.NewInt(1000)), 21000, gp, nil)
		tx, err = me.txOps.Signer(me.from, tx)
		require.NoError(t, err)
		err = me.httpCli.SendTransaction(context.Background(), tx)
		require.NoError(t, err)
		nonce++
	}
}

func TestMatch(t *testing.T) {
	tool := &testTool{T: t}

	book := NewDepthBook()
	me, err := NewMatchEngine(nil, nil, book, config, nil, nil)
	require.NoError(t, err)

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(100, 100, true), nil))
	// buy
	// 100, 100
	// ------
	// sell

	// taker filled 100,10
	// maker filled
	// 100,10
	mr, err := me.Match(newTestOrder(100, 10, false), nil)
	// buy
	// 100, 90
	// ------
	// sell
	require.NoError(t, err)
	require.Equal(t, 1, len(mr.MatchedRecords))
	require.Equal(t, "100", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "10", mr.MatchedRecords[0].Fill.Amount.String())

	require.Equal(t, "0", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "10", mr.TakerOrder.FrozenAmount.String())
	require.Equal(t, "90", mr.MatchedRecords[0].Maker.LeftAmount.String())

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(120, 5, false), nil))
	// buy
	// 100, 90
	// ------
	// 120, 5
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(110, 6, false), nil))
	// buy
	// 100, 90
	// ------
	// 110, 6
	// 120, 5
	// sell

	// taker filled 121,10
	// maker filled
	// 110, 6
	// 120, 5
	mr, err = me.Match(newTestOrder(121, 10, true), nil)
	// buy
	// 100, 90
	// ------
	// 120, 1
	// sell
	require.NoError(t, err)
	require.Equal(t, 2, len(mr.MatchedRecords))
	require.Equal(t, "110", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "6", mr.MatchedRecords[0].Fill.Amount.String())
	require.Equal(t, "0", mr.MatchedRecords[0].Maker.LeftAmount.String())

	require.Equal(t, "120", mr.MatchedRecords[1].Fill.Price.String())
	require.Equal(t, "4", mr.MatchedRecords[1].Fill.Amount.String())
	require.Equal(t, "1", mr.MatchedRecords[1].Maker.LeftAmount.String())

	require.Equal(t, "0", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "10", mr.TakerOrder.FrozenAmount.String())

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(99, 10, true), nil))
	// buy
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(97, 100, true), nil))
	// buy
	// 97, 100
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(98, 15, true), nil))
	// buy
	// 97, 100
	// 98, 15
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// taker filled 98,115
	// maker filled
	// 100, 90
	// 99, 10
	// 98, 15
	mr, err = me.Match(newTestOrder(98, 130, false), nil)
	// buy
	// 97, 100
	// ------
	// 98, 15
	// 120, 1
	// sell
	require.NoError(t, err)
	require.Equal(t, 3, len(mr.MatchedRecords))
	require.Equal(t, "100", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "90", mr.MatchedRecords[0].Fill.Amount.String())

	require.Equal(t, "99", mr.MatchedRecords[1].Fill.Price.String())
	require.Equal(t, "10", mr.MatchedRecords[1].Fill.Amount.String())

	require.Equal(t, "98", mr.MatchedRecords[2].Fill.Price.String())
	require.Equal(t, "15", mr.MatchedRecords[2].Fill.Amount.String())

	require.Equal(t, "15", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "115", mr.TakerOrder.FrozenAmount.String())
}

func TestBalance(t *testing.T) {
	var config = Config
	book := NewDepthBook()
	me, err := NewMatchEngine(nil, nil, book, config, nil, nil)
	require.NoError(t, err)

	banlance, err := me.contracts.PerpetualV1.GetAccountBalance(nil, addrBob)
	require.NoError(t, err)
	t.Logf("bob balance: %v", banlance)

	banlance, err = me.contracts.PerpetualV1.GetAccountBalance(nil, addrAlice)
	require.NoError(t, err)
	t.Logf("alice balance: %v", banlance)

	banlance, err = me.contracts.PerpetualV1.GetAccountBalance(nil, addrTuring)
	require.NoError(t, err)
	t.Logf("turing balance: %v", banlance)

	banlance, err = me.contracts.PerpetualV1.GetAccountBalance(nil, addrCaptain)
	require.NoError(t, err)
	t.Logf("captain balance: %v", banlance)
}

func TestDeposit(t *testing.T) {
	var config = Config
	book := NewDepthBook()
	me, err := NewMatchEngine(nil, nil, book, config, nil, nil)
	require.NoError(t, err)

	price, err := me.contracts.P1MakerOracle.GetPrice(&bind.CallOpts{
		From: common.HexToAddress(config.PerpetualV1ContractAddress),
	})
	require.NoError(t, err)
	t.Logf("price: %v", price)

	addr, err := me.contracts.PerpetualV1.GetTokenContract(nil)
	require.NoError(t, err)
	t.Logf("token contract: %v", addr.Hex())

	erc20c, err := contracts.NewTestToken(addr, me.httpCli)
	require.NoError(t, err)

	_ = erc20c

	accounts := []struct {
		Address common.Address
		PrivKey string
	}{
		{
			addrAlice, privKeyAlice,
		},
		{
			addrBob, privKeyBob,
		},
		{
			addrTuring, privKeyTuring,
		},
		{
			addrCaptain, privKeyCaptain,
		},
	}

	for i, user := range accounts {
		banlance, err := me.contracts.PerpetualV1.GetAccountBalance(nil, user.Address)
		require.NoError(t, err)
		t.Logf("acc %d balance: %v", i, banlance)

		b, err := erc20c.BalanceOf(nil, user.Address)
		require.NoError(t, err)
		t.Logf("erc20 balance: %v", b)

		priv, err := crypto.HexToECDSA(user.PrivKey)
		txOps, _ := bind.NewKeyedTransactorWithChainID(priv, me.chainID)
		txOps.GasLimit = 1000000
		tx, err := erc20c.Approve(txOps, me.contracts.Addresses.PerpetualV1, big.NewInt(math.MaxInt))
		require.NoError(t, err)
		t.Logf("approve tx: %v", tx.Hash().Hex())

		time.Sleep(3 * time.Second)

		privAdmin, err := crypto.HexToECDSA(config.PrivKeyHex)
		adminTxOps, _ := bind.NewKeyedTransactorWithChainID(privAdmin, me.chainID)
		adminTxOps.GasLimit = 1000000
		tx, err = erc20c.Mint(adminTxOps, user.Address, big.NewInt(1000_0000_0000))
		require.NoError(t, err)
		t.Logf("mint tx: %v", tx.Hash().Hex())

		tx, err = me.contracts.PerpetualV1.Deposit(txOps, user.Address, big.NewInt(1000_0000_0000))
		require.NoError(t, err)
		t.Logf("transfer tx: %v", tx.Hash().Hex())

		time.Sleep(3 * time.Second)
	}

	for i, user := range accounts {
		banlance, err := me.contracts.PerpetualV1.GetAccountBalance(nil, user.Address)
		require.NoError(t, err)
		t.Logf("acc %d balance: %v", i, banlance)

	}
}

//func TestTransaction(t *testing.T) {
//	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "match")
//	var options []log.Option
//	options = append(options, log.AllowDebug())
//	logger = log.NewFilter(logger, options...)
//
//	book := NewDepthBook()
//	me, err := NewMatchEngine(nil, book, config, nil, logger)
//	require.NoError(t, err)
//
//	tool := &testTool{T: t}
//
//	// order1
//	price, ok := big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order1 := newTestBigOrder(price, big.NewInt(1), true, addrCaptain, privKeyCaptain)
//	// no match
//	tool.requireNoMatch(me.MatchAndTrade(order1))
//
//	// order2
//	price, ok = big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order3 := newTestBigOrder(price, big.NewInt(1), false, addrBob, privKeyBob)
//
//	mr, err := me.matchAndTrade(order3)
//	require.NoError(t, err)
//	require.NotNil(t, mr.Tx)
//	require.Nil(t, mr.OnChain)
//
//	tx := mr.Tx
//	txBz, err := tx.MarshalBinary()
//	require.NoError(t, err)
//
//	t.Logf("txBz: %v", hex.EncodeToString(txBz))
//
//	var evmTx *msgEthereumTx = new(msgEthereumTx)
//	err = rlp.DecodeBytes(txBz, evmTx.Data)
//	require.NoError(t, err)
//
//	rawBz, err := rlp.EncodeToBytes(evmTx)
//	require.NoError(t, err)
//	t.Logf("rawBz: %v", hex.EncodeToString(rawBz))
//}

type msgEthereumTx struct {
	Data TxData
}

//func TestMatchAndTrade(t *testing.T) {
//	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "match")
//	var options []log.Option
//	options = append(options, log.AllowDebug())
//	logger = log.NewFilter(logger, options...)
//
//	tool := &testTool{T: t}
//
//	book := NewDepthBook()
//	me, err := NewMatchEngine(nil, book, config, nil, logger)
//	require.NoError(t, err)
//
//	// order1
//	price, ok := big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order1 := newTestBigOrder(price, big.NewInt(1), true, addrCaptain, privKeyCaptain)
//	// no match
//	tool.requireNoMatch(me.MatchAndTrade(order1))
//
//	// order2
//	price, ok = big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order2 := newTestBigOrder(price, big.NewInt(1), true, addrAlice, privKeyAlice)
//	// no match
//	tool.requireNoMatch(me.MatchAndTrade(order2))
//
//	// order3
//	price, ok = big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order3 := newTestBigOrder(price, big.NewInt(3), false, addrBob, privKeyBob)
//
//	mr, err := me.MatchAndTrade(order3)
//	require.NoError(t, err)
//	require.Equal(t, 2, len(mr.MatchedRecords))
//
//	isOnChain := <-mr.OnChain
//	require.True(t, isOnChain)
//
//	// order4
//	price, ok = big.NewInt(0).SetString("18200000000000000000000", 10)
//	require.True(t, ok)
//	order4 := newTestBigOrder(price, big.NewInt(1), true, addrCaptain, privKeyCaptain)
//
//	mr, err = me.MatchAndTrade(order4)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(mr.MatchedRecords))
//
//	isOnChain = <-mr.OnChain
//	require.True(t, isOnChain)
//}

func TestTxReceipt(t *testing.T) {
	book := NewDepthBook()
	me, err := NewMatchEngine(nil, nil, book, config, nil, nil)
	require.NoError(t, err)
	rec, err := me.httpCli.TransactionReceipt(context.Background(),
		common.HexToHash("0xf6cc010deb24d3c78f5d34119541d5f5417f06da5f88989f0c34b858370e0d52"),
	)
	require.NoError(t, err)
	t.Logf("rec: %v", rec)
}

func (tool *testTool) requireNoMatch(mr *MatchResult, err error) {
	require.NoError(tool, err)
	if mr != nil {
		require.Equal(tool, 0, len(mr.MatchedRecords))
	}
}

func newTestOrder(price, amount uint64, isBuy bool) *WrapOrder {
	o := &WrapOrder{}
	o.LimitPrice = big.NewInt(0).SetUint64(price)
	o.Amount = big.NewInt(0).SetUint64(amount)
	o.LeftAmount = big.NewInt(0).SetUint64(amount)
	o.FrozenAmount = big.NewInt(0)
	o.TriggerPrice = big.NewInt(0)
	o.LimitFee = big.NewInt(0)
	// time.Now().Unix()*2 to avoid to be pruned
	// rand.Int63() to avoid repeated orderHash
	o.Expiration = big.NewInt(time.Now().Unix()*2 + rand.Int63())
	if isBuy {
		o.Flags[31] = 1
	}
	return o
}

func newTestBigOrder(price, amount *big.Int, isBuy bool, maker common.Address, privKey string) *WrapOrder {
	o := &WrapOrder{}
	o.LimitPrice = big.NewInt(0).Set(price)
	o.Amount = big.NewInt(0).Set(amount)
	o.LeftAmount = big.NewInt(0).Set(amount)
	o.FrozenAmount = big.NewInt(0)
	o.TriggerPrice = big.NewInt(0)
	o.LimitFee = big.NewInt(0)
	// time.Now().Unix()*2 to avoid to be pruned
	// rand.Int63() to avoid repeated orderHash
	o.Expiration = big.NewInt(time.Now().Unix()*2 + rand.Int63())
	if !isBuy {
		o.Flags[31] = 1
	}

	o.Maker = maker
	chainid, err := strconv.Atoi(config.ChainID)
	if err != nil {
		panic(err)
	}
	sig, err := signOrder(o.P1Order, privKey, int64(chainid), config.P1OrdersContractAddress)
	if err != nil {
		panic(err)
	}
	o.Sig = sig

	return o
}
