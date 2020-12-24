package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const CodeType
const (
	CodeAddrAndProductAllRequired   uint32 = 64000
	codeInvalidTokenPair            uint32 = 64001
	codeTokenPairNotFound           uint32 = 64002
	codeBalanceNotEnough            uint32 = 64003
	codeInvalidAsset                uint32 = 64004
	codeUnknownOperator             uint32 = 64005
	codeExistOperator               uint32 = 64006
	codeInvalidWebsiteLength        uint32 = 64007
	codeInvalidWebsiteURL           uint32 = 64008
	CodeTokenPairSaveFailed         uint32 = 64009
	CodeInsufficientFeeCoins        uint32 = 64010
	CodeTokenPairAlreadyExist       uint32 = 64011
	CodeMustTokenPairOwner          uint32 = 64012
	CodeDepositOnlySupportDenom     uint32 = 64013
	CodeInsufficientDepositCoins    uint32 = 64014
	CodeWithdrawOnlySupportDenom    uint32 = 64015
	CodeInsufficientWithdrawCoins   uint32 = 64016
	CodeInvalidAddress              uint32 = 64017
	CodeUnauthorized                uint32 = 64018
	CodeInvalidCoins                uint32 = 64019
	CodeRepeatedTransferOwner       uint32 = 64020
	CodeDepositFailed               uint32 = 64021
	CodeWithdrawFailed              uint32 = 64022
	CodeGetConfirmOwnershipNotExist uint32 = 64023
	CodeTokenPairIsRequired         uint32 = 64024
	CodeIsTokenPairLocked           uint32 = 64025
	CodeDexUnknownMsgType           uint32 = 64026
	CodeDexUnknownQueryType         uint32 = 64027
	CodeInitPriceIsNotPositive      uint32 = 64028
	CodeAddressIsRequired           uint32 = 64029
	CodeTokenOfPairNotExist         uint32 = 64030
	CodeIsTransferringOwner         uint32 = 64031
	CodeTransferOwnerExpired        uint32 = 64032
	CodeUnauthorizedOperator        uint32 = 64033
)

// Addr and Product All Required
func ErrAddrAndProductAllRequired() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeAddrAndProductAllRequired, "bad request: address、base_asset and quote_asset could not be "+" at the same time")}
}

// ErrInvalidTokenPair returns invalid product error
func ErrInvalidTokenPair(tokenPair string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeInvalidTokenPair, fmt.Sprintf("invalid tokenpair: %s", tokenPair))}
}

// ErrTokenPairNotFound returns token pair not found error
func ErrTokenPairNotFound(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeTokenPairNotFound, fmt.Sprintf("tokenpair not found: %s", product))}
}

// ErrInvalidBalanceNotEnough returns invalid balance not enough error
func ErrBalanceNotEnough(proposer string, initialDeposit string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeBalanceNotEnough, fmt.Sprintf("failed to submit proposal because proposer %s didn't have enough coins to pay for the initial deposit %s", proposer, initialDeposit))}
}

// ErrInvalidAsset returns invalid asset error
func ErrInvalidAsset(localMinDeposit string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeInvalidAsset, fmt.Sprintf("failed to submit proposal because initial deposit should be more than %s", localMinDeposit))}
}

func ErrUnknownOperator(addr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeUnknownOperator, fmt.Sprintf("unknown dex operator with address %s", addr.String()))}
}

func ErrExistOperator(addr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeExistOperator, fmt.Sprintf("dex operator already exists with address %s", addr.String()))}
}

func ErrInvalidWebsiteLength(got, max int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeInvalidWebsiteLength, fmt.Sprintf("invalid website length, got length %v, max is %v", got, max))}
}

func ErrInvalidWebsiteURL(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, codeInvalidWebsiteURL, fmt.Sprintf("invalid website URL: %s", msg))}
}

func ErrTokenPairSaveFailed(err string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenPairSaveFailed, fmt.Sprintf("failed to SaveTokenPair: %s", err))}
}

func ErrInsufficientFeeCoins(fee sdk.Coins) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInsufficientFeeCoins, fmt.Sprintf("insufficient fee coins(need %s)", fee.String()))}
}

// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset string, quoteAsset string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenPairAlreadyExist, fmt.Sprintf("the token pair exists with %s and %s", baseAsset, quoteAsset))}
}

func ErrMustTokenPairOwner(addr string, tokenPair string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMustTokenPairOwner, fmt.Sprintf("%s is not the owner of product: %s", addr, tokenPair))}
}

func ErrDepositOnlySupportDenom(denom string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeDepositOnlySupportDenom, fmt.Sprintf("deposits only support %s token", denom))}
}

func ErrInsufficientDepositCoins(depositCoins string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInsufficientDepositCoins, fmt.Sprintf("insufficient deposit coins(need %s)", depositCoins))}
}

func ErrWithdrawOnlySupportDenom(denom string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeWithdrawOnlySupportDenom, fmt.Sprintf("failed to withdraws because deposits only support %s token", denom))}
}

func ErrInsufficientWithdrawCoins(depositCoins string, amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInsufficientWithdrawCoins, fmt.Sprintf("failed to withdraws because deposits:%s is less than withdraw:%s", depositCoins, amount))}
}

func ErrInvalidAddress(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("there is no withdrawing for address: %s", msg))}
}

func ErrUnauthorized(address, product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeUnauthorized, fmt.Sprintf("%s is not the owner of product(%s)", address, product))}
}

func ErrInvalidCoins() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "invalid coins")}
}

func ErrRepeatedTransferOwner(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeRepeatedTransferOwner, fmt.Sprintf("repeated transfer-ownership of product(%s)", product))}
}

func ErrDepositFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeDepositFailed, fmt.Sprintf("deposit occur error: %s", msg))}
}

func ErrWithdrawFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeWithdrawFailed, fmt.Sprintf("withdraw occur error: %s", msg))}
}

func ErrGetConfirmOwnershipNotExist(address string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeGetConfirmOwnershipNotExist, fmt.Sprintf("no transfer-ownership of list (%s) to confirm", address))}
}

func ErrTokenPairIsRequired() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenPairIsRequired, "token pair is required")}
}

func ErrIsTokenPairLocked(tokenPairName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsTokenPairLocked, fmt.Sprintf("unexpected state, the trading pair (%s) is locked", tokenPairName))}
}

func ErrDexUnknownMsgType(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeDexUnknownMsgType, fmt.Sprintf("unrecognized dex message type: %T", msg))}
}

func ErrDexUnknownQueryType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeDexUnknownQueryType, "unknown dex query endpoint")}
}

func ErrInitPriceIsNotPositive() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInitPriceIsNotPositive, "invalid init price number")}
}

func ErrAddressIsRequired(addrType string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeAddressIsRequired, fmt.Sprintf("missing %s address", addrType))}
}

func ErrTokenOfPairNotExist(baseAsset string, quoteAsset string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenOfPairNotExist, fmt.Sprintf("the token of pair is not exists with %s and %s", baseAsset, quoteAsset))}
}

func ErrIsTransferringOwner(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsTransferringOwner, fmt.Sprintf("the product(%s) is transferring ownership, not allowed to be deposited", product))}
}

func ErrIsTransferOwnerExpired(time string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTransferOwnerExpired, fmt.Sprintf("transfer-ownership is expired, expire time (%s)", time))}
}

func ErrUnauthorizedOperator(operator, owner string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeUnauthorizedOperator, fmt.Sprintf("%s is not the owner of operator(%s)", owner, operator))}
}
