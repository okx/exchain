//go:build rocksdb
// +build rocksdb

package main

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/system"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func dbXenCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xen",
		Short: "Statistics" + system.ChainName + " data",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			// {home}/data/*
			fromDir := ctx.Config.DBDir()
			filter("state", fromDir)

			return nil
		},
	}
	cmd.Flags().String(flagRedisAddr, ":6379", "redis addr")
	cmd.Flags().String(flagRedisPassword, "", "redis password")
	cmd.Flags().Int64(flagXenStartHeight, 0, "start height")
	cmd.Flags().Int64(flagXenEndHeight, 0, "end height")
	viper.BindPFlag(flagRedisAddr, cmd.Flags().Lookup(flagRedisAddr))
	viper.BindPFlag(flagRedisPassword, cmd.Flags().Lookup(flagRedisPassword))
	viper.BindPFlag(flagXenStartHeight, cmd.Flags().Lookup(flagXenStartHeight))
	viper.BindPFlag(flagXenEndHeight, cmd.Flags().Lookup(flagXenEndHeight))

	return cmd
}

func filter(name, fromDir string) {
	rdb, err := dbm.NewRocksDB(name, fromDir)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	startHight := viper.GetInt64(flagXenStartHeight)
	endHeight := viper.GetInt64(flagXenEndHeight)
	for i := startHight; i < endHeight; i++ {
		//		abciResponse, err := state.LoadABCIResponses(rdb, startHight)
		//		if err != nil {
		//			panic(err)
		//		}
		//		for i, tx := range abciResponse.DeliverTxs{

	}
}

//func filterLog(){
//	{
//		logsLen := len(rd.Logs)
//		for i := 0; i < logsLen; i++ {
//			if rd.Logs[i].Address.String() == "0x1cC4D981e897A3D2E7785093A648c0a75fAd0453" && // xen contract
//				len(rd.Logs[i].Topics) > 0 {
//				if rd.Logs[i].Topics[0].String() == "0xe9149e1b5059238baed02fa659dbf4bd932fbcf760a431330df4d934bc942f37" { // claimRank
//					term := big.NewInt(0).SetBytes(rd.Logs[i].Data[:32]).Int64()
//					statistics.GetInstance().SaveMintAsync(&statistics.XenMint{
//						Height:    height,
//						BlockTime: blockTime,
//						TxHash:    rd.TxHash.String(),
//						TxSender:  strings.ToLower(sender),
//						UserAddr:  hexutil.Encode(rd.Logs[i].Topics[1][12:]),
//						Term:      term,
//						Rank:      big.NewInt(0).SetBytes(rd.Logs[i].Data[32:]).String(),
//						To:        to,
//					})
//					//log.Printf("giskook %s, txsender %s,userAddress %s, term %v\n",
//					//	rd.TxHash.String(), strings.ToLower(sender), hexutil.Encode(rd.Logs[i].Topics[1][12:]), big.NewInt(0).SetBytes(rd.Logs[i].Data[:32]).Uint64())
//				} else if rd.Logs[i].Topics[0].String() == "0xd74752b13281df13701575f3a507e9b1242e0b5fb040143211c481c1fce573a6" { // claimMintRewardAndShare & claimMintReward
//					statistics.GetInstance().SaveClaimAsync(&statistics.XenClaimReward{
//						Height:       height,
//						BlockTime:    blockTime,
//						TxHash:       rd.TxHash.String(),
//						TxSender:     sender,
//						UserAddr:     hexutil.Encode(rd.Logs[i].Topics[1][12:]),
//						RewardAmount: big.NewInt(0).SetBytes(rd.Logs[i].Data[:]).String(),
//						To:           to,
//					})
//					//log.Printf("giskook %s, txsender %s,userAddress %s, reword %v\n",
//					//	rd.TxHash.String(), strings.ToLower(sender), hexutil.Encode(rd.Logs[i].Topics[1][12:]), big.NewInt(0).SetBytes(rd.Logs[i].Data[:]).Uint64())
//				}
//			}
//		}
//	}
//}
