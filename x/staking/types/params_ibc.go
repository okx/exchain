package types

import "github.com/okex/exchain/x/params"

var (
	_ params.ParamSet = KeyHistoricalEntriesParamsSet{}
)

type KeyHistoricalEntriesParamsSet struct {
	HistoricalEntries uint32 `protobuf:"varint,4,opt,name=historical_entries,json=historicalEntries,proto3" json:"historical_entries,omitempty" yaml:"historical_entries"`
}

func (p KeyHistoricalEntriesParamsSet) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyHistoricalEntries, Value: &p.HistoricalEntries, ValidatorFn: validateHistoricalEntries},
	}
}

func KeyHistoricalEntriesParams(p uint32) params.ParamSet {
	return &KeyHistoricalEntriesParamsSet{HistoricalEntries: p}
}
