package mock

import (
	"github.com/okex/exchain/libs/tendermint/crypto"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

var _ tmtypes.PrivValidator = PV{}

// MockPV implements PrivValidator without any safety or persistence.
// Only use it for testing.
type PV struct {
	PrivKey ed25519.PrivKeyEd25519
}

func NewPV() PV {
	return PV{ed25519.GenPrivKey()}
}

// GetPubKey implements PrivValidator interface
func (pv PV) GetPubKey() (crypto.PubKey, error) {
	//return cryptocodec.ToTmPubKeyInterface(pv.PrivKey.PubKey())
	return pv.PrivKey.PubKey(), nil
}

// SignVote implements PrivValidator interface
func (pv PV) SignVote(chainID string, vote *tmtypes.Vote) error {
	signBytes := tmtypes.VoteSignBytes(chainID, vote)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	return nil
}

// SignProposal implements PrivValidator interface
func (pv PV) SignProposal(chainID string, proposal *tmtypes.Proposal) error {
	signBytes := tmtypes.ProposalSignBytes(chainID, proposal)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	proposal.Signature = sig
	return nil
}
