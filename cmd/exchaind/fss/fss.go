package fss

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return fssCmd
}

var fssCmd = &cobra.Command{
	Use:   "fss",
	Short: "FSS is an auxiliary fast storage system to IAVL",
	Long: `IAVL fast storage related commands:
This command include a set of command of the IAVL fast storage.
include create sub command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}
