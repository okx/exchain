package types

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	DelPk1       = ed25519.GenPrivKey().PubKey()
	DelPk2       = ed25519.GenPrivKey().PubKey()
	DelPk3       = ed25519.GenPrivKey().PubKey()
	DelAddr1     = sdk.AccAddress(DelPk1.Address())
	DelAddr2     = sdk.AccAddress(DelPk2.Address())
	DelAddr3     = sdk.AccAddress(DelPk3.Address())
	EmptyDelAddr sdk.AccAddress

	ValPk1       = ed25519.GenPrivKey().PubKey()
	ValPk2       = ed25519.GenPrivKey().PubKey()
	ValPk3       = ed25519.GenPrivKey().PubKey()
	ValAddr1     = sdk.ValAddress(ValPk1.Address())
	valAddr2     = sdk.ValAddress(ValPk2.Address())
	valAddr3     = sdk.ValAddress(ValPk3.Address())
	EmptyValAddr sdk.ValAddress

	emptyPubkey crypto.PubKey
)
