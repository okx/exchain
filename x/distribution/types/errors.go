// nolint
package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	CodeInvalidDistributionType             uint32 = 67819
	CodeEmptyDelegationDistInfo             uint32 = 67820
	CodeEmptyValidatorDistInfo              uint32 = 67821
	CodeEmptyDelegationVoteValidator        uint32 = 67822
	CodeZeroDelegationShares                uint32 = 67823
	CodeNotSupportWithdrawDelegationRewards uint32 = 67824
	CodeNotSupportDistributionProposal      uint32 = 67825
)

func ErrInvalidDistributionType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDistributionType, "invalid change distribution type")
}

func ErrCodeEmptyDelegationDistInfo() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyDelegationDistInfo, "no delegation distribution info")
}

func ErrCodeEmptyValidatorDistInfo() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyValidatorDistInfo, "no validator distribution info")
}

func ErrCodeCodeEmptyDelegationVoteValidator() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyDelegationVoteValidator, "delegation not vote validator")
}

func ErrCodeZeroDelegationShares() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeZeroDelegationShares, "zero delegation shares")
}

func ErrCodeNotSupportWithdrawDelegationRewards() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotSupportWithdrawDelegationRewards, "not support withdraw delegation rewards")
}

// ErrCodeNotSupportDistributionProposal returns an error when not support
func ErrCodeNotSupportDistributionProposal() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotSupportDistributionProposal, "not support distribution proposal")
}
