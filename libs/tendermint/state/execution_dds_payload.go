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
	treeDeltaMap := iavl.TreeDeltaMap{}
	err := treeDeltaMap.UnmarshalFromAmino(nil, input)
	return treeDeltaMap, err
}

func marshalTreeDeltaMap(deltaMap interface{}) ([]byte, error) {
	dm, ok := deltaMap.(iavl.TreeDeltaMap)
	if !ok {
		return nil, fmt.Errorf("failed marshal TreeDeltaMap")
	}
	return dm.MarshalToAmino(nil)
}

type DeltaInfo struct {
	abciResponses *ABCIResponses
	treeDeltaMap  interface{}
	watchData     interface{}

	marshalWatchData func() ([]byte, error)
}

// for upload
func (info *DeltaInfo) dataInfo2Bytes() (types.DeltaPayload, error) {
	var err error
	pl := types.DeltaPayload{}
	pl.ABCIRsp, err = info.abciResponses.MarshalToAmino(cdc)
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
	if pl == nil {
		return fmt.Errorf("Failed bytes to delta info: empty delta payload. ")
	}
	var err error
	ar := &ABCIResponses{}
	err = ar.UnmarshalFromAmino(nil, pl.ABCIRsp)
	if err != nil {
		return err
	}
	info.abciResponses = ar

	info.treeDeltaMap, err = unmarshalTreeDeltaMap(pl.DeltasBytes)
	if err != nil {
		return err
	}
	info.watchData, err = unmarshalWatchData(pl.WatchBytes)
	if err != nil {
		return err
	}

	return err
}
