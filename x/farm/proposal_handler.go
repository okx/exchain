package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
	govTypes "github.com/okex/okexchain/x/gov/types"
)

// NewManageWhiteListProposalHandler handles "gov" type message in "farm"
func NewManageWhiteListProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch contentType := proposal.Content.(type) {
		case types.ManageWhiteListProposal:
			return handleManageWhiteListProposal(ctx, k, proposal)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized param proposal content type: %s", contentType))
		}
	}
}

// TODO
func handleManageWhiteListProposal(ctx sdk.Context, keeper *Keeper, proposal *govTypes.Proposal) sdk.Error {
	return nil
}
