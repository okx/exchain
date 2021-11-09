package websocket

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"strings"
)

func subscriptionTopic2Query(topic *SubscriptionTopic) (channel, query string) {
	s, e := topic.ToString()
	if e == nil {
		query = fmt.Sprintf("tm.event='NewBlock' AND %s='%s'", rpcChannelKey, s)
	} else {
		query = fmt.Sprintf("tm.event='NewBlock' AND %s EXISTS", rpcChannelKey)
	}
	return s, query
}

func query2SubscriptionTopic(query string) *SubscriptionTopic {
	subQuery := strings.Split(query, "AND")
	if len(subQuery) == 2 {
		backendQuery := subQuery[1]
		items := strings.Split(backendQuery, "=")
		if len(items) == 2 {
			topicStr := strings.ReplaceAll(items[1], "'", "")
			topic := FormSubscriptionTopic(topicStr)
			return topic
		}
	}

	return nil
}

func gzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return io.ReadAll(reader)
}
