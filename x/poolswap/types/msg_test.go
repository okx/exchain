package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const addrStr = "1212121212121212123412121212121212121234"

func TestMsgCreateExchange(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	testToken := InitPoolToken(TestBasePooledToken)
	msg := NewMsgCreateExchange(testToken.Symbol, addr)
	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "create_exchange", msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgCreateExchange{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
}

func TestMsgCreateExchangeInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	tests := []struct {
		testCase         string
		symbol           string
		addr             sdk.AccAddress
		exceptResultCode sdk.CodeType
	}{
		{"success", "xxx", addr, sdk.CodeOK},
		{"nil addr", "xxx", nil, sdk.CodeInvalidAddress},
		{"invalid token", "1ab", addr, sdk.CodeUnknownRequest},
	}
	for _, testCase := range tests {
		msg := NewMsgCreateExchange(testCase.symbol, testCase.addr)
		err := msg.ValidateBasic()
		if err == nil && testCase.exceptResultCode == sdk.CodeOK {
			continue
		}
		require.Equal(t, testCase.exceptResultCode, err.Code())
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
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPair())
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
		maxBaseAmount    sdk.DecCoin
		quoteAmount      sdk.DecCoin
		deadLine         int64
		addr             sdk.AccAddress
		exceptResultCode sdk.CodeType
	}{
		{"success", minLiquidity, maxBaseAmount, quoteAmount, deadLine, addr, 0},
		{"tokens must be positive", minLiquidity, maxBaseAmount, notPositiveQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MaxBaseAmount", minLiquidity, invalidMaxBaseAmount, quoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid QuoteAmount", minLiquidity, maxBaseAmount, invalidQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"quote token only supports native token", minLiquidity, maxBaseAmount, notNativeQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"empty sender", minLiquidity, maxBaseAmount, quoteAmount, deadLine, nil, sdk.CodeInvalidAddress},
	}
	for _, testCase := range tests {
		fmt.Println(testCase.testCase)
		msg := NewMsgAddLiquidity(testCase.minLiquidity, testCase.maxBaseAmount, testCase.quoteAmount, testCase.deadLine, testCase.addr)
		err := msg.ValidateBasic()
		if err == nil && testCase.exceptResultCode == sdk.CodeOK {
			continue
		}
		require.Equal(t, testCase.exceptResultCode, err.Code())
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
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPair())
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
		minBaseAmount    sdk.DecCoin
		minQuoteAmount   sdk.DecCoin
		deadLine         int64
		addr             sdk.AccAddress
		exceptResultCode sdk.CodeType
	}{
		{"success", liquidity, minBaseAmount, minQuoteAmount, deadLine, addr, 0},
		{"empty sender", liquidity, minBaseAmount, minQuoteAmount, deadLine, nil, sdk.CodeInvalidAddress},
		{"coins must be positive", notPositiveLiquidity, minBaseAmount, minQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MinBaseAmount", liquidity, invalidMinBaseAmount, minQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"invalid MinQuoteAmount", liquidity, minBaseAmount, invalidMinQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},
		{"quote token only supports native token", liquidity, minBaseAmount, notNativeQuoteAmount, deadLine, addr, sdk.CodeUnknownRequest},

	}
	for _, testCase := range tests {
		msg := NewMsgRemoveLiquidity(testCase.liquidity, testCase.minBaseAmount, testCase.minQuoteAmount, testCase.deadLine, testCase.addr)
		err := msg.ValidateBasic()
		if err == nil && testCase.exceptResultCode == sdk.CodeOK {
			continue
		}
		require.Equal(t, testCase.exceptResultCode, err.Code())
	}
}


func TestMsgTokenToNativeToken(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	soldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	msg := NewMsgTokenToNativeToken(soldTokenAmount, minBoughtTokenAmount, deadLine, addr, addr)

	require.Nil(t, msg.ValidateBasic())
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, TypeMsgTokenSwap, msg.Type())

	bytesMsg := msg.GetSignBytes()
	resMsg := &MsgTokenToNativeToken{}
	err = json.Unmarshal(bytesMsg, resMsg)
	require.Nil(t, err)
	resAddr := msg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
	expectTokenPair := TestBasePooledToken + "_" + TestQuotePooledToken
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPair())
	msg = NewMsgTokenToNativeToken(minBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr)
	require.Equal(t, expectTokenPair, msg.GetSwapTokenPair())
}

func TestMsgTokenToNativeTokenInvalid(t *testing.T) {
	addr, err := hex.DecodeString(addrStr)
	require.Nil(t, err)
	minBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	deadLine := time.Now().Unix()
	soldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	invalidMinBoughtTokenAmount := sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1))
	invalidMinBoughtTokenAmount.Denom = "1aaa"
	invalidSoldTokenAmount := sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(2))
	invalidSoldTokenAmount.Denom = "1sdf"
	notNativeSoldTokenAmount := sdk.NewDecCoinFromDec("abc", sdk.NewDec(2))

	tests := []struct {
		testCase             string
		minBoughtTokenAmount sdk.DecCoin
		soldTokenAmount      sdk.DecCoin
		deadLine             int64
		recipient            sdk.AccAddress
		addr                 sdk.AccAddress
		exceptResultCode     sdk.CodeType
	}{
		{"success", minBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr, sdk.CodeOK},
		{"empty sender", minBoughtTokenAmount, soldTokenAmount, deadLine, addr, nil, sdk.CodeInvalidAddress},
		{"empty recipient", minBoughtTokenAmount, soldTokenAmount, deadLine, nil, addr, sdk.CodeInvalidAddress},
		{"both token to sell and token to buy do not contain native token", minBoughtTokenAmount, notNativeSoldTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},
		{"invalid SoldTokenAmount", soldTokenAmount, invalidSoldTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},
		{"invalid MinBoughtTokenAmount", invalidMinBoughtTokenAmount, soldTokenAmount, deadLine, addr, addr, sdk.CodeUnknownRequest},

	}
	for _, testCase := range tests {
		msg := NewMsgTokenToNativeToken(testCase.soldTokenAmount, testCase.minBoughtTokenAmount, testCase.deadLine, testCase.recipient, testCase.addr)
		err := msg.ValidateBasic()
		if err == nil && testCase.exceptResultCode == sdk.CodeOK {
			continue
		}
		require.Equal(t, testCase.exceptResultCode, err.Code())
	}
}
