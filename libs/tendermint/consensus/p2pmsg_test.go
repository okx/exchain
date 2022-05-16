package consensus

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	cosmossdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/autofile"
	"github.com/okex/exchain/libs/tendermint/p2p/conn"
	tmtime "github.com/okex/exchain/libs/tendermint/types/time"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/okex/exchain/libs/tendermint/types"
)

func TestP2PMsgProcess(t *testing.T) {
	// create private key
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
		t.Fatal(err)
	}
	vote.Signature = sig

	// prepare msg bytes
	msginfo := msgInfo{&VoteMessage{vote}, ""}
	msgBytes := cdc.MustMarshalBinaryBare(msginfo)

	// init wal
	wal := initWal()
	err = wal.Start()
	defer func() {
		wal.Stop()
		wal.Wait()
	}()

	buf := bytes.Buffer{}
	for i := 0; i < 10; i++ {
		buf.Write(msgBytes)
	}
	msg := conn.PacketMsg{ChannelID: 0x7f, EOF: 45, Bytes: buf.Bytes()}
	bz, err := cdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		t.Fatal(err)
	}

	err = parseMsg(bz, wal)
	if err != nil {
		t.Fatal(err)
	}
}

func initWal() *BaseWAL {
	walDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(walDir)

	walFile := filepath.Join(walDir, "wal")

	// this magic number 4K can truncate the content when RotateFile.
	// defaultHeadSizeLimit(10M) is hard to simulate.
	// this magic number 1 * time.Millisecond make RotateFile check frequently.
	// defaultGroupCheckDuration(5s) is hard to simulate.
	wal, err := NewWAL(walFile,
		autofile.GroupHeadSizeLimit(4096),
		autofile.GroupCheckDuration(1*time.Millisecond),
	)
	fmt.Println(err)
	return wal
}

func parseMsg(bz []byte, wal *BaseWAL) error {
	var packet conn.Packet
	var bufTmp = bytes.NewBuffer(bz)

	packet, _, err := conn.UnmarshalPacketFromAminoReader(bufTmp, int64(bufTmp.Len()))
	if err != nil {
		return err
	}

	switch pkt := packet.(type) {
	case conn.PacketMsg:
		msg, err := DecodeMsg(pkt.Bytes)
		if err != nil {
			return err
		}

		if err = msg.ValidateBasic(); err != nil {
			return err
		}

		wal.Write(msg)
	}
	return nil
}
