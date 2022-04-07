package types

// TokenMapping defines a mapping between native denom and contract
type TokenMapping struct {
	Denom    string `json:"denom"`
	Contract string `json:"contract"`
}

// GenesisState defines the erc20 module genesis state
type GenesisState struct {
	Params            Params         `json:"params"`
	ExternalContracts []TokenMapping `json:"external_contracts"`
	AutoContracts     []TokenMapping `json:"auto_contracts"`
}

// DefaultGenesisState sets default erc20 genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
