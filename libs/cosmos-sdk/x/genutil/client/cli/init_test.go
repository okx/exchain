package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	tcmd "github.com/okex/exchain/libs/tendermint/cmd/tendermint/commands"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/tests"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/x/genutil"
)

var testMbm = module.NewBasicManager(genutil.AppModuleBasic{})

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
	cmd := InitCmd(ctx, cdc, testMbm, home)

	require.NoError(t, cmd.RunE(nil, []string{"appnode-test"}))
}

func setupClientHome(t *testing.T) func() {
	clientDir, cleanup := tests.NewTestCaseDir(t)
	viper.Set(flagClientHome, clientDir)
	return cleanup
}

func TestEmptyState(t *testing.T) {
	defer server.SetupViper(t)()
	defer setupClientHome(t)()

	home, cleanup := tests.NewTestCaseDir(t)
	defer cleanup()

	logger := log.NewNopLogger()
	cfg, err := tcmd.ParseConfig()
	require.Nil(t, err)

	ctx := server.NewContext(cfg, logger)
	cdc := makeCodec()

	cmd := InitCmd(ctx, cdc, testMbm, home)
	require.NoError(t, cmd.RunE(nil, []string{"appnode-test"}))

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmd = server.ExportCmd(ctx, cdc, nil)

	err = cmd.RunE(nil, nil)
	require.NoError(t, err)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	out := <-outC

	require.Contains(t, out, "genesis_time")
	require.Contains(t, out, "chain_id")
	require.Contains(t, out, "consensus_params")
	require.Contains(t, out, "app_hash")
	require.Contains(t, out, "app_state")
}


func TestInitNodeValidatorFiles(t *testing.T) {
	home, cleanup := tests.NewTestCaseDir(t)
	defer cleanup()
	viper.Set(cli.HomeFlag, home)
	viper.Set(flags.FlagName, "moniker")
	cfg, err := tcmd.ParseConfig()
	require.Nil(t, err)
	nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(cfg)
	require.Nil(t, err)
	require.NotEqual(t, "", nodeID)
	require.NotEqual(t, 0, len(valPubKey.Bytes()))
}

// custom tx codec
func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
