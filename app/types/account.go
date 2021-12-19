package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ exported.Account = (*EthAccount)(nil)
var _ exported.GenesisAccount = (*EthAccount)(nil)

func init() {
	authtypes.RegisterAccountTypeCodec(&EthAccount{}, EthAccountName)
}

// ----------------------------------------------------------------------------
// Main OKExChain account
// ----------------------------------------------------------------------------

// EthAccount implements the auth.Account interface and embeds an
// auth.BaseAccount type. It is compatible with the auth.AccountKeeper.
type EthAccount struct {
	*authtypes.BaseAccount `json:"base_account" yaml:"base_account"`
	CodeHash               []byte `json:"code_hash" yaml:"code_hash"`
	StateRoot ethcmn.Hash  `json:"state_root" yaml:"state_root"`	// merkle root of the storage trie
}

// ProtoAccount defines the prototype function for BaseAccount used for an
// AccountKeeper.
func ProtoAccount() exported.Account {
	return &EthAccount{
		BaseAccount: &auth.BaseAccount{},
		CodeHash:    ethcrypto.Keccak256(nil),
		StateRoot: ethcmn.Hash{},
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

// TODO: remove on SDK v0.40

// Balance returns the balance of an account.
func (acc EthAccount) Balance(denom string) sdk.Dec {
	return acc.GetCoins().AmountOf(denom)
}

// SetBalance sets an account's balance of the given coin denomination.
//
// CONTRACT: assumes the denomination is valid.
func (acc *EthAccount) SetBalance(denom string, amt sdk.Dec) {
	coins := acc.GetCoins()
	diff := amt.Sub(coins.AmountOf(denom))
	switch {
	case diff.IsPositive():
		// Increase coins to amount
		coins = coins.Add(sdk.NewCoin(denom, diff))
	case diff.IsNegative():
		// Decrease coins to amount
		coins = coins.Sub(sdk.NewCoins(sdk.NewCoin(denom, diff.Neg())))
	default:
		return
	}

	if err := acc.SetCoins(coins); err != nil {
		panic(fmt.Errorf("could not set %s coins for address %s: %w", denom, acc.EthAddress().String(), err))
	}
}

type ethermintAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	EthAddress    string         `json:"eth_address" yaml:"eth_address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	CodeHash      string         `json:"code_hash" yaml:"code_hash"`
	StorageRoot string `json:"storage_root" yaml:"storage_root"`
}

// MarshalYAML returns the YAML representation of an account.
func (acc EthAccount) MarshalYAML() (interface{}, error) {
	alias := ethermintAccountPretty{
		Address:       acc.Address,
		EthAddress:    acc.EthAddress().String(),
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
		StorageRoot: acc.StateRoot.String(),
	}

	var err error

	if acc.PubKey != nil {
		alias.PubKey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of an EthAccount.
func (acc EthAccount) MarshalJSON() ([]byte, error) {
	var ethAddress = ""

	if acc.BaseAccount != nil && acc.Address != nil {
		ethAddress = acc.EthAddress().String()
	}

	alias := ethermintAccountPretty{
		Address:       acc.Address,
		EthAddress:    ethAddress,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
		StorageRoot: acc.StateRoot.String(),
	}

	var err error

	if acc.PubKey != nil {
		alias.PubKey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into an EthAccount.
func (acc *EthAccount) UnmarshalJSON(bz []byte) error {
	var (
		alias ethermintAccountPretty
		err   error
	)

	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	switch {
	case !alias.Address.Empty() && alias.EthAddress != "":
		// Both addresses provided. Verify correctness
		ethAddress := ethcmn.HexToAddress(alias.EthAddress)
		ethAddressFromAccAddress := ethcmn.BytesToAddress(alias.Address.Bytes())

		if !bytes.Equal(ethAddress.Bytes(), alias.Address.Bytes()) {
			err = sdkerrors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"expected %s, got %s",
				ethAddressFromAccAddress.String(), ethAddress.String(),
			)
		}

	case !alias.Address.Empty() && alias.EthAddress == "":
		// unmarshal sdk.AccAddress only. Do nothing here
	case alias.Address.Empty() && alias.EthAddress != "":
		// retrieve sdk.AccAddress from ethereum address
		ethAddress := ethcmn.HexToAddress(alias.EthAddress)
		alias.Address = sdk.AccAddress(ethAddress.Bytes())
	case alias.Address.Empty() && alias.EthAddress == "":
		err = sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"account must contain address in Ethereum Hex or Cosmos Bech32 format",
		)
	}

	if err != nil {
		return err
	}

	acc.BaseAccount = &authtypes.BaseAccount{
		Coins:         alias.Coins,
		Address:       alias.Address,
		AccountNumber: alias.AccountNumber,
		Sequence:      alias.Sequence,
	}
	acc.CodeHash = ethcmn.Hex2Bytes(alias.CodeHash)
	acc.StateRoot = ethcmn.HexToHash(alias.StorageRoot)

	if alias.PubKey != "" {
		acc.BaseAccount.PubKey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, alias.PubKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// String implements the fmt.Stringer interface
func (acc EthAccount) String() string {
	out, _ := yaml.Marshal(acc)
	return string(out)
}

func (acc EthAccount)  GetStorageRoot() ethcmn.Hash {
	return acc.StateRoot
}

func (acc EthAccount) IsEthAccount() bool {
	return true
}

type EthAccountPretty struct {
	authtypes.BaseAccountPretty `json:"base_account_pretty" yaml:"base_account_pretty"`
	CodeHash               []byte `json:"code_hash" yaml:"code_hash"`
	StateRoot string  `json:"state_root" yaml:"state_root"`	// merkle root of the storage trie
}

func (alia EthAccountPretty) Pretty2Acc() (EthAccount, error) {
	bsAcc, err := alia.BaseAccountPretty.Pretty2Acc()
	if err != nil {
		return EthAccount{}, err
	}

	ethAcc := EthAccount{
		BaseAccount: &bsAcc,
		CodeHash: alia.CodeHash,
		StateRoot: ethcmn.HexToHash(alia.StateRoot),
	}

	return ethAcc, nil
}

func (acc EthAccount) GetPrettyAccount() (EthAccountPretty, error) {
	if acc.BaseAccount == nil {
		return EthAccountPretty{}, errors.New("nil base account")
	}

	baseAccPretty, err := acc.BaseAccount.GetPrettyAccount()
	if err != nil {
		return EthAccountPretty{}, err
	}

	ethAccPretty := EthAccountPretty{
		BaseAccountPretty: baseAccPretty,
		CodeHash:          acc.CodeHash,
		StateRoot: acc.StateRoot.String(),
	}

	return ethAccPretty, nil
}

// RLPEncodeToBytes returns the rlp representation of an account.
func (acc EthAccount) RLPEncodeToBytes() ([]byte, error) {
	alias, err := acc.GetPrettyAccount()
	if err != nil {
		return nil, err
	}

	return rlp.EncodeToBytes(alias)
}

// RLPDecodeBytes reduction account from rlp encode bytes
func (acc *EthAccount) RLPDecodeBytes(data []byte) error {
	var alia EthAccountPretty
	err := rlp.DecodeBytes(data, &alia)
	if err != nil {
		return err
	}
	*acc, err = alia.Pretty2Acc()
	return err
}

func (acc *EthAccount) EncodeRLP(w io.Writer) error {
	alias, err := acc.GetPrettyAccount()
	if err != nil {
		return err
	}

	if err = rlp.Encode(w, exported.EthAcc); err != nil {
		return err
	}
	return rlp.Encode(w, alias)
}

func (acc *EthAccount) DecodeRLP(s *rlp.Stream) error {
	var alia EthAccountPretty
	err := s.Decode(&alia)
	if err != nil {
		return err
	}

	*acc, err = alia.Pretty2Acc()

	return err
}

func (acc *EthAccount) Copy() *EthAccount {
	return &EthAccount{
		BaseAccount: acc.BaseAccount.Copy(),
		CodeHash:    acc.CodeHash,
		StateRoot:   acc.StateRoot,
	}
}
