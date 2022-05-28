package types

import (
	"bytes"
	"errors"
	"fmt"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
)

func (sh SignedHeader) ValidateBasicForIBC(chainID string) error {
	if sh.Header == nil {
		return errors.New("missing header")
	}
	if sh.Commit == nil {
		return errors.New("missing commit")
	}

	if err := sh.Header.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid header: %w", err)
	}
	if err := sh.Commit.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid commit: %w", err)
	}

	if sh.ChainID != chainID {
		return fmt.Errorf("header belongs to another chain %q, not %q", sh.ChainID, chainID)
	}

	// Make sure the header is consistent with the commit.
	if sh.Commit.Height != sh.Height {
		return fmt.Errorf("header and commit height mismatch: %d vs %d", sh.Height, sh.Commit.Height)
	}
	if hhash, chash := sh.PureIBCHash(), sh.Commit.BlockID.Hash; !bytes.Equal(hhash, chash) {
		return fmt.Errorf("commit signs block %X, header is block %X", chash, hhash)
	}
	return nil
}

func (h *Header) PureIBCHash() tmbytes.HexBytes {
	if h == nil || len(h.ValidatorsHash) == 0 {
		return nil
	}
	return h.IBCHash()
}

func (commit *Commit) IBCVoteSignBytes(chainID string, valIdx int) []byte {
	return commit.GetVote(valIdx).ibcSignBytes(chainID)
}
