package types

import "fmt"

const (
	// SubModuleName defines the IBC client name
	SubModuleName string = "client"

	// RouterKey is the message route for IBC client
	RouterKey string = SubModuleName

	// QuerierRoute is the querier route for IBC client
	QuerierRoute string = SubModuleName

	// KeyNextClientSequence is the key used to store the next client sequence in
	// the keeper.
	KeyNextClientSequence = "nextClientSequence"
)

// FormatClientIdentifier returns the client identifier with the sequence appended.
// This is a SDK specific format not enforced by IBC protocol.
func FormatClientIdentifier(clientType string, sequence uint64) string {
	return fmt.Sprintf("%s-%d", clientType, sequence)
}

