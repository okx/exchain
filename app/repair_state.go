package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	"github.com/okex/exchain/libs/iavl"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mock"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/proxy"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/store"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const (
	applicationDB = "application"
	blockStoreDB  = "blockstore"
	stateDB       = "state"

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
	orgIgnoreSmbCheck := sm.IgnoreSmbCheck
	orgIgnoreVersionCheck := iavl.GetIgnoreVersionCheck()
	iavl.EnableAsyncCommit = false
	RepairState(ctx, true)
	//set original config
	sm.SetIgnoreSmbCheck(orgIgnoreSmbCheck)
	iavl.SetIgnoreVersionCheck(orgIgnoreVersionCheck)
	iavl.EnableAsyncCommit = viper.GetBool(iavl.FlagIavlEnableAsyncCommit)
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
		latestVersion := repairApp.getLatestVersion()
		if types.HigherThanMars(latestVersion) {
			lastMptVersion := int64(repairApp.EvmKeeper.GetLatestStoredBlockHeight())
			if lastMptVersion < latestVersion {
				latestVersion = lastMptVersion
			}
		}
		startVersion = latestVersion - 2
	}
	if startVersion <= 0 {
		panic("height too low, please restart from height 0 with genesis file")
	}

	if commitVersion < startVersion {
		startVersion = commitVersion
	}

	err = repairApp.LoadStartVersion(startVersion)
	panicError(err)

	sdk.TrieDirtyDisabled = true
	repairApp.EvmKeeper.SetTargetMptVersion(startVersion)

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
	// construct state for repair
	state = constructStartState(state, stateStoreDB, startHeight)
	ctx.Logger.Debug("constructStartState", "state", fmt.Sprintf("%+v", state))
	var err error
	// repair state
	blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{})
	blockExec.SetIsAsyncDeliverTx(viper.GetBool(sm.FlagParalleledTx))
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
			state.NextValidators = state.NextValidators.Copy()
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
