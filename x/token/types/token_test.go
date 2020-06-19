package types

import (
	"encoding/json"
	"testing"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestAccountResponse(t *testing.T) {
	accResp := AccountResponse{
		Address:    "address",
		Currencies: []CoinInfo{},
	}
	accResp1 := NewAccountResponse("address")
	require.EqualValues(t, accResp, accResp1)
}

func TestCoinInfo(t *testing.T) {
	coinInfo := CoinInfo{
		Symbol:    "btc",
		Available: "1000001",
		Locked:    "8888",
	}

	coinInfo1 := NewCoinInfo("btc", "1000001", "8888")
	require.EqualValues(t, coinInfo, *coinInfo1)
}

func TestCurrency(t *testing.T) {
	testCase := []struct {
		currency    Currency
		expectedStr string
	}{
		{Currency{
			Description: "my currency",
			Symbol:      common.NativeToken,
			TotalSupply: sdk.NewDec(10000000),
		}, `{"description":"my currency","symbol":"` + common.NativeToken + `","total_supply":"10000000.00000000"}`},
		{Currency{
			Description: common.NativeToken,
			Symbol:      common.NativeToken,
			TotalSupply: sdk.NewDec(10000),
		}, `{"description":"` + common.NativeToken + `","symbol":"` + common.NativeToken + `","total_supply":"10000.00000000"}`},
	}
	for _, currencyCase := range testCase {
		b, err := json.Marshal(currencyCase.currency)
		require.Nil(t, err)
		require.EqualValues(t, string(b), currencyCase.currency.String())
		require.EqualValues(t, currencyCase.expectedStr, currencyCase.currency.String())
	}
}

func TestToken(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("okchain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3")
	require.Nil(t, err)

	testCase := []struct {
		token       Token
		expectedStr string
	}{
		{Token{
			Description:         "my token",
			Symbol:              common.NativeToken,
			OriginalSymbol:      common.NativeToken,
			WholeName:           "btc",
			OriginalTotalSupply: sdk.NewDec(1000000),
			TotalSupply:         sdk.ZeroDec(),
			Type:                0,
			Owner:               nil,
			Mintable:            false,
		}, `{"description":"my token","symbol":"` + common.NativeToken + `","original_symbol":"` + common.NativeToken + `","whole_name":"btc","original_total_supply":"1000000.00000000","total_supply":"0.00000000","type":0,"owner":"","mintable":false}`},
		{Token{
			Description:         "okblockchain coin",
			Symbol:              common.NativeToken,
			OriginalSymbol:      common.NativeToken,
			WholeName:           "ok coin",
			OriginalTotalSupply: sdk.NewDec(1000000000),
			TotalSupply:         sdk.ZeroDec(),
			Type:                0,
			Owner:               addr,
			Mintable:            true,
		}, `{"description":"okblockchain coin","symbol":"` + common.NativeToken + `","original_symbol":"` + common.NativeToken + `","whole_name":"ok coin","original_total_supply":"1000000000.00000000","total_supply":"0.00000000","type":0,"owner":"okchain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3","mintable":true}`},
	}
	for _, tokenCase := range testCase {
		b, err := json.Marshal(tokenCase.token)
		require.Nil(t, err)
		require.EqualValues(t, string(b), tokenCase.token.String())
		require.EqualValues(t, tokenCase.expectedStr, tokenCase.token.String())
	}
}

func TestKeys(t *testing.T) {
	symbol := common.NativeToken
	b := GetTokenAddress(symbol)
	require.EqualValues(t, b, append(TokenKey, []byte(symbol)...))

	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	b = GetLockAddress(addr)
	require.EqualValues(t, b, append(LockKey, addr.Bytes()...))

	//asset := "btc"
	//b = keyDexListAsset(asset)
	//require.EqualValues(t, b, []byte(fmt.Sprintf("asset:%s", asset)))
}
