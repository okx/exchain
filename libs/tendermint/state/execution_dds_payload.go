package state

import "github.com/okex/exchain/libs/tendermint/types"

var(
	unmarshalTreeDeltaMap func([]byte) (interface{}, error)
	unmarshalWatchData func([]byte) (interface{}, error)

	marshalTreeDeltaMap func(interface{}) ([]byte, error)
	marshalWatchData func(interface{}) ([]byte, error)
)

type DeltaInfo struct {
	abciResponses *ABCIResponses
	treeDeltaMap interface{}
	watchData interface{}
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

	pl.WatchBytes, err = marshalWatchData(info.watchData)
	if err != nil {
		return pl, err
	}

	return pl, err
}


func (dd *DeltaInfo) bytes2DeltaInfo(pl *types.DeltaPayload) (error) {
	var err error
	err = types.Json.Unmarshal(pl.ABCIRsp, &dd.abciResponses)
	if err != nil {
		return err
	}

	dd.treeDeltaMap, err = unmarshalTreeDeltaMap(pl.DeltasBytes)
	if err != nil {
		return err
	}
	dd.watchData, err = unmarshalWatchData(pl.WatchBytes)
	if err != nil {
		return err
	}

	return err
}


func bytes2DeltaInfo(pl *types.DeltaPayload) (*DeltaInfo, error) {
	var err error
	dd := &DeltaInfo{}
	err = types.Json.Unmarshal(pl.ABCIRsp, &dd.abciResponses)
	if err != nil {
		return nil, err
	}

	dd.treeDeltaMap, err = unmarshalTreeDeltaMap(pl.DeltasBytes)
	if err != nil {
		return nil, err
	}
	dd.watchData, err = unmarshalWatchData(pl.WatchBytes)
	if err != nil {
		return nil, err
	}

	return dd, err
}
