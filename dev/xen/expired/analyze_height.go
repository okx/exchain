package expired

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	analyzeHeightOutputFile = "./analyze_height.csv"
	flagRedisAddrHeight     = "redis_addr_height"
	flagRedisPassWordHeight = "redis_auth_height"
)

func init() {
	analyzeHeightCmd.Flags().String(flagRedisAddrHeight, ":6379", "redis addr")
	analyzeHeightCmd.Flags().String(flagRedisPassWordHeight, "", "redis password")
	viper.BindPFlag(flagRedisAddrHeight, analyzeHeightCmd.Flags().Lookup(flagRedisAddrHeight))
	viper.BindPFlag(flagRedisPassWordHeight, analyzeHeightCmd.Flags().Lookup(flagRedisPassWordHeight))
}

func AnalyzeHeightCommand() *cobra.Command {
	return analyzeHeightCmd
}

var analyzeHeightCmd = &cobra.Command{
	Use:   "analyzeHeight",
	Short: "analyze height",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimHeight()
		return nil
	},
}

func scanClaimHeight() {
	filename := filepath.Join(viper.GetString(flagOutputDir), analyzeHeightOutputFile)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	global.RedisAddr = viper.GetString(flagRedisAddrHeight)
	global.RedisPassword = viper.GetString(flagRedisPassWordHeight)
	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	_, err = db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}

	ret := make(map[int64]int)

	curse := viper.GetInt(flagCursor)
	for {
		curseValues, err := redis.Values(db.Do("SCAN", curse))
		if err != nil {
			panic(err)
		}
		if len(curseValues) > 0 {
			curse, err = redis.Int(curseValues[0], err)
			if err != nil {
				panic(fmt.Sprintf("got curse error %v %v", curseValues[0], err))
			}

			values, err := redis.Values(curseValues[1], err)
			if err != nil {
				panic(fmt.Sprintf("get values error %v %v", curseValues[1], err))
			}
			for _, v := range values {
				useraddr, err := redis.String(v, nil)
				if err != nil {
					panic(err)
				}
				if strings.Contains(useraddr, "0x") {
					claims := getClaim(useraddr)
					ret[claims.Height]++
				}
			}

			if curse == 0 {
				printHeightStatistics(ret, f)
				return
			}
		}
	}
}

func printHeightStatistics(toStatistics map[int64]int, f *os.File) {
	sortedTo := make([]int64, 0, len(toStatistics))
	for k, _ := range toStatistics {
		sortedTo = append(sortedTo, k)
	}
	sort.Slice(sortedTo, func(i, j int) bool {
		return sortedTo[i] < sortedTo[j]
	})

	total := 0
	counter := 1
	for _, v := range sortedTo {
		total += toStatistics[v]
		line := fmt.Sprintf("%v,%v,%v,%v\n", counter, v, toStatistics[v], total)
		_, err := f.WriteString(line)
		if err != nil {
			panic(err)
		}
		counter++
	}
}
