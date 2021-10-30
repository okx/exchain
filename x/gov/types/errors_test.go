package types

import (
	sdkerror "github.com/okex/exchain/dependence/cosmos-sdk/types/errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	require.Equal(t, BaseGovError+1, ErrUnknownProposal( 0).(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+7, ErrInvalidateProposalStatus().(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+8, ErrInitialDepositNotEnough("").(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+9, ErrInvalidProposer().(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+10, ErrInvalidHeight(0, 0, 0).(*sdkerror.Error).ABCICode())
}
