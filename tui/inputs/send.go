package inputs

func Send(message string) (value string, err error) {
	return SendPrefixed(message, "")
}
