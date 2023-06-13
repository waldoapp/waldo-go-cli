package lib

func CompactMap[T any](slice []T, fn func(T) bool) []T {
	var result []T

	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}

	return result
}

func Map[T any, U any](slice []T, fn func(T) U) []U {
	var result []U

	for _, item := range slice {
		result = append(result, fn(item))
	}

	return result
}
