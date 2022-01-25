package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	WorkspaceServerConfigName = "exchaind"
	WorkspaceClientConfigName = "exchaincli"
	WorkspaceBinaryDir        = "binary"
)

func BootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap <config>",
		Short: "fast bootstrap a testnet network for common testing from the instance.config.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("should specify the configuration file")
			}
			wd, _ := cmd.Flags().GetString(FlagWorkdir)
			workspace := fmt.Sprintf("%s/%s", wd, DefaultWorkspace)

			if _, err := os.Stat(workspace); !os.IsNotExist(err) {
				f, err := cmd.Flags().GetBool(FlagBootstrapReplace)
				if err != nil {
					return err
				}
				if !f {
					return fmt.Errorf("workspace %s is already exists", workspace)
				}
			}

			// read the instance configuration from file
			path := args[0]
			var err error
			instance, err := ReadInstanceConfig(path)
			if err != nil {
				return err
			}
			defer func() {
				if err != nil {
					for _, node := range instance.Nodes {
						os.RemoveAll(filepath.Join(workspace, node.Name))
					}
				}
			}()

			ctxs := make([]*Context, 0)
			// 1. read the node and prepare workspace && 2. build binary
			for _, node := range instance.Nodes {
				// create new context
				ctx := NewContext()
				ctxs = append(ctxs, ctx)
				ctx.Name = node.Name
				ctx.Root = filepath.Join(workspace, node.Name)
				ctx.ClientConfigName = WorkspaceClientConfigName
				ctx.ServerConfigName = WorkspaceServerConfigName

				// generate both server and client genesis configuration
				err = GenGenesisConfig(ctx)
				if err != nil {
					return err
				}

				// try to build binary use specific path
				repo, _ := cmd.Flags().GetString(FlagGitRepo)
				err = BuildBinary(wd, repo, node.Branch, filepath.Join(workspace, node.Name, WorkspaceBinaryDir))
				if err != nil {
					return err
				}
			}

			// 3. init the genesis config

			// 4. rewrite every node configuration

			// 5. fire now !

			// 6. print all configuration of the cluster

			return nil
		},
	}
	cmd.Flags().Bool(FlagBootstrapReplace, false, "replace old workspace")
	return cmd
}

func GenerateAndSaveFlags(instance *Instance, nodeName string) []string {
	flags := []string{}

	return flags
}

func ReadInstanceConfig(path string) (instance *Instance, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = fmt.Errorf("configuration file %s is not exisits ", path)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	instance = NewInstance()
	err = json.Unmarshal(buffer, instance)
	if err != nil {
		return
	}
	if !instance.Unique() {
		err = errors.New("node name is not unique please check")
		return
	}
	return instance, nil
}
