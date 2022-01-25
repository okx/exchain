package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean [workspace]",
		Short: "clean current workspace",
		Run: func(cmd *cobra.Command, args []string) {
			workspace := DefaultWorkspace
			if len(args) == 1 {
				workspace = args[0]
			}
			path := GetWorkspacePath(cmd, workspace)
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Printf("clean failed because of %v \n", err)
				return
			}
			fmt.Println("clean successfully")
		},
	}
}
