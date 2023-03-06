package client

import (
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/x/gov/client/rest"
)

// function to create the rest handler
type RESTHandlerFn func(context.CLIContext) rest.ProposalRESTHandler

// function to create the cli handler
type CLIHandlerFn func(proxy *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command

// The combined type for a proposal handler for both cli and rest
type ProposalHandler struct {
	CLIHandler  CLIHandlerFn
	RESTHandler RESTHandlerFn
}

// NewProposalHandler creates a new ProposalHandler object
func NewProposalHandler(cliHandler CLIHandlerFn, restHandler RESTHandlerFn) ProposalHandler {
	return ProposalHandler{
		CLIHandler:  cliHandler,
		RESTHandler: restHandler,
	}
}
