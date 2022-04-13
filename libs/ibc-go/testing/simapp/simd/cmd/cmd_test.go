package cmd_test

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/tests"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/genutil/client/cli"
	tcmd "github.com/okex/exchain/libs/tendermint/cmd/tendermint/commands"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	flagOverwrite  = "overwrite"
	flagClientHome = "home-client"
)

//func TestInitCmd(t *testing.T) {
//	rootCmd, _ := cmd.NewRootCmd()
//	rootCmd.SetArgs([]string{
//		"init",        // Test the init cmd
//		"simapp-test", // Moniker
//		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
//	})
//
//	require.NoError(t, svrcmd.Execute(rootCmd, simapp.DefaultNodeHome))
//}
func TestInitCmd(t *testing.T) {
	defer server.SetupViper(t)()
	defer setupClientHome(t)()
	home, cleanup := tests.NewTestCaseDir(t)
	defer cleanup()

	logger := log.NewNopLogger()
	cfg, err := tcmd.ParseConfig()
	require.Nil(t, err)

	ctx := server.NewContext(cfg, logger)
	cdc := makeCodec()
	cmd := cli.InitCmd(ctx, cdc, testMbm, home)

	require.NoError(t, cmd.RunE(nil, []string{"appnode-test"}))
}

func setupClientHome(t *testing.T) func() {
	clientDir, cleanup := tests.NewTestCaseDir(t)
	viper.Set(flagClientHome, clientDir)
	viper.Set(flagOverwrite, true)
	return cleanup
}

// custom tx codec
func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
