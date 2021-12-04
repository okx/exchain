package datacenter_cgi

import (
	"bytes"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"io/ioutil"
	"net/http"
)


// DataCenterMsg is the format of Data send to DataCenter
type DataCenterMsg struct {
	Height    int64  `json:"height"`
	Block     []byte `json:"block"`
	Delta     []byte `json:"delta"`
	WatchData []byte `json:"watch_data"`
}

// sendToDatacenter send bcBlockResponseMessage to DataCenter
func sendToDatacenter(logger log.Logger, block *types.Block, deltas *types.Deltas, wd *types.WatchData) {
	var blockBytes, deltaBytes, wdBytes []byte
	var err error
	if block != nil {
		if blockBytes, err = block.Marshal(); err != nil {
			return
		}
	}
	if deltas != nil {
		if deltaBytes, err = deltas.Marshal(); err != nil {
			return
		}
	}
	if wd != nil && wd.Size() > 0{
		wdBytes = wd.WatchDataByte
	}

	msg := DataCenterMsg{block.Height, blockBytes, deltaBytes, wdBytes}
	msgBody, err := types.Json.Marshal(&msg)
	if err != nil {
		return
	}
	response, err := http.Post(types.GetCenterUrl() + "save", "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		logger.Error("sendToDatacenter err ,", err)
		return
	}
	defer response.Body.Close()
}

// getDataFromDatacenter send bcBlockResponseMessage to DataCenter
func getDeltaFromDatacenter(logger log.Logger, height int64) (*types.Deltas, error) {
	msg := DataCenterMsg{Height: height}
	msgBody, err := types.Json.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(types.GetCenterUrl() + "loadDelta", "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		logger.Error("getDataFromDatacenter err ,", err)
		return nil, err
	}

	defer response.Body.Close()
	rlt, _ := ioutil.ReadAll(response.Body)
	logger.Info("GetDataFromDatacenter", "height", height, "len", len(rlt))

	delta := &types.Deltas{}
	if delta.Unmarshal(rlt) != nil {
		return nil, err
	}

	return delta, nil
}
