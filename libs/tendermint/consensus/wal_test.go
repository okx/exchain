package consensus

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	// "sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/tendermint/consensus/types"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	"github.com/okex/exchain/libs/tendermint/libs/autofile"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	tmtime "github.com/okex/exchain/libs/tendermint/types/time"
)

const (
	walTestFlushInterval = time.Duration(100) * time.Millisecond
)

func TestWALTruncate(t *testing.T) {
	walDir, err := ioutil.TempDir("", "wal")
	require.NoError(t, err)
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
	require.NoError(t, err)
	wal.SetLogger(log.TestingLogger())
	err = wal.Start()
	require.NoError(t, err)
	defer func() {
		wal.Stop()
		// wait for the wal to finish shutting down so we
		// can safely remove the directory
		wal.Wait()
	}()

	// 60 block's size nearly 70K, greater than group's headBuf size(4096 * 10),
	// when headBuf is full, truncate content will Flush to the file. at this
	// time, RotateFile is called, truncate content exist in each file.
	err = WALGenerateNBlocks(t, wal.Group(), 60)
	require.NoError(t, err)

	time.Sleep(1 * time.Millisecond) //wait groupCheckDuration, make sure RotateFile run

	wal.FlushAndSync()

	h := int64(50)
	gr, found, err := wal.SearchForEndHeight(h, &WALSearchOptions{})
	assert.NoError(t, err, "expected not to err on height %d", h)
	assert.True(t, found, "expected to find end height for %d", h)
	assert.NotNil(t, gr)
	defer gr.Close()

	dec := NewWALDecoder(gr)
	msg, err := dec.Decode()
	assert.NoError(t, err, "expected to decode a message")
	rs, ok := msg.Msg.(tmtypes.EventDataRoundState)
	assert.True(t, ok, "expected message of type EventDataRoundState")
	assert.Equal(t, rs.Height, h+1, "wrong height")
}

func TestWALEncoderDecoder(t *testing.T) {
	now := tmtime.Now()
	msgs := []TimedWALMessage{
		{Time: now, Msg: EndHeightMessage{0}},
		{Time: now, Msg: timeoutInfo{Duration: time.Second, Height: 1, Round: 1, Step: types.RoundStepPropose}},
	}

	b := new(bytes.Buffer)

	for _, msg := range msgs {
		msg := msg

		b.Reset()

		enc := NewWALEncoder(b)
		err := enc.Encode(&msg)
		require.NoError(t, err)

		c := new(bytes.Buffer)
		c.Reset()
		enc1 := NewWALEncoder(c)
		err = encodeOld(enc1, &msg)
		require.NoError(t, err)
		require.Equal(t, c.Bytes(), b.Bytes())

		dec := NewWALDecoder(b)
		decoded, err := dec.Decode()
		require.NoError(t, err)

		assert.Equal(t, msg.Time.UTC(), decoded.Time)
		assert.Equal(t, msg.Msg, decoded.Msg)
	}
}

func BenchmarkWALEncode(b *testing.B) {
	cdc.RegisterConcrete([]byte{}, "tendermint/wal/byte", nil)
	size := 1024 * 1024

	data := nBytes(1024 * 2)
	msgs := []TimedWALMessage{}
	for i := 0; i < size; i++ {
		msgs = append(msgs) //timeoutInfo{Duration: time.Second, Height: 1, Round: 1, Step: types.RoundStepPropose}
	}
	buffer := new(bytes.Buffer)
	msg := TimedWALMessage{Msg: data, Time: time.Now().Round(time.Second).UTC()}
	b.Run("Encode", func(b *testing.B) {
		b.ResetTimer()
		enc := NewWALEncoder(buffer)
		for i := 0; i < b.N; i++ {
			buffer.Reset()
			err := enc.Encode(&msg)
			if err != nil {
				b.Fatal(err)
			}
		}
		b.ReportAllocs()
	})
	b.Run("EncodeOld", func(b *testing.B) {
		b.ResetTimer()
		enc := NewWALEncoder(buffer)
		for i := 0; i < b.N; i++ {
			buffer.Reset()
			err := encodeOld(enc, &msg)
			if err != nil {
				b.Fatal(err)
			}
		}
		b.ReportAllocs()
	})
}

