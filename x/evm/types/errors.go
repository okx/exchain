package types

import (
	"errors"
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

// NOTE: We can't use 1 since that error code is reserved for internal errors.
const (
	DefaultCodespace string = ModuleName
)

var (
	// ErrInvalidState returns an error resulting from an invalid Storage State.
	ErrInvalidState = sdkerrors.Register(ModuleName, 2, "invalid storage state")

	// ErrChainConfigNotFound returns an error if the chain config cannot be found on the store.
	ErrChainConfigNotFound = sdkerrors.Register(ModuleName, 3, "chain configuration not found")

	// ErrInvalidChainConfig returns an error resulting from an invalid ChainConfig.
	ErrInvalidChainConfig = sdkerrors.Register(ModuleName, 4, "invalid chain configuration")

	// ErrCreateDisabled returns an error if the EnableCreate parameter is false.
	ErrCreateDisabled = sdkerrors.Register(ModuleName, 5, "EVM Create operation is disabled")

	// ErrCallDisabled returns an error if the EnableCall parameter is false.
	ErrCallDisabled = sdkerrors.Register(ModuleName, 6, "EVM Call operation is disabled")

	// ErrKeyNotFound returns an error if the target key not found in database.
	ErrKeyNotFound = sdkerrors.Register(ModuleName, 8, "Key not found in database")

	// ErrStrConvertFailed returns an error if failed to convert string
	ErrStrConvertFailed = sdkerrors.Register(ModuleName, 9, "Failed to convert string")

	// ErrUnexpectedProposalType returns an error when the proposal type is not supported in evm module
	ErrUnexpectedProposalType = sdkerrors.Register(ModuleName, 10, "Unsupported proposal type of evm module")

	// ErrEmptyAddressList returns an error if the address list is empty
	ErrEmptyAddressList = sdkerrors.Register(ModuleName, 11, "Empty account address list")

	// ErrDuplicatedAddr returns an error if the address is duplicated in address list
	ErrDuplicatedAddr = sdkerrors.Register(ModuleName, 12, "Duplicated address in address list")

	// ErrDuplicatedAddr returns an error if the address is duplicated in address list
	ErrOperation = sdkerrors.Register(ModuleName, 16, "Special contract method blocked operation can not change blocked contract list")

	// ErrDuplicatedMethod returns an error if the contract method is duplicated
	ErrDuplicatedMethod = sdkerrors.Register(ModuleName, 17, "Duplicated contract method in address list")

	// ErrEmptyMethod returns an error if the contract method is empty
	ErrEmptyMethod = sdkerrors.Register(ModuleName, 18, "Empty contract method blocked is not allowed")

	// ErrEmptyAddressBlockedContract returns an error if the contract method is empty
	ErrEmptyAddressBlockedContract = sdkerrors.Register(ModuleName, 19, "Empty address in contract method blocked list is not allowed")

	CodeSpaceEvmCallFailed = uint32(7)

	ErrorHexData = "HexData"

	ErrorContractMethodBlockedIsNotExist = errors.New("it's not exist in contract method blocked list")
)

const (
	CodeGetChainHeightFailed   uint32 = 62004
	CodeGetBlockTxHashesFailed uint32 = 62005
)

// ErrOversizeAddrList returns an error when the length of address list in the proposal is larger than the max limitation
func ErrOversizeAddrList(length int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			13,
			fmt.Sprintf("failed. the length of address list in the proposal %d is larger than the max limitation %d",
				length, maxAddressListLength,
			))}
}

// ErrUnauthorizedAccount returns an error when an account not in the whitelist tries to create a contract
func ErrUnauthorizedAccount(distributorAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			14,
			fmt.Sprintf("failed. the account %s is not allowed to deploy a contract", ethcmn.BytesToAddress(distributorAddr)))}
}

// ErrCallBlockedContract returns an error when the blocked contract is invoked
func ErrCallBlockedContract(descriptor string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			15,
			descriptor,
		),
	}
}

// ErrBlockedContractMethodIsNotExist returns an error when the blocked contract method is not exist
func ErrBlockedContractMethodIsNotExist(address sdk.Address, err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			20,
			fmt.Sprintf("Delete contract(%s) method failed: %s", address, err.Error()),
		),
	}
}

type ErrContractBlockedVerify struct {
	Descriptor string
}

func ErrRegisteredContract(contract string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			21,
			fmt.Sprintf("the contract is already registered: %s", contract),
		),
	}
}
