package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/spf13/cobra"
	"io"
	"strings"
	"testing"

	"github.com/okex/exchain/libs/tendermint/libs/cli"

	tmcfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	// todo testdata?
	//	"github.com/okex/exchain/libs/cosmos-sdk/testutil/testdata"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/x/genutil"
	genutilcli "github.com/okex/exchain/libs/cosmos-sdk/x/genutil/client/cli"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp"
	simcmd "github.com/okex/exchain/libs/ibc-go/testing/simapp/simd/cmd"
)

func CreateDefaultTendermintConfig(rootDir string) (*tmcfg.Config, error) {
	conf := tmcfg.DefaultConfig()
	conf.SetRoot(rootDir)
	tmcfg.EnsureRoot(rootDir)

	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, nil
}

// BufferReader is implemented by types that read from a string buffer.
type BufferReader interface {
	io.Reader
	Reset(string)
}

// BufferWriter is implemented by types that write to a buffer.
type BufferWriter interface {
	io.Writer
	Reset()
	Bytes() []byte
	String() string
}

// ApplyMockIO replaces stdin/out/err with buffers that can be used during testing.
// Returns an input BufferReader and an output BufferWriter.
func ApplyMockIO(c *cobra.Command) (BufferReader, BufferWriter) {
	mockIn := strings.NewReader("")
	mockOut := bytes.NewBufferString("")

	c.SetIn(mockIn)
	c.SetOut(mockOut)
	c.SetErr(mockOut)

	return mockIn, mockOut
}

func ExecInitCmd(testMbm module.BasicManager, home string, cdc codec.Codec) error {
	logger := log.NewNopLogger()
	cfg, err := CreateDefaultTendermintConfig(home)
	if err != nil {
		return err
	}

	serverCtx := server.NewContext(cfg, logger)
	cmd := genutilcli.InitCmd(serverCtx, &cdc, testMbm, home)
	//	clientCtx := client.Context{}.WithCodec(cdc).WithHomeDir(home)
	clientCtx := clientCtx.CLIContext{HomeDir: home}.WithCodec(&cdc)

	_, out := ApplyMockIO(cmd)
	clientCtx = clientCtx.WithOutput(out)

	ctx := context.Background()
	// todo no ClientContextKey. what does it for?
	//	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	//	ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)

	cmd.SetArgs([]string{"appnode-test", fmt.Sprintf("--%s=%s", cli.HomeFlag, home)})

	return cmd.ExecuteContext(ctx)
}

// KeyTestPubAddr generates a new secp256k1 keypair.
func KeyTestPubAddr() (secp256k1.PrivKeySecp256k1, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())

	return key, pub, addr
}

var testMbm = module.NewBasicManager(genutil.AppModuleBasic{})

func TestAddGenesisAccountCmd(t *testing.T) {
	_, _, addr1 := KeyTestPubAddr()
	tests := []struct {
		name      string
		addr      string
		denom     string
		expectErr bool
	}{
		{
			name:      "invalid address",
			addr:      "",
			denom:     "1000atom",
			expectErr: true,
		},
		{
			name:      "valid address",
			addr:      addr1.String(),
			denom:     "1000atom",
			expectErr: false,
		},
		{
			name:      "multiple denoms",
			addr:      addr1.String(),
			denom:     "1000atom, 2000stake",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			// todo no ClientContextKey
			//logger := log.NewNopLogger()
			//cfg, err := CreateDefaultTendermintConfig(home)
			//require.NoError(t, err)

			appCodec := simapp.MakeTestEncodingConfig().Marshaler
			err := ExecInitCmd(testMbm, home, appCodec)
			require.NoError(t, err)

			// todo no ClientContextKey
			//serverCtx := server.NewContext( cfg, logger)
			//clientCtx := client.Context{}.WithJSONCodec(appCodec).WithHomeDir(home)

			ctx := context.Background()
			// todo no ClientContextKey
			//	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
			//	ctx = context.WithValue(ctx, server.ServerContextKey, serverCtx)

			cmd := simcmd.AddGenesisAccountCmd(home)
			cmd.SetArgs([]string{
				tc.addr,
				tc.denom,
				fmt.Sprintf("--%s=home", flags.FlagHome)})

			if tc.expectErr {
				require.Error(t, cmd.ExecuteContext(ctx))
			} else {
				require.NoError(t, cmd.ExecuteContext(ctx))
			}
		})
	}
}
