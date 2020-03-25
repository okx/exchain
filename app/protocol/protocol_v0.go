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
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	"github.com/okex/okchain/app/utils"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/debug"
	"github.com/okex/okchain/x/dex"
	dexClient "github.com/okex/okchain/x/dex/client"
	distr "github.com/okex/okchain/x/distribution"
	"github.com/okex/okchain/x/genutil"
	"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/gov/keeper"
	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/params"
	paramsclient "github.com/okex/okchain/x/params/client"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/stream"
	"github.com/okex/okchain/x/token"
	"github.com/okex/okchain/x/upgrade"
	upgradeClient "github.com/okex/okchain/x/upgrade/client"
)

var (
	// check the implements of ProtocolV0
	_ Protocol = (*ProtocolV0)(nil)

	// default home directories for okchaincli
	DefaultCLIHome = os.ExpandEnv("$HOME/.okchaincli")

	// default home directories for okchaind
	DefaultNodeHome = os.ExpandEnv("$HOME/.okchaind")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
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
			dexClient.DelistProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},

		// okchain extended
		token.AppModuleBasic{},
		dex.AppModuleBasic{},
		order.AppModuleBasic{},
		backend.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		stream.AppModuleBasic{},
		//for test
		debug.AppModuleBasic{},
	)

	// module account permissions
	// for bankKeeper and supplyKeeper
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
	}
)

type ProtocolV0 struct {
	parent  Parent
	version uint64
	cdc     *codec.Codec
	logger  log.Logger

	/*----------- necessary part for cm36 --------------*/
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
	protocolKeeper proto.ProtocolKeeper
	backendKeeper  backend.Keeper
	streamKeeper   stream.Keeper
	upgradeKeeper  upgrade.Keeper
	/******** for test **********/
	debugKeeper debug.Keeper
	/****************************/

	stopped bool

	anteHandler sdk.AnteHandler // ante handler for fee and auth
	router      sdk.Router      // handle any kind of message
	queryRouter sdk.QueryRouter // router for redirecting query calls

	// the module manager
	mm *module.Manager
}

// create new protocol_v0
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

// get the version of the protocol
func (p *ProtocolV0) GetVersion() uint64 {
	return p.version
}

// export the app state && validators
func (p *ProtocolV0) ExportAppStateAndValidators(ctx sdk.Context) (appState json.RawMessage,
	validators []types.GenesisValidator, err error) {
	// get all exsiting accounts
	var exportedAccounts []ExportedAccount
	appendAccFn := func(acc exported.Account) (stop bool) {
		exportedAcc := NewExportedAccount(acc)
		exportedAccounts = append(exportedAccounts, exportedAcc)
		return
	}

	p.accountKeeper.IterateAccounts(ctx, appendAccFn)
	// make state 2 export
	exportedState := ExportState{
		Accounts: exportedAccounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
		GovData:  gov.DefaultGenesisState(),
	}

	if appState, err = codec.MarshalJSONIndent(p.cdc, exportedState); err != nil {
		return nil, nil, err
	}

	return
}

// load the protocol 2 the app
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

// nothing 2 do
func (p *ProtocolV0) Init() {}

// get tx codec
func (p *ProtocolV0) GetCodec() *codec.Codec {
	if p.cdc == nil {
		panic("Invalid cdc from ProtocolV0")
	}
	return p.cdc
}

//  gracefully stop OKChain
func (p *ProtocolV0) CheckStopped() {
	if p.stopped {
		p.logger.Info("OKChain is going to exit")
		server.Stop()
		p.logger.Info("OKChain was stopped")
		select {}
	}
}

// get backend keeper
func (p *ProtocolV0) GetBackendKeeper() backend.Keeper {
	return p.backendKeeper
}

// get stream keeper
func (p *ProtocolV0) GetStreamKeeper() stream.Keeper {
	return p.streamKeeper
}

// get crisis keeper
func (p *ProtocolV0) GetCrisisKeeper() crisis.Keeper {
	return p.crisisKeeper
}

