package mpt

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	iavlstore "github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
)

func cleanIavlStoreCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean-iavl",
		Short: "clean up migrated iavl store",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- clean state start ---------")
			clean(ctx)
			log.Println("--------- clean state end ---------")
		},
	}
	return cmd
}

func clean(ctx *server.Context) {
	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	db, err := sdk.NewLevelDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	rs := rootmulti.NewStore(db)
	latestVersion := rs.GetLatestVersion()

	// 1.clean acc store
	fmt.Println("Start to clean account store")
	if err = deleteIavlStore(db, []byte(iavlAccKey), latestVersion, iavlstore.IavlCacheSize); err != nil {
		fmt.Println("fail to clean iavl store: ", err)
	}

	// 2.clean evm store
	fmt.Println("Start to clean evm store")
	if err = deleteIavlStore(db, []byte(iavlEvmKey), latestVersion, iavlstore.IavlCacheSize); err != nil {
		fmt.Println("fail to clean iavl store: ", err)
	}
}

func deleteIavlStore(db dbm.DB, prefix []byte, maxVersion int64, cacheSize int) error {
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return err
	}

	// delete verion [from, to)
	return tree.DeleteVersionsRange(0, maxVersion+1, true)
}
