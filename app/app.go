package app

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"

	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"

	"io"
	"math/big"
	"os"
	"sort"
	"sync"

	"github.com/okex/exchain/app/ante"
	okexchaincodec "github.com/okex/exchain/app/codec"
	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/refund"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/app/utils/sanity"
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	upgradetypes "github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	capabilityModule "github.com/okex/exchain/libs/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/crisis"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	govclient "github.com/okex/exchain/libs/cosmos-sdk/x/mint/client"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	ibctransferkeeper "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	ibcclient "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmos "github.com/okex/exchain/libs/tendermint/libs/os"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/common/analyzer"
	commonversion "github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/dex"
	dexclient "github.com/okex/exchain/x/dex/client"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/erc20"
	erc20client "github.com/okex/exchain/x/erc20/client"
	"github.com/okex/exchain/x/evidence"
	"github.com/okex/exchain/x/evm"
	evmclient "github.com/okex/exchain/x/evm/client"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/farm"
	farmclient "github.com/okex/exchain/x/farm/client"
	"github.com/okex/exchain/x/genutil"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/gov/keeper"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/params"
	paramsclient "github.com/okex/exchain/x/params/client"
	"github.com/okex/exchain/x/slashing"
	"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/token"

	"github.com/spf13/viper"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

func init() {
	// set the address prefixes
	config := sdk.GetConfig()
	okexchain.SetBech32Prefixes(config)
	okexchain.SetBip44CoinType(config)
}

const (
	appName = "OKExChain"
)

var (
	// DefaultCLIHome sets the default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.exchaincli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.exchaind")

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		supply.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distr.ProposalHandler,
			dexclient.DelistProposalHandler, farmclient.ManageWhiteListProposalHandler,
			evmclient.ManageContractDeploymentWhitelistProposalHandler,
			evmclient.ManageContractBlockedListProposalHandler,
			evmclient.ManageContractMethodBlockedListProposalHandler,
			govclient.ManageTreasuresProposalHandler,
			erc20client.TokenMappingProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		evidence.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evm.AppModuleBasic{},
		token.AppModuleBasic{},
		dex.AppModuleBasic{},
		order.AppModuleBasic{},
		ammswap.AppModuleBasic{},
		farm.AppModuleBasic{},
		capabilityModule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},
		erc20.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:       nil,
		distr.ModuleName:            nil,
		mint.ModuleName:             {supply.Minter},
		staking.BondedPoolName:      {supply.Burner, supply.Staking},
		staking.NotBondedPoolName:   {supply.Burner, supply.Staking},
		gov.ModuleName:              nil,
		token.ModuleName:            {supply.Minter, supply.Burner},
		dex.ModuleName:              nil,
		order.ModuleName:            nil,
		ammswap.ModuleName:          {supply.Minter, supply.Burner},
		farm.ModuleName:             nil,
		farm.YieldFarmingAccount:    nil,
		farm.MintFarmingAccount:     {supply.Burner},
		ibctransfertypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		erc20.ModuleName:            {authtypes.Minter, authtypes.Burner},
	}

	GlobalGpIndex = GasPriceIndex{}

	onceLog sync.Once
)

var _ simapp.App = (*OKExChainApp)(nil)

// OKExChainApp implements an extended ABCI application. It is an application
// that may process transactions through Ethereum's EVM running atop of
// Tendermint consensus.
type OKExChainApp struct {
	*bam.BaseApp

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	MintKeeper     mint.Keeper
	DistrKeeper    distr.Keeper
	GovKeeper      gov.Keeper
	CrisisKeeper   crisis.Keeper
	UpgradeKeeper  upgrade.Keeper
	ParamsKeeper   params.Keeper
	EvidenceKeeper evidence.Keeper
	EvmKeeper      *evm.Keeper
	TokenKeeper    token.Keeper
	DexKeeper      dex.Keeper
	OrderKeeper    order.Keeper
	SwapKeeper     ammswap.Keeper
	FarmKeeper     farm.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	blockGasPrice []*big.Int

	configurator module.Configurator
	// ibc
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper  capabilitykeeper.ScopedKeeper
	TransferKeeper       ibctransferkeeper.Keeper
	CapabilityKeeper     *capabilitykeeper.Keeper
	IBCKeeper            *ibc.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	marshal              *codec.CodecProxy
	heightTasks          map[int64]*upgradetypes.HeightTasks
	Erc20Keeper          erc20.Keeper
}

