package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/tendermint/go-amino"

	"gopkg.in/yaml.v2"

	"github.com/okex/exchain/libs/tendermint/crypto"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
)

var (
	_ authexported.GenesisAccount = (*ModuleAccount)(nil)
	_ exported.ModuleAccountI     = (*ModuleAccount)(nil)
)

func init() {
	// Register the ModuleAccount type as a GenesisAccount so that when no
	// concrete GenesisAccount types exist and **default** genesis state is used,
	// the genesis state will serialize correctly.
	authtypes.RegisterAccountTypeCodec(&ModuleAccount{}, "cosmos-sdk/ModuleAccount")
}

// ModuleAccount defines an account for modules that holds coins on a pool
type ModuleAccount struct {
	*authtypes.BaseAccount

	Name        string   `json:"name" yaml:"name"`               // name of the module
	Permissions []string `json:"permissions" yaml:"permissions"` // permissions of module account
}

func (acc *ModuleAccount) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
	var dataLen uint64 = 0

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		if pbType != amino.Typ3_ByteLength {
			return fmt.Errorf("invalid type byte: %v", pbType)
		}
		data = data[1:]

		var n int
		dataLen, n, err = amino.DecodeUvarint(data)
		if err != nil {
			return err
		}

		data = data[n:]
		if len(data) < int(dataLen) {
			return fmt.Errorf("invalid data length: %v", dataLen)
		}
		subData := data[:dataLen]

		switch pos {
		case 1:
			base := new(authtypes.BaseAccount)
			err = base.UnmarshalFromAmino(cdc, subData)
			if err != nil {
				return err
			}
			acc.BaseAccount = base
		case 2:
			acc.Name = string(subData)
		case 3:
			acc.Permissions = append(acc.Permissions, string(subData))
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (acc ModuleAccount) Copy() interface{} {
	return NewModuleAccount(authtypes.NewBaseAccount(acc.Address, acc.Coins, acc.PubKey, acc.AccountNumber, acc.Sequence), acc.Name, acc.Permissions...)
}

func (acc ModuleAccount) AminoSize(cdc *amino.Codec) int {
	size := 0
	if acc.BaseAccount != nil {
		baccSize := acc.BaseAccount.AminoSize(cdc)
		size += 1 + amino.UvarintSize(uint64(baccSize)) + baccSize
	}
	if acc.Name != "" {
		size += 1 + amino.EncodedStringSize(acc.Name)
	}
	for _, p := range acc.Permissions {
		size += 1 + amino.EncodedStringSize(p)
	}
	return size
}

func (acc ModuleAccount) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(acc.AminoSize(cdc))
	err := acc.MarshalAminoTo(cdc, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (acc ModuleAccount) MarshalAminoTo(cdc *amino.Codec, buf *bytes.Buffer) error {
	// field 1
	if acc.BaseAccount != nil {
		const pbKey = 1<<3 | 2
		buf.WriteByte(pbKey)
		baccSize := acc.BaseAccount.AminoSize(cdc)
		err := amino.EncodeUvarintToBuffer(buf, uint64(baccSize))
		if err != nil {
			return err
		}
		lenBeforeData := buf.Len()
		err = acc.BaseAccount.MarshalAminoTo(cdc, buf)
		if err != nil {
			return err
		}
		if buf.Len()-lenBeforeData != baccSize {
			return amino.NewSizerError(baccSize, buf.Len()-lenBeforeData, baccSize)
		}
	}

	// field 2
	if acc.Name != "" {
		const pbKey = 2<<3 | 2
		err := amino.EncodeStringWithKeyToBuffer(buf, acc.Name, pbKey)
		if err != nil {
			return err
		}
	}

	// field 3
	for _, perm := range acc.Permissions {
		const pbKey = 3<<3 | 2
		err := amino.EncodeStringWithKeyToBuffer(buf, perm, pbKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewModuleAddress creates an AccAddress from the hash of the module's name
func NewModuleAddress(name string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(name)))
}

// NewEmptyModuleAccount creates a empty ModuleAccount from a string
func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := authtypes.NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: &baseAcc,
		Name:        name,
		Permissions: permissions,
	}
}

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba *authtypes.BaseAccount,
	name string, permissions ...string) *ModuleAccount {

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: ba,
		Name:        name,
		Permissions: permissions,
	}
}

// HasPermission returns whether or not the module account has permission.
func (ma ModuleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma ModuleAccount) GetName() string {
	return ma.Name
}

// GetPermissions returns permissions granted to the module account
func (ma ModuleAccount) GetPermissions() []string {
	return ma.Permissions
}

// SetPubKey - Implements Account
func (ma ModuleAccount) SetPubKey(pubKey crypto.PubKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// SetSequence - Implements Account
func (ma ModuleAccount) SetSequence(seq uint64) error {
	return fmt.Errorf("not supported for module accounts")
}

// Validate checks for errors on the account fields
func (ma ModuleAccount) Validate() error {
	if strings.TrimSpace(ma.Name) == "" {
		return errors.New("module account name cannot be blank")
	}
	if !ma.Address.Equals(sdk.AccAddress(crypto.AddressHash([]byte(ma.Name)))) {
		return fmt.Errorf("address %s cannot be derived from the module name '%s'", ma.Address, ma.Name)
	}

	return ma.BaseAccount.Validate()
}

type moduleAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	Name          string         `json:"name" yaml:"name"`
	Permissions   []string       `json:"permissions" yaml:"permissions"`
}

func (ma ModuleAccount) String() string {
	out, _ := ma.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (ma ModuleAccount) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(moduleAccountPretty{
		Address:       ma.Address,
		Coins:         ma.Coins,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Sequence:      ma.Sequence,
		Name:          ma.Name,
		Permissions:   ma.Permissions,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// MarshalJSON returns the JSON representation of a ModuleAccount.
func (ma ModuleAccount) MarshalJSON() ([]byte, error) {
	return codec.Cdc.MarshalJSON(moduleAccountPretty{
		Address:       ma.Address,
		Coins:         ma.Coins,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Sequence:      ma.Sequence,
		Name:          ma.Name,
		Permissions:   ma.Permissions,
	})
}

// UnmarshalJSON unmarshals raw JSON bytes into a ModuleAccount.
func (ma *ModuleAccount) UnmarshalJSON(bz []byte) error {
	var alias moduleAccountPretty
	if err := codec.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	ma.BaseAccount = authtypes.NewBaseAccount(alias.Address, alias.Coins, nil, alias.AccountNumber, alias.Sequence)
	ma.Name = alias.Name
	ma.Permissions = alias.Permissions

	return nil
}
