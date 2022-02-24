package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"sync"

	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tendermintTypes "github.com/okex/exchain/libs/tendermint/types"
)

var (
	regexChainID     = `[a-z]*`
	regexSeparator   = `-{1}`
	regexEpoch       = `[1-9][0-9]*`
	ethermintChainID = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, regexChainID, regexSeparator, regexEpoch))
)

const mainnetChainId = "exchain-66"
const testnetChainId = "exchain-65"

var (
	chainIdSetOnce sync.Once
	chainId        string
	chainIdEpoch   *big.Int
)

// IsValidChainID returns false if the given chain identifier is incorrectly formatted.
func IsValidChainID(chainID string) bool {
	if len(chainID) > 48 {
		return false
	}

	return ethermintChainID.MatchString(chainID)
}
func isMainNetChainID(chainID string) bool {
	return chainID == mainnetChainId
}
func isTestNetChainID(chainID string) bool {
	return chainID == testnetChainId
}

func SetChainId(chainid string) error {
	epoch, err := ParseChainID(chainid)
	if err != nil {
		return err
	}
	chainIdSetOnce.Do(func() {
		chainId = chainid
		chainIdEpoch = epoch
	})
	return nil
}
func GetChainId() string {
	if chainIdEpoch == nil {
		panic(sdkerrors.Wrap(ErrInvalidChainID, "ChainId is not Init, need to call SetChainId when the app start"))
	}
	return chainId
}
func GetChainIdEpoch() *big.Int {
	if chainIdEpoch == nil {
		panic(sdkerrors.Wrap(ErrInvalidChainID, "ChainId is not Init, need to call SetChainId when the app start"))
	}
	return chainIdEpoch
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
	if isMainNetChainID(chainID) && !tendermintTypes.IsMainNet() {
		return fmt.Errorf("Must use <make mainnet> to rebuild if chain-id is <%s>, Current GenesisHeight is <%d>", chainID, tendermintTypes.GetStartBlockHeight())
	}
	if isTestNetChainID(chainID) && !tendermintTypes.IsTestNet() {
		return fmt.Errorf("Must use <make testnet> to rebuild if chain-id is <%s>, Current GenesisHeight is <%d>", chainID, tendermintTypes.GetStartBlockHeight())
	}
	return nil
}
