package types

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/keys/ed25519"
	cryptotypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/types"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
)

// nolint:deadcode,unused,varcheck
var (
	delPk1       = ed25519.GenPrivKey().PubKey()
	delPk2       = ed25519.GenPrivKey().PubKey()
	delPk3       = ed25519.GenPrivKey().PubKey()
	delAddr1     = sdk.AccAddress(delPk1.Address())
	delAddr2     = sdk.AccAddress(delPk2.Address())
	delAddr3     = sdk.AccAddress(delPk3.Address())
	emptyDelAddr sdk.AccAddress

	valPk1       = ed25519.GenPrivKey().PubKey()
	valPk2       = ed25519.GenPrivKey().PubKey()
	valPk3       = ed25519.GenPrivKey().PubKey()
	valAddr1     = sdk.ValAddress(valPk1.Address())
	valAddr2     = sdk.ValAddress(valPk2.Address())
	valAddr3     = sdk.ValAddress(valPk3.Address())
	emptyValAddr sdk.ValAddress

	emptyPubkey cryptotypes.PubKey
)
