package expired

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	flagRedisAddr     = "redis_addr"
	flagRedisPassWord = "redis_auth"
)

func init() {
	expiredRedisV2Cmd.Flags().String(flagRedisAddr, ":6379", "redis addr")
	expiredRedisV2Cmd.Flags().String(flagRedisPassWord, "", "redis password")
	viper.BindPFlag(flagRedisAddr, expiredCmd.Flags().Lookup(flagRedisAddr))
	viper.BindPFlag(flagRedisPassWord, expiredCmd.Flags().Lookup(flagRedisPassWord))
}

func RedisV2Command() *cobra.Command {
	return expiredRedisCmd
}

var expiredRedisV2Cmd = &cobra.Command{
	Use:   "expired_redis_parallel",
	Short: "get the expired xen",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimRedisV2()
		return nil
	},
}

func scanClaimRedisV2() {
	ttl := viper.GetInt64(flagTTL)
	log.Println("ttl", ttl)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	var claims rediscli.XenMint

	counter := 0
	curse := viper.GetInt(flagCursor)
	for {
		_, err := db.Do("SELECT", 0)
		if err != nil {
			panic(err)
		}
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

func filterReward(userAddr string) *rediscli.XenClaimReward {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	_, err := db.Do("SELECT", 1)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.Do("SELECT", 0)
		db.Close()
	}()

	ua := "r" + userAddr
	exists, _ := redis.Int(db.Do("EXISTS", ua))
	if exists == 0 {
		return nil
	}
	content, err := redis.StringMap(db.Do("HGETALL", ua))
	if err != nil {
		panic(err)
	}
	var reward rediscli.XenClaimReward
	for key, value := range content {
		parseReward(&reward, key, value)
	}
	return &reward
}

func parseReward(reward *rediscli.XenClaimReward, key, value string) {
	switch key {
	case "height":
		height, _ := strconv.Atoi(value)
		reward.Height = int64(height)
	case "txhash":
		reward.TxHash = value
	case "btime":
		utc, _ := strconv.Atoi(value)
		tim := time.Unix(int64(utc), 0)
		reward.BlockTime = tim
	}
}
