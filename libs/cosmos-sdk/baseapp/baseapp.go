package baseapp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/trace"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mempool"
	tmhttp "github.com/okex/exchain/libs/tendermint/rpc/client/http"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	runTxModeCheck          runTxMode = iota // Check a transaction
	runTxModeReCheck                         // Recheck a (pending) transaction after a commit
	runTxModeSimulate                        // Simulate a transaction
	runTxModeDeliver                         // Deliver a transaction
	runTxModeDeliverInAsync                  //Deliver a transaction in Aysnc
	runTxModeTrace                           // Trace a transaction
	runTxModeWrappedCheck

	// MainStoreKey is the string representation of the main store
	MainStoreKey = "main"

	// LatestSimulateTxHeight is the height to simulate tx based on the state of latest block height
	// only for runTxModeSimulate
	LatestSimulateTxHeight = 0
)

var (
	_ abci.Application = (*BaseApp)(nil)

	// mainConsensusParamsKey defines a key to store the consensus params in the
	// main store.
	mainConsensusParamsKey = []byte("consensus_params")

	globalMempool            mempool.Mempool
	mempoolEnableSort        = false
	mempoolEnableRecheck     = true
	mempoolEnablePendingPool = false
)

func GetGlobalMempool() mempool.Mempool {
	return globalMempool
}

func IsMempoolEnableSort() bool {
	return mempoolEnableSort
}

func IsMempoolEnableRecheck() bool {
	return cfg.DynamicConfig.GetMempoolRecheck()
}

func IsMempoolEnablePendingPool() bool {
	return mempoolEnablePendingPool
}

func SetGlobalMempool(mempool mempool.Mempool, enableSort bool, enablePendingPool bool) {
	globalMempool = mempool
	mempoolEnableSort = enableSort
	mempoolEnablePendingPool = enablePendingPool
}

type (
	// Enum mode for app.runTx
	runTxMode uint8

	// StoreLoader defines a customizable function to control how we load the CommitMultiStore
	// from disk. This is useful for state migration, when loading a datastore written with
	// an older version of the software. In particular, if a module changed the substore key name
	// (or removed a substore) between two versions of the software.
	StoreLoader func(ms sdk.CommitMultiStore) error
)

func (m runTxMode) String() (res string) {
	switch m {
	case runTxModeCheck:
		res = "ModeCheck"
	case runTxModeReCheck:
		res = "ModeReCheck"
	case runTxModeSimulate:
		res = "ModeSimulate"
	case runTxModeDeliver:
		res = "ModeDeliver"
	case runTxModeDeliverInAsync:
		res = "ModeDeliverInAsync"
	case runTxModeWrappedCheck:
		res = "ModeWrappedCheck"
	default:
		res = "Unknown"
	}

	return res
}

// BaseApp reflects the ABCI application implementation.
type BaseApp struct { // nolint: maligned
	// initialized on creation
	logger      log.Logger
	name        string               // application name from abci.Info
	db          dbm.DB               // common DB backend
	cms         sdk.CommitMultiStore // Main (uncached) state
	storeLoader StoreLoader          // function to handle store loading, may be overridden with SetStoreLoader()
	router      sdk.Router           // handle any kind of message
	queryRouter sdk.QueryRouter      // router for redirecting query calls

	// txDecoder returns a cosmos-sdk/types.Tx interface that definitely is an StdTx or a MsgEthereumTx
	txDecoder sdk.TxDecoder

	// set upon LoadVersion or LoadLatestVersion.
	baseKey *sdk.KVStoreKey // Main KVStore in cms

	anteHandler      sdk.AnteHandler      // ante handler for fee and auth
	GasRefundHandler sdk.GasRefundHandler // gas refund handler for gas refund
	AccHandler       sdk.AccHandler       // account handler for cm tx nonce

	initChainer    sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker   sdk.BeginBlocker // logic to run before any txs
	endBlocker     sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	addrPeerFilter sdk.PeerFilter   // filter peers by address and port
	idPeerFilter   sdk.PeerFilter   // filter peers by node ID
	fauxMerkleMode bool             // if true, IAVL MountStores uses MountStoresDB for simulation speed.

	getTxFee                     sdk.GetTxFeeHandler
	updateFeeCollectorAccHandler sdk.UpdateFeeCollectorAccHandler
	logFix                       sdk.LogFix

	// volatile states:
	//
	// checkState is set on InitChain and reset on Commit
	// deliverState is set on InitChain and BeginBlock and set to nil on Commit
	checkState   *state // for CheckTx
	deliverState *state // for DeliverTx

	// an inter-block write-through cache provided to the context during deliverState
	interBlockCache sdk.MultiStorePersistentCache

	// absent validators from begin block
	voteInfos []abci.VoteInfo

	// consensus params
	// TODO: Move this in the future to baseapp param store on main store.
	consensusParams *abci.ConsensusParams

	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. This is mainly used for DoS and spam prevention.
	minGasPrices sdk.DecCoins

	// flag for sealing options and parameters to a BaseApp
	sealed bool

	// block height at which to halt the chain and gracefully shutdown
	haltHeight uint64

	// minimum block time (in Unix seconds) at which to halt the chain and gracefully shutdown
	haltTime uint64

	// application's version string
	appVersion string

	// trace set will return full stack traces for errors in ABCI Log field
	trace bool

	// start record handle
	startLog recordHandle

	// end record handle
	endLog recordHandle

	parallelTxManage *parallelTxManager

	chainCache *sdk.Cache
	blockCache *sdk.Cache

	checkTxNum        int64
	wrappedCheckTxNum int64
	anteTracer        *trace.Tracer
}

