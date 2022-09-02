package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain/app/logevents"
	"github.com/okex/exchain/cmd/exchaind/mpt"

	"github.com/okex/exchain/app/rpc"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	clientkeys "github.com/okex/exchain/libs/cosmos-sdk/client/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/cmd/client"
	"github.com/okex/exchain/x/genutil"
	genutilcli "github.com/okex/exchain/x/genutil/client/cli"
	genutiltypes "github.com/okex/exchain/x/genutil/types"
	"github.com/okex/exchain/x/staking"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false

	codecProxy, registry := codec.MakeCodecSuit(app.ModuleBasics)

	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	keys.CryptoCdc = codecProxy.GetCdc()
	genutil.ModuleCdc = codecProxy.GetCdc()
	genutiltypes.ModuleCdc = codecProxy.GetCdc()
	clientkeys.KeysCdc = codecProxy.GetCdc()

	config := sdk.GetConfig()
	okexchain.SetBech32Prefixes(config)
	okexchain.SetBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "exchaind",
		Short:             "ExChain App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		client.ValidateChainID(
			genutilcli.InitCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(ctx, codecProxy.GetCdc(), auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, codecProxy.GetCdc()),
		genutilcli.GenTxCmd(
			ctx, codecProxy.GetCdc(), app.ModuleBasics, staking.AppModuleBasic{}, auth.GenesisAccountIterator{},
			app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics),
		client.TestnetCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics, auth.GenesisAccountIterator{}),
		replayCmd(ctx, client.RegisterAppFlag, codecProxy, newApp, registry, registerRoutes),
		repairStateCmd(ctx),
		displayStateCmd(ctx),
		mpt.MptCmd(ctx),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		AddGenesisAccountCmd(ctx, codecProxy.GetCdc(), app.DefaultNodeHome, app.DefaultCLIHome),
		flags.NewCompletionCmd(rootCmd, true),
		dataCmd(ctx),
		exportAppCmd(ctx),
		iaviewerCmd(ctx, codecProxy.GetCdc()),
		subscribeCmd(codecProxy.GetCdc()),
	)

	subFunc := func(logger log.Logger) log.Subscriber {
		return logevents.NewProvider(logger)
	}
	// Tendermint node base commands
	server.AddCommands(ctx, codecProxy, registry, rootCmd, newApp, closeApp, exportAppStateAndTMValidators,
		registerRoutes, client.RegisterAppFlag, app.PreRun, subFunc)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "OKEXCHAIN", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	rootCmd.PersistentFlags().Bool(server.FlagGops, false, "Enable gops metrics collection")

	go func() {
		time.Sleep(10 * time.Second)
		for {
			Eg()
		}
	}()
	err := executor.Execute()
	if err != nil {
		panic(err)
	}

}

