package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	cmstore "github.com/okex/exchain/libs/cosmos-sdk/store"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb/util"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/node"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/store"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/dex"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/evidence"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/params"
	"github.com/okex/exchain/x/slashing"
	"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/token"
)

const (
	flagHeight    = "height"
	flagPruning   = "enable_pruning"
	flagDBBackend = "db_backend"

	blockDBName = "blockstore"
	stateDBName = "state"
	appDBName   = "application"
)

var wg sync.WaitGroup

func dataCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "modify data or query data in database",
	}

	cmd.AddCommand(
		pruningCmd(ctx),
		queryCmd(ctx),
		dbConvertCmd(ctx),
	)

	return cmd
}

func pruningCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune-compact",
		Short: "Prune and Compact blocks and application states",
	}

	cmd.AddCommand(pruneAllCmd(ctx),
		pruneAppCmd(ctx),
		pruneBlockCmd(ctx),
	)

	cmd.PersistentFlags().Int64P(flagHeight, "r", 0, "Removes block or state up to (but not including) a height")
	cmd.PersistentFlags().BoolP(flagPruning, "p", true, "Enable pruning")
	cmd.PersistentFlags().String(flagDBBackend, "goleveldb", "Database backend: goleveldb | rocksdb")
	return cmd
}

func pruneAllCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Compact both application states and blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			if err := checkBackend(dbm.BackendType(ctx.Config.DBBackend)); err != nil {
				return err
			}

			blockStoreDB := initDB(config, blockDBName)
			stateDB := initDB(config, stateDBName)
			appDB := initDB(config, appDBName)

			if viper.GetBool(flagPruning) {
				baseHeight, retainHeight := getPruneBlockParams(blockStoreDB)

				log.Println("--------- pruning start... ---------")
				wg.Add(3)
				go pruneBlocks(blockStoreDB, baseHeight, retainHeight)
				go pruneStates(stateDB, baseHeight, retainHeight)
				go pruneApp(appDB, baseHeight, retainHeight)
				wg.Wait()
				log.Println("--------- pruning end!!!   ---------")
			}

			log.Println("--------- compact start... ---------")
			wg.Add(3)
			go compactDB(blockStoreDB, blockDBName, dbm.BackendType(ctx.Config.DBBackend))
			go compactDB(stateDB, stateDBName, dbm.BackendType(ctx.Config.DBBackend))
			go compactDB(appDB, appDBName, dbm.BackendType(ctx.Config.DBBackend))
			wg.Wait()
			log.Println("--------- compact end!!!   ---------")

			return nil
		},
	}

	return cmd
}

func pruneAppCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Compact while pruning application state",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			if err := checkBackend(dbm.BackendType(ctx.Config.DBBackend)); err != nil {
				return err
			}

			appDB := initDB(config, appDBName)

			if viper.GetBool(flagPruning) {
				retainHeight := getPruneAppParams(appDB)

				wg.Add(1)
				go pruneApp(appDB, 1, retainHeight)
				wg.Wait()
			}

			log.Println("--------- compact start ---------")
			wg.Add(1)
			go compactDB(appDB, appDBName, dbm.BackendType(ctx.Config.DBBackend))
			wg.Wait()
			log.Println("--------- compact end ---------")

			return nil
		},
	}

	return cmd
}

func pruneBlockCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Compact while pruning blocks and states",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			if err := checkBackend(dbm.BackendType(ctx.Config.DBBackend)); err != nil {
				return err
			}

			blockStoreDB := initDB(config, blockDBName)
			stateDB := initDB(config, stateDBName)

			if viper.GetBool(flagPruning) {
				baseHeight, retainHeight := getPruneBlockParams(blockStoreDB)

				log.Println("--------- pruning start... ---------")
				wg.Add(2)
				go pruneBlocks(blockStoreDB, baseHeight, retainHeight)
				go pruneStates(stateDB, baseHeight, retainHeight)
				wg.Wait()
				log.Println("--------- pruning end!!!   ---------")
			}

			log.Println("--------- compact start... ---------")
			wg.Add(2)
			go compactDB(blockStoreDB, blockDBName, dbm.BackendType(ctx.Config.DBBackend))
			go compactDB(stateDB, stateDBName, dbm.BackendType(ctx.Config.DBBackend))
			wg.Wait()
			log.Println("--------- compact end!!!   ---------")

			return nil
		},
	}

	return cmd
}

func dbConvertCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert oec data from goleveldb to rocksdb",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			// {home}/data/*
			fromDir := ctx.Config.DBDir()
			toDir := filepath.Join(ctx.Config.RootDir, "data_convert")
			if _, err := os.Stat(toDir); os.IsNotExist(err) {
				err := os.MkdirAll(toDir, 0700)
				if err != nil {
					return fmt.Errorf("could not create directory %v: %w", toDir, err)
				}
			}

			fromFs, err := ioutil.ReadDir(fromDir)
			if err != nil {
				return err
			}
			for _, f := range fromFs {
				str := strings.Split(f.Name(), ".")
				if f.IsDir() && len(str) == 2 && str[1] == "db" {
					wg.Add(1)
					go func(name, fromDir, toDir string) {
						defer wg.Done()
						LtoR(name, fromDir, toDir)
					}(str[0], fromDir, toDir)
				} else {
					// cp
					err := exec.Command("cp", "-r", filepath.Join(fromDir, f.Name()), toDir).Run()
					if err != nil {
						panic("Execute Command failed:" + err.Error())
					}
				}
			}
			wg.Wait()

			// mv
			err = exec.Command("mv", fromDir, filepath.Join(ctx.Config.RootDir, "data_backup")).Run()
			if err != nil {
				panic("Execute Command failed:" + err.Error())
			}
			err = exec.Command("mv", toDir, filepath.Join(ctx.Config.RootDir, "data")).Run()
			if err != nil {
				panic("Execute Command failed:" + err.Error())
			}

			return nil
		},
	}

	return cmd
}

func getPruneBlockParams(blockStoreDB dbm.DB) (baseHeight, retainHeight int64) {
	baseHeight, size := getBlockInfo(blockStoreDB)

	retainHeight = viper.GetInt64(flagHeight)
	if retainHeight >= baseHeight+size-1 || retainHeight <= baseHeight {
		retainHeight = baseHeight + size - 2
	}

	return
}

func getPruneAppParams(appDB dbm.DB) (retainHeight int64) {
	rs := initAppStore(appDB)
	latestV := rs.GetLatestVersion()

	retainHeight = viper.GetInt64(flagHeight)
	if retainHeight >= latestV || retainHeight <= 1 {
		retainHeight = latestV - 1
	}

	return
}

func checkBackend(dbType dbm.BackendType) error {
	if _, ok := backends[dbType]; !ok {
		keys := make([]string, len(backends))
		i := 0
		for k := range backends {
			keys[i] = string(k)
			i++
		}
		return fmt.Errorf("unknown db_backend %s, expected <%s>", dbType, strings.Join(keys, " , "))
	}

	return nil
}

func initDB(config *cfg.Config, dbName string) dbm.DB {
	if dbName != blockDBName && dbName != stateDBName && dbName != appDBName {
		panic(fmt.Sprintf("unknow db name:%s", dbName))
	}

	db, err := node.DefaultDBProvider(&node.DBContext{dbName, config})
	panicError(err)

	return db
}

// pruneBlocks deletes blocks between the given heights (including from, excluding to).
func pruneBlocks(blockStoreDB dbm.DB, baseHeight, retainHeight int64) {
	defer wg.Done()

	log.Printf("Prune blocks [%d,%d)...", baseHeight, retainHeight)
	if retainHeight <= baseHeight {
		return
	}

	baseHeightBefore, sizeBefore := getBlockInfo(blockStoreDB)
	start := time.Now()
	_, err := store.NewBlockStore(blockStoreDB).PruneBlocks(retainHeight)
	if err != nil {
		panic(fmt.Errorf("failed to prune block store: %w", err))
	}

	baseHeightAfter, sizeAfter := getBlockInfo(blockStoreDB)
	log.Printf("Block db info [baseHeight,size]: [%d,%d] --> [%d,%d]\n", baseHeightBefore, sizeBefore, baseHeightAfter, sizeAfter)
	log.Printf("Prune blocks done in %v \n", time.Since(start))
}

// pruneStates deletes states between the given heights (including from, excluding to).
func pruneStates(stateDB dbm.DB, from, to int64) {
	defer wg.Done()

	log.Printf("Prune states [%d,%d)...", from, to)
	if to <= from {
		return
	}

	start := time.Now()
	if err := sm.PruneStates(stateDB, from, to); err != nil {
		panic(fmt.Errorf("failed to prune state database: %w", err))
	}

	log.Printf("Prune states done in %v \n", time.Since(start))
}