type recordHandle func(string)

// NewBaseApp returns a reference to an initialized BaseApp. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
//
// NOTE: The db is used to store the version number for now.
func NewBaseApp(
	name string, logger log.Logger, db dbm.DB, txDecoder sdk.TxDecoder, options ...func(*BaseApp),
) *BaseApp {

	app := &BaseApp{
		logger:         logger,
		name:           name,
		db:             db,
		cms:            store.NewCommitMultiStore(db),
		storeLoader:    DefaultStoreLoader,
		router:         NewRouter(),
		queryRouter:    NewQueryRouter(),
		fauxMerkleMode: false,
		trace:          false,

		parallelTxManage: newParallelTxManager(),
		chainCache:       sdk.NewChainCache(),
		txDecoder:        txDecoder,
		anteTracer:       trace.NewTracer(trace.AnteChainDetail),
	}

	for _, option := range options {
		option(app)
	}

	if app.interBlockCache != nil {
		app.cms.SetInterBlockCache(app.interBlockCache)
	}
	app.cms.SetLogger(app.logger)

	app.parallelTxManage.workgroup.Start()

	return app
}

// Name returns the name of the BaseApp.
func (app *BaseApp) Name() string {
	return app.name
}

// AppVersion returns the application's version string.
func (app *BaseApp) AppVersion() string {
	return app.appVersion
}

// Logger returns the logger of the BaseApp.
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// SetStartLogHandler set the startLog of the BaseApp.
func (app *BaseApp) SetStartLogHandler(handle recordHandle) {
	app.startLog = handle
}

