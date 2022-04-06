package keeper_test

import (
	"testing"
	"time"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	evm_types "github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/suite"
)

type TreasuresTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *app.OKExChainApp
	stateDB *evm_types.CommitStateDB
	codec   *codec.Codec

	handler sdk.Handler
}

func (suite *TreasuresTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.codec = codec.New()

	suite.app.MintKeeper.SetParams(suite.ctx, types.DefaultParams())
	suite.app.MintKeeper.SetMinter(suite.ctx, types.InitialMinterCustom())
}

var (
	treasure1 = types.NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 2))
	treasure2 = types.NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 2))
	treasure3 = types.NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 2))
	treasure4 = types.NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 2))
	treasure5 = types.NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(1, 1))
	treasures = []types.Treasure{*treasure1, *treasure2, *treasure3, *treasure4, *treasure5}
)

func TestTreasuresTestSuite(t *testing.T) {
	suite.Run(t, new(TreasuresTestSuite))
}

func (suite *TreasuresTestSuite) TestGetSetTreasures() {
	input := []types.Treasure{}
	testCases := []struct {
		msg      string
		prepare  func()
		expected []types.Treasure
	}{
		{
			msg: "set one treasure into empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
			},
			expected: []types.Treasure{treasures[1]},
		},
		{
			msg: "set one treasure into db which has one",
			prepare: func() {
				input = []types.Treasure{treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1]})
			},
			expected: []types.Treasure{treasures[0]},
		},
		{
			msg: "set one treasure(exist) into db which has one",
			prepare: func() {
				input = []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			expected: []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}},
		},
		{
			msg: "set one treasure into db which has multi",
			prepare: func() {
				input = []types.Treasure{treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[2]})
			},
			expected: []types.Treasure{treasures[0]},
		},
		{
			msg: "set one treasure(exist) into db which has multi",
			prepare: func() {
				input = []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[2], treasures[0]})
			},
			expected: []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}},
		},
		{
			msg: "set multi treasure into empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[2], treasures[0], treasures[3], treasures[4]}
			},
			expected: []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1], treasures[0]},
		},
		{
			msg: "set multi treasure into db which has one",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[2], treasures[3], treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			expected: []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1]},
		},
		{
			msg: "set multi treasure(part exist) into db which has one",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[2], types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}, treasures[3], treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			expected: []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1], types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}},
		}, {
			msg: "set multi treasure into db which has multi",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[2]})
			},
			expected: []types.Treasure{treasures[3], treasures[0]},
		},
		{
			msg: "set multi treasure(part exist) into db which has multi",
			prepare: func() {
				input = []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}, treasures[1], treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[2], treasures[0]})
			},
			expected: []types.Treasure{treasures[3], treasures[1], types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}},
		},
		{
			msg: "set multi treasure(all exist) into db which has multi",
			prepare: func() {
				input = []types.Treasure{types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}, treasures[1], treasures[3]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[3], treasures[0]})
			},
			expected: []types.Treasure{treasures[3], treasures[1], types.Treasure{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(4, 1)}},
		},
	}

	for _, tc := range testCases {
		tc.prepare()
		suite.app.MintKeeper.SetTreasures(suite.ctx, input)
		actual := suite.app.MintKeeper.GetTreasures(suite.ctx)
		suite.Require().Equal(tc.expected, actual, tc.msg)
		suite.app.MintKeeper.SetTreasures(suite.ctx, make([]types.Treasure, 0))
	}
}

