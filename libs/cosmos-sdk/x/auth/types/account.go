package types

import (
	"bytes"
	"errors"
	"time"

	"github.com/tendermint/go-amino"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"gopkg.in/yaml.v2"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	cryptoamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
)

//-----------------------------------------------------------------------------
// BaseAccount

var _ exported.Account = (*BaseAccount)(nil)
var _ exported.GenesisAccount = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
}

func (acc BaseAccount) Copy() interface{} {
	return NewBaseAccount(acc.Address, acc.Coins, acc.PubKey, acc.AccountNumber, acc.Sequence)
}

func (acc BaseAccount) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [5]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2, 4 << 3, 5 << 3}
	for pos := 1; pos < 6; pos++ {
		lBeforeKey := buf.Len()
		var noWrite bool
		err := buf.WriteByte(fieldKeysType[pos-1])
		if err != nil {
			return nil, err
		}

		switch pos {
		case 1:
			addressLen := len(acc.Address)
			if addressLen == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarint(&buf, uint64(addressLen))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(acc.Address)
			if err != nil {
				return nil, err
			}
		case 2:
			coinsLen := len(acc.Coins)
			if coinsLen == 0 {
				noWrite = true
				break
			}
			if coinsLen == 1 {
				data, err := acc.Coins[0].MarshalToAmino()
				if err != nil {
					return nil, err
				}
				err = amino.EncodeUvarint(&buf, uint64(len(data)))
				if err != nil {
					return nil, err
				}
				_, err = buf.Write(data)
				if err != nil {
					return nil, err
				}
			} else {
				buf.Truncate(lBeforeKey)
				for _, coin := range acc.Coins {
					err := buf.WriteByte(fieldKeysType[pos-1])
					if err != nil {
						return nil, err
					}
					data, err := coin.MarshalToAmino()
					if err != nil {
						return nil, err
					}
					err = amino.EncodeUvarint(&buf, uint64(len(data)))
					if err != nil {
						return nil, err
					}
					_, err = buf.Write(data)
					if err != nil {
						return nil, err
					}
				}
			}
		case 3:
			if acc.PubKey == nil {
				noWrite = true
				break
			}
			data, err := cryptoamino.MarshalPubKeyToAminoWithTypePrefix(acc.PubKey)
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarint(&buf, uint64(len(data)))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(data)
			if err != nil {
				return nil, err
			}
		case 4:
			if acc.AccountNumber == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarint(&buf, acc.AccountNumber)
			if err != nil {
				return nil, err
			}
		case 5:
			if acc.Sequence == 0 {
				noWrite = true
				break
			}
			err := amino.EncodeUvarint(&buf, acc.Sequence)
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}

		if noWrite {
			buf.Truncate(lBeforeKey)
		}
	}
	return buf.Bytes(), nil
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address sdk.AccAddress, coins sdk.Coins,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() exported.Account {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr sdk.AccAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && acc.Address != nil &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

type baseAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	if acc.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}
