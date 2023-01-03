package main

import (
	"github.com/okex/exchain/dev/xen/expired"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "xen",
		Short: "xen related command",
	}
	rootCmd.AddCommand(expired.Command())
	rootCmd.AddCommand(expired.RedisCommand())
	rootCmd.AddCommand(expired.RedisV2Command())
	rootCmd.AddCommand(expired.RedisV3Command())

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
