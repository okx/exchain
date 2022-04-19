package types

import (
	cryptotypes "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"time"
)

type AccountI interface {
	Copy() AccountI
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error
	GetPubKey() cryptotypes.PubKey
	SetPubKey(cryptotypes.PubKey) error
	GetAccountNumber() uint64
	SetAccountNumber(uint64) error
	GetSequence() uint64
	SetSequence(uint64) error
	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error
	SpendableCoins(blockTime time.Time) sdk.Coins
	String() string
}
