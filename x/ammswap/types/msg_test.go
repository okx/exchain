package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/okex/okexchain/x/common"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const addrStr = "1212121212121212123412121212121212121234"

func testCode(t *testing.T, err sdk.Error, expectedCode uint32) {
	if expectedCode != 0 {
		require.NotNil(t, err)
	}else {
		require.Nil(t, err)
	}
}

func TestMsgCreateExchange(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	testToken := InitPoolToken(TestBasePooledToken)
	msg := NewMsgCreateExchange(testToken.Symbol, TestQuotePooledToken, addr)
	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "create_exchange", msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgCreateExchange{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)

	expectTokenPair := TestBasePooledToken + "_" + TestQuotePooledToken
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPairName())
}

func TestMsgCreateExchangeInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	tests := []struct {
		testCase         string
		symbol0           string
		symbol1           string
		addr             sdk.AccAddress
		exceptResultCode uint32
	}{
		{"success", "aaa", common.NativeToken, addr, sdk.CodeOK},
		{"success", "aaa", "bbb", addr, sdk.CodeOK},
		{"success", "bbb", "aaa", addr, sdk.CodeOK},
		{"nil addr", "aaa", common.NativeToken, nil, sdk.CodeInvalidAddress},
		{"invalid token", "1ab",common.NativeToken, addr, sdk.CodeInvalidCoins},
		{"invalid token", common.NativeToken, common.NativeToken, addr, sdk.CodeInvalidCoins},
		//{"The lexicographic order of BaseTokenName must be less than QuoteTokenName", "xxb", addr, sdk.CodeUnknownRequest},

	}
	for i, testCase := range tests {
		msg := NewMsgCreateExchange(testCase.symbol0, testCase.symbol1, testCase.addr)
		err := msg.ValidateBasic()
		fmt.Println(i, err)
		testCode(t, err, testCase.exceptResultCode)
	}
}

func TestMsgAddLiquidity(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minLiquidity := sdk.NewDec(1)
	maxBaseAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(10000))
	quoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(10000))
	deadLine := time.Now().Unix()
	msg := NewMsgAddLiquidity(minLiquidity, maxBaseAmount, quoteAmount, deadLine, addr)
	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "add_liquidity", msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgAddLiquidity{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)

	expectTokenPair := TestBasePooledToken + "_" + TestQuotePooledToken
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPairName())
}

func TestMsgAddLiquidityInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minLiquidity := sdk.NewDec(1)
	maxBaseAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(10000))
	quoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(10000))
	notPositiveQuoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0))
	invalidMaxBaseAmount := sdk.NewDecCoinFromDec("bsa", sdk.NewDec(10000))
	invalidMaxBaseAmount.Denom = "1add"
	invalidQuoteAmount := sdk.NewDecCoinFromDec("bsa", sdk.NewDec(10000))
	invalidQuoteAmount.Denom = "1dfdf"
	notNativeQuoteAmount := sdk.NewDecCoinFromDec("abc", sdk.NewDec(10000))
	deadLine := time.Now().Unix()

	tests := []struct {
		testCase         string
		minLiquidity     sdk.Dec
		maxBaseAmount    sdk.SysCoin
		quoteAmount      sdk.SysCoin
		deadLine         int64
		addr             sdk.AccAddress
		exceptResultCode uint32
	}{
		{"success", minLiquidity, maxBaseAmount, quoteAmount, deadLine, addr, sdk.CodeOK},
		{"tokens must be positive", minLiquidity, maxBaseAmount, notPositiveQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MaxBaseAmount", minLiquidity, invalidMaxBaseAmount, quoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid QuoteAmount", minLiquidity, maxBaseAmount, invalidQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"success(quote token supports any type of tokens)", minLiquidity, maxBaseAmount, notNativeQuoteAmount, deadLine, addr, sdk.CodeOK},
		{"empty sender", minLiquidity, maxBaseAmount, quoteAmount, deadLine, nil, sdk.CodeInvalidAddress},
		{"invalid token", minLiquidity, maxBaseAmount, maxBaseAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"The lexicographic order of BaseTokenName must be less than QuoteTokenName", minLiquidity, quoteAmount, maxBaseAmount, deadLine, addr, sdk.CodeUnknownRequest},
	}
	for i, testCase := range tests {
		fmt.Println(testCase.testCase)
		msg := NewMsgAddLiquidity(testCase.minLiquidity, testCase.maxBaseAmount, testCase.quoteAmount, testCase.deadLine, testCase.addr)
		err := msg.ValidateBasic()
		fmt.Println(i, err)
		testCode(t, err, testCase.exceptResultCode)
	}
}

