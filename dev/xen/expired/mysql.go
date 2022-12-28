package expired

import (
	"fmt"
	"github.com/okex/exchain/x/evm/statistics/mysqldb"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

const (
	flagMysqlConnection = "mysql"
	flagOutputDir       = "output"
	flagOffset          = "offset"
	flagLimit           = "limit"
	flagTTL             = "ttl"
	xenExpiredAddr      = "xen_expired.csv"
)

func Command() *cobra.Command {
	return expiredCmd
}

var expiredCmd = &cobra.Command{
	Use:   "expired_mysql",
	Short: "go the expired fss",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaim()
		return nil
	},
}

func init() {
	expiredCmd.Flags().String(flagMysqlConnection, "okc:okcpassword@(localhost:3306)/xen_stats?charset=utf8mb4&parseTime=True&loc=Local", "MySQL Database connection info")
	expiredCmd.Flags().String(flagOutputDir, "./", "useraddr output dir")
	expiredCmd.Flags().Int(flagOffset, 0, "the table claim's ID")
	expiredCmd.Flags().Int(flagLimit, 0, "the table claim's ID")
	expiredCmd.Flags().Int(flagTTL, 8, "after term time, we wait for ttl days")
	viper.BindPFlag(flagMysqlConnection, expiredCmd.Flags().Lookup(flagMysqlConnection))
	viper.BindPFlag(flagOutputDir, expiredCmd.Flags().Lookup(flagOutputDir))
	viper.BindPFlag(flagOffset, expiredCmd.Flags().Lookup(flagOffset))
	viper.BindPFlag(flagLimit, expiredCmd.Flags().Lookup(flagLimit))
	viper.BindPFlag(flagTTL, expiredCmd.Flags().Lookup(flagTTL))
}

func scanClaim() {
	offset := viper.GetInt(flagOffset)
	limit := viper.GetInt(flagLimit)
	ttl := viper.GetInt64(flagTTL)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	mysqldb.GetInstance().Init()
	db := mysqldb.GetInstance().GetGormDB()
	var claims []model.Claim
	tx := db.Table("claim").Select(
		"id", "block_time", "txhash", "useraddr", "term").Where(
		"reward=0 and id > ? and limit ?", offset, limit).Scan(&claims)
	if tx.Error != nil {
		panic(tx.Error)
	}
	for _, v := range claims {
		if time.Now().Unix() > v.BlockTime.Add(time.Duration(*v.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
			line := fmt.Sprintf("%v,%v,%v\n", v.ID, *v.Txhash, *v.Useraddr)
			_, err = f.WriteString(line)
			if err != nil {
				panic(err)
			}
		}
	}
}
