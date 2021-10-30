package cli

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/dependence/cosmos-sdk/client/flags"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/types"

	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	"github.com/okex/exchain/dependence/cosmos-sdk/server"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/types/module"
	"github.com/okex/exchain/x/genutil"
)

const (
	flagOverwrite  = "overwrite"
	flagClientHome = "home-client"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))
	return err
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager,
	defaultNodeHome string) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators' and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(flags.FlagChainID)

			nodeID, _, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = args[0]
			genFile := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			appState, err := codec.MarshalJSONIndent(cdc, mbm.DefaultGenesis())
			if err != nil {
				return err
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return err
				}
			}

			genDoc.ChainID, genDoc.Validators, genDoc.AppState = chainID, nil, appState
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return err
			}
			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "testchain-1", "genesis file chain-id, it's necessary to be provided in the format like \"[chain name]-[positive integer]\"")

	return cmd
}
