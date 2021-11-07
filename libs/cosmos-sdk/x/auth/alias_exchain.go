package auth

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"
)

type (
	Account = exported.Account
	ObserverI = keeper.ObserverI
)
