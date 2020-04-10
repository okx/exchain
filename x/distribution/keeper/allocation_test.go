package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	//init
	ctx, _, k, sk, _ := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create one validator
	msg := staking.NewMsgCreateValidator(valOpAddr1, valConsPk1, staking.Description{}, NewDecCoin(1))
	require.True(t, sh(ctx, msg).IsOK())
	val := sk.Validator(ctx, valOpAddr1)

	// allocate tokens
	tokens := NewDecCoins(1, 8)
	k.AllocateTokensToValidator(ctx, val, tokens)

	// check commission
	expected := NewDecCoins(1, 8)
	require.Equal(t, expected, k.GetValidatorAccumulatedCommission(ctx, val.GetOperator()))
}

type testParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        int64
	expected   [4]int64
}

func getTestParams() []testParam {
	return []testParam{
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
	tests := getTestParams()

	//init
	valOpAddrs, valConsPks, valConsAddrs := GetAddrs()
	ctx, ak, k, sk, supplyKeeper := CreateTestInputDefault(t, false, 1000)
	sh := staking.NewHandler(sk)

	// create four validators
	for i := int64(0); i < 4; i++ {
		msg := staking.NewMsgCreateValidator(valOpAddrs[i], valConsPks[i], staking.Description{}, NewDecCoin(i+1))
		require.True(t, sh(ctx, msg).IsOK())
		// assert initial state: zero current rewards
		require.True(t, k.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]).IsZero())
	}

	//start test
	for _, test := range tests {
		//set the fee
		feeCollector := supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
		require.NotNil(t, feeCollector)
		fees, _ := NewDecCoin(test.fee).TruncateDecimal()
		feeCoins := sdk.NewCoins(fees)
		err := feeCollector.SetCoins(feeCoins)
		require.NoError(t, err)
		ak.SetAccount(ctx, feeCollector)

		// crate votes info
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

		// allocate the tokens
		k.AllocateTokens(ctx, test.totalPower, valConsAddrs[0], votes)

		remainCommission := feeCoins
		// check the results
		for i := int64(0); i < int64(len(test.isVote)); i++ {
			expectedCommssion := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(test.expected[i]))
			if !expectedCommssion.IsZero() {
				require.Equal(t, expectedCommssion, k.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]))
				commission, err := k.WithdrawValidatorCommission(ctx, valOpAddrs[i])
				require.NoError(t, err)
				expectedWithdrawCommission, _ := expectedCommssion.TruncateDecimal()
				require.Equal(t, expectedWithdrawCommission, commission)
				remainCommission = remainCommission.Sub(commission)
			}
			if test.isJailed[i] {
				sk.Unjail(ctx, valConsAddrs[i])
			}
		}
		feeCollector = supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
		require.Equal(t, true, feeCollector.GetCoins().IsEqual(remainCommission))
	}
}
