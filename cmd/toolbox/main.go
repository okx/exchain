package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "toolbox",
		Short: "command line tools for bootstrap a new testnet cluster and do some automation",
	}

	wd, _ := os.Getwd()

	root.PersistentFlags().String(FlagWorkdir, wd, "specific the workdir for this tool")
	root.PersistentFlags().String(FlagGitRepo, "https://github.com/okex/exchain.git", "specific the exchain repo")
	root.PersistentFlags().String(FlagGitBranch, "dev", "specific the exchain repo")

	root.AddCommand(
		ScoffldCommand(),
		BuildCommand(),
		BootstrapCommand(),
		StartCommand(),
		CleanCommand(),
		KillCommand(),
	)

	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
