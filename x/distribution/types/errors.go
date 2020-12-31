// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = ModuleName

	CodeNilDelegatorAddr                            uint32 = 67800
	CodeNoValidatorCommission                       uint32 = 67801
	CodeSetWithdrawAddrDisabled                     uint32 = 67802
	CodeInvalidRoute                                uint32 = 67803
	CodeWithdrawValidatorRewardsAndCommissionFailed uint32 = 67804
	CodeAccAddressFromBech32Failed                  uint32 = 67805
	CodeValAddressFromBech32                        uint32 = 67806
	CodeSendCoinsFromModuleToAccountFailed          uint32 = 67807
	CodeWithdrawValidatorCommissionFailed           uint32 = 67808
	CodeUnknownDistributionMsgType                  uint32 = 67809
	CodeUnknownDistributionCommunityPoolProposaType uint32 = 67810
	CodeUnknownDistributionQueryType                uint32 = 67811
	CodeUnknownDistributionParamType                uint32 = 67812
	CodeWithdrawAddrInBlacklist                     uint32 = 67813
	CodeNilWithdrawAddr                             uint32 = 67814
	CodeNilValidatorAddr                            uint32 = 67815
	CodeBadDistribution                             uint32 = 67816
	CodeInvalidProposalAmount                       uint32 = 67817
	CodeEmptyProposalRecipient                      uint32 = 67818
)

func ErrNilDelegatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNilDelegatorAddr, "delegator address is empty")
}

func ErrNoValidatorCommission() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNoValidatorCommission, "no validator commission to withdraw")
}

func ErrSetWithdrawAddrDisabled() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSetWithdrawAddrDisabled, "set withdraw address disabled")
}

func ErrSendCoinsFromModuleToAccountFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromModuleToAccountFailed, "send coins from module to account failed")
}

func ErrWithdrawValidatorCommissionFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawValidatorCommissionFailed, "withdraw validator commission failed")
}

func ErrUnknownDistributionMsgType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownDistributionMsgType, "unknown distribution message type")
}

func ErrUnknownDistributionCommunityPoolProposaType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownDistributionCommunityPoolProposaType, "unknown community pool proposal type")
}

func ErrUnknownDistributionQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownDistributionQueryType, "unknown distribution query type")
}

func ErrUnknownDistributionParamType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownDistributionParamType, "unknown distribution param type")
}

func ErrWithdrawAddrInblacklist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawAddrInBlacklist, "withdraw address in black list")
}

func ErrNilWithdrawAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNilWithdrawAddr, "withdraw address is empty")
}

func ErrNilValidatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNilValidatorAddr, "validator address is empty")
}

func ErrBadDistribution() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBadDistribution, "community pool does not have sufficient coins to distribute")
}

func ErrInvalidProposalAmount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProposalAmount, "invalid community pool spend proposal amount")
}

func ErrEmptyProposalRecipient() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyProposalRecipient, "invalid community pool spend proposal recipient")
}