// get staking keeper
func (p *ProtocolV0) GetStakingKeeper() staking.Keeper {
	return p.stakingKeeper
}

// get distr keeper
func (p *ProtocolV0) GetDistrKeeper() distr.Keeper {
	return p.distrKeeper
}

// get slashing keeper
func (p *ProtocolV0) GetSlashingKeeper() slashing.Keeper {
	return p.slashingKeeper
}

// get token keeper
func (p *ProtocolV0) GetTokenKeeper() token.Keeper {
	return p.tokenKeeper
}

// get the map of KVStoreKeys
func (p *ProtocolV0) GetKVStoreKeysMap() map[string]*sdk.KVStoreKey {
	return p.keys
}

// get the map of TransientStoreKeys
func (p *ProtocolV0) GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey {
	return p.tkeys
}

// create the main codec of each module 2 ProtocolV0
func (p *ProtocolV0) setCodec() {
	p.cdc = MakeCodec()
}

// create all keepers declared in the ProtocolV0 struct
func (p *ProtocolV0) produceKeepers() {
	// get config
	appConfig, err := config.ParseConfig()
	if err != nil {
		p.logger.Error(fmt.Sprintf("the config of OKChain was parsed error : %s", err.Error()))
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

	// 2.add keepers
	p.accountKeeper = auth.NewAccountKeeper(p.cdc, p.keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	p.bankKeeper = bank.NewBaseKeeper(p.accountKeeper, bankSubspace, bank.DefaultCodespace, p.moduleAccountAddrs())
	p.paramsKeeper.SetBankKeeper(p.bankKeeper)
	p.supplyKeeper = supply.NewKeeper(p.cdc, p.keys[supply.StoreKey], p.accountKeeper, p.bankKeeper, maccPerms)
	// rollback to cosmos staking module
	//stakingKeeper := staking.NewKeeper(
	//	p.cdc, p.keys[staking.StoreKey], p.keys[staking.DelegatorPoolKey], p.keys[staking.RedelegationKeyM],
	//	p.keys[staking.RedelegationActonKey], p.keys[staking.UnbondingKey], p.tkeys[staking.TStoreKey],
	//	p.supplyKeeper, stakingSubspace, staking.DefaultCodespace,
	//)
	stakingKeeper := staking.NewKeeper(p.cdc, p.keys[staking.StoreKey], p.tkeys[staking.TStoreKey],
		p.supplyKeeper, stakingSubspace, staking.DefaultCodespace)

	p.paramsKeeper.SetStakingKeeper(stakingKeeper)
	p.mintKeeper = mint.NewKeeper(
		p.cdc, p.keys[mint.StoreKey], mintSubspace, &stakingKeeper, p.supplyKeeper, auth.FeeCollectorName,
	)

	// rollback to cosmos distr module
	//p.distrKeeper = distr.NewKeeper(p.cdc, p.keys[distr.StoreKey],
	//	p.keys[distr.ValidatorsSnapshotKey], p.keys[distr.DelegationSnapshotKey],
	//	distrSubspace, &stakingKeeper, p.supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName)
	p.distrKeeper = distr.NewKeeper(p.cdc, p.keys[distr.StoreKey],
		distrSubspace, &stakingKeeper, p.supplyKeeper,
		distr.DefaultCodespace, auth.FeeCollectorName, p.moduleAccountAddrs(),
	)

	p.slashingKeeper = slashing.NewKeeper(
		p.cdc, p.keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)

	p.crisisKeeper = crisis.NewKeeper(crisisSubspace, p.invCheckPeriod, p.supplyKeeper, auth.FeeCollectorName)

	p.tokenKeeper = token.NewKeeper(
		p.bankKeeper, p.paramsKeeper, tokenSubspace, auth.FeeCollectorName, p.supplyKeeper,
		p.keys[token.StoreKey], p.keys[token.KeyFreeze], p.keys[token.KeyLock],
		p.cdc, appConfig.BackendConfig.EnableBackend, false,
	)

	p.dexKeeper = dex.NewKeeper(auth.FeeCollectorName, p.supplyKeeper, dexSubspace, p.tokenKeeper, &stakingKeeper, p.bankKeeper,
		p.keys[dex.StoreKey], p.keys[dex.TokenPairStoreKey], p.cdc)

	p.orderKeeper = order.NewKeeper(
		p.tokenKeeper, p.supplyKeeper, p.paramsKeeper, p.dexKeeper, orderSubspace, auth.FeeCollectorName,
		p.keys[order.OrderStoreKey],
		p.cdc, appConfig.BackendConfig.EnableBackend, orderMetrics,
	)

	p.streamKeeper = stream.NewKeeper(p.orderKeeper, p.tokenKeeper, p.dexKeeper, p.accountKeeper, p.cdc, p.logger,
		appConfig.StreamConfig, streamMetrics)

	p.backendKeeper = backend.NewKeeper(p.orderKeeper, p.tokenKeeper, p.dexKeeper, p.streamKeeper.GetMarketKeeper(),
		p.cdc, p.logger, appConfig.BackendConfig)

	// 3.register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(&p.paramsKeeper)).
		//AddRoute(token.RouterKey, token.NewDexProposalHandler(&p.tokenKeeper)).
		AddRoute(dex.RouterKey, dex.NewProposalHandler(&p.dexKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewAppUpgradeProposalHandler(&p.upgradeKeeper))
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
	//p.tokenKeeper.SetGovKeeper(p.govKeeper)
	p.dexKeeper.SetGovKeeper(p.govKeeper)
	// 4.register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	p.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(p.distrKeeper.Hooks(), p.slashingKeeper.Hooks()),
	)

	p.upgradeKeeper = upgrade.NewKeeper(
		p.cdc, p.keys[upgrade.StoreKey], p.protocolKeeper, p.stakingKeeper, p.bankKeeper, upgradeSubspace,
	)
	// 5.for test
	p.debugKeeper = debug.NewDebugKeeper(
		p.cdc, p.keys[debug.StoreKey],
		p.orderKeeper,
		p.stakingKeeper, p.tokenKeeper, p.supplyKeeper, auth.FeeCollectorName, p.Stop,
	)

}

// return all the module account addresses
func (p *ProtocolV0) moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[p.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// set module.Manager in protocolV0
func (p *ProtocolV0) setManager() {

	// DEMO Manager
	p.mm = module.NewManager(
		genaccounts.NewAppModule(p.accountKeeper),
		genutil.NewAppModule(p.accountKeeper, p.stakingKeeper, p.parent.DeliverTx),
		auth.NewAppModule(p.accountKeeper),
		bank.NewAppModule(p.bankKeeper, p.accountKeeper),
		crisis.NewAppModule(&p.crisisKeeper),
		supply.NewAppModule(p.supplyKeeper, p.accountKeeper),
		params.NewAppModule(p.paramsKeeper),
		mint.NewAppModule(p.mintKeeper), // mining disabled
		slashing.NewAppModule(p.slashingKeeper, p.stakingKeeper),

		// rollback to staking/distr module
		staking.NewAppModule(p.stakingKeeper, p.distrKeeper, p.accountKeeper, p.supplyKeeper),

		// rollback to cosmos distr module
		//distr.NewAppModule(version.ProtocolVersionV0, p.distrKeeper, p.supplyKeeper),
		distr.NewAppModule(p.distrKeeper, p.supplyKeeper),
		gov.NewAppModule(version.ProtocolVersionV0, p.govKeeper, p.supplyKeeper),
		order.NewAppModule(version.ProtocolVersionV0, p.orderKeeper, p.supplyKeeper),
		// okchain extended
		token.NewAppModule(version.ProtocolVersionV0, p.tokenKeeper, p.supplyKeeper),

		// TODO
		dex.NewAppModule(version.ProtocolVersionV0, p.dexKeeper, p.supplyKeeper),

		backend.NewAppModule(p.backendKeeper),
		stream.NewAppModule(p.streamKeeper),
		upgrade.NewAppModule(p.upgradeKeeper),
		/******** for test **********/
		debug.NewAppModule(p.debugKeeper),
		/****************************/
	)

	// ORDER SETTING
	p.mm.SetOrderBeginBlockers(
		/******** for test **********/
		debug.ModuleName,
		/****************************/
		order.ModuleName,
		token.ModuleName,
		dex.ModuleName,
		mint.ModuleName,
		distr.ModuleName,
		slashing.ModuleName,
		staking.ModuleName,
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
		order.ModuleName,
		token.ModuleName,
		dex.ModuleName,
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		supply.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName,
		params.ModuleName,
		upgrade.ModuleName,
	)
}

// register Routers by Manager
func (p *ProtocolV0) registerRouters() {
	p.mm.RegisterInvariants(&p.crisisKeeper)
	p.mm.RegisterRoutes(p.router, p.queryRouter)
	p.parent.SetRouter(p.router, p.queryRouter)
}

// set ante handler
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

// InitChainer(hook) initializes application state at genesis
func (p *ProtocolV0) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	p.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	var accGenesisState genaccounts.GenesisState
	p.cdc.MustUnmarshalJSON(genesisState[genaccounts.ModuleName], &accGenesisState)

	var acc auth.Account
	if len(accGenesisState) > 0 {
		acc = accGenesisState[0].ToAccount()
	}

	if err := token.IssueOKT(ctx, p.tokenKeeper, genesisState[token.ModuleName], acc); err != nil {
		panic(err)
	}
	return p.mm.InitGenesis(ctx, genesisState)

}

// BeginBlocker(hook) set function 2 BaseApp
func (p *ProtocolV0) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return p.mm.BeginBlock(ctx, req)
}

// EndBlocker(hook) set function 2 BaseApp
func (p *ProtocolV0) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return p.mm.EndBlock(ctx, req)
}

