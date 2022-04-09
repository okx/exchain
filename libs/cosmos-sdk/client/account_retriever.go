package client

import (
	clictx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
)

// Account defines a read-only version of the auth module's AccountI.
type Account interface {
	GetAddress() sdk.AccAddress
	GetPubKey() crypto.PubKey// can return nil.
	GetAccountNumber() uint64
	GetSequence() uint64
}

// AccountRetriever defines the interfaces required by transactions to
// ensure an account exists and to be able to query for account fields necessary
// for signing.
type AccountRetriever interface {
	GetAccount(clientCtx clictx.CLIContext, addr sdk.AccAddress) (Account, error)
	GetAccountWithHeight(clientCtx clictx.CLIContext, addr sdk.AccAddress) (Account, int64, error)
	EnsureExists(clientCtx clictx.CLIContext, addr sdk.AccAddress) error
	GetAccountNumberSequence(clientCtx clictx.CLIContext, addr sdk.AccAddress) (accNum uint64, accSeq uint64, err error)
}
