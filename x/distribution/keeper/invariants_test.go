package keeper

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

type testInvariantParam struct {
	totalCommission sdk.SysCoins
	commissions     []sdk.SysCoins
	expected        bool
}

func getTestInvariantParams() []testInvariantParam {
	return []testInvariantParam{
		{ // when commission is zero
			nil,
			[]sdk.SysCoins{NewTestSysCoins(0, 0)},
			true,
		},
		{ // when withdraw commission failed
			NewTestSysCoins(5, 1),
			[]sdk.SysCoins{NewTestSysCoins(15, 1)},
			false,
		},
		{ // when the sum of commission is not equal to distribution account
			NewTestSysCoins(30, 1),
			[]sdk.SysCoins{NewTestSysCoins(15, 1), NewTestSysCoins(20, 1)},
			false,
		},
		{ // when the sum of commission is not equal to distribution account
			NewTestSysCoins(45, 1),
			[]sdk.SysCoins{NewTestSysCoins(15, 1), NewTestSysCoins(20, 1)},
			false,
		},
		{ // success
			NewTestSysCoins(45, 1),
			[]sdk.SysCoins{NewTestSysCoins(15, 1), NewTestSysCoins(30, 1)},
			true,
		},
	}
}

func TestInvariants(t *testing.T) {
	valOpAddrs, _, _ := GetTestAddrs()
	tests := getTestInvariantParams()
	for _, test := range tests {
		ctx, ak, keeper, sk, supplyKeeper := CreateTestInputDefault(t, false, 1000)
		acc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
		err := acc.SetCoins(test.totalCommission)
		require.NoError(t, err)
		ak.SetAccount(ctx, acc, false)
		for i, commission := range test.commissions {
			val := sk.Validator(ctx, valOpAddrs[i])
			keeper.AllocateTokensToValidator(ctx, val, commission)
		}

		var invariants [3]sdk.Invariant
		invariants[0] = NonNegativeCommissionsInvariant(keeper)
		invariants[1] = CanWithdrawInvariant(keeper)
		invariants[2] = ModuleAccountInvariant(keeper)
		count := 0
		for _, invariant := range invariants {
			if _, broken := invariant(ctx); broken {
				count++
			}
		}
		isSuccess := true
		if count != 0 {
			isSuccess = false
		}
		require.Equal(t, test.expected, isSuccess)
	}
}
