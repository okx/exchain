package v017

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	v017evm "github.com/okex/exchain/x/evm/legacy/v0_17"
	v017staking "github.com/okex/exchain/x/staking/legacy/v0_17"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestMigrate tests v017
func TestMigrate(t *testing.T) {
	v016Codec := codec.New()
	codec.RegisterCrypto(v016Codec)

	v017Codec := codec.New()
	codec.RegisterCrypto(v017Codec)

	appState := genutil.AppMap{
		"staking": []byte(`{"params":{"bond_denom":"okt","epoch":252,"max_bonded_validators":21,"max_validators_to_add_shares":30,"min_delegation":"0.000100000000000000","min_self_delegation":"10000.000000000000000000","unbonding_time":"1209600000000000"}}`),
		"evm":     []byte(`{"params":{"enable_call":true,"enable_create":true,"evm_denom":"okt","extra_eips":null}}`),
	}
	statsMigrate := Migrate(appState)

	// stakingState
	var stakingState v017staking.GenesisState
	v017Codec.MustUnmarshalJSON(statsMigrate[v017staking.ModuleName], &stakingState)
	require.NotNil(t, stakingState)

	// evmState
	var evmState v017evm.GenesisState
	v017Codec.MustUnmarshalJSON(statsMigrate[v017evm.ModuleName], &evmState)
	require.True(t, evmState.Params.EnableContractDeploymentWhitelist)
	require.True(t, evmState.Params.EnableContractBlockedList)
}
