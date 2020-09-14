package channels

import "fmt"

func GetSpotMatchKey(args ...string) string {
	return GetSpotKey(
		PUBLICCHANNEL,
		"matches",
		args,
	)
}

func GetCSpotMatchKey(args ...string) string {
	return fmt.Sprintf("%s:", GetSpotKey(
		PUBLICCHANNELC,
		"matches",
		args,
	))
}
