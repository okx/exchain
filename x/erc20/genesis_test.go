package erc20_test

import (
	"github.com/okex/exchain/x/erc20"
	"github.com/okex/exchain/x/erc20/types"
)

func (suite *Erc20TestSuite) TestInitGenesis() {
	testCases := []struct {
		name     string
		malleate func()
		genState types.GenesisState
		expPanic bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesisState(),
			false,
		},
		{
			"Wrong denom in external token mapping",
			func() {},
			types.GenesisState{
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong denom in auto token mapping",
			func() {},
			types.GenesisState{
				AutoContracts: []types.TokenMapping{
					{
						Denom:    "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong contract in external token mapping",
			func() {},
			types.GenesisState{
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x00000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong contract in auto token mapping",
			func() {},
			types.GenesisState{
				AutoContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x00000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Correct token mapping",
			func() {},
			types.GenesisState{
				Params: types.DefaultParams(),
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()
			if tc.expPanic {
				suite.Require().Panics(
					func() {
						erc20.InitGenesis(suite.ctx, suite.app.Erc20Keeper, tc.genState)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						erc20.InitGenesis(suite.ctx, suite.app.Erc20Keeper, tc.genState)
					},
				)
			}
		})
	}
}

func (suite *Erc20TestSuite) TestExportGenesis() {
	genesisState := erc20.ExportGenesis(suite.ctx, suite.app.Erc20Keeper)
	suite.Require().Equal(genesisState.Params.IbcTimeout, types.DefaultParams().IbcTimeout)
	suite.Require().Equal(genesisState.Params.EnableAutoDeployment, types.DefaultParams().EnableAutoDeployment)
}
