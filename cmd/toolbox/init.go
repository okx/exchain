package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

func ScoffldCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scoffld [number]",
		Short: "scoffld a new specificed number node cluster config",
		RunE: func(cmd *cobra.Command, args []string) error {
			number := 4
			if len(args) == 1 {
				n, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				number = n
			}
			wd, _ := cmd.Flags().GetString(FlagWorkdir)
			path := filepath.Join(wd, "scoffld.json")
			instance := NewInstance()
			instance.Name = "scoffld-configuration"
			instance.Description = "scoffld local cluster config"
			instance.Workspace = filepath.Join(wd, "workspace")
			os.MkdirAll(filepath.Join(wd, "workspace"), 0755)
			instance.Network.Validators = []string{}
			for i := 0; i < number; i++ {
				name := fmt.Sprintf("n%d", i)
				instance.Network.Validators = append(instance.Network.Validators, name)
			}
			for _, name := range instance.Network.Validators {
				instance.Nodes = append(instance.Nodes, Node{
					Name:   name,
					Branch: "dev",
					Flags:  []string{},
				})
			}
			document, _ := json.MarshalIndent(instance, "", "\t")
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			file.Write(document)
			return nil
		},
	}

	return cmd
}
