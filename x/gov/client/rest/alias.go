package rest

import (
	sdkGovRest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

type (
	// ProposalRESTHandler is alias of cm gov ProposalRESTHandler
	ProposalRESTHandler = sdkGovRest.ProposalRESTHandler
)

var (
	// RegisterRoutes is alias of cm gov RegisterRoutes
	RegisterRoutes = sdkGovRest.RegisterRoutes
)
