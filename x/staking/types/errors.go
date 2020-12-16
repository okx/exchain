// nolint
package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = ModuleName

	CodeInvalidValidator  						uint32 = 67000
	CodeInvalidDelegation 						uint32 = 67001
	CodeInvalidInput      						uint32 = 67002

	CodeInvalidAddress          				uint32 = 67003
	CodeUnknownRequest           				uint32 = 67004
	CodeInvalidMinSelfDelegation 				uint32 = 67005
	CodeInvalidProxy             				uint32 = 67006
	CodeInvalidShareAdding       				uint32 = 67007
	CodeParseHTTPArgsWithLimit					uint32 = 67008
	CodeInvalidValidateBasic	 				uint32 = 67009
	CodeAddressNotEqual			 				uint32 = 67010
	CodeInternalError			 				uint32 = 67011
	CodeGetConsPubKeyBech32Failed 				uint32 = 67012
	CodeUnknownStakingQueryType		 			uint32 = 67013
	CodeDescriptionLengthBiggerThanLimit 		uint32 = 67014
	CodeAddSharesAsMinSelfDelegationFailed		uint32 = 67015
	CodeValidatorUpdateDescriptionFailed		uint32 = 67016
	CodeBondedPoolOrNotBondedIsNotExist			uint32 = 67017
)

// ErrNilValidatorAddr returns an error when an empty validator address appears
func ErrNilValidatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "validator address is invalid input")
}

// ErrBadValidatorAddr returns an error when an invalid validator address appears
func ErrBadValidatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, "validator address is invalid")
}

// ErrNoValidatorFound returns an error when a validator doesn't exist
func ErrNoValidatorFound(valAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidValidator, fmt.Sprintf("validator %s does not exist", valAddr))}
}

// ErrValidatorOwnerExists returns an error when the validator address has been registered
func ErrValidatorOwnerExists() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator,
		"validator already exist for this operator address, must use new validator operator address")
}

// ErrValidatorPubKeyExists returns an error when the validator consensus pubkey has been registered
func ErrValidatorPubKeyExists() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator,
		"validator already exist for this pubkey, must use new validator pubkey")
}

// ErrValidatorPubKeyTypeNotSupported returns an error when the type of  pubkey was not supported
func ErrValidatorPubKeyTypeNotSupported(keyType string,
	supportedTypes []string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator,
		fmt.Sprintf("validator pubkey type %s is not supported, must use %s", keyType, strings.Join(supportedTypes, ",")))
}

// ErrDescriptionLength returns an error when the description of validator has a wrong length
func ErrDescriptionLength(descriptor string, got, max int) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator,
		fmt.Sprintf("bad description length for %v, got length %v, max is %v", descriptor, got, max))
}

// ErrCommissionNegative returns an error when the commission is not positive
func ErrCommissionNegative() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission must be positive")
}

// ErrCommissionHuge returns an error when the commission is greater than 100%
func ErrCommissionHuge() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission cannot be more than 100%")
}

// ErrCommissionGTMaxRate returns an error when the commission rate is greater than the max rate
func ErrCommissionGTMaxRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission cannot be more than the max rate")
}

// ErrCommissionUpdateTime returns an error when the commission is remodified within 24 hours
func ErrCommissionUpdateTime() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission cannot be changed more than once in 24h")
}

// ErrCommissionChangeRateNegative returns an error when the commission change rate is not positive
func ErrCommissionChangeRateNegative() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission change rate must be positive")
}

// ErrCommissionChangeRateGTMaxRate returns an error when the commission change rate is greater than the max rate
func ErrCommissionChangeRateGTMaxRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission change rate cannot be more than the max rate")
}

// ErrCommissionGTMaxChangeRate returns an error when the new rate % points change is greater than the max change rate
func ErrCommissionGTMaxChangeRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "commission cannot be changed more than max change rate")
}

// ErrMinSelfDelegationInvalid returns an error when the msd isn't positive
func ErrMinSelfDelegationInvalid() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidValidator, "minimum self delegation must be a positive integer")
}

func ErrDescriptionIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "description must be included")
}

// ErrNilDelegatorAddr returns an error when the delegator address is nil
func ErrNilDelegatorAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "delegator address is nil")
}

// ErrWrongOperationAddr returns an error when the address is not expected
func ErrWrongOperationAddr() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidInput, "wrong operation addr")
}

// ErrBadDenom returns an error when the coin denomination is invalid
func ErrBadDenom() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "invalid coin denomination")}
}

// ErrBadDelegationAmount returns an error when the amount of delegation isn't positive
func ErrBadDelegationAmount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "amount must be > 0")
}

// ErrNoUnbondingDelegation returns an error when the unbonding delegation doesn't exist
func ErrNoUnbondingDelegation() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "no unbonding delegation found")
}

// ErrAddSharesToDismission returns an error when a zero-msd validator becomes the shares adding target
func ErrAddSharesToDismission(valAddr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidShareAdding,
		fmt.Sprintf("failed. destroyed validator %s isn't allowed to add shares to. please get rid of it from the "+
			"shares adding list by adding shares to other validators again or unbond all delegated tokens", valAddr))
}

// ErrNoAvailableValsToAddShares returns an error when none of the validators in shares adding list is available
func ErrNoAvailableValsToAddShares() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidShareAdding,
		"failed. there's no available validators among the shares adding list")
}

// ErrDelegatorNotAProxy returns an error when the target delegator to bind is not registered as a proxy
func ErrDelegatorNotAProxy(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. target address haven't reg as a proxy %s", delegator))}
}

