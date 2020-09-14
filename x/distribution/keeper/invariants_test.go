package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

type testInvariantParam struct {
	totalCommission sdk.DecCoins
	commissions     []sdk.DecCoins
	expected        bool
}

func getTestInvariantParams() []testInvariantParam {
	return []testInvariantParam{
		{ // when commission is zero
			nil,
			[]sdk.DecCoins{NewTestDecCoins(0, 0)},
			true,
		},
		{ // when withdraw commission failed
			NewTestDecCoins(5, 1),
			[]sdk.DecCoins{NewTestDecCoins(15, 1)},
			false,
		},
		{ // when the sum of commission is not equal to distribution account
			NewTestDecCoins(30, 1),
			[]sdk.DecCoins{NewTestDecCoins(15, 1), NewTestDecCoins(20, 1)},
			false,
		},
		{ // when the sum of commission is not equal to distribution account
			NewTestDecCoins(45, 1),
			[]sdk.DecCoins{NewTestDecCoins(15, 1), NewTestDecCoins(20, 1)},
			false,
		},
		{ // success
			NewTestDecCoins(45, 1),
			[]sdk.DecCoins{NewTestDecCoins(15, 1), NewTestDecCoins(30, 1)},
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
		ak.SetAccount(ctx, acc)
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
