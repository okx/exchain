package server

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"

	cmn "github.com/okex/exchain/libs/tendermint/libs/os"
	"github.com/spf13/cobra"
)

// exchain full-node start flags
const (
	FlagListenAddr         = "rest.laddr"
	FlagUlockKey           = "rest.unlock_key"
	FlagUlockKeyHome       = "rest.unlock_key_home"
	FlagRestPathPrefix     = "rest.path_prefix"
	FlagCORS               = "cors"
	FlagMaxOpenConnections = "max-open"
	FlagHookstartInProcess = "startInProcess"
	FlagWebsocket          = "wsport"
	FlagWsMaxConnections   = "ws.max_connections"
	FlagWsSubChannelLength = "ws.sub_channel_length"
)

//module hook

type fnHookstartInProcess func(ctx *Context) error

type serverHookTable struct {
	hookTable map[string]interface{}
}

var gSrvHookTable = serverHookTable{make(map[string]interface{})}

func InstallHookEx(flag string, hooker fnHookstartInProcess) {
	gSrvHookTable.hookTable[flag] = hooker
}

//call hooker function
func callHooker(flag string, args ...interface{}) error {
	params := make([]interface{}, 0)
	switch flag {
	case FlagHookstartInProcess:
		{
			//none hook func, return nil
			function, ok := gSrvHookTable.hookTable[FlagHookstartInProcess]
			if !ok {
				return nil
			}
			params = append(params, args...)
			if len(params) != 1 {
				return errors.New("too many or less parameter called, want 1")
			}

			//param type check
			p1, ok := params[0].(*Context)
			if !ok {
				return errors.New("wrong param 1 type. want *Context, got" + reflect.TypeOf(params[0]).String())
			}

			//get hook function and call it
			caller := function.(fnHookstartInProcess)
			return caller(p1)
		}
	default:
		break
	}
	return nil
}

//end of hook

func setPID(ctx *Context) {
	pid := os.Getpid()
	f, err := os.OpenFile(filepath.Join(ctx.Config.RootDir, "config", "pid"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		cmn.Exit(err.Error())
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Println(err.Error())
	}
	writer.Flush()
}

// StopCmd stop the node gracefully
// Tendermint.
func StopCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node gracefully",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(filepath.Join(ctx.Config.RootDir, "config", "pid"))
			if err != nil {
				errStr := fmt.Sprintf("%s Please finish the process of exchaind through kill -2 pid to stop gracefully", err.Error())
				cmn.Exit(errStr)
			}
			defer f.Close()
			in := bufio.NewScanner(f)
			in.Scan()
			pid, err := strconv.Atoi(in.Text())
			if err != nil {
				errStr := fmt.Sprintf("%s Please finish the process of exchaind through kill -2 pid to stop gracefully", err.Error())
				cmn.Exit(errStr)
			}
			process, err := os.FindProcess(pid)
			if err != nil {
				cmn.Exit(err.Error())
			}
			err = process.Signal(os.Interrupt)
			if err != nil {
				cmn.Exit(err.Error())
			}
			fmt.Println("pid", pid, "has been sent SIGINT")
			return nil
		},
	}
	return cmd
}

var sem *nodeSemaphore

type nodeSemaphore struct {
	done chan struct{}
}

func Stop() {
	sem.done <- struct{}{}
}

// registerRestServerFlags registers the flags required for rest server
func registerRestServerFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(FlagListenAddr, "tcp://0.0.0.0:26659", "EVM RPC and cosmos-sdk REST API listen address.")
	cmd.Flags().String(FlagUlockKey, "", "Select the keys to unlock on the RPC server")
	cmd.Flags().String(FlagUlockKeyHome, os.ExpandEnv("$HOME/.exchaincli"), "The keybase home path")
	cmd.Flags().String(FlagRestPathPrefix, "exchain", "Path prefix for registering rest api route.")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(FlagCORS, "", "Set the rest-server domains that can make CORS requests (* for all)")
	cmd.Flags().Int(FlagMaxOpenConnections, 1000, "The number of maximum open connections of rest-server")
	cmd.Flags().String(FlagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().Int(FlagWsMaxConnections, 20000, "the max capacity number of websocket client connections")
	cmd.Flags().Int(FlagWsSubChannelLength, 100, "the length of subscription channel")
	cmd.Flags().String(flags.FlagChainID, "", "Chain ID of tendermint node for web3")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block) for web3")
	return cmd
}

func nodeModeCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node-mode",
		Short: "exchaind start --node-mode help info",
		Long: `There are three node modes that can be set when the exchaind start
set --node-mode=rpc to manage the following flags:
	--disable-checktx-mutex=true
	--disable-query-mutex=true
	--enable-bloom-filter=true
	--fast-lru=10000
	--fast-query=true
	--iavl-enable-async-commit=true
	--max-open=20000
	--mempool.enable_pending_pool=true
	--cors=*

set --node-mode=validator to manage the following flags:
	--disable-checktx-mutex=true
	--disable-query-mutex=true
	--enable-dynamic-gp=false
	--iavl-enable-async-commit=true
	--iavl-cache-size=10000000
	--pruning=everything

set --node-mode=archive to manage the following flags:
	--pruning=nothing
	--disable-checktx-mutex=true
	--disable-query-mutex=true
	--enable-bloom-filter=true
	--fast-lru=10000
	--iavl-enable-async-commit=true
	--max-open=20000
	--cors=*`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}
