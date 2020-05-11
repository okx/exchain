package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	err := ErrInitialDepositNotEnough(DefaultCodespace, "")
	require.Equal(t, CodeInitialDepositNotEnough, err.Code())

	err = ErrUnknownProposal(DefaultCodespace, 0)
	require.Equal(t, CodeUnknownProposal, err.Code())

	err = ErrInvalidateProposalStatus(DefaultCodespace, "")
	require.Equal(t, CodeInvalidProposalStatus, err.Code())

	err = ErrInvalidHeight(DefaultCodespace, 100, 100, 100)
	require.Equal(t, CodeInvalidHeight, err.Code())

	err = ErrInvalidProposer(DefaultCodespace, "")
	require.Equal(t, CodeInvalidProposer, err.Code())

}
