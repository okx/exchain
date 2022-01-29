package sanity

import (
	"fmt"
	apptype "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/state"
	ttypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"testing"
)

// universeFlag used to build command
type universeFlag interface {
	// add flag to cmd
	add(cmd *cobra.Command)
	// args get formatted flags
	args() string
	// changed If the user set the value (or if left to default)
	changed() bool
}

// boolFlag bool type flag
type boolFlag struct {
	Name    string
	Default bool
	Changed bool
	Value   bool
}

func (bf *boolFlag) add(cmd *cobra.Command) {
	cmd.Flags().Bool(bf.Name, bf.Default, "")
	viper.BindPFlag(bf.Name, cmd.Flags().Lookup(bf.Name))
}

func (bf *boolFlag) args() string {
	return fmt.Sprintf("--%v=%v", bf.Name, bf.Value)
}

func (bf *boolFlag) changed() bool {
	return bf.Changed
}

// stringFlag string type flag
type stringFlag struct {
	Name    string
	Default string
	Changed bool
	Value   string
}

func (sf *stringFlag) add(cmd *cobra.Command) {
	cmd.Flags().String(sf.Name, sf.Default, "")
	viper.BindPFlag(sf.Name, cmd.Flags().Lookup(sf.Name))
}

func (sf *stringFlag) args() string {
	return fmt.Sprintf("--%v=%v", sf.Name, sf.Value)
}

func (sf *stringFlag) changed() bool {
	return sf.Changed
}

// getCommand build command by flags
func getCommand(flags []universeFlag) *cobra.Command {
	cmd := &cobra.Command{}
	var args []string
	for _, v := range flags {
		v.add(cmd)
		if v.changed() {
			args = append(args, v.args())
		}
	}
	cmd.ParseFlags(args)

	cmd.Execute()
	return cmd
}

func getCommandUserSet() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    "user-set",
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    "user-not-set",
			Default: false,
		},
	})
}

func Test_checkUserSetFlag(t *testing.T) {
	type args struct {
		cmd    *cobra.Command
		inFlag string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1. user set", args: args{cmd: getCommandUserSet(), inFlag: "user-set"}, want: true},
		{name: "2. user not set", args: args{cmd: getCommandUserSet(), inFlag: "user-not-set"}, want: false},
		{name: "3. flag not exist", args: args{cmd: getCommandUserSet(), inFlag: "flag-not-exist"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkUserSetFlag(tt.args.cmd, tt.args.inFlag); got != tt.want {
				t.Errorf("checkUserSetFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCommandSwitchNormal() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    "mixed",
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    "pure",
			Default: false,
		},
	})
}

func getCommandNodeModeRpcParalleledTx() *cobra.Command {
	return getCommand([]universeFlag{
		&stringFlag{
			Name:    apptype.FlagNodeMode,
			Default: "",
			Changed: true,
			Value:   string(apptype.RpcNode),
		},
		&boolFlag{
			Name:    state.FlagParalleledTx,
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func getCommandNodeModeRpcPruningNothing() *cobra.Command {
	return getCommand([]universeFlag{
		&stringFlag{
			Name:    apptype.FlagNodeMode,
			Default: "",
			Changed: true,
			Value:   string(apptype.RpcNode),
		},
		&stringFlag{
			Name:    server.FlagPruning,
			Default: types.PruningOptionDefault,
			Changed: true,
			Value:   types.PruningOptionNothing,
		},
	})
}

func getCommandFastQueryParalleledTx() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    watcher.FlagFastQuery,
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    state.FlagParalleledTx,
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func getCommandFastQueryPruningNothing() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    watcher.FlagFastQuery,
			Default: false,
			Changed: true,
			Value:   true,
		},
		&stringFlag{
			Name:    server.FlagPruning,
			Default: "",
			Changed: true,
			Value:   types.PruningOptionNothing,
		},
	})
}

func getCommandEnablePreruntxDownloadDelta() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    consensus.EnablePrerunTx,
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    ttypes.FlagDownloadDDS,
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func getCommandUploadDeltaParalleledTx() *cobra.Command {
	return getCommand([]universeFlag{
		&boolFlag{
			Name:    ttypes.FlagUploadDDS,
			Default: false,
			Changed: true,
			Value:   true,
		},
		&boolFlag{
			Name:    state.FlagParalleledTx,
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func TestCheckStart(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1. conflicts --fast-query and --paralleled-tx", args: args{cmd: getCommandFastQueryParalleledTx()}, wantErr: true},
		{name: "2. conflicts --fast-query and --pruning=nothing", args: args{cmd: getCommandFastQueryPruningNothing()}, wantErr: true},
		{name: "3. conflicts --enable-preruntx and --download-delta", args: args{cmd: getCommandEnablePreruntxDownloadDelta()}, wantErr: true},
		{name: "4. conflicts --upload-delta and --paralleled-tx=true", args: args{cmd: getCommandUploadDeltaParalleledTx()}, wantErr: true},
		{name: "5. conflicts --node-mod=rpc and --paralleled-tx=true", args: args{cmd: getCommandNodeModeRpcParalleledTx()}, wantErr: true},
		{name: "6. conflicts --node-mod=rpc and --pruning=nothing", args: args{cmd: getCommandNodeModeRpcPruningNothing()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = CheckStart(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("CheckStart() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(err)
		})
	}
}

func getCommandNodeModeArchiveFastQuery() *cobra.Command {
	return getCommand([]universeFlag{
		&stringFlag{
			Name:    apptype.FlagNodeMode,
			Default: "",
			Changed: true,
			Value:   string(apptype.ArchiveNode),
		},
		&boolFlag{
			Name:    watcher.FlagFastQuery,
			Default: false,
			Changed: true,
			Value:   true,
		},
	})
}

func TestCheckStartArchive(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1. conflicts --node-mod=archive and --fast-query", args: args{cmd: getCommandNodeModeArchiveFastQuery()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = CheckStart(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("CheckStart() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(err)
		})
	}
}
