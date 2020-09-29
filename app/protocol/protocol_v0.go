package protocol

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/app/utils"
	"github.com/okex/okexchain/x/ammswap"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/common/proto"
	"github.com/okex/okexchain/x/common/version"
	"github.com/okex/okexchain/x/debug"
	"github.com/okex/okexchain/x/dex"
	dexClient "github.com/okex/okexchain/x/dex/client"
	distr "github.com/okex/okexchain/x/distribution"
	"github.com/okex/okexchain/x/farm"
	"github.com/okex/okexchain/x/genutil"
	"github.com/okex/okexchain/x/gov"
	"github.com/okex/okexchain/x/gov/keeper"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/params"
	paramsclient "github.com/okex/okexchain/x/params/client"
	"github.com/okex/okexchain/x/staking"
	"github.com/okex/okexchain/x/stream"
	"github.com/okex/okexchain/x/token"
	"github.com/okex/okexchain/x/upgrade"
	upgradeClient "github.com/okex/okexchain/x/upgrade/client"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	// check the implements of ProtocolV0
	_ Protocol = (*ProtocolV0)(nil)

	// DefaultCLIHome is the directory for okexchaincli
	DefaultCLIHome = os.ExpandEnv("$HOME/.okexchaincli")
	// DefaultNodeHome is the directory for okexchaind
	DefaultNodeHome = os.ExpandEnv("$HOME/.okexchaind")

	// ModuleBasics is in charge of setting up basic, non-dependant module elements,
	// such as codec registration and genesis verification
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			upgradeClient.ProposalHandler, paramsclient.ProposalHandler,
			dexClient.DelistProposalHandler, distr.ProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},

		// okexchain extended
		token.AppModuleBasic{},
		dex.AppModuleBasic{},
		order.AppModuleBasic{},
		backend.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		stream.AppModuleBasic{},
		debug.AppModuleBasic{},
		ammswap.AppModuleBasic{},
		farm.AppModuleBasic{},
	)

	// module account permissions for bankKeeper and supplyKeeper
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            nil,
		token.ModuleName:          {supply.Minter, supply.Burner},
		order.ModuleName:          nil,
		backend.ModuleName:        nil,
		dex.ModuleName:            nil,
		ammswap.ModuleName:        {supply.Minter, supply.Burner},
		farm.ModuleName:           nil,
		farm.YieldFarmingName:     nil,
	}
)

// ProtocolV0 is the struct of the original protocol of okexchain
type ProtocolV0 struct {
	parent         Parent
	version        uint64
	cdc            *codec.Codec
	logger         log.Logger
	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	tokenKeeper    token.Keeper
	dexKeeper      dex.Keeper
	orderKeeper    order.Keeper
	swapKeeper     ammswap.Keeper
	protocolKeeper proto.ProtocolKeeper
	backendKeeper  backend.Keeper
	streamKeeper   stream.Keeper
	upgradeKeeper  upgrade.Keeper
	debugKeeper    debug.Keeper
	farmKeeper     farm.Keeper

	stopped     bool
	anteHandler sdk.AnteHandler // ante handler for fee and auth
	router      sdk.Router      // handle any kind of message
	queryRouter sdk.QueryRouter // router for redirecting query calls

	// the module manager
	mm *module.Manager
}

// NewProtocolV0 creates a new instance of NewProtocolV0
func NewProtocolV0(
	parent Parent, version uint64, log log.Logger, invCheckPeriod uint, pk proto.ProtocolKeeper,
) *ProtocolV0 {
	return &ProtocolV0{
		parent:         parent,
		version:        version,
		logger:         log,
		invCheckPeriod: invCheckPeriod,
		protocolKeeper: pk,
		keys:           kvStoreKeysMap,
		tkeys:          transientStoreKeysMap,
		router:         baseapp.NewRouter(),
		queryRouter:    baseapp.NewQueryRouter(),
	}
}

// GetVersion gets the version of this protocol
func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

// LoadContext updates the context for the app after the upgrade of protocol
func (p *ProtocolV0) LoadContext() {
	p.logger.Debug("Protocol V0: LoadContext")
	p.setCodec()
	p.produceKeepers()
	p.setManager()
	p.registerRouters()
	p.setAnteHandler()

	p.parent.PushInitChainer(p.InitChainer)
	p.parent.PushBeginBlocker(p.BeginBlocker)
	p.parent.PushEndBlocker(p.EndBlocker)
}

// GetCodec gets tx codec
func (p *ProtocolV0) GetCodec() *codec.Codec {
	if p.cdc == nil {
		panic("Invalid cdc from ProtocolV0")
	}
	return p.cdc
}

