package keeper_test

import (
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evidence/internal/types"
	"github.com/okex/exchain/x/staking"
	stakingtypes "github.com/okex/exchain/x/staking/types"

	"github.com/okex/exchain/libs/tendermint/crypto"
)

const EPOCH = 252

func newTestMsgCreateValidator(address sdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) stakingtypes.MsgCreateValidator {
	//	commission := staking.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	msd := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(10000))
	return staking.NewMsgCreateValidator(address, pubKey,
		staking.NewDescription("my moniker", "my identity", "my website", "my details"), msd,
	)
}

func (suite *KeeperTestSuite) TestHandleDoubleSign() {
	ctx := suite.ctx.WithIsCheckTx(false).WithBlockHeight(EPOCH)
	suite.populateValidators(ctx)

	power := sdk.NewIntFromUint64(10000)
	stakingParams := suite.app.StakingKeeper.GetParams(ctx)
	amt := power
	operatorAddr, val := valAddresses[0], pubkeys[0]

	// create validator
	res, err := staking.NewHandler(suite.app.StakingKeeper)(ctx, newTestMsgCreateValidator(operatorAddr, val, amt))
	suite.NoError(err)
	suite.NotNil(res)

	// execute end-blocker and verify validator attributes
	staking.EndBlocker(ctx, suite.app.StakingKeeper)
	suite.Equal(
		suite.app.BankKeeper.GetCoins(ctx, sdk.AccAddress(operatorAddr)),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt.Sub(amt))),
	)

	// handle a signature to set signing info
	suite.app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), amt.Int64(), true)

	// double sign less than max age
	evidence := types.Equivocation{
		Height:           0,
		Time:             time.Unix(0, 0),
		Power:            power.Int64(),
		ConsensusAddress: sdk.ConsAddress(val.Address()),
	}
	suite.keeper.HandleDoubleSign(ctx, evidence)

	// should be jailed and tombstoned
	suite.True(suite.app.StakingKeeper.Validator(ctx, operatorAddr).IsJailed())
	suite.True(suite.app.SlashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(val.Address())))

	// submit duplicate evidence
	suite.keeper.HandleDoubleSign(ctx, evidence)

	// jump to past the unbonding period
	ctx = ctx.WithBlockTime(time.Unix(1, 0).Add(stakingParams.UnbondingTime))

	// require we cannot unjail
	suite.Error(suite.app.SlashingKeeper.Unjail(ctx, operatorAddr))

	// require we be able to unbond now
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	msgDestroy := stakingtypes.NewMsgDestroyValidator(sdk.AccAddress(operatorAddr))
	res, err = staking.NewHandler(suite.app.StakingKeeper)(ctx, msgDestroy)
	suite.NoError(err)
	suite.NotNil(res)
}

func (suite *KeeperTestSuite) TestHandleDoubleSign_TooOld() {
	ctx := suite.ctx.WithIsCheckTx(false).WithBlockHeight(EPOCH).WithBlockTime(time.Now())
	suite.populateValidators(ctx)

	power := sdk.NewIntFromUint64(10000)
	//stakingParams := suite.app.StakingKeeper.GetParams(ctx)
	amt := power
	operatorAddr, val := valAddresses[0], pubkeys[0]

	// create validator
	res, err := staking.NewHandler(suite.app.StakingKeeper)(ctx, newTestMsgCreateValidator(operatorAddr, val, amt))
	suite.NoError(err)
	suite.NotNil(res)

	// execute end-blocker and verify validator attributes
	staking.EndBlocker(ctx, suite.app.StakingKeeper)
	suite.Equal(
		suite.app.BankKeeper.GetCoins(ctx, sdk.AccAddress(operatorAddr)),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt.Sub(amt))),
	)

	evidence := types.Equivocation{
		Height:           0,
		Time:             ctx.BlockTime(),
		Power:            power.Int64(),
		ConsensusAddress: sdk.ConsAddress(val.Address()),
	}
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(suite.app.EvidenceKeeper.MaxEvidenceAge(ctx) + 1))
	suite.keeper.HandleDoubleSign(ctx, evidence)

	suite.False(suite.app.StakingKeeper.Validator(ctx, operatorAddr).IsJailed())
	suite.False(suite.app.SlashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(val.Address())))
}