func encodeOld(enc *WALEncoder, v *TimedWALMessage) error {
	data := cdc.MustMarshalBinaryBare(v)

	crc := crc32.Checksum(data, crc32c)
	length := uint32(len(data))
	if length > maxMsgSizeBytes {
		return fmt.Errorf("msg is too big: %d bytes, max: %d bytes", length, maxMsgSizeBytes)
	}
	totalLength := 8 + int(length)

	msg := make([]byte, totalLength)
	binary.BigEndian.PutUint32(msg[0:4], crc)
	binary.BigEndian.PutUint32(msg[4:8], length)
	copy(msg[8:], data)

	_, err := enc.wr.Write(msg)
	return err
}

func TestWALEncode(t *testing.T) {
	cdc.RegisterConcrete([]byte{}, "tendermint/wal/byte", nil)
	size := 1024 * 1024

	data := nBytes(1024 * 2)
	msgs := []TimedWALMessage{}
	for i := 0; i < size; i++ {
		msgs = append(msgs, TimedWALMessage{Msg: data, Time: time.Now().Round(time.Second).UTC()}) //timeoutInfo{Duration: time.Second, Height: 1, Round: 1, Step: types.RoundStepPropose}
	}
	b := new(bytes.Buffer)
	enc := NewWALEncoder(b)

	//debug.SetGCPercent(-1)
	startM := getMemStats()
	for i := 0; i < size; i++ {
		b.Reset()
		err := enc.Encode(&msgs[i])
		require.NoError(t, err)
	}

	encodeM := getMemStats()
	t.Logf("After Encode   :GC<enable> MEMORY increase:%dMB,GC increase:%d", int64(encodeM.Alloc/1024/1024)-int64(startM.Alloc/1024/1024), encodeM.NumGC-startM.NumGC)
	for i := 0; i < size; i++ {
		b.Reset()
		err := encodeOld(enc, &msgs[i])
		require.NoError(t, err)
	}

	encodeOldM := getMemStats()
	t.Logf("After EncodeOld:GC<enable> MEMORY increase:%dMB,GC increase:%d", int64(encodeOldM.Alloc/1024/1024)-int64(encodeM.Alloc/1024/1024), encodeOldM.NumGC-encodeM.NumGC)

	//GC disable
	debug.SetGCPercent(-1)
	startM1 := getMemStats()
	for i := 0; i < size; i++ {
		b.Reset()
		err := enc.Encode(&msgs[i])
		require.NoError(t, err)
	}

	encodeM1 := getMemStats()
	t.Logf("After Encode   :GC<disable> MEMORY increase:%dMB,GC increase:%d", int64(encodeM1.Alloc/1024/1024)-int64(startM1.Alloc/1024/1024), encodeM1.NumGC-startM1.NumGC)
	for i := 0; i < size; i++ {
		b.Reset()
		err := encodeOld(enc, &msgs[i])
		require.NoError(t, err)
	}

	encodeOldM1 := getMemStats()
	t.Logf("After EncodeOld:GC<disable> MEMORY increase:%dMB,GC increase:%d", int64(encodeOldM1.Alloc/1024/1024)-int64(encodeM1.Alloc/1024/1024), encodeOldM1.NumGC-encodeM1.NumGC)

}

func getMemStats() (m runtime.MemStats) {
	runtime.ReadMemStats(&m)
	return m
}

func TestWALWrite(t *testing.T) {
	walDir, err := ioutil.TempDir("", "wal")
	require.NoError(t, err)
	defer os.RemoveAll(walDir)
	walFile := filepath.Join(walDir, "wal")

	wal, err := NewWAL(walFile)
	require.NoError(t, err)
	err = wal.Start()
	require.NoError(t, err)
	defer func() {
		wal.Stop()
		// wait for the wal to finish shutting down so we
		// can safely remove the directory
		wal.Wait()
	}()

	// 1) Write returns an error if msg is too big
	msg := &BlockPartMessage{
		Height: 1,
		Round:  1,
		Part: &tmtypes.Part{
			Index: 1,
			Bytes: make([]byte, 1),
			Proof: merkle.SimpleProof{
				Total:    1,
				Index:    1,
				LeafHash: make([]byte, maxMsgSizeBytes-30),
			},
		},
	}
	err = wal.Write(msg)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "msg is too big")
	}
}

