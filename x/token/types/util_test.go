package types

import (
	"testing"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestAmountToCoins(t *testing.T) {
	coinStr := "2btc,1" + common.NativeToken
	coins, err := sdk.ParseDecCoins(coinStr)
	require.Nil(t, err)
	expectedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc", sdk.NewDec(2)),
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1)),
	}
	require.EqualValues(t, expectedCoins, coins)
}

func TestStrToTransfers(t *testing.T) {
	//coinStr := `[{"to": "cosmos18ragjd23yv4ctjg3vadh43q5zf8z0hafm4qjrf", "amount": "1BNB,2BTC"},
	//{"to": "cosmos18ragjd23yv4ctjg3vadh43q5zf8z0hafm4qjrf", "amount": "1OKB,2BTC"}]`
	coinStr := `[{"to":"okexchain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3","amount":"1` + common.NativeToken + `"}]`
	coinStrError := `[{"to":"kochain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3","amount":"1` + common.NativeToken + `"}]`
	addr, err := sdk.AccAddressFromBech32("okexchain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3")
	require.Nil(t, err)
	_, err = StrToTransfers(coinStrError)
	require.Error(t, err)
	transfers, err := StrToTransfers(coinStr)
	require.Nil(t, err)
	transfer := []TransferUnit{
		{
			To: addr,
			Coins: []sdk.DecCoin{
				sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1)),
			},
		},
	}
	require.EqualValues(t, transfer, transfers)

	coinStr = `[{"to":"okexchain1dfpljpe0g0206jch32fx95lyagq3z5ws2vgwx3",amount":"1"}]`
	_, err = StrToTransfers(coinStr)
	require.Error(t, err)
}

func TestMergeCoinInfo(t *testing.T) {

	//availableCoins, freezeCoins, lockCoins sdk.DecCoins
	availableCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
		sdk.NewDecCoinFromDec("bnb", sdk.NewDec(100)),
		sdk.NewDecCoinFromDec("btc", sdk.NewDec(100)),
	}

	lockedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc", sdk.NewDec(100)),
		sdk.NewDecCoinFromDec("abc", sdk.NewDec(100)),
	}

	coinsInfo := MergeCoinInfo(availableCoins, lockedCoins)
	expectedCoinsInfo := CoinsInfo{
		CoinInfo{"abc", "0", "100.00000000"},
		CoinInfo{"bnb", "100.00000000", "0"},
		CoinInfo{"btc", "100.00000000", "100.00000000"},
		CoinInfo{common.NativeToken, "100.00000000", "0"},
	}
	require.EqualValues(t, expectedCoinsInfo, coinsInfo)
}

func TestDecAccount_String(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	dec := sdk.MustNewDecFromStr("0.2")
	decCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, dec),
	}
	decAccount := DecAccount{
		Address:       addr,
		Coins:         decCoins,
		PubKey:        pubKey,
		AccountNumber: 1,
		Sequence:      1,
	}

	expectedStr := `Account:
 Address:       ` + addr.String() + `
 Pubkey:        ` + sdk.MustBech32ifyAccPub(pubKey) + `
 Coins:         0.20000000` + common.NativeToken + `
 AccountNumber: 1
 Sequence:      1`

	decStr := decAccount.String()
	require.EqualValues(t, decStr, expectedStr)
}

func TestBaseAccountToDecAccount(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}

	baseAccount := auth.BaseAccount{
		Address:       addr,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: 1,
		Sequence:      1,
	}

	dec := sdk.MustNewDecFromStr("100.00000000")
	decCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, dec),
	}

	expectedDecAccount := DecAccount{
		Address:       addr,
		Coins:         decCoins,
		PubKey:        pubKey,
		AccountNumber: 1,
		Sequence:      1,
	}

	decAccount := BaseAccountToDecAccount(baseAccount)
	require.EqualValues(t, decAccount, expectedDecAccount)
}

func TestValidCoinName(t *testing.T) {
	coinName := "abf.s0fa"
	valid := sdk.ValidateDenom(coinName)
	require.Error(t, valid)
}

func TestValidOriginalSymbol(t *testing.T) {
	name := "abc"
	require.True(t, ValidOriginalSymbol(name))
	name = notAllowedPrefix
	require.False(t, ValidOriginalSymbol(name))
	name = notAllowedPrefix + "abc"
	require.False(t, ValidOriginalSymbol(name))
	name = notAllowedPrefix + "-abc"
	require.False(t, ValidOriginalSymbol(name))
	name = notAllowedPrefix + "-abc-1af"
	require.False(t, ValidOriginalSymbol(name))
	name = notAllowedPrefix + "1"
	require.False(t, ValidOriginalSymbol(name))
	name = "abc-1fa"
	require.False(t, ValidOriginalSymbol(name))
}

func TestValidateDenom(t *testing.T) {
	name := "abc"
	require.Nil(t, sdk.ValidateDenom(name))

	name = notAllowedPrefix + "-abc"
	require.Nil(t, sdk.ValidateDenom(name))

	name = notAllowedPrefix + "-abc-1af"
	require.Nil(t, sdk.ValidateDenom(name))

	name = notAllowedPrefix + "-abcde-aaa"
	require.Nil(t, sdk.ValidateDenom(name))

	name = notAllowedPrefix + "f-abcde-aaa"
	require.NotNil(t, sdk.ValidateDenom(name))

	name = "pool-abcde-aaa"
	require.NotNil(t, sdk.ValidateDenom(name))

	name = "token-abcde-aaa"
	require.NotNil(t, sdk.ValidateDenom(name))

	name = "pool-abc"
	require.Nil(t, sdk.ValidateDenom(name))
}