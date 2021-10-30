package v018

import (
	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/genutil"
	v016evm "github.com/okex/exchain/x/evm/legacy/v0_16"
	v018evm "github.com/okex/exchain/x/evm/legacy/v0_18"
	v011staking "github.com/okex/exchain/x/staking/legacy/v0_11"
	v018staking "github.com/okex/exchain/x/staking/legacy/v0_18"
)

// Migrate migrates exported state from v0.16 to a v0.17 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v016Codec := codec.New()
	codec.RegisterCrypto(v016Codec)

	v018Codec := codec.New()
	codec.RegisterCrypto(v018Codec)

	// migrate auth state
	if appState[v018evm.ModuleName] != nil {
		var evmState v016evm.GenesisState
		v016Codec.MustUnmarshalJSON(appState[v018evm.ModuleName], &evmState)

		delete(appState, v018evm.ModuleName) // delete old key in case the name changed
		appState[v018evm.ModuleName] = v018Codec.MustMarshalJSON(v018evm.Migrate(evmState))
	}

	// migrate statking state
	if appState[v018staking.ModuleName] != nil {
		var stakingState v011staking.GenesisState
		v016Codec.MustUnmarshalJSON(appState[v018staking.ModuleName], &stakingState)

		delete(appState, v018staking.ModuleName) // delete old key in case the name changed
		appState[v018staking.ModuleName] = v018Codec.MustMarshalJSON(v018staking.Migrate(stakingState))
	}

	return appState
}