// ErrNeverProxied returns an error when a delegator who's not a proxy trys to unreg
func ErrNeverProxied(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		fmt.Sprintf("failed. delegator %s has never registered as a proxy", delegator))}
}

// ErrAlreadyProxied returns an error when a proxy trys to reg the second time
func ErrAlreadyProxied(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		fmt.Sprintf("failed. delegator %s has already registered as a proxy", delegator))}
}

// ErrAddSharesDuringProxy returns an error when a delegator who has bound tries to add shares to validators by itself
func ErrAddSharesDuringProxy(delegator string, proxy string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. ban from adding shares to validators before unbinding proxy relationship between %s and %s", delegator, proxy))}
}

// ErrDoubleProxy returns an error when a delegator trys to bind more than one proxy
func ErrDoubleProxy(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. proxy isn't allowed to bind with other proxy %s", delegator))}
}

// ErrNotFoundProxy returns an error when the proxy doesn't exist
func ErrNotFoundProxy(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. no proxy with %s", delegator))}
}

// ErrInvalidDelegation returns an error when the delegation is invalid
func ErrInvalidDelegation(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. invalid delegation on %s", delegator))}
}

// ErrNilValidatorAddrs returns an error when there's no target validator address in MsgAddShares
func ErrNilValidatorAddrs() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidInput,
		"failed. validator addresses are nil")}
}

// ErrExceedValidatorAddrs returns an error when the number of target validators exceeds the max limit
func ErrExceedValidatorAddrs(num int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidInput,
		fmt.Sprintf("failed. the number of validator addresses is over the limit %d", num))}
}

// ErrNoDelegationToAddShares returns an error when there's no delegation to support adding shares to validators
func ErrNoDelegationToAddShares(delegator string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. there's no delegation of %s", delegator))}
}

// ErrNotInDelegating returns an error when the UndelegationInfo doesn't exist during it's unbonding period
func ErrNotInDelegating(addr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. the addr %s is not in the status of undelegating", addr))
}

// ErrInsufficientDelegation returns an error when the delegation left is not enough for unbonding
func ErrInsufficientDelegation(quantity, delLeft string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. insufficient delegation. [delegation left]:%s, [quantity to unbond]:%s", delLeft, quantity))
}

// ErrInsufficientQuantity returns an error when the quantity is less than the min delegation limit
func ErrInsufficientQuantity(quantity, minLimit string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. insufficient quantity. [min limit]:%s, [quantity]:%s", minLimit, quantity))
}

// ErrMoreMinSelfDelegation returns an error when the msd doesn't match the rest of shares on a validator
func ErrMoreMinSelfDelegation(valAddr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidMinSelfDelegation,
		fmt.Sprintf("failed. min self delegation of %s is more than its shares", valAddr))
}

// ErrNoMinSelfDelegation returns an error when the msd has already been unbonded
func ErrNoMinSelfDelegation(valAddr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidMinSelfDelegation,
		fmt.Sprintf("failed. there's no min self delegation on %s", valAddr))
}

// ErrBadUnDelegationAmount returns an error when the amount of delegation is not positive
func ErrBadUnDelegationAmount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		"failed. amount must be greater than 0")
}

// ErrInvalidProxyUpdating returns an error when the total delegated tokens on a proxy are going to be negative
func ErrInvalidProxyUpdating() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		"failed. the total delegated tokens on the proxy will be negative after this update")
}

// ErrInvalidProxyWithdrawTotal returns an error when proxy withdraws total tokens
func ErrInvalidProxyWithdrawTotal(addr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		fmt.Sprintf("failed. proxy %s has to unreg before withdrawing total tokens", addr))
}

// ErrAlreadyAddedShares returns an error when a delegator tries to bind proxy after adding shares
func ErrAlreadyAddedShares(delAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		fmt.Sprintf("failed. delegator %s isn't allowed to bind proxy while it has added shares. please unbond the delegation first", delAddr))}
}

// ErrNoDelegatorExisted returns an error when the info if a certain delegator doesn't exist
func ErrNoDelegatorExisted(delAddr string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		fmt.Sprintf("failed. delegator %s doesn't exist", delAddr))
}

// ErrTargetValsDuplicate returns an error when the target validators in voting list are duplicate
func ErrTargetValsDuplicate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidShareAdding,
		"failed. duplicate target validators")
}

// ErrAlreadyBound returns an error when a delegator keeps binding a proxy before proxy register
func ErrAlreadyBound(delAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidProxy,
		fmt.Sprintf("failed. %s has already bound a proxy. it's necessary to unbind before proxy register",
			delAddr))}
}

func ErrCodeInternalError() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternalError, "occur internal error")
}

func ErrGetConsPubKeyBech32() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetConsPubKeyBech32Failed, "get cons public key bech32 failed")
}

func ErrUnknownRequest() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "unknown request")
}

func ErrUnknownStakingQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownStakingQueryType, "unknown staking query type")
}

func ErrDescriptionLengthBiggerThanLimit() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDescriptionLengthBiggerThanLimit, "description's length is bigger than limit")
}

func ErrAddSharesAsMinSelfDelegationFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeAddSharesAsMinSelfDelegationFailed, fmt.Sprintf("transfer or add shares failed: ", message))
}

func ErrValidatorUpdateDescriptionFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorUpdateDescriptionFailed, "validator update description failed")
}

func ErrBondedPoolOrNotBondedIsNotExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBondedPoolOrNotBondedIsNotExist, "bonded pool or not bonded pool is empty")
}