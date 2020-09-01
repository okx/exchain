package websocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionTopic2Query(t *testing.T) {
	topic := SubscriptionTopic{
		Channel: "dex_spot/ticker",
		Filter:  "tbtc_tusdk",
	}

	rpcchannel, query := subscriptionTopic2Query(&topic)
	require.Equal(t, rpcchannel, "dex_spot/ticker:tbtc_tusdk")

	newTopic := query2SubscriptionTopic(query)

	require.Equal(t, newTopic.Channel, topic.Channel)
	require.Equal(t, newTopic.Filter, topic.Filter)
}
