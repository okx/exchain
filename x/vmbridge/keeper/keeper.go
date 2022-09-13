package keeper

import (
	"github.com/okex/exchain/x/vmbridge/types"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

type Keeper struct {
	cdc *codec.CodecProxy

	logger log.Logger

	evmKeeper  EVMKeeper
	wasmKeeper WASMKeeper
}

func NewKeeper(cdc *codec.CodecProxy, logger log.Logger, evmKeeper EVMKeeper, wasmKeeper WASMKeeper) *Keeper {
	logger = logger.With("module", types.ModuleName)
	return &Keeper{cdc: cdc, logger: logger, evmKeeper: evmKeeper, wasmKeeper: wasmKeeper}
}

func (k Keeper) Logger() log.Logger {
	return k.logger
}

func (k Keeper) getAminoCodec() *codec.Codec {
	return k.cdc.GetCdc()
}

func (k Keeper) GetProtoCodec() *codec.ProtoCodec {
	return k.cdc.GetProtocMarshal()
}