// SetStopLogHandler set the endLog of the BaseApp.
func (app *BaseApp) SetEndLogHandler(handle recordHandle) {
	app.endLog = handle
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...sdk.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			if !app.fauxMerkleMode {
				app.MountStore(key, sdk.StoreTypeIAVL)
			} else {
				// StoreTypeDB doesn't do anything upon commit, and it doesn't
				// retain history, but it's useful for faster simulation.
				app.MountStore(key, sdk.StoreTypeDB)
			}

		case *sdk.TransientStoreKey:
			app.MountStore(key, sdk.StoreTypeTransient)

		default:
			panic("Unrecognized store key type " + reflect.TypeOf(key).Name())
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountKVStores(keys map[string]*sdk.KVStoreKey) {
	for _, key := range keys {
		if !app.fauxMerkleMode {
			app.MountStore(key, sdk.StoreTypeIAVL)
		} else {
			// StoreTypeDB doesn't do anything upon commit, and it doesn't
			// retain history, but it's useful for faster simulation.
			app.MountStore(key, sdk.StoreTypeDB)
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountTransientStores(keys map[string]*sdk.TransientStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeTransient)
	}
}

// MountStoreWithDB mounts a store to the provided key in the BaseApp
// multistore, using a specified DB.
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion(baseKey *sdk.KVStoreKey) error {
	err := app.storeLoader(app.cms)
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// DefaultStoreLoader will be used by default and loads the latest version
func DefaultStoreLoader(ms sdk.CommitMultiStore) error {
	return ms.LoadLatestVersion()
}

// StoreLoaderWithUpgrade is used to prepare baseapp with a fixed StoreLoader
// pattern. This is useful in test cases, or with custom upgrade loading logic.
func StoreLoaderWithUpgrade(upgrades *storetypes.StoreUpgrades) StoreLoader {
	return func(ms sdk.CommitMultiStore) error {
		return ms.LoadLatestVersionAndUpgrade(upgrades)
	}
}

// UpgradeableStoreLoader can be configured by SetStoreLoader() to check for the
// existence of a given upgrade file - json encoded StoreUpgrades data.
//
// If not file is present, it will peform the default load (no upgrades to store).
//
// If the file is present, it will parse the file and execute those upgrades
// (rename or delete stores), while loading the data. It will also delete the
// upgrade file upon successful load, so that the upgrade is only applied once,
// and not re-applied on next restart
//
// This is useful for in place migrations when a store key is renamed between
// two versions of the software. (TODO: this code will move to x/upgrades
// when PR #4233 is merged, here mainly to help test the design)
func UpgradeableStoreLoader(upgradeInfoPath string) StoreLoader {
	return func(ms sdk.CommitMultiStore) error {
		_, err := os.Stat(upgradeInfoPath)
		if os.IsNotExist(err) {
			return DefaultStoreLoader(ms)
		} else if err != nil {
			return err
		}

		// there is a migration file, let's execute
		data, err := ioutil.ReadFile(upgradeInfoPath)
		if err != nil {
			return fmt.Errorf("cannot read upgrade file %s: %v", upgradeInfoPath, err)
		}

		var upgrades storetypes.StoreUpgrades
		err = json.Unmarshal(data, &upgrades)
		if err != nil {
			return fmt.Errorf("cannot parse upgrade file: %v", err)
		}

		err = ms.LoadLatestVersionAndUpgrade(&upgrades)
		if err != nil {
			return fmt.Errorf("load and upgrade database: %v", err)
		}

		// if we have a successful load, we delete the file
		err = os.Remove(upgradeInfoPath)
		if err != nil {
			return fmt.Errorf("deleting upgrade file %s: %v", upgradeInfoPath, err)
		}
		return nil
	}
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running baseapp.
func (app *BaseApp) LoadVersion(version int64, baseKey *sdk.KVStoreKey) error {
	err := app.cms.LoadVersion(version)
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// GetCommitVersion loads the latest committed version.
func (app *BaseApp) GetCommitVersion() (int64, error) {
	return app.cms.GetCommitVersion()
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromMainStore(baseKey *sdk.KVStoreKey) error {
	mainStore := app.cms.GetKVStore(baseKey)
	if mainStore == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}

	// memoize baseKey
	if app.baseKey != nil {
		panic("app.baseKey expected to be nil; duplicate init?")
	}
	app.baseKey = baseKey

	// Load the consensus params from the main store. If the consensus params are
	// nil, it will be saved later during InitChain.
	//
	// TODO: assert that InitChain hasn't yet been called.
	consensusParamsBz := mainStore.Get(mainConsensusParamsKey)
	if consensusParamsBz != nil {
		var consensusParams = &abci.ConsensusParams{}

		err := proto.Unmarshal(consensusParamsBz, consensusParams)
		if err != nil {
			panic(err)
		}

		app.setConsensusParams(consensusParams)
	}

	// needed for the export command which inits from store but never calls initchain
	app.setCheckState(abci.Header{})
	app.Seal()

	return nil
}

func (app *BaseApp) setMinGasPrices(gasPrices sdk.DecCoins) {
	app.minGasPrices = gasPrices
}

func (app *BaseApp) setHaltHeight(haltHeight uint64) {
	app.haltHeight = haltHeight
}

func (app *BaseApp) setHaltTime(haltTime uint64) {
	app.haltTime = haltTime
}

func (app *BaseApp) setInterBlockCache(cache sdk.MultiStorePersistentCache) {
	app.interBlockCache = cache
}

func (app *BaseApp) setTrace(trace bool) {
	app.trace = trace
}

// Router returns the router of the BaseApp.
func (app *BaseApp) Router() sdk.Router {
	if app.sealed {
		// We cannot return a Router when the app is sealed because we can't have
		// any routes modified which would cause unexpected routing behavior.
		panic("Router() on sealed BaseApp")
	}
	return app.router
}

// QueryRouter returns the QueryRouter of a BaseApp.
func (app *BaseApp) QueryRouter() sdk.QueryRouter { return app.queryRouter }

// Seal seals a BaseApp. It prohibits any further modifications to a BaseApp.
func (app *BaseApp) Seal() { app.sealed = true }

// IsSealed returns true if the BaseApp is sealed and false otherwise.
func (app *BaseApp) IsSealed() bool { return app.sealed }

// setCheckState sets the BaseApp's checkState with a cache-wrapped multi-store
// (i.e. a CacheMultiStore) and a new Context with the cache-wrapped multi-store,
// provided header, and minimum gas prices set. It is set on InitChain and reset
// on Commit.
func (app *BaseApp) setCheckState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, true, app.logger).WithMinGasPrices(app.minGasPrices),
	}
}

// setDeliverState sets the BaseApp's deliverState with a cache-wrapped multi-store
// (i.e. a CacheMultiStore) and a new Context with the cache-wrapped multi-store,
// and provided header. It is set on InitChain and BeginBlock and set to nil on
// Commit.
func (app *BaseApp) setDeliverState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.logger),
	}
}

