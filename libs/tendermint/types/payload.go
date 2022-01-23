package types

type DeltaPayload struct {
	ABCIRsp     []byte
	DeltasBytes []byte
	WatchBytes  []byte
}

