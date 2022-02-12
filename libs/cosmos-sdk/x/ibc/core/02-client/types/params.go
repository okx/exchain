package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"strings"
)


var (
	// DefaultAllowedClients are "06-solomachine" and "07-tendermint"
	//DefaultAllowedClients = []string{exported.Solomachine, exported.Tendermint}

	// KeyAllowedClients is store's key for AllowedClients Params
	KeyAllowedClients = []byte("AllowedClients")
)


// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}


// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyAllowedClients, p.AllowedClients, validateClients),
	}
}

func validateClients(i interface{}) error {
	clients, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for i, clientType := range clients {
		if strings.TrimSpace(clientType) == "" {
			return fmt.Errorf("client type %d cannot be blank", i)
		}
	}

	return nil
}
