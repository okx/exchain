package consensus


import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	cosmossdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/p2p/conn"
	tmtime "github.com/okex/exchain/libs/tendermint/types/time"
	"testing"

	"github.com/okex/exchain/libs/tendermint/types"
)

func TestP2PMsgProcess(t *testing.T) {
	// init (privateKey, publicKey)
	privKey, _ := ethsecp256k1.GenerateKey()
	addr := cosmossdk.AccAddress(privKey.PubKey().Address())

	hash := []byte("0xff")
	header := types.PartSetHeader{
		Total: 1,
		Hash:  hash,
	}

	// construct vote message
	vote := &types.Vote{
		ValidatorAddress: types.Address(addr),
		ValidatorIndex:   1,
		Height:           int64(1),
		Round:            0,
		Timestamp:        tmtime.Now(),
		Type:             types.PrevoteType,
		BlockID:          types.BlockID{Hash: hash, PartsHeader: header},
	}
	signBytes := vote.SignBytes("65")
	sig, err := privKey.Sign(signBytes)
	if err != nil {
		fmt.Println(err)
	}
	vote.Signature = sig

	//
	msginfo := msgInfo{&VoteMessage{vote}, ""}
	msgBytes := cdc.MustMarshalBinaryBare(msginfo)

	buf := bytes.Buffer{}
	for i := 0; i < 10; i++ {
		buf.Write(msgBytes)
	}
	msg := conn.PacketMsg{0x7f, 45, buf.Bytes()}
	bz, err := cdc.MarshalBinaryLengthPrefixed(msg)

	//
	var packet conn.Packet
	var bufTmp = bytes.NewBuffer(bz)

	packet, _, err = conn.UnmarshalPacketFromAminoReader(bufTmp, int64(bufTmp.Len()))
	if err != nil {
		t.Fatal(err)
	}
	_ = packet

}
