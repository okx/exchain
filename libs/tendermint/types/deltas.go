package types

// state-delta mode
// 0 same as no state-delta
// 1 product delta and save into deltastore.db
// 2 consume delta and save into deltastore.db; if get no delta, do as 1
const (
	// for getting flag of delta-mode
	FlagStateDelta = "state-sync-mode"

	// delta-mode
	NoDelta      = "na"
	ProductDelta = "producer"
	ConsumeDelta = "consumer"

	// data-center
	FlagDataCenter = "data-center-mode"
	DataCenterUrl  = "data-center-url"
	DataCenterStr  = "dataCenter"

	// fast-query
	FlagFastQuery = "fast-query"
)

type BlockDelta struct {
	Block  *Block  `json:"block"`
	Deltas *Deltas `json:"deltas"`
	Height int64   `json:"height"`
}

// Deltas defines the ABCIResponse and state delta
type Deltas struct {
	ABCIRsp     []byte `json:"abci_rsp"`
	DeltasBytes []byte `json:"deltas_bytes"`
	Height      int64  `json:"height"`
}

// Size returns size of the deltas in bytes.
func (d *Deltas) Size() int {
	if d == nil {
		return 0
	}
	return len(d.ABCIRsp) + len(d.DeltasBytes)
}

// Marshal returns the amino encoding.
func (d *Deltas) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(d)
}

// Unmarshal deserializes from amino encoded form.
func (d *Deltas) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, d)
}

// WatchData defines the batch in watchDB and accounts for delete
type WatchData struct {
	//DirtyAccount []byte `json:"dirty_account"`
	//Batches      []byte `json:"batches"`
	WatchDataByte []byte `json:"watch_data_byte"`
	Height        int64  `json:"height"`
}

// Size returns size of the deltas in bytes.
func (wd *WatchData) Size() int {
	if wd == nil {
		return 0
	}
	//	return len(wd.DirtyAccount) + len(wd.Batches)
	return len(wd.WatchDataByte)
}

// Marshal returns the amino encoding.
func (wd *WatchData) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(wd)
}

// Unmarshal deserializes from amino encoded form.
func (wd *WatchData) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, wd)
}