// CheckStopped gives a quick check whether okexchain needs stopped
func (p *ProtocolV0) CheckStopped() {
	if p.stopped {
		p.logger.Info("OKExChain is going to exit")
		server.Stop()
		p.logger.Info("OKExChain was stopped")
		select {}
	}
}

// GetBackendKeeper gets backend keeper
func (p *ProtocolV0) GetBackendKeeper() backend.Keeper {
	return p.backendKeeper
}

// GetStreamKeeper gets stream keeper
func (p *ProtocolV0) GetStreamKeeper() stream.Keeper {
	return p.streamKeeper
}

// GetCrisisKeeper gets crisis keeper
func (p *ProtocolV0) GetCrisisKeeper() crisis.Keeper {
	return p.crisisKeeper
}

// GetStakingKeeper gets staking keeper
func (p *ProtocolV0) GetStakingKeeper() staking.Keeper {
	return p.stakingKeeper
}

// GetDistrKeeper gets distr keeper
func (p *ProtocolV0) GetDistrKeeper() distr.Keeper {
	return p.distrKeeper
}

// GetSlashingKeeper gets slashing keeper
func (p *ProtocolV0) GetSlashingKeeper() slashing.Keeper {
	return p.slashingKeeper
}

// GetTokenKeeper gets token keeper
func (p *ProtocolV0) GetTokenKeeper() token.Keeper {
	return p.tokenKeeper
}

// GetKVStoreKeysMap gets the map of kv store keys
func (p *ProtocolV0) GetKVStoreKeysMap() map[string]*sdk.KVStoreKey {
	return p.keys
}

// GetTransientStoreKeysMap gets the map of transient store keys
func (p *ProtocolV0) GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey {
	return p.tkeys
}

// nolint
func (p *ProtocolV0) Init() {}

func (p *ProtocolV0) setCodec() {
	p.cdc = MakeCodec()
}

