package simapp

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	cosmoscryptocodec "github.com/okex/exchain/libs/cosmos-sdk/crypto/ibc-codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

func makeCodec(bm module.BasicManager) types.InterfaceRegistry {
	// cdc := codec.NewLegacyAmino()

	// bm.RegisterLegacyAminoCodec(cdc)
	// std.RegisterLegacyAminoCodec(cdc)
	interfaceReg := types.NewInterfaceRegistry()
	bm.RegisterInterfaces(interfaceReg)
	cosmoscryptocodec.RegisterInterfaces(interfaceReg)

	return interfaceReg
}

func TestGetSimulationLog(t *testing.T) {
	//cdc := makeCodec(ModuleBasics)
	cdc := codec.NewCodecProxy(codec.NewProtoCodec(makeCodec(ModuleBasics)), codec.New())

	decoders := make(sdk.StoreDecoderRegistry)
	decoders[authtypes.StoreKey] = func(cdc *codec.Codec, kvAs, kvBs tmkv.Pair) string { return "10" }

	tests := []struct {
		store       string
		kvPairs     []tmkv.Pair
		expectedLog string
	}{
		{
			"Empty",
			[]tmkv.Pair{{}},
			"",
		},
		{
			authtypes.StoreKey,
			// todo old one is MustMarshal. does it want to test pb codec?
			[]tmkv.Pair{{Key: authtypes.GlobalAccountNumberKey, Value: cdc.GetCdc().MustMarshalBinaryBare(uint64(10))}},
			"10",
		},
		{
			"OtherStore",
			[]tmkv.Pair{{Key: []byte("key"), Value: []byte("value")}},
			fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", []byte("key"), []byte("value"), []byte("key"), []byte("value")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.store, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, GetSimulationLog(tt.store, decoders, tt.kvPairs, tt.kvPairs), tt.store)
		})
	}
}
