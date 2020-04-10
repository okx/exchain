package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	//init
	ctx, _, k, sk, _ := CreateTestInputDefault(t, false, 1000)
	val := sk.Validator(ctx, valOpAddr1)

	// allocate tokens
	tokens := NewTestDecCoins(1, 8)
	k.AllocateTokensToValidator(ctx, val, tokens)

	// check commissions
	expected := NewTestDecCoins(1, 8)
	require.Equal(t, expected, k.GetValidatorAccumulatedCommission(ctx, val.GetOperator()))
}

type testAllocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        int64
	expected   [4]int64
}

func getTestAllocationParams() []testAllocationParam {
	return []testAllocationParam{
		{ //test the case when fee is zero
			10,
			[]bool{true, true, true, true}, []bool{false, false, false, false},
			0, [4]int64{0, 0, 0, 0},
		},
		{ //test the case where total power is zero
			0,
			[]bool{true, true, true, true}, []bool{false, false, false, false},
			120, [4]int64{0, 0, 0, 0},
		},
		{ //test the case where just the part of vals has voted
			10,
			[]bool{true, true, false, false}, []bool{false, false, false, false},
			120, [4]int64{24, 33, 27, 36},
		},
		{ //test the case where two vals is jailed
			10,
			[]bool{true, true, false, false}, []bool{false, true, true, false},
			120, [4]int64{48, 0, 0, 72},
		},
	}
}

func TestAllocateTokensToManyValidators(t *testing.T) {
	tests := getTestAllocationParams()
	valOpAddrs, _, valConsAddrs := GetTestAddrs()

	//start test
	for _, test := range tests {
		ctx, ak, k, sk, supplyKeeper := CreateTestInputDefault(t, false, 1000)
		//set the fee
		feeCoins, _ := NewTestDecCoins(test.fee, 0).TruncateDecimal()
		setTestFees(t, ctx, k, ak, feeCoins)
		// crate votes info
		votes := createTestVotes(ctx, sk, test)

		// allocate the tokens
		k.AllocateTokens(ctx, test.totalPower, valConsAddrs[0], votes)

		remain := feeCoins
		// check the results
		for i := int64(0); i < int64(len(test.isVote)); i++ {
			expectedCommssion := NewTestDecCoins(test.expected[i], 0)
			if !expectedCommssion.IsZero() {
				require.Equal(t, expectedCommssion, k.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]))
				commission, err := k.WithdrawValidatorCommission(ctx, valOpAddrs[i])
				require.NoError(t, err)
				expectedWithdrawCommission, _ := expectedCommssion.TruncateDecimal()
				require.Equal(t, expectedWithdrawCommission, commission)
				remain = remain.Sub(commission)
			}
		}
		//TODO when rollback the community pool
		//require.True(t, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins().IsEqual(remain))
		//require.True(t, supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName).GetCoins().IsZero())
		require.Equal(t, true, supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName).GetCoins().IsEqual(remain))
	}
}

func setTestFees(t *testing.T, ctx sdk.Context, k Keeper, ak auth.AccountKeeper, fees sdk.DecCoins) {
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
