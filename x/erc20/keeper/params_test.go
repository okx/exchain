package keeper_test

import (
	"errors"

	erc20Keeper "github.com/okex/exchain/x/erc20/keeper"
	"github.com/okex/exchain/x/erc20/types"
)

func (suite *KeeperTestSuite) TestGetSourceChannelID() {

	testCases := []struct {
		name          string
		ibcDenom      string
		expectedError error
		postCheck     func(channelID string)
	}{
		{
			"wrong ibc denom",
			"test",
			errors.New("ibc denom is invalid: test is invalid"),
			func(channelID string) {},
		},
		{
			"correct ibc denom",
			types.IbcDenomDefaultValue,
			nil,
			func(channelID string) {
				suite.Require().Equal(channelID, "channel-0")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			// Create erc20 Keeper with mock transfer keeper
			erc20Keeper := erc20Keeper.NewKeeper(
				suite.app.Cdc,
				suite.app.GetKey(types.StoreKey),
				suite.app.GetSubspace(types.ModuleName),
				suite.app.AccountKeeper,
				suite.app.SupplyKeeper,
				suite.app.BankKeeper,
				suite.app.EvmKeeper,
				IbcKeeperMock{},
			)
			suite.app.Erc20Keeper = erc20Keeper

			channelId, err := suite.app.Erc20Keeper.GetSourceChannelID(suite.ctx, tc.ibcDenom)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
				tc.postCheck(channelId)
			}
		})
	}
}
