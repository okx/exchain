// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace            string          = ModuleName
	CodeInvalidInput            uint32          = 103
	CodeNoValidatorCommission   uint32          = 105
	CodeSetWithdrawAddrDisabled uint32          = 106
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
