package types

import (
	"errors"
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"
)

type AminoInterchainAccount struct {
	*auth.BaseAccount `protobuf:"bytes,1,opt,name=base_account,json=baseAccount,proto3,embedded=base_account" json:"base_account,omitempty" yaml:"base_account"`
	AccountOwner      string `protobuf:"bytes,2,opt,name=account_owner,json=accountOwner,proto3" json:"account_owner,omitempty" yaml:"account_owner"`
}

// AminoInterchainAccountI wraps the authtypes.AccountI interface
type AminoInterchainAccountI interface {
	auth.Account
}

// AminoInterchainAccountPretty defines an unexported struct used for encoding the AminoInterchainAccount details
type aminoInterchainAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	AccountOwner  string         `json:"account_owner" yaml:"account_owner"`
}

// NewAminoInterchainAccount creates and returns a new AminoInterchainAccount type
func NewAminoInterchainAccount(ba *auth.BaseAccount, accountOwner string) *AminoInterchainAccount {
	return &AminoInterchainAccount{
		BaseAccount:  ba,
		AccountOwner: accountOwner,
	}
}

// SetPubKey implements the authtypes.AccountI interface
func (ia AminoInterchainAccount) SetPubKey(pubKey crypto.PubKey) error {
	return sdkerrors.Wrap(ErrUnsupported, "cannot set public key for interchain account")
}

// SetSequence implements the authtypes.AccountI interface
func (ia AminoInterchainAccount) SetSequence(seq uint64) error {
	return sdkerrors.Wrap(ErrUnsupported, "cannot set sequence number for interchain account")
}

// Validate implements basic validation of the AminoInterchainAccount
func (ia AminoInterchainAccount) Validate() error {
	if strings.TrimSpace(ia.AccountOwner) == "" {
		return sdkerrors.Wrap(ErrInvalidAccountAddress, "AccountOwner cannot be empty")
	}
	if ia.PubKey != nil {
		return sdkerrors.Wrap(ErrInvalidPubKey, "pubkey must be nil")
	}
	if ia.Sequence != 0 {
		return sdkerrors.Wrap(ErrInvalidSequence, "sequence must be nil")
	}
	return ia.BaseAccount.Validate()
}

// String returns a string representation of the AminoInterchainAccount
func (ia AminoInterchainAccount) String() string {
	out, _ := ia.MarshalYAML()
	return string(out)
}

// MarshalYAML returns the YAML representation of the AminoInterchainAccount
func (ia AminoInterchainAccount) MarshalYAML() ([]byte, error) {
	accAddr := ia.Address

	bz, err := yaml.Marshal(aminoInterchainAccountPretty{
		Address:       accAddr,
		PubKey:        "",
		AccountNumber: ia.AccountNumber,
		Sequence:      ia.Sequence,
		AccountOwner:  ia.AccountOwner,
	})
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (m *AminoInterchainAccount) GetAddress() sdk.AccAddress {
	accAddr := m.BaseAccount.Address
	return accAddr
}

func (m *AminoInterchainAccount) SetAddress(address sdk.AccAddress) error {
	if len(m.BaseAccount.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	m.BaseAccount.Address = address
	return nil
}

func (m *AminoInterchainAccount) GetAccountNumber() uint64 {
	return m.BaseAccount.AccountNumber
}

func (m *AminoInterchainAccount) SetAccountNumber(u uint64) error {
	m.BaseAccount.AccountNumber = u
	return nil
}

func (m *AminoInterchainAccount) GetSequence() uint64 {
	return m.BaseAccount.Sequence
}

func (m *AminoInterchainAccount) Copy() sdk.Account {
	//TODO implement me
	cp := m.BaseAccount.Copy().(*auth.BaseAccount)
	return NewAminoInterchainAccount(cp, m.AccountOwner)
}
