package types

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func TestUnmarshalChainConfigFromAmino(t *testing.T) {
	config := &ChainConfig{
		HomesteadBlock:      sdk.OneInt(),
		DAOForkBlock:        sdk.OneInt(),
		DAOForkSupport:      true,
		EIP150Block:         sdk.OneInt(),
		EIP150Hash:          defaultEIP150Hash,
		EIP155Block:         sdk.OneInt(),
		EIP158Block:         sdk.OneInt(),
		ByzantiumBlock:      sdk.OneInt(),
		ConstantinopleBlock: sdk.OneInt(),
		PetersburgBlock:     sdk.OneInt(),
		IstanbulBlock:       sdk.OneInt(),
		MuirGlacierBlock:    sdk.ZeroInt(),
		YoloV2Block:         sdk.OneInt(),
		EWASMBlock:          sdk.OneInt(),
	}
	cdc := amino.NewCodec()
	RegisterCodec(cdc)

	data, err := cdc.MarshalBinaryBare(config)
	require.NoError(t, err)

	var configFromAmino ChainConfig
	err = cdc.UnmarshalBinaryBare(data, &configFromAmino)
	require.NoError(t, err)

	var configFromUnmarshaller ChainConfig
	configi, err := cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(data, &configFromUnmarshaller)
	require.NoError(t, err)
	configFromUnmarshaller = configi.(ChainConfig)

	require.EqualValues(t, configFromAmino, configFromUnmarshaller)
}

func BenchmarkUnmarshalChainConfigFromAmino(b *testing.B) {
	config := &ChainConfig{
		HomesteadBlock:      sdk.OneInt(),
		DAOForkBlock:        sdk.OneInt(),
		EIP150Block:         sdk.OneInt(),
		EIP150Hash:          defaultEIP150Hash,
		EIP155Block:         sdk.OneInt(),
		EIP158Block:         sdk.OneInt(),
		ByzantiumBlock:      sdk.OneInt(),
		ConstantinopleBlock: sdk.OneInt(),
		PetersburgBlock:     sdk.OneInt(),
		IstanbulBlock:       sdk.OneInt(),
		MuirGlacierBlock:    sdk.OneInt(),
		YoloV2Block:         sdk.OneInt(),
		EWASMBlock:          sdk.OneInt(),
	}
	cdc := amino.NewCodec()
	RegisterCodec(cdc)

	data, _ := cdc.MarshalBinaryBare(config)

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("amino", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var config ChainConfig
			_ = cdc.UnmarshalBinaryBare(data, &config)
		}
	})

	b.Run("unmarshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var config ChainConfig
			_config, _ := cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(data, &config)
			config = _config.(ChainConfig)
		}
	})
}
