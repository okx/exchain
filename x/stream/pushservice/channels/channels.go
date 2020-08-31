package channels

import (
	"fmt"
	"strings"
)

/*
	example:
		1. P3K_spot_instruments = "P3K:spot:instruments"
		2. P3A:futures:position:BTC-USD-170928:6810000:0
*/
const (
	PRODUCTS  = "P3K"
	PRODUCTSC = "P3KC"

	PUBLICCHANNEL  = "P3P"
	PUBLICCHANNELC = "P3C"

	PRIVATECHANNEL  = "P3A"
	PRIVATECHANNELC = "P3AC"

	DEPTHCHANNEL  = "P3D"
	DEPTHCHANNELC = "P3DC"
)

func getKey(channel, service, op string, args []string) string {
	if args == nil || len(args) == 0 {
		return fmt.Sprintf("%s:%s:%s", channel, service, op)
	}
	return fmt.Sprintf("%s:%s:%s:%s", channel, service, op, strings.Join(args, ":"))
}

func GetSpotKey(channel, op string, args []string) string {
	return getKey(channel, "dex_spot", op, args)
}
