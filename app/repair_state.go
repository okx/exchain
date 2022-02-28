package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/global"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mock"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/proxy"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/state/txindex"
	"github.com/okex/exchain/libs/tendermint/state/txindex/kv"
	"github.com/okex/exchain/libs/tendermint/state/txindex/null"
	"github.com/okex/exchain/libs/tendermint/store"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	applicationDB = "application"
	blockStoreDB  = "blockstore"
	stateDB       = "state"
	txIndexDB     = "tx_index"

	FlagStartHeight       string = "start-height"
	FlagEnableRepairState string = "enable-repair-state"
)

type repairApp struct {
	db dbm.DB
	*OKExChainApp
}

func (app *repairApp) getLatestVersion() int64 {
	rs := rootmulti.NewStore(app.db)
	return rs.GetLatestVersion()
}

func repairStateOnStart(ctx *server.Context) {
	// set flag
	orgIgnoreSmbCheck := sm.IgnoreSmbCheck
	orgIgnoreVersionCheck := iavl.GetIgnoreVersionCheck()
	orgEnableFlatKV := viper.GetBool(flatkv.FlagEnable)
	iavl.EnableAsyncCommit = false
	viper.Set(flatkv.FlagEnable, false)

	// repair state
	RepairState(ctx, true)

	//set original flag
	sm.SetIgnoreSmbCheck(orgIgnoreSmbCheck)
	iavl.SetIgnoreVersionCheck(orgIgnoreVersionCheck)
	iavl.EnableAsyncCommit = viper.GetBool(iavl.FlagIavlEnableAsyncCommit)
	viper.Set(flatkv.FlagEnable, orgEnableFlatKV)
	// load latest block height
	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	rmLockByDir(dataDir)
}

func RepairState(ctx *server.Context, onStart bool) {
	sm.SetIgnoreSmbCheck(true)
	iavl.SetIgnoreVersionCheck(true)

	// load latest block height
	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	latestBlockHeight := latestBlockHeight(dataDir)
	startBlockHeight := types.GetStartBlockHeight()
	fmt.Println("LLLL", latestBlockHeight, types.GetStartBlockHeight())
	if latestBlockHeight <= startBlockHeight+2 {
		log.Println(fmt.Sprintf("There is no need to repair data. The latest block height is %d, start block height is %d", latestBlockHeight, startBlockHeight))
		return
	}

	// create proxy app
	proxyApp, repairApp, err := createRepairApp(ctx)
	panicError(err)

	// get async commit version
	commitVersion, err := repairApp.GetCommitVersion()
	log.Println(fmt.Sprintf("repair state latestBlockHeight = %d \t commitVersion = %d", latestBlockHeight, commitVersion))
	panicError(err)

	if onStart && commitVersion == latestBlockHeight {
		log.Println("no need to repair state on start")
		return
	}

	// load state
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, _, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	// load start version
	startVersion := viper.GetInt64(FlagStartHeight)
	if startVersion == 0 {
		if onStart {
			startVersion = commitVersion
		} else {
			startVersion = commitVersion - 2 // case: state machine broken
		}
	}
	if startVersion <= 0 {
		panic("height too low, please restart from height 0 with genesis file")
	}

	err = repairApp.LoadStartVersion(startVersion)
	panicError(err)

	// repair data by apply the latest two blocks
	doRepair(ctx, state, stateStoreDB, proxyApp, startVersion, latestBlockHeight, dataDir)
}

func createRepairApp(ctx *server.Context) (proxy.AppConns, *repairApp, error) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	repairApp := newRepairApp(ctx.Logger, db, nil)

	clientCreator := proxy.NewLocalClientCreator(repairApp)
	// Create the proxyApp and establish connections to the ABCI app (consensus, mempool, query).
	proxyApp, err := createAndStartProxyAppConns(clientCreator)
	return proxyApp, repairApp, err
}

func newRepairApp(logger tmlog.Logger, db dbm.DB, traceStore io.Writer) *repairApp {
	return &repairApp{db, NewOKExChainApp(
		logger,
		db,
		traceStore,
		false,
		map[int64]bool{},
		0,
	)}
}

