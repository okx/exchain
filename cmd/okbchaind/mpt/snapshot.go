package mpt

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/rootmulti"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	cfg "github.com/okx/okbchain/libs/tendermint/config"
	tmflags "github.com/okx/okbchain/libs/tendermint/libs/cli/flags"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/spf13/cobra"
	stdlog "log"
	"os"
	"path/filepath"
)

func genSnapCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genSnap",
		Short: "generate mpt store's snapshot",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			stdlog.Println("--------- generate snapshot start ---------")
			genSnapshot(ctx)
			stdlog.Println("--------- generate snapshot end ---------")
		},
	}
	return cmd
}

func genSnapshot(ctx *server.Context) {
	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	db, err := sdk.NewDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	mpt.SetSnapshotRebuild(true)
	mpt.AccountStateRootRetriever = accountStateRootRetriever{}
	rs := rootmulti.NewStore(db)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	const logLevel = "main:info,iavl:info,*:error,state:info,provider:info,root-multi:info"
	logger, err = tmflags.ParseLogLevel(logLevel, logger, cfg.DefaultLogLevel())
	rs.SetLogger(logger)
	rs.MountStoreWithDB(sdk.NewKVStoreKey(mpt.StoreKey), sdk.StoreTypeMPT, nil)
	rs.LoadLatestVersion()
}

type accountStateRootRetriever struct{}

func (a accountStateRootRetriever) RetrieveStateRoot(bz []byte) common.Hash {
	acc := DecodeAccount("", bz)
	return acc.GetStateRoot()
}

func (a accountStateRootRetriever) ModifyAccStateRoot(before []byte, rootHash common.Hash) []byte {
	//TODO implement me
	panic("implement me")
}

func (a accountStateRootRetriever) GetAccStateRoot(rootBytes []byte) common.Hash {
	//TODO implement me
	panic("implement me")
}
