package state

import (
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/types"
)

func unmarshalTreeDeltaMap(input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("failed unmarshal TreeDeltaMap: empty data")
	}
	var treeDeltaMap iavl.TreeDeltaMap
	err := types.Json.Unmarshal(input, &treeDeltaMap)
	return treeDeltaMap, err
}

func marshalTreeDeltaMap(deltaMap interface{}) ([]byte, error) {
	return types.Json.Marshal(deltaMap.(iavl.TreeDeltaMap))
}

type DeltaInfo struct {
	abciResponses *ABCIResponses
	treeDeltaMap  interface{}
	watchData     interface{}

	marshalWatchData    func() ([]byte, error)
}

// for upload
func (info *DeltaInfo) dataInfo2Bytes() (types.DeltaPayload, error) {
	var err error
	pl := types.DeltaPayload{}
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

func (info *DeltaInfo) bytes2DeltaInfo(pl *types.DeltaPayload) error {
	var err error
	err = types.Json.Unmarshal(pl.ABCIRsp, &info.abciResponses)
	if err != nil {
		return err
	}

	info.treeDeltaMap, err = unmarshalTreeDeltaMap(pl.DeltasBytes)
	if err != nil {
		return err
	}
	//info.watchData, err = unmarshalData(pl.WatchBytes)
	//if err != nil {
	//	return err
	//}

	return err
}