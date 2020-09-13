package v0_11

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	v010dex "github.com/okex/okexchain/x/dex/legacy/v0_10"
	v011dex "github.com/okex/okexchain/x/dex/legacy/v0_11"
	v010order "github.com/okex/okexchain/x/order/legacy/v0_10"
	v011order "github.com/okex/okexchain/x/order/legacy/v0_11"
	v010staking "github.com/okex/okexchain/x/staking/legacy/v0_10"
	v011staking "github.com/okex/okexchain/x/staking/legacy/v0_11"
	v010token "github.com/okex/okexchain/x/token/legacy/v0_10"
	v011token "github.com/okex/okexchain/x/token/legacy/v0_11"
)

// Migrate migrates exported state from v0.10.x to a v0.11.0 genesis state
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v010Codec := codec.New()
	codec.RegisterCrypto(v010Codec)

	v011Codec := codec.New()
	codec.RegisterCrypto(v011Codec)

	// migrate dex state
	if appState[v010dex.ModuleName] != nil {
		var dexGenState v010dex.GenesisState
		v010Codec.MustUnmarshalJSON(appState[v010dex.ModuleName], &dexGenState)

		delete(appState, v010dex.ModuleName) // delete old key in case the name changed
		appState[v011dex.ModuleName] = v011Codec.MustMarshalJSON(v011dex.Migrate(dexGenState))
	}

	// migrate order state
	if appState[v010order.ModuleName] != nil {
		var orderGenState v010order.GenesisState
		v010Codec.MustUnmarshalJSON(appState[v010order.ModuleName], &orderGenState)

		delete(appState, v010order.ModuleName) // delete old key in case the name changed
		appState[v011order.ModuleName] = v011Codec.MustMarshalJSON(v011order.Migrate(orderGenState))
	}

	// migrate staking state
	if appState[v010staking.ModuleName] != nil {
		var stakingGenState v010staking.GenesisState
		v010Codec.MustUnmarshalJSON(appState[v010staking.ModuleName], &stakingGenState)

		delete(appState, v010staking.ModuleName)
		appState[v011staking.ModuleName] = v011Codec.MustMarshalJSON(v011staking.Migrate(stakingGenState))
	}

	// migrate token state
	if appState[v010token.ModuleName] != nil {
		var tokenGenState v010token.GenesisState
		v010Codec.MustUnmarshalJSON(appState[v010token.ModuleName], &tokenGenState)

		delete(appState, v010token.ModuleName)
		appState[v011token.ModuleName] = v011Codec.MustMarshalJSON(v011token.Migrate(tokenGenState))
	}

	return appState
}
