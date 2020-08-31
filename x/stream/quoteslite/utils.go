package quoteslite

import (
	"fmt"
	okex "github.com/okex/okchain/x/stream/quoteslite/okwebsocket"
	"strings"
)

func subscriptionTopic2Query(topic *okex.SubscriptionTopic) (channel, query string) {
	s, e := topic.ToString()
	if e == nil {
		query = fmt.Sprintf("tm.event='NewBlock' AND %s='%s'", rpcChannelKey, s)
	} else {
		query = fmt.Sprintf("tm.event='NewBlock' AND %s EXISTS", rpcChannelKey)
	}
	return s, query
}

func query2SubscriptionTopic(query string) *okex.SubscriptionTopic {

	subQuerys := strings.Split(query, "AND")
	if subQuerys != nil && len(subQuerys) == 2 {
		backendQuery := subQuerys[1]
		items := strings.Split(backendQuery, "=")
		if items != nil && len(items) == 2 {
			topicStr := strings.Replace(items[1], "'", "", -1)
			topic := okex.FormSubscriptionTopic(topicStr)
			return topic
		}
	}

	return nil
}
