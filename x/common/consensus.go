package common

import (
	"fmt"

	"github.com/okx/okbchain/x/params/subspace"
)

type ConsensusType string

const (
	PoA ConsensusType = "poa"
	DPoS ConsensusType = "dpos"
)

func ValidateConsensusType(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(ConsensusType)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v != PoA && v != DPoS {
			return fmt.Errorf("%s must be %s or %s: %s", param, PoA, DPoS, v)
		}

		return nil
	}
}
