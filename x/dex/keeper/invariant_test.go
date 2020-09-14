package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/dex/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAccountInvariant(t *testing.T) {

	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	accounts := testInput.TestAddrs
	keeper.SetParams(ctx, *types.DefaultParams())

	builtInTP := GetBuiltInTokenPair()
	builtInTP.Owner = accounts[0]
	err := keeper.SaveTokenPair(ctx, builtInTP)
	require.Nil(t, err)

	// deposit xxb_okt 100 okt
	depositMsg := types.NewMsgDeposit(builtInTP.Name(),
		sdk.NewDecCoin(builtInTP.QuoteAssetSymbol, sdk.NewInt(100)), accounts[0])

	err = keeper.Deposit(ctx, builtInTP.Name(), depositMsg.Depositor, depositMsg.Amount)
	require.Nil(t, err)

	// module acount balance 100okt
	// xxb_okt deposits 100 okt. withdraw info 0 okt
	invariant := ModuleAccountInvariant(keeper, keeper.supplyKeeper)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// withdraw xxb_okt 50 okt
	WithdrawMsg := types.NewMsgWithdraw(builtInTP.Name(),
		sdk.NewDecCoin(builtInTP.QuoteAssetSymbol, sdk.NewInt(50)), accounts[0])

	err = keeper.Withdraw(ctx, builtInTP.Name(), WithdrawMsg.Depositor, WithdrawMsg.Amount)
	require.Nil(t, err)

	// module acount balance 100okt
	// xxb_okt deposits 50 okt. withdraw info 50 okt
	invariant = ModuleAccountInvariant(keeper, keeper.supplyKeeper)
	_, broken = invariant(ctx)
	require.False(t, broken)

}
