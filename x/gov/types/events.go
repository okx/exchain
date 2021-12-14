package types

const (
	// AttributeKeyProposalStatus defines the proposal status attribute in gov
	AttributeKeyProposalStatus = "proposal_status"

	EventTypeSubmitProposal    = "submit_proposal"
	EventTypeProposalDeposit   = "proposal_deposit"
	EventTypeProposalVote      = "proposal_vote"
	EventTypeProposalVoteTally = "proposal_vote_tally"
	EventTypeInactiveProposal  = "inactive_proposal"
	EventTypeActiveProposal    = "active_proposal"

	AttributeKeyProposalResult     = "proposal_result"
	AttributeKeyProposalLog        = "proposal_result_log"
	AttributeKeyOption             = "option"
	AttributeKeyProposalID         = "proposal_id"
	AttributeKeyVotingPeriodStart  = "voting_period_start"
	AttributeValueCategory         = "governance"
	AttributeValueProposalDropped  = "proposal_dropped"  // didn't meet min deposit
	AttributeValueProposalPassed   = "proposal_passed"   // met vote quorum
	AttributeValueProposalRejected = "proposal_rejected" // didn't meet vote quorum
	AttributeValueProposalFailed   = "proposal_failed"   // error on proposal handler
)
