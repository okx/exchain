package transfer

import (
	"encoding/json"
	"fmt"

	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp/adapter"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

var (
	_ porttypes.Middleware = (*TransferModule)(nil)
)

type TransferModule struct {
	transfer.AppModule

	TKeeper keeper.Keeper
}

func (am TransferModule) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	//TODO implement me
	panic("implement me")
}

func (am TransferModule) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	//TODO implement me
	panic("implement me")
}

func (am TransferModule) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	//TODO implement me
	panic("implement me")
}

func TNewTransferModule(k keeper.Keeper, m *codec.CodecProxy) *TransferModule {
	ret := &TransferModule{}

	ret.AppModule = transfer.NewAppModule(k, m)
	ret.TKeeper = k
	return ret
}
func (am TransferModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.initGenesis(ctx, data)
}

func (am TransferModule) initGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	adapter.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	am.TKeeper.InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}

// ValidateGenesis performs genesis state validation for the mint module.
func (t TransferModule) ValidateGenesis(data json.RawMessage) error {
	if nil == data {
		return nil
	}
	var genState types.GenesisState
	if err := adapter.ModuleCdc.UnmarshalJSON(data, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// ExportGenesis returns the exported genesis state as raw bytes for the ibc-transfer
// module.
func (am TransferModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.exportGenesis(ctx)
}

func (am TransferModule) exportGenesis(ctx sdk.Context) json.RawMessage {
	gs := am.TKeeper.ExportGenesis(ctx)
	return adapter.ModuleCdc.MustMarshalJSON(gs)
}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// transfer module.
func (am TransferModule) DefaultGenesis() json.RawMessage {
	state := types.DefaultGenesisState()
	state.Params.SendEnabled = true
	state.Params.ReceiveEnabled = true
	return adapter.ModuleCdc.MustMarshalJSON(state)
}
