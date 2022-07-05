package distribution

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
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

	// delegation not exist
	_, err := handler(ctx, msg)
	require.Equal(t, types.ErrCodeEmptyDelegationDistInfo(), err)
	ctx.SetBlockTime(time.Now())
	// deposit and add shares
	keeper.DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)

	keeper.DoDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr2, valOpAddrs2)

	// delegation not exist
	_, err = handler(ctx, msg)
	require.Equal(t, types.ErrCodeNotSupportDistributionMethod(), err)

	// deposit and add shares ErrCodeNotSupportDistributionMethod
	keeper.DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	// delegation not exist
	_, err = handler(ctx, msg)
	require.Equal(t, types.ErrCodeNotSupportDistributionMethod(), err)

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
	err = feeCollector.SetCoins(blockRewardValueTokens)
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
	testAllocationParam := testAllocationParam{
		10,
		[]bool{true, true, true, true}, []bool{false, false, false, false},
		nil,
	}
	votes := createTestVotes(ctx, sk, testAllocationParam)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	require.Nil(t, err)

	_, err = handler(ctx, msg)
	require.Nil(t, err)
}
