package mempool

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"sync"
)

var globalRecord *recorder

type sendStatus struct {
	PeerKey          string
	PeerHeight       int64
	SuccessSendCount int64
	FailSendCount    int64
}

type recorder struct {
	logger        log.Logger
	body          sync.Map
	currentHeight int64
}

func GetGlobalRecord(l log.Logger) *recorder {
	if globalRecord == nil {
		globalRecord = &recorder{
			logger: l,
		}
	}
	globalRecord.logger = l
	return globalRecord
}

func (s *recorder) DoLog() {
	if s.currentHeight == 0 {
		s.logger.Info("mp broadcast log is empty, no mp tx broadcast send log")
		return
	}

	s.logger.Info(fmt.Sprintf("mp broadcast log height :%d, detail : %s", s.currentHeight, s.Detail()))
	//height is useless, delete it
	s.body.Delete(s.currentHeight)
	s.currentHeight = 0
}

func (s *recorder) AddPeer(peer p2p.Peer, success bool, txHeight, peerHeight int64) {
	if txHeight > s.currentHeight {
		s.currentHeight = txHeight
	}
	var peerMap sync.Map
	addr, _ := peer.NodeInfo().NetAddress()
	peerKey := addr.String()

	sendTmp := &sendStatus{
		PeerKey:    peerKey,
		PeerHeight: peerHeight,
	}
	if success {
		sendTmp.SuccessSendCount++
	} else {
		sendTmp.FailSendCount++
	}
	peerMap.Store(peerKey, sendTmp)

	if v, ok := s.body.Load(s.currentHeight); !ok {
		s.body.Store(s.currentHeight, peerMap)
	} else {
		//txHeight exist
		peerMap, ok := v.(sync.Map)
		if !ok {
			return
		}
		if sendInfoTmp, ok := peerMap.Load(peerKey); !ok {
			peerMap.Store(peerKey, sendTmp)
			s.body.Store(s.currentHeight, peerMap)
		} else {
			sendInfo, ok := sendInfoTmp.(*sendStatus)
			if !ok {
				return
			}
			sendInfo.PeerHeight = peerHeight
			if success {
				sendInfo.SuccessSendCount++
			} else {
				sendInfo.FailSendCount++
			}

			peerMap.Store(peerKey, sendInfo)
			s.body.Store(s.currentHeight, peerMap)
		}
	}
}

func (s *recorder) DelPeer(peer p2p.Peer) {
	//delete peer from current height
	if v, ok := s.body.Load(s.currentHeight); ok {
		peerMap, ok := v.(sync.Map)
		if !ok {
			return
		}
		addr, _ := peer.NodeInfo().NetAddress()
		peerKey := addr.String()
		peerMap.Delete(peerKey)
	}
}

func (s *recorder) Detail() string {
	var res string
	var sends []sendStatus

	if v, ok := s.body.Load(s.currentHeight); !ok {
		res = fmt.Sprintf("log record curret height : %d has no tx broadcast info", s.currentHeight)
	} else {
		peerMap, ok := v.(sync.Map)
		if !ok {
			res = "peerMap type wrong"
			return res
		}
		peerMap.Range(func(k, v interface{}) bool {
			info, ok := v.(*sendStatus)
			if !ok {
				res += "peer sendInfo is wrong"
				return false
			}

			sends = append(sends, sendStatus{
				PeerKey:          info.PeerKey,
				PeerHeight:       s.currentHeight,
				SuccessSendCount: info.SuccessSendCount,
				FailSendCount:    info.FailSendCount,
			})

			return true
		})
	}

	if len(res) != 0 {
		return res
	}
	sendsJ, _ := json.Marshal(sends)
	res = string(sendsJ)
	return res
}
