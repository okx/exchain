package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/okex/exchain/x/wasm/keeper/testdata"

	okexchaincodec "github.com/okex/exchain/app/codec"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	cosmoscryptocodec "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	authkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/crisis"
	crisistypes "github.com/okex/exchain/libs/cosmos-sdk/x/crisis"
	"github.com/okex/exchain/libs/cosmos-sdk/x/evidence"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	"github.com/okex/exchain/libs/cosmos-sdk/x/slashing"
	slashingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/slashing"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	"github.com/okex/exchain/x/ammswap"
	dex "github.com/okex/exchain/x/dex/types"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/erc20"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/staking"
	token "github.com/okex/exchain/x/token/types"

	//upgradeclient "github.com/okex/exchain/libs/cosmos-sdk/x/upgrade/client"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	ibctransfertypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	ibchost "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	ibckeeper "github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	tmproto "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/distribution"
	distrclient "github.com/okex/exchain/x/distribution/client"
	distributionkeeper "github.com/okex/exchain/x/distribution/keeper"
	distributiontypes "github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/gov"
	govkeeper "github.com/okex/exchain/x/gov/keeper"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/okex/exchain/x/params"
	paramproposal "github.com/okex/exchain/x/params"
	paramskeeper "github.com/okex/exchain/x/params"
	paramstypes "github.com/okex/exchain/x/params"
	paramsclient "github.com/okex/exchain/x/params/client"
	stakingkeeper "github.com/okex/exchain/x/staking/keeper"
	stakingtypes "github.com/okex/exchain/x/staking/types"
	"github.com/okex/exchain/x/wasm/keeper/wasmtesting"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/stretchr/testify/require"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry interfacetypes.InterfaceRegistry
	Marshaler         codec.CodecProxy
	TxConfig          client.TxConfig
	Amino             *codec.Codec
}

var moduleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	supply.AppModuleBasic{},
	gov.NewAppModuleBasic(
		paramsclient.ProposalHandler, distrclient.ChangeDistributionTypeProposalHandler,
		paramsclient.ProposalHandler, distrclient.CommunityPoolSpendProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
)

func MakeTestCodec(t testing.TB) codec.CodecProxy {
	return MakeEncodingConfig(t).Marshaler
}

func MakeEncodingConfig(_ testing.TB) EncodingConfig {
	codecProxy, interfaceReg := okexchaincodec.MakeCodecSuit(moduleBasics)
	txConfig := ibc_tx.NewTxConfig(codecProxy.GetProtocMarshal(), ibc_tx.DefaultSignModes)
	encodingConfig := EncodingConfig{
		InterfaceRegistry: interfaceReg,
		Marshaler:         *codecProxy,
		Amino:             codecProxy.GetCdc(),
		TxConfig:          txConfig}
	amino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	cosmoscryptocodec.PubKeyRegisterInterfaces(interfaceReg)
	// add wasmd types
	types.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)

	return encodingConfig
}

var TestingStakeParams = stakingtypes.Params{
	UnbondingTime:      100,
	MaxValidators:      10,
	Epoch:              10,
	MaxValsToAddShares: 1,
	MinDelegation:      sdk.OneDec(),
	MinSelfDelegation:  sdk.OneDec(),
	HistoricalEntries:  10,
}

type TestFaucet struct {
	t                testing.TB
	bankKeeper       bank.Keeper
	supplyKeeper     supply.Keeper
	sender           sdk.AccAddress
	balance          sdk.Coins
	minterModuleName string
}

func NewTestFaucet(t testing.TB, ctx sdk.Context, bankKeeper bank.Keeper, supplyKeeper supply.Keeper, minterModuleName string, initialAmount ...sdk.Coin) *TestFaucet {
	require.NotEmpty(t, initialAmount)
	r := &TestFaucet{t: t, bankKeeper: bankKeeper, minterModuleName: minterModuleName, supplyKeeper: supplyKeeper}
	_, _, addr := keyPubAddr()
	r.sender = addr
	r.Mint(ctx, addr, initialAmount...)
	r.balance = initialAmount
	return r
}

