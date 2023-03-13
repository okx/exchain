package types

import (
	"github.com/okx/okbchain/libs/tendermint/global"
	"math/rand"
	"testing"
	"time"

	"github.com/okx/okbchain/libs/cosmos-sdk/x/gov/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
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
		name        string
		title       string
		description string
		action      string
		extra       string
		err         error
	}{
		{
			"no proposal title",
			"",
			"description",
			"",
			"",
			govtypes.ErrInvalidProposalContent("title is required"),
		},
		{
			"gt max proposal title length",
			RandStr(types.MaxTitleLength + 1),
			"description",
			"",
			"",
			govtypes.ErrInvalidProposalContent("title length is bigger than max title length"),
		},
		{
			"gt max proposal title length",
			RandStr(types.MaxTitleLength),
			"",
			"",
			"",
			govtypes.ErrInvalidProposalContent("description is required"),
		},
		{
			"gt max proposal description length",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength + 1),
			"",
			"",
			govtypes.ErrInvalidProposalContent("description length is bigger than max description length"),
		},
		{
			"no action",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			"",
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's action is required"),
		},
		{
			"action too large",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxExtraActionLength + 1),
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's action length is bigger than max length"),
		},
		{
			"no extra body",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxExtraActionLength),
			"",
			govtypes.ErrInvalidProposalContent("extra proposal's extra is required"),
		},
		{
			"extra too large",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxTitleLength),
			RandStr(govtypes.MaxExtraBodyLength + 1),
			govtypes.ErrInvalidProposalContent("extra proposal's extra body length is bigger than max length"),
		},
		{
			"unknown action",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			RandStr(govtypes.MaxTitleLength),
			RandStr(govtypes.MaxExtraBodyLength),
			ErrUnknownExtraProposalAction,
		},
		{
			"ActionNextBlockUpdate, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionNextBlockUpdate,
			"{dfafdasf}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionNextBlockUpdate, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionNextBlockUpdate,
			"{\"block_nudm\":200}",
			ErrCodeInvalidHeight,
		},
		{
			"ActionNextBlockUpdate, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionNextBlockUpdate,
			"{\"block_num\":100}",
			ErrCodeInvalidHeight,
		},
		{
			"ActionNextBlockUpdate, ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionNextBlockUpdate,
			"{\"block_num\":101}",
			nil,
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{dfafdasf}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coins\":{\"denom\":\"okb\",\"amount\":\"1.000000000000000000\"}}",
			ErrExtraProposalParams("coin is nil"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okb\",\"aamount\":\"1.000000000000000000\"}}",
			ErrExtraProposalParams("coin is nil"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"ddenom\":\"okb\",\"amount\":\"1.000000000000000000\"}}",
			ErrExtraProposalParams("coin is nil"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":[{}]}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionMintedPerBlock, error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coins\":[]}",
			ErrExtraProposalParams("coin is nil"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":[]}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionMintedPerBlock, error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okb\",\"amount\":\"-1.000000000000000000\"}}",
			ErrExtraProposalParams("coin is negative"),
		},
		{
			"ActionMintedPerBlock, error",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okbb\",\"amount\":\"-1.000000000000000000\"}}",
			ErrExtraProposalParams("coin is nil"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okbb\",\"amount\":-1}}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionMintedPerBlock, error json",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okbb\",\"amount\":\"dfads\"}}",
			ErrExtraProposalParams("parse json error"),
		},
		{
			"ActionMintedPerBlock, ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okb\",\"amount\":\"1.000000000000000000\"}}",
			nil,
		},
		{
			"ActionMintedPerBlock, ok",
			RandStr(types.MaxTitleLength),
			RandStr(types.MaxDescriptionLength),
			ActionMintedPerBlock,
			"{\"coin\":{\"denom\":\"okb\",\"amount\":\"0.000000000000000000\"}}",
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			global.SetGlobalHeight(100)
			proposal := NewExtraProposal(tc.title, tc.description, tc.action, tc.extra)
			require.Equal(suite.T(), tc.title, proposal.GetTitle())
			require.Equal(suite.T(), tc.description, proposal.GetDescription())
			require.Equal(suite.T(), RouterKey, proposal.ProposalRoute())
			require.Equal(suite.T(), ProposalTypeExtra, proposal.ProposalType())
			require.NotPanics(suite.T(), func() {
				_ = proposal.String()
			})

			err := proposal.ValidateBasic()
			require.Equal(suite.T(), tc.err, err)
		})
	}
}
