package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

func TestParseABCILog(t *testing.T) {
	logs := `[{"log":"","msg_index":1,"success":true}]`

	res, err := ParseABCILogs(logs)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, res[0].Log, "")
	require.Equal(t, res[0].MsgIndex, uint16(1))
}

func TestABCIMessageLog(t *testing.T) {
	events := Events{NewEvent("transfer", NewAttribute("sender", "foo"))}
	msgLog := NewABCIMessageLog(0, "", events)

	msgLogs := ABCIMessageLogs{msgLog}
	bz, err := codec.Cdc.MarshalJSON(msgLogs)
	require.NoError(t, err)
	require.Equal(t, string(bz), msgLogs.String())
}

func TestABCIMessageLogJson(t *testing.T) {
	events := Events{NewEvent("transfer", NewAttribute("sender", "foo"))}
	msgLog := NewABCIMessageLog(0, "", events)

	tests := []ABCIMessageLogs{
		nil,
		{},
		{msgLog},
		{
			msgLog,
			NewABCIMessageLog(1000, "log", nil),
			NewABCIMessageLog(0, "log", Events{}),
			NewABCIMessageLog(1000, "",
				Events{
					Event{},
					NewEvent("", NewAttribute("", "")),
					NewEvent(""),
					NewEvent("type", NewAttribute("key", "value"), NewAttribute("", "")),
				}),
		},
	}

	for i, msgLogs := range tests {
		bz, err := codec.Cdc.MarshalJSON(msgLogs)
		require.NoError(t, err)

		nbz, err := msgLogs.MarshalToJson()
		require.NoError(t, err)
		require.EqualValues(t, bz, nbz)

		t.Log(fmt.Sprintf("%d passed", i))
	}
}

func BenchmarkABCIMessageLogJson(b *testing.B) {
	events := Events{NewEvent("transfer", NewAttribute("sender", "foo")), NewEvent("type", NewAttribute("key", "value"), NewAttribute("", ""))}
	msgLogs := ABCIMessageLogs{NewABCIMessageLog(1000, "test log", events)}
	b.ResetTimer()

	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = codec.Cdc.MarshalJSON(msgLogs)
		}
	})

	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = msgLogs.MarshalToJson()
		}
	})
}
