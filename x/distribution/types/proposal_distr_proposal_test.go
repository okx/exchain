package types

import (
	"github.com/okex/exchain/x/gov/types"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func TestNewChangeDistributionTypeProposal(t *testing.T) {
	title := RandStr(types.MaxTitleLength)
	description := RandStr(types.MaxDescriptionLength)
	distrType := DistributionTypeOffChain
	proposal := NewChangeDistributionTypeProposal(title, description, distrType)

	//expect success
	require.Equal(t, title, proposal.GetTitle())
	require.Equal(t, description, proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, ProposalTypeChangeDistributionType, proposal.ProposalType())
	require.Nil(t, proposal.ValidateBasic())
	require.NotPanics(t, func() {
		_ = proposal.String()
	})

	//expect failed,Title is nill
	proposal.Title = ""
	require.Error(t, proposal.ValidateBasic())

	//expect failed,Title is nill
	proposal.Title = RandStr(types.MaxTitleLength + 1)
	require.Error(t, proposal.ValidateBasic())

	//expect failed,Title len lg MaxTitleLength
	proposal.Title = RandStr(types.MaxTitleLength + 1)
	require.Error(t, proposal.ValidateBasic())

	//expect failed,Description is nill
	proposal.Title = RandStr(types.MaxTitleLength)
	proposal.Description = ""
	require.Error(t, proposal.ValidateBasic())

	//expect failed,Description lg MaxDescriptionLength
	proposal.Title = RandStr(types.MaxTitleLength)
	proposal.Description = RandStr(types.MaxDescriptionLength + 1)
	require.Error(t, proposal.ValidateBasic())

	//expect failed, type is 2
	proposal.Title = RandStr(types.MaxTitleLength)
	proposal.Description = RandStr(types.MaxDescriptionLength)
	proposal.Type = 2
	require.Error(t, proposal.ValidateBasic())

}
