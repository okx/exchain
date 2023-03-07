package staking

import (
	"testing"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/global"
	keep "github.com/okx/okbchain/x/staking/keeper"
	"github.com/okx/okbchain/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) TestEditValidatorCommission() {
	testCases := []struct {
		title              string
		setMilestoneHeight func()
		newRate            string
		setBlockTime       func(ctx *sdk.Context)
		handlerErrType     int
		err                [5]error
	}{
		{
			"default ok",
			func() {},
			"0.5",
			func(ctx *sdk.Context) {
				ctx.SetBlockTime(time.Now())
				ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			},
			0,
			[5]error{nil, nil, nil, nil},
		},
		{
			"-0.5",
			func() {},
			"-0.5",
			func(ctx *sdk.Context) {
				ctx.SetBlockTime(time.Now())
				ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			},
			0,
			[5]error{types.ErrInvalidCommissionRate(), types.ErrInvalidCommissionRate(),
				types.ErrCommissionNegative(), types.ErrCommissionNegative()},
		},
		{
			"do not set block time",
			func() {},
			"0.5",
			func(ctx *sdk.Context) {},
			0,
			[5]error{nil, nil, nil, types.ErrCommissionUpdateTime()},
		},
	}

	for _, tc := range testCases {
		global.SetGlobalHeight(0)
		suite.Run(tc.title, func() {
			ctx, _, mKeeper := CreateTestInput(suite.T(), false, SufficientInitPower)
			tc.setMilestoneHeight()
			keeper := mKeeper.Keeper
			_ = setInstantUnbondPeriod(keeper, ctx)
			handler := NewHandler(keeper)

			newRate, _ := sdk.NewDecFromStr(tc.newRate)
			msgEditValidator := NewMsgEditValidatorCommissionRate(sdk.ValAddress(keep.Addrs[0]), newRate)
			err := msgEditValidator.ValidateBasic()
			require.Equal(suite.T(), tc.err[0], err)

			// validator not exist
			got, err := handler(ctx, msgEditValidator)
			if tc.handlerErrType == 0 {
				require.Equal(suite.T(), ErrNoValidatorFound(msgEditValidator.ValidatorAddress.String()), err)
			} else {
				require.NotNil(suite.T(), err)
			}

			//create validator
			validatorAddr := sdk.ValAddress(keep.Addrs[0])
			msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], DefaultMSD)
			got, err = handler(ctx, msgCreateValidator)
			require.Nil(suite.T(), err, "expected create-validator to be ok, got %v", got)

			// must end-block
			updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
			require.Equal(suite.T(), 1, len(updates))
			SimpleCheckValidator(suite.T(), ctx, keeper, validatorAddr, DefaultMSD, sdk.Bonded,
				SharesFromDefaultMSD, false)

			// normal rate
			newRate, _ = sdk.NewDecFromStr(tc.newRate)
			msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
			err = msgEditValidator.ValidateBasic()
			require.Equal(suite.T(), tc.err[1], err)
			got, err = handler(ctx, msgEditValidator)
			if tc.handlerErrType == 0 {
				require.Equal(suite.T(), tc.err[2], err)
			} else {
				require.NotNil(suite.T(), err)
			}

			tc.setBlockTime(&ctx)
			msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
			got, err = handler(ctx, msgEditValidator)
			if tc.handlerErrType == 0 {
				require.Equal(suite.T(), tc.err[3], err)
			} else {
				require.NotNil(suite.T(), err)
			}
		})
	}
}
