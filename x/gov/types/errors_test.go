package types

import (
	sdkerror "github.com/cosmos/cosmos-sdk/types/errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	require.Equal(t, BaseGovError+1, ErrUnknownProposal( 0).(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+11, ErrInvalidateProposalStatus().(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+12, ErrInitialDepositNotEnough("").(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+13, ErrInvalidProposer().(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+14, ErrInvalidHeight(0, 0, 0).(*sdkerror.Error).ABCICode())
}
