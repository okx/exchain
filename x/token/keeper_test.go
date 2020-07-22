package token

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token/types"
)

func TestKeeper_GetFeeDetailList(t *testing.T) {
	_, keeper, _ := getMockDexApp(t, 0)
	keeper.GetFeeDetailList()
}

func TestKeeper_UpdateTokenSupply(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	token := types.Token{
		Description:         "new token",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		WholeName:           "ok coin",
		OriginalTotalSupply: sdk.NewDec(10000),
		Owner:               []byte("gyl"),
		Mintable:            true,
	}
	keeper.NewToken(ctx, token)

	info := keeper.GetTokenInfo(ctx, "xxbToken")
	require.EqualValues(t, "", info.Symbol)

	info = keeper.GetTokenInfo(ctx, common.NativeToken)
	require.EqualValues(t, common.NativeToken, info.Symbol)
	tokens := keeper.GetTokensInfo(ctx)
	require.Equal(t, info, tokens[0])
	require.Equal(t, len(tokens), len(keeper.GetCurrenciesInfo(ctx)))

	name, flag := addTokenSuffix(ctx, keeper, common.NativeToken)
	require.Equal(t, true, flag)
	require.NotEqual(t, common.NativeToken, name)

	name, flag = addTokenSuffix(ctx, keeper, common.NativeToken+"@#$")
	require.Equal(t, false, flag)
	require.Equal(t, "", name)

	token = types.Token{
		Description:         "new token",
		Symbol:              common.NativeToken,
		OriginalSymbol:      common.NativeToken,
		WholeName:           "ok coin",
		OriginalTotalSupply: sdk.NewDec(10000),
		Owner:               []byte("gyl"),
		Mintable:            true,
	}
	store := ctx.KVStore(keeper.tokenStoreKey)
	store.Set(types.TokenNumberKey, []byte("0"))
	keeper.NewToken(ctx, token)
}

func TestKeeper_AddFeeDetail(t *testing.T) {
	mapp, keeper, addrs := getMockDexApp(t, 2)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.enableBackend = true

	fee, err := sdk.ParseDecCoins(fmt.Sprintf("100%s", common.NativeToken))
	require.Nil(t, err)

	feeType := "test fee type"
	keeper.AddFeeDetail(ctx, addrs[0].String(), fee, feeType, "")
	keeper.AddFeeDetail(ctx, addrs[1].String(), fee, feeType, "")

	feeList := keeper.GetFeeDetailList()
	require.Equal(t, 2, len(feeList))

	require.Equal(t, fee.String(), feeList[0].Fee)
	require.Equal(t, addrs[0].String(), feeList[0].Address)
	require.Equal(t, addrs[1].String(), feeList[1].Address)
}

func TestKeeper_LockCoins(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	expectedLockedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}

	err := keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, expectedLockedCoins, types.LockCoinsTypeQuantity)
	require.NoError(t, err)
	lockedCoins := keeper.GetLockedCoins(ctx, testAccounts[0].baseAccount.Address)
	require.EqualValues(t, expectedLockedCoins, lockedCoins)
	err = keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, expectedLockedCoins, types.LockCoinsTypeQuantity)
	require.NoError(t, err)
	BigCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100000000000)),
	}
	err = keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, BigCoins, types.LockCoinsTypeQuantity)
	require.Error(t, err)
	_, lockStoreKeyNum := keeper.getNumKeys(ctx)
	require.Equal(t, int64(1), lockStoreKeyNum)

	locks := keeper.GetAllLockedCoins(ctx)
	require.NotNil(t, locks)
	require.EqualValues(t, common.NativeToken, locks[0].Coins[0].Denom)
}

