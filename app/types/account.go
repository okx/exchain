package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	"github.com/tendermint/go-amino"
)

var _ exported.Account = (*EthAccount)(nil)
var _ exported.GenesisAccount = (*EthAccount)(nil)
var emptyCodeHash = crypto.Keccak256(nil)

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
	CodeHash               []byte      `json:"code_hash" yaml:"code_hash"`
	StateRoot              ethcmn.Hash `json:"state_root" yaml:"state_root"` // merkle root of the storage trie
}

func (acc *EthAccount) UnmarshalFromAmino(cdc *amino.Codec, data []byte) error {
	var dataLen uint64 = 0
	var baseAccountFlag bool

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		// all EthAccount fields are (2)
		if pbType != amino.Typ3_ByteLength {
			return fmt.Errorf("invalid pbType: %v", pbType)
		}
		data = data[1:]

		var n int
		dataLen, n, err = amino.DecodeUvarint(data)
		if err != nil {
			return err
		}

		data = data[n:]
		if len(data) < int(dataLen) {
			return fmt.Errorf("not enough data for field %d", pos)
		}
		subData := data[:dataLen]

		switch pos {
		case 1:
			baseAccountFlag = true
			if acc.BaseAccount == nil {
				acc.BaseAccount = &auth.BaseAccount{}
			} else {
				*acc.BaseAccount = auth.BaseAccount{}
			}
			err = acc.BaseAccount.UnmarshalFromAmino(cdc, subData)
			if err != nil {
				return err
			}
		case 2:
			acc.CodeHash = make([]byte, len(subData))
			copy(acc.CodeHash, subData)
		case 3:
			acc.StateRoot.SetBytes(subData)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	if !baseAccountFlag {
		acc.BaseAccount = nil
	}
	return nil
}

type componentAccount struct {
	ethAccount  EthAccount
	baseAccount authtypes.BaseAccount
}

func (acc EthAccount) Copy() sdk.Account {
	// we need only allocate one object on the heap with componentAccount
	var cacc componentAccount

	cacc.baseAccount.Address = acc.Address
	cacc.baseAccount.Coins = acc.Coins
	cacc.baseAccount.PubKey = acc.PubKey
	cacc.baseAccount.AccountNumber = acc.AccountNumber
	cacc.baseAccount.Sequence = acc.Sequence

	cacc.ethAccount.BaseAccount = &cacc.baseAccount
	cacc.ethAccount.CodeHash = acc.CodeHash
	cacc.ethAccount.StateRoot = acc.StateRoot

	return &cacc.ethAccount
}

func (acc EthAccount) DeepCopy() sdk.Account {
	newAccount := ProtoAccount().(*EthAccount)
	buff, err := acc.MarshalJSON()
	if err != nil {
		return nil
	}
	err = newAccount.UnmarshalJSON(buff)
	if err != nil {
		return nil
	}
	return newAccount
}

func (acc EthAccount) AminoSize(cdc *amino.Codec) int {
	size := 0
	if acc.BaseAccount != nil {
		baccSize := acc.BaseAccount.AminoSize(cdc)
		size += 1 + amino.UvarintSize(uint64(baccSize)) + baccSize
	}
	if len(acc.CodeHash) != 0 {
		size += 1 + amino.ByteSliceSize(acc.CodeHash)
	}
	size += 1 + amino.ByteSliceSize(acc.StateRoot.Bytes())
	return size
}

func (acc EthAccount) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(acc.AminoSize(cdc))
	err := acc.MarshalAminoTo(cdc, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (acc EthAccount) MarshalAminoTo(cdc *amino.Codec, buf *bytes.Buffer) error {
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
			return amino.NewSizerError(acc.BaseAccount, baccSize, buf.Len()-lenBeforeData)
		}
	}

	// field 2
	if len(acc.CodeHash) != 0 {
		const pbKey = 2<<3 | 2
		err := amino.EncodeByteSliceWithKeyToBuffer(buf, acc.CodeHash, pbKey)
		if err != nil {
			return err
		}
	}
	const pbKey = 3<<3 | 2
	err := amino.EncodeByteSliceWithKeyToBuffer(buf, acc.StateRoot.Bytes(), pbKey)
	if err != nil {
		return err
	}

	return nil
}

// ProtoAccount defines the prototype function for BaseAccount used for an
// AccountKeeper.
func ProtoAccount() exported.Account {
	return &EthAccount{
		BaseAccount: &auth.BaseAccount{},
		CodeHash:    ethcrypto.Keccak256(nil),
		StateRoot:   types.EmptyRootHash,
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
	StateRoot     string         `josn:"state_root" yaml:"state_root"`
}

// MarshalYAML returns the YAML representation of an account.
func (acc EthAccount) MarshalYAML() (interface{}, error) {
	ethAddress := ""
	if !sdk.IsWasmAddress(acc.Address) {
		ethAddress = acc.EthAddress().String()
	}
	alias := ethermintAccountPretty{
		Address:       acc.Address,
		EthAddress:    ethAddress,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
		StateRoot:     acc.StateRoot.String(),
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
		if !sdk.IsWasmAddress(acc.Address) {
			ethAddress = acc.EthAddress().String()
		}
	}

	alias := ethermintAccountPretty{
		Address:       acc.Address,
		EthAddress:    ethAddress,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
		StateRoot:     acc.StateRoot.String(),
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
	acc.StateRoot = ethcmn.HexToHash(alias.StateRoot)

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

// IsContract returns if the account contains contract code.
func (acc EthAccount) IsContract() bool {
	return !bytes.Equal(acc.CodeHash, emptyCodeHash)
}

func (acc EthAccount) GetStateRoot() ethcmn.Hash {
	return acc.StateRoot
}