// produceKeepers initializes all keepers declared in the ProtocolV0 struct
func (p *ProtocolV0) produceKeepers() {
	// get config
	appConfig, err := config.ParseConfig()
	if err != nil {
		p.logger.Error(fmt.Sprintf("the config of OKExChain was parsed error : %s", err.Error()))
		panic(err)
	}

	// 1.init params keeper and subspaces
	p.paramsKeeper = params.NewKeeper(
		p.cdc, p.keys[params.StoreKey], p.tkeys[params.TStoreKey], params.DefaultCodespace,
	)
	authSubspace := p.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := p.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := p.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := p.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := p.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := p.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := p.paramsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := p.paramsKeeper.Subspace(crisis.DefaultParamspace)
	tokenSubspace := p.paramsKeeper.Subspace(token.DefaultParamspace)
	orderSubspace := p.paramsKeeper.Subspace(order.DefaultParamspace)
	upgradeSubspace := p.paramsKeeper.Subspace(upgrade.DefaultParamspace)
	dexSubspace := p.paramsKeeper.Subspace(dex.DefaultParamspace)
	swapSubspace := p.paramsKeeper.Subspace(ammswap.DefaultParamspace)
	farmSubspace := p.paramsKeeper.Subspace(farm.DefaultParamspace)

	// 2.add keepers
	p.accountKeeper = auth.NewAccountKeeper(p.cdc, p.keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	p.bankKeeper = bank.NewBaseKeeper(p.accountKeeper, bankSubspace, bank.DefaultCodespace, p.moduleAccountAddrs())
	p.paramsKeeper.SetBankKeeper(p.bankKeeper)
	p.supplyKeeper = supply.NewKeeper(p.cdc, p.keys[supply.StoreKey], p.accountKeeper, p.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(p.cdc, p.keys[staking.StoreKey], p.tkeys[staking.TStoreKey],
		p.supplyKeeper, stakingSubspace, staking.DefaultCodespace)

	p.paramsKeeper.SetStakingKeeper(stakingKeeper)
	p.mintKeeper = mint.NewKeeper(
		p.cdc,
		p.keys[mint.StoreKey],
		mintSubspace, &stakingKeeper,
		p.supplyKeeper,
		auth.FeeCollectorName,
		farm.ModuleName,
	)

	p.distrKeeper = distr.NewKeeper(p.cdc, p.keys[distr.StoreKey],
		distrSubspace, &stakingKeeper, p.supplyKeeper,
		distr.DefaultCodespace, auth.FeeCollectorName, p.moduleAccountAddrs(),
	)

	p.slashingKeeper = slashing.NewKeeper(
		p.cdc, p.keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)

	p.crisisKeeper = crisis.NewKeeper(crisisSubspace, p.invCheckPeriod, p.supplyKeeper, auth.FeeCollectorName)

	p.tokenKeeper = token.NewKeeper(
		p.bankKeeper, tokenSubspace, auth.FeeCollectorName, p.supplyKeeper,
		p.keys[token.StoreKey], p.keys[token.KeyLock],
		p.cdc, appConfig.BackendConfig.EnableBackend)

	p.dexKeeper = dex.NewKeeper(auth.FeeCollectorName, p.supplyKeeper, dexSubspace, p.tokenKeeper, &stakingKeeper,
		p.bankKeeper, p.keys[dex.StoreKey], p.keys[dex.TokenPairStoreKey], p.cdc)

	p.orderKeeper = order.NewKeeper(
		p.tokenKeeper, p.supplyKeeper, p.dexKeeper, orderSubspace, auth.FeeCollectorName,
		p.keys[order.OrderStoreKey], p.cdc, appConfig.BackendConfig.EnableBackend, orderMetrics,
	)

	p.swapKeeper = ammswap.NewKeeper(p.supplyKeeper, p.tokenKeeper, p.cdc, p.keys[ammswap.StoreKey], swapSubspace)

	p.streamKeeper = stream.NewKeeper(p.orderKeeper, p.tokenKeeper, &p.dexKeeper, &p.accountKeeper,
		p.cdc, p.logger, appConfig, streamMetrics)

	p.backendKeeper = backend.NewKeeper(p.orderKeeper, p.tokenKeeper, &p.dexKeeper, p.streamKeeper.GetMarketKeeper(),
		p.cdc, p.logger, appConfig.BackendConfig)
	// 3.register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(&p.paramsKeeper)).
		AddRoute(dex.RouterKey, dex.NewProposalHandler(&p.dexKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewAppUpgradeProposalHandler(&p.upgradeKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(p.distrKeeper))
	govProposalHandlerRouter := keeper.NewProposalHandlerRouter()
	govProposalHandlerRouter.AddRoute(params.RouterKey, &p.paramsKeeper).
		AddRoute(dex.RouterKey, &p.dexKeeper).
		AddRoute(upgrade.RouterKey, &p.upgradeKeeper)
	p.govKeeper = gov.NewKeeper(
		p.cdc, p.keys[gov.StoreKey], p.paramsKeeper, govSubspace,
		p.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
		p.bankKeeper, govProposalHandlerRouter, auth.FeeCollectorName,
	)
	p.paramsKeeper.SetGovKeeper(p.govKeeper)
	p.dexKeeper.SetGovKeeper(p.govKeeper)
	// 4.register the staking hooks
	p.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()),
	)
	p.upgradeKeeper = upgrade.NewKeeper(
		p.cdc, p.keys[upgrade.StoreKey], p.protocolKeeper, p.stakingKeeper, p.bankKeeper, upgradeSubspace,
	)
	p.debugKeeper = debug.NewDebugKeeper(p.cdc, p.keys[debug.StoreKey], p.orderKeeper, p.stakingKeeper, auth.FeeCollectorName, p.Stop)
	p.farmKeeper = farm.NewKeeper(auth.FeeCollectorName, p.supplyKeeper,
		p.tokenKeeper, p.swapKeeper,
		farmSubspace, p.keys[farm.StoreKey], p.cdc)
}

// moduleAccountAddrs returns all the module account addresses
func (p *ProtocolV0) moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// setManager sets module.Manager in protocolV0
func (p *ProtocolV0) setManager() {

	p.mm = module.NewManager(
		genaccounts.NewAppModule(p.accountKeeper),
		genutil.NewAppModule(p.accountKeeper, p.stakingKeeper, p.parent.DeliverTx),
		auth.NewAppModule(p.accountKeeper),
		bank.NewAppModule(p.bankKeeper, p.accountKeeper),
		crisis.NewAppModule(&p.crisisKeeper),
		supply.NewAppModule(p.supplyKeeper, p.accountKeeper),
		params.NewAppModule(p.paramsKeeper),
		mint.NewAppModule(p.mintKeeper),
		slashing.NewAppModule(p.slashingKeeper, p.stakingKeeper),
		staking.NewAppModule(p.stakingKeeper, p.accountKeeper, p.supplyKeeper),
		distr.NewAppModule(p.distrKeeper, p.supplyKeeper),
		gov.NewAppModule(version.ProtocolVersionV0, p.govKeeper, p.supplyKeeper),
		order.NewAppModule(version.ProtocolVersionV0, p.orderKeeper, p.supplyKeeper),
		token.NewAppModule(version.ProtocolVersionV0, p.tokenKeeper, p.supplyKeeper),
		ammswap.NewAppModule(p.swapKeeper),

		// TODO
		dex.NewAppModule(version.ProtocolVersionV0, p.dexKeeper, p.supplyKeeper),
		backend.NewAppModule(p.backendKeeper),
		stream.NewAppModule(p.streamKeeper),
		upgrade.NewAppModule(p.upgradeKeeper),
		debug.NewAppModule(p.debugKeeper),
		farm.NewAppModule(p.farmKeeper),
	)

	// ORDER SETTING
	p.mm.SetOrderBeginBlockers(
		stream.ModuleName,
		order.ModuleName,
		token.ModuleName,
		dex.ModuleName,
		mint.ModuleName,
		distr.ModuleName,
		slashing.ModuleName,
		staking.ModuleName,
		farm.ModuleName,
	)

	p.mm.SetOrderEndBlockers(
		crisis.ModuleName,
		gov.ModuleName,
		dex.ModuleName,
		order.ModuleName,
		staking.ModuleName,
		backend.ModuleName,
		stream.ModuleName,
		upgrade.ModuleName,
	)

	p.mm.SetOrderInitGenesis(
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		token.ModuleName,
		dex.ModuleName,
		order.ModuleName,
		ammswap.ModuleName,
		upgrade.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName,
		params.ModuleName,
		farm.ModuleName,
	)
}

// registerRouters registers Routers by Manager
func (p *ProtocolV0) registerRouters() {
	p.mm.RegisterInvariants(&p.crisisKeeper)
	p.mm.RegisterRoutes(p.router, p.queryRouter)
	p.parent.SetRouter(p.router, p.queryRouter)
}

// setAnteHandler sets ante handler
func (p *ProtocolV0) setAnteHandler() {
	p.anteHandler = auth.NewAnteHandler(
		p.accountKeeper,
		p.supplyKeeper,
		auth.DefaultSigVerificationGasConsumer,
		validateMsgHook(p.orderKeeper),
		isSystemFreeHook,
	)
	p.parent.PushAnteHandler(p.anteHandler)
}

// InitChainer initializes application state at genesis as a hook
func (p *ProtocolV0) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	p.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return p.mm.InitGenesis(ctx, genesisState)
}

// BeginBlocker set function to BaseApp as a hook
func (p *ProtocolV0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return p.mm.BeginBlock(ctx, req)
}

// EndBlocker sets function to BaseApp as a hook
func (p *ProtocolV0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return p.mm.EndBlock(ctx, req)
}

// Stop makes okexchain exit gracefully
func (p *ProtocolV0) Stop() {
	p.logger.Info(fmt.Sprintf("[%s]%s", utils.GoID, "OKExChain stops notification."))
	p.stopped = true
}

// MakeCodec registers codec from all the modules
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc
}

