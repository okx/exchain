package client

import (
	sdkGovClient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

type (
	// ProposalHandler is alias of cm gov ProposalHandler
	ProposalHandler = sdkGovClient.ProposalHandler
)

var (
	// NewProposalHandler is alias of cm gov NewProposalHandler
	NewProposalHandler = sdkGovClient.NewProposalHandler
)
