package main

import (
	"fmt"
	"log"
	"path/filepath"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/app"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagDisplayVersion string = "version"
	FlagDisplayAddress string = "address"
)

func displayStateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display-state",
		Short: "display account or contract state",
	}

	cmd.AddCommand(
		displayAccount(ctx),
		displayContract(ctx),
	)

	return cmd
}

func displayAccount(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "display account info at given height",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display account start ---------")
			displayAccountState(ctx)
			log.Println("--------- display account end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayAddress, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	cmd.Flags().String(sdk.FlagDBBackend, tmtypes.DBBackend, "Database backend: goleveldb | rocksdb")

	return cmd
}

func displayContract(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract",
		Short: "display contract state info at given height",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display contract state start ---------")
			displayContractState(ctx)
			log.Println("--------- display contract state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayAddress, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	cmd.Flags().String(sdk.FlagDBBackend, tmtypes.DBBackend, "Database backend: goleveldb | rocksdb")

	return cmd
}

func displayAccountState(ctx *server.Context) {
	dispApp := newDisplayApp(ctx)

	// load start version
	displayVersion := viper.GetInt64(FlagDisplayVersion)

	err := dispApp.LoadHeight(displayVersion)
	panicError(err)

	accountAddr := viper.GetString(FlagDisplayAddress)
	accAddr, err := sdk.AccAddressFromBech32(accountAddr)
	if err != nil {
		panic("Fail to parser AccAddress from : " + accountAddr)
	}

	// init deliver state
	dispApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: displayVersion + 1}})

	acc := dispApp.AccountKeeper.GetAccount(dispApp.GetDeliverStateCtx(), accAddr)
	fmt.Println("account is: ", acc.String())
}

func displayContractState(ctx *server.Context) {
	dispApp := newDisplayApp(ctx)

	// load start version
	displayVersion := viper.GetInt64(FlagDisplayVersion)

	err := dispApp.LoadHeight(displayVersion)
	panicError(err)

	contractAddr := viper.GetString(FlagDisplayAddress)
	addr := ethcmn.HexToAddress(contractAddr)

	// init deliver state
	dispApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: displayVersion + 1}})

	_ = dispApp.EvmKeeper.ForEachStorage(dispApp.GetDeliverStateCtx(), addr, func(key, value ethcmn.Hash) bool {
		fmt.Println("Key is: ", key.String(), ", value is: ", value.String())
		return false
	})
}

func newDisplayApp(ctx *server.Context) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	return app.NewOKExChainApp(
		ctx.Logger,
		db,
		nil,
		false,
		map[int64]bool{},
		0,
	)
}