// setTraceState sets the BaseApp's traceState with a cache-wrapped multi-store
// (i.e. a CacheMultiStore) and a new Context with the cache-wrapped multi-store,
// and provided header. It is set at the start of trace tx
func (app *BaseApp) newTraceState(header abci.Header, height int64) (*state, error) {
	ms, err := app.cms.CacheMultiStoreWithVersion(height)
	if err != nil {
		return nil, err
	}
	return &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.logger),
	}, nil
}

// setConsensusParams memoizes the consensus params.
func (app *BaseApp) setConsensusParams(consensusParams *abci.ConsensusParams) {
	app.consensusParams = consensusParams
}

// setConsensusParams stores the consensus params to the main store.
func (app *BaseApp) storeConsensusParams(consensusParams *abci.ConsensusParams) {
	consensusParamsBz, err := proto.Marshal(consensusParams)
	if err != nil {
		panic(err)
	}
	mainStore := app.cms.GetKVStore(app.baseKey)
	mainStore.Set(mainConsensusParamsKey, consensusParamsBz)
}

// getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// if maximum block gas is less than negative one and returns zero if negative
// one.
func (app *BaseApp) getMaximumBlockGas() uint64 {
	if app.consensusParams == nil || app.consensusParams.Block == nil {
		return 0
	}

	maxGas := app.consensusParams.Block.MaxGas
	switch {
	case maxGas < -1:
		panic(fmt.Sprintf("invalid maximum block gas: %d", maxGas))

	case maxGas == -1:
		return 0

	default:
		return uint64(maxGas)
	}
}

func (app *BaseApp) validateHeight(req abci.RequestBeginBlock) error {
	if req.Header.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Header.Height)
	}

	prevHeight := app.LastBlockHeight()
	if req.Header.Height != prevHeight+1 {
		return fmt.Errorf("invalid height: %d; expected: %d", req.Header.Height, prevHeight+1)
	}

	return nil
}

