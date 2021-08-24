package main

import (
	"fmt"
	"log"
	"sync"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	cmstore "github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb/util"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
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
	flagStart   = "start"
	flagEnd     = "end"
	flagPruning = "pruning"
)

var wg sync.WaitGroup

func pruningCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pruning",
		Short: "Pruning blocks",
	}

	pruningAppStateCmd := &cobra.Command{
		Use:   "state",
		Short: "Pruning application state",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			blockStoreDB, _, appDB, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			if viper.GetBool(flagPruning) {
				log.Println("--------- pruning start ---------")
				blockStore := store.NewBlockStore(blockStoreDB)
				baseHeight := blockStore.Base()
				size := blockStore.Size()
				retainHeight := baseHeight + size - 2
				log.Printf("Data info: baseHeight=%d, size=%d, retainHeight=%d\n", baseHeight, size, retainHeight)

				start := viper.GetInt64(flagStart)
				if start < baseHeight || start >= retainHeight {
					start = baseHeight
				}
				end := viper.GetInt64(flagEnd)
				if end <= start || end >= retainHeight || end <= baseHeight {
					end = retainHeight
				}
				log.Printf("Pruning info: start=%d, end=%d\n", start, end)

				wg.Add(1)
				go pruneApp(appDB, start, end)
				wg.Wait()
				log.Println("--------- pruning end ---------")
			}

			// sync before compact
			log.Println("--------- compact start ---------")
			wg.Add(1)
			go compactDB(appDB)
			wg.Wait()
			log.Println("--------- compact end ---------")

			return nil
		},
	}

	pruningBlockStateCmd := &cobra.Command{
		Use:   "block",
		Short: "Pruning blocks and states",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			log.Println("--------- pruning start ---------")
			blockStoreDB, stateDB, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			if viper.GetBool(flagPruning) {
				blockStore := store.NewBlockStore(blockStoreDB)
				baseHeight := blockStore.Base()
				size := blockStore.Size()
				retainHeight := baseHeight + size - 2
				log.Printf("Data info: baseHeight=%d, size=%d, retainHeight=%d\n", baseHeight, size, retainHeight)

				start := viper.GetInt64(flagStart)
				if start < baseHeight || start >= retainHeight {
					start = baseHeight
				}
				end := viper.GetInt64(flagEnd)
				if end <= start || end >= retainHeight || end <= baseHeight {
					end = retainHeight
				}
				log.Printf("Pruning info: start=%d, end=%d\n", start, end)

				wg.Add(2)
				go pruneBlocks(blockStore, start, end)
				go pruneStates(stateDB, start, end)
				wg.Wait()
			}

			// sync before compact
			log.Println("--------- compact start ---------")
			wg.Add(2)
			go compactDB(blockStoreDB)
			go compactDB(stateDB)
			wg.Wait()
			log.Println("--------- compact end ---------")

			log.Println("--------- pruning end ---------")
			return nil
		},
	}

	cmd.PersistentFlags().Int64P(flagStart, "s", -1, "Pruning from the start height")
	cmd.PersistentFlags().Int64P(flagEnd, "e", -1, "Pruning to the end height")
	cmd.PersistentFlags().BoolP(flagPruning, "p", false, "enable Pruning")

	cmd.AddCommand(pruningAppStateCmd, pruningBlockStateCmd)
	return cmd
}

func initDBs(config *cfg.Config, dbProvider node.DBProvider) (blockStoreDB, stateDB, appDB dbm.DB, err error) {
	blockStoreDB, err = dbProvider(&node.DBContext{"blockstore", config})
	if err != nil {
		return
	}

	stateDB, err = dbProvider(&node.DBContext{"state", config})
	if err != nil {
		return
	}

	appDB, err = dbProvider(&node.DBContext{"application", config})
	if err != nil {
		return
	}

	return
}

// pruneBlocks deletes blocks between the given heights (including from, excluding to).
func pruneBlocks(blockStore *store.BlockStore, from, to int64) {
	defer wg.Done()

	log.Printf("Prune blocks [%d,%d)...", from, to)
	if to <= from {
		return
	}

	pruned, err := blockStore.PruneRange(from, to)
	if err != nil {
		panic(fmt.Errorf("failed to prune block store: %w", err))
	}

	log.Printf("Prune blocks end: pruned: %d, new base: %d, block len:%d\n", pruned, blockStore.Base(), blockStore.Size())
}

// pruneStates deletes states between the given heights (including from, excluding to).
func pruneStates(stateDB dbm.DB, from, to int64) {
	defer wg.Done()
	log.Printf("Prune states [%d,%d)...", from, to)
	if to <= from {
		return
	}

	if err := sm.PruneStates(stateDB, from, to); err != nil {
		panic(fmt.Errorf("failed to prune state database: %w", err))
	}
	log.Println("Prune states end!")
}

// pruneApp deletes app states between the given heights (including from, excluding to).
func pruneApp(appDB dbm.DB, from, to int64) {
	defer wg.Done()
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

	// remained heights
	newVersion := make([]int64, 0)
	for _, v := range versions {
		if v >= to || v < from {
			newVersion = append(newVersion, v)
			continue
		}
		pruneHeights = append(pruneHeights, v)
	}
	log.Printf("Prune app store: LatestVersion=%d,Versions=%v PruneHeights=%v", latestV, newVersion, pruneHeights)

	for key, store := range rs.GetStores() {
		if store.GetStoreType() == types.StoreTypeIAVL {
			// If the store is wrapped with an inter-block cache, we must first unwrap
			// it to get the underlying IAVL store.
			store = rs.GetCommitKVStore(key)

			if err := store.(*iavl.Store).DeleteVersions(pruneHeights...); err != nil {
				log.Printf("failed to delete version: %s", err)
			}
		}
	}

	pruneHeights = make([]int64, 0)
	rs.FlushPruneHeights(pruneHeights, newVersion)
	log.Println("Prune app store end!")
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

func compactDB(db dbm.DB) {
	defer wg.Done()
	err := db.(*dbm.GoLevelDB).DB().CompactRange(util.Range{})
	panicError(err)
}
