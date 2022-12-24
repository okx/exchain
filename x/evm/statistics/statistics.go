package statistics

import (
	"log"
	"sync"
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
	config       *Config
	chanXenMint  chan *XenMint
	chanXenClaim chan *XenClaimReward
	exit         chan struct{}
	initOnce     sync.Once
}

func (s *statistics) Init(config *Config) {
	s.initOnce.Do(func() {
		s.config = config
		s.chanXenMint = make(chan *XenMint, config.XenMintChanSize)
		s.chanXenClaim = make(chan *XenClaimReward, config.XenClaimChanSize)
		s.exit = make(chan struct{})
	})
}

func (s *statistics) SaveMintAsync(mint *XenMint) {
	//	s.chanXenMint <- mint
}

func (s *statistics) SaveClaimAsync(claim *XenClaimReward) {
	//	s.chanXenClaim <- claim
}

func (s *statistics) Do() {
	go s.doMint()
	go s.doClaim()
}

func (s *statistics) doMint() {
	for {
		select {
		case mint := <-s.chanXenMint:
			log.Println(mint)
			return
		case <-s.exit:
			return
		}
	}
}

func (s *statistics) doClaim() {
	for {
		select {
		case claim := <-s.chanXenClaim:
			log.Println(claim)
			return
		case <-s.exit:
			return
		}
	}
}
