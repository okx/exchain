package main

import (
	"github.com/okex/exchain/dev/xen/expired"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	rootCmd.AddCommand(expired.OutdatedCommand())
	rootCmd.AddCommand(expired.AnalyzeToCommand())
	rootCmd.AddCommand(expired.AnalyzeHeightCommand())
	rootCmd.AddCommand(expired.CoinToolsIndexCmd())

	rootCmd.PersistentFlags().String(expired.FlagRedisCommon, ":6379", "redis addr")
	rootCmd.PersistentFlags().String(expired.FlagRedisAuthCommon, "", "redis password")
	viper.BindPFlag(expired.FlagRedisCommon, rootCmd.PersistentFlags().Lookup(expired.FlagRedisCommon))
	viper.BindPFlag(expired.FlagRedisAuthCommon, rootCmd.PersistentFlags().Lookup(expired.FlagRedisAuthCommon))

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
