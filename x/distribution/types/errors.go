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

func ErrNilDelegatorAddr(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "delegator address is nil")
}
func ErrNilWithdrawAddr(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "withdraw address is nil")
}
func ErrNilValidatorAddr(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "validator address is nil")
}
func ErrNoValidatorCommission(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeNoValidatorCommission, "no validator commission to withdraw")
}
func ErrSetWithdrawAddrDisabled(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeSetWithdrawAddrDisabled, "set withdraw address disabled")
}
func ErrBadDistribution(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "community pool does not have sufficient coins to distribute")
}
func ErrInvalidProposalAmount(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "invalid community pool spend proposal amount")
}
func ErrEmptyProposalRecipient(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "invalid community pool spend proposal recipient")
}
func ErrSendCoinsFromModuleToAccountFailed(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeSendCoinsFromModuleToAccountFailed, "invalid withdrawAddr or commission")
}
func ErrUnknownRequest(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnknownRequest, "incorrectly formatted request data")
}
func ErrUnauthorized(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnknownRequest, "blacklisted from receiving external funds")
}
func ErrSetWithdrawAddrFailed(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeSetWithdrawAddrFailed, "delegators addr or withdraw addr is invalid")
}
func ERRWithdrawValidatorCommissionFailed(codespace string) sdk.Error {
	return sdkerrors.New(codespace, CodeWithdrawValidatorCommissionFailed, "withdraw validator commission failed")
}