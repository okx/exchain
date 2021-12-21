package types

import "fmt"

// GenesisState - all nameservice state that must be provided at genesis
type GenesisState struct {
	// TODO: Fill out what is needed by the module for genesis
	WhoisRecords []Whois `json:"whois_records"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState( /* TODO: Fill out with what is needed for genesis state */ ) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		WhoisRecords: nil,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []Whois{},
	}
}

// ValidateGenesis validates the nameservice genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	for _, record := range data.WhoisRecords {
		if record.Creator == nil {
			return fmt.Errorf("invalid WhoisRecord: Creator: %s. Error: Missing Creator", record.Creator)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Value", record.Value)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecord: Price: %s. Error: Missing Price", record.Price)
		}
	}
	return nil
}
