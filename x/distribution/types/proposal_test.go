package types

import (
	"github.com/okex/exchain/x/gov/types"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/stretchr/testify/require"
)

func TestNewCommunityPoolSpendProposal(t *testing.T) {
	title := "Withdraw coins"
	description := "Want to get some coins as reward"
	recipient := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
	proposal := NewCommunityPoolSpendProposal(title, description, recipient, sdk.NewCoins(amount))

	require.Equal(t, title, proposal.GetTitle())
	require.Equal(t, description, proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, ProposalTypeCommunityPoolSpend, proposal.ProposalType())
	require.Nil(t, proposal.ValidateBasic())
	require.NotPanics(t, func() {
		_ = proposal.String()
	})

	proposal.Title = ""
	require.Error(t, proposal.ValidateBasic())
	proposal.Title = title
	proposal.Amount = sdk.SysCoins{sdk.SysCoin{Denom: "UNKNOWN", Amount: sdk.OneDec()}}
	require.Error(t, proposal.ValidateBasic())
	proposal.Amount = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()))
	proposal.Recipient = nil
	require.Error(t, proposal.ValidateBasic())
}

//TODO migrate to a public module
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
