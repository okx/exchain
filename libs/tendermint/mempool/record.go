package mempool

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"sync"
)

var globalRecord *record

type sendStatus struct {
	SuccessSendCount int64
	FailSendCount    int64
	PeerHeight       int64
}

type record struct {
	logger        log.Logger
	body          sync.Map
	currentHeight int64
}

func GetGlobalRecord(l log.Logger) *record {
	if globalRecord == nil {
		globalRecord = &record{
			logger: l,
		}
	}
	globalRecord.logger = l
	return globalRecord
}

func (s *record) DoLog() {
	s.logger.Info(fmt.Sprintf("damoen log height :%d, detail : %s",  s.Detail()))
	//height is useless, delete it
	s.body.Delete(s.currentHeight)
}

func (s *record) AddPeer(peer p2p.Peer, success bool, txHeight, peerHeight int64) {
	if txHeight > s.currentHeight {
		s.currentHeight = txHeight
	}
	sendTmp := &sendStatus{
		PeerHeight: peerHeight,
	}
	if success {
		sendTmp.SuccessSendCount++
	} else {
		sendTmp.FailSendCount++
	}
	addr, _ := peer.NodeInfo().NetAddress()
	peerKey := addr.String()
	if v, ok := s.body.Load(s.currentHeight); !ok {
		var peerMap sync.Map
		peerMap.Store(peerKey, sendTmp)
		s.body.Store(s.currentHeight, peerMap)
	} else {
		//txHeight exist
		peerMap, ok := v.(sync.Map)
		if !ok {
			return
		}
		if sendInfoTmp, ok := peerMap.Load(peerKey); !ok {
			//peer not exist, store
			s.body.Store(s.currentHeight, sendTmp)
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
			s.body.Store(s.currentHeight, sendInfo)
		}
	}
}

func (s *record) DelPeer(peer p2p.Peer) {
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

func (s *record) Detail() string {
	var res string
	if v, ok := s.body.Load(s.currentHeight); !ok {
		res = fmt.Sprintf("log record curret height : %d has no tx broadcast info", s.currentHeight)
	} else {
		peerMap, ok := v.(sync.Map)
		if !ok {
			res = ""
			return res
		}
		peerMap.Range(func(k, v interface{}) bool {
			addr, ok := k.(string)
			if !ok {
				res += "peer addr is wrong"
				return false
			}
			res += " <peer : " + addr
			info, ok := v.(*sendStatus)
			if !ok {
				res += "peer sendInfo is wrong"
				return false
			}
			res += fmt.Sprintf(" , SuccessSendCount : %d", info.SuccessSendCount)
			res += fmt.Sprintf(" , FailSendCount : %d", info.FailSendCount)
			res += fmt.Sprintf(" , TxHeight : %d", s.currentHeight)
			res += fmt.Sprintf(" , PeerHeight : %d> ", info.PeerHeight)
			return true
		})
	}
	return res
}
