package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/global"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/libs/tendermint/types/time"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	gotime "time"
)

func TestDecay(t *testing.T) {
	now := time.Now().Unix()
	after := time.Now().AddDate(0, 0, 52*7).Unix()

	tokens := sdk.NewDec(1000)
	nowDec, err := calculateWeight(now, tokens, 1)
	require.NoError(t, err)
	afterDec, err := calculateWeight(after, tokens, 1)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDec(2), afterDec.Quo(nowDec))
}

type ProposalSuite struct {
	suite.Suite
}

func TestProposalSuite(t *testing.T) {
	suite.Run(t, new(ProposalSuite))
}

func (suite *ProposalSuite) TestNewChangeDistributionTypeProposal() {
	testCases := []struct {
		title         string
		curHeight     int64
		upgradeHeight int64
		quo           int64
	}{
		{"default", 100, 0, 2},
		{"set upgrade height, not reached height", 100, 100, 2},
		{"set upgrade height, reached height", 101, 100, 1},
	}
	formatTime, _ := gotime.Parse("2006-01-02 15:04:05", "2023-06-01 00:00:00")
	require.Equal(suite.T(), formatTime.Unix(), fixedTimeStamp)

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			global.SetGlobalHeight(tc.curHeight)
			tmtypes.InitMilestoneVenus6Height(tc.upgradeHeight)
			tokens := sdk.NewDec(1000)
			nowDec, err := calculateWeight(time.Now().Unix(), tokens, tc.curHeight)
			require.NoError(suite.T(), err)
			afterDec, err := calculateWeight(time.Now().AddDate(0, 0, 52*7).Unix(), tokens, tc.curHeight)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), sdk.NewDec(tc.quo), afterDec.Quo(nowDec))
		})
	}
}
