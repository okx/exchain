package types

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type ValStatus byte

const (
	Created    ValStatus = 0x00
	Destroying ValStatus = 0x01
	Destroyed  ValStatus = 0x02

	ValStatusCreated    = "Created"
	ValStatusDestroying = "Destroying"
	ValStatusDestroyed  = "Destroied"
)

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	Address             sdk.ConsAddress `json:"address" yaml:"address"`                             // validator consensus address
	StartHeight         int64           `json:"start_height" yaml:"start_height"`                   // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64           `json:"index_offset" yaml:"index_offset"`                   // index offset into signed block bit array
	JailedUntil         time.Time       `json:"jailed_until" yaml:"jailed_until"`                   // timestamp validator cannot be unjailed until
	Tombstoned          bool            `json:"tombstoned" yaml:"tombstoned"`                       // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64           `json:"missed_blocks_counter" yaml:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
	ValidatorStatus     ValStatus       `json:"validator_status" yaml:"validator_status"`
}

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
func NewValidatorSigningInfo(
	condAddr sdk.ConsAddress, startHeight, indexOffset int64,
	jailedUntil time.Time, tombstoned bool, missedBlocksCounter int64, validatorStatus ValStatus,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:             condAddr,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		ValidatorStatus:     validatorStatus,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}
