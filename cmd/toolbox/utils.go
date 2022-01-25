package main

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

func GetWorkspacePath(cmd *cobra.Command, workspace string) string {
	wd, _ := cmd.Flags().GetString(FlagWorkdir)
	return filepath.Join(wd, workspace)
}
