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
	analyzeOutputFile   = "./analyze_to.csv"
	flagRedisAddrTo     = "redis_addr_to"
	flagRedisPassWordTo = "redis_auth_to"
)

func init() {
	analyzeToCmd.Flags().String(flagRedisAddrTo, ":6379", "redis addr")
	analyzeToCmd.Flags().String(flagRedisPassWordTo, "", "redis password")
	viper.BindPFlag(flagRedisAddrTo, analyzeToCmd.Flags().Lookup(flagRedisAddrTo))
	viper.BindPFlag(flagRedisPassWordTo, analyzeToCmd.Flags().Lookup(flagRedisPassWordTo))
}

func AnalyzeToCommand() *cobra.Command {
	return analyzeToCmd
}

var analyzeToCmd = &cobra.Command{
	Use:   "analyzeTo",
	Short: "analyze to address",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimTo()
		return nil
	},
}

func scanClaimTo() {
	filename := filepath.Join(viper.GetString(flagOutputDir), analyzeOutputFile)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	global.RedisAddr = viper.GetString(flagRedisAddrTo)
	global.RedisPassword = viper.GetString(flagRedisPassWordTo)
	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	_, err = db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]int)

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
					ret[claims.To]++
				}
			}

			if curse == 0 {
				printStatistics(ret, f)
				return
			}
		}
	}
	printStatistics(ret, f)
}

func getReward(key string) *rediscli.XenClaimReward {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 1)
	if err != nil {
		panic(err)
	}
	var ret rediscli.XenClaimReward

	content, err := redis.StringMap(db.Do("HGETALL", key))
	if err != nil || len(content) == 0 {
		panic(fmt.Sprintf("%v map %v error %v", key, len(content), err))
	}
	for key, value := range content {
		parseReward(&ret, key, value)
	}
	tmp := strings.Split(key, "_")
	ret.UserAddr = tmp[0][1:]

	return &ret
}

func getClaim(key string) *rediscli.XenMint {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}
	var ret rediscli.XenMint

	content, err := redis.StringMap(db.Do("HGETALL", key))
	if err != nil || len(content) == 0 {
		panic(fmt.Sprintf("%v map %v error %v", key, len(content), err))
	}
	for key, value := range content {
		parseClaim(&ret, key, value)
	}
	tmp := strings.Split(key, "_")
	ret.UserAddr = tmp[0][1:]

	return &ret
}

func printStatistics(toStatistics map[string]int, f *os.File) {
	sortedTo := make([]string, 0, len(toStatistics))
	for k, _ := range toStatistics {
		sortedTo = append(sortedTo, k)
	}
	sort.Slice(sortedTo, func(i, j int) bool {
		return toStatistics[sortedTo[i]] > toStatistics[sortedTo[j]]
	})

	counter := 1
	for _, v := range sortedTo {
		line := fmt.Sprintf("%v,%v,%v\n", counter, v, toStatistics[v])
		_, err := f.WriteString(line)
		if err != nil {
			panic(err)
		}
		counter++
	}
}
