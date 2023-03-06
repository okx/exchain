package app

import (
	"encoding/json"
	"fmt"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/prefix"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/module"
	"github.com/okx/okbchain/x/wasm"
	wasmtypes "github.com/okx/okbchain/x/wasm/types"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	dbm "github.com/okx/okbchain/libs/tm-db"

	"github.com/okx/okbchain/libs/cosmos-sdk/baseapp"
	"github.com/okx/okbchain/libs/cosmos-sdk/simapp"
	"github.com/okx/okbchain/libs/cosmos-sdk/simapp/helpers"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	distr "github.com/okx/okbchain/libs/cosmos-sdk/x/distribution"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/simulation"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/slashing"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply"
	"github.com/okx/okbchain/x/gov"
	"github.com/okx/okbchain/x/params"
	"github.com/okx/okbchain/x/staking"
)

func init() {
	simapp.GetSimulatorFlags()
}

type storeKeysPrefixes struct {
	A        sdk.StoreKey
	B        sdk.StoreKey
	Prefixes [][]byte
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, appName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application import/export simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, appName, app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	// nolint: dogsled
	_, newDB, newDir, _, _, err := simapp.SetupSimulation("leveldb-app-sim-2", "Simulation-2")
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewOKExChainApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, appName, newApp.Name())

	var genesisState map[string]json.RawMessage
	err = app.Codec().UnmarshalJSON(appState, &genesisState)
	require.NoError(t, err)

	ctxA := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []storeKeysPrefixes{
		{app.keys[baseapp.MainStoreKey], newApp.keys[baseapp.MainStoreKey], [][]byte{}},
		{app.keys[auth.StoreKey], newApp.keys[auth.StoreKey], [][]byte{}},
		{app.keys[staking.StoreKey], newApp.keys[staking.StoreKey],
			[][]byte{}}, // ordering may change but it doesn't matter
		{app.keys[slashing.StoreKey], newApp.keys[slashing.StoreKey], [][]byte{}},
		{app.keys[mint.StoreKey], newApp.keys[mint.StoreKey], [][]byte{}},
		{app.keys[distr.StoreKey], newApp.keys[distr.StoreKey], [][]byte{}},
		{app.keys[supply.StoreKey], newApp.keys[supply.StoreKey], [][]byte{}},
		{app.keys[params.StoreKey], newApp.keys[params.StoreKey], [][]byte{}},
		{app.keys[gov.StoreKey], newApp.keys[gov.StoreKey], [][]byte{}},
	}

	// reset contract code index in source DB for comparison with dest DB
	dropContractHistory := func(s store.KVStore, keys ...[]byte) {
		for _, key := range keys {
			prefixStore := prefix.NewStore(s, key)
			iter := prefixStore.Iterator(nil, nil)
			for ; iter.Valid(); iter.Next() {
				prefixStore.Delete(iter.Key())
			}
			iter.Close()
		}
	}
	prefixes := [][]byte{wasmtypes.ContractCodeHistoryElementPrefix, wasmtypes.ContractByCodeIDAndCreatedSecondaryIndexPrefix}
	dropContractHistory(ctxA.KVStore(app.keys[wasm.StoreKey]), prefixes...)
	dropContractHistory(ctxB.KVStore(newApp.keys[wasm.StoreKey]), prefixes...)

	normalizeContractInfo := func(ctx sdk.Context, app *OKExChainApp) {
		var index uint64
		app.WasmKeeper.IterateContractInfo(ctx, func(address sdk.AccAddress, info wasmtypes.ContractInfo) bool {
			created := &wasmtypes.AbsoluteTxPosition{
				BlockHeight: uint64(0),
				TxIndex:     index,
			}
			info.Created = created
			store := ctx.KVStore(app.keys[wasm.StoreKey])
			store.Set(wasmtypes.GetContractAddressKey(address), app.marshal.GetProtocMarshal().MustMarshal(&info))
			index++
			return false
		})
	}
	normalizeContractInfo(ctxA, app)
	normalizeContractInfo(ctxB, newApp)

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("compared %d key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Equal(t, len(failedKVAs), 0, simapp.GetSimulationLog(skp.A.Name(), app.SimulationManager().StoreDecoders, app.Codec(), failedKVAs, failedKVBs))
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	config, db, dir, logger, skip, err := simapp.SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation after import")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, appName, app.Name())

	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		simapp.SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	if stopEarly {
		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	// nolint: dosgsled
	_, newDB, newDir, _, _, err := simapp.SetupSimulation("leveldb-app-sim-2", "Simulation-2")
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewOKExChainApp(log.NewNopLogger(), newDB, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, appName, newApp.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		simapp.SimulationOperations(newApp, newApp.Codec(), config),
		newApp.ModuleAccountAddrs(), config,
	)
	require.NoError(t, err)
}

func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = helpers.SimAppChainID

	numTimesToRunPerSeed := 2
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	config.Seed = rand.Int63()

	for i := 0; i < numTimesToRunPerSeed; i++ {
		var logger log.Logger
		if simapp.FlagVerboseValue {
			logger = log.TestingLogger()
		} else {
			logger = log.NewNopLogger()
		}

		db := dbm.NewMemDB()

		app := NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, simapp.FlagPeriodValue, interBlockCacheOpt())

		fmt.Printf(
			"running non-determinism simulation; seed %d: attempt: %d/%d\n",
			config.Seed, i+1, numTimesToRunPerSeed,
		)

		_, _, err := simulation.SimulateFromSeed(
			t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
			simapp.SimulationOperations(app, app.Codec(), config),
			app.ModuleAccountAddrs(), config,
		)
		require.NoError(t, err)

		if config.Commit {
			simapp.PrintStats(db)
		}

		appHash := app.LastCommitID().Hash
		appHashList[i] = appHash

		if i != 0 {
			require.Equal(
				t, appHashList[0], appHashList[i],
				"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numTimesToRunPerSeed,
			)
		}
	}
}

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(codec *codec.Codec, manager *module.SimulationManager) simulation.AppStateFn {
	// quick hack to setup app state genesis with our app modules
	simapp.ModuleBasics = ModuleBasics
	if simapp.FlagGenesisTimeValue == 0 { // always set to have a block time
		simapp.FlagGenesisTimeValue = time.Now().Unix()
	}
	return simapp.AppStateFn(codec, manager)
}