// exit gracefully
func (p *ProtocolV0) Stop() {
	p.logger.Info(fmt.Sprintf("[%s]%s", utils.GoId, "OKChain stops notification."))
	p.stopped = true
}

// make codec from all the modules
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc
}

func validateMsgHook(orderKeeper order.Keeper) auth.ValidateMsgHandler {

	return func(newCtx sdk.Context, msgs []sdk.Msg) sdk.Result {

		for _, msg := range msgs {

			switch msg.(type) {
			case order.MsgNewOrders:
				return order.ValidateMsgNewOrders(newCtx, orderKeeper, msg.(order.MsgNewOrders))
			case order.MsgCancelOrders:
				return order.ValidateMsgCancelOrders(newCtx, orderKeeper, msg.(order.MsgCancelOrders))
			}

		}
		return sdk.Result{}
	}
}

func isSystemFreeHook(ctx sdk.Context, msgs []sdk.Msg) bool {
	if ctx.BlockHeight() < 1 {
		return true
	}

	for _, msg := range msgs {
		switch msg.(type) {
		case order.MsgNewOrders, order.MsgCancelOrders:
		case dex.MsgList, dex.MsgDelist, dex.MsgTransferOwnership:
			return true
		default:
			return false
		}
	}

	return true
}

// export
func (p *ProtocolV0) ExportGenesis(ctx sdk.Context) map[string]json.RawMessage {
	return p.mm.ExportGenesis(ctx)
}

// set logger
func (p *ProtocolV0) SetLogger(log log.Logger) Protocol {
	p.logger = log
	return p
}

// set parent implement
func (p *ProtocolV0) SetParent(parent Parent) Protocol {
	p.parent = parent
	return p
}

func (p *ProtocolV0) GetParent() Parent {
	if p.parent == nil {
		panic("parent is nil in protocol")
	}
	return p.parent
}