// pruneApp deletes app states between the given heights (including from, excluding to).
func pruneApp(appDB dbm.DB, from, to int64) {
	defer wg.Done()

	log.Printf("Prune applcation [%d,%d)...", from, to)
	if to <= from {
		return
	}

	rs := initAppStore(appDB)
	latestV := rs.GetLatestVersion()
	if to > latestV {
		return
	}
	versions := rs.GetVersions()
	if len(versions) == 0 {
		return
	}
	pruneHeights := rs.GetPruningHeights()

	newVersions := make([]int64, 0)
	newPruneHeights := make([]int64, 0)
	deleteVersions := make([]int64, 0)

	for _, v := range pruneHeights {
		if v >= to || v < from {
			newPruneHeights = append(newPruneHeights, v)
			continue
		}
		deleteVersions = append(deleteVersions, v)
	}

	for _, v := range versions {
		if v >= to || v < from {
			newVersions = append(newVersions, v)
			continue
		}
		deleteVersions = append(deleteVersions, v)
	}
	log.Printf("Prune application: Versions=%v, PruneVersions=%v", len(versions)+len(pruneHeights), len(deleteVersions))

	keysNumBefore, kvSizeBefore := calcKeysNum(appDB)
	start := time.Now()
	for key, store := range rs.GetStores() {
		if store.GetStoreType() == types.StoreTypeIAVL {
			// If the store is wrapped with an inter-block cache, we must first unwrap
			// it to get the underlying IAVL store.
			store = rs.GetCommitKVStore(key)

			if err := store.(*iavl.Store).DeleteVersions(deleteVersions...); err != nil {
				log.Printf("failed to delete version: %s", err)
			}
		}
	}

	rs.FlushPruneHeights(newPruneHeights, newVersions)

	keysNumAfter, kvSizeAfter := calcKeysNum(appDB)
	log.Printf("Application db key info [keysNum,kvSize]: [%d,%d] --> [%d,%d]\n", keysNumBefore, kvSizeBefore, keysNumAfter, kvSizeAfter)
	log.Printf("Prune application done in %v \n", time.Since(start))
}

func initAppStore(appDB dbm.DB) *rootmulti.Store {
	cms := cmstore.NewCommitMultiStore(appDB)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, upgrade.StoreKey, evidence.StoreKey,
		evm.StoreKey, token.StoreKey, token.KeyLock, dex.StoreKey, dex.TokenPairStoreKey,
		order.OrderStoreKey, ammswap.StoreKey, farm.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)

	}
	for _, key := range tkeys {
		cms.MountStoreWithDB(key, sdk.StoreTypeTransient, nil)
	}

	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	rs, ok := cms.(*rootmulti.Store)
	if !ok {
		panic("cms of from app is not rootmulti store")
	}

	return rs
}

func compactDB(db dbm.DB, name string, dbType dbm.BackendType) {
	defer wg.Done()

	log.Printf("Compact %s... \n", name)
	start := time.Now()

	if dbCompactor, ok := backends[dbType]; !ok {
		panic(fmt.Sprintf("Unknown db_backend %s, ", dbType))
	} else {
		dbCompactor(db)
	}

	log.Printf("Compact %s done in %v \n", name, time.Since(start))
}

func calcKeysNum(db dbm.DB) (keys, kvSize uint64) {
	iter, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; iter.Valid(); iter.Next() {
		keys++
		kvSize += uint64(len(iter.Key())) + uint64(len(iter.Value()))
	}
	iter.Close()
	return
}

func getBlockInfo(blockStoreDB dbm.DB) (baseHeight, size int64) {
	blockStore := store.NewBlockStore(blockStoreDB)
	baseHeight = blockStore.Base()
	size = blockStore.Size()
	return
}

func queryCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query blocks and states in database",
	}

	queryBlockState := &cobra.Command{
		Use:   "block",
		Short: "Query blocks heights in db",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			blockStoreDB := initDB(config, blockDBName)
			blockStore := store.NewBlockStore(blockStoreDB)
			fmt.Printf("[%d ~ %d]\n", blockStore.Base(), blockStore.Height())

			return nil
		},
	}

	queryAppState := &cobra.Command{
		Use:   "state",
		Short: "Query application states version in db",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			appStateDB := initDB(config, appDBName)

			rs := initAppStore(appStateDB)
			versions := rs.GetVersions()
			pruneHeights := rs.GetPruningHeights()

			fmt.Printf("%v\n", append(pruneHeights, versions...))
			return nil
		},
	}

	cmd.AddCommand(queryBlockState, queryAppState)

	return cmd
}

type dbCompactor func(dbm.DB)

var backends = map[dbm.BackendType]dbCompactor{}

func registerDBCompactor(dbType dbm.BackendType, compactor dbCompactor) {
	if _, ok := backends[dbType]; ok {
		return
	}
	backends[dbType] = compactor
}

func init() {
	dbCompactor := func(db dbm.DB) {
		if ldb, ok := db.(*dbm.GoLevelDB); ok {
			if err := ldb.DB().CompactRange(util.Range{}); err != nil {
				panic(err)
			}
		}
	}

	registerDBCompactor(dbm.GoLevelDBBackend, dbCompactor)
}
