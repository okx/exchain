package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	SpotCodespace sdk.CodespaceType = "spot"

	CodeEmptyProduct    CodeType = 62001
	CodeInvalidProduct  CodeType = 62002
	CodeInvalidWithdraw CodeType = 62003
)

// ErrEmptySymbol returns an error with an empty symbol of token
func ErrEmptyProduct(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProduct, "failed. product cannot be empty")
}

// ErrExistingProduct returns an error with an existing product
func ErrExistingProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProduct,
		fmt.Sprintf("failed. product %s already exists", product))
}

// ErrExistingProduct returns an error with the nonexistent product
func ErrNonexistentProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProduct,
		fmt.Sprintf("failed. product %s does not exist", product))
}

// ErrInvalidProduct returns an error with an invalid product
func ErrInvalidProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProduct,
		fmt.Sprintf("failed. product %s is invalid", product))
}

// ErrProductUnauthorizedIdentity returns an error with the unauthorized identity of the owner
func ErrProductUnauthorizedIdentity(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorizedIdentity,
		fmt.Sprintf("failed. not the owner of product: %s", product))
}

// ErrInvalidDepositsToken returns an error with the invalid deposits info
func ErrInvalidDepositsToken(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSymbol,
		fmt.Sprintf("failed. deposits only support %s token", symbol))
}

// ErrInvalidWithdrawAmount returns an error with the invalid withdraw info
func ErrInvalidWithdrawAmount(codespace sdk.CodespaceType, depositsAmount, withdrawAmount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidWithdraw,
		fmt.Sprintf("failed. deposits:%s is less than withdraw:%s", depositsAmount, withdrawAmount))
}

// ErrProductWithdraw returns an error when withdraw failed
func ErrProductWithdraw(codespace sdk.CodespaceType, deposits, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadJSONUnmarshaling, fmt.Sprintf("failed. withdraw deposits:%s error: %s", deposits, msg))
}

// ErrSaveProduct returns an error when save product to db failed
func ErrSaveProduct(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadJSONUnmarshaling, fmt.Sprintf("failed. save product error: %s", msg))
}