func (f *TestFaucet) Mint(parentCtx sdk.Context, addr sdk.AccAddress, amounts ...sdk.Coin) {
	require.NotEmpty(f.t, amounts)
	ctx := parentCtx.SetEventManager(sdk.NewEventManager()) // discard all faucet related events

	err := f.supplyKeeper.MintCoins(*ctx, f.minterModuleName, amounts)
	require.NoError(f.t, err)
	err = f.supplyKeeper.SendCoinsFromModuleToAccount(*ctx, f.minterModuleName, addr, amounts)
	require.NoError(f.t, err)
	f.balance = f.balance.Add(amounts...)
}

func (f *TestFaucet) Fund(parentCtx sdk.Context, receiver sdk.AccAddress, amounts ...sdk.Coin) {
	require.NotEmpty(f.t, amounts)
	// ensure faucet is always filled
	if !f.balance.IsAllGTE(amounts) {
		f.Mint(parentCtx, f.sender, amounts...)
	}
	ctx := parentCtx.SetEventManager(sdk.NewEventManager()) // discard all faucet related events
	err := f.bankKeeper.SendCoins(*ctx, f.sender, receiver, amounts)
	require.NoError(f.t, err)
	f.balance = f.balance.Sub(amounts)
}

func (f *TestFaucet) NewFundedAccount(ctx sdk.Context, amounts ...sdk.Coin) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	f.Fund(ctx, addr, amounts...)
	return addr
}

type TestKeepers struct {
	AccountKeeper  authkeeper.AccountKeeper
	supplyKeepr    supply.Keeper
	StakingKeeper  stakingkeeper.Keeper
	DistKeeper     distributionkeeper.Keeper
	BankKeeper     bank.Keeper
	GovKeeper      govkeeper.Keeper
	ContractKeeper types.ContractOpsKeeper
	WasmKeeper     *Keeper
	IBCKeeper      *ibckeeper.Keeper
	Router         *baseapp.Router
	EncodingConfig EncodingConfig
	Faucet         *TestFaucet
	MultiStore     sdk.CommitMultiStore
}

// CreateDefaultTestInput common settings for CreateTestInput
func CreateDefaultTestInput(t testing.TB) (sdk.Context, TestKeepers) {
	return CreateTestInput(t, false, "staking")
}

