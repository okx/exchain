package channels

import "fmt"

func GetSpotDepthKey(productId string) string {
	return GetSpotKey(
		DEPTHCHANNEL,
		"depth",
		[]string{productId},
	)
}

func GetCSpotDepthKey(productId string) string {
	return fmt.Sprintf("%s:", GetSpotKey(
		DEPTHCHANNELC,
		"depth",
		[]string{productId},
	))
}
