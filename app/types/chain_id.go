package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tendermintTypes "github.com/okex/exchain/libs/tendermint/types"
)

var (
	regexChainID     = `[a-z]*`
	regexSeparator   = `-{1}`
	regexEpoch       = `[1-9][0-9]*`
	ethermintChainID = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, regexChainID, regexSeparator, regexEpoch))
)

const mainnet_chain_id = "exchain-66"
const testnet_chain_id = "exchain-65"

// IsValidChainID returns false if the given chain identifier is incorrectly formatted.
func IsValidChainID(chainID string) bool {
	if len(chainID) > 48 {
		return false
	}

	return ethermintChainID.MatchString(chainID)
}
func IsMainNetChainID(chainID string) bool {
	return chainID == mainnet_chain_id
}
func IsTestNetChainID(chainID string) bool {
	return chainID == testnet_chain_id
}

// ParseChainID parses a string chain identifier's epoch to an Ethereum-compatible
// chain-id in *big.Int format. The function returns an error if the chain-id has an invalid format
func ParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)
	if len(chainID) > 48 {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "chain-id '%s' cannot exceed 48 chars", chainID)
	}

	matches := ethermintChainID.FindStringSubmatch(chainID)
	if matches == nil || len(matches) != 3 || matches[1] == "" {
		return nil, sdkerrors.Wrap(ErrInvalidChainID, chainID)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(matches[2], 10)
	if !ok {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "epoch %s must be base-10 integer format", matches[2])
	}

	return chainIDInt, nil
}

func IsValidateChainIdWithGenesisHeight(chainID string) error {
	if tendermintTypes.IsMainNet() && !IsMainNetChainID(chainID) {
		return fmt.Errorf("Must use <make mainnet> to rebuild if chain-id is <%s>, Current GenesisHeight is <%d>", chainID, tendermintTypes.GetStartBlockHeight())
	} else if tendermintTypes.IsTestNet() && !IsMainNetChainID(chainID) {
		return fmt.Errorf("Must use <make testnet> to rebuild if chain-id is <%s>, Current GenesisHeight is <%d>", chainID, tendermintTypes.GetStartBlockHeight())
	} else {
		return nil
	}
}
