package channels

import "fmt"

func GetSpotDepthKey(productID string) string {
	return GetSpotKey(
		DEPTHCHANNEL,
		"depth",
		[]string{productID},
	)
}

func GetCSpotDepthKey(productID string) string {
	return fmt.Sprintf("%s:", GetSpotKey(
		DEPTHCHANNELC,
		"depth",
		[]string{productID},
	))
}
