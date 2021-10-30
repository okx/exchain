package v018

import (
	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/genutil"
	v018evm "github.com/okex/exchain/x/evm/legacy/v0_18"
	v018staking "github.com/okex/exchain/x/staking/legacy/v0_18"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestMigrate tests v017
func TestMigrate(t *testing.T) {
	v016Codec := codec.New()
	codec.RegisterCrypto(v016Codec)

	v018Codec := codec.New()
	codec.RegisterCrypto(v018Codec)

	appState := genutil.AppMap{
		"staking": []byte(`{"params":{"bond_denom":"okt","epoch":252,"max_bonded_validators":21,"max_validators_to_add_shares":30,"min_delegation":"0.000100000000000000","min_self_delegation":"10000.000000000000000000","unbonding_time":"1209600000000000"}}`),
		"evm":     []byte(`{"params":{"enable_call":true,"enable_create":true,"evm_denom":"okt","extra_eips":null}}`),
	}
	statsMigrate := Migrate(appState)

	// stakingState
	var stakingState v018staking.GenesisState
	v018Codec.MustUnmarshalJSON(statsMigrate[v018staking.ModuleName], &stakingState)
	require.NotNil(t, stakingState)

	// evmState
	var evmState v018evm.GenesisState
	v018Codec.MustUnmarshalJSON(statsMigrate[v018evm.ModuleName], &evmState)
	require.True(t, evmState.Params.EnableContractDeploymentWhitelist)
	require.True(t, evmState.Params.EnableContractBlockedList)
}
