package types

import (
	"encoding/json"
	"errors"
	"fmt"

	stakingtypes "github.com/okex/exchain/x/staking/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GenesisState defines the raw genesis transaction in JSON
type GenesisState struct {
	GenTxs []json.RawMessage `json:"gentxs" yaml:"gentxs"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(genTxs []json.RawMessage) GenesisState {
	return GenesisState{
		GenTxs: genTxs,
	}
}

// ValidateGenesis validates GenTx transactions
func ValidateGenesis(genesisState GenesisState) error {
	for i, genTx := range genesisState.GenTxs {
		var tx authtypes.StdTx
		if err := ModuleCdc.UnmarshalJSON(genTx, &tx); err != nil {
			return err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		// TODO: abstract back to staking
		if _, ok := msgs[0].(stakingtypes.MsgCreateValidator); !ok {
			return fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}
	return nil
}
