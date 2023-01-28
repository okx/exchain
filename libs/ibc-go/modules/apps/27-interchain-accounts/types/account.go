package types

import (
	"regexp"

	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

// DefaultMaxAddrLength defines the default maximum character length used in validation of addresses
var DefaultMaxAddrLength = 128

// isValidAddr defines a regular expression to check if the provided string consists of
// strictly alphanumeric characters and is non empty.
var isValidAddr = regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString

// InterchainAccountI wraps the authtypes.AccountI interface
type InterchainAccountI interface {
	authtypes.Account
}

// interchainAccountPretty defines an unexported struct used for encoding the InterchainAccount details
type interchainAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	AccountOwner  string         `json:"account_owner" yaml:"account_owner"`
}

// GenerateAddress returns an sdk.AccAddress derived using the provided module account address and connection and port identifiers.
// The sdk.AccAddress returned is a sub-address of the module account, using the host chain connection ID and controller chain's port ID as the derivation key
func GenerateAddress(moduleAccAddr sdk.AccAddress, connectionID, portID string) sdk.AccAddress {
	return sdk.AccAddress(sdk.Derive(moduleAccAddr, []byte(connectionID+portID)))
}

// ValidateAccountAddress performs basic validation of interchain account addresses, enforcing constraints
// on address length and character set
func ValidateAccountAddress(addr string) error {
	if !isValidAddr(addr) || len(addr) > DefaultMaxAddrLength {
		return sdkerrors.Wrapf(
			ErrInvalidAccountAddress,
			"address must contain strictly alphanumeric characters, not exceeding %d characters in length",
			DefaultMaxAddrLength,
		)
	}

	return nil
}
