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

const (
	minHeight = 15405261
	maxHeight = 16359432
)

func init() {
	expiredRedisV3Cmd.Flags().String(flagRedisAddr, ":6379", "redis addr")
	expiredRedisV3Cmd.Flags().String(flagRedisPassWord, "", "redis password")
	viper.BindPFlag(flagRedisAddr, expiredRedisV3Cmd.Flags().Lookup(flagRedisAddr))
	viper.BindPFlag(flagRedisPassWord, expiredRedisV3Cmd.Flags().Lookup(flagRedisPassWord))
}

func RedisV3Command() *cobra.Command {
	return expiredRedisV3Cmd
}

var expiredRedisV3Cmd = &cobra.Command{
	Use:   "expired_redis_v3",
	Short: "get the expired xen",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimRedisV3()
		return nil
	},
}

func scanClaimRedisV3() {
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

	_, err = db.Do("SELECT", 3)
	if err != nil {
		panic(err)
	}
	var claims rediscli.XenMint

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
					content, err := redis.StringMap(db.Do("HGETALL", useraddr))
					if err != nil {
						panic(err)
					}
					for key, value := range content {
						parseClaim(&claims, key, value)
					}
					claims.UserAddr = useraddr
					if time.Now().Unix() > claims.BlockTime.Add(time.Duration(claims.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
						reward := filterReward(claims.UserAddr[1:])
						if reward == nil ||
							(reward != nil && reward.BlockTime.Unix() < claims.BlockTime.Add(time.Duration(claims.Term)*time.Duration(24)*time.Hour).Unix()) {
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

			if curse == 0 {
				return
			}
		}
	}
}
