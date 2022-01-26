package state

import (
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/types"
)

func unmarshalTreeDeltaMap(input []byte) (iavl.TreeDeltaMap, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("failed to unmarshal TreeDeltaMap: empty data")
	}
	var treeDeltaMap iavl.TreeDeltaMap
	err := types.Json.Unmarshal(input, &treeDeltaMap)
	return treeDeltaMap, err
}

func marshalTreeDeltaMap(deltaMap iavl.TreeDeltaMap) ([]byte, error) {
	return types.Json.Marshal(deltaMap)
}

type DeltaInfo struct {
	abciResponses *ABCIResponses
	treeDeltaMap  iavl.TreeDeltaMap
	//watchData     interface{}

	marshalWatchData    func() ([]byte, error)
}

// for upload
func (info *DeltaInfo) dataInfo2Bytes() (pl types.DeltaPayload, err error) {
	pl.ABCIRsp, err = types.Json.Marshal(info.abciResponses)
	if err != nil {
		return pl, err
	}

	pl.DeltasBytes, err = marshalTreeDeltaMap(info.treeDeltaMap)
	if err != nil {
		return pl, err
	}

	pl.WatchBytes, err = info.marshalWatchData()
	if err != nil {
		return pl, err
	}

	return pl, err
}

// for download
func (info *DeltaInfo) bytes2DeltaInfo(pl *types.DeltaPayload) (err error) {
	err = types.Json.Unmarshal(pl.ABCIRsp, &info.abciResponses)
	if err != nil {
		return err
	}

	info.treeDeltaMap, err = unmarshalTreeDeltaMap(pl.DeltasBytes)
	if err != nil {
		return err
	}

	// todo unmarshal watchData in download thread
	//info.watchData, err = unmarshalData(pl.WatchBytes)
	//if err != nil {
	//	return err
	//}

	return err
}