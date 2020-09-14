package channels

func GetSpotAccountKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNEL,
		"account",
		args,
	)
}

func GetSpotOrderKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNEL,
		"order",
		args,
	)
}

func GetSpotDealKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNEL,
		"deal",
		args,
	)
}

func GetCSpotAccountKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNELC,
		"account",
		args,
	)
}

func GetCSpotOrderKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNELC,
		"order",
		args,
	)
}

func GetCSpotDealKey(args ...string) string {
	return GetSpotKey(
		PRIVATECHANNELC,
		"deal",
		args,
	)
}
