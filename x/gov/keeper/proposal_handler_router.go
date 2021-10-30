package keeper

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"

	"github.com/okex/exchain/x/gov/types"
)

var (
	_ ProposalHandlerRouter = (*proposalHandlerRouter)(nil)
)

// ProposalHandlerRouter defines the interface of managing the proposal handler of proposal
type ProposalHandlerRouter interface {
	AddRoute(r string, mp ProposalHandler) (mpr ProposalHandlerRouter)
	HasRoute(r string) bool
	GetRoute(path string) (h ProposalHandler)
	Seal()
}

// ProposalHandler defines the interface handler in different periods of proposal
type ProposalHandler interface {
	GetMinDeposit(ctx sdk.Context, content types.Content) sdk.SysCoins
	GetMaxDepositPeriod(ctx sdk.Context, content types.Content) time.Duration
	GetVotingPeriod(ctx sdk.Context, content types.Content) time.Duration
	CheckMsgSubmitProposal(ctx sdk.Context, msg types.MsgSubmitProposal) sdk.Error
	AfterSubmitProposalHandler(ctx sdk.Context, proposal types.Proposal)
	VoteHandler(ctx sdk.Context, proposal types.Proposal, vote types.Vote) (string, sdk.Error)
	AfterDepositPeriodPassed(ctx sdk.Context, proposal types.Proposal)
	RejectedHandler(ctx sdk.Context, content types.Content)
}

type proposalHandlerRouter struct {
	routes map[string]ProposalHandler
	sealed bool
}

// nolint
func NewProposalHandlerRouter() ProposalHandlerRouter {
	return &proposalHandlerRouter{
		routes: make(map[string]ProposalHandler),
	}
}

// Seal seals the router which prohibits any subsequent route handlers to be
// added. Seal will panic if called more than once.
func (phr *proposalHandlerRouter) Seal() {
	if phr.sealed {
		panic("router already sealed")
	}
	phr.sealed = true
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
func (phr *proposalHandlerRouter) AddRoute(path string, mp ProposalHandler) ProposalHandlerRouter {
	if phr.sealed {
		panic("router sealed; cannot add route handler")
	}

	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if phr.HasRoute(path) {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	phr.routes[path] = mp
	return phr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (phr *proposalHandlerRouter) HasRoute(path string) bool {
	return phr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (phr *proposalHandlerRouter) GetRoute(path string) ProposalHandler {
	if !phr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return phr.routes[path]
}
