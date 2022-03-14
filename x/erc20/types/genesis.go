package types

// GenesisState defines the erc20 module genesis state
type GenesisState struct {
	Params Params `json:"params"`
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
