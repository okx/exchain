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

	CodeInvalidAddress			 uint32 = BaseGovError
	CodeUnknownProposal          uint32 = BaseGovError+1
	CodeInvalidContent           uint32 = BaseGovError+6
	CodeInvalidProposalType      uint32 = BaseGovError+7
	CodeInvalidVote              uint32 = BaseGovError+8
	CodeInvalidGenesis           uint32 = BaseGovError+9
	CodeProposalHandlerNotExists uint32 = BaseGovError+10
	CodeInvalidProposalStatus   uint32 = BaseGovError+11
	CodeInitialDepositNotEnough uint32 = BaseGovError+12
	CodeInvalidProposer         uint32 = BaseGovError+13
	CodeInvalidHeight           uint32 = BaseGovError+14
	CodeInvalidCoins			uint32 = BaseGovError+17
	CodeUnknownParamType		uint32 = BaseGovError+18
)

func ErrInvalidAddress(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address %s", address))
}

func ErrInvalidProposalContent() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidContent, fmt.Sprintf("invalid proposal content"))
}

func ErrInvalidProposalType(proposalType string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProposalType, fmt.Sprintf("proposal type '%s' is not valid", proposalType))
}

func ErrInvalidVote(voteOption VoteOption) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption.String()))
}

func ErrInvalidGenesis() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidVote, "initial proposal ID hasn't been set")
}

func ErrNoProposalHandlerExists(content interface{}) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeProposalHandlerNotExists, fmt.Sprintf("'%T' does not have a corresponding handler", content))
}

func ErrUnknownProposal(proposalID uint64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownProposal, fmt.Sprintf("unknown proposal with id %d", proposalID))
}

func ErrInvalidateProposalStatus() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProposalStatus, "The status of proposal is can not be voted.")
}

func ErrInitialDepositNotEnough(initDeposit string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInitialDepositNotEnough,
		fmt.Sprintf("InitialDeposit must be greater than or equal to %s", initDeposit))
}

func ErrInvalidProposer() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProposer, "invalid proposer")
}

func ErrInvalidHeight(h, ch, max uint64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidHeight,
		fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.",
			h, ch, ch, max))
}

func ErrInvalidCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "invalide coins")
}

func ErrUnknownGovParamType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownParamType, "unkonwn gov param type")
}