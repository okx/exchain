package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/gov/types"
)

func TestProposalHandlerRouter_AddRoute(t *testing.T) {
	// nolint
	_, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	govProposalHandlerRouter := NewProposalHandlerRouter()

	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute("@###", keeper)
	})

	govProposalHandlerRouter.AddRoute(types.RouterKey, keeper)

	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute(types.RouterKey, keeper)
	})

	govProposalHandlerRouter.Seal()
	require.Panics(t, func() {
		govProposalHandlerRouter.AddRoute(types.RouterKey, keeper)
	})
}

func TestProposalHandlerRouter_GetRoute(t *testing.T) {
	govProposalHandlerRouter := NewProposalHandlerRouter()
	require.Panics(t, func() {
		govProposalHandlerRouter.GetRoute(types.RouterKey)
	})
}

func TestProposalHandlerRouter_Seal(t *testing.T) {
	govProposalHandlerRouter := NewProposalHandlerRouter()
	govProposalHandlerRouter.Seal()
	require.Panics(t, func() {
		govProposalHandlerRouter.Seal()
	})
}
