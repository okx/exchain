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
	"strings"
	"time"
)

func init() {
	expiredOutdatedCmd.Flags().String(flagRedisAddr, ":6379", "redis addr")
	expiredOutdatedCmd.Flags().String(flagRedisPassWord, "", "redis password")
	viper.BindPFlag(flagRedisAddr, expiredOutdatedCmd.Flags().Lookup(flagRedisAddr))
	viper.BindPFlag(flagRedisPassWord, expiredOutdatedCmd.Flags().Lookup(flagRedisPassWord))
}

func OutdatedCommand() *cobra.Command {
	return expiredOutdatedCmd
}

var expiredOutdatedCmd = &cobra.Command{
	Use:   "expired_outdated",
	Short: "get the expired xen",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanOutdated()
		return nil
	},
}

func scanOutdated() {
	ttl := viper.GetInt64(flagTTL)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	global.RedisAddr = viper.GetString(flagRedisAddr)
	global.RedisPassword = viper.GetString(flagRedisPassWord)
	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	_, err = db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}

	counter := 0
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
					eqUserAddr := checkMintRewardIfEqualReturnUserAddr(useraddr)
					if eqUserAddr != "" {
						claims := getLatestClaimEx(eqUserAddr)
						reward := getLatestReward(eqUserAddr)
						if reward.BlockTime.Unix() > claims.BlockTime.Add(time.Duration(claims.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
							counter++
							line := fmt.Sprintf("%v,%v,%v\n", counter, claims.TxHash, claims.UserAddr)
							_, err = f.WriteString(line)
							if err != nil {
								panic(err)
							}
						}
					}
				}
			}
		}

		if curse == 0 {
			return
		}
	}
}

func checkMintRewardIfEqualReturnUserAddr(userAddr string) string {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}
	countMint, err := redis.Int(db.Do("ZCOUNT", userAddr, minHeight, maxHeight))
	if err != nil {
		panic(err)
	}

	_, err = db.Do("SELECT", 3)
	if err != nil {
		panic(err)
	}

	rewardUserAddr := "r" + userAddr[1:]
	countReward, err := redis.Int(db.Do("ZCOUNT", rewardUserAddr, minHeight, maxHeight))
	if err != nil {
		panic(err)
	}
	if countReward > countMint {
		panic(fmt.Sprintf("%v impossible!", userAddr))
	}
	if countReward != countMint {
		return ""
	}

	return userAddr[1:]
}

func getLatestClaimEx(userAddr string) *rediscli.XenMint {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}

	v, err := redis.Strings(db.Do("ZREVRANGE", "c"+userAddr, 0, -1))
	if err != nil {
		panic(err)
	}
	if len(v) == 0 {
		return nil
	}
	latestRewardKey := "c" + userAddr + "_" + v[0]

	content, err := redis.StringMap(db.Do("HGETALL", latestRewardKey))
	if err != nil {
		panic(err)
	}
	var mint rediscli.XenMint
	for key, value := range content {
		parseClaim(&mint, key, value)
	}

	return &mint
}
