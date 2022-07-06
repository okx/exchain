package distribution

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWithdrawDelegatorReward(t *testing.T) {
	ctx, _, dk, sk, _ := keeper.CreateTestInputDefault(t, false, 10)
	handler := NewHandler(dk)
	delAddr1 := keeper.TestDelAddrs[0]
	valAddr1 := keeper.TestValAddrs[0]
	valOpAddrs := []sdk.ValAddress{valAddr1}

	delAddr2 := keeper.TestDelAddrs[1]
	valAddr2 := keeper.TestValAddrs[1]
	valOpAddrs2 := []sdk.ValAddress{valAddr1, valAddr2}

	msg := NewMsgWithdrawDelegatorReward(delAddr1, valAddr1)
	msg2 := NewMsgWithdrawDelegatorReward(delAddr2, valAddr2)

	// not support msg type
	_, err := handler(ctx, msg)
	require.Equal(t, types.ErrUnknownDistributionMsgType(), err)
	ctx.SetBlockTime(time.Now())
	// deposit and add shares
	keeper.DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)

	keeper.DoDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr2, valOpAddrs2)

	// delegation not exist
	_, err = handler(ctx, msg)
	require.Equal(t, types.ErrUnknownDistributionMsgType(), err)

	// deposit and add shares ErrCodeNotSupportDistributionMethod
	keeper.DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	// delegation not exist
	_, err = handler(ctx, msg)
	require.Equal(t, types.ErrUnknownDistributionMsgType(), err)

	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
	keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOnChain)

	// withdraw del1 ok
	keeper.DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr1, valOpAddrs2)
	_, err = handler(ctx, msg)
	require.Nil(t, err)

	// withdraw del2 ok
	_, err = handler(ctx, msg2)
	require.Nil(t, err)
}

type allocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        sdk.SysCoins
}

func createVotes(ctx sdk.Context, sk staking.Keeper, test allocationParam) []abci.VoteInfo {
	var votes []abci.VoteInfo
	for i := int64(0); i < int64(len(test.isVote)); i++ {
		if test.isJailed[i] {
			sk.Jail(ctx, keeper.TestConsAddrs[i])
		}
		abciVal := abci.Validator{Address: keeper.TestConsAddrs[i], Power: i + 1}
		if test.isVote[i] {
			votes = append(votes, abci.VoteInfo{Validator: abciVal, SignedLastBlock: true})
		}
	}
	return votes
}

func TestWithdrawValidatorCommission(t *testing.T) {
	ctx, ak, dk, sk, supplyKeeper := keeper.CreateTestInputDefault(t, false, 10)
	_ = sk
	handler := NewHandler(dk)
	//delAddr1 := keeper.TestDelAddrs[0]
	valAddr1 := keeper.TestValAddrs[0]
	//valOpAddrs := []sdk.ValAddress{valAddr1}

	msg := NewMsgWithdrawValidatorCommission(valAddr1)

	// ErrNoValidatorCommission
	_, err := handler(ctx, msg)
	require.Equal(t, types.ErrNoValidatorCommission(), err)

	//ok
	//newBlockAndAllocateReward(t)
	staking.EndBlocker(ctx, sk)
	feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	require.NotNil(t, feeCollector)
	err = feeCollector.SetCoins(sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}})
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
	allocationParam := allocationParam{
		10,
		[]bool{true, true, true, true}, []bool{false, false, false, false},
		nil,
	}
	votes := createVotes(ctx, sk, allocationParam)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	require.Nil(t, err)

	_, err = handler(ctx, msg)
	require.Nil(t, err)
}
