// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace          							string          = ModuleName
	CodeInvalidInput									uint32          = 67800
	CodeNoValidatorCommission							uint32          = 67801
	CodeSetWithdrawAddrDisabled							uint32          = 67802
	CodeInvalideData									uint32		  = 67803
	CodeInvalideRoute									uint32		  = 67804
	CodeUnmarshalJSONFailed								uint32		  = 67805
	CodeInvalideBasic									uint32		  = 67806
	CodeWithdrawValidatorRewardsAndCommissionFailed		uint32		  = 67807
	CodeAccAddressFromBech32Failed						uint32		  = 67808
	CodeValAddressFromBech32							uint32		  = 67809
	CodeReadRESTReqFailed								uint32		  = 67810
	CodeSendCoinsFromModuleToAccountFailed				uint32		  = 67811
	CodeUnknownRequest									uint32		  = 67812
	CodeUnauthorized									uint32		  = 67813
	CodeSetWithdrawAddrFailed							uint32		  = 67814
	CodeWithdrawValidatorCommissionFailed				uint32		  = 67815
)

func ErrNilDelegatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "delegator address is nil")
}
func ErrNilWithdrawAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "withdraw address is nil")
}
func ErrNilValidatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "validator address is nil")
}
func ErrNoValidatorCommission() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNoValidatorCommission, "no validator commission to withdraw")
}
func ErrSetWithdrawAddrDisabled() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSetWithdrawAddrDisabled, "set withdraw address disabled")
}
func ErrBadDistribution() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "community pool does not have sufficient coins to distribute")
}
func ErrInvalidProposalAmount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "invalid community pool spend proposal amount")
}
func ErrEmptyProposalRecipient() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "invalid community pool spend proposal recipient")
}
func ErrSendCoinsFromModuleToAccountFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromModuleToAccountFailed, "invalid withdrawAddr or commission")
}
func ErrUnknownRequest() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "incorrectly formatted request data")
}
func ErrUnauthorized() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "blacklisted from receiving external funds")
}
func ErrSetWithdrawAddrFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSetWithdrawAddrFailed, "delegators addr or withdraw addr is invalid")
}
func ERRWithdrawValidatorCommissionFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawValidatorCommissionFailed, "withdraw validator commission failed")
}