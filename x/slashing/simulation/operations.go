package simulation

import (
	"math/rand"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	simappparams "github.com/okex/exchain/libs/cosmos-sdk/simapp/params"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/x/slashing/internal/keeper"
	"github.com/okex/exchain/x/slashing/internal/types"
	stakingkeeper "github.com/okex/exchain/x/staking/keeper"
)

// Simulation operation weights constants
const (
	OpWeightMsgUnjail = "op_weight_msg_unjail"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak types.AccountKeeper,
	k keeper.Keeper, sk stakingkeeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgUnjail int
	appParams.GetOrGenerate(cdc, OpWeightMsgUnjail, &weightMsgUnjail, nil,
		func(_ *rand.Rand) {
			weightMsgUnjail = simappparams.DefaultWeightMsgUnjail
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUnjail,
			SimulateMsgUnjail(ak, k, sk),
		),
	}
}

// SimulateMsgUnjail generates a MsgUnjail with random values
// nolint: funlen
func SimulateMsgUnjail(ak types.AccountKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simulation.Operation { // nolint:interfacer
	return nil
	//todo:disable follow code just for now
	//return func(
	//	r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
	//	accs []simulation.Account, chainID string,
	//) (simulation.OperationMsg, []simulation.FutureOperation, error) {
	//
	//	validator, ok := stakingkeeper.RandomValidator(r, sk, ctx)
	//	if !ok {
	//		return simulation.NoOpMsg(types.ModuleName), nil, nil // skip
	//	}
	//
	//	simAccount, found := simulation.FindAccount(accs, sdk.AccAddress(validator.GetOperator()))
	//	if !found {
	//		return simulation.NoOpMsg(types.ModuleName), nil, nil // skip
	//	}
	//
	//	if !validator.IsJailed() {
	//		// TODO: due to this condition this message is almost, if not always, skipped !
	//		return simulation.NoOpMsg(types.ModuleName), nil, nil
	//	}
	//
	//	consAddr := sdk.ConsAddress(validator.GetConsPubKey().Address())
	//	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	//	if !found {
	//		return simulation.NoOpMsg(types.ModuleName), nil, nil // skip
	//	}
	//
	//	selfDel := sk.Delegation(ctx, simAccount.Address, validator.GetOperator())
	//	if selfDel == nil {
	//		return simulation.NoOpMsg(types.ModuleName), nil, nil // skip
	//	}
	//
	//	account := ak.GetAccount(ctx, sdk.AccAddress(validator.GetOperator()))
	//	fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
	//	if err != nil {
	//		return simulation.NoOpMsg(types.ModuleName), nil, err
	//	}
	//
	//	msg := types.NewMsgUnjail(validator.GetOperator())
	//
	//	tx := helpers.GenTx(
	//		[]sdk.Msg{msg},
	//		fees,
	//		helpers.DefaultGenTxGas,
	//		chainID,
	//		[]uint64{account.GetAccountNumber()},
	//		[]uint64{account.GetSequence()},
	//		simAccount.PrivKey,
	//	)
	//
	//	_, res, err := app.Deliver(tx)
	//
	//	// result should fail if:
	//	// - validator cannot be unjailed due to tombstone
	//	// - validator is still in jailed period
	//	// - self delegation too low
	//	if info.Tombstoned ||
	//		ctx.BlockHeader().Time.Before(info.JailedUntil) ||
	//		validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
	//		if res != nil && err == nil {
	//			if info.Tombstoned {
	//				return simulation.NewOperationMsg(msg, true, ""), nil, errors.New("validator should not have been unjailed if validator tombstoned")
	//			}
	//			if ctx.BlockHeader().Time.Before(info.JailedUntil) {
	//				return simulation.NewOperationMsg(msg, true, ""), nil, errors.New("validator unjailed while validator still in jail period")
	//			}
	//			if validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
	//				return simulation.NewOperationMsg(msg, true, ""), nil, errors.New("validator unjailed even though self-delegation too low")
	//			}
	//		}
	//		// msg failed as expected
	//		return simulation.NewOperationMsg(msg, false, ""), nil, nil
	//	}
	//
	//	if err != nil {
	//		return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
	//	}
	//
	//	return simulation.NewOperationMsg(msg, true, ""), nil, nil
	//}
}
