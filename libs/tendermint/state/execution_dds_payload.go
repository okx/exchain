package state

import (
	"fmt"

	"github.com/okx/okbchain/libs/tendermint/types"
)

func unmarshalTreeDeltaMap(input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("failed unmarshal TreeDeltaMap: empty data")
	}
	treeDeltaMap := types.NewTreeDelta()
	err := treeDeltaMap.Unmarshal(input)
	return treeDeltaMap, err
}

func marshalTreeDeltaMap(deltaMap interface{}) ([]byte, error) {
	dm, ok := deltaMap.(*types.TreeDelta)
	if !ok {
		return nil, fmt.Errorf("failed marshal TreeDeltaMap")
	}
	return dm.Marshal(), nil
}

type DeltaInfo struct {
	from          string
	deltaLen      int
	deltaHeight   int64
	abciResponses *ABCIResponses
	treeDeltaMap  interface{}
	watchData     interface{}
	wasmWatchData interface{}

	marshalWatchData     func() ([]byte, error)
	wasmMarshalWatchData func() ([]byte, error)
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

	pl.WasmWatchBytes, err = info.wasmMarshalWatchData()
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
	if types.FastQuery {
		info.watchData, err = evmWatchDataManager.UnmarshalWatchData(pl.WatchBytes)
		if err != nil {
			return err
		}
	}
	info.wasmWatchData, err = wasmWatchDataManager.UnmarshalWatchData(pl.WasmWatchBytes)
	if err != nil {
		return err
	}

	return err
}
