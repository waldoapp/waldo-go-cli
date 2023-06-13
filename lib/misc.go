package lib

func AppendIfNotEmpty(result *string, key, value, kvsep, rsep string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}

	if len(*result) > 0 {
		*result += rsep
	}

	*result += key + kvsep + value
}
