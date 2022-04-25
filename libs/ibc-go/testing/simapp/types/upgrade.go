package types

import "fmt"

const (
	KeyUpgradedIBCState = "upgradedIBCState"

	// KeyUpgradedClient is the sub-key under which upgraded client state will be stored
	KeyUpgradedClient = "upgradedClient"

	// KeyUpgradedConsState is the sub-key under which upgraded consensus state will be stored
	KeyUpgradedConsState = "upgradedConsState"
)

func UpgradedClientKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedClient))
}

func UpgradedConsStateKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedConsState))
}
