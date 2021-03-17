package types

import (
	"strings"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	govtypes "github.com/okex/okexchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

const (
	expectedTitle          = "default title"
	expectedDescription    = "default description"
	expectedProposalString = `ManageContractDeploymentWhitelistProposal:
 Title:					default title
 Description:        	default description
 Type:                	ManageContractDeploymentWhitelist
 DistributorAddr:		okexchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqupa6dx
 IsAdded:				true`
)

func TestProposal(t *testing.T) {
	addr := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	proposal := NewManageContractDeploymentWhitelistProposal(
		expectedTitle,
		expectedDescription,
		addr,
		true,
	)

	require.Equal(t, expectedTitle, proposal.GetTitle())
	require.Equal(t, expectedDescription, proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, proposalTypeManageContractDeploymentWhitelist, proposal.ProposalType())
	require.Equal(t, expectedProposalString, proposal.String())

	// validateBasic check
	require.NoError(t, proposal.ValidateBasic())
	// check for error
	// 1. empty title
	proposal.Title = ""
	require.Error(t, proposal.ValidateBasic())
	// 2. overlong title
	var b strings.Builder
	for i := 0; i < govtypes.MaxTitleLength+1; i++ {
		b.WriteByte('a')
	}
	proposal.Title = b.String()
	require.Error(t, proposal.ValidateBasic())
	// 3. empty description
	proposal.Description = ""
	proposal.Title = expectedTitle
	require.Error(t, proposal.ValidateBasic())
	// 4. overlong description
	b.Reset()
	for i := 0; i < govtypes.MaxDescriptionLength+1; i++ {
		b.WriteByte('a')
	}
	proposal.Description = b.String()
	require.Error(t, proposal.ValidateBasic())
	// 5. empty address
	proposal.DistributorAddr = nil
	proposal.Description = expectedDescription
	require.Error(t, proposal.ValidateBasic())
}