// NewOKExChainApp returns a reference to a new initialized OKExChain application.
func NewOKExChainApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint,
	baseAppOptions ...func(*bam.BaseApp),
) *OKExChainApp {
	logger.Info("Starting OEC",
		"GenesisHeight", tmtypes.GetStartBlockHeight(),
		"MercuryHeight", tmtypes.GetMercuryHeight(),
		"VenusHeight", tmtypes.GetVenusHeight(),
	)
	onceLog.Do(func() {
		iavl.SetLogger(logger.With("module", "iavl"))
		logStartingFlags(logger)
	})

	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)

	// NOTE we use custom OKExChain transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := bam.NewBaseApp(appName, logger, db, evm.TxDecoder(cdc), baseAppOptions...)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	bApp.SetStartLogHandler(analyzer.StartTxLog)
	bApp.SetEndLogHandler(analyzer.StopTxLog)

	bApp.SetInterfaceRegistry(interfaceReg)

	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	bApp.SetTxDecoder(func(txBytes []byte, height ...int64) (ret sdk.Tx, err error) {
		return evm.TxDecoder(cdc, codecProxy)(txBytes, height...)
	})

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, upgrade.StoreKey, evidence.StoreKey,
		evm.StoreKey, token.StoreKey, token.KeyLock, dex.StoreKey, dex.TokenPairStoreKey,
		order.OrderStoreKey, ammswap.StoreKey, farm.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		host.StoreKey,
		erc20.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &OKExChainApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
		heightTasks:    make(map[int64]*upgradetypes.HeightTasks),
	}
	bApp.SetInterceptors(makeInterceptors(cdc))

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace)
	app.subspaces[crisis.ModuleName] = app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.ParamsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[evm.ModuleName] = app.ParamsKeeper.Subspace(evm.DefaultParamspace)
	app.subspaces[token.ModuleName] = app.ParamsKeeper.Subspace(token.DefaultParamspace)
	app.subspaces[dex.ModuleName] = app.ParamsKeeper.Subspace(dex.DefaultParamspace)
	app.subspaces[order.ModuleName] = app.ParamsKeeper.Subspace(order.DefaultParamspace)
	app.subspaces[ammswap.ModuleName] = app.ParamsKeeper.Subspace(ammswap.DefaultParamspace)
	app.subspaces[farm.ModuleName] = app.ParamsKeeper.Subspace(farm.DefaultParamspace)
	app.subspaces[host.ModuleName] = app.ParamsKeeper.Subspace(host.ModuleName)
	app.subspaces[ibctransfertypes.ModuleName] = app.ParamsKeeper.Subspace(ibctransfertypes.ModuleName)
	app.subspaces[erc20.ModuleName] = app.ParamsKeeper.Subspace(erc20.DefaultParamspace)

	//proxy := codec.NewMarshalProxy(cc, cdc)
	app.marshal = codecProxy
	// use custom OKExChain account for contracts
	app.AccountKeeper = auth.NewAccountKeeper(
		cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], okexchain.ProtoAccount,
	)

	bankKeeper := bank.NewBaseKeeperWithMarshal(
		&app.AccountKeeper, codecProxy, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs(),
	)
	app.BankKeeper = &bankKeeper
	app.ParamsKeeper.SetBankKeeper(app.BankKeeper)
	app.SupplyKeeper = supply.NewKeeper(
		cdc, keys[supply.StoreKey], &app.AccountKeeper, app.BankKeeper, maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		codecProxy, keys[staking.StoreKey], app.SupplyKeeper, app.subspaces[staking.ModuleName],
	)
	app.ParamsKeeper.SetStakingKeeper(stakingKeeper)
	app.MintKeeper = mint.NewKeeper(
		cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper,
		app.SupplyKeeper, auth.FeeCollectorName, farm.MintFarmingAccount,
	)
	app.DistrKeeper = distr.NewKeeper(
		cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashing.NewKeeper(
		cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.CrisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName,
	)
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.marshal.GetCdc())
	app.ParamsKeeper.RegisterSignal(evmtypes.SetEvmParamsNeedUpdate)
	app.EvmKeeper = evm.NewKeeper(
		app.marshal.GetCdc(), keys[evm.StoreKey], app.subspaces[evm.ModuleName], &app.AccountKeeper, app.SupplyKeeper, app.BankKeeper, logger)
	(&bankKeeper).SetInnerTxKeeper(app.EvmKeeper)

	app.TokenKeeper = token.NewKeeper(app.BankKeeper, app.subspaces[token.ModuleName], auth.FeeCollectorName, app.SupplyKeeper,
		keys[token.StoreKey], keys[token.KeyLock], app.marshal.GetCdc(), false, &app.AccountKeeper)

	app.DexKeeper = dex.NewKeeper(auth.FeeCollectorName, app.SupplyKeeper, app.subspaces[dex.ModuleName], app.TokenKeeper, &stakingKeeper,
		app.BankKeeper, app.keys[dex.StoreKey], app.keys[dex.TokenPairStoreKey], app.marshal.GetCdc())

	app.OrderKeeper = order.NewKeeper(
		app.TokenKeeper, app.SupplyKeeper, app.DexKeeper, app.subspaces[order.ModuleName], auth.FeeCollectorName,
		app.keys[order.OrderStoreKey], app.marshal.GetCdc(), false, orderMetrics)

	app.SwapKeeper = ammswap.NewKeeper(app.SupplyKeeper, app.TokenKeeper, app.marshal.GetCdc(), app.keys[ammswap.StoreKey], app.subspaces[ammswap.ModuleName])

	app.FarmKeeper = farm.NewKeeper(auth.FeeCollectorName, app.SupplyKeeper, app.TokenKeeper, app.SwapKeeper, *app.EvmKeeper, app.subspaces[farm.StoreKey],
		app.keys[farm.StoreKey], app.marshal.GetCdc())

	// create evidence keeper with router
	evidenceKeeper := evidence.NewKeeper(
		cdc, keys[evidence.StoreKey], app.subspaces[evidence.ModuleName], &app.StakingKeeper, app.SlashingKeeper,
	)
	evidenceRouter := evidence.NewRouter()
	evidenceKeeper.SetRouter(evidenceRouter)
	app.EvidenceKeeper = *evidenceKeeper

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(codecProxy, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(host.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	scopedIBCMockKeeper := app.CapabilityKeeper.ScopeToModule("mock")

	app.IBCKeeper = ibc.NewKeeper(
		codecProxy, keys[host.StoreKey], app.GetSubspace(host.ModuleName), &stakingKeeper, &scopedIBCKeeper, interfaceReg,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		codecProxy, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.SupplyKeeper, app.SupplyKeeper, scopedTransferKeeper, interfaceReg,
	)
	ibctransfertypes.SetMarshal(codecProxy)

	app.Erc20Keeper = erc20.NewKeeper(app.marshal.GetCdc(), app.keys[erc20.ModuleName], app.subspaces[erc20.ModuleName],
		app.AccountKeeper, app.SupplyKeeper, app.BankKeeper, app.EvmKeeper, app.TransferKeeper)

	// register the proposal types
	// 3.register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(&app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(dex.RouterKey, dex.NewProposalHandler(&app.DexKeeper)).
		AddRoute(farm.RouterKey, farm.NewManageWhiteListProposalHandler(&app.FarmKeeper)).
		AddRoute(evm.RouterKey, evm.NewManageContractDeploymentWhitelistProposalHandler(app.EvmKeeper)).
		AddRoute(mint.RouterKey, mint.NewManageTreasuresProposalHandler(&app.MintKeeper)).
		AddRoute(host.RouterKey, ibcclient.NewClientUpdateProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(erc20.RouterKey, erc20.NewProposalHandler(&app.Erc20Keeper))
	govProposalHandlerRouter := keeper.NewProposalHandlerRouter()
	govProposalHandlerRouter.AddRoute(params.RouterKey, &app.ParamsKeeper).
		AddRoute(dex.RouterKey, &app.DexKeeper).
		AddRoute(farm.RouterKey, &app.FarmKeeper).
		AddRoute(evm.RouterKey, app.EvmKeeper).
		AddRoute(mint.RouterKey, &app.MintKeeper).
		AddRoute(erc20.RouterKey, &app.Erc20Keeper)
	app.GovKeeper = gov.NewKeeper(
		app.marshal.GetCdc(), app.keys[gov.StoreKey], app.ParamsKeeper, app.subspaces[gov.DefaultParamspace],
		app.SupplyKeeper, &stakingKeeper, gov.DefaultParamspace, govRouter,
		app.BankKeeper, govProposalHandlerRouter, auth.FeeCollectorName,
	)
	app.ParamsKeeper.SetGovKeeper(app.GovKeeper)
	app.DexKeeper.SetGovKeeper(app.GovKeeper)
	app.FarmKeeper.SetGovKeeper(app.GovKeeper)
	app.EvmKeeper.SetGovKeeper(app.GovKeeper)
	app.MintKeeper.SetGovKeeper(app.GovKeeper)
	app.Erc20Keeper.SetGovKeeper(app.GovKeeper)

	// Set EVM hooks
	app.EvmKeeper.SetHooks(evm.NewLogProcessEvmHook(erc20.NewSendToIbcEventHandler(app.Erc20Keeper)))
	// Set IBC hooks
	app.TransferKeeper = *app.TransferKeeper.SetHooks(erc20.NewIBCTransferHooks(app.Erc20Keeper))
	transferModule := transfer.NewAppModule(app.TransferKeeper, codecProxy)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	//ibcRouter.AddRoute(ibcmock.ModuleName, mockModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper),
		distr.NewAppModule(app.DistrKeeper, app.SupplyKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		evm.NewAppModule(app.EvmKeeper, &app.AccountKeeper),
		token.NewAppModule(commonversion.ProtocolVersionV0, app.TokenKeeper, app.SupplyKeeper),
		dex.NewAppModule(commonversion.ProtocolVersionV0, app.DexKeeper, app.SupplyKeeper),
		order.NewAppModule(commonversion.ProtocolVersionV0, app.OrderKeeper, app.SupplyKeeper),
		ammswap.NewAppModule(app.SwapKeeper),
		farm.NewAppModule(app.FarmKeeper),
		params.NewAppModule(app.ParamsKeeper),
		// ibc
		ibc.NewAppModule(app.IBCKeeper),
		capabilityModule.NewAppModule(codecProxy, *app.CapabilityKeeper),
		transferModule,
		erc20.NewAppModule(app.Erc20Keeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		bank.ModuleName,
		capabilitytypes.ModuleName,
		order.ModuleName,
		token.ModuleName,
		dex.ModuleName,
		mint.ModuleName,
		distr.ModuleName,
		slashing.ModuleName,
		staking.ModuleName,
		farm.ModuleName,
		evidence.ModuleName,
		evm.ModuleName,
		host.ModuleName,
		ibctransfertypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisis.ModuleName,
		gov.ModuleName,
		dex.ModuleName,
		order.ModuleName,
		staking.ModuleName,
		evm.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, supply.ModuleName,
		token.ModuleName, dex.ModuleName, order.ModuleName, ammswap.ModuleName, farm.ModuleName,
		ibctransfertypes.ModuleName,
		host.ModuleName,
		evm.ModuleName, crisis.ModuleName, genutil.ModuleName, params.ModuleName, evidence.ModuleName,
		erc20.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	app.configurator = module.NewConfigurator(app.Codec(), app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)
	app.setupUpgradeModules()

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		distr.NewAppModule(app.DistrKeeper, app.SupplyKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper), // NOTE: only used for simulation to generate randomized param change proposals
		ibc.NewAppModule(app.IBCKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.AccountKeeper, app.EvmKeeper, app.SupplyKeeper, validateMsgHook(app.OrderKeeper)))
	app.SetEndBlocker(app.EndBlocker)
	app.SetGasRefundHandler(refund.NewGasRefundHandler(app.AccountKeeper, app.SupplyKeeper))
	app.SetAccHandler(NewAccHandler(app.AccountKeeper))
	app.SetParallelTxHandlers(updateFeeCollectorHandler(app.BankKeeper, app.SupplyKeeper), evmTxFeeHandler(), fixLogForParallelTxHandler(app.EvmKeeper))

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	app.ScopedIBCMockKeeper = scopedIBCMockKeeper

	return app
}

func (app *OKExChainApp) SetOption(req abci.RequestSetOption) (res abci.ResponseSetOption) {
	if req.Key == "CheckChainID" {
		if err := okexchain.IsValidateChainIdWithGenesisHeight(req.Value); err != nil {
			app.Logger().Error(err.Error())
			panic(err)
		}
		err := okexchain.SetChainId(req.Value)
		if err != nil {
			app.Logger().Error(err.Error())
			panic(err)
		}
	}
	return app.BaseApp.SetOption(req)
}

func (app *OKExChainApp) LoadStartVersion(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// Name returns the name of the App
func (app *OKExChainApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker updates every begin block
func (app *OKExChainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker updates every end block
func (app *OKExChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	if appconfig.GetOecConfig().GetEnableDynamicGp() {
		GlobalGpIndex = CalBlockGasPriceIndex(app.blockGasPrice, appconfig.GetOecConfig().GetDynamicGpWeight())
		app.blockGasPrice = app.blockGasPrice[:0]
	}

	return app.mm.EndBlock(ctx, req)
}

// InitChainer updates at chain initialization
func (app *OKExChainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

	var genesisState simapp.GenesisState
	app.marshal.GetCdc().MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads state at a particular height
func (app *OKExChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *OKExChainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// SimulationManager implements the SimulationApp interface
func (app *OKExChainApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OKExChainApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// Codec returns OKExChain's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OKExChainApp) Codec() *codec.Codec {
	return app.marshal.GetCdc()
}

func (app *OKExChainApp) Marshal() *codec.CodecProxy {
	return app.marshal
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OKExChainApp) GetSubspace(moduleName string) params.Subspace {
	return app.subspaces[moduleName]
}
func (app *OKExChainApp) setupUpgradeModules() {
	heightTasks, pip, pR := app.CollectUpgradeModules(app.mm)
	app.heightTasks = heightTasks
	if pip != nil {
		app.GetCMS().SetPruneHeightFilterPipeline(pip)
		app.GetCMS().SetCommitHeightFilterPipeline(pip)
	}
	vs := app.subspaces
	for k, vv := range pR {
		supace, exist := vs[k]
		if !exist {
			continue
		}
		vs[k] = supace.LazyWithKeyTable(subspace.NewKeyTable(vv.ParamSetPairs()...))
	}
}

var protoCodec = encoding.GetCodec(proto.Name)

func makeInterceptors(cdc *codec.Codec) map[string]bam.Interceptor {
	m := make(map[string]bam.Interceptor)
	m["/cosmos.tx.v1beta1.Service/Simulate"] = bam.NewFunctionInterceptor(func(req *abci.RequestQuery) error {
		req.Path = "app/simulate"
		return nil
	}, func(resp *abci.ResponseQuery) {})

	m["/cosmos.bank.v1beta1.Query/AllBalances"] = bam.NewFunctionInterceptor(func(req *abci.RequestQuery) error {
		req.Path = "custom/bank/grpc_balances"
		return nil
	}, func(resp *abci.ResponseQuery) {})

	m["/cosmos.staking.v1beta1.Query/Params"] = bam.NewFunctionInterceptor(func(req *abci.RequestQuery) error {
		req.Path = "custom/staking/parameters"
		return nil
	}, func(resp *abci.ResponseQuery) {})

	return m
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

func validateMsgHook(orderKeeper order.Keeper) ante.ValidateMsgHandler {
	return func(newCtx sdk.Context, msgs []sdk.Msg) error {

		wrongMsgErr := sdk.ErrUnknownRequest(
			"It is not allowed that a transaction with more than one message contains order or evm message")
		var err error

		for _, msg := range msgs {
			switch assertedMsg := msg.(type) {
			case order.MsgNewOrders:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
				_, err = order.ValidateMsgNewOrders(newCtx, orderKeeper, assertedMsg)
			case order.MsgCancelOrders:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
				err = order.ValidateMsgCancelOrders(newCtx, orderKeeper, assertedMsg)
			case *evmtypes.MsgEthereumTx:
				if len(msgs) > 1 {
					return wrongMsgErr
				}
			}

			if err != nil {
				return err
			}
		}
		return nil
	}
}

func NewAccHandler(ak auth.AccountKeeper) sdk.AccHandler {
	return func(
		ctx sdk.Context, addr sdk.AccAddress,
	) uint64 {
		return ak.GetAccount(ctx, addr).GetSequence()
	}
}

func PreRun(ctx *server.Context) error {
	// set the dynamic config
	appconfig.RegisterDynamicConfig(ctx.Logger.With("module", "config"))

	// check start flag conflicts
	err := sanity.CheckStart()
	if err != nil {
		return err
	}

	// set config by node mode
	setNodeConfig(ctx)

	//download pprof
	appconfig.PprofDownload(ctx)

	// pruning options
	_, err = server.GetPruningOptionsFromFlags()
	if err != nil {
		return err
	}
	// repair state on start
	if viper.GetBool(FlagEnableRepairState) {
		repairStateOnStart(ctx)
	}

	// init tx signature cache
	tmtypes.InitSignatureCache()
	return nil
}

func (o *OKExChainApp) CollectUpgradeModules(m *module.Manager) (map[int64]*upgradetypes.HeightTasks, types.HeightFilterPipeline, map[string]params.ParamSet) {
	hm := make(map[int64]*upgradetypes.HeightTasks)
	hStoreInfoModule := make(map[int64]map[string]struct{})
	paramsRet := make(map[string]params.ParamSet)
	for _, mm := range m.Modules {
		if ada, ok := mm.(upgradetypes.UpgradeModule); ok {
			set := ada.RegisterParam()
			if set != nil {
				if _, exist := paramsRet[ada.ModuleName()]; !exist {
					paramsRet[ada.ModuleName()] = set
				}
			}
			h := ada.UpgradeHeight()
			if h <= 0 {
				continue
			}
			t := ada.RegisterTask()
			if t == nil {
				continue
			}
			if err := t.ValidateBasic(); nil != err {
				panic(err)
			}
			storeInfoModule := hStoreInfoModule[h]
			if storeInfoModule == nil {
				storeInfoModule = make(map[string]struct{})
				hStoreInfoModule[h] = storeInfoModule
			}
			names := ada.BlockStoreModules()
			for _, n := range names {
				storeInfoModule[n] = struct{}{}
			}
			taskList := hm[h]
			if taskList == nil {
				v := make(upgradetypes.HeightTasks, 0)
				taskList = &v
				hm[h] = taskList
			}
			*taskList = append(*taskList, t)
		}
	}
	for _, v := range hm {
		sort.Sort(*v)
	}
	var pip types.HeightFilterPipeline
	for height, mm := range hStoreInfoModule {
		f := func(h int64) func(str string) bool {
			if h >= height {
				// next
				return nil
			}
			// filter block module
			return func(str string) bool {
				_, exist := mm[str]
				return exist
			}
		}
		if pip == nil {
			pip = f
		} else {
			pip = types.LinkPipeline(f, pip)
		}
	}
	return hm, pip, paramsRet
}
