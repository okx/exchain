package keeper_test

import sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"

var (
	// The default power validators are initialized to have within tests
	InitTokens = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
)
