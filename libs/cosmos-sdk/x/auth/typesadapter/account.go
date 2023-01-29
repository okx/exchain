package types

import (
	"bytes"
	"errors"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	cryptotypes "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// NewBaseAccountWithAddress - returns a new base account with a given address
// leaving AccountNumber and Sequence to zero.
func NewBaseAccountWithAddress(addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		Address: addr.String(),
	}
}

// NewBaseAccount creates a new BaseAccount object
//nolint:interfacer
func NewBaseAccount(address sdk.AccAddress, pubKey cryptotypes.PubKey, accountNumber, sequence uint64) *BaseAccount {
	acc := &BaseAccount{
		Address:       address.String(),
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}

	err := acc.SetPubKey(pubKey)
	if err != nil {
		panic(err)
	}

	return acc
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &acc)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.Address == "" || acc.PubKey == nil {
		return nil
	}

	accAddr, err := sdk.AccAddressFromBech32(acc.Address)
	if err != nil {
		return err
	}

	if !bytes.Equal(acc.GetPubKey().Address().Bytes(), accAddr.Bytes()) {
		return errors.New("account address and pubkey address do not match")
	}

	return nil
}

// GetPubKey - Implements sdk.AccountI.
func (acc BaseAccount) GetPubKey() (pk cryptotypes.PubKey) {
	if acc.PubKey == nil {
		return nil
	}
	content, ok := acc.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil
	}
	return content
}

// SetPubKey - Implements sdk.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey cryptotypes.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
		return nil
	}
	any, err := codectypes.NewAnyWithValue(pubKey)
	if err == nil {
		acc.PubKey = any
	}
	return err
}
