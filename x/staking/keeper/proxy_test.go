package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestVoteValidatorsAndWithdrawVote(t *testing.T) {
	ctx, _, mkeeper := CreateTestInput(t, false, 0)
	keeper := mkeeper.Keeper
	valsOld := createVals(ctx, 4, keeper)

	// never vote before
	dlgAddr := addrDels[0]
	lastVals, lastVotes := keeper.GetLastValsVotedExisted(ctx, dlgAddr)
	require.Nil(t, lastVals)
	require.True(t, lastVotes.IsZero())

	// withdraw the votes last time
	keeper.WithdrawLastVotes(ctx, dlgAddr, lastVals, lastVotes)

	// votes validators
	voteOrig := sdk.NewDec(10000)
	_, e := keeper.VoteValidators(ctx, dlgAddr, valsOld, voteOrig)
	require.Nil(t, e)

	// check valsOld status
	valsNew := getVals(ctx, valsOld, keeper, t)
	for i := 0; i < 4; i++ {
		require.True(t, valsNew[i].DelegatorShares.GT(valsOld[i].DelegatorShares),
			valsNew[i].Standardize().String(), valsOld[i].Standardize().String())

		// check votes
		vote, found := keeper.GetVote(ctx, dlgAddr, valsNew[i].OperatorAddress)
		require.True(t, found)
		require.True(t, vote.GT(lastVotes), vote)
	}

	// standarlize
	sVals := valsNew.Standardize()
	require.NotNil(t, sVals)
	r, err := sVals.MarshalYAML()
	require.Nil(t, err)
	require.Contains(t, r, "Operator Address")
}

func createVals(ctx sdk.Context, num int, keeper Keeper) types.Validators {
	vals := make(types.Validators, num)
	for i := 0; i < num; i++ {
		vals[i] = types.NewValidator(addrVals[i], PKs[i], types.Description{})
		keeper.SetValidator(ctx, vals[i])
	}

	return vals
}

func getVals(ctx sdk.Context, valOld types.Validators, keeper Keeper, t *testing.T) types.Validators {
	lenVals := len(valOld)
	gotVals := make(types.Validators, lenVals)
	for i := 0; i < lenVals; i++ {
		val, found := keeper.GetValidator(ctx, valOld[i].OperatorAddress)
		require.True(t, found)
		gotVals[i] = val
	}
	return gotVals
}
