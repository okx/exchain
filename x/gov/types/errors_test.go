package types

import (
	sdkerror "github.com/cosmos/cosmos-sdk/types/errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	require.Equal(t, BaseGovError+1, ErrUnknownProposal(DefaultParamspace, 0).(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+10, ErrInvalidateProposalStatus(DefaultParamspace, "").(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+12, ErrInitialDepositNotEnough(DefaultParamspace, "").(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+13, ErrInvalidProposer(DefaultParamspace, "").(*sdkerror.Error).ABCICode())
	require.Equal(t, BaseGovError+14, ErrInvalidHeight(DefaultParamspace, 0, 0, 0).(*sdkerror.Error).ABCICode())
}
