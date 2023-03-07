// nolint
package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

const (
	CodeInvalidDistributionType                   uint32 = 67819
	CodeEmptyDelegationDistInfo                   uint32 = 67820
	CodeEmptyValidatorDistInfo                    uint32 = 67821
	CodeEmptyDelegationVoteValidator              uint32 = 67822
	CodeZeroDelegationShares                      uint32 = 67823
	CodeNotSupportWithdrawDelegationRewards       uint32 = 67824
	CodeNotSupportDistributionProposal            uint32 = 67825
	CodeDisabledWithdrawRewards                   uint32 = 67826
	CodeNotSupportWithdrawRewardEnabledProposal   uint32 = 67827
	CodeProposerMustBeValidator                   uint32 = 67828
	CodeNotSupportRewardTruncatePrecisionProposal uint32 = 67829
	CodeOutOfRangeRewardTruncatePrecision         uint32 = 67830
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

func ErrCodeEmptyDelegationVoteValidator() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyDelegationVoteValidator, "delegation not vote the validator")
}

func ErrCodeZeroDelegationShares() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeZeroDelegationShares, "zero delegation shares")
}

func ErrCodeNotSupportWithdrawDelegationRewards() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotSupportWithdrawDelegationRewards, "not support withdraw delegation rewards")
}

func ErrCodeNotSupportDistributionProposal() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotSupportDistributionProposal, "do not support, distribution proposal invalid")
}

func ErrCodeDisabledWithdrawRewards() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDisabledWithdrawRewards, "withdraw rewards disabled")
}

func ErrCodeProposerMustBeValidator() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeProposerMustBeValidator, "the proposal of proposer must be validator")
}

func ErrCodeRewardTruncatePrecision() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOutOfRangeRewardTruncatePrecision, "reward truncate precision out of range [0,18]")
}
