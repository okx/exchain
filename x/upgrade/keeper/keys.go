package keeper

import (
	"bytes"
	"fmt"
	"strconv"
)

// keys
var (
	signalKey         = "s/%s/%s"      // s/<protocolVersion>/<switchVoterAddress>
	proposalIDKey     = "p/%s"         // p/<proposalId>
	successVersionKey = "success/%s"   // success/<protocolVersion>
	failedVersionKey  = "failed/%s/%s" // failed/<protocolVersion>/<proposalId>
	signalPrefixKey   = "s/%s"
)

func GetSignalKey(versionID uint64, switchVoterAddr string) []byte {
	return []byte(fmt.Sprintf(signalKey, UintToHexString(versionID), switchVoterAddr))
}

func UintToHexString(i uint64) string {
	hex := strconv.FormatUint(i, 16)
	var stringBuild bytes.Buffer
	for i := 0; i < 16-len(hex); i++ {
		stringBuild.Write([]byte("0"))
	}
	stringBuild.Write([]byte(hex))
	return stringBuild.String()
}

func GetProposalIDKey(proposalID uint64) []byte {
	return []byte(fmt.Sprintf(proposalIDKey, UintToHexString(proposalID)))
}

func GetSuccessVersionKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(successVersionKey, UintToHexString(versionID)))
}

func GetFailedVersionKey(versionID uint64, proposalID uint64) []byte {
	return []byte(fmt.Sprintf(failedVersionKey, UintToHexString(versionID), UintToHexString(proposalID)))
}

func GetSignalPrefixKey(versionID uint64) []byte {
	return []byte(fmt.Sprintf(signalPrefixKey, UintToHexString(versionID)))
}
