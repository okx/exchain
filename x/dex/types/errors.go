package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const CodeType
const (
	CodeAddrAndProductAllRequired           uint32 = 64000
	codeInvalidTokenPair                    uint32 = 64001
	codeTokenPairNotFound                   uint32 = 64002
	codeDelistOwnerNotMatch                 uint32 = 64003
	codeInvalidBalanceNotEnough             uint32 = 64004
	codeInvalidAsset                        uint32 = 64005
	codeUnknownOperator                     uint32 = 64006
	codeExistOperator                       uint32 = 64007
	codeInvalidWebsiteLength                uint32 = 64008
	codeInvalidWebsiteURL                   uint32 = 64009
	CodeTokenPairIsInvalid                  uint32 = 64010
	CodeTokenPairSaveFailed                 uint32 = 64011
	CodeInsufficientFeeCoins                uint32 = 64012
	CodeTokenPairAlreadyExist               uint32 = 64013
	CodeMustTokenPairOwner                  uint32 = 64014
	CodeDepositOnlySupportDefaultBondDenom  uint32 = 64015
	CodeInsufficientDepositCoins            uint32 = 64016
	CodeWithdrawOnlySupportDefaultBondDenom uint32 = 64017
	CodeInsufficientWithdrawCoins           uint32 = 64018
	CodeWithdrawDepositsError               uint32 = 64019
	CodeMustOperatorOwner                   uint32 = 64020
	CodeInvalidProduct                      uint32 = 64021
	CodeInvalidAddress                      uint32 = 64022
	CodeUnknownRequest                      uint32 = 64023
	CodeInternal                            uint32 = 64024
	CodeUnauthorized                        uint32 = 64025
	CodeInvalidCoins                        uint32 = 64027
	CodeBlockTimeAfterFailed                uint32 = 64028
	CodeDepositFailed                       uint32 = 64029
	CodeWithdrawFailed                      uint32 = 64030
	CodeGetConfirmOwnershipNotExist         uint32 = 64031
	CodeIsNil                               uint32 = 64032
	CodeIsTokenPairLocked                   uint32 = 64033
	CodeUnknownMsgType                      uint32 = 64034
	CodeUnknownQueryType                    uint32 = 64035

	CodeInitPriceIsNotPositive uint32 = 64037
	CodeProductIsEmpty         uint32 = 64038
)

// CodeType to Message
func codeToDefaultMsg(code uint32) string {
	switch code {
	case codeInvalidTokenPair:
		return "invalid tokenpair"
	case codeTokenPairNotFound:
		return "tokenpair not found"
	case codeDelistOwnerNotMatch:
		return "tokenpair delistor should be it's owner "
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}

// Addr and Product All Required
func ErrAddrAndProductAllRequired() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeAddrAndProductAllRequired, "bad request: address、base_asset and quote_asset could not be empty at the same time")
}

// invalid tokenpair
func ErrTokenPairIsInvalid() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairIsInvalid, "the nil pointer is not expected")
}

// ErrTokenPairNotFound returns token pair not found error
func ErrTokenPairNotFound() sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeTokenPairNotFound, codeToDefaultMsg(codeTokenPairNotFound))
}

// ErrDelistOwnerNotMatch returns delist owner not match error
func ErrDelistOwnerNotMatch(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeDelistOwnerNotMatch, fmt.Sprintf(codeToDefaultMsg(codeDelistOwnerNotMatch)+": %s", msg))
}

// ErrInvalidBalanceNotEnough returns invalid balance not enough error
func ErrInvalidBalanceNotEnough(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeInvalidBalanceNotEnough, message)
}

// ErrInvalidAsset returns invalid asset error
func ErrInvalidAsset(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeInvalidAsset, message)
}

func ErrUnknownOperator(addr sdk.AccAddress) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeUnknownOperator, fmt.Sprintf("unknown dex operator with address %s", addr.String()))
}

func ErrExistOperator(addr sdk.AccAddress) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeExistOperator, fmt.Sprintf("dex operator already exists with address %s", addr.String()))
}