func doRepair(ctx *server.Context, state sm.State, stateStoreDB dbm.DB,
	proxyApp proxy.AppConns, startHeight, latestHeight int64, dataDir string) {
	stateCopy := state.Copy()
	ctx.Logger.Debug("stateCopy", "state", fmt.Sprintf("%+v", stateCopy))
	// construct state for repair
	state = constructStartState(state, stateStoreDB, startHeight)
	ctx.Logger.Debug("constructStartState", "state", fmt.Sprintf("%+v", state))
	var err error
	// repair state
	eventBus := types.NewEventBus()
	err = startEventBusAndIndexerService(ctx.Config, eventBus, ctx.Logger)
	panicError(err)
	blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{})
	blockExec.SetIsAsyncDeliverTx(viper.GetBool(sm.FlagParalleledTx))
	blockExec.SetEventBus(eventBus)
	global.SetGlobalHeight(startHeight + 1)
	for height := startHeight + 1; height <= latestHeight; height++ {
		repairBlock, repairBlockMeta := loadBlock(height, dataDir)
		state, _, err = blockExec.ApplyBlock(state, repairBlockMeta.BlockID, repairBlock)
		panicError(err)
		// use stateCopy to correct the repaired state
		if state.LastBlockHeight == stateCopy.LastBlockHeight {
			state.LastHeightConsensusParamsChanged = stateCopy.LastHeightConsensusParamsChanged
			state.LastHeightValidatorsChanged = stateCopy.LastHeightValidatorsChanged
			state.LastValidators = stateCopy.LastValidators.Copy()
			state.Validators = stateCopy.Validators.Copy()
			state.NextValidators = stateCopy.NextValidators.Copy()
			sm.SaveState(stateStoreDB, state)
		}
		ctx.Logger.Debug("repairedState", "state", fmt.Sprintf("%+v", state))
		res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
		panicError(err)
		repairedBlockHeight := res.LastBlockHeight
		repairedAppHash := res.LastBlockAppHash
		log.Println("Repaired block height", repairedBlockHeight)
		log.Println("Repaired app hash", fmt.Sprintf("%X", repairedAppHash))
	}
}

func startEventBusAndIndexerService(config *cfg.Config, eventBus *types.EventBus, logger tmlog.Logger) error {
	eventBus.SetLogger(logger.With("module", "events"))
	if err := eventBus.Start(); err != nil {
		return err
	}
	// Transaction indexing
	var txIndexer txindex.TxIndexer
	switch config.TxIndex.Indexer {
	case "kv":
		store, err := openDB(txIndexDB, filepath.Join(config.RootDir, "data"))
		if err != nil {
			return err
		}
		switch {
		case config.TxIndex.IndexKeys != "":
			txIndexer = kv.NewTxIndex(store, kv.IndexEvents(splitAndTrimEmpty(config.TxIndex.IndexKeys, ",", " ")))
		case config.TxIndex.IndexAllKeys:
			txIndexer = kv.NewTxIndex(store, kv.IndexAllEvents())
		default:
			txIndexer = kv.NewTxIndex(store)
		}
	default:
		txIndexer = &null.TxIndex{}
	}

	indexerService := txindex.NewIndexerService(txIndexer, eventBus)
	indexerService.SetLogger(logger.With("module", "txindex"))
	if err := indexerService.Start(); err != nil {
		return err
	}
	return nil
}

// splitAndTrimEmpty slices s into all subslices separated by sep and returns a
// slice of the string s with all leading and trailing Unicode code points
// contained in cutset removed. If sep is empty, SplitAndTrim splits after each
// UTF-8 sequence. First part is equivalent to strings.SplitN with a count of
// -1.  also filter out empty strings, only return non-empty strings.
func splitAndTrimEmpty(s, sep, cutset string) []string {
	if s == "" {
		return []string{}
	}

	spl := strings.Split(s, sep)
	nonEmptyStrings := make([]string, 0, len(spl))
	for i := 0; i < len(spl); i++ {
		element := strings.Trim(spl[i], cutset)
		if element != "" {
			nonEmptyStrings = append(nonEmptyStrings, element)
		}
	}
	return nonEmptyStrings
}

func constructStartState(state sm.State, stateStoreDB dbm.DB, startHeight int64) sm.State {
	stateCopy := state.Copy()
	validators, err := sm.LoadValidators(stateStoreDB, startHeight)
	lastValidators, err := sm.LoadValidators(stateStoreDB, startHeight-1)
	if err != nil {
		return stateCopy
	}
	nextValidators, err := sm.LoadValidators(stateStoreDB, startHeight+1)
	if err != nil {
		return stateCopy
	}
	consensusParams, err := sm.LoadConsensusParams(stateStoreDB, startHeight+1)
	if err != nil {
		return stateCopy
	}
	stateCopy.Validators = validators
	stateCopy.LastValidators = lastValidators
	stateCopy.NextValidators = nextValidators
	stateCopy.ConsensusParams = consensusParams
	stateCopy.LastBlockHeight = startHeight
	return stateCopy
}

func loadBlock(height int64, dataDir string) (*types.Block, *types.BlockMeta) {
	storeDB, err := openDB(blockStoreDB, dataDir)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	panicError(err)
	block := blockStore.LoadBlock(height)
	meta := blockStore.LoadBlockMeta(height)
	return block, meta
}

func latestBlockHeight(dataDir string) int64 {
	storeDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	return blockStore.Height()
}

func rmLockByDir(dataDir string) {
	files, _ := ioutil.ReadDir(dataDir)
	for _, f := range files {
		if f.IsDir() {
			os.Remove(filepath.Join(dataDir, f.Name(), "LOCK"))
		}
	}
}

// panic if error is not nil
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func openDB(dbName string, dataDir string) (db dbm.DB, err error) {
	return sdk.NewLevelDB(dbName, dataDir)
}

func createAndStartProxyAppConns(clientCreator proxy.ClientCreator) (proxy.AppConns, error) {
	proxyApp := proxy.NewAppConns(clientCreator)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("error starting proxy app connections: %v", err)
	}
	return proxyApp, nil
}