func validateMsgHook(orderKeeper order.Keeper) auth.ValidateMsgHandler {
	return func(newCtx sdk.Context, msgs []sdk.Msg) (res sdk.Result) {

		wrongMsgRes := sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  "It is not allowed that a transaction with more than one message contains placeOrder or cancelOrder message",
		}

		for _, msg := range msgs {
			switch assertedMsg := msg.(type) {
			case order.MsgNewOrders:
				if len(msgs) > 1 {
					res = wrongMsgRes
					break
				}
				res = order.ValidateMsgNewOrders(newCtx, orderKeeper, assertedMsg)
			case order.MsgCancelOrders:
				if len(msgs) > 1 {
					res = wrongMsgRes
					break
				}
				res = order.ValidateMsgCancelOrders(newCtx, orderKeeper, assertedMsg)
			}

			if !res.IsOK() {
				break
			}
		}
		return
	}
}

func isSystemFreeHook(ctx sdk.Context, msgs []sdk.Msg) bool {
	if ctx.BlockHeight() < 1 {
		return true
	}

	return false
}

// ExportGenesis exports the genesis state for whole protocol
func (p *ProtocolV0) ExportGenesis(ctx sdk.Context) map[string]json.RawMessage {
	return p.mm.ExportGenesis(ctx)
}

// SetLogger sets logger
func (p *ProtocolV0) SetLogger(log log.Logger) Protocol {
	p.logger = log
	return p
}

// SetParent sets parent implement
func (p *ProtocolV0) SetParent(parent Parent) Protocol {
	p.parent = parent
	return p
}

// GetParent gets parent implement
func (p *ProtocolV0) GetParent() Parent {
	if p.parent == nil {
		panic("parent is nil in protocol")
	}
	return p.parent
}
