package keeper

import (
	"testing"

	"github.com/okex/okexchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestHooks(t *testing.T) {
	ctx, ak, k, _, supplyKeeper := CreateTestInputDefault(t, false, 1000)
	hook := k.Hooks()

	// test AfterValidatorCreated
	hook.AfterValidatorCreated(ctx, valOpAddr1)
	require.True(t, k.GetValidatorAccumulatedCommission(ctx, valOpAddr1).IsZero())

	// test AfterValidatorRemoved
	acc := ak.GetAccount(ctx, supplyKeeper.GetModuleAddress(types.ModuleName))
	err := acc.SetCoins(NewTestSysCoins(123, 2))
	require.NoError(t, err)
	ak.SetAccount(ctx, acc)
	k.SetValidatorAccumulatedCommission(ctx, valOpAddr1, NewTestSysCoins(123,2))
	hook.AfterValidatorRemoved(ctx, nil, valOpAddr1)
	require.True(t, ctx.KVStore(k.storeKey).Get(valOpAddr1) == nil)

	// test to promote the coverage
	hook.AfterValidatorDestroyed(ctx, valConsAddr1, valOpAddr1)
	hook.BeforeValidatorModified(ctx, valOpAddr1)
	hook.AfterValidatorBonded(ctx, valConsAddr1, valOpAddr1)
	hook.AfterValidatorBeginUnbonding(ctx, valConsAddr1, valOpAddr1)
}