func TestMsgRemoveLiquidity(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	liquidity, err := sdk.NewDecFromStr("0.01")
	require.Nil(t, err)
	minBaseAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	minQuoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	msg := NewMsgRemoveLiquidity(liquidity, minBaseAmount, minQuoteAmount, deadLine, addr)

	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "remove_liquidity", msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgRemoveLiquidity{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)

	expectTokenPair := TestBasePooledToken + "_" + TestQuotePooledToken
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPairName())
}

func TestMsgRemoveLiquidityInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	liquidity, err := sdk.NewDecFromStr("0.01")
	require.Nil(t, err)
	minBaseAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	minQuoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	notPositiveLiquidity := sdk.NewDec(0)
	invalidMinBaseAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	invalidMinBaseAmount.Denom = "1sss"
	invalidMinQuoteAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(1))
	invalidMinQuoteAmount.Denom = "1sss"
	notNativeQuoteAmount := sdk.NewDecCoinFromDec("sss", sdk.NewDec(1))

	tests := []struct {
		testCase         string
		liquidity        sdk.Dec
		minBaseAmount    sdk.SysCoin
		minQuoteAmount   sdk.SysCoin
		deadLine         int64
		addr             sdk.AccAddress
		exceptResultCode uint32
	}{
		{"success", liquidity, minBaseAmount, minQuoteAmount, deadLine, addr, sdk.CodeOK},
		{"empty sender", liquidity, minBaseAmount, minQuoteAmount, deadLine, nil, sdk.CodeInvalidAddress},
		{"coins must be positive", notPositiveLiquidity, minBaseAmount, minQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MinBaseAmount", liquidity, invalidMinBaseAmount, minQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MinQuoteAmount", liquidity, minBaseAmount, invalidMinQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"success(quote token supports any type of tokens)", liquidity, minBaseAmount, notNativeQuoteAmount, deadLine, addr, sdk.CodeOK},
		{"invalid token", liquidity, minBaseAmount, minBaseAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"The lexicographic order of BaseTokenName must be less than QuoteTokenName", liquidity, minQuoteAmount, minBaseAmount, deadLine, addr, sdk.CodeUnknownRequest},


	}
	for _, testCase := range tests {
		msg := NewMsgRemoveLiquidity(testCase.liquidity, testCase.minBaseAmount, testCase.minQuoteAmount, testCase.deadLine, testCase.addr)
		err := msg.ValidateBasic()
		testCode(t, err, testCase.exceptResultCode)
	}
}

func TestMsgTokenToToken(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	soldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	msg := NewMsgTokenToToken(soldTokenAmount, minBoughtTokenAmount, deadLine, addr, addr)

	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, TypeMsgTokenSwap, msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgTokenToToken{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
	expectTokenPair := TestBasePooledToken + "_" + TestQuotePooledToken
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPairName())
	msg = NewMsgTokenToToken(minBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr)
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPairName())
}

func TestMsgTokenToTokenInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	soldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	invalidMinBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	invalidMinBoughtTokenAmount.Denom = "1aaa"
	invalidSoldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	invalidSoldTokenAmount.Denom = "1sdf"
	invalidSoldTokenAmount2 := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0))
	invalidSoldTokenAmount2.Denom = "aaa"
	notNativeSoldTokenAmount := sdk.NewDecCoinFromDec("abc", sdk.NewDec(2))

	tests := []struct {
		testCase             string
		minBoughtTokenAmount sdk.SysCoin
		soldTokenAmount      sdk.SysCoin
		deadLine             int64
		recipient            sdk.AccAddress
		addr                 sdk.AccAddress
		exceptResultCode     uint32
	}{
		{"success", minBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr, sdk.CodeOK},
		{"empty sender", minBoughtTokenAmount, soldTokenAmount, deadLine, addr, nil, sdk.CodeInvalidAddress},
		{"empty recipient", minBoughtTokenAmount, soldTokenAmount, deadLine, nil, addr, sdk.CodeInvalidAddress},
		{"success(both token to sell and token to buy do not contain native token)", minBoughtTokenAmount, notNativeSoldTokenAmount, deadLine, addr, addr, sdk.CodeOK},
		{"invalid SoldTokenAmount", soldTokenAmount, invalidSoldTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},
		{"invalid MinBoughtTokenAmount", invalidMinBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},
		{"invalid token", minBoughtTokenAmount, minBoughtTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},
		{"invalid SoldTokenAmount(zero)", minBoughtTokenAmount, invalidSoldTokenAmount2, deadLine, addr, addr, sdk.CodeUnknownRequest},


	}
	for _, testCase := range tests {
		msg := NewMsgTokenToToken(testCase.soldTokenAmount, testCase.minBoughtTokenAmount, testCase.deadLine, testCase.recipient, testCase.addr)
		err := msg.ValidateBasic()
		testCode(t, err, testCase.exceptResultCode)
	}
}
