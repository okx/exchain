package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/spf13/cobra"
)

func StartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start action with its subcommand",
	}
	cmd.AddCommand(
		StartClusterCommand(),
	)
	return cmd
}

func StartClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster <config>",
		Short: "start a local cluster network default use current work dir as workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			instance := NewInstance()
			blob, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			err = json.Unmarshal(blob, instance)
			if err != nil {
				return err
			}

			p2pPorts := map[string]int{}
			rpcPorts := map[string]int{}
			restPorts := map[string]int{}
			for i, node := range instance.Nodes {
				p2pPorts[node.Name] = instance.Network.P2PBase + 100*i
				rpcPorts[node.Name] = instance.Network.RpcBase + 100*i
				restPorts[node.Name] = instance.Network.RpcBase + 1*i
			}

			// the seed node derive from config
			seeds := []string{}
			seedMap := map[string]bool{}
			for _, name := range instance.Network.Seeds {
				nodeKeyFile := filepath.Join(instance.Workspace, name, "exchaind", "config", "node_key.json")
				if _, err := os.Stat(nodeKeyFile); os.IsNotExist(err) {
					return err
				}
				key, err := p2p.LoadNodeKey(nodeKeyFile)
				if err != nil {
					return err
				}
				seedMap[name] = true
				seeds = append(seeds, fmt.Sprintf("%s@%s:%d", key.ID(), instance.Network.IP, p2pPorts[name]))
			}

			whitelist := []string{}
			for _, name := range instance.Network.Whitelist {
				nodeKeyFile := filepath.Join(instance.Workspace, name, "exchaind", "config", "node_key.json")
				if _, err := os.Stat(nodeKeyFile); os.IsNotExist(err) {
					return err
				}
				key, err := p2p.LoadNodeKey(nodeKeyFile)
				if err != nil {
					return err
				}
				whitelist = append(whitelist, hexutil.Encode(key.PubKey().Bytes()))
			}

			for _, node := range instance.Nodes {
				flags := []string{}
				binary := filepath.Join(instance.Workspace, node.Name, "binary", "exchaind")
				loggerFile := filepath.Join(instance.Workspace, fmt.Sprintf("%s.log", node.Name))
				if _, ok := seedMap[node.Name]; !ok {
					// set seed
					// set whitelist
					flags = []string{
						binary,
						"start",
						"--home",
						filepath.Join(instance.Workspace, node.Name, "exchaind"),
						"--p2p.seed_mode",
						"false",
						"--p2p.seeds",
						safeStringFlag(strings.Join(seeds, ",")),
						"--p2p.allow_duplicate_ip",
						"--enable-dynamic-gp=false",
						"--p2p.pex=false",
						"--p2p.addr_book_strict=false",
						"--p2p.laddr",
						safeStringFlag(fmt.Sprintf("tcp://%s:%d", instance.Network.IP, p2pPorts[node.Name])),
						"--rpc.laddr",
						safeStringFlag(fmt.Sprintf("tcp://%s:%d", instance.Network.IP, rpcPorts[node.Name])),
						"--consensus.timeout_commit",
						"600ms",
						"--log_level",
						safeStringFlag("main:debug,*:error,consensus:error,state:info,ante:info,txdecoder:info"),
						"--chain-id",
						instance.Network.ChainID,
						"--upload-delta=false",
						"--enable-gid",
						"--append-pid=true",
						"--elapsed DeliverTxs=0,Round=1,CommitRound=1,Produce=1",
						"--rest.laddr",
						safeStringFlag(fmt.Sprintf("tcp://localhost:%d", restPorts[node.Name])),
						"--enable-preruntx=false",
						fmt.Sprintf("--consensus-role=%s", node.Name),
						"--trace=true",
						"--keyring-backend test",
					}
					if node.Wtx {
						flags = append(flags, "--enable-wtx", "true")
						if node.White {
							flags = append(flags, "--mempool.node_key_whitelist", strings.Join(whitelist, ","))
						}
					} else {
						flags = append(flags, "--enable-wtx", "false")
					}
					flags = append(flags, ">", loggerFile, "2>&1 &")
				} else {
					flags = []string{
						binary,
						"start",
						"--home",
						filepath.Join(instance.Workspace, node.Name, "exchaind"),
						"--p2p.seed_mode",
						"true",
						"--p2p.allow_duplicate_ip",
						"--enable-dynamic-gp=false",
						"--p2p.pex=false",
						"--p2p.addr_book_strict=false",
						"--p2p.laddr",
						safeStringFlag(fmt.Sprintf("tcp://%s:%d", instance.Network.IP, p2pPorts[node.Name])),
						"--rpc.laddr",
						safeStringFlag(fmt.Sprintf("tcp://%s:%d", instance.Network.IP, rpcPorts[node.Name])),
						"--consensus.timeout_commit",
						"600ms",
						"--log_level",
						safeStringFlag("main:debug,*:error,consensus:error,state:info,ante:info,txdecoder:info"),
						"--chain-id",
						safeStringFlag(instance.Network.ChainID),
						"--upload-delta=false",
						"--enable-gid",
						"--append-pid=true",
						"--elapsed DeliverTxs=0,Round=1,CommitRound=1,Produce=1",
						"--rest.laddr",
						safeStringFlag(fmt.Sprintf("tcp://localhost:%d", restPorts[node.Name])),
						"--enable-preruntx=false",
						fmt.Sprintf("--consensus-role=%s", node.Name),
						"--trace=true",
						"--keyring-backend test",
					}
					flags = append(flags, ">", loggerFile, "2>&1 &")
				}
				cmd.Printf("node cmd : \n nohup %s \n", strings.Join(flags, " "))
				executable := exec.Command("nohup", flags...)
				executable.Start()
				cmd.Printf("node: %s started", node.Name)
			}

			return nil
		},
	}

	return cmd
}

func safeStringFlag(flag string) string {
	return fmt.Sprintf("\"%s\"", flag)
}