func (suite *TreasuresTestSuite) TestAllocateTokenToTreasure() {
	input := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(5, 1)))
	testCases := []struct {
		msg        string
		treasures  func() []types.Treasure
		prepare    func(trs []types.Treasure)
		expectfunc func(trs []types.Treasure, remain sdk.Coins, err error, msg string)
	}{
		{
			msg: "0.5 coin allocate 0%(one)",
			treasures: func() []types.Treasure {
				return []types.Treasure{{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(0, 2)}}
			},
			prepare: func(trs []types.Treasure) {
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(50, 2))), remain, msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(50, 2).MulTruncate(trs[i].Proportion))).String(), acc.GetCoins().String(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
		{
			msg: "0.5 coin allocate 50%(one)",
			treasures: func() []types.Treasure {
				return []types.Treasure{{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(50, 2)}}
			},
			prepare: func(trs []types.Treasure) {
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(25, 2))), remain, msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(25, 2))), acc.GetCoins(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
		{
			msg: "0.5 coin allocate 50%(multi)",
			treasures: func() []types.Treasure {
				return []types.Treasure{
					{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(30, 2)},
					{Address: sdk.AccAddress([]byte{0x02}), Proportion: sdk.NewDecWithPrec(20, 2)},
				}
			},
			prepare: func(trs []types.Treasure) {
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(25, 2))), remain, msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(50, 2).MulTruncate(trs[i].Proportion))), acc.GetCoins(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
		{
			msg: "0.5 coin allocate 60%(multi)",
			treasures: func() []types.Treasure {
				return []types.Treasure{
					{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(30, 2)},
					{Address: sdk.AccAddress([]byte{0x02}), Proportion: sdk.NewDecWithPrec(30, 2)},
				}
			},
			prepare: func(trs []types.Treasure) {
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(20, 2))), remain, msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(50, 2).MulTruncate(trs[i].Proportion))), acc.GetCoins(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
		{
			msg: "0.5 coin allocate 100%(multi)",
			treasures: func() []types.Treasure {
				return []types.Treasure{
					{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(60, 2)},
					{Address: sdk.AccAddress([]byte{0x02}), Proportion: sdk.NewDecWithPrec(40, 2)},
				}
			},
			prepare: func(trs []types.Treasure) {
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(0, 2))).String(), remain.String(), msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(50, 2).MulTruncate(trs[i].Proportion))), acc.GetCoins(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
		{
			msg: "0.7 coin allocate 60%(multi)",
			treasures: func() []types.Treasure {
				return []types.Treasure{
					{Address: sdk.AccAddress([]byte{0x01}), Proportion: sdk.NewDecWithPrec(30, 2)},
					{Address: sdk.AccAddress([]byte{0x02}), Proportion: sdk.NewDecWithPrec(30, 2)},
				}
			},
			prepare: func(trs []types.Treasure) {
				input = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(7, 1)))
				err := suite.app.MintKeeper.MintCoins(suite.ctx, input)
				suite.Require().NoError(err)
				suite.app.MintKeeper.SetTreasures(suite.ctx, trs)
			},
			expectfunc: func(trs []types.Treasure, remain sdk.Coins, err error, msg string) {
				suite.Require().NoError(err, msg)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(28, 2))), remain, msg)
				for i, _ := range trs {
					acc := suite.app.AccountKeeper.GetAccount(suite.ctx, trs[i].Address)
					suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(70, 2).MulTruncate(trs[i].Proportion))), acc.GetCoins(), msg)
					err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroDec())))
					suite.Require().NoError(err, msg)
					suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
				}
			},
		},
	}

	for _, tc := range testCases {
		input = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecWithPrec(5, 1)))
		//prepare environment
		treasuresInput := tc.treasures()
		tc.prepare(treasuresInput)

		// handler test case
		remain, err := suite.app.MintKeeper.AllocateTokenToTreasure(suite.ctx, input)

		// verify
		tc.expectfunc(treasuresInput, remain, err, tc.msg)

		// reset environment
		suite.app.MintKeeper.SetTreasures(suite.ctx, make([]types.Treasure, 0))
	}
}

