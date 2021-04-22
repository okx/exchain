package v017

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	v016evm "github.com/okex/exchain/x/evm/legacy/v0_16"
	v017evm "github.com/okex/exchain/x/evm/legacy/v0_17"
	v011staking "github.com/okex/exchain/x/staking/legacy/v0_11"
	v017staking "github.com/okex/exchain/x/staking/legacy/v0_17"
)

// Migrate migrates exported state from v0.16 to a v0.17 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v016Codec := codec.New()
	codec.RegisterCrypto(v016Codec)

	v017Codec := codec.New()
	codec.RegisterCrypto(v017Codec)

	// migrate auth state
	if appState[v017evm.ModuleName] != nil {
		var evmState v016evm.GenesisState
		v016Codec.MustUnmarshalJSON(appState[v017evm.ModuleName], &evmState)

		delete(appState, v017evm.ModuleName) // delete old key in case the name changed
		appState[v017evm.ModuleName] = v017Codec.MustMarshalJSON(v017evm.Migrate(evmState))
	}

	// migrate auth state
	if appState[v017staking.ModuleName] != nil {
		var stakingState v011staking.GenesisState
		v016Codec.MustUnmarshalJSON(appState[v017staking.ModuleName], &stakingState)

		delete(appState, v017staking.ModuleName) // delete old key in case the name changed
		appState[v017staking.ModuleName] = v017Codec.MustMarshalJSON(v017staking.Migrate(stakingState))
	}

	return appState
}