func TestWALSearchForEndHeight(t *testing.T) {
	walBody, err := WALWithNBlocks(t, 6)
	if err != nil {
		t.Fatal(err)
	}
	walFile := tempWALWithData(walBody)

	wal, err := NewWAL(walFile)
	require.NoError(t, err)
	wal.SetLogger(log.TestingLogger())

	h := int64(3)
	gr, found, err := wal.SearchForEndHeight(h, &WALSearchOptions{})
	assert.NoError(t, err, "expected not to err on height %d", h)
	assert.True(t, found, "expected to find end height for %d", h)
	assert.NotNil(t, gr)
	defer gr.Close()

	dec := NewWALDecoder(gr)
	msg, err := dec.Decode()
	assert.NoError(t, err, "expected to decode a message")
	rs, ok := msg.Msg.(tmtypes.EventDataRoundState)
	assert.True(t, ok, "expected message of type EventDataRoundState")
	assert.Equal(t, rs.Height, h+1, "wrong height")
}

func TestWALPeriodicSync(t *testing.T) {
	walDir, err := ioutil.TempDir("", "wal")
	require.NoError(t, err)
	defer os.RemoveAll(walDir)

	walFile := filepath.Join(walDir, "wal")
	wal, err := NewWAL(walFile, autofile.GroupCheckDuration(1*time.Millisecond))
	require.NoError(t, err)

	wal.SetFlushInterval(walTestFlushInterval)
	wal.SetLogger(log.TestingLogger())

	// Generate some data
	err = WALGenerateNBlocks(t, wal.Group(), 5)
	require.NoError(t, err)

	// We should have data in the buffer now
	assert.NotZero(t, wal.Group().Buffered())

	require.NoError(t, wal.Start())
	defer func() {
		wal.Stop()
		wal.Wait()
	}()

	time.Sleep(walTestFlushInterval + (10 * time.Millisecond))

	// The data should have been flushed by the periodic sync
	assert.Zero(t, wal.Group().Buffered())

	h := int64(4)
	gr, found, err := wal.SearchForEndHeight(h, &WALSearchOptions{})
	assert.NoError(t, err, "expected not to err on height %d", h)
	assert.True(t, found, "expected to find end height for %d", h)
	assert.NotNil(t, gr)
	if gr != nil {
		gr.Close()
	}
}

/*
var initOnce sync.Once

func registerInterfacesOnce() {
	initOnce.Do(func() {
		var _ = wire.RegisterInterface(
			struct{ WALMessage }{},
			wire.ConcreteType{[]byte{}, 0x10},
		)
	})
}
*/

func nBytes(n int) []byte {
	buf := make([]byte, n)
	n, _ = rand.Read(buf)
	return buf[:n]
}

func benchmarkWalDecode(b *testing.B, n int) {
	// registerInterfacesOnce()

	buf := new(bytes.Buffer)
	enc := NewWALEncoder(buf)

	data := nBytes(n)
	enc.Encode(&TimedWALMessage{Msg: data, Time: time.Now().Round(time.Second).UTC()})

	encoded := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.Write(encoded)
		dec := NewWALDecoder(buf)
		if _, err := dec.Decode(); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
}

func BenchmarkWalDecode512B(b *testing.B) {
	benchmarkWalDecode(b, 512)
}

func BenchmarkWalDecode10KB(b *testing.B) {
	benchmarkWalDecode(b, 10*1024)
}
func BenchmarkWalDecode100KB(b *testing.B) {
	benchmarkWalDecode(b, 100*1024)
}
func BenchmarkWalDecode1MB(b *testing.B) {
	benchmarkWalDecode(b, 1024*1024)
}
func BenchmarkWalDecode10MB(b *testing.B) {
	benchmarkWalDecode(b, 10*1024*1024)
}
func BenchmarkWalDecode100MB(b *testing.B) {
	benchmarkWalDecode(b, 100*1024*1024)
}
func BenchmarkWalDecode1GB(b *testing.B) {
	benchmarkWalDecode(b, 1024*1024*1024)
}
