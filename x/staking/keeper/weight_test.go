package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
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
	nowDec, err := calculateWeight(now, tokens, false)
	require.NoError(t, err)
	afterDec, err := calculateWeight(after, tokens, false)
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
		title    string
		upgraded bool
		quo      int64
	}{
		{"set upgrade height, not reached height", false, 2},
		{"set upgrade height, reached height", true, 1},
	}
	formatTime, _ := gotime.Parse("2006-01-02 15:04:05", "2023-06-01 00:00:00")
	require.Equal(suite.T(), formatTime.Unix(), fixedTimeStamp)

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			tokens := sdk.NewDec(1000)
			nowDec, err := calculateWeight(time.Now().Unix(), tokens, tc.upgraded)
			require.NoError(suite.T(), err)
			afterDec, err := calculateWeight(time.Now().AddDate(0, 0, 52*7).Unix(), tokens, tc.upgraded)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), sdk.NewDec(tc.quo), afterDec.Quo(nowDec))
		})
	}
}
