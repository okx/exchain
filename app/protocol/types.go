package protocol

import (
	"encoding/json"

	"github.com/okex/okchain/x/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/okex/okchain/x/staking"

	"github.com/cosmos/cosmos-sdk/x/crisis"

	"github.com/okex/okchain/x/backend"
	distr "github.com/okex/okchain/x/distribution"
	"github.com/okex/okchain/x/stream"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Protocol shows the expected behavior for any protocol version
type Protocol interface {
	GetVersion() uint64

	// load base installation for each protocol
	LoadContext()
	Init()
	GetCodec() *codec.Codec

	// gracefully stop okchaind
	CheckStopped()

	// setter
	SetLogger(log log.Logger) Protocol
	SetParent(parent Parent) Protocol

	//getter
	GetParent() Parent

	// get specific keeper
	GetBackendKeeper() backend.Keeper
	GetStreamKeeper() stream.Keeper
	GetCrisisKeeper() crisis.Keeper
	GetStakingKeeper() staking.Keeper
	GetDistrKeeper() distr.Keeper
	GetSlashingKeeper() slashing.Keeper
	GetTokenKeeper() token.Keeper

	// fit cm36
	GetKVStoreKeysMap() map[string]*sdk.KVStoreKey
	GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey
	ExportGenesis(ctx sdk.Context) map[string]json.RawMessage
}

// Parent shows the expected behavior of BaseApp(hooks)
type Parent interface {
	DeliverTx(abci.RequestDeliverTx) abci.ResponseDeliverTx
	PushInitChainer(initChainer sdk.InitChainer)
	PushBeginBlocker(beginBlocker sdk.BeginBlocker)
	PushEndBlocker(endBlocker sdk.EndBlocker)
	PushAnteHandler(ah sdk.AnteHandler)
	SetRouter(router sdk.Router, queryRouter sdk.QueryRouter)
}
