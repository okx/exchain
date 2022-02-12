package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
	"strings"
)


var (
	// DefaultIBCVersion represents the latest supported version of IBC used
	// in connection version negotiation. The current version supports only
	// ORDERED and UNORDERED channels and requires at least one channel type
	// to be agreed upon.
	DefaultIBCVersion = NewVersion(DefaultIBCVersionIdentifier, []string{"ORDER_ORDERED", "ORDER_UNORDERED"})

	// DefaultIBCVersionIdentifier is the IBC v1.0.0 protocol version identifier
	DefaultIBCVersionIdentifier = "1"

	// AllowNilFeatureSet is a helper map to indicate if a specified version
	// identifier is allowed to have a nil feature set. Any versions supported,
	// but not included in the map default to not supporting nil feature sets.
	allowNilFeatureSet = map[string]bool{
		DefaultIBCVersionIdentifier: false,
	}
)

// ProtoVersionsToExported converts a slice of the Version proto definition to
// the Version interface.
func ProtoVersionsToExported(versions []*Version) []exported.Version {
	exportedVersions := make([]exported.Version, len(versions))
	for i := range versions {
		exportedVersions[i] = versions[i]
	}

	return exportedVersions
}



var _ exported.Version = &Version{}

// NewVersion returns a new instance of Version.
func NewVersion(identifier string, features []string) *Version {
	return &Version{
		Identifier: identifier,
		Features:   features,
	}
}

// GetIdentifier implements the VersionI interface
func (version Version) GetIdentifier() string {
	return version.Identifier
}

// GetFeatures implements the VersionI interface
func (version Version) GetFeatures() []string {
	return version.Features
}

// ValidateVersion does basic validation of the version identifier and
// features. It unmarshals the version string into a Version object.
func ValidateVersion(version *Version) error {
	if version == nil {
		return sdkerrors.Wrap(ErrInvalidVersion, "version cannot be nil")
	}
	if strings.TrimSpace(version.Identifier) == "" {
		return sdkerrors.Wrap(ErrInvalidVersion, "version identifier cannot be blank")
	}
	for i, feature := range version.Features {
		if strings.TrimSpace(feature) == "" {
			return sdkerrors.Wrapf(ErrInvalidVersion, "feature cannot be blank, index %d", i)
		}
	}

	return nil
}

// VerifyProposedVersion verifies that the entire feature set in the
// proposed version is supported by this chain. If the feature set is
// empty it verifies that this is allowed for the specified version
// identifier.
func (version Version) VerifyProposedVersion(proposedVersion exported.Version) error {
	if proposedVersion.GetIdentifier() != version.GetIdentifier() {
		return sdkerrors.Wrapf(
			ErrVersionNegotiationFailed,
			"proposed version identifier does not equal supported version identifier (%s != %s)", proposedVersion.GetIdentifier(), version.GetIdentifier(),
		)
	}

	if len(proposedVersion.GetFeatures()) == 0 && !allowNilFeatureSet[proposedVersion.GetIdentifier()] {
		return sdkerrors.Wrapf(
			ErrVersionNegotiationFailed,
			"nil feature sets are not supported for version identifier (%s)", proposedVersion.GetIdentifier(),
		)
	}

	for _, proposedFeature := range proposedVersion.GetFeatures() {
		if !contains(proposedFeature, version.GetFeatures()) {
			return sdkerrors.Wrapf(
				ErrVersionNegotiationFailed,
				"proposed feature (%s) is not a supported feature set (%s)", proposedFeature, version.GetFeatures(),
			)
		}
	}

	return nil
}
// contains returns true if the provided string element exists within the
// string set.
func contains(elem string, set []string) bool {
	for _, element := range set {
		if elem == element {
			return true
		}
	}

	return false
}
