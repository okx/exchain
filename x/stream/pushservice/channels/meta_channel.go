package channels

func GetSpotMetaKey() string {
	return GetSpotKey(
		PRODUCTS,
		"instruments",
		nil,
	)
}

func GetCSpotMetaKey() string {
	return GetSpotKey(
		PRODUCTSC,
		"instruments",
		nil,
	)
}
