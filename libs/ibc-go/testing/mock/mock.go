package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/spf13/cobra"

	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

const (
	ModuleName = "mock"

	Version = "mock-version"
)

var (
	MockAcknowledgement             = channeltypes.NewResultAcknowledgement([]byte("mock acknowledgement"))
	MockFailAcknowledgement         = channeltypes.NewErrorAcknowledgement("mock failed acknowledgement")
	MockPacketData                  = []byte("mock packet data")
	MockFailPacketData              = []byte("mock failed packet data")
	MockAsyncPacketData             = []byte("mock async packet data")
	MockRecvCanaryCapabilityName    = "mock receive canary capability name"
	MockAckCanaryCapabilityName     = "mock acknowledgement canary capability name"
	MockTimeoutCanaryCapabilityName = "mock timeout canary capability name"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

//var _ porttypes.IBCModule = AppModule{}

// Expected Interface
// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
	IsBound(ctx sdk.Context, portID string) bool
}

// AppModuleBasic is the mock AppModuleBasic.
type AppModuleBasic struct{}

func (a AppModuleBasic) RegisterCodec(c *codec.Codec) {
	return
}

// Name implements AppModuleBasic interface.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic interface.
//func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces implements AppModuleBasic interface.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {}

// DefaultGenesis implements AppModuleBasic interface.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	r, _ := json.Marshal([]byte{})
	return r
}

// ValidateGenesis implements the AppModuleBasic interface.
func (AppModuleBasic) ValidateGenesis(json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes implements AppModuleBasic interface.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx context.CLIContext, rtr *mux.Router) {}

// RegisterGRPCGatewayRoutes implements AppModuleBasic interface.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(_ context.CLIContext, _ *runtime.ServeMux) {}

// GetTxCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetTxCmd(proxy *codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetQueryCmd(codec *codec.Codec) *cobra.Command {
	return nil
}

// AppModule represents the AppModule for the mock module.
type AppModule struct {
	AppModuleBasic
	scopedKeeper capabilitykeeper.ScopedKeeper
	portKeeper   PortKeeper
	ibcApps      []*MockIBCApp
}

func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// NewAppModule returns a mock AppModule instance.
func NewAppModule(sk capabilitykeeper.ScopedKeeper, pk PortKeeper) AppModule {
	return AppModule{
		scopedKeeper: sk,
		portKeeper:   pk,
	}
}

// RegisterInvariants implements the AppModule interface.
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route implements the AppModule interface.
func (am AppModule) Route() string {
	return sdk.NewRoute(ModuleName, nil).Path()
}

// QuerierRoute implements the AppModule interface.
func (AppModule) QuerierRoute() string {
	return ""
}

// LegacyQuerierHandler implements the AppModule interface.
// func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
// 	return nil
// }

// RegisterServices implements the AppModule interface.
func (am AppModule) RegisterServices(module.Configurator) {}

// InitGenesis implements the AppModule interface.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// bind mock port ID
	//cap := am.portKeeper.BindPort(ctx, ModuleName)
	//am.scopedKeeper.ClaimCapability(ctx, cap, host.PortPath(ModuleName))

	for _, ibcApp := range am.ibcApps {
		if ibcApp.PortID != "" && !am.portKeeper.IsBound(ctx, ibcApp.PortID) {
			// bind mock portID
			cap := am.portKeeper.BindPort(ctx, ibcApp.PortID)
			ibcApp.ScopedKeeper.ClaimCapability(ctx, cap, host.PortPath(ibcApp.PortID))
		}
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis implements the AppModule interface.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// OnChanOpenInit implements the IBCModule interface.
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context, _ channeltypes.Order, _ []string, portID string,
	channelID string, chanCap *capabilitytypes.Capability, _ channeltypes.Counterparty, v string,
) (string, error) {
	// Claim channel capability passed back by IBC module
	if err := am.scopedKeeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	return "", nil
}

// OnChanOpenTry implements the IBCModule interface.
func (am AppModule) OnChanOpenTry(
	ctx sdk.Context, _ channeltypes.Order, _ []string, portID string,
	channelID string, chanCap *capabilitytypes.Capability, _ channeltypes.Counterparty, _, _ string,
) (string, error) {
	// Claim channel capability passed back by IBC module
	if err := am.scopedKeeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	return "", nil
}

// OnChanOpenAck implements the IBCModule interface.
func (am AppModule) OnChanOpenAck(sdk.Context, string, string, string, string) error {
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface.
func (am AppModule) OnChanOpenConfirm(sdk.Context, string, string) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface.
func (am AppModule) OnChanCloseInit(sdk.Context, string, string) error {
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface.
func (am AppModule) OnChanCloseConfirm(sdk.Context, string, string) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface.
func (am AppModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	// set state by claiming capability to check if revert happens return
	_, err := am.scopedKeeper.NewCapability(ctx, MockRecvCanaryCapabilityName+strconv.Itoa(int(packet.GetSequence())))
	if err != nil {
		// application callback called twice on same packet sequence
		// must never occur
		panic(err)
	}
	if bytes.Equal(MockPacketData, packet.GetData()) {
		return MockAcknowledgement
	} else if bytes.Equal(MockAsyncPacketData, packet.GetData()) {
		return nil
	}

	return MockFailAcknowledgement
}

// OnAcknowledgementPacket implements the IBCModule interface.
func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, _ []byte, _ sdk.AccAddress) error {
	_, err := am.scopedKeeper.NewCapability(ctx, MockAckCanaryCapabilityName+strconv.Itoa(int(packet.GetSequence())))
	if err != nil {
		// application callback called twice on same packet sequence
		// must never occur
		panic(err)
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface.
func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, _ sdk.AccAddress) error {
	_, err := am.scopedKeeper.NewCapability(ctx, MockTimeoutCanaryCapabilityName+strconv.Itoa(int(packet.GetSequence())))
	if err != nil {
		// application callback called twice on same packet sequence
		// must never occur
		panic(err)
	}

	return nil
}

// NegotiateAppVersion implements the IBCModule interface.
func (am AppModule) NegotiateAppVersion(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionID string,
	portID string,
	counterparty channeltypes.Counterparty,
	proposedVersion string,
) (string, error) {
	if proposedVersion != Version { // allow testing of error scenarios
		return "", errors.New("failed to negotiate app version")
	}

	return Version, nil
}