// CreateTestInput encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func CreateTestInput(t testing.TB, isCheckTx bool, supportedFeatures string, opts ...Option) (sdk.Context, TestKeepers) {
	// Load default wasm config
	return createTestInput(t, isCheckTx, supportedFeatures, types.DefaultWasmConfig(), dbm.NewMemDB(), opts...)
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func createTestInput(
	t testing.TB,
	isCheckTx bool,
	supportedFeatures string,
	wasmConfig types.WasmConfig,
	db dbm.DB,
	opts ...Option,
) (sdk.Context, TestKeepers) {
	tempDir := t.TempDir()
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, upgrade.StoreKey, evidence.StoreKey,
		evm.StoreKey, token.StoreKey, token.KeyLock, dex.StoreKey, dex.TokenPairStoreKey,
		order.OrderStoreKey, ammswap.StoreKey, farm.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		ibchost.StoreKey,
		erc20.StoreKey,
		mpt.StoreKey,
		types.StoreKey,
	)
	ms := store.NewCommitMultiStore(db)
	for _, v := range keys {
		ms.MountStoreWithDB(v, sdk.StoreTypeIAVL, db)
	}
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeTransient, db)
	}

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, sdk.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger())
	ctx = types.WithTXCounter(ctx, 0)

	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	for _, m := range []string{authtypes.ModuleName,
		bank.ModuleName,
		stakingtypes.ModuleName,
		mint.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		ibchost.ModuleName,
		govtypes.ModuleName,
		types.ModuleName,
	} {
		paramsKeeper.Subspace(m)
	}
	subspace := func(m string) paramstypes.Subspace {
		r, ok := paramsKeeper.GetSubspace(m)
		require.True(t, ok)
		return r
	}
	maccPerms := map[string][]string{ // module account permissions
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
		types.ModuleName:            nil,
	}
	accountKeeper := auth.NewAccountKeeper(legacyAmino, keys[authtypes.StoreKey], keys[mpt.StoreKey], subspace(authtypes.ModuleName), okexchain.ProtoAccount)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	bankKeeper := bank.NewBaseKeeperWithMarshal(
		&accountKeeper, &appCodec, subspace(bank.ModuleName), blockedAddrs,
	)
	bankKeeper.SetSendEnabled(ctx, true)

	supplyKeeper := supply.NewKeeper(
		legacyAmino, keys[supply.StoreKey], &accountKeeper, bank.NewBankKeeperAdapter(bankKeeper), maccPerms,
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		&appCodec,
		keys[stakingtypes.StoreKey],
		supplyKeeper,
		subspace(stakingtypes.ModuleName),
	)
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	distKeeper := distributionkeeper.NewKeeper(
		legacyAmino,
		keys[distributiontypes.StoreKey],
		subspace(distributiontypes.ModuleName),
		stakingKeeper,
		supplyKeeper,
		authtypes.FeeCollectorName,
		blockedAddrs,
	)
	distKeeper.SetParams(ctx, distributiontypes.DefaultParams())
	stakingKeeper.SetHooks(distKeeper.Hooks())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distributiontypes.InitialFeePool())

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		keys[upgradetypes.StoreKey],
		legacyAmino,
	)

	faucet := NewTestFaucet(t, ctx, bankKeeper, supplyKeeper, mint.ModuleName, sdk.NewCoin("stake", sdk.NewInt(100_000_000_000)))

	// set some funds ot pay out validatores, based on code from:
	// https://github.com/okex/exchain/libs/cosmos-sdk/blob/fea231556aee4d549d7551a6190389c4328194eb/x/distribution/keeper/keeper_test.go#L50-L57
	distrAcc := distKeeper.GetDistributionAccount(ctx)
	faucet.Fund(ctx, distrAcc.GetAddress(), sdk.NewCoin("stake", sdk.NewInt(2000000)))

	supplyKeeper.SetModuleAccount(ctx, distrAcc)

	capabilityKeeper := capabilitykeeper.NewKeeper(
		&appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedWasmKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)

	ibcKeeper := ibckeeper.NewKeeper(
		&appCodec,
		keys[ibchost.StoreKey],
		subspace(ibchost.ModuleName),
		stakingKeeper,
		upgradeKeeper,
		&scopedIBCKeeper,
		encodingConfig.InterfaceRegistry,
	)

	router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(bank.RouterKey, bh)
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(stakingtypes.RouterKey, sh)
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(distributiontypes.RouterKey, dh)

	querier := baseapp.NewGRPCQueryRouter()
	querier.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetAddressVerifier(types.VerifyAddressLen())

	keeper := NewKeeper(
		&appCodec,
		keys[types.StoreKey],
		subspace(types.ModuleName),
		&accountKeeper,
		bank.NewBankKeeperAdapter(bankKeeper),
		ibcKeeper.ChannelKeeper,
		&ibcKeeper.PortKeeper,
		scopedWasmKeeper,
		wasmtesting.MockIBCTransferKeeper{},
		msgRouter,
		querier,
		tempDir,
		wasmConfig,
		supportedFeatures,
		opts...,
	)
	keeper.SetParams(ctx, types.DefaultParams())
	// add wasm handler so we can loop-back (contracts calling contracts)
	contractKeeper := NewDefaultPermissionKeeper(&keeper)
	router.AddRoute(types.RouterKey, TestHandler(contractKeeper))

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bank.NewAppModule(bankKeeper, accountKeeper, supplyKeeper),
		staking.NewAppModule(stakingKeeper, accountKeeper, supplyKeeper),
		distribution.NewAppModule(distKeeper, supplyKeeper),
		supply.NewAppModule(supplyKeeper, accountKeeper),
	)
	configurator := module.NewConfigurator(legacyAmino, msgRouter, querier)
	am.RegisterServices(configurator)
	types.RegisterMsgServer(msgRouter, NewMsgServerImpl(NewDefaultPermissionKeeper(keeper)))
	types.RegisterQueryServer(querier, NewGrpcQuerier(appCodec, keys[types.ModuleName], keeper, keeper.queryGasLimit))

	govRouter := gov.NewRouter().
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(&paramsKeeper)).
		AddRoute(distributiontypes.RouterKey, distribution.NewDistributionProposalHandler(distKeeper))
	//AddRoute(types.RouterKey, NewWasmProposalHandler(&keeper, types.EnableAllProposals))

	govProposalHandlerRouter := govkeeper.NewProposalHandlerRouter()
	govProposalHandlerRouter.AddRoute(params.RouterKey, &paramsKeeper)

	govKeeper := gov.NewKeeper(
		legacyAmino, keys[govtypes.StoreKey], paramsKeeper, subspace(govtypes.ModuleName),
		supplyKeeper, &stakingKeeper, gov.DefaultParamspace, govRouter,
		bankKeeper, govProposalHandlerRouter, auth.FeeCollectorName,
	)

	//govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	//govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	//govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	//govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())

	keepers := TestKeepers{
		AccountKeeper:  accountKeeper,
		StakingKeeper:  stakingKeeper,
		supplyKeepr:    supplyKeeper,
		DistKeeper:     distKeeper,
		ContractKeeper: contractKeeper,
		WasmKeeper:     &keeper,
		BankKeeper:     bankKeeper,
		GovKeeper:      govKeeper,
		IBCKeeper:      ibcKeeper,
		Router:         router,
		EncodingConfig: encodingConfig,
		Faucet:         faucet,
		MultiStore:     ms,
	}
	return ctx, keepers
}

// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k types.ContractOpsKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = *ctx.SetEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case *types.MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case *types.MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleStoreCode(ctx sdk.Context, k types.ContractOpsKeeper, msg *types.MsgStoreCode) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	codeID, err := k.Create(ctx, senderAddr, msg.WASMByteCode, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", codeID)),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleInstantiate(ctx sdk.Context, k types.ContractOpsKeeper, msg *types.MsgInstantiateContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
	}

	contractAddr, _, err := k.Instantiate(ctx, msg.CodeID, senderAddr, adminAddr, msg.Msg, msg.Label, sdk.CoinAdaptersToCoins(msg.Funds))
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleExecute(ctx sdk.Context, k types.ContractOpsKeeper, msg *types.MsgExecuteContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "admin")
	}
	data, err := k.Execute(ctx, contractAddr, senderAddr, msg.Msg, sdk.CoinAdaptersToCoins(msg.Funds))
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   data,
		Events: ctx.EventManager().Events(),
	}, nil
}

var PubKeyCache = make(map[string]crypto.PubKey)

func RandomAccountAddress(_ testing.TB) sdk.AccAddress {
	_, pub, addr := keyPubAddr()
	PubKeyCache[addr.String()] = pub
	return addr
}

func RandomBech32AccountAddress(t testing.TB) string {
	return RandomAccountAddress(t).String()
}

type ExampleContract struct {
	InitialAmount sdk.Coins
	Creator       crypto.PrivKey
	CreatorAddr   sdk.AccAddress
	CodeID        uint64
}

func StoreHackatomExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/hackatom.wasm")
}

func StoreBurnerExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/burner.wasm")
}

func StoreIBCReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "./testdata/ibc_reflect.wasm")
}

func StoreReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) uint64 {
	_, _, creatorAddr := keyPubAddr()
	codeID, err := keepers.ContractKeeper.Create(ctx, creatorAddr, testdata.ReflectContractWasm(), nil)
	require.NoError(t, err)
	return codeID
}

func StoreExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers, wasmFile string) ExampleContract {
	anyAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000))
	creator, _, creatorAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, keepers.supplyKeepr, creatorAddr, anyAmount)

	wasmCode, err := ioutil.ReadFile(wasmFile)
	require.NoError(t, err)

	codeID, err := keepers.ContractKeeper.Create(ctx, creatorAddr, wasmCode, nil)
	require.NoError(t, err)
	return ExampleContract{anyAmount, creator, creatorAddr, codeID}
}

var wasmIdent = []byte("\x00\x61\x73\x6D")

type ExampleContractInstance struct {
	ExampleContract
	Contract sdk.AccAddress
}

// SeedNewContractInstance sets the mock wasmerEngine in keeper and calls store + instantiate to init the contract's metadata
func SeedNewContractInstance(t testing.TB, ctx sdk.Context, keepers TestKeepers, mock types.WasmerEngine) ExampleContractInstance {
	t.Helper()
	exampleContract := StoreRandomContract(t, ctx, keepers, mock)
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, exampleContract.CodeID, exampleContract.CreatorAddr, exampleContract.CreatorAddr, []byte(`{}`), "", nil)
	require.NoError(t, err)
	return ExampleContractInstance{
		ExampleContract: exampleContract,
		Contract:        contractAddr,
	}
}

