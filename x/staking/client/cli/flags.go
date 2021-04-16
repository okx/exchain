package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/okex/exchain/x/staking/types"
)

// nolint
const (
	FlagPubKey = "pubkey"

	FlagMoniker  = "moniker"
	FlagIdentity = "identity"
	FlagWebsite  = "website"
	FlagDetails  = "details"

	//FlagCommissionRate          = "commission-rate"
	//FlagCommissionMaxRate       = "commission-max-rate"
	//FlagCommissionMaxChangeRate = "commission-max-change-rate"

	//FlagMinSelfDelegation = "min-self-delegation"

	FlagNodeID = "node-id"
	FlagIP     = "ip"
)

// common flagsets to add to various functions
var (
	FsPk                = flag.NewFlagSet("", flag.ContinueOnError)
	fsDescriptionCreate = flag.NewFlagSet("", flag.ContinueOnError)
	//FsCommissionCreate  = flag.NewFlagSet("", flag.ContinueOnError)
	//fsCommissionUpdate  = flag.NewFlagSet("", flag.ContinueOnError)
	//FsMinSelfDelegation = flag.NewFlagSet("", flag.ContinueOnError)
	fsDescriptionEdit = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded PubKey of the validator")
	fsDescriptionCreate.String(FlagMoniker, "", "The validator's name")
	fsDescriptionCreate.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fsDescriptionCreate.String(FlagWebsite, "", "The validator's (optional) website")
	fsDescriptionCreate.String(FlagDetails, "", "The validator's (optional) details")
	//fsCommissionUpdate.String(FlagCommissionRate, "", "The new commission rate percentage")
	//FsCommissionCreate.String(FlagCommissionRate, "", "The initial commission rate percentage")
	//FsCommissionCreate.String(FlagCommissionMaxRate, "", "The maximum commission rate percentage")
	//FsCommissionCreate.String(FlagCommissionMaxChangeRate, "", "The maximum commission change rate percentage (per day)")
	//FsMinSelfDelegation.String(FlagMinSelfDelegation, fmt.Sprintf("0.001%s", sdk.DefaultBondDenom),
	//	"The minimum self delegation required on the validator")
	fsDescriptionEdit.String(FlagMoniker, types.DoNotModifyDesc, "The validator's name")
	fsDescriptionEdit.String(FlagIdentity, types.DoNotModifyDesc,
		"The (optional) identity signature (ex. UPort or Keybase)")
	fsDescriptionEdit.String(FlagWebsite, types.DoNotModifyDesc, "The validator's (optional) website")
	fsDescriptionEdit.String(FlagDetails, types.DoNotModifyDesc, "The validator's (optional) details")
}
