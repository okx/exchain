package types

type SimulateData struct {
	TxBytes        []byte `json:"tx"`
	OverridesBytes []byte `json:"overrides"`
}
