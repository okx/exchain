package mempool

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"sync"
	"time"
)

var globalRecord *record

type sendStatus struct {
	SuccessSendCount int64
	FailSendCount    int64
	TxHeight         int64
	PeerHeight       int64
}

type record struct {
	lock   sync.RWMutex
	logger log.Logger
	body   map[p2p.Peer]*sendStatus `json:"detailInfo"`
}

func GetGlobalRecord(l log.Logger) *record {
	if globalRecord == nil {
		globalRecord = &record{
			logger : l,
			body: make(map[p2p.Peer]*sendStatus),
		}
		//采取定期打印log的方式
		go globalRecord.GoLog()

	}
	return globalRecord
}

func (s *record) GoLog()  {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <- ticker.C:
			s.logger.Info(fmt.Sprintf("damoen log : %s"), s.Detail())
		}
	}
}

func (s *record) AddPeer(peer p2p.Peer, success bool, txHeight, peerHeight int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.body[peer]; !ok {
		s.body[peer] = &sendStatus{
			TxHeight:   txHeight,
			PeerHeight: peerHeight,
		}
	} else {
		s.body[peer].TxHeight = txHeight
		s.body[peer].PeerHeight = peerHeight
	}

	if success {
		s.body[peer].SuccessSendCount++
	} else {
		s.body[peer].FailSendCount++
	}

}

func (s *record) DelPeer(peer p2p.Peer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.body[peer]; !ok {
		delete(s.body, peer)
	}
}

func (s *record) Detail() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var res string
	for k, v := range s.body {
		addr, _ := k.NodeInfo().NetAddress()
		res += " <peer : " + addr.String()
		res += fmt.Sprintf(" , SuccessSendCount : %d", v.SuccessSendCount)
		res += fmt.Sprintf(" , FailSendCount : %d", v.FailSendCount)
		res += fmt.Sprintf(" , TxHeight : %d", v.TxHeight)
		res += fmt.Sprintf(" , PeerHeight : %d> ", v.PeerHeight)
	}

	return res
}
