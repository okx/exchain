package privval

import (
	"fmt"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/types"
)

func DefaultValidationRequestHandler(
	privVal types.PrivValidator,
	req SignerMessage,
	chainID string,
) (SignerMessage, error) {
	var res SignerMessage
	var err error

	switch r := req.(type) {
	case *PubKeyRequest:
		var pubKey crypto.PubKey
		pubKey, err = privVal.GetPubKey()
		if err != nil {
			res = &PubKeyResponse{nil, &RemoteSignerError{0, err.Error()}}
		} else {
			res = &PubKeyResponse{pubKey, nil}
		}

	case *SignVoteRequest:
		err = privVal.SignVote(chainID, r.Vote)
		if err != nil {
			res = &SignedVoteResponse{nil, &RemoteSignerError{0, err.Error()}}
		} else {
			res = &SignedVoteResponse{r.Vote, nil}
		}

	case *SignProposalRequest:
		err = privVal.SignProposal(chainID, r.Proposal)
		if err != nil {
			res = &SignedProposalResponse{nil, &RemoteSignerError{0, err.Error()}}
		} else {
			res = &SignedProposalResponse{r.Proposal, nil}
		}

	case *PingRequest:
		err, res = nil, &PingResponse{}

	default:
		err = fmt.Errorf("unknown msg: %v", r)
	}

	return res, err
}
