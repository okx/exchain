package ut

import (
	"testing"

	"github.com/okx/okbchain/x/gov/keeper"
	"github.com/okx/okbchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestProposalHandlerRouter_AddRoute(t *testing.T) {
	// nolint
	_, _, k, _, _ := CreateTestInput(t, false, 1000)
	govProposalHandlerRouter := keeper.NewProposalHandlerRouter()

	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute("@###", k)
	})

	govProposalHandlerRouter.AddRoute(types.RouterKey, k)

	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute(types.RouterKey, k)
	})

	govProposalHandlerRouter.Seal()
	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute(types.RouterKey, k)
	})
}

func TestProposalHandlerRouter_GetRoute(t *testing.T) {
	govProposalHandlerRouter := keeper.NewProposalHandlerRouter()
	require.Panics(t, func() {
		govProposalHandlerRouter.GetRoute(types.RouterKey)
	})
}

func TestProposalHandlerRouter_Seal(t *testing.T) {
	govProposalHandlerRouter := keeper.NewProposalHandlerRouter()
	govProposalHandlerRouter.Seal()
	require.Panics(t, func() {
		govProposalHandlerRouter.Seal()
	})
}
