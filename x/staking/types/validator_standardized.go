package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Export returns the exported format of validator in genesis export
func (v Validator) Export() ValidatorExported {
	consPkStr, err := Bech32ifyConsPub(v.ConsPubKey)
	if err != nil {
		panic(fmt.Sprintf("bech32 of consensus pubkey error: %s", err.Error()))
	}

	return ValidatorExported{
		v.OperatorAddress,
		consPkStr,
		v.Jailed,
		v.Status,
		v.DelegatorShares,
		v.Description,
		v.UnbondingHeight,
		v.UnbondingCompletionTime,
		v.MinSelfDelegation,
	}
}

// Standardize converts inner struct Validator to StandardizedValidator which is used to display
func (v Validator) Standardize() StandardizedValidator {
	return StandardizedValidator{
		v.OperatorAddress,
		v.ConsPubKey,
		v.Jailed,
		v.Status,
		v.DelegatorShares,
		v.Description,
		v.UnbondingHeight,
		v.UnbondingCompletionTime,
		v.MinSelfDelegation,
	}
}

// Export returns the exported format of Validators in genesis export
func (v Validators) Export() []ValidatorExported {
	valsLen := len(v)
	valExported := make([]ValidatorExported, valsLen)
	for i := 0; i < valsLen; i++ {
		valExported[i] = v[i].Export()
	}

	return valExported
}

// Standardize converts inner struct Validators to StandardizedValidators which is used to display
func (v Validators) Standardize() StandardizedValidators {
	n := len(v)
	standardizedVals := make(StandardizedValidators, n)
	for i := 0; i < n; i++ {
		standardizedVals[i] = v[i].Standardize()
	}
	return standardizedVals
}

// StandardizedValidator is just a copy of Validator in cosmos sdk
// The field "DelegatorShares"/"MinSelfDelegation" is treated by descending power 8 to decimal
type StandardizedValidator struct {
	OperatorAddress         sdk.ValAddress `json:"operator_address" yaml:"operator_address"`
	ConsPubKey              crypto.PubKey  `json:"consensus_pubkey" yaml:"consensus_pubkey"`
	Jailed                  bool           `json:"jailed" yaml:"jailed"`
	Status                  sdk.BondStatus `json:"status" yaml:"status"`
	DelegatorShares         sdk.Dec        `json:"delegator_shares" yaml:"delegator_shares"`
	Description             Description    `json:"description" yaml:"description"`
	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`
	UnbondingCompletionTime time.Time      `json:"unbonding_time" yaml:"unbonding_time"`
	MinSelfDelegation       sdk.Dec        `json:"min_self_delegation" yaml:"min_self_delegation"`
}

// String returns a human readable string representation of a StandardizeValidator
func (sv StandardizedValidator) String() string {
	bechConsPubkey, err := Bech32ifyConsPub(sv.ConsPubKey)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`Validator
  Operator Address:           %s
  Validator Consensus Pubkey: %s
  Jailed:                     %v
  Status:                     %s
  Delegator Shares:           %s
  Description:                %s
  Unbonding Height:           %d
  Unbonding Completion Time:  %v
  Minimum Self Delegation:    %v`,
		sv.OperatorAddress, bechConsPubkey, sv.Jailed, sv.Status,
		sv.DelegatorShares, sv.Description, sv.UnbondingHeight,
		sv.UnbondingCompletionTime, sv.MinSelfDelegation)
}

// MarshalYAML implememts the text format for yaml marshaling
func (sv StandardizedValidator) MarshalYAML() (interface{}, error) {
	return sv.String(), nil
}

// StandardizedValidators is the type alias of the StandardizedValidator slice
type StandardizedValidators []StandardizedValidator

// MarshalYAML implememts the text format for yaml marshaling
func (svs StandardizedValidators) MarshalYAML() (interface{}, error) {
	return svs.String(), nil
}

// String returns a human readable string representation of StandardizeValidators
func (svs StandardizedValidators) String() string {
	var output string
	n := len(svs)
	for i := 0; i < n; i++ {
		output += svs[i].String() + "\n"
	}
	return strings.TrimSpace(output)
}

// shares from the shares adding convert to the consensus power
func sharesToConsensusPower(shares Shares) int64 {
	return shares.QuoInt(sdk.PowerReduction).Int64()
}

// PotentialConsensusPowerByShares gets potential consensus-engine power based on shares
func (v Validator) PotentialConsensusPowerByShares() int64 {
	return sharesToConsensusPower(v.DelegatorShares)
}

// ConsensusPowerByShares gets the consensus-engine power
func (v Validator) ConsensusPowerByShares() int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPowerByShares()
	}
	return 0
}

// ABCIValidatorUpdateByShares returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power based on shares
func (v Validator) ABCIValidatorUpdateByShares() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  v.ConsensusPowerByShares(),
	}
}
