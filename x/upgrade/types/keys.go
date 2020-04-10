package types

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

const (
	// ModuleName is the name of the upgrade module
	ModuleName = "upgrade"
	// StoreKey is the string store representation
	StoreKey = ModuleName
	// RouterKey is the msg router key for the upgrade module
	RouterKey = ModuleName
	// QuerierRoute is the querier route for the upgrade module
	QuerierRoute = ModuleName
	// DefaultParamspace is the default paramspace for the upgrade module
	DefaultParamspace = ModuleName
)

// keys
var (
	signalKey         = "s/%s/%s"      // s/<protocolVersion>/<switchVoterAddress>
	proposalIDKey     = "p/%s"         // p/<proposalId>
	successVersionKey = "success/%s"   // success/<protocolVersion>
	failedVersionKey  = "failed/%s/%s" // failed/<protocolVersion>/<proposalId>
	signalPrefixKey   = "s/%s"
)

// GetSignalKey gets signal store key
func GetSignalKey(versionID uint64, switchVoterAddr string) []byte {
	return []byte(fmt.Sprintf(signalKey, uintToHexString(versionID), switchVoterAddr))
}

func uintToHexString(i uint64) string {
	hex := strconv.FormatUint(i, 16)
	var stringBuild bytes.Buffer
	for i := 0; i < 16-len(hex); i++ {
		if _, err := stringBuild.Write([]byte("0")); err != nil {
			log.Println(err)
		}

	}
	if _, err := stringBuild.Write([]byte(hex)); err != nil {
		log.Println(err)
	}

	return stringBuild.String()
}

// GetProposalIDKey gets proposal ID store key
func GetProposalIDKey(proposalID uint64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, uintToHexString(proposalID)))
}

// GetSuccessVersionKey gets successful version store key
func GetSuccessVersionKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(successVersionKey, uintToHexString(versionID)))
}

// GetFailedVersionKey gets failed version store key
func GetFailedVersionKey(versionID uint64, proposalID uint64) []byte {
	return []byte(fmt.Sprintf(failedVersionKey, uintToHexString(versionID), uintToHexString(proposalID)))
}

// GetSignalPrefixKey gets signal prefix store key
func GetSignalPrefixKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(signalPrefixKey, uintToHexString(versionID)))
}
