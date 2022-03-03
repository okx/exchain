package mint_test

import (
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	treasure1 = types.NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 2))
	treasure2 = types.NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 2))
	treasure3 = types.NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 2))
	treasure4 = types.NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 2))
	treasure5 = types.NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 2))
	treasures = []types.Treasure{*treasure1, *treasure2, *treasure3, *treasure4, *treasure5}
)

type MintTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	govHandler govtypes.Handler
	querier    sdk.Querier
	app        *app.OKExChainApp
	codec      *codec.Codec
}

func (suite *MintTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})

	suite.govHandler = mint.NewManageTreasuresProposalHandler(&suite.app.MintKeeper)
	suite.querier = mint.NewQuerier(suite.app.MintKeeper)
	suite.codec = codec.New()
}

func TestMintTestSuite(t *testing.T) {
	suite.Run(t, new(MintTestSuite))
}

func (suite *MintTestSuite) TestTreasuresProposal() {
	proposal := types.NewManageTreasuresProposal(
		"default title",
		"default description",
		treasures,
		true,
	)
	govProposal := govtypes.Proposal{
		Content: proposal,
	}
	passfunc := func(err error, trs []types.Treasure, msg string) {
		suite.Require().NoError(err, msg)
	}
	treasuresError := func(err error, trs []types.Treasure, msg string) {
		suite.Require().Error(err)
		suite.Require().Contains(err.Error(), "treasure proportion should non-negative and less than one", msg)
	}
	sumProportionError := func(err error, trs []types.Treasure, msg string) {
		suite.Require().Error(err)
		suite.Require().Contains(err.Error(), "the sum of treasure proportion should non-negative and less than one", msg)
	}
	unexistError := func(err error, trs []types.Treasure, msg string) {
		suite.Require().Error(err)
		suite.Require().Contains(err.Error(), "because it's not exist from treasures", msg)
	}
	testCases := []struct {
		msg             string
		expectfunc      func(err error, trs []types.Treasure, msg string)
		prepare         func()
		targetTreasures []types.Treasure
	}{
		{
			"add one into empty",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{treasures[0]}
			},
			[]types.Treasure{treasures[0]},
		},
		{
			"add one into one",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})

			},
			[]types.Treasure{treasures[1], treasures[0]},
		},
		{
			"add one into multi",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[1], treasures[0]},
		},
		{
			"add multi into multi",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{treasures[1], treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[3], treasures[2], treasures[1], treasures[0]},
		},
		{
			"add multi into one",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{treasures[1], treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})

			},
			[]types.Treasure{treasures[3], treasures[1], treasures[0]},
		},
		{
			"update one into one",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			[]types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
		},
		{
			"update one into multi",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
		},
		{
			"update multi into multi",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{{Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 2)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
		},
		{
			"update/insert multi into multi",
			passfunc,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 2)}, treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[3], {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 2)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
		},
		{
			"delete one from one",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3]})
			},
			[]types.Treasure{},
		},
		{
			"delete one from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{treasures[0]},
		},
		{
			"delete multi from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[3], treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
		{
			"delete multi from multi more",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[3], treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0], treasures[1]})
			},
			[]types.Treasure{treasures[1]},
		},
		{
			"add multi(negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(more negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(all negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(8, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(input sum proportion more than one) into multi",
			sumProportionError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(1, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"add multi(result sum proportion more than one) into multi",
			sumProportionError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})

			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(more negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(-16, 2)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(all negative) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(-16, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(more more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 0)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(all more than one) into multi",
			treasuresError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(16, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(input the sum proportion all more than one) into multi",
			sumProportionError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(21, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[0]},
		},
		{
			"update multi(result the sum proportion all more than one) into multi",
			sumProportionError,
			func() {
				proposal.IsAdded = true
				proposal.Treasures = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(20, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[1], treasures[2]})
			},
			[]types.Treasure{treasures[2], treasures[1], treasures[0]},
		},
		{
			"delete multi(unexist) from multi",
			unexistError,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[4], treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{treasures[3], treasures[0]},
		},
		{
			"delete multi(part unexist) from multi",
			unexistError,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[4], treasures[2], treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{treasures[3], treasures[0]},
		},
		{
			"delete multi(all unexist) from multi",
			unexistError,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{treasures[3], treasures[0]},
		},
		{
			"delete multi(negative) from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{{Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(-2, 0)}, treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
		{
			"delete multi(more negative) from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{{Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(-2, 0)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
		{
			"delete multi(more than one) from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{{Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
		{
			"delete multi(more more than one) from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{{Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
		{
			"delete multi(the sum proportion more than one) from multi",
			passfunc,
			func() {
				proposal.IsAdded = false
				proposal.Treasures = []types.Treasure{{Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(80, 2)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(22, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[0]})
			},
			[]types.Treasure{},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()
			govProposal.Content = proposal

			err := suite.govHandler(suite.ctx, &govProposal)
			tc.expectfunc(err, tc.targetTreasures, tc.msg)

			// check the whitelist with target address list
			actual := suite.app.MintKeeper.GetTreasures(suite.ctx)
			suite.Require().Equal(len(tc.targetTreasures), len(actual), tc.msg)

			for i, _ := range actual {
				suite.Require().Equal(tc.targetTreasures[i], actual[i], tc.msg)
			}

			// reset data
			suite.app.MintKeeper.SetTreasures(suite.ctx, make([]types.Treasure, 0))
		})
	}
}