func TestKeeper_UnlockCoins(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	lockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}

	err := keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, lockCoins, types.LockCoinsTypeQuantity)
	require.NoError(t, err)

	unlockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(10)),
	}

	err = keeper.UnlockCoins(ctx, testAccounts[0].baseAccount.Address, unlockCoins, types.LockCoinsTypeQuantity)
	require.NoError(t, err)
	err = keeper.UnlockCoins(ctx, []byte(""), unlockCoins, types.LockCoinsTypeQuantity)
	require.Error(t, err)

	biglockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(999999999)),
	}

	err = keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, biglockCoins, types.LockCoinsTypeQuantity)
	require.Error(t, err)
	bigunlockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(10000000000000)),
	}
	err = keeper.UnlockCoins(ctx, testAccounts[0].baseAccount.Address, bigunlockCoins, types.LockCoinsTypeQuantity)
	require.Error(t, err)

	//expectedLockCoins
	expectedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(90)),
	}
	require.EqualValues(t, keeper.GetLockedCoins(ctx, testAccounts[0].baseAccount.Address), expectedCoins)

	err = keeper.UnlockCoins(ctx, testAccounts[0].baseAccount.Address, keeper.GetLockedCoins(ctx, testAccounts[0].baseAccount.Address), types.LockCoinsTypeQuantity)
	require.NoError(t, err)
	require.EqualValues(t, keeper.GetLockedCoins(ctx, testAccounts[0].baseAccount.Address), sdk.DecCoins(nil))
}

func TestKeeper_BurnLockedCoins(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	initCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000000000))}
	genAccs, testAccounts := CreateGenAccounts(1, initCoins)
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(initCoins))

	lockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}
	err := keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, lockCoins, types.LockCoinsTypeQuantity)
	require.NoError(t, err)
	burnLockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(5)),
	}
	err = keeper.BalanceAccount(ctx, testAccounts[0].baseAccount.Address, burnLockCoins, sdk.ZeroFee().ToCoins())
	require.Nil(t, err)

	expectedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(95)),
	}
	require.EqualValues(t, keeper.GetLockedCoins(ctx, testAccounts[0].baseAccount.Address), expectedCoins)
	err = keeper.BalanceAccount(ctx, testAccounts[0].baseAccount.Address, expectedCoins, sdk.ZeroFee().ToCoins())
	require.Nil(t, err)

	expectedCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(10000000000000)),
	}
	err = keeper.BalanceAccount(ctx, testAccounts[0].baseAccount.Address, expectedCoins, sdk.ZeroFee().ToCoins())
	require.NotNil(t, err)
}

func TestKeeper_ReceiveLockedCoins(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, testAccounts := CreateGenAccounts(1,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	lockCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}
	err := keeper.LockCoins(ctx, testAccounts[0].baseAccount.Address, lockCoins, types.LockCoinsTypeQuantity)
	require.Nil(t, err)

	err = keeper.BalanceAccount(ctx, testAccounts[0].baseAccount.Address, sdk.ZeroFee().ToCoins(), lockCoins)
	require.Nil(t, err)
	coins := keeper.GetCoins(ctx, testAccounts[0].baseAccount.Address)

	expectedCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000)),
	}
	require.EqualValues(t, expectedCoins, coins)
}

func TestKeeper_SetCoins(t *testing.T) {
	mapp, keeper, _ := getMockDexApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	genAccs, testAccounts := CreateGenAccounts(2,
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000)),
		})
	mock.SetGenesis(mapp.App, types.DecAccountArrToBaseAccountArr(genAccs))

	// issue raw token
	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}
	err := keeper.bankKeeper.SetCoins(ctx, testAccounts[0].baseAccount.Address, coins)
	require.Nil(t, err)

	require.EqualValues(t, keeper.GetCoins(ctx, testAccounts[0].baseAccount.Address), coins)

	info := keeper.GetCoinsInfo(ctx, testAccounts[0].baseAccount.Address)
	require.NotNil(t, info)

	keeper.SetParams(ctx, types.DefaultParams())
	pas := keeper.GetParams(ctx)
	require.NotNil(t, pas)

	//send coin
	coins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1)),
	}
	err1 := keeper.SendCoinsFromAccountToAccount(ctx, testAccounts[0].baseAccount.Address,
		testAccounts[1].baseAccount.Address, coins)
	require.Nil(t, err1)

	require.EqualValues(t, "1001.00000000", keeper.GetCoinsInfo(ctx,
		testAccounts[1].baseAccount.Address)[0].Available)
}
