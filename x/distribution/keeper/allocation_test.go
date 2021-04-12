package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	//init
	ctx, _, k, sk, _ := CreateTestInputDefault(t, false, 1000)
	val := sk.Validator(ctx, valOpAddr1)

	// allocate tokens
	tokens := NewTestSysCoins(1, sdk.Precision)
	k.AllocateTokensToValidator(ctx, val, tokens)

	// check commissions
	expected := NewTestSysCoins(1, sdk.Precision)
	require.Equal(t, expected, k.GetValidatorAccumulatedCommission(ctx, val.GetOperator()))
}

type testAllocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        sdk.SysCoins
}

func getTestAllocationParams() []testAllocationParam {
	return []testAllocationParam{
		{ //test the case when fee is zero
			10,
			[]bool{true, true, true, true}, []bool{false, false, false, false},
			nil,
		},
		{ //test the case where total power is zero
			0,
			[]bool{true, true, true, true}, []bool{false, false, false, false},
			NewTestSysCoins(123, 2),
		},
		{ //test the case where just the part of vals has voted
			10,
			[]bool{true, true, false, false}, []bool{false, false, false, false},
			NewTestSysCoins(123, 2),
		},
		{ //test the case where two vals is jailed
			10,
			[]bool{true, true, false, false}, []bool{false, true, true, false},
			NewTestSysCoins(123, 2),
		},
	}
}

func TestAllocateTokensToManyValidators(t *testing.T) {
	tests := getTestAllocationParams()
	_, _, valConsAddrs := GetTestAddrs()

	//start test
	for _, test := range tests {
		ctx, ak, k, sk, _ := CreateTestInputDefault(t, false, 1000)
		//set the fee
		setTestFees(t, ctx, k, ak, test.fee)
		// crate votes info
		votes := createTestVotes(ctx, sk, test)

		// allocate the tokens
		k.AllocateTokens(ctx, test.totalPower, valConsAddrs[0], votes)
		commissions := NewTestSysCoins(0, 0)
		k.IterateValidatorAccumulatedCommissions(ctx,
			func(val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
				commissions = commissions.Add(commission...)
				return false
			})
		totalCommissions := k.GetDistributionAccount(ctx).GetCoins()
		communityCoins := k.GetFeePoolCommunityCoins(ctx)
		require.Equal(t, totalCommissions, communityCoins.Add(commissions...))
		require.Equal(t, test.fee, totalCommissions)
	}
}

func setTestFees(t *testing.T, ctx sdk.Context, k Keeper, ak auth.AccountKeeper, fees sdk.SysCoins) {
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	require.NotNil(t, feeCollector)
	err := feeCollector.SetCoins(fees)
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
}

func createTestVotes(ctx sdk.Context, sk staking.Keeper, test testAllocationParam) []abci.VoteInfo {
	_, valConsPks, valConsAddrs := GetTestAddrs()
	var votes []abci.VoteInfo
	for i := int64(0); i < int64(len(test.isVote)); i++ {
		if test.isJailed[i] {
			sk.Jail(ctx, valConsAddrs[i])
		}
		abciVal := abci.Validator{Address: valConsPks[i].Address(), Power: i + 1}
		if test.isVote[i] {
			votes = append(votes, abci.VoteInfo{Validator: abciVal, SignedLastBlock: true})
		}
	}
	return votes
}
