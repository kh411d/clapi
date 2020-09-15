package api

func getClaps() map[string]interface{} {
	return map[string]interface{}{
		"StatusCode": 200,
		"Body":       "Get clap",
	}
}

func addClap() map[string]interface{} {
	return map[string]interface{}{
		"StatusCode": 200,
		"Body":       "Adding clap",
	}
}