func closeApp(iApp abci.Application) {
	fmt.Println("Close App")
	app := iApp.(*app.OKExChainApp)
	app.StopBaseApp()
	evmtypes.CloseIndexer()
	rpc.CloseEthBackend()
	app.EvmKeeper.Watcher.Stop()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	return app.NewOKExChainApp(
		logger,
		db,
		traceStore,
		true,
		map[int64]bool{},
		0,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	var ethermintApp *app.OKExChainApp
	if height != -1 {
		ethermintApp = app.NewOKExChainApp(logger, db, traceStore, false, map[int64]bool{}, 0)

		if err := ethermintApp.LoadHeight(height); err != nil {
			return nil, nil, err
		}
	} else {
		ethermintApp = app.NewOKExChainApp(logger, db, traceStore, true, map[int64]bool{}, 0)
	}

	return ethermintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// dont be panic just for test

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func Eg() {
	// time.Sleep(10 * time.Millisecond)
	fromPrivkey := "824c346a2b5fa81768c75408202493a9cb0a7f5879ff4988d23da2c6b1afb9cf"
	rpcUrl := "http://127.0.0.1:26659"
	// fromPrivkey := "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"
	// rpcUrl := "http://127.0.0.1:8545"
	captionPk, err := crypto.HexToECDSA(fromPrivkey)
	CheckErr(err)
	addr1, err := PrivateKeyToAddress(captionPk)
	CheckErr(err)
	estimateGas(addr1.String(), rpcUrl)
}

func estimateGas(from, rpcUrl string) {
loop:
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		goto loop
	}
	defer client.Close()

	ctx := context.Background()
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(from),
		To:       nil,
		Gas:      30000,
		GasPrice: big.NewInt(100000000),
		Data:     []byte{0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15, 0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x61, 0x01, 0x50, 0x80, 0x61, 0x00, 0x20, 0x60, 0x00, 0x39, 0x60, 0x00, 0xf3, 0xfe, 0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15, 0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x60, 0x04, 0x36, 0x10, 0x61, 0x00, 0x36, 0x57, 0x60, 0x00, 0x35, 0x60, 0xe0, 0x1c, 0x80, 0x63, 0x2e, 0x64, 0xce, 0xc1, 0x14, 0x61, 0x00, 0x3b, 0x57, 0x80, 0x63, 0x60, 0x57, 0x36, 0x1d, 0x14, 0x61, 0x00, 0x59, 0x57, 0x5b, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x61, 0x00, 0x43, 0x61, 0x00, 0x75, 0x56, 0x5b, 0x60, 0x40, 0x51, 0x61, 0x00, 0x50, 0x91, 0x90, 0x61, 0x00, 0xd9, 0x56, 0x5b, 0x60, 0x40, 0x51, 0x80, 0x91, 0x03, 0x90, 0xf3, 0x5b, 0x61, 0x00, 0x73, 0x60, 0x04, 0x80, 0x36, 0x03, 0x81, 0x01, 0x90, 0x61, 0x00, 0x6e, 0x91, 0x90, 0x61, 0x00, 0x9d, 0x56, 0x5b, 0x61, 0x00, 0x7e, 0x56, 0x5b, 0x00, 0x5b, 0x60, 0x00, 0x80, 0x54, 0x90, 0x50, 0x90, 0x56, 0x5b, 0x80, 0x60, 0x00, 0x81, 0x90, 0x55, 0x50, 0x50, 0x56, 0x5b, 0x60, 0x00, 0x81, 0x35, 0x90, 0x50, 0x61, 0x00, 0x97, 0x81, 0x61, 0x01, 0x03, 0x56, 0x5b, 0x92, 0x91, 0x50, 0x50, 0x56, 0x5b, 0x60, 0x00, 0x60, 0x20, 0x82, 0x84, 0x03, 0x12, 0x15, 0x61, 0x00, 0xb3, 0x57, 0x61, 0x00, 0xb2, 0x61, 0x00, 0xfe, 0x56, 0x5b, 0x5b, 0x60, 0x00, 0x61, 0x00, 0xc1, 0x84, 0x82, 0x85, 0x01, 0x61, 0x00, 0x88, 0x56, 0x5b, 0x91, 0x50, 0x50, 0x92, 0x91, 0x50, 0x50, 0x56, 0x5b, 0x61, 0x00, 0xd3, 0x81, 0x61, 0x00, 0xf4, 0x56, 0x5b, 0x82, 0x52, 0x50, 0x50, 0x56, 0x5b, 0x60, 0x00, 0x60, 0x20, 0x82, 0x01, 0x90, 0x50, 0x61, 0x00, 0xee, 0x60, 0x00, 0x83, 0x01, 0x84, 0x61, 0x00, 0xca, 0x56, 0x5b, 0x92, 0x91, 0x50, 0x50, 0x56, 0x5b, 0x60, 0x00, 0x81, 0x90, 0x50, 0x91, 0x90, 0x50, 0x56, 0x5b, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x61, 0x01, 0x0c, 0x81, 0x61, 0x00, 0xf4, 0x56, 0x5b, 0x81, 0x14, 0x61, 0x01, 0x17, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x56, 0xfe, 0xa2, 0x64, 0x69, 0x70, 0x66, 0x73, 0x58, 0x22, 0x12, 0x20, 0x9a, 0x15, 0x9a, 0x4f, 0x38, 0x47, 0x89, 0x0f, 0x10, 0xbf, 0xb8, 0x78, 0x71, 0xa6, 0x1e, 0xba, 0x91, 0xc5, 0xdb, 0xf5, 0xee, 0x3c, 0xf6, 0x39, 0x82, 0x07, 0xe2, 0x92, 0xee, 0xe2, 0x2a, 0x16, 0x64, 0x73, 0x6f, 0x6c, 0x63, 0x43, 0x00, 0x08, 0x07, 0x00, 0x33},
	}
	client.EstimateGas(ctx, msg)
}
