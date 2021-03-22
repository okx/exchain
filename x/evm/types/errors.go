package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcmn "github.com/ethereum/go-ethereum/common"
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

	CodeSpaceEvmCallFailed = uint32(7)

	ErrorHexData = "HexData"
)

// ErrDeployerAlreadyExists returns an error when duplicated deployer address will be added
func ErrDeployerAlreadyExists(distributorAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			13,
			fmt.Sprintf("failed. deployer %s is already in the whitelist", distributorAddr.String()))}
}

// ErrDeployerNotExists returns an error when a deployer address not in the whitelist will be deleted
func ErrDeployerNotExists(distributorAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			14,
			fmt.Sprintf("failed. deployer %s is not in the whitelist", distributorAddr.String()))}
}

// ErrDeployerUnqualified returns an error when a deployer not in the whitelist tries to create a contract
func ErrDeployerUnqualified(distributorAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			15,
			fmt.Sprintf("failed. unqualified deployer %s for a contract deployment", distributorAddr.String()))}
}

// ErrContractAlreadyExists returns an error when duplicated contract will be added into blocked list
func ErrContractAlreadyExists(contractAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			16,
			fmt.Sprintf("failed. contract %s is already in the blocked list", contractAddr.String()))}
}

// ErrContractNotExists returns an error when a contract not in the blocked list will be deleted
func ErrContractNotExists(contractAddr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			17,
			fmt.Sprintf("failed. contract %s is not in the blocked list", contractAddr.String()))}
}

// ErrCallBlockedContract returns an error when the blocked contract is invoked
func ErrCallBlockedContract(contractAddr ethcmn.Address) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			18,
			fmt.Sprintf("failed. the contract %s is not allowed to invoke", contractAddr.Hex()),
		),
	}
}
