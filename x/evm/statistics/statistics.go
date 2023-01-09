package statistics

import (
	"github.com/okex/exchain/x/evm/statistics/mysqldb"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"log"
	"sync"
	"time"
)

var stats *statistics

func init() {
	stats = &statistics{}
}

func GetInstance() *statistics {
	return stats
}

type Config struct {
	XenMintChanSize  int
	XenClaimChanSize int
}

type statistics struct {
	config             *Config
	chanXenMint        chan *XenMint
	chanXenClaimReward chan *XenClaimReward
	exit               chan struct{}
	initOnce           sync.Once
}

func (s *statistics) Init(config *Config) {
	s.initOnce.Do(func() {
		mysqldb.GetInstance().Init()
		mysqldb.GetInstance().GetLatestHeightAndDeleteHeight()

		rediscli.GetInstance().Init()
		rediscli.GetInstance().UpdateHeight()
		s.config = config
		s.chanXenMint = make(chan *XenMint, config.XenMintChanSize)
		s.chanXenClaimReward = make(chan *XenClaimReward, config.XenClaimChanSize)
		s.exit = make(chan struct{})
	})
}

func (s *statistics) SaveMintAsync(mint *XenMint) {
	s.chanXenMint <- mint
}

func (s *statistics) SaveClaimAsync(claim *XenClaimReward) {
	s.chanXenClaimReward <- claim
}

func (s *statistics) Do() {
	go s.doMint()
	go s.doClaim()
}

func (s *statistics) DeleteFrom(height int64) {
	mysqldb.GetInstance().DeleteFromHeight(height)
}

func (s *statistics) doMint() {
	var reward int64 = 0
	var rewardID int64 = 0
	for {
		select {
		case mint := <-s.chanXenMint:
			if mint != nil {
				rediscli.GetInstance().InsertClaim(&rediscli.XenMint{
					Height:    mint.Height,
					BlockTime: mint.BlockTime,
					TxHash:    mint.TxHash,
					TxSender:  mint.TxSender,
					UserAddr:  mint.UserAddr,
					Term:      mint.Term,
					Rank:      mint.Rank,
					To:        mint.To,
				})
				mysqldb.GetInstance().InsertClaim(model.Claim{
					Height:    &mint.Height,
					BlockTime: &mint.BlockTime,
					Txhash:    &mint.TxHash,
					Txsender:  &mint.TxSender,
					Useraddr:  &mint.UserAddr,
					Term:      &mint.Term,
					Rank:      &mint.Rank,
					Reward:    &reward,
					RewardID:  &rewardID,
				})
			}
		case <-s.exit:
			return
		}
	}
}

func (s *statistics) doClaim() {
	for {
		select {
		case claim := <-s.chanXenClaimReward:
			if claim != nil {
				rediscli.GetInstance().InsertReward(&rediscli.XenClaimReward{
					Height:       claim.Height,
					BlockTime:    claim.BlockTime,
					TxHash:       claim.TxHash,
					TxSender:     claim.TxSender,
					UserAddr:     claim.UserAddr,
					RewardAmount: claim.RewardAmount,
					To:           claim.To,
				})
				mysqldb.GetInstance().InsertReward(model.Reward{
					Height:    &claim.Height,
					BlockTime: &claim.BlockTime,
					Txhash:    &claim.TxHash,
					Txsender:  &claim.TxSender,
					Useraddr:  &claim.UserAddr,
					Amount:    &claim.RewardAmount,
				})
			}
		case <-s.exit:
			return
		}
	}
}

func (s *statistics) Close() {
	close(s.chanXenMint)
	close(s.chanXenClaimReward)
	time.Sleep(10 * time.Second)
	log.Printf("chan xen min %v chan xen reward %v\n", len(s.chanXenMint), len(s.chanXenClaimReward))

	close(s.exit)
}
