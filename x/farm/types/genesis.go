package types

// GenesisState - all farm state that must be provided at genesis
type GenesisState struct {
	// TODO: Fill out what is needed by the module for genesis
	Params Params `json:"params" yaml:"params"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params /* TODO: Fill out with what is needed for genesis state */) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		Params: params,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state, these values will be initialized but empty
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the farm genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	return nil
}
