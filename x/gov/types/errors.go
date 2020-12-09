//nolint
package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = "gov"
	BaseGovError     uint32 = 68000

	CodeUnknownProposal          uint32 = BaseGovError+1
	CodeInactiveProposal         uint32 = BaseGovError+2
	CodeAlreadyActiveProposal    uint32 = BaseGovError+3
	CodeAlreadyFinishedProposal  uint32 = BaseGovError+4
	CodeAddressNotStaked         uint32 = BaseGovError+5
	CodeInvalidContent           uint32 = BaseGovError+6
	CodeInvalidProposalType      uint32 = BaseGovError+7
	CodeInvalidVote              uint32 = BaseGovError+8
	CodeInvalidGenesis           uint32 = BaseGovError+9
	CodeProposalHandlerNotExists uint32 = BaseGovError+10
	CodeInvalidProposalStatus   uint32 = BaseGovError+11
	CodeInitialDepositNotEnough uint32 = BaseGovError+12
	CodeInvalidProposer         uint32 = BaseGovError+13
	CodeInvalidHeight           uint32 = BaseGovError+14
	CodeInsufficientCoins		uint32 = BaseGovError+15
	CodeUnknownRequest			uint32 = BaseGovError+16
	CodeInvalidCoins			uint32 = BaseGovError+17

)

func ErrInactiveProposal(codespace string, proposalID uint64) sdk.Error {
	return sdkerrors.New(codespace, CodeInactiveProposal, fmt.Sprintf("inactive proposal with id %d", proposalID))
}

func ErrAlreadyActiveProposal(codespace string, proposalID uint64) sdk.Error {
	return sdkerrors.New(codespace, CodeAlreadyActiveProposal, fmt.Sprintf("proposal %d has been already active", proposalID))
}

func ErrAlreadyFinishedProposal(codespace string, proposalID uint64) sdk.Error {
	return sdkerrors.New(codespace, CodeAlreadyFinishedProposal, fmt.Sprintf("proposal %d has already passed its voting period", proposalID))
}

func ErrAddressNotStaked(codespace string, address sdk.AccAddress) sdk.Error {
	return sdkerrors.New(codespace, CodeAddressNotStaked, fmt.Sprintf("address %s is not staked and is thus ineligible to vote", address))
}

func ErrInvalidProposalContent(cs string, msg string) sdk.Error {
	return sdkerrors.New(cs, CodeInvalidContent, fmt.Sprintf("invalid proposal content: %s", msg))
}

func ErrInvalidProposalType(codespace string, proposalType string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidProposalType, fmt.Sprintf("proposal type '%s' is not valid", proposalType))
}

func ErrInvalidVote(codespace string, voteOption VoteOption) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption.String()))
}

func ErrInvalidGenesis(codespace string, msg string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidVote, msg)
}

func ErrNoProposalHandlerExists(codespace string, content interface{}) sdk.Error {
	return sdkerrors.New(codespace, CodeProposalHandlerNotExists, fmt.Sprintf("'%T' does not have a corresponding handler", content))
}

func ErrUnknownProposal(codespace string, proposalID uint64) sdk.Error {
	return sdkerrors.New(codespace, CodeUnknownProposal, fmt.Sprintf("unknown proposal with id %d", proposalID))
}

func ErrInvalidateProposalStatus(codespace string, msg string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidProposalStatus, msg)
}

func ErrInitialDepositNotEnough(codespace string, initDeposit string) sdk.Error {
	return sdkerrors.New(codespace, CodeInitialDepositNotEnough,
		fmt.Sprintf("InitialDeposit must be greater than or equal to %s", initDeposit))
}

func ErrInvalidProposer(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidProposer, message)
}

func ErrInvalidHeight(codespace string, h, ch, max uint64) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidHeight,
		fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.",
			h, ch, ch, max))
}

func ErrInsufficientCoins(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeInsufficientCoins, message)
}

func ErrUnknownRequest(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnknownRequest, message)
}

func ErrInvalidCoins(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, message)
}