package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/gov/types"
	"github.com/okex/okexchain/x/params"
	paramsTypes "github.com/okex/okexchain/x/params/types"
)

func TestKeeper_AddDeposit(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// deposit on proposal which is not exist
	err := keeper.AddDeposit(ctx, 0, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.NotNil(t, err)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	// nil address deposit
	err = keeper.AddDeposit(ctx, proposalID, sdk.AccAddress{},
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.NotNil(t, err)

	// deposit on proposal whose status is not DepositPeriod
	proposal.Status = types.StatusPassed
	keeper.SetProposal(ctx, proposal)
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.NotNil(t, err)

	proposal.Status = types.StatusDepositPeriod
	keeper.SetProposal(ctx, proposal)
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.Nil(t, err)

	// change old deposit and activate proposal
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)}, "")
	require.Nil(t, err)

	// deposit on proposal which registered proposal handler router
	paramsChanges := []params.ParamChange{{Subspace: "staking", Key: "MaxValidators", Value: "105"}}
	content = paramsTypes.NewParameterChangeProposal("Test", "", paramsChanges, 1)
	proposal, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID = proposal.ProposalID
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)}, "")
	require.Nil(t, err)
}

func TestKeeper_GetDeposit(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.Nil(t, err)

	expectedDeposit := types.Deposit{
		ProposalID: proposalID,
		Depositor:  Addrs[0],
		Amount:     sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)},
	}
	deposit, found := keeper.GetDeposit(ctx, proposalID, Addrs[0])
	require.True(t, found)
	require.True(t, deposit.Equals(expectedDeposit))

	// get deposit from db
	deposit, found = keeper.GetDeposit(ctx, proposalID, Addrs[0])
	require.True(t, found)
	require.True(t, deposit.Equals(expectedDeposit))
}

func TestKeeper_GetDeposits(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.Nil(t, err)

	err = keeper.AddDeposit(ctx, proposalID, Addrs[1],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}, "")
	require.Nil(t, err)

	expectedDeposits := types.Deposits{
		{
			ProposalID: proposalID,
			Depositor:  Addrs[0],
			Amount:     sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)},
		},
		{
			ProposalID: proposalID,
			Depositor:  Addrs[1],
			Amount:     sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)},
		},
	}
	deposits := keeper.GetDeposits(ctx, proposalID)
	require.Equal(t, expectedDeposits, deposits)

	// get deposits from db
	deposits = keeper.GetDeposits(ctx, proposalID)
	require.Equal(t, expectedDeposits, deposits)
}

func TestKeeper_DistributeDeposits(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	amount1 := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0], amount1, "")
	require.Nil(t, err)

	amount2 := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}
	err = keeper.AddDeposit(ctx, proposalID, Addrs[1], amount2, "")
	require.Nil(t, err)

	moduleAccBalance := keeper.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins()
	require.Equal(t, amount1.Add(amount2), moduleAccBalance)

	// after DistributeDeposits
	keeper.DistributeDeposits(ctx, proposalID)
	moduleAccBalance = keeper.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins()
	require.Equal(t, sdk.Coins(nil), moduleAccBalance)
	feeCollectorBalance := keeper.SupplyKeeper().GetModuleAccount(ctx, keeper.feeCollectorName).GetCoins()
	require.Equal(t, amount1.Add(amount2), feeCollectorBalance)
}

func TestKeeper_RefundDeposits(t *testing.T) {
	ctx, accKeeper, keeper, _, _ := CreateTestInput(t, false, 1000)
	amount := accKeeper.GetAccount(ctx, Addrs[0]).GetCoins().AmountOf(sdk.DefaultBondDenom)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	amount1 := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0], amount1, "")
	require.Nil(t, err)

	amount2 := sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}
	err = keeper.AddDeposit(ctx, proposalID, Addrs[1], amount2, "")
	require.Nil(t, err)

	moduleAccBalance := keeper.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins()
	require.Equal(t, amount1.Add(amount2), moduleAccBalance)

	// after RefundDeposits
	keeper.RefundDeposits(ctx, proposalID)
	moduleAccBalance = keeper.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins()
	require.Equal(t, sdk.Coins(nil), moduleAccBalance)

	require.Equal(t, amount, sdk.NewDec(1000))

	require.Equal(t, accKeeper.GetAccount(ctx, Addrs[1]).GetCoins(),
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1000)})

	// refund panic
	content = types.NewTextProposal("Test", "description")
	proposal, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID = proposal.ProposalID

	amount1 = sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 40)}
	err = keeper.AddDeposit(ctx, proposalID, Addrs[0], amount1, "")
	require.Nil(t, err)

	err = keeper.SupplyKeeper().SendCoinsFromModuleToModule(ctx, types.ModuleName, keeper.feeCollectorName,
		amount1)
	require.Nil(t, err)
	require.Panics(t, func() {
		keeper.RefundDeposits(ctx, proposalID)
	})
}
