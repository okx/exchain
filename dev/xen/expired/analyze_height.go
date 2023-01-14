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
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(2)
	ret := make(map[int64]int)
	heightBlockTime := make(map[int64]int64)

	claimsRet := make(map[int64]int)
	claimsHeightBlockTime := make(map[int64]int64)

	go func() {
		defer wg.Done()
		getMints(ret, heightBlockTime)
	}()
	go func() {
		defer wg.Done()
		getClaims(claimsRet, claimsHeightBlockTime)
	}()
	wg.Wait()

	printHeightStatistics(ret, heightBlockTime, claimsRet, claimsHeightBlockTime, f)
}

func getMints(ret map[int64]int, heightBlockTime map[int64]int64) {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	// mint counter
	_, err := db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}

	var curse = 0
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
				return ret, heightBlockTime
			}
		}
	}
}

func getClaims(ret map[int64]int, heightBlockTime map[int64]int64) {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	// mint counter
	_, err := db.Do("SELECT", 1)
	if err != nil {
		panic(err)
	}

	var curse = 0
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
					claims := getReward(useraddr)
					ret[claims.Height]++
					heightBlockTime[claims.Height] = claims.BlockTime.UnixNano()
				}
			}

			if curse == 0 {
				return ret, heightBlockTime
			}
		}
	}
}

// toStatistics height <-> count
// mintHeightBlockTime  height <-> blockTime
// claimStatistics  height <-> count
// claimHeightBlockTime height <-> blockTime
func printHeightStatistics(toStatistics map[int64]int, mintHeightBlockTime map[int64]int64,
	claimStatistics map[int64]int, claimHeightBlockTime map[int64]int64, f *os.File) {
	// mint
	sortedTo := make([]int64, 0, len(toStatistics))
	for k, _ := range toStatistics {
		sortedTo = append(sortedTo, k)
	}
	sort.Slice(sortedTo, func(i, j int) bool {
		return sortedTo[i] < sortedTo[j]
	})

	// claim
	claimHeight := make([]int64, 0, len(claimStatistics))
	for k, _ := range claimStatistics {
		claimHeight = append(sortedTo, k)
	}
	sort.Slice(claimHeight, func(i, j int) bool {
		return claimHeight[i] < claimHeight[j]
	})
	minHeight := min(sortedTo[0], claimHeight[0])
	maxHeight := max(sortedTo[len(sortedTo)-1], claimHeight[len(claimHeight)-1])

	counter := 1

	// toStatistics height <-> count
	// mintHeightBlockTime  height <-> blockTime
	// claimStatistics  height <-> count
	// claimHeightBlockTime height <-> blockTime
	var index int
	var height int64
	var blockTime int64
	var mintCount int
	var claimCount int
	var totalMint int
	var totalClaim int
	var ok bool
	for i := minHeight; i <= maxHeight; i++ {
		ok = false
		blockTime, ok = mintHeightBlockTime[i]
		if !ok {
			blockTime, ok = claimHeightBlockTime[i]
		}
		if ok {
			index = counter
			height = i
			mintCount = toStatistics[i]
			claimCount = claimStatistics[i]
			totalMint += mintCount
			totalClaim += claimCount
			line := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", index, height, blockTime, mintCount, claimCount, totalMint, totalClaim)
			_, err := f.WriteString(line)
			if err != nil {
				panic(err)
			}

			counter++
		}
	}
}

func min(i, j int64) int64 {
	if i < j {
		return i
	}
	return j
}

func max(i, j int64) int64 {
	if i > j {
		return i
	}
	return j
}
