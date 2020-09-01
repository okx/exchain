package websocket

import "strings"

type SubscriptionTopic struct {
	Channel string
	Filter  string `default:""`
}

func (st *SubscriptionTopic) NeedLogin() bool {
	return st.Channel == DexSpotAccount || st.Channel == DexSpotOrder
}

func (st *SubscriptionTopic) ToString() (topic string, err error) {
	if len(st.Channel) == 0 {
		return "", errSubscribeParams
	}

	if len(st.Filter) > 0 {
		return st.Channel + ":" + st.Filter, nil
	} else {
		return st.Channel, nil
	}
}

func FormSubscriptionTopic(str string) *SubscriptionTopic {
	idx := strings.Index(str, ":")
	st := SubscriptionTopic{}
	if idx >= 0 {
		st.Channel = str[:idx]
		st.Filter = str[idx+1:]
	} else {
		st.Channel = str
		st.Filter = ""
	}

	return &st
}