// validateBasicTxMsgs executes basic validator calls for messages.
func validateBasicTxMsgs(msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must contain at least one message")
	}

	for _, msg := range msgs {
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the applications's deliverState if app is in runTxModeDeliver,
// otherwise it returns the application's checkstate.
func (app *BaseApp) getState(mode runTxMode) *state {
	if mode == runTxModeDeliver || mode == runTxModeDeliverInAsync {
		return app.deliverState
	}

	return app.checkState
}

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getContextForTx(mode runTxMode, txBytes []byte) sdk.Context {
	ctx := app.getState(mode).ctx.
		WithTxBytes(txBytes).
		WithVoteInfos(app.voteInfos).
		WithConsensusParams(app.consensusParams)

	if mode == runTxModeReCheck {
		ctx = ctx.WithIsReCheckTx(true)
	}

	if mode == runTxModeWrappedCheck {
		ctx = ctx.WithIsWrappedCheckTx(true)
	}

	if mode == runTxModeSimulate {
		ctx, _ = ctx.CacheContext()
	}

	if mode == runTxModeDeliver {
		if s := app.parallelTxManage.GetTxSignCache(txBytes); s != nil {
			ctx = ctx.WithSigCache(s)
		}
	}

	if app.parallelTxManage.isAsyncDeliverTx && mode == runTxModeDeliverInAsync {
		ctx = ctx.WithAsync()
		if s, ok := app.parallelTxManage.txStatus[string(txBytes)]; ok && s.signCache != nil {
			ctx = ctx.WithSigCache(s.signCache)
		}
		ctx = ctx.WithTxBytes(getRealTxByte(txBytes))
	}

	return ctx
}

// retrieve the context for simulating the tx w/ txBytes
func (app *BaseApp) getContextForSimTx(txBytes []byte, height int64) (sdk.Context, error) {
	cms, ok := app.cms.(*rootmulti.Store)
	if !ok {
		return sdk.Context{}, fmt.Errorf("get context for simulate tx failed")
	}

	ms, err := cms.CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{}, err
	}

	abciHeader, err := GetABCIHeader(height)
	if err != nil {
		return sdk.Context{}, err
	}

	simState := &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, abciHeader, true, app.logger).WithMinGasPrices(app.minGasPrices),
	}

	ctx := simState.ctx.WithTxBytes(txBytes)

	return ctx, nil
}
func GetABCITx(hash []byte) (*ctypes.ResultTx, error) {
	laddr := viper.GetString("rpc.laddr")
	splits := strings.Split(laddr, ":")
	if len(splits) < 2 {
		return nil, fmt.Errorf("get tx failed!")
	}

	rpcCli, err := tmhttp.New(fmt.Sprintf("tcp://127.0.0.1:%s", splits[len(splits)-1]), "/websocket")
	if err != nil {
		return nil, fmt.Errorf("get tx failed!")
	}

	tx, err := rpcCli.Tx(hash, false)
	if err != nil {
		return nil, fmt.Errorf("get ABCI tx failed!")
	}

	return tx, nil
}
func GetABCIBlock(height int64) (*ctypes.ResultBlock, error) {
	laddr := viper.GetString("rpc.laddr")
	splits := strings.Split(laddr, ":")
	if len(splits) < 2 {
		return nil, fmt.Errorf("get tendermint Block failed!")
	}

	rpcCli, err := tmhttp.New(fmt.Sprintf("tcp://127.0.0.1:%s", splits[len(splits)-1]), "/websocket")
	if err != nil {
		return nil, fmt.Errorf("get tendermint Block failed!")
	}

	block, err := rpcCli.Block(&height)
	if err != nil {
		return nil, fmt.Errorf("get tendermint Block failed!")
	}
	return block, nil
}
func GetABCIHeader(height int64) (abci.Header, error) {
	block, err := GetABCIBlock(height)
	if err != nil {
		return abci.Header{}, fmt.Errorf("get ABCI header failed!")
	}
	return blockHeaderToABCIHeader(block.Block.Header), nil
}

func blockHeaderToABCIHeader(header tmtypes.Header) abci.Header {
	return abci.Header{
		Version: abci.Version{
			Block: uint64(header.Version.Block),
			App:   uint64(header.Version.App),
		},
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,
		LastBlockId: abci.BlockID{
			Hash: header.LastBlockID.Hash,
			PartsHeader: abci.PartSetHeader{
				Total: int32(header.LastBlockID.PartsHeader.Total),
				Hash:  header.LastBlockID.PartsHeader.Hash,
			},
		},
		LastCommitHash:     header.LastCommitHash,
		DataHash:           header.DataHash,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,
		EvidenceHash:       header.EvidenceHash,
		ProposerAddress:    header.ProposerAddress,
	}
}

// cacheTxContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (sdk.Context, sdk.CacheMultiStore) {
	ms := ctx.MultiStore()
	// TODO: https://github.com/cosmos/cosmos-sdk/issues/2824
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.SetTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmtypes.Tx(txBytes).Hash(ctx.BlockHeight())),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