func ErrInvalidWebsiteLength(got, max int) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeInvalidWebsiteLength, fmt.Sprintf("invalid website length, got length %v, max is %v", got, max))
}

func ErrInvalidWebsiteURL(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeInvalidWebsiteURL, fmt.Sprintf("invalid website URL: %s", msg))
}

// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset string, quoteAsset string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairAlreadyExist, fmt.Sprintf("failed. the token pair exists with %s and %s", baseAsset, quoteAsset))
}

// ErrInvalidTokenPair returns invalid product error
func ErrInvalidTokenPair(tokenPair string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, codeInvalidTokenPair, fmt.Sprintf(codeToDefaultMsg(codeInvalidTokenPair)+": %s", tokenPair))
}
func ErrTokenPairSaveFailed(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairSaveFailed, fmt.Sprintf("failed to SaveTokenPair: %s", msg))
}
func ErrInsufficientFeeCoins(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientFeeCoins, fmt.Sprintf("insufficient fee coins(need %s)", msg))
}
func ErrMustTokenPairOwner(addr string, tokenPair string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMustTokenPairOwner, fmt.Sprintf("failed because %s is not the owner of product:%s", addr, tokenPair))
}
func ErrDepositOnlySupportDefaultBondDenom(defaultBondDenom string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDepositOnlySupportDefaultBondDenom, fmt.Sprintf("failed to deposit because deposits only support %s token", defaultBondDenom))
}
func ErrInsufficientDepositCoins(msg string, depositCoins string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientDepositCoins, fmt.Sprintf("failed: %s, because insufficient deposit coins(need %s)", msg, depositCoins))
}
func ErrWithdrawOnlySupportDefaultBondDenom(defaultBondDenom string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawOnlySupportDefaultBondDenom, fmt.Sprintf("failed to withdraws because deposits only support %s token", defaultBondDenom))
}
func ErrInsufficientWithdrawCoins(depositCoins string, amount string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientWithdrawCoins, fmt.Sprintf("failed to withdraws because deposits:%s is less than withdraw:%s", depositCoins, amount))
}
func ErrWithdrawDepositsError(depositCoins string, msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawDepositsError, fmt.Sprintf("withdraw deposits:%s error:%s", depositCoins, msg))
}
func ErrMustOperatorOwnerOwner(addr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMustOperatorOwner, fmt.Sprintf("failed because %s is not the owner of operator", addr))
}
func ErrInvalidProduct() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProduct, "invalid product")
}

func ErrInvalidAddress(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("there is no withdrawing for address %s", message))
}

func ErrUnknownRequest(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, fmt.Sprintf("unknown request: %s", message))
}

func ErrInternal(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternal, fmt.Sprintf("internal error: %s", message))
}

func ErrUnauthorized(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnauthorized, fmt.Sprintf("%s is not the owner of product", address))
}

func ErrInvalidCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "invalid coins")
}

func ErrBlockTimeAfterFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBlockTimeAfterFailed, fmt.Sprintf("block time after failed: %s", message))
}

func ErrDepositFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDepositFailed, fmt.Sprintf("deposit occur error: %s", message))
}

func ErrWithdrawFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeWithdrawFailed, fmt.Sprintf("withdraw occur error: %s", message))
}

func ErrGetConfirmOwnershipNotExist(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetConfirmOwnershipNotExist, fmt.Sprintf("no transfer-ownership of list (%s) to confirm", address))
}

func ErrIsNil() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsNil, "param is nil")
}

func ErrIsTokenPairLocked(tokenPairName string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsTokenPairLocked, fmt.Sprintf("unexpected state, the trading pair (%s) is locked", tokenPairName))
}

func ErrUnknownMsgType(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownMsgType, fmt.Sprintf("unrecognized dex message type: %T", msg))
}

func ErrUnknownQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownQueryType, "unknown dex query endpoint")
}

func ErrInitPriceIsNotPositive() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInitPriceIsNotPositive, "invalid init price number")
}

func ErrProductIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeProductIsEmpty, "product param is empty")
}
