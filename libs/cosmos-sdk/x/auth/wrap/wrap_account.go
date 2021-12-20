package wrap

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/pkg/errors"
)

type WrapAccount struct {
	RealAcc exported.Account
}

func (acc *WrapAccount) DecodeRLP(s *rlp.Stream) error {
	var kind uint
	err := s.Decode(&kind)
	if err != nil {
		return err
	}

	if varFunc, ok  := exported.ConcreteAccount[kind]; ok {
		data, err := s.Raw()
		if err != nil {
			return err
		}

		puppet := varFunc()
		err = puppet.RLPDecodeBytes(data)
		acc.RealAcc = puppet
		return err
	}

	return errors.New("Unknown Account type")
}