func (suite *TreasuresTestSuite) TestUpdateTreasures() {
	input := []types.Treasure{}
	normal := func(isPass bool, expected []types.Treasure, msg string, err error) {
		if isPass {
			suite.Require().NoError(err, msg)
			actual := suite.app.MintKeeper.GetTreasures(suite.ctx)
			suite.Require().Equal(expected, actual, msg)
		} else {
			suite.Require().Error(err, msg)
		}
	}
	treasureError := func(isPass bool, expected []types.Treasure, msg string, err error) {
		suite.Require().False(isPass, msg)
		suite.Require().Error(err, msg)
		suite.Require().Contains(err.Error(), "treasure proportion should non-negative and less than one", msg)
	}
	sumProportionError := func(isPass bool, expected []types.Treasure, msg string, err error) {
		suite.Require().False(isPass, msg)
		suite.Require().Error(err, msg)
		suite.Require().Contains(err.Error(), "the sum of treasure proportion should non-negative and less than one", msg)
	}
	testCases := []struct {
		msg        string
		prepare    func()
		isPass     bool
		expected   []types.Treasure
		expectFunc func(isPass bool, expected []types.Treasure, msg string, err error)
	}{
		{
			msg: "insert one treasure into empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[1]},
			expectFunc: normal,
		},
		{
			msg: "insert one treasure(proportion is negative) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure(proportion is more than one) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(2)}}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure into db which has one",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[1], treasures[0]},
			expectFunc: normal,
		},
		{
			msg: "insert one treasure(proportion is negative) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure(proportion is more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure(sum proportion is more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert one treasure into db which has multi",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2], treasures[1], treasures[0]},
			expectFunc: normal,
		},
		{
			msg: "insert one treasure(proportion is negative) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure(proportion is more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert one treasure(sum proportion is more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert multi treasure into empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[3], treasures[4], treasures[2], treasures[0]}
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1], treasures[0]},
			expectFunc: normal,
		},
		{
			msg: "insert multi treasure(part negative) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, treasures[3], treasures[4], treasures[2], treasures[0]}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all negative) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[3].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[4].Address, Proportion: sdk.NewDec(-1)}}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(more than one) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(2)}, treasures[2]}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all more than one) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(2)}, {Address: treasures[2].Address, Proportion: sdk.NewDec(2)}}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(the sum proportion more than one) into empty db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(5, 1)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(6, 1)}}
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert multi treasure into db which has one",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[3], treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1], treasures[0]},
			expectFunc: normal,
		},
		{
			msg: "insert multi treasure(negative) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, treasures[3], treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(part negative) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[3].Address, Proportion: sdk.NewDec(-1)}, treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all negative) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[3].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[4].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[3], treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(part more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(the sum proportion more than one) into db which has one",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(2, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert multi treasure into db which has one (the result treasures's sum proportion more than one)",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert multi treasure into db which has multi",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[3], treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[4], treasures[3], treasures[2], treasures[1], treasures[0]},
			expectFunc: normal,
		},
		{
			msg: "insert multi treasure(negative) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, treasures[3], treasures[4], treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(part negative) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[3].Address, Proportion: sdk.NewDec(-1)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all negative) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[3].Address, Proportion: sdk.NewDec(-1)}, {Address: treasures[4].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[3], treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(part more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[4]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(all more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "insert multi treasure(the sum proportion more than one) into db which has multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[4].Address, Proportion: sdk.NewDecWithPrec(2, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "insert multi treasure into db which has one (the result treasures's sum proportion more than multi)",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[3].Address, Proportion: sdk.NewDecWithPrec(8, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "update one treasure with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
			expectFunc: normal,
		},
		{
			msg: "update one treasure(negative) with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update one treasure(more than one) with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDec(3)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update one treasure with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2], {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
			expectFunc: normal,
		},
		{
			msg: "update one treasure(negative) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDec(-1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update one treasure(more than one) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDec(5)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update one treasure(the sum proportion more than one) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "update multi treasure with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}, treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2], {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
			expectFunc: normal,
		},
		{
			msg: "update multi treasure(negative) with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(more than one) with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(the sum proportion more than one) with one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(1, 0)}, treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "update multi treasure(part) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[3]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[3], {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
			expectFunc: normal,
		},
		{
			msg: "update multi treasure(part negtive) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(-10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[3]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(part more than one) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[3]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(part the sum proportion more than one) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(2, 1)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[3]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
		{
			msg: "update multi treasure(all) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{{Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}, {Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 2)}},
			expectFunc: normal,
		},
		{
			msg: "update multi treasure(all negative1) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(-10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(all negative2) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-8, 2)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(all more than one 1) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(all more than one 2) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(10, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: treasureError,
		},
		{
			msg: "update multi treasure(all the sum proportion more than one) with multi db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 1)}, {Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(23, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: sumProportionError,
		},
	}

	for _, tc := range testCases {
		//prepare environment
		tc.prepare()

		// handler test case
		err := suite.app.MintKeeper.UpdateTreasures(suite.ctx, input)
		//verify test case expect
		tc.expectFunc(tc.isPass, tc.expected, tc.msg, err)

		// reset environment
		suite.app.MintKeeper.SetTreasures(suite.ctx, make([]types.Treasure, 0))
	}
}

func (suite *TreasuresTestSuite) TestDeleteTreasures() {
	input := []types.Treasure{}
	normal := func(isPass bool, expected []types.Treasure, msg string, err error) {
		if isPass {
			suite.Require().NoError(err)
			actual := suite.app.MintKeeper.GetTreasures(suite.ctx)
			suite.Require().Equal(expected, actual, msg)
		} else {
			suite.Require().Error(err)
		}
	}
	unexistError := func(isPass bool, expected []types.Treasure, msg string, err error) {
		suite.Require().False(isPass, msg)
		suite.Require().Error(err, msg)
		suite.Require().Contains(err.Error(), "because it's not exist from treasures", msg)
	}
	testCases := []struct {
		msg        string
		prepare    func()
		isPass     bool
		expected   []types.Treasure
		expectFunc func(isPass bool, expected []types.Treasure, msg string, err error)
	}{
		{
			msg: "delete one treasure from empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasure from empty db",
			prepare: func() {
				input = []types.Treasure{treasures[1], treasures[0]}
			},
			isPass:     false,
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasure from one db",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1]})
			},
			isPass:     true,
			expected:   nil,
			expectFunc: normal,
		},
		{
			msg: "delete one treasure(not exist) from one db",
			prepare: func() {
				input = []types.Treasure{treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasure(negative) from one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   nil,
			expectFunc: normal,
		},
		{
			msg: "delete one treasure(negative,not exist) from one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasure(more than one) from one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     true,
			expected:   nil,
			expectFunc: normal,
		},
		{
			msg: "delete one treasure(more than one,not exist) from one db",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasures from multi",
			prepare: func() {
				input = []types.Treasure{treasures[0]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[1]},
			expectFunc: normal,
		},
		{
			msg: "delete one treasures(no exist) from multi",
			prepare: func() {
				input = []types.Treasure{treasures[2]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasures(negative) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[1]},
			expectFunc: normal,
		},
		{
			msg: "delete one treasures(negative,not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete one treasures(more than one) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[1]},
			expectFunc: normal,
		},
		{
			msg: "delete one treasures(more than one,not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[2].Address, Proportion: sdk.NewDecWithPrec(2, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     false,
			expected:   []types.Treasure{},
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures from multi",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0]})
			},
			isPass:     true,
			expected:   nil,
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part) from multi",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,part negative) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,all negative) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, {Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,part more than one) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,all more than one) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(8, 0)}, {Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(8, 0)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,the sum proportion more than one) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(83, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1], treasures[0], treasures[2]})
			},
			isPass:     true,
			expected:   []types.Treasure{treasures[2]},
			expectFunc: normal,
		},
		{
			msg: "delete multi treasures(part,part not exist) from multi",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,negative part not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,more than one part not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,the sumproportion more than one, part not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(83, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[0], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,all not exist) from multi",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,negative,all not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(-1, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,more than one all not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 0)}, treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures(part,the sumproportion more than one, all not exist) from multi",
			prepare: func() {
				input = []types.Treasure{{Address: treasures[0].Address, Proportion: sdk.NewDecWithPrec(2, 1)}, {Address: treasures[1].Address, Proportion: sdk.NewDecWithPrec(83, 2)}}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[3], treasures[2]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
		{
			msg: "delete multi treasures from one",
			prepare: func() {
				input = []types.Treasure{treasures[0], treasures[1]}
				suite.app.MintKeeper.SetTreasures(suite.ctx, []types.Treasure{treasures[1]})
			},
			isPass:     false,
			expected:   nil,
			expectFunc: unexistError,
		},
	}

	for _, tc := range testCases {
		//prepare environment
		tc.prepare()

		// handler test case
		err := suite.app.MintKeeper.DeleteTreasures(suite.ctx, input)
		//verify test case expect
		tc.expectFunc(tc.isPass, tc.expected, tc.msg, err)

		// reset environment
		suite.app.MintKeeper.SetTreasures(suite.ctx, make([]types.Treasure, 0))
	}
}