// StoreRandomContract sets the mock wasmerEngine in keeper and calls store
func StoreRandomContract(t testing.TB, ctx sdk.Context, keepers TestKeepers, mock types.WasmerEngine) ExampleContract {
	return StoreRandomContractWithAccessConfig(t, ctx, keepers, mock, nil)
}

func StoreRandomContractWithAccessConfig(
	t testing.TB, ctx sdk.Context,
	keepers TestKeepers,
	mock types.WasmerEngine,
	cfg *types.AccessConfig,
) ExampleContract {
	t.Helper()
	anyAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000))
	creator, _, creatorAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, keepers.supplyKeepr, creatorAddr, anyAmount)
	keepers.WasmKeeper.wasmVM = mock
	wasmCode := append(wasmIdent, rand.Bytes(10)...) //nolint:gocritic
	codeID, err := keepers.ContractKeeper.Create(ctx, creatorAddr, wasmCode, cfg)
	require.NoError(t, err)
	exampleContract := ExampleContract{InitialAmount: anyAmount, Creator: creator, CreatorAddr: creatorAddr, CodeID: codeID}
	return exampleContract
}

type HackatomExampleInstance struct {
	ExampleContract
	Contract        sdk.AccAddress
	Verifier        crypto.PrivKey
	VerifierAddr    sdk.AccAddress
	Beneficiary     crypto.PrivKey
	BeneficiaryAddr sdk.AccAddress
}

// InstantiateHackatomExampleContract load and instantiate the "./testdata/hackatom.wasm" contract
func InstantiateHackatomExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) HackatomExampleInstance {
	contract := StoreHackatomExampleContract(t, ctx, keepers)

	verifier, _, verifierAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, keepers.supplyKeepr, verifierAddr, contract.InitialAmount)

	beneficiary, _, beneficiaryAddr := keyPubAddr()
	initMsgBz := HackatomExampleInitMsg{
		Verifier:    verifierAddr,
		Beneficiary: beneficiaryAddr,
	}.GetBytes(t)
	initialAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 100))

	adminAddr := contract.CreatorAddr
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, contract.CodeID, contract.CreatorAddr, adminAddr, initMsgBz, "demo contract to query", initialAmount)
	require.NoError(t, err)
	return HackatomExampleInstance{
		ExampleContract: contract,
		Contract:        contractAddr,
		Verifier:        verifier,
		VerifierAddr:    verifierAddr,
		Beneficiary:     beneficiary,
		BeneficiaryAddr: beneficiaryAddr,
	}
}

type HackatomExampleInitMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}

func (m HackatomExampleInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type IBCReflectExampleInstance struct {
	Contract      sdk.AccAddress
	Admin         sdk.AccAddress
	CodeID        uint64
	ReflectCodeID uint64
}

// InstantiateIBCReflectContract load and instantiate the "./testdata/ibc_reflect.wasm" contract
func InstantiateIBCReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) IBCReflectExampleInstance {
	reflectID := StoreReflectContract(t, ctx, keepers)
	ibcReflectID := StoreIBCReflectContract(t, ctx, keepers).CodeID

	initMsgBz := IBCReflectInitMsg{
		ReflectCodeID: reflectID,
	}.GetBytes(t)
	adminAddr := RandomAccountAddress(t)

	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, ibcReflectID, adminAddr, adminAddr, initMsgBz, "ibc-reflect-factory", nil)
	require.NoError(t, err)
	return IBCReflectExampleInstance{
		Admin:         adminAddr,
		Contract:      contractAddr,
		CodeID:        ibcReflectID,
		ReflectCodeID: reflectID,
	}
}

type IBCReflectInitMsg struct {
	ReflectCodeID uint64 `json:"reflect_code_id"`
}

func (m IBCReflectInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type BurnerExampleInitMsg struct {
	Payout sdk.AccAddress `json:"payout"`
}

func (m BurnerExampleInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

func fundAccounts(t testing.TB, ctx sdk.Context, am authkeeper.AccountKeeper, bank bank.Keeper, supplyKeeper supply.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	acc := am.NewAccountWithAddress(ctx, addr)
	am.SetAccount(ctx, acc)
	NewTestFaucet(t, ctx, bank, supplyKeeper, mint.ModuleName, coins...).Fund(ctx, addr, coins...)
}

var keyCounter uint64

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