func (app *BaseApp) pin(tag string, start bool, mode runTxMode) {

	if mode != runTxModeDeliver {
		return
	}

	if app.startLog != nil {
		if start {
			app.startLog(tag)
		} else {
			app.endLog(tag)
		}
	}
}

// runMsgs iterates through a list of messages and executes them with the provided
// Context and execution mode. Messages will only be executed during simulation
// and DeliverTx. An error is returned if any single message fails or if a
// Handler does not exist for a given message route. Otherwise, a reference to a
// Result is returned. The caller must not commit state if an error is returned.
func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode runTxMode) (*sdk.Result, error) {
	msgLogs := make(sdk.ABCIMessageLogs, 0, len(msgs))
	data := make([]byte, 0, len(msgs))
	events := sdk.EmptyEvents()

	// NOTE: GasWanted is determined by the AnteHandler and GasUsed by the GasMeter.
	for i, msg := range msgs {
		// skip actual execution for (Re)CheckTx mode
		if mode == runTxModeCheck || mode == runTxModeReCheck || mode == runTxModeWrappedCheck {
			break
		}

		msgRoute := msg.Route()
		handler := app.router.Route(ctx, msgRoute)
		if handler == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s; message index: %d", msgRoute, i)
		}

		msgResult, err := handler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to execute message; message index: %d", i)
		}

		msgEvents := sdk.Events{
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type())),
		}
		//app.pin("AppendEvents", true, mode)

		msgEvents = msgEvents.AppendEvents(msgResult.Events)

		// append message events, data and logs
		//
		// Note: Each message result's data must be length-prefixed in order to
		// separate each result.
		events = events.AppendEvents(msgEvents)
		data = append(data, msgResult.Data...)
		msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint16(i), msgResult.Log, msgEvents))
		//app.pin("AppendEvents", false, mode)

	}

	return &sdk.Result{
		Data:   data,
		Log:    strings.TrimSpace(msgLogs.String()),
		Events: events,
	}, nil
}

func (app *BaseApp) Export(toApp *BaseApp, version int64) error {
	fromCms, ok := app.cms.(*rootmulti.Store)
	if !ok {
		return fmt.Errorf("cms of from app is not rootmulti store")
	}

	toCms, ok := toApp.cms.(*rootmulti.Store)
	if !ok {
		return fmt.Errorf("cms of to app is not rootmulti store")
	}

	return fromCms.Export(toCms, version)
}

func (app *BaseApp) StopStore() {
	app.cms.StopStore()
}

func (app *BaseApp) GetTxInfo(ctx sdk.Context, tx sdk.Tx) mempool.ExTxInfo {
	exTxInfo := tx.GetTxInfo(ctx)
	if exTxInfo.Nonce == 0 && exTxInfo.Sender != "" && app.AccHandler != nil {
		addr, _ := sdk.AccAddressFromBech32(exTxInfo.Sender)
		exTxInfo.Nonce = app.AccHandler(ctx, addr)

		if app.anteHandler != nil && exTxInfo.Nonce > 0 {
			exTxInfo.Nonce -= 1 // in ante handler logical, the nonce will incress one
		}
	}

	return exTxInfo
}

func (app *BaseApp) GetRawTxInfo(rawTx tmtypes.Tx) mempool.ExTxInfo {
	tx, err := app.txDecoder(rawTx)
	if err != nil {
		return mempool.ExTxInfo{}
	}

	return app.GetTxInfo(app.checkState.ctx, tx)
}

func (app *BaseApp) GetTxHistoryGasUsed(rawTx tmtypes.Tx) int64 {
	tx, err := app.txDecoder(rawTx)
	if err != nil {
		return -1
	}

	txFnSig, toDeployContractSize := tx.GetTxFnSignatureInfo()
	if txFnSig == nil {
		return -1
	}

	db := InstanceOfHistoryGasUsedRecordDB()
	data, err := db.Get(txFnSig)
	if err != nil || len(data) == 0 {
		return -1
	}

	if toDeployContractSize > 0 {
		// if deploy contract case, the history gas used value is unit gas used
		return int64(binary.BigEndian.Uint64(data))*int64(toDeployContractSize) + int64(1000)
	}

	return int64(binary.BigEndian.Uint64(data))
}
