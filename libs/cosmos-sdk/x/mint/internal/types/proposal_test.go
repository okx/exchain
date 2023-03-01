package types

import (
	"math/rand"
	"testing"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
	"github.com/okex/exchain/libs/tendermint/global"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	exgovtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProposalSuite struct {
	suite.Suite
}

func TestProposalSuite(t *testing.T) {
	suite.Run(t, new(ProposalSuite))
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func (suite *ProposalSuite) TestNewChangeDistributionTypeProposal() {
	testCases := []struct {
		title               string
		setMilestoneHeight  func()
		proposalTitle       string
		proposalDescription string
		blockNum            uint64
		err                 error
	}{
		{
			"no proposal title",
			func() {
				global.SetGlobalHeight(0)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(0)
			},
			"",
			"description",
			0,
			exgovtypes.ErrInvalidProposalContent("title is required"),
		},
		{
			"gt max proposal title length",
			func() {
				global.SetGlobalHeight(0)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(0)
			},
			RandStr(types.MaxTitleLength + 1),
			"description",
			0,
			exgovtypes.ErrInvalidProposalContent("title length is bigger than max title length"),
		},
		{
			"gt max proposal title length",
			func() {
				global.SetGlobalHeight(0)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(0)
			},
			RandStr(types.MaxTitleLength),
			"",
			0,
			exgovtypes.ErrInvalidProposalContent("description is required"),
		},
		{
			"gt max proposal description length",
			func() {
				global.SetGlobalHeight(0)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(0)
			},
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength + 1),
			0,
			exgovtypes.ErrInvalidProposalContent("description length is bigger than max description length"),
		},
		{
			"invalid height",
			func() {
				global.SetGlobalHeight(100)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(-1)
			},
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			100,
			ErrCodeInvalidHeight,
		},
		{
			"valid height",
			func() {
				global.SetGlobalHeight(100)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(-1)
			},
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			101,
			nil,
		},
		{
			"ok",
			func() {
				global.SetGlobalHeight(0)
				tmtypes.UnittestOnlySetMilestoneVenus5Height(0)
			},
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			0,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			tc.setMilestoneHeight()
			title := tc.proposalTitle
			description := tc.proposalDescription
			proposal := NewModifyNextBlockUpdateProposal(title, description, tc.blockNum)

			require.Equal(suite.T(), title, proposal.GetTitle())
			require.Equal(suite.T(), description, proposal.GetDescription())
			require.Equal(suite.T(), RouterKey, proposal.ProposalRoute())
			require.Equal(suite.T(), proposalTypeModifyNextBlockUpdate, proposal.ProposalType())
			require.NotPanics(suite.T(), func() {
				_ = proposal.String()
			})

			err := proposal.ValidateBasic()
			require.Equal(suite.T(), tc.err, err)
		})
	}
}
