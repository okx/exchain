package tendermint

import (
	"bytes"
	"encoding/hex"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/libs/autofile"
	"github.com/okex/exchain/libs/tendermint/p2p/conn"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestP2PMsgProcess(t *testing.T) {
	// wal
	walDir,_ := os.Getwd()
	defer os.RemoveAll(walDir)
	walFile := filepath.Join(walDir, "wal")
	wal, err := consensus.NewWAL(walFile,
		autofile.GroupHeadSizeLimit(4096),
		autofile.GroupCheckDuration(1*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = wal.Start()
	defer func() {
		wal.Stop()
		wal.Wait()
	}()

	// prepare msg bytes
	msginfo := "b05b4f2c082210011abb01c1da2d940ab4010802100222480a20f6c1c76201852e755fb4ba22f07ed027c335c48e6d8f7269d4e7f028a53f96d71224080d12205350e84037be21d7ffeb6eb9e541d10a2c5e906a34da005e88fa07989e5b89fe2a0c0896d488940610b08ef1d1013214006e126d9b97ccd450361317d8caa4c5d24e84cf4240209ab24a1d1d6dc6fd81e306b36b365c050e6fb9a1c583613976e8ad4e931d7dbe6a1ad4c31e7611ca5896f69125df10503e9f6fdeaf52ca50e086e4b64d3401"
	msgBytes, err := hex.DecodeString(msginfo)
	if err != nil {
		t.Fatal(err)
	}

	err = parseMsg(msgBytes, wal)
	if err != nil {
		t.Fatal(err)
	}
}

func parseMsg(bz []byte, wal *consensus.BaseWAL) error {
	var packet conn.Packet
	var bufTmp = bytes.NewBuffer(bz)

	packet, _, err := conn.UnmarshalPacketFromAminoReader(bufTmp, int64(bufTmp.Len()))
	if err != nil {
		return err
	}

	switch pkt := packet.(type) {
	case conn.PacketMsg:
		msg, err := consensus.DecodeMsg(pkt.Bytes)
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
